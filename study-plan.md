# Interactive Go Backend Course: Four Weeks with Codex

## Course promise

This is a compact, hands-on Go course built in the **same progression as the Backend from first
principles lecture series**. The lecture notes and transcripts are the primary curriculum. The
Let’s Go book is supporting material when you need to see an idea expressed in a concrete Go
application.

You will not pass by watching videos or memorizing APIs. In every session you will inspect code,
explain execution, write a small amount manually, diagnose an intentional defect, test a change,
and review AI-generated code before accepting it.

The goal is to make you able to read, reason about, diagnose, and modify Go backends with AI
assistance while keeping the engineering judgment yourself.

## Course sources

Primary sources:

- Lecture notes: resources/backend-from-first-principles/notes/
- Lecture transcripts: resources/backend-from-first-principles/transcripts/

Supporting source only:

- [Let’s Go](<resources/Alex Edwards - Let's Go (2022, Alex Edwards) - libgen.li.md>)

When a lecture note is unclear, open its matching VTT transcript. When a Go implementation detail is
unclear, use the supporting book section and then current Go documentation. Do not reverse that
order: the course follows the lecture curriculum, not the book table of contents.

## Time commitment and format

- **Length:** four weeks.
- **Sessions:** five guided sessions per week.
- **Daily commitment:** 90–120 minutes for a core session, plus an optional 20–30 minute recall or
  cleanup session.
- **Ratio:** about 30% lecture/reading, 70% tracing, coding, debugging, testing, and review.
- **One evolving project:** a realistic Go backend named **TraceLog**.

TraceLog lets a user create study/work records, attach tags and optional file metadata, retrieve
records, and receive a simple status update. It evolves from an in-memory HTTP service to a
PostgreSQL-backed, authenticated, observable service with one background task. We will discuss
caching, queueing, search, object storage, real-time updates, and scale at the right time; we will
not add infrastructure merely to check a box.

## How I will teach you

At the start of each session, tell me the session number and send your progress note. I will then:

1. Give a five-minute model of the topic and the minimum vocabulary.
2. Generate or point you to a small focused codebase or feature for inspection.
3. Ask you to predict behavior before running code.
4. Give you one manual coding exercise that is small enough to understand fully.
5. Give you a code-reading/explanation task and one intentional bug.
6. Ask you to propose a test and a smallest safe fix.
7. Review your reasoning before giving a solution.
8. Help you review any AI-generated diff against Go conventions, tests, security, and scope.
9. End with a short checkpoint and update your progress tracker.

### Hint and solution rules

| If you ask for… | I will do |
| --- | --- |
| An explanation | Explain the code path, data flow, call stack, error path, and relevant Go behavior. I will not hide difficult details behind “the framework handles it.” |
| A hint | Start with a conceptual hint, then a trace hint, then pseudocode. I will show code only after you have made a reasonable attempt or explicitly ask for it. |
| A bug | Give a reproducible symptom and enough context, not the location or answer. I will ask for a hypothesis and a test before confirming it. |
| A solution | First ask you to state the invariant and expected behavior. Then show the smallest implementation and explain why alternatives are unnecessary. |
| An AI review | Separate blockers from optional improvements, cite code evidence, and never invent findings. |
| A feature | Help you define behavior, boundaries, tests, and a small plan before generating a diff. |

### How you should interact with Codex

Use these commands or equivalent natural-language requests:

~~~text
Session 3. My confidence is 2/5. I completed the routing exercise but I cannot explain why the
handler receives a pointer to Request. Teach today’s lesson using the course rules.
~~~

~~~text
Do not edit yet. Trace this code from route registration to response. Explain each function, the
data it reads/writes, errors it can return, and tests that cover it.
~~~

~~~text
Give me one intentional bug for the current session. Tell me the symptom and how to reproduce it,
but do not reveal the location or solution. Review my diagnosis after I respond.
~~~

~~~text
Review this generated diff before I accept it. Check behavior, Go idioms, validation, authorization,
context cancellation, resource cleanup, tests, security, dependencies, and over-engineering.
~~~

## AI-agent rules for the entire course

Never ask an agent to “build the whole feature” without a scoped behavior and test plan. Before you
accept generated code:

- Ask for a read-only explanation or trace first.
- Require the agent to list files, callers, assumptions, and tests before editing.
- Compare unfamiliar standard-library behavior against Go documentation.
- Inspect the diff yourself. Ask what it did not change and why.
- Write or agree on a failing test before accepting a behavior-changing fix.
- Prefer one vertical change over a cross-repository rewrite.
- Challenge new dependencies, interfaces, frameworks, queues, caches, and abstractions.
- Debug failures from evidence: reproduce, observe, narrow the cause, fix the root, add a regression
  test. Do not ask the agent to blindly rewrite the codebase.

## Active learning and tool-proficiency contract

Every session treats tools and libraries as part of the engineering work, not as magic hidden behind
an agent. Each detailed session now has a source-aware tool register and a small tool lab.

### Accurate tool provenance and required depth

| Label in a session | Meaning | Required outcome |
| --- | --- | --- |
| Lecture-demonstrated | The lecturer visibly uses or explicitly names it for the topic. | Use or inspect it in one bounded exercise when it is safe and relevant. |
| Lecture-mentioned | The lecture uses it as an example or comparison. | Explain its role, trade-offs, and place in the system; do not install it by reflex. |
| Course-supporting | A Go tool, standard-library package, terminal client, or browser feature needed to practise the lesson. | Use it manually and relate its output to the code path. |

For a hands-on tool, “I know it” means you can explain what job it does, run one focused operation,
inspect its raw output, diagnose one controlled failure or surprising result, and say whether the
current TraceLog stage needs it. Recognition-only technologies are deliberately deferred, not ignored.

### Non-negotiable working method

1. Start with the behavior: write the request, expected result, invariant, or failing test before
   asking for implementation.
2. Make a real manual attempt in every session: one code path, test, debugging step, or small
   feature change. Inspection code may be generated as a learning target, but it is not accepted
   until you can trace it.
3. Before accepting any line, explain its inputs, outputs, mutation, error path, dependency, and any
   resource, concurrency, or security effect that applies.
4. Use documentation or a tiny experiment for behavior you cannot justify. “The agent said so” is
   never evidence.
5. Review generated diffs line by line. State what behavior they change, what tests prove it, and
   what they intentionally leave unchanged.

### Strict teaching and hint ladder

I will begin with a question, a prediction, or a test rather than a complete patch. If you are stuck,
I will escalate only as far as needed: (1) a Socratic question, (2) a trace or counterexample,
(3) the relevant documentation location, (4) one local code clue or pseudocode, then (5) the
smallest solution after a genuine attempt or an explicit request. I will give you bugs as symptoms
and reproduction steps first, require a hypothesis and a regression test, and not let a broad
“rewrite it” replace diagnosis.

### Library and tool adoption gate

Before we add a dependency, service, or framework, you must answer:

1. What concrete problem in the current code does it solve?
2. What can the Go standard library or existing code do instead?
3. What is the smallest API surface we would actually use?
4. What version, security, configuration, lifecycle, and test obligations does it add?
5. Which documentation page and focused test will verify our use of it?

We will use a tool when it helps us observe or build the current behavior. We will not install a
library merely because a lecture mentioned it or an agent suggested it.

---

## Lecture-led course map

| Week / phase | Lecture progression preserved | Why topics are combined | Depth |
| --- | --- | --- | --- |
| 1. Request foundations | L01–L09 | backend model → HTTP → routing → data representation → identity/validation | Deep practice for Go basics, HTTP, routing, serialization, validation |
| 2. Application and data | L10–L16 | layers/context → API design → PostgreSQL → caching/jobs/search → fault handling | Deep practice for packages, interfaces, errors, SQL, transactions |
| 3. Production behavior | L17–L22 | config → observability → shutdown/security → scale → concurrency | Deep practice for context, goroutine lifecycle, security, diagnosis |
| 4. Delivery and judgment | L23–L27 | storage → real time → testing → twelve-factor application | Working familiarity for storage/realtime; deep practice for tests/review/production judgment |

The sequence is never reversed. For example, we do not discuss a background worker before we can
trace an HTTP handler, and we do not add caching before we understand the PostgreSQL query it would
avoid.

## What to understand deeply versus recognize

### Understand deeply

- Go execution flow, functions/methods, values/pointers, structs, slices/maps, packages, errors,
  defer, interfaces at real call sites, context, and goroutine ownership.
- Request/response flow, routing, serialization, validation, authorization, middleware order, and
  error translation.
- PostgreSQL query boundaries, parameterization, connection pooling, transactions, data constraints,
  and useful tests.
- Security boundaries, cancellation, resource cleanup, observability, graceful shutdown, and the
  process for reviewing AI-generated changes.

### Recognize and look up when necessary

- Exact standard-library function signatures and uncommon syntax.
- Exact HTTP status codes and headers beyond common API use.
- Cache library settings, queue technology APIs, search engine mappings, cloud object-storage SDK
  calls, WebSocket/SSE APIs, metrics SDKs, TLS settings, and deployment manifests.
- Advanced generic patterns, database tuning syntax, profiling flags, and distributed-system
  technology choices.

The rule is: remember the concept, lifecycle, constraints, and questions. Look up implementation
details just before you need them.

---

# Week 1 — Request foundations: L01–L09

## Week purpose

You will learn what happens between a client request and a Go handler, how Go code executes, and how
to turn untrusted JSON into a validated application operation. This is the base for every later
lecture.

## Sessions and lecture order

| Session | Lectures | Teaching focus | Supporting book material |
| --- | --- | --- | --- |
| 1 | L01 Roadmap; L02 Backend engineer path; L03 What is a backend; L04 First principles | backend mental model, Go program flow, packages, functions, values | §§1, 1.1, 2.1, 2.2 |
| 2 | L05 HTTP | client/server flow, methods, headers, status, request/response | §§2.2, 2.4, 2.5 |
| 3 | L06 Routing; L07 Serialization/deserialization | handlers, routing, JSON, structs, slices, maps, pointers | §§2.3, 2.9 |
| 4 | L08 Authentication/authorization | identity versus permission, trust boundaries, methods/interfaces | §§9.1–9.3, 11.1–11.6 as reference only |
| 5 | L09 Validation/transformation | validation, normalization, error responses, test-first behavior | §§8.1–8.6 |

Primary files:

- notes/01-1-roadmap-for-backend-from-first-principles.md through
  notes/09-9-validations-and-transformations-for-backend-engineers.md
- Matching transcripts 01 through 09 when you need detail.

## What I will teach

- A backend is a program that receives requests, applies behavior and policy, reads/writes data, and
  responds. The handler is only one stage in that flow.
- Go code execution starts at main, builds dependencies, registers routes, and waits for requests.
  A function call, method call, pointer, and returned error all affect that flow.
- HTTP gives request method/path/headers/body and response status/headers/body. Routing selects a
  handler; routing does not validate, authenticate, or authorize.
- Deserialization turns JSON into Go values. Those values are untrusted until validation and
  normalization complete.
- Authentication asks who the caller is. Authorization asks whether that caller may do this action on
  this resource. Validation asks whether the input has an acceptable shape.

## Go concepts you must understand

- Variables, types, zero values, conditionals, loops, functions, multiple return values, and errors.
- Structs, fields, slices, maps, pointers, methods, and when mutation is visible to a caller.
- Package imports and exported versus unexported names.
- The basic HTTP handler interface and how a function can satisfy it.
- JSON struct fields/tags at a recognition level; look up exact encoding options.

## Codebase I will generate for inspection

I will generate TraceLog version 1: a deliberately small in-memory API with:

- GET /healthz
- GET /entries
- POST /entries
- GET /entries/{id}

It will use only standard net/http, encoding/json, and simple structs. You will read it before
changing it. The code will contain no database, framework, generic repository, or authentication
library.

## Manual coding, reading, bugs, and features

| Type | Assignment |
| --- | --- |
| Manual coding | Write a small Entry struct, a validation function, and one handler branch manually. |
| Code reading | Trace POST /entries from route registration to JSON response. Explain each function and every possible returned status. |
| Intentional bug | Diagnose one of: ignored JSON decode error, missing return after an error response, map mutation through a nil map, incorrect route match, or validation occurring after storage. |
| Feature addition | Add durationMinutes and tags, then normalize tags by trimming, lowercasing, deduplicating, and bounding them. |
| Tests | Write handler tests for valid JSON, malformed JSON, missing title, duplicate tag, unknown ID, and wrong method. |
| AI audit | Review an AI-generated handler for response double-writes, unbounded body input, hidden errors, incorrect content type, and invented abstractions. |

## Week 1 practical Codex prompt

~~~text
Session 3. Generate the smallest standard-library Go TraceLog API needed for GET /entries,
POST /entries, and GET /entries/{id}. Before writing code, list the package layout, route behavior,
data invariants, error cases, and tests. Use no database, router library, or interface unless it has
a concrete job. After generating it, give me a code-reading task and do not reveal the answer.
~~~

## Week 1 checkpoint

Without opening the code, explain the lifecycle of POST /entries for:

1. valid JSON;
2. malformed JSON;
3. missing title;
4. an unknown route;
5. a caller that uses the wrong HTTP method.

Then explain the difference among routing, deserialization, validation, authentication, and
authorization. You pass when your explanation identifies where each decision belongs.

---

# Week 2 — Application boundaries and data: L10–L16

## Week purpose

You will turn the in-memory service into a small PostgreSQL-backed backend and learn to reason about
layers, contexts, errors, transactions, cache decisions, background work, search, and fault
tolerance without overbuilding it.

## Sessions and lecture order

| Session | Lectures | Teaching focus | Supporting book material |
| --- | --- | --- | --- |
| 6 | L10 Controllers, services, repositories, middleware, request context | package boundaries, composition, interfaces, request-scoped context | §§3.1–3.5, 6.1, 12.1–12.2 |
| 7 | L11 Complete REST API design | API contract, resource behavior, pagination, error semantics | §§2.4–2.5, 7.1–7.2 |
| 8 | L12 PostgreSQL | schema, queries, pools, transactions, constraints | §§4.1–4.9 |
| 9 | L13 Caching; L14 Task queues/background jobs | cache-aside judgment, async work, idempotency, failure choices | §§3.1–3.2 as supporting patterns |
| 10 | L15 Search; L16 Fault tolerance/error handling | search boundaries, errors, retries, degradation, root-cause fixes | §§3.4, 4.9, 14.1–14.2 |

Primary files:

- notes/10-10-what-are-controllers-services-repositories-middlewares-and-request-context.md
- notes/11-11-complete-rest-api-design.md
- notes/12-12-mastering-databases-with-postgres.md
- notes/13-13-caching-the-secret-behind-it-all.md
- notes/14-14-task-queues-and-background-jobs.md
- notes/15-15-full-text-search-using-elasticsearch-for-blazingly-fast-search.md
- notes/16-16-error-handling-and-building-fault-tolerant-systems.md

## What I will teach

- A handler translates HTTP into an application operation. Storage code owns SQL. Middlewares handle
  cross-cutting concerns. This is a tool for clarity, not a mandatory number of packages.
- Interfaces are behavioral contracts. Make one only at a genuine dependency or substitution
  boundary; a one-implementation interface is often needless ceremony.
- Context transports cancellation, deadlines, and request-scoped values. It should not transport
  configuration, a database pool, or arbitrary global state.
- A database/sql handle is a pool. A query must use parameters, context, bounded results, and cleanup.
  Transactions protect an invariant that must change together.
- Caching, queues, and search are trade-offs. Start with direct PostgreSQL work until a measured or
  specified problem requires more infrastructure.
- Fault tolerance means predictable failure behavior, not hiding every error or adding retries
  everywhere.

## Go concepts you must understand

- Package direction and import boundaries.
- Interfaces and composition at call sites.
- Error wrapping/classification, errors.Is, errors.As, defer, Close, and transaction rollback/commit.
- Context propagation.
- SQL query/row/result lifecycle and basic PostgreSQL concepts.
- Recognition only: cache client APIs, external queues, search engine client APIs.

## Codebase I will generate for inspection

I will generate TraceLog version 2:

- PostgreSQL migrations for users, entries, and tags.
- A PostgreSQL storage implementation alongside a narrow in-memory test double when one is justified.
- A single application service that creates an entry and tags atomically.
- List pagination with a fixed maximum limit and stable sort order.
- Request logging and error translation middleware.

## Manual coding, reading, bugs, and features

| Type | Assignment |
| --- | --- |
| Manual coding | Write one PostgreSQL query with parameters and context; implement one transaction boundary. |
| Code reading | Draw the call chain for create entry: route → middleware → handler → application behavior → database → response. |
| Intentional bug | Diagnose one of: string-built SQL, missing rows.Close, transaction never rolled back, N+1 query, context ignored, or duplicate entry reported as a 500. |
| Feature addition | Add tags and tag filtering with pagination. Explain why this does or does not justify a cache. |
| Tests | Write a PostgreSQL integration test that proves a failed tag insert leaves no partial entry. |
| AI audit | Review generated SQL for parameters, transactions, cleanup, cardinality, index assumptions, pagination limits, and database error leakage. |

## Week 2 practical Codex prompt

~~~text
Session 8. I have this PostgreSQL schema and TraceLog handler. Do not edit yet. Trace the code,
identify the true storage boundary, list every query and transaction, and propose a test matrix for
not found, duplicate data, cancellation, query failure, and partial transaction failure. Then suggest
the smallest safe implementation plan.
~~~

## Week 2 checkpoint

Explain:

1. why a database/sql handle is not one physical connection;
2. why application validation does not replace database constraints;
3. when a transaction is necessary;
4. why a cache or task queue is not automatically an improvement;
5. where a context cancellation should reach.

You pass when you can use a concrete TraceLog request to answer each question.

---

# Week 3 — Production behavior: L17–L22

## Week purpose

You will learn how a backend behaves under failure, load, shutdown, and attack. The practical focus is
not “add every production system,” but trace lifetimes and defend trust boundaries in Go code.

## Sessions and lecture order

| Session | Lectures | Teaching focus | Supporting book material |
| --- | --- | --- | --- |
| 11 | L17 Configuration management | environment/config boundaries, secrets, startup validation | §§3.1–3.2 |
| 12 | L18 Logging, monitoring, observability | structured logs, metrics/traces concepts, diagnosis | §§3.2, 6.3 |
| 13 | L19 Graceful shutdown; L20 Backend security | lifecycle, timeouts, cancellation, threat model, access control | §§6.2, 6.4, 10.1–10.4, 11.1–11.7 |
| 14 | L21.1 Scaling | bottlenecks, load, state, horizontal/vertical trade-offs | §§4.4, 10.4 |
| 15 | L21.2 Performance; L22 Concurrency/parallelism | measurement, goroutines, channels, synchronization, race safety | §§6.5, 14.8 |

Primary files:

- notes/17-17-production-grade-configuration-management.md
- notes/18-18-logging-monitoring-and-observability.md
- notes/19-19-graceful-shutdown.md
- notes/20-20-backend-security-everything-you-need-to-know.md
- notes/21-21-1-backend-scaling-and-performance-engineering-part-1.md
- notes/22-21-2-backend-scaling-and-performance-engineering-part-2.md
- notes/23-22-concurrency-parallelism-io-bound-vs-cpu-bound.md

## What I will teach

- Configuration is an explicit, validated contract. Secrets are never source code, logs, or error
  messages.
- Logs answer “what happened”; metrics quantify “how often/how much”; traces connect work across
  boundaries. For this course, structured logs are enough to implement; metrics/traces are concepts
  to recognize.
- Graceful shutdown means stop accepting new work, cancel or drain owned work, wait with a deadline,
  close dependencies, and exit predictably.
- Security is a boundary mindset: validate input, parameterize operations, enforce authorization at
  point of access, avoid secret leakage, and test abusive paths.
- Scaling/performance start with a measured bottleneck. Concurrency is independent progress;
  parallelism is simultaneous CPU work.
- Every goroutine must have an owner, inputs, outputs, an exit condition, cancellation behavior, and
  an error decision. A channel is not a default replacement for a mutex or a durable queue.

## Go concepts you must understand

- Context cancellation, deadlines, and propagation to database/network work.
- Goroutines, channels, select, mutexes, WaitGroup, race detector, and ownership.
- Error logging versus client-visible errors.
- Signals/shutdown behavior at a conceptual and implementation level.
- Recognition only: advanced metric/tracing SDKs, autoscaling configuration, profiling options, and
  distributed rate limiting.

## Codebase I will generate for inspection

I will generate TraceLog version 3:

- Config loaded from environment with startup validation.
- Structured request/error logging using log/slog.
- Health and readiness endpoints.
- An in-process bounded summary worker triggered after an entry is created.
- Graceful shutdown with context cancellation and a deadline.
- A deliberately small security boundary around authenticated owner access.

The worker is intentionally in-process. We will assess whether an external queue is needed; we will
not add one by default.

## Manual coding, reading, bugs, and features

| Type | Assignment |
| --- | --- |
| Manual coding | Add a cancellation-aware worker loop or one context-aware database operation. |
| Code reading | Identify every goroutine and write owner, start, stop, shared state, cancellation, and error paths. |
| Intentional bug | Diagnose one of: leaked goroutine, unbounded channel, worker using a request context after response, data race, shutdown that closes DB too early, secret in a log, or missing object-level authorization. |
| Feature addition | Add a weekly summary background job and a readiness endpoint. State why it remains in-process. |
| Tests | Run go test -race; write a shutdown/cancellation test and one two-user authorization test. |
| AI audit | Review an AI-generated concurrency/security diff for lifecycle gaps, races, hidden retries, weak permissions, secret leakage, and unsupported performance claims. |

## Week 3 practical Codex prompt

~~~text
Session 15. Audit this concurrent Go code. For every goroutine, identify owner, trigger, inputs,
outputs, shared state, cancellation path, exit condition, error path, and test strategy. Flag leaks,
races, deadlocks, unbounded work, and work that outlives a request. Do not rewrite it yet.
~~~

## Week 3 checkpoint

Given a handler that launches a goroutine, explain:

1. why it might be unsafe even if it makes the response faster;
2. what context it should use;
3. who stops it;
4. what happens when the server shuts down;
5. how the race detector or a test could reveal a bug.

Then threat-model one TraceLog endpoint with assets, actors, trust boundaries, abuse cases, and tests.

---

# Week 4 — Delivery and engineering judgment: L23–L27

## Week purpose

You will cover the remaining backend capabilities without pretending you must operate every platform.
You will then focus deeply on testing, reviewing, and delivering a small production-oriented backend.

## Sessions and lecture order

| Session | Lectures | Teaching focus | Supporting book material |
| --- | --- | --- | --- |
| 16 | L23 Object storage part 1; L24 Object storage part 2 | large-file flow, metadata, streaming, multipart and authorization concepts | templates/static-file examples only as reference |
| 17 | L25 Real-time backends | polling versus SSE/WebSockets, connection lifecycle, backpressure | no book-first prerequisite |
| 18 | L26 Testing, mocks, TDD, coverage | test pyramid, fakes/mocks, handler/integration tests, coverage limits | §§14.1–14.8 |
| 19 | L27 Twelve-factor app | config, statelessness, disposability, logs, dev/prod parity | §§3.1, 3.2, 10.4 |
| 20 | Capstone assessment | audit, bug diagnosis, small feature delivery, final review | §§17.1–17.6 as optional exercises |

Primary files:

- notes/24-23-object-storage-everything-you-need-to-know-part-1.md
- notes/25-24-object-storage-everything-you-need-to-know-part-2.md
- notes/26-25-real-time-backends.md
- notes/27-26-testing-for-backend-engineers-mocks-tdd-and-coverage.md
- notes/28-27-the-twelve-factor-app.md

## What I will teach

- Object storage is for blobs; databases normally store metadata. Understand signed-upload style flows,
  streaming, size/content validation, and authorization. Do not memorize a cloud SDK.
- Real-time delivery has connection lifetime, reconnection, ordering, backpressure, and authorization
  consequences. Recognize SSE/WebSocket patterns; build only a tiny local status stream if useful.
- Tests are evidence, not a coverage ritual. Unit tests prove small rules; handler tests prove HTTP
  behavior; integration tests prove database interactions; mocks are only useful at a real boundary.
- Twelve-factor principles make deployment behavior explicit: config, stateless processes, backing
  services, logs, disposability, and environment parity.
- Production API quality means stable request/response behavior, bounded work, useful errors, access
  control, tests, observability, and a deliberate scope.

## Go concepts you must understand

- Test tables, httptest, integration test boundaries, mocks/fakes, and coverage interpretation.
- IO/streaming and resource cleanup at a conceptual code-reading level.
- Long-lived connection and goroutine lifecycle at a code-reading level.
- Recognition only: cloud object-storage SDKs, message brokers, full WebSocket frameworks, and
  deployment platforms.

## Codebase I will generate for inspection

I will generate TraceLog version 4 or a focused review branch containing:

- Attachment metadata plus a fake object-storage boundary; no real cloud account required.
- A tiny server-sent event or polling-status example to inspect.
- A test suite with unit, handler, integration, and intentionally poor mock tests.
- A production review document listing configuration, security, dependencies, operational behavior,
  and known scope limits.

## Manual coding, reading, bugs, and features

| Type | Assignment |
| --- | --- |
| Manual coding | Add one focused handler test and one integration test; write a small file-metadata validation rule. |
| Code reading | Compare a synchronous upload path, a background upload path, and an SSE endpoint. Explain ownership and backpressure. |
| Intentional bug | Diagnose one of: uploaded file accepted before authorization, unbounded file body, leaked SSE goroutine, mock that proves nothing, flaky integration test, or coverage increase with no behavior assertion. |
| Feature addition | Add attachment metadata or a status event with bounded behavior. Do not add a real cloud SDK unless a requirement needs it. |
| Tests | Create a test matrix for file metadata, unauthorized access, worker shutdown, handler errors, and database failure. |
| AI audit | Review an AI-generated feature for hidden infrastructure, weak cleanup, brittle tests, missing authorization, and unneeded dependencies. |

## Week 4 practical Codex prompt

~~~text
Session 18. Review this Go test suite as evidence, not as a coverage score. Classify each test as
unit, handler, integration, or implementation-coupled. Identify missing negative cases, false
confidence, fragile mocks, and the smallest high-value tests to add. Do not change production code.
~~~

## Week 4 checkpoint

Explain the difference among a unit test, handler test, integration test, and end-to-end test using
TraceLog. Then explain why object storage, task queues, caching, and real-time connections are
architectural trade-offs rather than mandatory backend features.

---

# Final assessment: TraceLog production slice

At the end of Week 4, complete this assessment with my guidance.

## Required behavior

TraceLog must have:

- An HTTP API with explicit routes, JSON decode/encode, validation, and useful error responses.
- PostgreSQL persistence with parameterized queries, connection-pool setup, constraints, pagination,
  and one meaningful transaction.
- Authenticated identity and owner authorization for protected records.
- Configuration validation, structured logs, health/readiness, and graceful shutdown.
- One bounded in-process background summary task.
- Tests for validation, handlers, authorization, PostgreSQL behavior, failure paths, and one
  concurrency/shutdown concern.
- A short README explaining local run/test commands, configuration, known limits, and how to verify
  the primary flow.

## Assessment tasks

1. **Reading:** I give you an unfamiliar TraceLog route. You trace it and identify data, error,
   authorization, context, and test paths.
2. **Bug diagnosis:** I introduce one defect. You reproduce it, form hypotheses, write/identify a
   test, find the root cause, and propose the smallest fix.
3. **Feature change:** You implement a bounded feature, such as tag filtering, an attachment-metadata
   endpoint, or a status field, with a plan and tests.
4. **AI review:** Codex generates a plausible diff. You audit it for correctness, security,
   lifecycle, tests, dependencies, Go idioms, and over-engineering.
5. **Design defense:** You explain why the project does not yet need Redis, a queue broker,
   Elasticsearch, microservices, or a custom framework.

## Final capstone prompt

~~~text
Act as my final Go-backend mentor. Do not edit yet. Read this TraceLog repository and give me:
1) an architecture map; 2) request traces for the main endpoints; 3) the highest-risk boundaries;
4) missing tests; 5) a security and concurrency review; and 6) three small, evidence-based
improvements. Cite files and lines. Then ask me to choose one improvement and require a test plan
before implementation.
~~~

---

# Progress tracker

Keep this simple. At the end of each session, append one row to a table in your personal notes or
send it to me.

| Session | Lecture(s) | Finished | Confidence 1–5 | One thing I can explain | One unanswered question | Evidence | Tool/library proof |
| --- | --- | --- | ---: | --- | --- | --- | --- |
| Example | L05 | yes | 3 | why HTTP is stateless | when to use cache-control | handler test passed | inspected a 405 in curl and Postman |

Evidence must be concrete: a test name, a command you ran, a request you traced, a bug you fixed, or
a diff you rejected. Tool/library proof should name the operation you performed and what it showed.
“Watched the lecture” is not evidence of understanding.

Use this confidence rule:

- **1:** I recognize the words but cannot trace code.
- **2:** I can follow an explanation but cannot predict behavior.
- **3:** I can trace a small example and write a guided change.
- **4:** I can diagnose a related defect and review an AI change.
- **5:** I can explain trade-offs and teach the concept using this codebase.

When you report a confidence of 1 or 2, I will give a smaller code-reading task instead of advancing
to more infrastructure.

---

# Final competency checklist

## Go code understanding

- [ ] I can trace Go program execution from main through package setup and a request handler.
- [ ] I can read and reason about functions, methods, return values, and control flow.
- [ ] I understand variables, types, zero values, pointers, slices, maps, structs, and mutation.
- [ ] I understand packages, imports, exported names, project organization, and composition.
- [ ] I understand interfaces as behavior contracts and can identify unnecessary abstractions.
- [ ] I can follow errors, error wrapping/classification, defer, cleanup, and failure paths.
- [ ] I understand goroutines, channels, mutexes, cancellation, race detection, and the difference
  between concurrency and parallelism.

## Backend request and API understanding

- [ ] I can trace HTTP request/response flow, methods, headers, bodies, status codes, and handlers.
- [ ] I can reason about routing, route parameters, serialization, and deserialization.
- [ ] I can distinguish authentication, authorization, validation, and transformation.
- [ ] I can explain controllers, services, repositories, middleware, and request context without
  treating them as compulsory layers.
- [ ] I can review API design for clear resource behavior, bounded pagination, error contracts, and
  production-oriented practices.

## Data and systems understanding

- [ ] I can reason about PostgreSQL schemas, SQL queries, parameterization, transactions, constraints,
  connection pooling, errors, and query lifecycle.
- [ ] I understand caching purpose, freshness, invalidation, and why it is not a default addition.
- [ ] I understand background jobs/task queues, idempotency, failure semantics, and why one worker can
  be sufficient at first.
- [ ] I recognize search-system trade-offs and when PostgreSQL search versus a search service might
  matter.
- [ ] I understand configuration management, secrets, logging, monitoring, observability, readiness,
  graceful shutdown, fault tolerance, scaling, and performance measurement.
- [ ] I understand object storage/large-file flows and real-time backend connection lifecycle at a
  working code-review level.
- [ ] I understand twelve-factor application principles and can apply the useful parts to a small Go
  service.

## Testing, security, and AI judgment

- [ ] I can write and interpret unit, handler, integration, mock/fake, TDD, coverage, and regression
  tests.
- [ ] I can diagnose bugs by reproducing, tracing, testing a hypothesis, fixing the root, and adding
  a regression test.
- [ ] I can audit authorization, input validation, SQL injection risk, secret handling, error
  leakage, resource cleanup, cancellation, and common boundary failures.
- [ ] I can ask Codex to explain code without hiding detail, review a diff, challenge assumptions,
  compare behavior against Go documentation, and reject over-engineered or unsafe generated code.
- [ ] I can use AI for small, well-tested implementation tasks without asking it to blindly rewrite a
  failure.

## Start here

Send me this message when you are ready:

~~~text
Start Session 1. My Go background is: [briefly describe it]. I have [90/120] minutes today.
Teach L01–L04 using TraceLog. Give me a short diagnostic first, then a code-reading task, one
manual exercise, one intentional bug, one test task, and a checkpoint. Do not give solutions until
I show my reasoning.
~~~
