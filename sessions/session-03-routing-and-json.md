# Session 3 — Routing and JSON Boundaries

[← Course overview](../study-plan.md) · [← Session 1](session-01-backend-foundations.md) · [← Session 2](session-02-http-server-flow.md)

## Session contract

**Time:** 100–120 minutes  
**Lecture sequence:** L06 — Routing; L07 — Serialization and deserialization  
**Project:** TraceLog version 1 — a small in-memory JSON API  
**Goal:** Learn to distinguish how a request reaches behavior from how bytes become Go values and Go values become response bytes.

Session 2 established HTTP messages and a single static endpoint. This session adds collection and resource routes plus JSON request and response bodies. We will make the boundaries explicit without prematurely introducing database layers, authentication, validation frameworks, or third-party routers.

### What we are not building yet

Do not add a routing framework, database, repository interface, controller/service hierarchy, authentication, authorization, full validation system, middleware framework, OpenAPI generator, pagination system, or asynchronous work. We need to see route selection and JSON decoding directly first.

---

## Source material

Read the routing lecture before the serialization lecture. Use the transcripts only to resolve a concept you cannot explain from the notes.

| Order | Primary material | What to focus on |
| --- | --- | --- |
| 1 | [L06 — Routing](../resources/backend-from-first-principles/notes/06-6-what-is-routing-in-backend-how-requests-find-their-way-home.md) · [transcript](../resources/backend-from-first-principles/transcripts/06-SubuU1iOC2s.en.vtt) | Method and path as route intent; static and dynamic routes; path parameters; query parameters; nested routes; versioning. |
| 2 | [L07 — Serialization and deserialization](../resources/backend-from-first-principles/notes/07-7-serialization-and-deserialization-for-backend-engineers.md) · [transcript](../resources/backend-from-first-principles/transcripts/07-vzg90tY3uM0.en.vtt) | Why network data needs a shared representation; JSON shape; converting between JSON bytes and Go values. |
| 3 | *Let’s Go* §§2.3 and 2.9 | Supporting Go examples for handlers and JSON. Keep the lecture sequence as the curriculum. |

---

## Outcomes

By the end of this session, you should be able to:

- Explain that routing selects a handler based on request method and path; it does not authenticate, validate, or persist data.
- Recognize static routes, dynamic path parameters, query parameters, nested routes, and API version prefixes.
- Trace GET /entries/42 from request path to path-parameter parsing to a response.
- Explain serialization as Go values becoming JSON bytes and deserialization as JSON bytes becoming Go values.
- Define a small request type and response type that make an API contract explicit.
- Recognize why JSON fields should be deliberately named and why API types should not casually expose internal fields.
- Handle malformed JSON as a client error without continuing into storage or a success response.
- Use handler tests to verify route selection, JSON body decoding, response status, headers, and JSON shape.
- Audit generated code for loose route matching, ignored decode errors, unbounded request bodies, incorrect response writes, and needless router abstractions.

You do not need to memorize every JSON option, URL-escaping rule, content-negotiation header, route-library API, or protocol format. You should know the data flow and where to read documentation before changing boundary behavior.

---

## Active practice and tool lab

### Source-aware tool register

| Tool or technology | Source and role | Depth for this session | Proof of competence |
| --- | --- | --- | --- |
| net/http, encoding/json, and net/http/httptest | Course-supporting Go standard-library tools for the routing and JSON work shown in L06–L07. | Hands-on. | Read the relevant package docs, write a focused route or decode test, and explain the observed response. |
| Browser fetch and DevTools | L06 uses a browser/React client to demonstrate backend requests. | Hands-on for request inspection. | Send a POST with valid JSON and one malformed body; compare what the browser shows with handler tests. |
| React, Angular, and Vue | L06–L07 name client ecosystems as examples. | Recognition only. | Explain that the backend contract is independent of which client calls it; do not start a frontend project. |
| gRPC, WebSockets, XML, and YAML | L07 contrasts other protocols and representations. | Recognition only. | State when JSON-over-HTTP is the chosen contract here and defer other transports until the later lectures. |

No router library is required. First inspect the installed Go documentation for ServeMux. If its
documented patterns fit the route table, use it; otherwise use the tiny explicit dispatcher already
allowed by this lesson.

### Manual-first rule

The generated baseline is only a reading target. You must manually add GET /entries/{id}, including
the happy-path, missing-entry, and malformed-id tests. Start by writing the route table and expected
JSON contract. Do not ask for code until you have explained where id text enters, where it becomes an
integer, and where the response becomes bytes.

### Tool drill

Run the following before changing route behavior:

~~~text
go doc net/http.ServeMux
go doc encoding/json.Decoder
go test ./...
~~~

Then use curl, Postman, or browser fetch to send one valid create request and one malformed JSON
request. Inspect the status, Content-Type, and raw body. Prove with an httptest case that malformed
JSON did not change the in-memory entry count. The manual client is an observation; the test is the
regression proof.

### Completion evidence

Explain every line of the route you added, including route registration, path extraction, parsing,
store lookup, error branch, JSON encoding, and the tests that prove it. Record one documentation
fact you verified rather than guessed.

---

## Agenda

| Time | Activity | Evidence you produce |
| --- | --- | --- |
| 0–10 min | Diagnostic and Session 2 recall | Short written answers |
| 10–30 min | L06 route map | A route table and flow diagram |
| 30–45 min | L07 JSON boundary map | A request and response contract |
| 45–65 min | Inspect TraceLog version 1 | An end-to-end code trace |
| 65–82 min | Manual route and JSON exercise | A test-first feature change |
| 82–97 min | Intentional bug investigation | A failing test and smallest fix |
| 97–112 min | AI-generated code audit | Evidence-based findings |
| 112–120 min | Checkpoint and progress log | Answers and Session 4 readiness |

---

## 1. Start with a diagnostic

Answer without searching first.

1. What is the difference between GET /entries and GET /entries/42?
2. Which part of a request normally chooses the handler: JSON body, method plus path, or response status?
3. What type is a path parameter before your handler parses it?
4. What is the difference between a path parameter and a query parameter?
5. What is the difference between deserializing JSON and validating business rules?
6. Why can valid JSON still be unsafe or unacceptable input?
7. Why is writing a Go struct directly to an HTTP connection not a portable API contract?
8. If JSON decoding fails, which later actions must not happen?
9. What could go wrong if a handler matches every path that merely starts with /entries?

Record the answers. Your uncertainty tells us where to slow down.

### Start-of-session prompt for Codex

> Teach Session 3 from sessions/session-03-routing-and-json.md. Use L06 then L07 in that order. Ask my diagnostic questions one at a time and make me distinguish route selection, JSON decoding, and semantic validation. Build only the small standard-library TraceLog JSON API described here. Do not add a router library, database, authentication, layers, or a validation framework.

---

## 2. Routing: how a request finds behavior

Routing is a mapping from an HTTP request target to a handler. At a minimum, the decision includes method and path.

~~~text
request method + path
  -> route pattern matches
  -> route extracts path values when applicable
  -> selected handler runs
  -> handler decides response
~~~

For TraceLog, the target API surface at the end of this session is:

| Method | Route | Meaning | Source of client data |
| --- | --- | --- | --- |
| GET | /healthz | Service liveness check | No application data |
| GET | /entries | List entries | Optional future query parameters |
| POST | /entries | Create an entry | JSON request body |
| GET | /entries/{id} | Retrieve one entry | id path parameter |

### Route vocabulary

| Term | Example | What you must understand |
| --- | --- | --- |
| Static route | /healthz | The path text is fixed. |
| Collection route | /entries | Represents a set of resources. |
| Resource route | /entries/42 | Represents one identified item. |
| Dynamic path parameter | {id} | A variable segment extracted from the path; it starts as text. |
| Query parameter | /entries?limit=20 | Optional modifier or filter supplied after the question mark. |
| Nested route | /users/42/entries | Shows a relationship or hierarchy in the resource path. |
| Version prefix | /api/v1/entries | Lets an API evolve while old clients transition. |

### The key separation

~~~text
routing answers: “Which handler should receive this request?”
deserialization answers: “How do these bytes become values?”
validation answers: “Is this decoded value acceptable?”
authentication answers: “Who is the caller?”
authorization answers: “May that caller perform this action?”
storage answers: “How is state read or changed?”
~~~

Do not allow a generated codebase to blur these questions into one anonymous helper. The functions can be small today, but their jobs should remain visible.

### Dynamic route exercise

For each request, predict the selected route, path values, and intended outcome.

| Request | Selected route | Path values | Outcome to expect |
| --- | --- | --- | --- |
| GET /entries | | | |
| POST /entries | | | |
| GET /entries/42 | | | |
| GET /entries/not-a-number | | | |
| DELETE /entries/42 | | | |
| GET /entries?limit=10 | | | |
| GET /entries-extra | | | |

Do not assume a router behaves a certain way. Check the local standard-library documentation or write a small handler test for the exact pattern you intend to use.

### Standard-library decision

Before implementation, ask Codex to inspect your installed Go version and the relevant net/http documentation. If its ServeMux supports the exact method-and-path patterns you need, use that small standard-library form. If it does not, use a short explicit dispatcher for this learning exercise rather than installing a router merely to avoid understanding matching.

The purpose is not to memorize route syntax. The purpose is to prove that the route contract matches requests and rejects near-misses.

---

## 3. Serialization and deserialization: the data boundary

Go structs and JSON documents are not the same thing. A Go program works with typed values. A client sends and receives bytes in a standardized representation.

~~~text
client JSON bytes
  -> request body
  -> JSON decoder
  -> Go request value
  -> application behavior
  -> Go response value
  -> JSON encoder
  -> response bytes
  -> client
~~~

### Example API shapes

A client can send:

~~~json
{
  "title": "Trace an HTTP request"
}
~~~

A successful creation can return:

~~~json
{
  "id": 1,
  "title": "Trace an HTTP request"
}
~~~

A client should not need to know how your in-memory map, future database row, lock, or internal error is represented. The JSON contract is deliberately smaller than the whole program.

### Request and response types

Use separate small types when the transport contract differs from the internal model. That is not premature architecture; it makes the API boundary explicit and prevents accidental field exposure.

~~~go
type createEntryInput struct {
    Title string `json:"title"`
}

type entryResponse struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
}
~~~

Questions to answer:

1. Which fields in this example are available to encoding/json?
2. Why are the field names capitalized in Go?
3. What API spelling does each JSON tag define?
4. Why might an internal field such as an owner identifier or database-only metadata not belong in entryResponse?
5. Why should an input type not be assumed valid merely because decoding succeeded?

### Boundary rules for this session

For a JSON-accepting handler:

1. Treat the request body as untrusted.
2. Define a maximum body size appropriate for this learning API.
3. Decode into a dedicated input type.
4. Reject malformed JSON and type mismatches with a client error.
5. For a strict API contract, reject unexpected fields and trailing JSON values; look up the exact decoder pattern before relying on it.
6. Do not create or mutate an Entry if decoding fails.
7. Set the JSON Content-Type before encoding a response.
8. Do not attempt to send a second error response after response output has started.

Formal semantic validation and normalization are the focus of Session 5. This session preserves the Session 1 invariant that an Entry title cannot be blank, but it does not build a broad validation framework.

### Understand versus recognize

| Topic | Understand deeply enough to reason about | Recognize and look up when needed |
| --- | --- | --- |
| Route matching | Method and exact path select behavior; near-miss paths require tests. | Every route-pattern grammar and router-library API. |
| Path values | They arrive as text and must be parsed and checked. | Escaping edge cases and complex path encodings. |
| Query values | They modify or filter a request and are still untrusted text. | Pagination cursor formats and caching interaction. |
| JSON decoder | Converts bytes to typed Go values and can fail. | Every decoder option and custom decoder implementation. |
| JSON encoder | Converts a response value to bytes; headers go first. | Streaming and custom marshaling details. |
| JSON tags | Explicitly state the external field names. | All tag options such as omission behavior. |
| Input versus response type | Keeps incoming and outgoing API contracts deliberate. | Schema generators and versioned DTO frameworks. |
| In-memory store | Holds test data for this session; a small lock prevents unsafe concurrent mutation. | Lock design and concurrent performance, taught later. |
| Error response | A failed parse must end the success path. | A complete production error-envelope specification. |

---

## 4. Inspect TraceLog version 1

Create a new directory named tracelog-session-03 by evolving the Session 2 program.

~~~text
tracelog-session-03/
  go.mod
  main.go
  server.go
  entry.go
  server_test.go
~~~

### Required behavior at the end of the session

- GET /healthz still behaves as in Session 2.
- GET /entries returns a JSON array.
- POST /entries accepts a small JSON object with title and returns 201 plus a JSON representation.
- GET /entries/{id} returns the matching JSON representation or 404 when absent.
- A malformed JSON request returns a client error and does not create state.
- A wrong method produces a useful method-related response.
- JSON responses state an application/json Content-Type.
- Handler tests cover routes through the actual returned handler.

### Deliberate design constraints

- Use net/http, encoding/json, net/http/httptest, and other standard-library packages only.
- Keep one concrete in-memory store. It may use a small mutex because HTTP handlers can run concurrently, but do not design a storage interface today.
- Keep package main and a small file layout. Package architecture comes later.
- Use dedicated input and response types only where they clarify the JSON contract.
- Use explicit parsing for the id path value and choose a documented response for an invalid id.
- Do not return raw internal errors to clients.
- Do not add a database, authentication, generic response system, reflection-based mapper, or router framework.
- Keep route registration and handler behavior easy to trace in one screen at a time.

### The full request trace to annotate

Trace a valid POST /entries before you run it:

~~~text
1. Client sends method, path, headers, and JSON bytes.
2. The mux selects the collection creation handler.
3. The handler enforces its HTTP and body-size assumptions.
4. JSON bytes decode into createEntryInput.
5. The existing title invariant is checked.
6. The concrete in-memory store assigns an id and stores the Entry.
7. The handler maps Entry to entryResponse.
8. Header and status are chosen.
9. The encoder writes JSON response bytes.
10. The client receives 201 and the representation.
~~~

For each step, mark:

- Input or output value
- Trusted or untrusted data
- Possible error
- Function responsible
- Test that would prove it

### Code-reading task

Before running the server, answer:

1. Where are routes registered?
2. What exact conditions select the GET /entries/{id} behavior?
3. How does the handler obtain the id text?
4. Where is the text converted into an integer?
5. What happens when conversion fails?
6. What code distinguishes a malformed JSON document from a missing entry?
7. Where does a Go Entry become an API response type?
8. Which write commits response headers and status?
9. Does any error path continue into storage or success encoding?
10. Why is the in-memory store protected or deliberately constrained?
11. Which parts should later move behind application or storage boundaries, and which should remain HTTP-specific?

### Prompt to generate the inspection code

> Create a new tracelog-session-03 directory by evolving my Session 2 health server into the smallest standard-library JSON API. Preserve GET /healthz. Add GET /entries, POST /entries, and a route pattern for GET /entries/{id}; use the installed Go version’s documented net/http routing support, or a tiny explicit dispatcher if necessary. Use a concrete in-memory store, small JSON input/output types, bounded request-body handling, clear JSON error behavior, and httptest tests. Before writing code, show the route table, JSON contracts, data invariants, error cases, file tree, and tests. Do not add a third-party router, database, authentication, interfaces, controller/service/repository layers, or generic response abstractions. Wait for my review before making changes.

---

## 5. Manual coding and test-first feature work

Start with a baseline that supports GET /entries and POST /entries. Then add GET /entries/{id} yourself.

### Feature contract: retrieve one entry

| Input | Expected result |
| --- | --- |
| GET /entries/1 where id 1 exists | 200 and one JSON entry |
| GET /entries/999 where no such entry exists | 404 and a safe error response |
| GET /entries/not-a-number | Your documented client-error choice, with a test |
| POST /entries/1 | Method-related error; it must not create anything |
| GET /entries-extra | It must not accidentally use the collection handler |

### Implementation order

1. Read the existing collection route tests.
2. Write a failing success test for GET /entries/{id}.
3. Write a failing not-found test.
4. Decide and test the malformed-id contract.
5. Register the route and parse the path value.
6. Ask the store for the entry.
7. Map it to a response type and encode JSON.
8. Format and run focused tests.
9. Run all tests.
10. Explain the entire path out loud without reading the code.

### JSON contract exercise

Write a test that sends a valid JSON create request and verifies all of the following:

- The response status is 201.
- The Content-Type describes JSON.
- The response decodes as JSON.
- The returned id is usable.
- The returned title matches the input.
- A later GET /entries/{id} returns the same entry.

Then write one test that sends malformed JSON. Prove that the in-memory entry count did not change.

### Read and explain exercise

For the following request, identify every boundary and transformation:

~~~text
POST /entries
Content-Type: application/json

{"title":"  Learn routes  "}
~~~

Answer:

1. What arrives as bytes?
2. What becomes a Go string?
3. What is structural decoding versus semantic validation?
4. At what later session should whitespace normalization become explicit?
5. What should the client receive if the JSON syntax is malformed?
6. What should the client receive if the JSON is syntactically valid but title is blank?
7. Why should the server not return the entire internal Entry just because it is convenient?

### Targeted implementation prompt

> Guide me through adding GET /entries/{id} to my Session 3 TraceLog API. Do not write code first. Ask me to define the route contract and tests. Then help me make one small change at a time: route registration, id parsing, not-found behavior, JSON response, and tests. Require me to explain each response status and stop if I propose a router, interface, or database without a current need.

---

## 6. Intentional bug lab

Choose one bug after the happy-path tests pass. Do not reveal its location to yourself in advance.

| Bug | Faulty behavior | Observable symptom |
| --- | --- | --- |
| Ignored decode error | The handler continues after JSON decoding fails. | A zero-value or partial entry is stored, or the handler reports false success. |
| Loose route match | A prefix test matches /entries-extra as if it were /entries. | An invalid path reaches the wrong behavior. |
| Parse fallback | An id parse error quietly becomes zero. | GET /entries/not-a-number can retrieve or report the wrong resource. |
| Late JSON header | The handler sets Content-Type after encoding starts. | Clients receive an incorrect or default content type. |
| Double response | The handler writes a JSON error and then encodes a success response. | Confusing body, warnings, or an untrustworthy contract. |
| Internal-field leak | The response serializes an internal-only field. | The public API exposes data it should not promise or reveal. |
| Unbounded body | The handler reads arbitrary request size without a limit. | A client can force disproportionate memory or work. |

### Debugging protocol

1. Reproduce with a focused httptest test before using a browser or manual client.
2. State the visible contract failure: wrong route, wrong status, wrong JSON, leaked field, state change, or resource risk.
3. Trace input bytes through route selection, parsing, decoding, storage, and encoding.
4. Identify the earliest point where the invariant was violated.
5. Add a regression test that would fail before the fix.
6. Make the smallest root-cause correction at the shared decision point.
7. Re-run formatting, focused tests, all tests, and a manual request if helpful.
8. Explain why the fix does not need a new framework or rewrite.

### Investigation questions

- Which request exactly triggers the bug?
- Which route did the server select?
- Which input values were raw text or bytes?
- Did decoding report an error, and what did the caller do with it?
- Did state change after a request that should fail?
- What is the first response write?
- Is the behavior caused by route matching, data conversion, validation, storage, or response encoding?

### Bug-lab prompt

> I intentionally broke my Session 3 TraceLog API. Be my debugging partner. First ask for the failing httptest case or exact request/response pair. Then guide me to trace routing, path parsing, JSON decoding, state mutation, and response writing in order. Make me state the violated contract and add a regression test. Do not rewrite the API or reveal a patch until I have a reasoned root-cause hypothesis.

---

## 7. Review AI-generated route and JSON code

An agent can make this API look plausible while violating its contract. Review generated changes before accepting them.

### Audit rubric

| Review question | Evidence of a good answer |
| --- | --- |
| Are method and path matched deliberately? | Each expected route works; near-misses and wrong methods have tests. |
| Are path parameters parsed safely? | Parse errors do not silently turn into a default identifier. |
| Is the body bounded and decoded explicitly? | Untrusted body input has a known size limit and decode errors stop the path. |
| Does JSON represent a deliberate public contract? | Input and response fields are explicit; internal fields are not leaked. |
| Are statuses and headers correct? | JSON success and error paths set headers before output and use meaningful statuses. |
| Can malformed input mutate state? | Tests prove it cannot. |
| Is shared in-memory state handled honestly? | A small guard exists or the limitation is explicitly constrained. |
| Is the design proportionate? | Standard library, concrete store, no speculative layers or router framework. |
| Are tests end-to-end at the handler boundary? | httptest exercises registered routes, status, headers, and decoded response bodies. |

### Common AI-generated-code failure modes

- Checking the wrong path with a broad prefix or string contains test.
- Assuming a dynamic route variable is already an integer.
- Ignoring the error from JSON decoding.
- Accepting arbitrary body size.
- Using the database/internal model as the public JSON response by default.
- Returning the same generic error for malformed JSON, missing resources, and unexpected server failures.
- Writing a response error and continuing to success code.
- Adding a router or reflection-based serializer before the standard library is understood.
- Creating an interface for the only in-memory store.
- Testing helpers while bypassing the handler and route registration that users actually invoke.

### Audit prompt

> Review this generated Session 3 Go diff before I accept it. Trace each supported request from method and path through route matching, path/query parsing, JSON decode or encode, state change, and response. Find incorrect route matches, ignored or unsafe input handling, leaked fields, status/header mistakes, missing tests, and unnecessary abstractions. Cite exact code. Separate must-fix defects from optional simplifications. Do not generate a replacement implementation unless I ask.

### Documentation comparison prompt

> I am using this net/http or encoding/json behavior: [paste the route pattern, PathValue use, decoder setting, or encoder call]. Check the relevant official Go documentation or local package docs. Explain what is guaranteed, what my code assumes, and one focused httptest case that verifies the assumption.

---

## 8. Checkpoint: route, decode, and explain

Complete this without notes first.

1. Trace a valid POST /entries from raw HTTP request through JSON response in order.
2. Explain the difference among a static route, dynamic route parameter, and query parameter.
3. Explain why a route parameter must be parsed even when it looks numeric in the URL.
4. Explain the difference between malformed JSON and a decoded but invalid title.
5. Predict the appropriate outcome for each request:

| Request | Expected category of result |
| --- | --- |
| GET /entries | |
| POST /entries with valid JSON | |
| POST /entries with malformed JSON | |
| GET /entries/999 | |
| GET /entries/not-a-number | |
| POST /entries/1 | |
| GET /entries-extra | |

6. Add a test that proves malformed JSON does not create an Entry.
7. Review an agent proposal to install a router library. Accept, reject, or defer it based on the codebase’s actual current needs.
8. Explain why entryResponse can be safer than encoding the internal Entry directly.

### Passing standard

You are ready for Session 4 when you can:

- Trace routing and JSON conversion without conflating them.
- Add a resource retrieval route with tests before implementation.
- Diagnose one input-boundary bug from evidence.
- Explain which data is untrusted at every step of a request.
- State why the API currently has no authentication and why that is a security boundary to add next.
- Name one encoding/json or route-pattern detail you would verify in documentation rather than guess.

If an answer is weak, repeat the one relevant experiment or test. Do not restart both lectures.

---

## 9. Update your progress log

~~~markdown
## Session 3 — Routing and JSON

- Date and actual time:
- Diagnostic answer I changed:
- Route table I can explain:
- TraceLog folder or commit:
- JSON request and response contract:
- Chosen bug, symptom, root cause, and regression test:
- Route or JSON detail I looked up:
- AI review finding:
- Confidence (1–5):
- Question to carry into Session 4:
~~~

---

## Ready for Session 4

The next lecture adds identity and permission checks. Before authentication or authorization can be correct, keep this Session 3 chain clear:

~~~text
route -> deserialize -> validate -> identify caller -> authorize action -> apply behavior -> serialize response
~~~

Bring your TraceLog version 1 code and prepare to answer:

- Which values came from an untrusted caller?
- Where could an authenticated identity be attached to a request?
- Why must authorization check both the caller and the specific requested resource?

> Teach Session 4 only after I show you the Session 3 checkpoint answers and can trace a malformed JSON request without confusing it with routing or validation.
