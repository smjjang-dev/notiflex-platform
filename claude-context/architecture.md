# Notiflex 아키텍처 스냅샷 — k3s 마이그레이션 시점 (2026-08-13)

이 문서는 AI가 매 대화에서 현재 아키텍처를 빠르게 파악할 수 있도록 한 페이지로 요약한다.

> **마이그레이션 안내**: 이 스냅샷은 ch9 완료 시점(GKE)이 아니라 GKE→k3s 마이그레이션 이후 상태를 기록한다. 마이그레이션 배경·단계별 결정은 `docs/infraPlan.md`, 전환 전 GKE 전용 원본 매니페스트는 `docs/gke-legacy/`를 참조한다. 현재는 **노트북 k3d 클러스터에서 검증 완료**, **Synology NAS 이관은 보류 중**이다.

## 3층 지식 구조

| 문서 | 역할 | 업데이트 주기 |
|------|------|------------|
| **CLAUDE.md** | 프로젝트 메타데이터 (클러스터 환경, 레지스트리, 저장소) | 초기 설정 시 |
| **claude-context/** | 현재 아키텍처 스냅샷 (토폴로지·컴포넌트·파이프라인) | 챕터/마이그레이션 단계 완료 시 |
| **docs/architecture-decisions.md** | 결정 누적 기록 (왜 이 도구를 선택했는가) | 결정 시점마다 |
| **docs/infraPlan.md** | GKE→k3s 마이그레이션 계획 및 Phase별 실측 기록 | Phase 진행 시마다 |

CLAUDE.md는 매 대화 자동 로드, claude-context는 AI가 참조 요청 시, ADR은 사람·AI가 결정 근거 검토 시 사용한다.

## 클러스터 토폴로지

| 항목 | 값 |
|------|-----|
| 클러스터(현재) | k3d `notiflex` (노트북, Docker Desktop 위 단일 노드) |
| 클러스터(목표) | Synology DS920+ 위 k3s 단일 노드, RAM 12GB (4GB+8GB) — 이관 보류 중 |
| 노드 구성 | 단일 노드 — GKE 노드풀(api-pool/worker-pool/ops-pool) 구분 및 관련 nodeSelector는 전부 제거됨 |
| k3s 기본 컴포넌트 | Traefik(Ingress/Gateway), CoreDNS, local-path-provisioner(기본 StorageClass), metrics-server |

## 컴포넌트 다이어그램

```
외부 클라이언트
    │
    ▼ HTTP :80 (Service) → :8000 (Traefik web entryPoint)
Traefik Gateway (gatewayClassName: traefik)
    │ HTTPRoute
    ▼
Service: notiflex-api (stable)
Service: notiflex-api-preview (canary)
    │
    ▼
Rollout: notiflex-api (Canary 전략, Argo Rollouts v1.9.1)
  ├── stable ReplicaSet (ghcr.io/smjjang-dev/notiflex/api)
  └── canary ReplicaSet (배포 시 생성)
    │
    ├── Secret 볼륨 (native K8s Secret)
    │     Secret: valkey / key: valkey-password → /mnt/secrets/valkey-password
    │
    ├── imagePullSecrets: ghcr-pull (private GHCR pull용, 네임스페이스별 재생성 필요)
    │
    ├── Valkey (Bitnami Helm 차트, ADR-008, helm-values/valkey.yaml)
    │     INCR 명령으로 분산 ID 생성
    │
    ├── Kafka (notiflex-kafka-kafka-bootstrap:9092, Strimzi KRaft, v4.3.0)
    │     notifications 토픽으로 이벤트 발행
    │
    └── Tempo (tempo.monitoring:4317)
          OTel OTLP gRPC로 트레이스 전송
```

## 배포 파이프라인

```
코드 변경 (app/)
    │
    ▼ git push → main
GitHub Actions CI (docker/login-action, GITHUB_TOKEN)
    │ docker build + push
    ▼
GHCR: ghcr.io/smjjang-dev/notiflex/api
    │ 매니페스트 이미지 태그 업데이트 → git push
    ▼
ArgoCD (v3.5.1, notiflex-smb / notiflex-enterprise Application)
  auto-sync + selfHeal
    │
    ▼
Argo Rollouts (Canary)
  20% → 30s → 50% → 30s → 80% → 30s → 100%
```

## 관측 가능성

| 도구 | 역할 | 네임스페이스 | 실측 메모리 |
|------|------|------------|------------|
| Prometheus | 메트릭 수집 (kube-prometheus-stack, retention 2d) | monitoring | ~506Mi |
| Grafana | 대시보드·데이터소스 통합 | monitoring | ~317Mi |
| Loki | 로그 저장 (SingleBinary, filesystem) | monitoring | ~170Mi |
| Fluent Bit | 로그 수집 DaemonSet → Loki | monitoring | ~6Mi |
| Alertmanager | 알림 라우팅 (PrometheusRule 연동) | monitoring | ~32Mi |
| Tempo | 분산 트레이싱 (OTLP gRPC) | monitoring | ~29Mi |

전체 노드 사용량(Kafka 포함, 앱 파드 제외): 실측 약 5.26Gi — 상세는 `docs/infraPlan.md` Phase 7/8 참조.

## 주요 네임스페이스

| 네임스페이스 | 주요 워크로드 |
|------------|-------------|
| notiflex | Rollout/notiflex-api, CronJob/notiflex-healthcheck — nodeSelector 없음(단일 노드) |
| enterprise | Rollout/notiflex-api — enterprise 테넌트, nodeSelector 없음 |
| kafka | Strimzi KRaft Kafka 4.3.0(combined controller+broker, replicas 1), KafkaTopic/notifications |
| argocd | ArgoCD v3.5.1 core 구성만(server/repo-server/application-controller/redis, dex·notifications·applicationset 제외) |
| argo-rollouts | Argo Rollouts v1.9.1 컨트롤러 |
| monitoring | kube-prometheus-stack, Loki, Fluent Bit, Tempo |

## 검증 완료 (노트북 k3d, 2026-08-13)

`curl http://localhost/health` → `{"status":"ok","version":"v0.3.1"}`, `curl http://localhost/id` → Valkey INCR 기반 ID 정상 증가, Kafka 이벤트 발행/수신 로그 확인, `notiflex-healthcheck` CronJob 정상 통과. Traefik→HTTPRoute→Service→Rollout→Valkey/Kafka 전체 체인이 실제로 동작함을 확인.

## 알려진 갭 (마이그레이션 범위 밖)

- enterprise 네임스페이스의 이미지 태그는 CI가 자동 관리하지 않음 — 수동 승격 대상 (원본 설계부터 동일)
