# Infra Migration Plan — GKE → k3s (노트북 검증 → Synology DS920+ 이관)

- 작성일: 2026-08-13
- 배경: 이 저장소는 원래 GKE(Google Kubernetes Engine) 위에 구축되어 있음. 이를 개인 소유 Synology DS920+ NAS(Docker/Container Manager)로 이전하고자 함.
- 결정된 방향: **B안 — NAS에 k3s(경량 K8s)를 올려 기존 K8s 매니페스트(Rollout, ArgoCD, Strimzi 등)를 최대한 재사용**. 순수 `docker-compose` 전면 재설계는 채택하지 않음.
- NAS 사양: DS920+, RAM 4GB(기본) + 8GB(증설) = 총 12GB. DSM 자체 오버헤드(~1.5~2GB) 제외 시 k3s 클러스터 실사용 가능 예산은 약 **9~10GB**로 가정.
- 진행 전략: 리스크가 큰 매니페스트 재작성/CI 전환 작업은 먼저 **노트북 Docker(k3s/k3d)** 에서 검증하고, 검증이 끝난 뒤 동일 Git 상태를 NAS의 새 k3s 클러스터에 GitOps로 재구성한다 (Git이 단일 진실 소스이므로 클러스터 간 수동 이관 최소화).

---

## 1. 현재 상태 — GCP/GKE 종속 지점 인벤토리

| 영역 | 파일 | GCP 종속 내용 |
|---|---|---|
| 인그레스/LB | `k8s/smb/gateway.yaml`, `k8s/smb/healthcheckpolicy.yaml` | `gatewayClassName: gke-l7-regional-external-managed` (GKE 관리형 L7 LB), `HealthCheckPolicy`는 `networking.gke.io/v1` 전용 CRD |
| 시크릿 | `k8s/smb/secret-provider.yaml`, `k8s/enterprise/secret-provider.yaml` | `provider: gke` + `secrets-store-gke.csi.k8s.io` CSI로 Google Secret Manager 마운트, Workload Identity 필요 |
| 이미지 레지스트리 | `k8s/smb/rollout.yaml`, `k8s/enterprise/rollout.yaml` | `image: asia-northeast3-docker.pkg.dev/.../notiflex/api:sha-...` (Artifact Registry) |
| 노드 스케줄링 | rollout.yaml ×2, `healthcheck-cronjob.yaml`, `helm-values/strimzi.yaml`, `helm-values/tempo.yaml` | `nodeSelector: cloud.google.com/gke-nodepool: {api-pool\|worker-pool\|ops-pool}` |
| 스토리지 | `k8s/kafka/kafka-cluster.yaml` | PVC storageClassName 미지정 → GKE 기본값(`standard-rwo`, GCE PD)에 암묵적 의존 |
| CI/CD | `.github/workflows/ci.yaml` | `google-github-actions/auth`(WIF) → Artifact Registry push → rollout.yaml sed 패치 후 git push(GitOps 트리거) |
| 문서 | `docs/architecture-decisions.md`(ADR-006/009/011), `CLAUDE.md`, `claude-context/architecture.md`, `ONBOARDING.md` | GCP 프로젝트 ID, 클러스터명, 리전, Gateway 고정 IP 등 |

종속성이 없는 부분: `app/`(Go 소스, scratch 기반 Dockerfile — 완전 이식 가능), `helm-values/loki.yaml`(이미 filesystem 스토리지), Prometheus/Grafana/Fluent Bit 값 파일, ArgoCD Application 정의 자체(K8s API만 있으면 동작).

---

## 2. Stage 1 — 노트북에서 검증

### Phase 0. 사전 준비
- 노트북 Docker 환경 확인(Docker Desktop/WSL2), k3s vs k3d 결정 (권장: **k3d** — Docker Desktop 위에서 바로 구동 가능)
- GHCR용 GitHub PAT / Actions 설정 준비

### Phase 1. 노트북 k3s 클러스터 구성
- k3d(or k3s)로 단일 노드 클러스터 기동, `kubectl get nodes`로 확인

### Phase 2. Gateway API → Traefik 전환
- Gateway API CRD 설치, Traefik의 experimental Gateway provider 활성화
- `gatewayClassName`을 `gke-l7-regional-external-managed` → `traefik`으로 교체
- `healthcheckpolicy.yaml`(GKE 전용 CRD) 삭제, 대신 컨테이너에 표준 `readinessProbe`/`livenessProbe`(`/health:8080`) 추가

### Phase 3. 시크릿 관리 전환
- `secret-provider.yaml`(SecretProviderClass, GKE) 삭제
- `valkey-secret.yaml` 패턴(순수 K8s Secret)을 smb/enterprise 양쪽에 통일 적용
- rollout.yaml의 CSI 볼륨 마운트 → 일반 `secret` 볼륨으로 교체(마운트 경로 동일 유지 → 앱 코드 변경 불필요)
- 홈랩 단일 사용자 환경이므로 Vault 등 별도 시크릿 매니저는 과설계로 판단, k3s `--secrets-encryption` 옵션으로 최소 방어선 확보(NAS 설치 시 적용)

### Phase 4. GHCR 기반 CI/CD 전환
- `google-github-actions/auth` 스텝 제거, `docker/login-action`으로 `ghcr.io` 로그인(`GITHUB_TOKEN` 사용, 별도 WIF 불필요)
- 이미지 태그: `ghcr.io/<github-id>/notiflex/api:sha-...`
- 기존 "sed로 rollout.yaml 패치 후 git push" GitOps 패턴은 그대로 유지(레지스트리만 교체)

### Phase 5. GKE nodeSelector 제거
- `rollout.yaml`(api-pool) ×2, `healthcheck-cronjob.yaml`(ops-pool), `kafka-cluster.yaml`(worker-pool), `helm-values/strimzi.yaml`, `helm-values/tempo.yaml`에서 `cloud.google.com/gke-nodepool` 셀렉터 전부 제거 (단일 노드이므로 불필요)

### Phase 6. 스토리지 전환 확인
- Kafka PVC는 k3s 기본 StorageClass(`local-path`)로 자동 동작하지만, `storageClassName: local-path`를 명시해 의도 명확화

### Phase 7. 리소스 예산 적용 및 실측 검증

**노트북 k3d 클러스터 실측치 (2026-08-13, Kafka 제외 상태)**

| 컴포넌트 | 요청 예상치 | 실측 메모리 | 비고 |
|---|---|---|---|
| k3s 시스템(Traefik/CoreDNS/local-path/metrics-server/svclb) | ~1.0Gi | ~143Mi | 예상보다 훨씬 가벼움 |
| ArgoCD(core만: server/repo-server/app-controller/redis) | ~800Mi | ~242Mi | core-only 트리밍 효과 큼 |
| Argo Rollouts controller | ~100Mi | ~25Mi | |
| Prometheus(retention 2일로 축소) | ~800Mi~1Gi | ~506Mi | |
| Grafana | ~200Mi | ~317Mi | 예상보다 다소 높음 |
| Alertmanager | ~50Mi | ~32Mi | |
| Loki(singlebinary, filesystem, gateway+canary 포함) | ~300Mi | ~170Mi | |
| Fluent Bit(1노드) | ~60Mi | ~6Mi | |
| Tempo | ~200Mi | ~29Mi | 트래픽 없는 유휴 상태라 낮음 |
| kube-state-metrics + node-exporter | (미산정) | ~36Mi | |
| **파드 합계(Kafka·앱 제외)** | ~3.5Gi | **~1.5Gi** | |
| **노드 전체 사용량**(`kubectl top node`) | - | **~3.95Gi / 7.7Gi (51%)** | 파드 합계와 약 2.4Gi 차이 = containerd/kubelet/Docker Desktop WSL2 VM 기반 오버헤드 |

**중요 발견**: 파드별 사용량 합계(~1.5Gi)와 노드 전체 사용량(~3.95Gi) 사이에 상당한 차이가 있음 — k3s 런타임 자체(containerd, kubelet) + 노트북은 Docker Desktop의 WSL2 가상머신 기반 오버헤드가 추가로 존재하기 때문. NAS는 Docker Desktop 가상화 계층 없이 k3s가 리눅스에 직접 설치되므로 이 오버헤드는 노트북보다 작을 것으로 예상되나, **정확한 값은 Phase 10 NAS 설치 후 재실측 필요**. 워크로드(파드) 자체의 실측 메모리는 예산표보다 전반적으로 여유 있음이 확인됨 — 이전 예상치(소계 ~3.8Gi)는 보수적으로 잡았던 것으로 판단.

### Phase 8. Kafka(Strimzi) combined 모드 배치 — 실측 결과 포함
- KRaft 모드, broker+controller combined 단일 파드(`KafkaNodePool.spec.roles: [controller, broker]`, replicas: 1) — 원본부터 이미 combined 구성이라 별도 축소 작업 불필요, 확인만 진행
- JVM 힙 `jvmOptions: {"-Xms": "512m", "-Xmx": "512m"}` 명시적 제한 추가 (`k8s/kafka/kafka-cluster.yaml`)
- **문제 + 해결**: 최신 Strimzi 오퍼레이터(1.1.0, helm 최신)가 원본의 `kafka.spec.kafka.version: 4.1.0`을 더 이상 지원하지 않음(지원 버전: 4.2.0/4.2.1/4.3.0) — 마이그레이션과 무관한 순수 버전 드리프트. `4.3.0`으로 상향 후 정상 기동
- 실측: `notiflex-kafka-controller-0` 387Mi, `notiflex-kafka-entity-operator` 150Mi, `strimzi-cluster-operator` 233Mi → 합계 **~770Mi**, 예산(~800Mi~1Gi)과 거의 일치
- `KafkaTopic/notifications` 정상 생성 확인(`READY: True`)
- 클러스터 전체 누적 사용량(Kafka 포함): 노드 전체 **~5.26Gi / 7.7Gi (68%)**

| 컴포넌트 | 요청 예상치 | 실측 메모리 |
|---|---|---|
| Strimzi 오퍼레이터 | ~256Mi | ~233Mi |
| Kafka(KRaft combined, heap 512Mi) + entity-operator | ~800Mi~1Gi | ~537Mi |
| **총합(전체 스택, Kafka 포함)** | ~5~5.5Gi | **~5.26Gi** |

### Phase 9. ArgoCD + Argo Rollouts 설치 및 검증(노트북)
- `argocd/root-app.yaml`의 `destination.server: https://kubernetes.default.svc`(in-cluster)는 변경 불필요
- `repoURL`(GitHub)도 그대로, poll/webhook 방식으로 동기화
- 카나리 전략(20→50→80→100%, 30초 pause) 동작 확인

### Phase 9.5. 앱 end-to-end 검증 (추가 작업)
- Valkey가 원본 저장소에도 매니페스트로 커밋되어 있지 않아(ADR-008: Bitnami Helm 차트) 앱이 CrashLoop 상태였음 — `helm-values/valkey.yaml` 신규 작성 후 `helm install valkey bitnami/valkey`로 배포(release replication 모드로 `valkey-primary` 서비스명이 앱 기대값과 일치)
- `curl http://localhost/health` → `{"status":"ok","version":"v0.3.1"}`, `curl http://localhost/id` → Valkey INCR 기반 ID 정상 증가 확인
- Kafka 이벤트 발행/수신 로그, `notiflex-healthcheck` CronJob 정상 통과까지 확인 — Traefik→HTTPRoute→Service→Rollout→Valkey/Kafka 전체 체인 검증 완료

---

## 3. Stage 2 — NAS(DS920+)로 이관 (보류 중, 2026-08-13 기준)

> 사용자 요청으로 NAS 실배포는 보류. 아래 Phase 10~12는 계획만 확정된 상태이며 재개 시 SSH 접속 정보가 필요하다.

### Phase 10. NAS에 k3s 설치
- Phase 7 실측 데이터를 근거로 예산표 최종 조정
- SSH로 `curl -sfL https://get.k3s.io | sh -s - server --write-kubeconfig-mode 644 --secrets-encryption --data-dir /volume1/docker/k3s`
- local-path-provisioner 데이터 경로도 `/volume1` 하위로 지정(시스템 파티션 보호)
- DSM이 80/443 포트를 이미 쓰는 경우 충돌 처리(포트 변경 또는 별도 포워딩)

### Phase 11. NAS에 ArgoCD 재설치 및 GitOps 동기화
- NAS 클러스터에 ArgoCD 설치 → 동일 Git 저장소 연결 → 검증된 최신 매니페스트 자동 pull
- Secret 실값은 Git에 없으므로 NAS에서 별도로 `kubectl create secret` 재생성 필요
- Kafka 토픽 데이터, 수동으로 만든 Grafana 대시보드, Prometheus/Loki 히스토리 등 Git에 없는 상태는 이관되지 않음 — 필요 시 별도 백업/복원

### Phase 12. 검증 및 컷오버
- NAS 클러스터 단독 정상 동작 확인 → DNS/포트포워딩을 NAS로 전환
- 안정화 후 노트북 클러스터 폐기, GKE 리소스 정리(포워딩 룰/디스크 등 쿼터 누수 항목 포함)

### Phase 13. 문서 업데이트
- `CLAUDE.md`, `claude-context/architecture.md`, `ONBOARDING.md`, `docs/architecture-decisions.md`(ADR-006/009/011)를 NAS/k3s 기준으로 갱신

---

## 4. 진행 순서 요약

`Phase 0 → 1 → 2 → 3 → 5 → 6 → 9(앱 먼저 뜨게) → 7(리소스 검증) → 4(CI/CD) → 8(Kafka) → 10 → 11 → 12 → 13`
