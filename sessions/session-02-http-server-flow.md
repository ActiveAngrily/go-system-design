# Session 2 — HTTP Requests, Responses, and Server Flow

[← Course overview](../study-plan.md) · [← Session 1](session-01-backend-foundations.md)

## Session contract

**Time:** 90–120 minutes  
**Lecture sequence:** L05 — Understanding HTTP for backend engineers  
**Project:** TraceLog version 0.5 — a minimal HTTP health endpoint  
**Goal:** Understand the HTTP boundary well enough to trace one request through a Go server and deliberately form a correct response.

Session 1 taught how a Go program starts and changes in-memory state. This session gives that program a network boundary. The program will still have no JSON API, dynamic routes, database, authentication, or framework. Those come in later sessions.

### What we are not building yet

Do not add an API router library, database, JSON request body, authentication, middleware stack, configuration framework, service/repository layers, or a production health-check system. A single static endpoint is enough to make HTTP behavior visible.

---

## Source material

Read the lecture note first. Open the transcript only to clarify a topic you cannot explain from the note.

| Primary material | What to focus on |
| --- | --- |
| [L05 — Understanding HTTP](../resources/backend-from-first-principles/notes/05-5-understanding-http-for-backend-engineers-where-it-all-starts.md) | Request and response structure, methods, paths, headers, bodies, status codes, CORS, and the practical HTTP lifecycle. |
| [L05 transcript](../resources/backend-from-first-principles/transcripts/05-a3C1DMswClQ.en.vtt) | Revisit the examples of status codes, HTTP method semantics, preflight requests, caching-related headers, and connections as needed. |
| *Let’s Go* §§2.2, 2.4, 2.5 | Supporting Go examples for a minimal server and handler only. Do not let the book change the lecture order. |

---

## Outcomes

By the end of this session, you should be able to:

- Describe an HTTP request as a method, target, headers, and optional body.
- Describe an HTTP response as a status, headers, and optional body.
- Explain the usual lifecycle from a client call to a Go handler and back.
- Recognize the meaning of GET, POST, PUT, PATCH, DELETE, and OPTIONS.
- Select and justify basic response statuses: 200, 201, 204, 400, 404, 405, 500, and 503.
- Read a Go HTTP handler signature and identify the request input and response output objects.
- Explain why headers and status must be decided before writing the response body.
- Test a handler in memory with net/http/httptest rather than requiring a live port.
- Recognize CORS, preflight, proxies, TLS, connection reuse, caching headers, and HTTP versions without needing to implement them today.

You are not expected to memorize the complete HTTP specification, every header, proxy behavior, or every status code. You should be able to trace normal behavior, locate the contract for an endpoint, and look up details before changing behavior.

---

## Active practice and tool lab

### Source-aware tool register

| Tool or technology | Source and role | Depth for this session | Proof of competence |
| --- | --- | --- | --- |
| Postman | Demonstrated in L05 as an HTTP client for sending and inspecting requests. | Hands-on if available. | Create a local TraceLog request, inspect status, headers, and body, then send a wrong method. |
| Browser DevTools Network panel and fetch | Demonstrated in L05 to show what the browser actually sends and receives. | Hands-on. | Trigger one local request, locate it in Network, and identify its method, status, response headers, and body. |
| curl | Course-supporting terminal counterpart to Postman. | Hands-on. | Reproduce the same successful and failed request with curl -i and compare the raw response. |
| React Query | Mentioned in L05 as a client-side caching example. | Recognition only. | Explain that it belongs on a client cache path, not inside this small Go health server; do not add it. |

Do not treat curl or Postman as a substitute for handler tests. They observe the live server;
httptest isolates and proves the handler contract.

### Manual-first rule

Before requesting implementation help, write the GET and POST /healthz contracts and make the
POST test fail. Then manually write the method branch, header decisions, and return after the error
response. You must be able to point to the first response write and explain why execution cannot
continue after it.

### Tool drill

1. Run the server locally and use curl -i for GET /healthz, POST /healthz, and GET /does-not-exist.
2. Repeat the first two requests in Postman if it is available; do not put credentials or secrets in
   a saved environment.
3. In a browser, use fetch against the local endpoint or a small local page, then inspect the request
   in Network. Note that curl does not enforce browser CORS rules; the browser is the correct place
   to observe that behavior.
4. Compare the handler test, curl output, Postman view, and browser view. They should describe the
   same HTTP contract in different forms.

### Completion evidence

Submit the two test names, one raw curl response, and a short explanation of which component chose
each status: the mux or the handler. If a tool result surprises you, diagnose it from the HTTP
message before asking an agent to change code.

---

## Agenda

| Time | Activity | Evidence you produce |
| --- | --- | --- |
| 0–10 min | Diagnostic and recall from Session 1 | Short written answers |
| 10–30 min | L05 synthesis | A request/response map |
| 30–45 min | Read a tiny Go server | A call-flow trace |
| 45–65 min | Build the health endpoint | Formatted code and passing handler tests |
| 65–80 min | Exercise response semantics | A correct 405 behavior and explanation |
| 80–95 min | Intentional bug investigation | Failing test, root cause, smallest fix |
| 95–110 min | AI-generated code audit | Evidence-based review |
| 110–120 min | Checkpoint and progress log | Answers and next-step questions |

If you have only 90 minutes, stop after the bug lab. Do the AI audit before Session 3.

---

## 1. Start with a diagnostic

Answer without searching first.

1. What information does a client send in an HTTP request?
2. What does a response status code communicate that a response body does not?
3. What is the difference between GET and POST at the level of intent?
4. Why should a handler return after writing an error response?
5. What happens if a handler writes a response body before setting its Content-Type header?
6. Which HTTP method should a server recognize when a browser performs a CORS preflight?
7. Why can a handler be tested without calling ListenAndServe?
8. In Session 1, the program called newEntry directly. What component will call behavior in an HTTP backend?

Save the answers in your progress log. We will revisit them at the checkpoint.

### Start-of-session prompt for Codex

> Teach Session 2 from sessions/session-02-http-server-flow.md. Start by asking my diagnostic questions one at a time and make me reason aloud before explaining. Use L05 as the primary source. Keep TraceLog to one small static HTTP endpoint; do not introduce JSON bodies, dynamic routing, databases, frameworks, authentication, or layers yet.

---

## 2. Build the HTTP mental model

Write the two messages in a form you can recognize in logs, browser tools, or tests.

### Request

~~~text
METHOD path?query HTTP/version
Header-Name: value
Header-Name: value

optional body bytes
~~~

### Response

~~~text
HTTP/version status-code reason
Header-Name: value
Header-Name: value

optional body bytes
~~~

The blank line separates headers from the body. Headers carry metadata. The body carries the representation or message. Status tells the caller the high-level result of processing the request.

### TraceLog health request

~~~text
client
  -> GET /healthz
  -> Go server accepts request
  -> mux selects healthHandler
  -> handler checks method
  -> handler chooses status, headers, and body
  -> server writes HTTP response
  -> client observes status, headers, and body
~~~

At this stage, the handler is both the HTTP boundary and the application behavior. That is acceptable because there is one trivial behavior. Later, the handler will translate HTTP into application operations instead of owning all business rules.

### Method intent: understand versus recognize

| Method | Understand today | Do not memorize today |
| --- | --- | --- |
| GET | Read or retrieve a representation; it should not intentionally create or mutate state. | Cache-control details and conditional request syntax. |
| POST | Submit data to create a resource or trigger a non-idempotent action. | Every POST convention. |
| PUT | Replace a whole resource representation when that is the defined contract. | Edge cases involving idempotency and conditional updates. |
| PATCH | Apply a partial change; often the better semantic match for partial updates. | Every patch document format. |
| DELETE | Request deletion of a resource. | Whether a particular API returns 200, 202, or 204. |
| OPTIONS | Discover supported communication options; browsers use it for some CORS preflight requests. | Full browser preflight rules. |

### Status choices: practical core

| Situation | Usual status | Why |
| --- | --- | --- |
| Health endpoint succeeded | 200 OK | The server successfully returned a representation. |
| A resource was created | 201 Created | The request created a new resource. |
| Request succeeded with no response body | 204 No Content | There is no representation to send. |
| Request body or parameters are malformed | 400 Bad Request | The request cannot be processed as sent. |
| Caller is not authenticated | 401 Unauthorized | The caller has not established usable identity. |
| Caller lacks permission | 403 Forbidden | Identity is known but not allowed. |
| Resource or route is absent | 404 Not Found | Nothing matches the requested target. |
| Endpoint exists but method is not allowed | 405 Method Not Allowed | The path is known, but the request method is wrong. |
| State conflicts with requested change | 409 Conflict | The requested action conflicts with current state. |
| Semantic request validation fails | 422 Unprocessable Content | The request has a usable shape but violates rules. |
| Unexpected server failure | 500 Internal Server Error | Server-side behavior failed unexpectedly. |
| Temporarily unavailable service | 503 Service Unavailable | The service cannot currently handle the request. |

A status is part of the API contract. Do not return 200 merely because an error body contains the word “error.”

### Recognition-only topics from L05

Know what these mean and why they matter, but defer implementation details until a feature needs them.

- **CORS:** a browser-enforced policy based on response headers. It is not authentication or authorization.
- **Preflight:** an OPTIONS request sent by browsers before some cross-origin requests.
- **TLS/HTTPS:** protects transport confidentiality and integrity; it does not replace application authorization.
- **Persistent connections:** a connection can carry multiple request/response exchanges.
- **HTTP/2 and HTTP/3:** newer transport behavior; application handler intent still begins with method, target, headers, and body.
- **Proxies and load balancers:** intermediary systems can add, remove, or interpret headers and may produce gateway-related errors.
- **ETag and conditional requests:** tools for caching and concurrent-update behavior; look up exact header semantics when implementing them.

---

## 3. Inspect TraceLog version 0.5

Ask me to generate a tiny server in a new directory named tracelog-session-02.

~~~text
tracelog-session-02/
  go.mod
  main.go
  server.go
  server_test.go
~~~

### Required behavior

- GET /healthz returns 200 with a small plain-text body.
- Any other method at /healthz returns 405.
- The 405 response includes Allow: GET.
- The successful response declares a plain-text Content-Type before writing its body.
- The server constructor returns an http.Handler so tests can call it without binding a port.
- main starts the server and does not hide a startup failure.
- Tests use net/http/httptest.

### Design constraints

- Use only the standard library.
- Use a local http.NewServeMux rather than a global mutable mux.
- A small healthHandler function is enough.
- Do not expose a package-global server, do not add an interface, and do not add a framework.
- Keep the server startup path distinct from the testable handler path.
- Do not write a response more than once.

### The minimum code shape to recognize

~~~text
main()
  -> newServer()
  -> http.ListenAndServe(address, handler)

test
  -> newServer()
  -> httptest.NewRequest(...)
  -> httptest.NewRecorder()
  -> handler.ServeHTTP(recorder, request)
  -> assertions on status, headers, and body
~~~

### Code-reading exercise

Before running the generated server, answer:

1. Which function creates the routing object?
2. Which function receives the request and response objects?
3. Which values in the request are trusted? Which are not?
4. Where does method validation happen?
5. In which order are response header, response status, and response body written?
6. How is a live network listener avoided in the unit test?
7. What does http.ResponseWriter represent at a high level?
8. What should happen after the handler sends a 405 response?
9. Which code path runs only when the executable is started, not when a handler test runs?

### Prompt to generate the inspection code

> Create a new tracelog-session-02 directory containing the smallest standard-library Go HTTP server for TraceLog. Use only go.mod, main.go, server.go, and server_test.go. Implement GET /healthz as plain text and return 405 plus Allow: GET for other methods. newServer must return an http.Handler so tests use httptest without a real port. Show the file tree, request contract, and tests before writing code. Do not add JSON, dynamic routing, a database, middleware, interfaces, frameworks, or configuration systems. Wait for my review after showing the diff.

---

## 4. Manual coding exercise

Implement or complete the health handler yourself. Write its contract before writing it.

### Contract

| Input | Expected output |
| --- | --- |
| GET /healthz | 200, text/plain response, small body such as ok |
| POST /healthz | 405, Allow: GET, error response |
| GET /unknown | 404 from the mux |
| Server cannot start | Startup failure is surfaced rather than ignored |

### Your implementation order

1. Write the GET test.
2. Write the POST-to-healthz test.
3. Implement the smallest handler that makes the tests pass.
4. Format the code.
5. Run the focused test and then all tests.
6. Start the server once and inspect actual HTTP messages with curl.
7. Explain why your handler returns on the error path.

### Verification tasks

Run the testable path first:

~~~text
go fmt ./...
go test ./...
~~~

Then inspect the live path in one terminal:

~~~text
go run .
~~~

From a second terminal:

~~~text
curl -i http://localhost:8080/healthz
curl -i -X POST http://localhost:8080/healthz
curl -i http://localhost:8080/does-not-exist
~~~

For each response, record:

- Status line
- Content-Type
- Allow header when applicable
- Body
- Which component made the decision: mux or handler

### Explanation questions

- Why is 405 more accurate than 404 for POST /healthz?
- Why should Content-Type be selected before a body write?
- Why is the request a pointer in the handler signature?
- Why can a test invoke ServeHTTP directly?
- What is the risk of treating an HTTP status as a cosmetic detail?
- Why is a local mux easier to test than a global shared mux?

### Focused explanation prompt

> Explain this Go HTTP handler line in the context of the TraceLog health endpoint: [paste line]. State what it reads, what it writes or mutates, whether it can end the response, and what test proves its behavior. Do not rewrite the server or introduce later-course abstractions.

---

## 5. Intentional bug lab

Choose one bug. I should inject the bug after the baseline tests pass. You should find it from the symptom and evidence.

| Bug | Faulty change | Expected symptom |
| --- | --- | ---|
| Missing return | Write an error response for a wrong method, then continue to write the success body. | A double-write warning, incorrect body, or misleading response behavior. |
| Late header | Set Content-Type after writing the body. | The client receives an unexpected or default content type. |
| Wrong method status | Return 200 with an error message for POST /healthz. | A client cannot reliably distinguish success from failure. |
| Missing Allow | Return 405 without identifying the allowed method. | The client has less useful recovery information. |
| Test the wrong target | Call a handler function directly in a way that bypasses mux behavior. | The test does not prove the actual request path. |

### Debugging protocol

1. Reproduce the behavior using the smallest test or curl command.
2. Describe the observed HTTP message, not just “it is broken.”
3. Find the first response write in the handler.
4. Trace whether code continues after that write.
5. State the endpoint contract that failed.
6. Add or tighten a test that proves the failure.
7. Make the smallest root-cause fix.
8. Run formatting, focused tests, all tests, and one live request.
9. Explain why a larger refactor is unnecessary.

### Questions I will ask before helping

- Which status, headers, and body did you expect?
- Which status, headers, and body did you actually receive?
- Where is the first response write?
- Could any later code still write the same response?
- What exact assertion would have caught this before manual testing?

### Bug-lab prompt

> I intentionally broke my Session 2 TraceLog HTTP server. Be a debugging partner, not an auto-fixer. Ask for the response observed by curl or the failing httptest assertion, then ask for the smallest relevant handler and test. Help me trace response writes and state the violated HTTP contract. Do not give a patch until I state a likely root cause and a regression test.

---

## 6. Review AI-generated HTTP code

Ask an agent for a solution, then audit it before accepting it.

### Review rubric

| Review question | Evidence of a good answer |
| --- | --- |
| Does the route have a clear method contract? | GET succeeds and unsupported methods result in 405. |
| Are status codes meaningful? | Success and errors are distinguishable to a client. |
| Are headers written in time? | Content-Type and Allow are set before body output. |
| Does an error path stop execution? | A return prevents a second response write. |
| Is the handler testable? | Tests exercise the returned handler with httptest. |
| Is the design proportionate? | One local mux and one handler; no framework or unnecessary interface. |
| Is global state avoided? | Tests do not depend on a shared global mux. |
| Are assumptions explicit? | Port, body format, and health semantics are not silently invented. |

### Common AI-generated-code failure modes

- Returning a JSON-like error body with HTTP 200.
- Calling http.Error and then continuing into the happy path.
- Setting headers after Write or Encode has already committed the response.
- Using a global default mux that leaks routes between tests.
- Ignoring ListenAndServe errors.
- Writing an elaborate router, middleware system, or generic response abstraction for one endpoint.
- Claiming CORS is authentication, or adding permissive CORS headers without a real cross-origin requirement.
- Testing an isolated helper while never testing the actual registered route.

### Audit prompt

> Review this generated Session 2 Go HTTP diff before I accept it. Verify route behavior, method handling, status codes, header order, response write paths, server-startup error handling, httptest coverage, and unnecessary abstractions. Cite exact lines and separate must-fix defects from optional simplifications. Do not replace the implementation unless I request it.

### Assumption-challenge prompt

> List every external behavior this HTTP handler assumes: request method, path, headers, body, response type, port, client behavior, and server lifecycle. Mark each as enforced, tested, documented, or merely assumed. Recommend only the smallest changes needed for this session.

---

## 7. Checkpoint: trace and change one behavior

Complete this without opening the generated code first.

1. Trace GET /healthz from the client command through the Go program and back to the client.
2. Explain the difference between a request path and a handler.
3. Explain why POST /healthz should be 405 rather than 404.
4. Put response headers, status, and body in the correct order and explain why.
5. Explain one way an HTTP response can accidentally be written twice.
6. Describe the difference between CORS and authentication.
7. Write a handler test that proves the Allow header is present for the wrong method.
8. Read an agent proposal to add a routing framework. Accept, reject, or defer it using the current requirements.

### Passing standard

You are ready for Session 3 when you can:

- Read a raw HTTP request and response at a useful high level.
- Predict the behavior of GET, wrong-method, and unknown-path requests.
- Explain the flow from a handler test through ServeHTTP.
- Diagnose a response-writing bug from evidence.
- State why this server is intentionally not a JSON API yet.
- Name one HTTP detail you would look up rather than guess.

If one answer is weak, repeat only the matching exercise. Do not restart the whole session.

---

## 8. Update your progress log

~~~markdown
## Session 2 — HTTP server flow

- Date and actual time:
- Diagnostic answer I changed:
- My request-to-response trace:
- TraceLog folder or commit:
- Chosen HTTP bug, symptom, root cause, and regression test:
- Status code I can justify:
- HTTP topic I recognize but will look up when needed:
- AI review finding:
- Confidence (1–5):
- Question to carry into Session 3:
~~~

---

## Ready for Session 3

The next session separates two ideas that are often confused:

~~~text
routing chooses a handler
serialization/deserialization converts data representations
~~~

Bring the Session 2 server and prepare to answer:

- How should GET /entries/42 reach a different behavior from GET /entries?
- How does JSON become a Go value without becoming trusted?
- What should a handler return when a request body is malformed?

> Teach Session 3 only after I show you my Session 2 progress log and can explain why headers must be written before the response body.
