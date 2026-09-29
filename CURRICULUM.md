# CURRICULUM.md — hookrelay 학습 커리큘럼

**현재 단계: 1 (Go 기초)**

각 단계는 "완료 기준"을 만족하면 끝난다. 끝나면 커밋하고 상단의 현재 단계를 갱신한다.

---

## 0. 환경 구성

- Go 설치, VS Code + Go 확장, Git 설치, GitHub 저장소 연결
- `go version`, `go run` 으로 hello world 실행

**완료 기준**: hello world가 실행되고, 첫 커밋이 GitHub에 올라가 있다.

---

## 1. Go 기초 (1~2주)

프로젝트 폴더 안 `learn/` 디렉터리에서 진행. 나중에 삭제하거나 .gitignore 처리 가능.

- 패키지, 변수, 타입, 함수, 다중 반환값
- 슬라이스, 맵, 구조체, 메서드
- 인터페이스 (작게 정의하고 암묵적으로 구현한다는 점)
- 에러 처리 (`error` 값, `errors.Is/As`, `fmt.Errorf("%w")`)
- 포인터
- goroutine, channel, `sync.WaitGroup`, `select`
- `context.Context` (취소와 타임아웃)
- `encoding/json` 직렬화·역직렬화
- 테스트 작성 (`_test.go`, 테이블 기반 테스트)

참고: A Tour of Go, Go by Example

**완료 기준**: 고루틴 3개가 채널로 결과를 보내고 메인이 받아서 출력하는 프로그램을 설명 없이 짤 수 있다.

---

## 2. 웹훅 개념 실습 (30분~1시간)

- `learn/receiver.go`: `/hook` POST를 받아 본문을 출력하고 200 응답
- `learn/sender.go`: receiver로 JSON POST
- receiver를 끄고 sender 실행 → 실패 확인
- sender에 재시도(3초 간격, 최대 5회) 추가
- receiver를 꺼둔 채 sender 실행 → 도중에 receiver 켜기 → 도착 확인

**완료 기준**: "웹훅이 뭐고, 왜 재시도가 필요한지"를 자기 말로 설명할 수 있다.

---

## 3. 접수 API + 메모리 큐 + 워커 1개

- `cmd/server/main.go`: `net/http` 서버 기동, graceful shutdown
- `POST /webhooks` — `{url, body}` 접수, ID 발급, 202 응답
- `GET /webhooks/{id}` — 상태 조회
- 메모리 큐(채널) + 워커 고루틴 1개가 꺼내서 전송
- 상태: `pending` → `delivering` → `delivered` / `failed`
- `slog`로 구조화 로깅

**완료 기준**: curl로 접수하면 receiver에 도착하고, 조회 API에 `delivered`가 보인다.

---

## 4. 워커 풀 + 재시도·백오프

- 워커 N개 (설정 가능)
- 지수 백오프 (1s, 2s, 4s, … 상한), 최대 시도 횟수
- 시도마다 `attempt` 기록 (시각, 응답 코드, 에러)
- 최종 실패 시 `dead` 상태 (DLQ 개념)
- 전송 타임아웃 (`context.WithTimeout`)
- 동시성 안전 (`sync.Mutex` 또는 채널로 상태 보호)
- 취소 API: `POST /webhooks/{id}/cancel`

**완료 기준**: receiver가 500을 돌려주는 상황에서 재시도 로그가 백오프 간격으로 찍히고, 최대 횟수 후 `dead`가 된다. `go test -race` 통과.

---

## 5. PostgreSQL 영속화

- Docker로 Postgres 실행
- `store` 인터페이스 정의 → 메모리 구현과 Postgres 구현(`pgx`) 두 개
- 마이그레이션 (goose 또는 golang-migrate)
- 서버 시작 시 `pending`/`delivering` 상태 작업을 큐에 복구
- 설정: 환경변수로 DB 접속 정보

**완료 기준**: 접수 → 서버 강제 종료 → 재시작 → 밀린 작업이 전송된다.

---

## 6. 서명·이력·운영 API

- HMAC-SHA256 서명 헤더 (`X-Hookrelay-Signature`), 수신 측 검증 예제
- 멱등성 키 헤더 (`X-Hookrelay-Delivery-Id`)
- `GET /webhooks` — 상태·기간 필터, 페이징
- `GET /webhooks/{id}/attempts` — 시도 이력
- `POST /webhooks/{id}/retry` — `dead` 작업 수동 재전송
- 간단한 API 키 인증 미들웨어

**완료 기준**: receiver가 서명을 검증하고, 잘못된 서명은 거부한다.

---

## 7. 테스트

- 핸들러 단위 테스트 (`httptest.NewRecorder`)
- 전송 로직 테스트 (`httptest.NewServer`로 가짜 수신 서버, 500·지연·정상 시나리오)
- 재시도 정책 테이블 기반 테스트
- Postgres 통합 테스트 (Docker 필요 시 build tag로 분리)

**완료 기준**: `go test ./...` 통과, 핵심 경로 커버리지 확인.

---

## 8. Docker·CI·배포·데모

- 멀티스테이지 Dockerfile, `docker compose` (server + postgres + receiver)
- GitHub Actions: lint(`go vet`, `staticcheck`) + test
- 클라우드 배포 (Fly.io 또는 VM)
- `cmd/receiver`: 실패 모드 옵션(`--fail-rate`, `--delay`) 지원
- README: 한 줄 소개, 아키텍처 그림, API 표, 데모 GIF(receiver 껐다 켜기 → 재전송)

**완료 기준**: 저장소 링크 하나로 남이 실행하고 이해할 수 있다.

---

## 이후 확장 (선택)

- SSE로 전송 상태 스트리밍
- 두 번째 작업 종류(이미지 리사이즈) 추가 → 큐가 범용임을 증명
- Prometheus 메트릭
