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

### Grafana에 Tempo 데이터소스 연결 (트레이스도 한 곳에서 보기)
- Loki와 동일한 이유로 Tempo도 데이터소스 미등록 상태였음 — `helm-values/kube-prometheus.yaml`의 `grafana.additionalDataSources`에 Tempo(`http://tempo.monitoring.svc.cluster.local:3200`) 추가로 해결.
- 검증: `/health`, `/id`에 실제 요청을 보낸 뒤 Tempo `/api/search`로 조회 → 방금 보낸 요청들이 트레이스로 잡히는 것 확인. 트레이스 상세(`/api/traces/{id}`)에서 span 이름·소요시간까지 정상 확인.
- **참고(사소한 개선 여지)**: `app/main.go`가 OTel Resource에 `service.name`을 명시적으로 설정하지 않아, Tempo에 `rootServiceName`이 `unknown_service:notiflex-api`로 잡힘 — 기능상 문제는 없으나 Tempo UI에서 서비스명으로 검색/필터링할 때 보기 불편할 수 있음.

### `service.name` 리소스 속성 추가
- `initTracer()`에서 `resource.Merge(..., resource.NewSchemaless(attribute.String("service.name", "notiflex-api")))`로 TracerProvider에 명시적 서비스명 부여.
- semconv 패키지(버전 고정 이슈 있음) 대신 `attribute.String` 직접 사용으로 의존성 추가 없이 처리.
- **1차 시도 실패 + 원인**: `resource.Merge(resource.Default(), myResource)` 순서로 넣었더니 여전히 `unknown_service:notiflex-api`로 나옴 — `resource.Merge(a, b)`는 키 충돌 시 **a의 값이 우선**하는데, `resource.Default()`가 이미 자체 `unknown_service:...` 서비스명을 갖고 있어서 내가 지정한 값이 조용히 버려지고 있었음(에러 리턴 없이). 인자 순서를 `resource.Merge(myResource, resource.Default())`로 바꿔서 해결.
- CI/CD 배포 후 Tempo `rootServiceName`이 `notiflex-api`로 정상 표기됨을 재확인.

### 알림(Alerting) 파이프라인 점검
- "문제 생기면 자동 알림 받을 수 있나" 질문 계기로 점검 — 두 가지 미비점 발견:
  1. `k8s/monitoring/pod-restart-alert.yaml`(PrometheusRule)이 저장소에만 있고 **클러스터에 한 번도 적용된 적이 없었음** — `kubectl apply`로 적용.
  2. Alertmanager의 receiver가 기본값 `"null"`(no-op)만 있어서, 규칙이 발동해도 실제로는 아무 데도 알림이 안 가는 상태였음.
- 적용 후 검증: Prometheus가 규칙 로드(`health: ok`, `state: inactive`=정상), Prometheus↔Alertmanager 연결(`activeAlertmanagers`) 확인 — **감지 파이프라인은 끝까지 정상 동작**, 마지막 "어디로 보낼지"만 비어있는 상태.
- 실제 알림 채널은 **Discord webhook**으로 정하고 연결은 보류했으나, 이후 **Slack webhook**으로 최종 결정하고 실제 연결 완료.

### Slack 알림 연결
- webhook URL은 git에 커밋하지 않고 `slack-webhook` Secret(`monitoring` 네임스페이스)으로 보관, Alertmanager 파드에 `/etc/alertmanager/secrets/slack-webhook/`로 마운트 → `global.slack_api_url_file`로 참조.
- `route.receiver`를 `null`→`slack`으로 변경(Watchdog은 계속 null 유지), `receivers`에 `slack_configs` 추가.
- 검증: Alertmanager API로 테스트 알림(`TestAlert`) 직접 POST → `active` 상태로 접수 확인, Slack 채널에 실제 도착 확인.
- **부수 발견**: 활성 알림 목록에 `KubeControllerManagerDown`/`KubeSchedulerDown`/`KubeProxyDown`이 이미 떠 있었음 — k3s는 controller-manager/scheduler/kube-proxy를 별도 파드가 아니라 `k3s server` 프로세스 안에 내장하므로 `absent(up{job="..."})` 기반 기본 규칙이 **k3s에서는 구조적으로 항상 오탐**(NAS에서도 동일하게 발생할 것). Slack 라우팅에서 이 3개만 `null`로 mute 처리.
- `NodeClockNotSynchronising`도 active였는데 이건 진짜 — 실측 `node_timex_sync_status: 0`, 오차 16초. **Docker Desktop WSL2 VM 시계 드리프트**로 추정(노트북 환경 한정, NAS에서는 미발생 예상) — mute하지 않고 그대로 둠.
- **트러블슈팅**: helm upgrade가 백그라운드 실행 중 타임아웃으로 끊기면서 릴리스가 `pending-upgrade`에 걸려 다음 upgrade가 `another operation in progress`로 실패 — `helm rollback kube-prometheus <이전 정상 리비전>`으로 `deployed` 상태 복구 후 재적용.
- Slack 테스트 알림 수신 확인(사용자 확인 완료).

### 알림 규칙 전수 점검 + `KubeJobFailed` 원인 규명
- Prometheus에 로드된 alerting 규칙 전수 확인 — 총 **156개**(kube-prometheus-stack 기본 제공 155개 + 커스텀 `PodRestartTooMany` 1개), 카테고리: kubernetes-apps/node-exporter/prometheus/etcd/kubelet/resources/storage/apiserver 등.
- `KubeJobFailed`가 `pending` 상태인 것을 확인해 원인 조사: `notiflex-healthcheck-297765{35,40,45}` 3개 Job이 `Failed` 상태로 남아있었음 — 전부 새벽 3시대(Valkey 비밀번호 불일치로 앱이 응답 못하던 시점)에 실패한 **이미 해결된 과거 문제의 잔여 기록**. CronJob의 `failedJobsHistoryLimit: 3`이 이걸 계속 보관하고 있었음.
- 3개 Job 삭제 → `kube_job_failed` 메트릭 즉시 사라짐 → 다음 평가 주기에 `KubeJobFailed`가 `inactive`로 정상 해소되는 것까지 확인.

## 시크릿 관리 정리 + Canary 자동 롤백 도입 (2026-08-13)

### 시크릿 유출 발견 및 정리
- `k8s/smb/valkey-secret.yaml`, `k8s/enterprise/valkey-secret.yaml`이 base64로 "인코딩"된 Valkey 비밀번호를 담은 채 git에 커밋되어 있었음(base64는 암호화가 아니라 평문이나 다름없음) — 이미 GitHub(private repo)에 두 커밋으로 push된 상태였음.
- 대응: ① 비밀번호 로테이션(클러스터에 즉시 반영, Valkey/앱 파드 재기동으로 검증) ② 두 파일 git 추적 해제 + `.gitignore`에 `*-secret.yaml` 추가 ③ `git filter-branch`로 히스토리에서 완전 제거 후 **origin에 force-push** ④ ONBOARDING.md에 imperative 재생성 명령 문서화(GHCR pull secret과 동일 패턴).
- **재발 방지 장치 추가**: `.githooks/pre-commit`(populated `kind: Secret` 매니페스트, 프라이빗 키, AWS/GitHub/Slack 토큰 패턴을 커밋 직전에 차단) + `.github/workflows/secret-scan.yaml`(기존 `ci.yaml`이 `paths: [app/**]`라 k8s 매니페스트 변경엔 반응 안 하는 사각지대를 메우는 경로 무제한 백스톱). 다른 클론에서 fresh clone → `git config core.hooksPath .githooks` 활성화까지 실제로 재현해 검증, CI 백스톱도 가짜 시크릿을 실제로 push해서 빨간 X로 실패하는 것까지 확인 후 되돌림.

### Canary 배포 에러율 기반 자동 롤백
- 기존 canary(`strategy.canary`, `setWeight`/`pause` 고정 스텝)는 실제 헬스/에러율과 무관하게 타이머로만 진행되고, `/health`가 의존성과 무관한 얕은 체크라 로직 버그가 있어도 그대로 100% 승격될 수 있는 구조였음(ADR-007 문서상 "Blue/Green"이라 적혀 있었지만 실제 매니페스트는 canary — 문서-구현 불일치도 함께 발견).
- 추가한 것: `app/main.go`에 `http_requests_total{path,code}` 카운터 + `/metrics`(prometheus/client_golang) → Service에 라벨/named port 부여 → `ServiceMonitor`(smb/enterprise 각각)로 kube-prometheus가 스크레이핑 → `ClusterAnalysisTemplate`(canary preview 서비스의 5xx 비율 쿼리, `failureLimit: 2`) → 양쪽 Rollout에 `strategy.canary.analysis`로 background analysis 연결(`startingStep: 1`).
- **해피패스 실증**: 실제 git push → CI 빌드/푸시/매니페스트 패치 → ArgoCD sync → 진짜 canary 롤아웃 발생 → `AnalysisRun`이 Prometheus에서 실제 에러율(0)을 측정하며 `Successful` → 100% 승격.
- **네거티브 테스트(자동 abort 실증)**: `/id`를 일부러 항상 500 반환하도록 수정해 배포. 1차 시도는 상태 확인하느라 시간을 끄는 사이 90초 canary 창이 끝나버려 에러 트래픽이 늦게 도착 — **버그 이미지가 그대로 100% 승격되어 실제로 프로덕션에 노출되는 사고**로 이어짐(`/id` 500 확인) → 즉시 이전 커밋으로 되돌려 복구. 2차 시도는 canary가 뜨기 전에 백그라운드 부하생성기를 미리 띄워두고 재시도 → 측정값이 0.40, 0.63, 0.93으로 치솟으며 `failureLimit: 2` 충족 → Rollout이 `phase: Degraded, abort: true`로 **자동 중단**, 두 파드 모두 마지막 정상 이미지로 유지된 채 버그 이미지는 끝까지 승격되지 않음을 확인. 이후 코드는 정상으로 되돌려 재배포.
- **부수 발견(ArgoCD prune 함정)**: 시크릿 정리 때 git에서 뺀 `valkey-secret.yaml`이, ArgoCD(`prune: true`+`selfHeal: true`)가 다음 sync에서 "git에 더 이상 없는 리소스"로 보고 라이브 Secret을 **즉시 삭제**해버림(예전에 ArgoCD가 한 번 적용한 리소스는 git에서 지운다고 안전해지지 않음) — canary 테스트 도중 이걸로 파드가 못 뜨는 걸 발견하고 즉시 재인지. 순수 `kubectl create`로 재생성(ArgoCD가 만든 적 없는 오브젝트라 다시는 안 지워짐)하고 ONBOARDING.md에 이 함정을 기록.

### Kafka 발행 비동기화 (요청 몰릴 때 지연 문제)
- `idHandler`가 `sarama.SyncProducer`로 Kafka produce+ack를 기다린 뒤 응답하고 있어, 동시 요청이 몰리면 이 대기시간이 그대로 누적되는 구조였음(Valkey INCR 자체는 서브밀리초라 병목 아님).
- `sarama.AsyncProducer`로 전환: `Input()`에 non-blocking send 후 즉시 응답, 내부 버퍼가 가득 차는 극단적 상황에서만 드롭+로그(기존에도 Kafka 실패는 로그만 남기고 무시하던 best-effort 정책 유지). `Errors()`/`Successes()`를 계속 비워주는 `drainKafkaProducer` goroutine 추가(안 비우면 결국 `Input()`도 막힘).
- 검증: 로컬 Docker 스모크 테스트(Valkey만) 통과 → 실제 canary 배포 → 클러스터 내 진짜 Kafka로 `/id` 호출 후 consumer 로그에 정상 수신 확인, 순차 요청 30개 전부 200.

### enterprise 네임스페이스 따라잡기
- CI가 `k8s/smb/rollout.yaml`만 자동 패치하는 기존 known gap 때문에 enterprise는 여전히 `v0.3.1`(위 기능들 전부 없는 구버전)에 머물러 있었음 — `k8s/enterprise/rollout.yaml`의 이미지 태그를 수동으로 `sha-4506e5c`(metrics+canary analysis+비동기 Kafka 전부 포함)로 맞춰 커밋/push.
- enterprise에서 canary+analysis 조합이 **처음으로 실전 동작** — 20% → analysis(`Successful`) → 100% 승격까지 확인, `/id`·`/version` 정상 응답.
