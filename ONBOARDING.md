# Notiflex Platform 온보딩 가이드

> **2026-08-13**: GKE에서 k3s로 마이그레이션됨. 아래는 현재(노트북 k3d) 기준 안내다. NAS(Synology DS920+) 이관은 보류 중 — 배경은 `docs/infraPlan.md` 참조. GKE 시절 원본 절차는 `docs/gke-legacy/`에 보관되어 있다.

## 빠른 시작

### 필수 도구

| 도구 | 버전 | 역할 |
|------|------|------|
| k3d (또는 k3s) | 5.x | 로컬/NAS Kubernetes 클러스터 |
| kubectl | 1.35+ | K8s 클러스터 운영 |
| helm | 3.x | Helm 차트 관리 |
| gh | 2.x | GitHub 저장소 관리 (GHCR push/pull 인증 포함) |

### 시크릿 유출 방지 훅 (최초 1회)

이 저장소는 `kind: Secret` 매니페스트에 실제 값이 채워진 채로 커밋되는 걸 막는 pre-commit 훅을 `.githooks/`에 두고 있다. clone 직후 한 번만 활성화하면 된다:

```bash
git config core.hooksPath .githooks
```

`.github/workflows/secret-scan.yaml`이 동일한 검사를 CI에서도 한 번 더 수행한다 (훅을 깜빡했거나 `--no-verify`로 우회한 경우의 백스톱).

### 클러스터 접근

```bash
# 노트북: k3d 클러스터 생성 (최초 1회)
k3d cluster create notiflex --servers 1 --agents 0 -p "80:80@loadbalancer" -p "443:443@loadbalancer"

# kubectl은 k3d가 자동으로 현재 컨텍스트를 설정한다 (별도 --context 불필요)
kubectl get nodes
```

### 주요 리소스 확인

```bash
# ArgoCD Application 상태
kubectl get app -n argocd

# 전체 Pod 상태
kubectl get pods -A

# Gateway/HTTPRoute 상태
kubectl get gateway,httproute -n notiflex

# API 테스트 (Traefik이 로컬 80 포트로 바인딩됨)
curl http://localhost/health
curl http://localhost/id
```

## 아키텍처 요약

`claude-context/architecture.md` 참조 — ch9 완료 시점 스냅샷.

핵심 흐름:
1. 코드 변경 → GitHub → ArgoCD → Argo Rollouts (Canary)
2. /id 호출 → Valkey INCR → Kafka 이벤트 → OTel 트레이스 → Tempo
3. PrometheusRule → Alertmanager 알림

## 의사결정 이력

`docs/architecture-decisions.md` — ADR-001~016 16개 결정 기록

| 챕터 | ADR 범위 | 핵심 결정 |
|------|----------|----------|
| ch3 | 001~002 | ArgoCD, GitHub Actions |
| ch4 | 003~005 | Prometheus+Grafana, Loki+Fluent Bit, PrometheusRule |
| ch5 | 006~007 | Gateway API, Blue/Green |
| ch6 | 008~010 | Valkey, Secret Manager CSI, Canary |
| ch7 | 011~013 | 노드풀 분리, App of Apps, 멀티테넌시 |
| ch8 | 014~016 | Kafka, Tempo, CronJob |

## 트러블슈팅

자주 발생하는 문제:
- ArgoCD Sync 실패 → `argocd.argoproj.io/refresh=hard` 어노테이션
- Valkey 연결 실패 → `kubectl get pod valkey-primary-0 -n notiflex`
- Kafka entity-operator CrashLoop → userOperator 섹션 제거
- GHCR 이미지 pull 401 → private 저장소면 각 네임스페이스에 `ghcr-pull` imagePullSecret 재생성 필요 (Secret은 git에 없음, `gh auth token`으로 생성)
- Valkey 인증 실패 / Secret 재생성 필요 → `valkey-secret.yaml`은 git에 없음(base64가 평문이나 마찬가지라 커밋 금지), 아래처럼 imperative하게 생성:
  ```bash
  kubectl create secret generic valkey \
    --from-literal=valkey-password="$(openssl rand -base64 24)" \
    -n notiflex --dry-run=client -o yaml | kubectl apply -f -
  # enterprise 네임스페이스도 동일 비밀번호로 맞춰줘야 함(같은 Valkey 인스턴스 공유)
  kubectl create secret generic valkey \
    --from-literal=valkey-password="<위와 동일한 값>" \
    -n enterprise --dry-run=client -o yaml | kubectl apply -f -
  # 비밀번호를 바꿨다면 valkey-primary-0을 재시작해야 새 값을 반영함
  kubectl delete pod valkey-primary-0 -n notiflex
  ```
- Strimzi가 `Unsupported Kafka.spec.kafka.version` 에러 → 오퍼레이터 버전이 올라가며 구버전 Kafka 지원이 빠질 수 있음, 오퍼레이터 로그의 지원 버전 목록 확인 후 `k8s/kafka/kafka-cluster.yaml`의 `version` 값 조정
- k3d에서 Fluent Bit DaemonSet `FailedMount: /etc/machine-id` → k3d 노드가 컨테이너라 발생하는 노트북 전용 이슈, `daemonSetVolumes`에서 `etcmachineid` 제거 후 설치(NAS 실기기에서는 발생하지 않음)
