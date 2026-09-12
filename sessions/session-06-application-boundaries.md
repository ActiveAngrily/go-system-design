# Session 6 — Handlers, Services, Stores, Middleware, and Request Context

[← Course overview](../study-plan.md) · [← Session 5](session-05-validation-and-transformation.md)

## Session contract

**Time:** 105–120 minutes  
**Lecture sequence:** L10 — Controllers, services, repositories, middleware, and request context  
**Project:** TraceLog version 2.1 — a small, testable boundary refactor  
**Goal:** Learn to trace a request across package boundaries and introduce only the abstractions that make the existing behavior clearer, safer to modify, and ready for real storage in Session 8.

Sessions 1–5 deliberately kept TraceLog flat so you could see every operation. The code now has enough independent responsibilities to separate:

~~~text
HTTP translation
application behavior
storage access
cross-cutting request work
request-scoped cancellation and identity
~~~

This is not a rule that every Go application needs five layers. It is a practical exercise in responsibility, dependency direction, and testability.

### What we are not building yet

Do not add:

- PostgreSQL, SQL, migrations, a cache, background jobs, search, or a queue.
- A framework, dependency-injection container, code generator, or generic repository base class.
- A controller, service, and repository package for every type just because those names exist.
- Interfaces with no actual substitution or behavior-driven test need.
- A global error framework that obscures normal handler behavior.
- A request context containing a database handle, configuration object, raw credential, mutable global state, or arbitrary optional parameters.

The in-memory store remains. Session 8 will introduce PostgreSQL after you can identify the actual storage boundary.

---

## Source material

Read L10 as one connected lesson. The book supports Go examples and package recognition but does not replace the lecture’s architecture progression.

| Primary material | What to focus on |
| --- | --- |
| [L10 — Controllers, services, repositories, middleware, and request context](../resources/backend-from-first-principles/notes/10-10-what-are-controllers-services-repositories-middlewares-and-request-context.md) | Request lifecycle; handler/controller responsibility; application behavior; storage ownership; reusable middleware; request context for identity, request IDs, deadlines, and cancellation. |
| [L10 transcript](../resources/backend-from-first-principles/transcripts/10-hyc-7w3pee8.en.vtt) | Revisit the lifecycle diagrams, reasons for separation, authentication context example, error middleware discussion, and cancellation/deadline discussion. |
| *Let’s Go* §§3.1–3.5, 6.1, 12.1–12.2 | Supporting material for packages, middleware, context, and error patterns. Treat older structural examples as reference, not a mandatory folder layout. |

---

## Outcomes

By the end of this session, you should be able to:

- Trace a TraceLog request from route registration through middleware, handler, application behavior, storage, and response.
- Explain what a handler, service, repository or store, middleware, and request context each own.
- Identify code that is HTTP-specific versus code that can run without HTTP.
- Explain package dependency direction and spot an import cycle before it happens.
- Use a narrow interface at a genuine consumer-side storage boundary and explain why it exists.
- Explain why Go interfaces describe behavior and usually belong near the code that consumes that behavior.
- Pass a request context from handler to application behavior to storage without replacing it with context.Background.
- Explain why verified identity can be carried in request context, but business behavior should receive the actor explicitly.
- Keep validation, authorization, and state changes visible in the correct layer.
- Test both a handler boundary and an application/storage boundary.
- Audit AI-generated architecture for import cycles, misplaced HTTP code, generic repositories, interface ceremony, context misuse, and error translation bugs.

You do not need to memorize a universal package layout, every Go project convention, dependency-injection framework, middleware library, context helper, or repository pattern. You should be able to justify each dependency and explain the request flow without guessing.

---

## Active practice and tool lab

### Source-aware tool register

| Tool or technology | Source and role | Depth for this session | Proof of competence |
| --- | --- | --- | --- |
| go list -deps | Course-supporting view of the actual package dependency graph behind L10’s boundary discussion. | Hands-on. | Run it, identify one direct and one transitive dependency, and compare the result with your hand-drawn package graph. |
| go doc context and go test | Course-supporting tools for request-context behavior and test-preserving refactoring. | Hands-on. | Verify one context behavior in documentation and run tests after every refactor slice. |
| go vet | Course-supporting static check for obvious Go issues. | Working familiarity. | Run it after the refactor and investigate any finding rather than suppressing it. |
| Middleware, repositories, interfaces, and dependency injection | L10 architecture concepts; no specific framework is required. | Hands-on with standard-library code; framework recognition only. | Justify each package and the one narrow interface from an actual consumer and test need. |

No framework, dependency-injection container, or generic repository is needed. A tool that draws or
lists dependencies should make the existing graph easier to audit, not be used to hide it.

### Manual-first rule

Do not ask an agent to mass-refactor the project. First draw the current request path and package
graph, record the HTTP contracts that must not change, and manually move one small responsibility:
the entry creation use case or one pure type/function. Make the project compile and tests pass before
the next move. You must explain each import and why context travels to the next layer.

### Tool drill

Run these from the refactored module:

~~~text
go list -deps ./...
go doc context.Context
go test ./...
go vet ./...
~~~

Use the output to answer: which package owns HTTP details, which owns entry behavior, which owns the
map-backed storage, and which package composes concrete dependencies? If a test or compiler error
appears, trace the first broken dependency or changed contract; do not request a blind rewrite.

### Completion evidence

You finish with a hand-drawn package graph, one manually completed refactor slice, passing handler
and service tests, a documented context-propagation assertion, and a line-by-line explanation of the
one interface that remains. If you cannot justify it, remove or defer it.

---

## Agenda

| Time | Activity | Evidence you produce |
| --- | --- | --- |
| 0–10 min | Diagnostic and Session 5 recall | Written answers |
| 10–28 min | L10 responsibility map | A responsibility and dependency table |
| 28–43 min | Package graph and context walkthrough | A directed import graph |
| 43–65 min | Inspect TraceLog version 2.1 | An end-to-end call trace |
| 65–83 min | Test-first refactor | One extracted application operation |
| 83–98 min | Intentional architecture bug lab | Root cause and regression check |
| 98–112 min | AI-generated-code audit | Evidence-based design review |
| 112–120 min | Checkpoint and progress log | Session 7 readiness |

---

## 1. Start with a diagnostic

Answer without searching first.

1. What part of a Go backend should know about HTTP status codes and JSON?
2. What part should know about how an Entry is stored in memory or a database?
3. What part should own the rule that a caller becomes the owner of a newly created Entry?
4. What is middleware useful for that a handler should not repeat route by route?
5. What is request context useful for?
6. What must not be put into request context?
7. What does a Go interface describe: an inheritance hierarchy or a set of behavior?
8. Why can an interface with one implementation be needless?
9. Why should an application package not import an HTTP package?
10. What happens if a handler calls context.Background instead of passing r.Context to downstream work?

Record your answers in your progress log. Uncertainty is useful evidence, especially here: architecture is easy to name and hard to reason about.

### Start-of-session prompt for Codex

> Teach Session 6 from sessions/session-06-application-boundaries.md. Use L10 as the primary source. Begin with the diagnostic questions one at a time. Help me refactor only the existing TraceLog responsibilities into a small Go package graph. Before suggesting any interface or package, ask what concrete behavior, dependency, or test need it serves. Do not introduce PostgreSQL, frameworks, generic repositories, dependency injection, or a new architecture for its own sake.

---

## 2. The request-lifecycle model

Use this as the reference flow for a protected create request.

~~~text
POST /entries
  -> route selects HTTP handler
  -> authentication middleware verifies credential
  -> middleware attaches verified principal to this request context
  -> handler decodes bounded JSON and calls prepareEntry
  -> handler extracts actor from request context
  -> handler calls entry service with r.Context, actor, and draft
  -> service assigns owner and applies application policy
  -> store creates Entry in memory
  -> service returns Entry or categorized error
  -> handler maps result to HTTP status and JSON response
~~~

The layer names are less important than the responsibilities and direction of data flow.

### Responsibility map

| Component | Owns | Must not own today |
| --- | --- | --- |
| Route registration | Connects method/path to a handler chain. | JSON parsing, business policy, or storage internals. |
| Middleware | Reusable request-wide work such as verified identity attachment. | Endpoint-specific Entry creation or storage rules. |
| Handler or controller | HTTP request decoding, calling application behavior, and HTTP response mapping. | Direct data structure internals, SQL, or complex cross-use-case policy. |
| Service or application use case | Behavior such as server-assigned ownership and calling the required storage operations. | HTTP status, ResponseWriter, raw JSON, or concrete SQL syntax. |
| Store or repository boundary | Reads and writes Entry state. | HTTP headers, JSON, caller tokens, or response formatting. |
| Concrete memory store | Implements current in-memory data access. | Application authorization policy or HTTP decisions. |
| Request context | Cancellation, deadline, and narrow request-scoped metadata. | Global configuration, database pools, raw credentials, persistent domain state, or ordinary optional arguments. |
| Composition root | Creates concrete dependencies and connects them at startup. | Endpoint behavior or business decisions. |

### Handler versus controller

The lecture uses handler and controller together. In this Go course, use **handler** for the function or method that speaks net/http. The important question is not the label; it is whether the code translates HTTP rather than becoming the whole application.

### Service versus helper

A service is not a class-shaped bucket for every function. Create one when a meaningful application operation needs to coordinate policy and dependencies.

For TraceLog, creating an entry is now a reasonable use case because it coordinates:

- A verified actor.
- A prepared draft from Session 5.
- Server-assigned ownership.
- A storage operation.
- A stable error outcome for the handler to map.

Trimming tags or encoding JSON is not a service responsibility.

### Repository or store

Repository is a common name for code that owns persistence operations. Go code may simply call it Store. The name is less important than this boundary:

~~~text
application behavior asks for entry storage capability
concrete storage implementation owns data-structure or database details
~~~

The service must not know whether an Entry sits in a map today or PostgreSQL later.

---

## 3. Package direction and composition

### Target package graph

Use a small package graph rather than a deep folder tree.

~~~text
cmd/tracelog
  -> internal/httpapi
  -> internal/entries
  -> standard library

cmd/tracelog
  -> internal/memory
  -> internal/entries
  -> standard library

internal/httpapi
  -> internal/entries

internal/memory
  -> internal/entries
~~~

The arrows point from the importing package toward the imported package.

The important rule is:

~~~text
entries does not import httpapi
entries does not import memory
memory does not import httpapi
httpapi does not import memory
main is allowed to know concrete pieces and wire them together
~~~

This makes main the **composition root**. It is the one place allowed to say, in effect, “use this memory store with this entry service and this HTTP server.”

### Suggested file layout

~~~text
tracelog-session-06/
  go.mod
  cmd/
    tracelog/
      main.go
  internal/
    entries/
      model.go
      service.go
      store.go
      service_test.go
    memory/
      store.go
    httpapi/
      server.go
      server_test.go
~~~

This is enough for the current responsibilities. Do not create folders called common, utils, helpers, base, shared, infrastructure, or repository merely to feel enterprise-ready.

### The narrow store interface

The consumer, entries.Service, needs storage behavior. It can define the smallest contract it consumes.

~~~go
type Store interface {
    Create(context.Context, Entry) (Entry, error)
    GetByID(context.Context, int) (Entry, error)
    ListByOwner(context.Context, string) ([]Entry, error)
}
~~~

Introduce this interface only when it has a concrete job:

- The in-memory store implements it now.
- A small test-only store can return controlled failures or record calls for service tests.
- PostgreSQL will implement the same consumer-needed behavior in Session 8 if the contract still fits.

Do not create a generic Repository interface with methods for every imagined data operation. Do not name it IEntryStore. Go interface names usually describe the capability, not an implementation prefix.

### Interface decision checklist

Before adding an interface, answer all five:

1. Which code consumes this behavior?
2. Which exact methods does it need now?
3. Is there a real substitution, such as memory versus PostgreSQL or a focused test double?
4. Would the code be clearer if it depended directly on the concrete type instead?
5. Can the interface live in the consumer package without creating an import cycle?

If you cannot answer these, keep the concrete dependency.

---

## 4. Context: cancellation, deadlines, and request values

### Two related but distinct uses

| Context use | Example | Rule |
| --- | --- | --- |
| Cancellation and deadlines | r.Context reaches a future database query or outbound call. | Pass it unchanged down the call chain. Do not replace it with context.Background. |
| Request-scoped value | Authentication middleware attaches a verified principal. | Use a private typed key and a narrow accessor. |
| Observability value | Future request ID used by logging and tracing. | Recognize today; implement with structured logging in Session 12. |

### Correct TraceLog flow

~~~text
HTTP handler receives r.Context()
  -> service.Create(r.Context(), actor, draft)
  -> store.Create(ctx, entry)
  -> future database call uses ExecContext or QueryContext with that same ctx
~~~

The service receives the actor explicitly because ownership is business data. The handler may obtain that actor from request context because middleware verified it. This avoids making every business function search for hidden context values.

### Context rules

- Use context.Context as the first parameter of a function when that operation can be cancelled, time-limited, or needs request-scoped metadata.
- Pass the received context onward. Do not start a new background context inside request work.
- Do not store an SQL database handle, logger, configuration, global store, or raw credential in context.
- Do not use a string context key in shared code. Use a private typed key.
- Do not use context values as a bag of optional function arguments.
- Do not retain a request context after its request completes.
- Do not assume an in-memory store makes context unimportant. The service boundary is where you preserve cancellation for Session 8 and later calls.

### Context propagation experiment

A test-only recording store can receive the context passed to Create. Use a marker in the request context inside the test, then assert the store can observe it. The marker is a test of propagation, not a production pattern for arbitrary values.

Then answer:

- Which code received the original request context?
- Which code forwarded it?
- Which code would break cancellation if it substituted context.Background?
- Which context values should the application service avoid reading implicitly?

---

## 5. Error boundaries and response mapping

Error categories can cross the service/store boundary without HTTP details leaking inward.

### Minimal categories for TraceLog

| Application meaning | Example internal category | Handler responsibility |
| --- | --- | --- |
| Entry absent | ErrNotFound | Map to 404 response. |
| Caller lacks access | ErrForbidden | Map to documented 403 or hiding 404 policy. |
| Store cannot complete operation | ErrUnavailable or wrapped cause | Map to safe 500 or temporary-failure policy. |
| Input invalid | Validation field errors from preparation | Map to 422 response before service call. |
| Missing identity | Authentication wrapper failure | Map to 401 before handler behavior. |

When a store wraps an error, use a stable cause that callers can recognize with errors.Is. The handler should decide HTTP status. The memory store and future PostgreSQL code must not call http.Error or know JSON response shapes.

### Error-handling rules for this session

- Keep errors close to the operation that failed.
- Add useful operation context when wrapping an error.
- Preserve a category only when a caller needs to make a decision.
- Do not convert every error into a generic string.
- Do not leak a database or implementation error to a client.
- Do not create a global error middleware merely to avoid a clear local response map.
- Session 10 will deepen fault handling, retries, and error policy. Today, preserve explainable categories.

---

## 6. Inspect TraceLog version 2.1

Create a new directory named tracelog-session-06 by refactoring Session 5 behavior without changing the published HTTP contract.

### Required behavior

- The existing routes, authentication behavior, ownership rules, validation rules, normalized JSON response, and tests remain valid.
- HTTP handlers no longer directly mutate the memory map.
- A small entries.Service owns entry-creation and retrieval behavior that needs actor or store coordination.
- The entries package defines a narrow Store contract only if its test and upcoming storage boundary justify it.
- The memory package implements the store behavior without importing HTTP code.
- Authentication remains one small HTTP middleware wrapper that attaches verified identity to request context.
- Handlers obtain the actor from context and pass it explicitly to the service.
- The service and store accept and forward context.Context.
- Main wires the concrete memory store, service, and HTTP server.
- Handler tests exercise actual routes and middleware; service tests exercise behavior with a focused test store.

### Refactor constraint

Your first goal is behavior preservation, not a feature. Existing Session 5 tests should remain green before adding any new test double or new package.

### End-to-end create trace to annotate

~~~text
POST /entries
Authorization: Bearer test-user-one
Content-Type: application/json

{"title":"Plan the storage boundary","durationMinutes":45,"tags":["go","architecture"]}
~~~

Fill in the actual function and package for each step:

| Step | Package | Function or type | HTTP-specific? | Can fail? |
| --- | --- | --- | --- | --- |
| Route registration | | | | |
| Authentication wrapper | | | | |
| Principal extraction | | | | |
| JSON decode | | | | |
| Preparation and field errors | | | | |
| Service create operation | | | | |
| Store create call | | | | |
| Error categorization | | | | |
| JSON response | | | | |

### Code-reading exercise

Before running code, answer:

1. Which package owns the Entry type and why?
2. Which package owns the Store interface and why?
3. Which package owns the concrete map and mutex?
4. Which package knows about http.ResponseWriter?
5. Which package knows about the test-only token strings?
6. Which function builds concrete dependencies?
7. Which function receives request context first?
8. Which functions pass context onward?
9. Can entries import httpapi without breaking the graph? Explain.
10. Which layer decides HTTP 404 versus 500?
11. Which layer assigns owner ID?
12. Which existing Session 5 behavior should not change merely because files moved?

### Prompt to generate the inspection code

> Create a new tracelog-session-06 directory by refactoring my Session 5 TraceLog API into the smallest clear Go package graph. Use cmd/tracelog as the composition root; internal/httpapi for net/http handlers and the existing test-only auth wrapper; internal/entries for Entry types, prepared drafts, a small service, and a narrow consumer-owned Store contract; and internal/memory for the concrete map-backed store. Preserve every current HTTP contract and test behavior. Pass request context from handler to service to store; extract a verified actor from request context but pass it explicitly to the service. Add a focused store test double only if it proves a service error or context test. Before writing code, show the package graph, responsibilities, interface justification, preserved behavior list, error mapping, file tree, and tests. Do not add PostgreSQL, frameworks, a dependency-injection container, generic repositories, extra interfaces, or broad middleware.

---

## 7. Manual coding and test-first refactor work

### Refactor in safe slices

Do not move every file first and hope it compiles. Work in this order.

1. Write down the existing HTTP contract and run all Session 5 tests.
2. Move or define the Entry, draft, and field-error types in entries.
3. Extract a pure service method for Create that accepts context, actor, and prepared draft.
4. Keep the existing handler response mapping while changing it to call the service.
5. Extract the concrete in-memory map operations into memory.
6. Add the narrow Store interface only when service tests need controlled storage behavior.
7. Move dependency construction into cmd/tracelog/main.go.
8. Re-run existing handler tests after every small refactor.
9. Add one service test for a real behavior or failure boundary.
10. Draw the final package graph and explain every import.

### Test matrix

| Test type | What it proves |
| --- | --- |
| Existing handler test | Method, auth middleware, JSON contract, validation, and status behavior still work end to end. |
| Service creation test | A verified actor becomes owner and prepared draft values reach storage. |
| Service forbidden/get test | Ownership policy stays in application behavior rather than handler-only code. |
| Store failure test | A controlled storage error is returned with a recognizable category for handler mapping. |
| Context propagation test | The context supplied to service reaches the store boundary. |
| Package compile test | go test ./... proves imports and package declarations form a valid graph. |

### Manual feature: add a service-level operation

Implement this behavior with guidance:

> Move entry creation so the handler no longer writes to the memory map. The service must assign owner ID from the explicit actor, then call Store.Create with the received context.

Start with tests:

1. A normal handler create request still returns normalized JSON and owner-scoped behavior.
2. A service test proves actor ID becomes Entry owner ID.
3. A test-only failing store proves the service returns a categorized storage failure.
4. A handler test proves a storage failure does not leak its internal text.

### Reasoning questions

- Why should the handler not call a map or future SQL query directly?
- Why should the service not receive http.ResponseWriter?
- Why does the Store interface belong with the service consumer instead of the memory implementation?
- What makes the test-only store a genuine reason to have an interface?
- Why should actor be explicit in service.Create even though middleware placed it in context?
- Why should a store accept context even when a map lookup cannot currently block?
- Why is moving code not enough to prove behavior has been preserved?
- Which current helper should remain a pure function rather than become a method on Service?

### Focused implementation prompt

> Guide me through the Session 6 TraceLog refactor in very small, test-preserving steps. First make me identify the current handler, application behavior, and storage responsibilities. Then help me extract only the Entry creation use case into an entries.Service, pass context and explicit actor, and place a narrow Store interface beside its consumer if a focused test double justifies it. Require go test ./... after each step. Do not add PostgreSQL, a framework, a generic repository, a dependency container, new interfaces without a caller, or a package per concept.

---

## 8. Intentional architecture bug lab

Choose one bug after the refactor’s baseline tests pass. Diagnose the root boundary violation, not merely the compiler symptom.

| Bug | Faulty behavior | Why it matters |
| --- | --- | --- |
| Import-cycle bug | entries imports httpapi to use Principal or response types. | Lower-level behavior depends on transport and Go rejects the cycle. |
| HTTP leakage bug | Store receives ResponseWriter or calls http.Error. | Storage cannot be reused or tested outside HTTP. |
| Service bypass bug | One handler still writes directly to memory store, skipping service policy. | Authorization or ownership rules become inconsistent across endpoints. |
| Context reset bug | Service or store uses context.Background instead of received ctx. | Cancellation, deadlines, and request metadata disappear. |
| Context-bag bug | Service reads arbitrary config, store, or raw credential from context. | Dependencies become hidden and hard to test. |
| Generic repository bug | An interface includes many unused CRUD methods for future needs. | API surface and test burden grow without a current benefit. |
| Error-mapping bug | Store error becomes a 500 even when it means not found or forbidden. | Client behavior and diagnosis become incorrect. |
| Middleware-order bug | Protected handler runs before identity middleware. | Handler sees missing or spoofable identity. |
| Global-state bug | Concrete store or authenticated principal is held in a mutable package global. | Tests leak state and request behavior becomes hard to reason about. |
| Validation migration bug | Refactor moves validation after the service/store write. | Session 5 invariant is broken by a structural change. |

### Debugging protocol

1. Reproduce with a compile error, focused unit test, or handler test.
2. State which responsibility or dependency rule is violated.
3. Draw the actual import or call edge that caused it.
4. Identify the lowest layer that should own the missing data or decision.
5. Add a test if the defect is behavioral; use the compiler as evidence if it is an import cycle.
6. Make the smallest move or dependency inversion that restores direction.
7. Run go fmt ./... and go test ./....
8. Explain why a new framework, package, or interface is not needed.

### Questions I will ask before helping

- What exact package imports what?
- Which layer needs this data, and why?
- Is the data HTTP transport detail, application behavior, storage detail, or request metadata?
- Could the dependency be passed explicitly instead?
- Does the service need an interface, or does the handler just need a concrete dependency?
- Did a refactor change behavior or only location?
- Which existing test would catch a policy bypass?

### Bug-lab prompt

> I intentionally introduced a Session 6 package-boundary, context, or layering bug. Be a debugging partner, not an auto-refactorer. Ask for the package graph, compile/test failure, and smallest relevant code first. Help me identify the violated dependency direction or responsibility, draw the call/import path, and add a behavioral regression test where applicable. Do not propose a redesign until I state the smallest boundary correction.

---

## 9. Audit AI-generated architecture

Architecture code can contain a lot of ceremony that looks professional while making a small Go service harder to trace.

### Audit rubric

| Review question | Evidence of a good answer |
| --- | --- |
| Does each package have one clear reason to exist? | HTTP, entries behavior, memory storage, and wiring have distinct responsibilities. |
| Do imports point inward? | Lower-level entries and memory packages do not import HTTP transport. |
| Is main the composition root? | Concrete memory store and HTTP server are wired in main, not hidden in globals. |
| Is the interface justified and narrow? | Consumer-owned methods correspond to service behavior and a real memory/test substitution exists. |
| Does the handler stay transport-focused? | It decodes, gets actor, calls service, and maps response. |
| Does the service stay transport-neutral? | No HTTP status, JSON, ResponseWriter, or raw header use. |
| Does the store stay storage-focused? | No authorization header, response encoding, or policy decisions. |
| Is context propagated correctly? | r.Context reaches service and store; no background replacement. |
| Is context used narrowly? | Verified request identity uses a typed key; business actor is explicit. |
| Are errors categorized without leaking internals? | Handler maps stable causes to HTTP behavior and tests safe output. |
| Is behavior preserved? | Existing endpoint tests remain and new service tests cover the extracted responsibility. |
| Is complexity proportionate? | No factory, base repository, interface-per-type, framework, or dependency container. |

### Common AI-generated-code failure modes

- Creating a controller, service, repository, interface, mock, and factory for one map operation.
- Moving HTTP response code into a service or store.
- Moving business authorization entirely into middleware even though it depends on a specific Entry.
- Putting raw Authorization tokens, configuration, database pools, or mutable maps in context.
- Calling context.Background in a request path.
- Using public string context keys that can collide.
- Giving interfaces names like IEntryRepository or adding every possible CRUD method.
- Defining an interface in the implementation package rather than at the consumer boundary.
- Adding an import cycle and then resolving it with a vague shared package that mixes all domains.
- Rewriting passing tests instead of preserving behavior during a refactor.
- Hiding all errors behind a generic internal-server-error branch.
- Adding a global error middleware before handlers have an intentional error-return design.
- Adding dependency injection machinery when ordinary constructor parameters and main wiring are enough.

### Audit prompt

> Review this generated Session 6 Go refactor before I accept it. Draw the package/import graph and trace POST /entries from route through middleware, handler, application behavior, store, and response. Check package direction, interface justification, concrete wiring, context propagation, explicit actor flow, error categorization, behavior preservation, tests, and over-engineering. Cite exact code. Classify findings as must-fix, safe simplification, or design decision requiring a later session. Do not write replacement code unless I ask.

### Simplification prompt

> Review this TraceLog architecture diff only for unnecessary complexity. Identify each package, interface, constructor, factory, mock, or middleware that has no present caller, substitution, test, or cross-cutting responsibility. For each, say what to delete or collapse and the smallest concrete replacement. Do not flag justified input validation, authorization, error handling, or context propagation as complexity.

---

## 10. Checkpoint: classify, trace, and justify

Complete this without notes first.

1. Draw the target package graph and explain why entries must not import httpapi.
2. Trace POST /entries across middleware, handler, service, store, and response.
3. Classify each responsibility:

| Responsibility | Handler | Service | Store | Middleware | Context | Main |
| --- | --- | --- | --- | --- | --- | --- |
| Decode JSON | | | | | | |
| Assign Entry owner | | | | | | |
| Write map or future SQL state | | | | | | |
| Verify test token | | | | | | |
| Carry request cancellation | | | | | | |
| Map not found to HTTP behavior | | | | | | |
| Construct concrete dependencies | | | | | | |

4. Explain the difference between an interface and a concrete implementation in TraceLog.
5. Explain why the Store interface is consumer-owned and narrow.
6. Explain why request context must reach the future PostgreSQL query even though current memory storage is fast.
7. Name three things that must not be placed in context.
8. Add a service test proving an actor owns a created Entry.
9. Add a handler test proving a wrapped store failure does not leak internal error text.
10. Review an agent proposal to add a generic repository base class. Accept, reject, or defer it with a concrete reason.

### Passing standard

You are ready for Session 7 when you can:

- Trace one request through every package without calling them vague layers.
- Explain why each current package and interface exists.
- Keep HTTP, behavior, storage, and wiring concerns separate.
- Pass cancellation context and an explicit actor correctly.
- Recognize an import cycle and choose a small fix.
- Preserve behavior with tests while refactoring.
- Reject abstraction added only for future possibilities.

If an answer is weak, repeat a narrow request trace or one refactor step. Do not add more layers to make the architecture sound clearer.

---

## 11. Update your progress log

~~~markdown
## Session 6 — Application boundaries

- Date and actual time:
- Diagnostic answer I changed:
- My package graph:
- TraceLog folder or commit:
- Interface justification:
- Request-context flow:
- Chosen bug, violated boundary, root cause, and regression check:
- AI architecture finding:
- One abstraction I rejected and why:
- Confidence (1–5):
- Question to carry into Session 7:
~~~

---

## Ready for Session 7

[Session 7 — Complete REST API Design](session-07-rest-api-design.md) studies the client-facing
contract. With the current boundaries, you should be able to change an API contract without confusing
route behavior, application policy, and storage.

Bring your TraceLog version 2.1 code and prepare to answer:

- Which API behavior is a public contract?
- Which response distinctions matter to a client?
- How will pagination and stable ordering affect handler, service, and store responsibilities?
- Which error categories should become explicit API behavior?
- Where should versioning, filtering, and resource semantics be decided?

> Teach Session 7 only after I show you the Session 6 checkpoint answers and can draw the POST /entries call chain from memory.
