# Notiflex Platform

## 프로젝트 개요

Notiflex — B2B 알림 SaaS 플랫폼. 기업 고객별 알림 채널(이메일, SMS, Slack)을 관리한다.

## 기술 스택

- **언어**: Go 표준 라이브러리 (외부 프레임워크 없음)
- **컨테이너**: scratch 베이스 이미지 (최소 크기)
- **인프라**: k3s (단일 노드) — 2026-08-13 GKE에서 마이그레이션. 상세 배경은 `docs/infraPlan.md` 참조

## 클러스터 설정

- **현재 환경**: 노트북 k3d 클러스터(`notiflex`) — 검증용, Docker Desktop 위에서 구동
- **목표 환경**: Synology DS920+ (RAM 12GB) 위 k3s — NAS 이관은 보류 중 (Phase 10~12)
- **이미지 레지스트리**: ghcr.io/smjjang-dev/notiflex (GHCR, private + imagePullSecrets)
- **Git 저장소**: https://github.com/smjjang-dev/notiflex-platform (private)
- **원본(책) 저장소**: https://github.com/sysnet4admin/notiflex-platform — GKE 기준, 더 이상 배포 대상 아님

## 행동 규칙

- 항상 현재 상태를 확인한 후 작업을 진행한다
- kubectl은 현재 컨텍스트(`k3d-notiflex`)를 사용한다 — GKE `--context` 지정 불필요
- 변경 전 현재 리소스 상태를 먼저 확인한다
- 인프라 변경 계획은 `docs/infraPlan.md`의 Phase 구조를 따른다
