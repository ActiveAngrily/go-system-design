# Go Backend Course Progress

## Completed: Session 1 — Backend Foundations and Go Program Flow

**Status:** Complete  
**Project:** `tracelog_session_01/`  
**Verified:** `go test ./...` passes.

### What I learned

- A backend follows: input → validation → business rules → state → output.
- A Go module is the project declared by `go.mod`; a package groups Go files; an executable starts at `main()` in `package main`.
- Structs hold related fields; functions define behavior; methods attach behavior to a type through a receiver.
- `*Entry` is a pointer to an `Entry`; constructors can return `nil, err` when no valid value exists.
- Maps use unique keys. Reading a nil map is safe; writing one panics. `AddTag` initializes `Tags` with `make(map[string]bool)` before writing.
- Expected failures are returned as `error`; panics signal broken program state or assumptions.

### Evidence

- Created the `tracelog_session_01` Go module.
- Implemented and tested blank-title rejection, title trimming, tag addition, and summary output.
- Reproduced the nil-map panic by temporarily removing the map-initialization guard, then restored the safe fix.
- Used `go mod init`, `go fmt`, `go test`, `go run`, `go list`, and `go doc strings.TrimSpace`.

### Current TraceLog shape

- `Entry` has `ID`, `Title`, and `Tags`.
- `NewEntry` trims and rejects blank titles, returning `(*Entry, error)`.
- `AddTag` is a pointer-receiver method and safely initializes a nil tag map.
- Tests pass.

### Resume point

Start **Session 2 — HTTP Server Flow** from `sessions/session-02-http-server-flow.md`.

Suggested prompt for the next chat:

> Continue my Go backend course from `progress.md`. I completed Session 1 and am ready for Session 2. Teach the theory before asking diagnostic questions, then guide me through the exercise one small step at a time.

## Completed: Session 2 — HTTP Server Flow

**Status:** Complete  
**Project:** `tracelog_session_02/`  
**Verified:** `go test ./...` passes; live GET, wrong-method, and unknown-path requests were checked with `curl -i`.

### What I learned

- An HTTP request contains a method, target path, headers, and an optional body.
- An HTTP response contains a status, headers, and an optional body.
- A router/mux selects a handler by request path; the handler reads `*http.Request` and writes through `http.ResponseWriter`.
- `405 Method Not Allowed` is correct when a route exists but does not support the requested method; `404 Not Found` means no route matches.
- Response headers must be set before the first status or body write because the first write commits the response and Go may infer `Content-Type`.
- `httptest.NewRequest`, `httptest.NewRecorder`, and `ServeHTTP` test a handler without opening a network port.
- CORS is a browser-enforced cross-origin reading policy, not authentication or authorization.

### Evidence

- Created `tracelog_session_02/go.mod`, `server.go`, `server_test.go`, and `main.go`.
- Implemented `GET /healthz` with `200`, plain text, and `ok` body.
- Implemented `405` for unsupported methods with `Allow: GET`.
- Verified `404` for an unknown path from the mux.
- Diagnosed and fixed a missing `return` after the error response; the regression test catches the incorrect `Method not allowed` plus `ok` body.
- Used `gofmt`, `go test`, `go run`, and `curl -i`.

### Session checkpoint

- Request-to-response trace: client → server listener → mux → handler → response writer → client.
- Status codes justified: `200` success, `404` unknown route, `405` unsupported method.
- AI/code review finding: keep the local standard-library mux and set `Allow` before `http.Error`; no framework or abstraction is needed.
- HTTP topic to look up when needed: detailed CORS preflight headers and browser behavior.
- Confidence (1–5): not recorded.
- Question for Session 3: how does routing distinguish paths such as `/entries` and `/entries/42`, and how does JSON become a Go value?
