# Notiflex 여정 기록

이 파일은 독자가 실제로 진행한 내용을 기록한다. AI가 각 챕터 완료 시 자동으로 업데이트한다.

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ✅ | 2026-04-30 | |
| ch2 | 2.3 gcloud 설정 | ✅ | 2026-04-30 | |
| ch2 | 2.4 GitHub 저장소 | ✅ | 2026-04-30 | |
| ch2 | 2.5 GKE 클러스터 | ✅ | 2026-04-30 | |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-04-30 | |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-04-30 | |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-04-30 | |
| ch3 | 3.3 기능 추가 | ✅ | 2026-04-30 | |
| ch3 | 3.4 CI | ✅ | 2026-04-30 | |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-04-30 | |
| ch4 | 4.2 메트릭 모니터링 | ✅ | 2026-04-30 | |
| ch4 | 4.3 로그 수집 | ✅ | 2026-04-30 | |
| ch4 | 4.4 알림 | ✅ | 2026-04-30 | |
| ch5 | 5.2 트래픽 관리 | ✅ | 2026-04-30 | |
| ch5 | 5.3 무중단 배포 | ✅ | 2026-04-30 | |
| ch5 | 5.4 ADR 기록 | ✅ | 2026-04-30 | |
| ch6 | 6.1 캐시 | ✅ | 2026-04-30 | |
| ch6 | 6.2 시크릿 관리 | ✅ | 2026-04-30 | |
| ch6 | 6.3 Canary 전환 | ✅ | 2026-04-30 | |
| ch6 | 6.4 아키텍처 스냅샷 | ✅ | 2026-04-30 | |
| ch7 | 7.2 멀티 노드풀 | ✅ | 2026-04-30 | |
| ch7 | 7.3 App of Apps | ✅ | 2026-04-30 | |
| ch7 | 7.4 멀티테넌시 | ✅ | 2026-04-30 | |
| ch8 | 8.1 메시징 | ✅ | 2026-04-30 | |
| ch8 | 8.2 트레이싱 | ✅ | 2026-04-30 | |
| ch8 | 8.3 CronJob | ✅ | 2026-04-30 | |
| ch9 | 9.1 저장소 분석 | ✅ | 2026-04-30 | |
| ch9 | 9.2 회고 | ✅ | 2026-04-30 | |
| ch9 | 9.3 온보딩 문서 | ✅ | 2026-04-30 | |
| ch9 | 9.4 GitAIOps 분석 | ✅ | 2026-04-30 | |
| ch9 | 9.5 마무리 | ✅ | 2026-04-30 | |

## 도구 선택 기록

독자가 3-프롬프트 패턴(탐색→비교→실행)에서 실제로 선택한 도구와 이유를 기록한다.

| 영역 | 선택 | 검토한 대안 | 선택 이유 |
|------|------|-----------|----------|
| GitOps (ch3.2) | ArgoCD v3.3.8 | Flux, Jenkins X | K8s 네이티브, Web UI, App of Apps 지원 |
| CI (ch3.4) | GitHub Actions | Jenkins, GitLab CI | 저장소 통합, WIF 지원, 무료 |
| 메트릭 (ch4.2) | Prometheus + Grafana | Datadog, New Relic | 오픈소스, kube-prometheus-stack 통합 |
| 로깅 (ch4.3) | Loki + Fluent Bit | ELK, Datadog Logs | Grafana 통합, 경량 인덱싱 |
| 알림 (ch4.4) | PrometheusRule | Grafana Alert | Prometheus 네이티브, PromQL 표현식 |
| 트래픽 관리 (ch5.2) | Gateway API (gke-l7-regional-external-managed) | Ingress, NGINX | GKE 네이티브, K8s 표준, HealthCheckPolicy |
| 배포 전략 (ch5.3) | Argo Rollouts Blue/Green | Flagger, Istio | ArgoCD 통합, 즉각 롤백, autoPromotion |
| 캐시 (ch6.1) | Valkey (Bitnami, standalone) | Redis, Memcached | Redis fork, BSD-3 라이선스, INCR 분산 ID |
| 시크릿 관리 (ch6.2) | GKE Secret Manager CSI + WI | K8s Secret, Vault | GKE 네이티브, SA 키 불필요, 파일 마운트 |
| 배포 전략 전환 (ch6.3) | Argo Rollouts Canary | Blue/Green 유지 | 트래픽 점진 이동, 운영 위험 최소화 |
| 노드 스케줄링 (ch7.2) | nodeSelector (cloud.google.com/gke-nodepool) | nodeAffinity, Taint/Toleration | GKE 자동 라벨, 단순 YAML, 역할별 노드풀 분리 |
| 멀티앱 관리 (ch7.3) | App of Apps (argocd/apps/ 디렉터리) | ApplicationSet, 개별 Application | 파일 추가만으로 앱 등록, Sync Wave 순서 보장 |
| 멀티테넌시 (ch7.4) | Namespace 분리 + per-tenant Rollout | 단일 namespace + 라벨 격리, vCluster | 강한 격리, ArgoCD App of Apps와 자연 결합, 테넌트별 독립 배포 |
| 배치 자동화 (ch8.3) | K8s CronJob | 외부 cron + 쿠버네티스 외부 트리거, Argo Workflows | 쿠버네티스 네이티브, ops-pool 배치, ArgoCD가 매니페스트로 관리 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | |
| Notiflex 이미지 | v0.3.1 | v0.1.0→v0.1.1→v0.2.0(Valkey)→v0.2.1(CSI)→v0.3.0(Kafka)→v0.3.1(OTel) |
| ArgoCD | v3.3.8 | |
| Kafka | 4.1.0 (Strimzi 1.0.0, KRaft) | |
| OTel SDK | - (Tempo 설치, SDK 적용) | |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium | 2 | Valkey, ArgoCD, monitoring |
| api-pool | e2-medium | 1 | notiflex-api (smb + enterprise) |
| worker-pool | e2-standard-2 | 1 | Kafka (ch8) |
| ops-pool | e2-small | 1 | Tempo, CronJob (ch8) |

## 트러블슈팅 이력

독자가 겪은 문제와 해결 방법을 기록한다. 같은 문제를 다시 겪지 않도록 한다.

| 챕터 | 문제 | 해결 |
|------|------|------|
| | | |

## GKE → k3s 마이그레이션 (2026-08-13)

ch9 완료 이후, GKE 원본을 개인 Synology DS920+ NAS(RAM 12GB)에서 운영하기 위해 k3s로 마이그레이션. 전체 계획과 Phase별 실측은 `docs/infraPlan.md` 참조. 노트북 k3d 클러스터에서 전 구간 검증 완료, NAS 이관(Phase 10~12)은 보류 중.

| Phase | 내용 | 상태 |
|------|------|------|
| 0~1 | 노트북 환경 준비, k3d 클러스터 구성 | ✅ |
| 2 | Gateway API → Traefik 전환 | ✅ |
| 3 | GKE Secret Manager CSI → 순수 K8s Secret | ✅ |
| 4 | Artifact Registry → GHCR CI/CD 전환 | ✅ |
| 5 | GKE nodeSelector 전부 제거 | ✅ |
| 6 | Kafka PVC storageClassName(local-path) 명시 | ✅ |
| 7 | 리소스 예산 실측 (kube-prometheus-stack/Loki/Tempo/Fluent Bit) | ✅ |
| 8 | Kafka(Strimzi) combined 모드 배치 | ✅ |
| 9 | ArgoCD + Argo Rollouts 설치·GitOps 검증 | ✅ |
| 9.5 | Valkey(Bitnami) 배포, 앱 end-to-end 테스트 | ✅ |
| 10~12 | NAS(Synology DS920+)에 k3s 설치 및 컷오버 | ⏸ 보류 |
| 13 | 문서 업데이트(CLAUDE.md/architecture.md/ONBOARDING.md/ADR) | ✅ |

**마이그레이션 중 발견한 이슈**(GKE와 무관한 순수 버전 드리프트/환경 이슈):
- Traefik의 `web` entryPoint가 내부 8000번 포트 사용 → Gateway `Listener.port`를 80→8000으로 조정
- 최신 Strimzi 오퍼레이터가 Kafka 4.1.0 미지원 → 4.3.0으로 상향
- k3d 노드(컨테이너)에 `/etc/machine-id` 파일이 없어 Fluent Bit DaemonSet 마운트 실패 → 노트북 전용 volumes 오버라이드로 우회(NAS에서는 미발생 예상)
- GHCR private 이미지 pull 401 → `imagePullSecrets`로 해결
- Valkey는 원본 저장소에도 매니페스트가 없었음(ADR-008에 따라 Bitnami Helm으로 별도 설치) → `helm-values/valkey.yaml` 신규 추가

## 마이그레이션 이후 후속 작업 (2026-08-13)

k3s 마이그레이션 검증이 끝난 뒤 실제 운영 관점에서 발견된 이슈 정리와, 이 파이프라인 위에서 진행한 첫 기능 추가.

### enterprise 네임스페이스 정상화
- `enterprise/rollout.yaml`이 참조하는 `ghcr.io/smjjang-dev/notiflex/api:v0.3.1` 태그가 실제로는 한 번도 push된 적이 없어 `ImagePullBackOff` 상태였음 — CI는 `k8s/smb/rollout.yaml`만 자동 패치하도록 되어 있어 enterprise는 처음부터 수동 승격 대상이었기 때문. `app/`을 직접 빌드해 해당 태그로 GHCR push하여 해결.
- `imagePullSecrets` 추가 이전에 생성됐던 구버전 ReplicaSet이 잔존해 재발 → 삭제로 정리 (smb 때와 동일 패턴).
- **추가 발견**: enterprise 전용 `valkey` Secret이 원본 저장소의 예전 비밀번호(`bTdMRTR1dVUyWQ==`)를 그대로 갖고 있었는데, 실제 접속 대상인 공유 Valkey 인스턴스(notiflex 네임스페이스)는 새 비밀번호로 떠 있어 `WRONGPASS` 에러 발생. `k8s/enterprise/valkey-secret.yaml`을 동일 값으로 수정 후 커밋 — **kubectl로 Secret을 직접 patch했더니 ArgoCD selfHeal이 git의 예전 값으로 즉시 되돌리는 것도 함께 확인**(GitOps 환경에서 out-of-band 변경이 무의미하다는 걸 실증).

### `/version` 엔드포인트 추가 — 첫 실전 CI/CD 사이클
- `app/main.go`에 `/version` 핸들러 추가(`version`/`commit`/`buildTime`/`goVersion` 반환), 값은 하드코딩이 아니라 `Dockerfile`의 `-ldflags`로 빌드 시점에 주입.
- `.github/workflows/ci.yaml`에 `--build-arg VERSION/COMMIT/BUILD_TIME` 추가.
- 코드 변경 → `git push` → GitHub Actions가 46초 만에 이미지 빌드/GHCR push + `rollout.yaml` 자동 패치·재커밋 → ArgoCD가 새 커밋 감지 후 자동 동기화 → Argo Rollouts가 카나리(20%→50%→80%→100%) 단계를 실제로 밟으며 무중단 전환 — **마이그레이션 이후 처음으로 전체 GitOps 파이프라인을 처음부터 끝까지 실전 검증**.
- 배포 완료 후 `curl http://localhost/version` → `{"version":"sha-29310c8","commit":"29310c84f866a372f0020da1dc4d2ccc0dd4f898","buildTime":"2026-08-13T05:34:08Z","goVersion":"go1.25.12"}` 확인.
- enterprise는 이 배포 대상에서 제외되어 여전히 `/version` 없는 `v0.3.1` 사용 중.

### Grafana를 Gateway로 상시 노출
- 기존 Prometheus/Grafana는 `kubectl port-forward`로만 접근 가능했음 — `/grafana` 경로로 항상 접속되도록 Gateway/HTTPRoute 추가.
- `k8s/smb/gateway.yaml`의 `allowedRoutes.namespaces.from`을 `Same`→`All`로 넓혀 다른 네임스페이스(monitoring)의 HTTPRoute도 이 Gateway에 붙을 수 있게 함.
- `k8s/monitoring/grafana-httproute.yaml` 신규 추가 — `/grafana` prefix를 `kube-prometheus-grafana:80`으로 라우팅.
- Grafana를 서브패스로 정상 서빙하려면 `server.root_url`/`serve_from_sub_path` 설정이 필요해 `helm-values/kube-prometheus.yaml`에 `grafana.ini` 섹션 추가.
- **재확인된 패턴**: `k8s/smb/`는 ArgoCD(`notiflex-smb`)가 관리 중이라, git에 커밋하지 않고 `kubectl apply`만 하면 selfHeal이 즉시 되돌림 — 반드시 git 커밋·push가 먼저.

### Prometheus도 서브패스로 노출
- `helm-values/kube-prometheus.yaml`에 `prometheus.prometheusSpec.routePrefix`/`externalUrl` 추가, `k8s/monitoring/prometheus-httproute.yaml` 신규로 `/prometheus` 라우팅.
- **사이드 이펙트 발견**: Grafana의 `serve_from_sub_path`가 브라우저 요청뿐 아니라 **내부 스크레이핑 요청까지** `root_url`(`http://localhost/grafana/...`)로 301 리다이렉트시켜서, Prometheus가 ClusterIP로 직접 찌르던 `/metrics` 스크레이핑이 전부 실패(`connect: connection refused`)하는 부작용 발생.
- `grafana.serviceMonitor.path`를 실제 리다이렉트 목적지인 `/grafana/metrics`로 맞춰서 해결. Prometheus 타겟 13개 전부 `up`, `count(up)` 쿼리로 실제 수집 확인.

### Grafana에 Loki 데이터소스 연결 (로그를 한 곳에서 보기)
- Grafana API로 확인해보니 `Prometheus`/`Alertmanager`만 등록되어 있고 **Loki 데이터소스가 없었음** — `helm-values/loki.yaml`의 `grafana.datasource` 설정은 현재 Loki 차트 버전에서 더 이상 지원하지 않는 죽은 설정(`helm show values grafana/loki`로 확인, 최신 차트엔 해당 키 자체가 없음).
- `helm-values/kube-prometheus.yaml`에 `grafana.additionalDataSources`로 Loki(`http://loki-gateway.monitoring.svc.cluster.local`)를 직접 등록해서 해결.
- 검증: Loki 라벨 조회로 로그가 들어오는 네임스페이스 확인(`argo-rollouts, argocd, enterprise, kafka, kube-system, monitoring, notiflex` 7개 전부), `{namespace="notiflex"}` 쿼리로 실제 healthcheck 로그 라인까지 확인.
- 접속: `/grafana/` 로그인 후 좌측 **Explore** 메뉴에서 데이터소스를 `Loki`로 선택 → LogQL로 조회(예: `{namespace="notiflex"}`).
- **교훈**: 서브패스 노출(`serve_from_sub_path`)은 외부 접근뿐 아니라 같은 클러스터 내부의 서비스 간 스크레이핑/헬스체크 경로에도 영향을 줄 수 있다 — 도입 시 내부 클라이언트(Prometheus 등)의 접근 경로도 같이 점검 필요.
