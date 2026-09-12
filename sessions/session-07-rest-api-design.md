# Session 7 — Complete REST API Design

[← Course overview](../study-plan.md) · [← Session 6](session-06-application-boundaries.md)

## Session contract

**Time:** 110–120 minutes  
**Lecture sequence:** L11 — Complete REST API Design  
**Project:** TraceLog version 2.2 — a contract-first, owner-scoped API with a bounded list response  
**Goal:** Treat an API as a durable public interface. Design and test a consistent resource contract
before changing Go code, then implement one small vertical feature that proves the contract can be
enforced through the existing handler, service, and store boundaries.

Session 6 gave TraceLog clear responsibility boundaries. This session gives those boundaries a
client-facing promise:

~~~text
client intent
  -> method + path + query + JSON request contract
  -> authentication and authorization
  -> handler translates HTTP
  -> service applies behavior
  -> store returns a stable result
  -> status + headers + JSON response contract
~~~

An API is not its routes alone, and it is not its Go implementation. It is the agreement a client
can rely on: accepted inputs, authentication and ownership rules, side effects, response shape,
status behavior, ordering, defaults, and error behavior.

### What we are not building yet

Do not add:

- PostgreSQL, SQL pagination, migrations, caching, queues, search, file storage, or real-time work.
- A router library, REST framework, generic CRUD framework, API gateway, or a package called common
  that hides contract decisions.
- Every possible CRUD endpoint merely because an Entry is a resource.
- A public version prefix in the running local app just to make it look production-ready.
- An OpenAPI generator, Swagger server, client SDK generator, or new Go dependency.
- Cursor pagination, arbitrary sort-field support, rich query-language parsing, or a reusable
  pagination framework.
- A generic response envelope for every endpoint.

We will design a larger API surface on paper, but implement only one narrow behavior: owner-scoped,
bounded, deterministic pagination for GET /entries. Session 8 will retain that contract while moving
the storage boundary to PostgreSQL.

---

## Source material

Read L11 as a connected lesson. It deliberately designs an interface before choosing language- or
framework-specific implementation.

| Primary material | What to focus on |
| --- | --- |
| [L11 — Complete REST API Design](../resources/backend-from-first-principles/notes/11-11-complete-rest-api-design.md) | REST-oriented interface design; URLs and resources; CRUD versus custom actions; pagination, filtering, sorting, defaults; response semantics; consistency; documentation. |
| [L11 transcript](../resources/backend-from-first-principles/transcripts/11-RG6q57DwV8Y.en.vtt) | Client/server and stateless ideas; URL anatomy; resource pluralization and dynamic IDs; version path examples; Insomnia request design; list pagination/filtering/sorting; PATCH, DELETE, archive/clone actions; OpenAPI/Swagger advice. |
| *Let’s Go* §§2.4–2.5, 7.1–7.2 | Supporting Go examples for request handling, response headers/statuses, templates/forms as contrasts, and testing. The book is not the API-design sequence. |

### Lecture concepts to carry forward

L11 uses these conventions and examples:

- A URL has scheme, authority/domain, path, query parameters, and sometimes a browser fragment.
  The fragment is client-side navigation, not application input for the server.
- An API may use an API subdomain and a path version such as v1. Those are compatibility choices,
  not a substitute for thinking through a contract.
- Resource path segments are conventionally lowercase plural nouns. A collection and a member still
  use the plural resource name: /entries and /entries/{id}.
- Query parameters describe a collection request, such as pagination, filtering, and sorting.
- Create, read, update, and delete follow recurring shapes. A server-side operation such as archive
  or clone can be a custom action rather than a disguised field update.
- Status describes what actually happened. POST create commonly returns 201; a custom POST action
  can return 200, or 201 if it creates a new resource.
- Documentation and repeatable API-client requests are part of the API experience, not afterthoughts.

---

## Outcomes

By the end of this session, you should be able to:

- Define an API contract as method, path, query, request body, authentication/authorization policy,
  side effect, response status/headers/body, and compatibility promise.
- Distinguish a collection, a member resource, a relationship, and a custom action.
- Design resource-oriented routes without leaking Go package names, database tables, or UI button
  names into the path.
- Explain why /entries/{id} is a member route while /entries?tag=go is a collection query.
- Recognize the L11 conventions for plural resource names, lower-case path segments, readable IDs or
  slugs, and a version path when a published compatibility policy needs one.
- Distinguish POST create, GET collection/member, PATCH partial update, PUT replacement, DELETE, and
  a custom action endpoint well enough to read and review Go handlers.
- Explain why an empty successful list is 200 with an empty array, while a missing requested member
  is normally 404.
- Design bounded page-and-limit pagination with defaults, a maximum limit, deterministic ordering,
  useful metadata, and tests.
- Explain why pagination, filtering, and sorting are API and performance decisions, not frontend
  decoration.
- Preserve owner authorization through a new list feature; query parameters never bypass ownership.
- Recognize a stable API error contract and avoid leaking implementation errors.
- Write a small API contract document and a minimal OpenAPI fragment without depending on a code
  generator.
- Use Insomnia, Postman, or curl to exercise the same contract that handler tests prove.
- Audit AI-generated API changes for ambiguous routes, wrong status semantics, broken pagination,
  accidental breaking changes, authorization bypasses, documentation drift, and needless frameworks.

You do not need to memorize every HTTP status, every OpenAPI keyword, every formal REST constraint,
every pagination strategy, or a universal rule for nested routes and versioning. You must understand
the client-visible behavior, be able to trace it into code, and know when to consult documentation
or write a focused test.

### Understand deeply versus recognize and look up

| Learn deeply now | Recognize and look up when needed |
| --- | --- |
| Resource versus action; collection versus member; method, path, input, output, and error contracts; stable list behavior; page/limit bounds; authorization after query parsing; status differences for success, bad input, missing member, and empty collection; backwards-compatibility thinking. | Exact RFC wording for safe/idempotent methods; complete HTTP status catalog; OpenAPI schema grammar; Swagger UI deployment; client SDK generation; cursor pagination; complex filters; sparse fieldsets; HATEOAS; GraphQL and RPC APIs. |

---

## Active practice and tool lab

### Source-aware tool register

| Tool or technology | Source and role | Depth for this session | Proof of competence |
| --- | --- | --- | --- |
| Insomnia | Demonstrated in L11 as a lightweight API client and interface-design workspace. | Hands-on if available. | Create a TraceLog collection with folders, a local base-URL environment, request descriptions, and saved safe examples for list behavior. |
| Postman | L11 names it as the better-known alternative to Insomnia. | Working familiarity. | Explain that either client can make and inspect the HTTP contract; use one rather than maintaining duplicate collections. |
| OpenAPI and Swagger tools | L11 recommends maintained interactive documentation and a playground. | Hands-on contract authoring; tool UI familiarity. | Write a small OpenAPI fragment by hand and, if Swagger Editor/UI is already available, load it and compare it against tests. Do not generate server code. |
| curl | Course-supporting raw HTTP client. | Hands-on. | Inspect a paginated response and an invalid-query response with curl -i. |
| net/http and net/http/httptest | Course-supporting standard-library implementation and proof tools. | Hands-on. | Add handler tests that make the pagination contract executable. |
| RPC and GraphQL | Mentioned in L11 as other API styles. | Recognition only. | Explain that this course chooses JSON-over-HTTP REST-style endpoints for TraceLog; do not add a second API style. |

No API-client collection, OpenAPI document, or browser tool replaces tests. The collection is a
human-facing contract and reproducible manual check. Handler and service tests are the executable
evidence that the server keeps its promise.

### Manual-first rule

Before an agent writes a pagination line, you must:

1. Write a route table and a response example in your own words.
2. State the default page, default limit, maximum limit, order, owner scope, and invalid-input rules.
3. Write the first failing handler test for GET /entries?page=1&limit=2.
4. Explain where query text is parsed, where owner policy is applied, where the slice/page is
   selected, and where JSON becomes the response.

You will manually implement one vertical slice. An agent may help inspect the existing code or review
your test plan, but it may not fill in the entire feature before you make an attempt.

### Tool drill

1. Create one API-client collection named TraceLog API. Use an environment variable for the local
   base URL only; do not save credentials or raw bearer tokens.
2. Add requests for the current health endpoint, list entries with defaults, list entries with
   page/limit, a page with no matching entries, and invalid page/limit input.
3. Write a short description beside each request: intent, expected status, expected response shape,
   and whether it can mutate state.
4. Use curl -i for one successful list and one invalid query. Compare its raw status/headers/body
   with what the client UI shows.
5. Write a minimal OpenAPI fragment for GET /entries after the route tests pass. If a Swagger UI or
   editor is available, use it only to inspect the contract; do not install a generator or accept
   generated Go code.

### Completion evidence

Finish the tool lab with:

- one organized collection or an equivalent curl request log;
- one raw response annotated with status, Content-Type, and response fields;
- a small API contract document or OpenAPI fragment;
- the exact test names that prove the list contract; and
- a sentence explaining why an empty list is not the same as a missing Entry.

---

## Agenda

| Time | Activity | Evidence you produce |
| --- | --- | --- |
| 0–12 min | Diagnostic and Session 6 recall | Written predictions about routes and statuses |
| 12–28 min | L11 synthesis | A resource map and URL anatomy annotation |
| 28–43 min | Contract-first design | A TraceLog route table and example list response |
| 43–58 min | Inspect TraceLog version 2.1 | A request-to-store code trace |
| 58–82 min | Manual pagination feature | Failing tests, small implementation, passing tests |
| 82–95 min | API-client and documentation drill | Saved requests and a contract fragment |
| 95–108 min | Intentional bug investigation | Root cause, regression test, smallest fix |
| 108–118 min | AI-generated API audit | Evidence-based review |
| 118–120 min | Checkpoint and progress log | Readiness decision for PostgreSQL |

If you have only 90 minutes, complete through the manual feature and its tests. Do the OpenAPI and
AI-review work before Session 8; do not skip the contract or bug exercise.

---

## 1. Start with a diagnostic

Answer without searching first. A wrong answer tells us which decision to slow down and test.

1. What client-visible information makes GET /entries an API contract rather than merely a Go
   handler function?
2. Why is GET /entries/42 a different request shape from GET /entries?tag=go?
3. If an authenticated user has no entries, should GET /entries return 200, 404, or 204? Why?
4. If GET /entries/42 requests an Entry that does not exist, which broad response should it receive?
5. What is the difference between POST /entries and POST /entries/{id}/archive?
6. Why might PATCH better communicate a partial Entry update than PUT?
7. Why must a server bound a client-supplied limit?
8. What order must a paginated response promise before a client can reliably request page two?
9. Does a POST method force a 201 response? Give an example.
10. What could break for a client if a raw JSON array becomes a pagination object?
11. Why is an API-client collection useful but insufficient as the test suite?
12. Which implementation detail must never appear in a public route: Go package name, resource name,
    or a member identifier?

Save your answers under “Session 7 diagnostic” in the progress log. Do not correct them yet.

### Start-of-session prompt for Codex

> Teach Session 7 from sessions/session-07-rest-api-design.md. Use L11 as the primary source. Ask the diagnostic questions one at a time and make me state my reasoning before explaining. Require me to write the TraceLog route table, list response example, pagination defaults, maximum limit, ordering rule, and first failing test before any implementation. Make me complete one manual vertical slice for GET /entries pagination. Use the hint ladder; do not generate a full API, router framework, OpenAPI generator, database, or versioned rewrite.

---

## 2. API design means designing an interface

The L11 distinction is essential:

~~~text
interface design:
  What can a client request?
  What does it mean?
  What inputs are accepted?
  What response and side effect can it rely on?

implementation:
  Which Go function, package, store, SQL query, or JSON encoder produces it?
~~~

Implementation can change behind a preserved contract. For example, Session 8 may replace an
in-memory list with a PostgreSQL query. A client should still receive the same documented list
shape, ordering, pagination behavior, authorization behavior, and error semantics unless we make a
deliberate compatibility change.

### The minimum contract checklist

For every endpoint, answer all of these before coding:

| Contract question | Example for GET /entries |
| --- | --- |
| Who can call it? | A verified TraceLog user. |
| What does the client send? | GET method, collection path, optional page and limit query values. |
| What is client-controlled? | Query strings and every header/body value until verified or validated. |
| What resource/behavior does it request? | A page of the caller’s Entries. |
| What defaults and bounds apply? | Page and limit have documented defaults and a fixed maximum. |
| What state may it mutate? | None. |
| What determines the result order? | An explicit stable course order, not map iteration. |
| What happens when there are no matching Entries? | 200 with an empty data array and valid metadata. |
| What can fail? | Bad query input, missing identity, internal store failure. |
| What output is guaranteed? | Documented JSON response, status, Content-Type, and no data belonging to another owner. |
| What tests prove it? | Handler tests for defaults, pages, invalid input, empty page, and owner scope. |

### REST-style mental model, without ceremony

L11 covers client/server separation, stateless requests, uniform resource-oriented interfaces, and
self-descriptive messages. You do not need to prove that TraceLog satisfies every formal REST
constraint. Use the useful operational questions:

- Can a request carry enough information for the server to process it without relying on hidden
  client conversation state?
- Does the method and path make the requested resource or action understandable?
- Do status, headers, and JSON describe the outcome clearly?
- Can a client learn the contract from documentation and examples without reading server internals?

Do not call an endpoint RESTful merely because its URL contains nouns. The behavior and consistency
matter more than a label.

---

## 3. Resource map and route decisions

### TraceLog resource map

Start by naming the things the API exposes, not functions you want the server to execute.

~~~text
TraceLog
  └── entries                         collection owned/scoped by a principal
      └── {id}                        one Entry
          └── archive                 a possible domain action on that Entry
~~~

The slash expresses a relationship in the path. Do not nest arbitrarily just because IDs exist. A
deep path should communicate a real scoping or ownership relationship, not duplicate information the
server already has.

### Target contract map

Design all rows. Implement only the row marked “manual implementation” today.

| Client intent | Method and path | Contract status | Why it has this shape |
| --- | --- | --- | --- |
| Health check | GET /healthz | Existing | A public operational endpoint. |
| List the caller’s Entries | GET /entries?page=&limit=&tag= | **Manual implementation:** page/limit; tag is design-only | Collection query parameters refine a collection request. |
| Create an Entry | POST /entries | Existing | A new Entry is created under the verified principal. |
| Read one Entry | GET /entries/{id} | Existing | The path identifies one member. |
| Partially update an Entry | PATCH /entries/{id} | Design-only | L11 uses PATCH for partial field updates rather than pretending every update replaces the resource. |
| Delete an Entry | DELETE /entries/{id} | Design-only | Deletion is a member operation; success often has no representation. |
| Archive an Entry | POST /entries/{id}/archive | Design-and-test stretch | Archive can trigger domain work beyond a field assignment, so it is a named action. |

Avoid route shapes such as:

~~~text
POST /createEntry
GET /getEntries
POST /entries/{id}/delete
GET /entries/list
POST /entryService/archiveEntry
~~~

They repeat the method’s job, leak implementation/UI names, or make a client learn a new arbitrary
verb convention for each endpoint.

### Methods communicate intent; result determines status

| Operation | Typical method | What you must understand | Do not memorize |
| --- | --- | --- | --- |
| Fetch collection/member | GET | Reads a representation; query parameters can refine a collection. | Every cache and conditional-header rule. |
| Create a resource | POST to collection | Commonly returns 201 with the created representation when creation succeeds. | Every Location-header convention; look it up when implementing one. |
| Partial update | PATCH to member | Represents changing selected fields; define which fields are allowed and how omitted fields behave. | A universal PATCH media type or merge algorithm. |
| Full replacement | PUT to member | Represents replacement semantics when the contract genuinely needs that. | Treating PUT and PATCH as interchangeable. |
| Delete | DELETE to member | A successful deletion can be 204 with no response body. | What a repeated delete must return in every product. |
| Domain action | POST to member/action | Makes a server-side action explicit. Archive may return 200; clone may return 201 if it creates a new resource. | Using POST as a catch-all for every endpoint. |

### Versioning is a compatibility decision

L11 shows route-based versioning such as an API host plus /v1. TraceLog is a local learning API with
no external clients, so retain the existing unversioned paths today. Write this policy in your
contract:

> A published breaking change requires an explicit compatibility decision. We will not add /v1
> merely to predict the future, and we will not silently break a published response shape.

Changing GET /entries from a raw JSON array to a pagination object is a breaking change for a real
consumer. It is acceptable in this controlled exercise only because you will update the contract and
tests together. Treat that fact as practice in compatibility judgment, not permission to change APIs
casually.

---

## 4. Design the list contract before code

### Required API design artifact

Create api-contract.md in your Session 7 project folder. Fill this table yourself before opening a
handler:

| Field | Your decision |
| --- | --- |
| Endpoint | |
| Caller and owner scope | |
| Query parameters and accepted values | |
| Defaults | |
| Maximum limit | |
| Stable ordering | |
| Success status and Content-Type | |
| Success JSON shape | |
| Empty-result behavior | |
| Invalid query behavior | |
| Authentication/authorization behavior | |
| No-state-change guarantee | |
| Tests that prove it | |
| Compatibility impact on earlier clients | |

After your attempt, compare it with this bounded course contract.

### TraceLog’s bounded course contract

~~~text
GET /entries?page=1&limit=20

Authentication:
  required; the result includes only Entries owned by the verified principal.

Query defaults:
  page defaults to 1.
  limit defaults to 20.

Query bounds:
  page must be a positive integer.
  limit must be a positive integer no greater than 100.

Order:
  Entries are in ascending numeric ID order for this in-memory exercise.
  The order is explicit because map iteration is not an API contract.

Success:
  200 and application/json.
  A list-specific pagination object:
  {
    "data": [...],
    "page": 1,
    "limit": 20,
    "total": 42,
    "totalPages": 3
  }

No matching Entries:
  200 with "data": [] and valid pagination metadata.

Bad page or limit:
  Existing safe client-error shape; no storage mutation and no internal error text.
~~~

The use of data is intentional for this list response. Do not turn it into a generic response
wrapper for health checks or every member endpoint. The metadata belongs to a paginated collection,
not to every JSON object in the service.

### Page selection model

For a 1-based page:

~~~text
start = (page - 1) * limit
end   = smaller of start + limit and total
data  = ordered matching entries from start up to end
~~~

You must validate and bound inputs before applying this model. You also must handle start at or past
total as a valid empty page, not as an out-of-range panic or a 404 collection response.

### Filtering and sorting: design now, keep scope small

L11 treats filtering and sorting as ordinary collection-query behavior. For TraceLog:

| Concern | Session 7 decision | Reason |
| --- | --- | --- |
| Tag filter | Design ?tag=go, but defer implementation unless pagination is complete early. | It demonstrates a query contract without turning this session into a search system. |
| Sort field/order | Do not accept arbitrary client fields today. Fix ID-ascending order explicitly. | A fixed order is testable and avoids reflection, SQL injection risk later, and arbitrary parsing. |
| Total count | Return it because L11’s page model uses it. | It helps a client decide whether more pages exist; Session 8 will make its database cost visible. |
| Cursor pagination | Recognition only. | Page-and-limit is enough to learn bounded list contracts before database-scale trade-offs. |

### Error and response rules

Preserve the stable error representation already established by the earlier sessions. The important
distinctions are:

| Situation | Expected broad result | Reason |
| --- | --- | --- |
| Missing/invalid authentication | 401 | No verified caller exists. |
| Caller requests another owner’s member Entry | Existing 403-or-404 policy | Preserve the deliberate Session 4 resource-disclosure policy. |
| page or limit cannot parse or violates the contract | 400 | The collection request is invalid. |
| Collection has zero matches | 200 plus empty data | The collection request succeeded. |
| GET /entries/{id} has no accessible matching member | 404 or the established access policy | A specific member was requested. |
| Wrong method for a route | 405 with useful allowed-method behavior | The path may exist, but the method is not accepted. |
| Unexpected storage failure | Safe 5xx behavior | Clients do not receive Go/map/future SQL internals. |

Do not use 404 as a decorative way to say “your filter had no results.” It changes client control
flow and makes list endpoints unpredictable.

---

## 5. Inspect TraceLog version 2.2

Before the manual feature, inspect a small baseline copied from Session 6. It should retain the
current package direction:

~~~text
tracelog-session-07/
  cmd/tracelog/
    main.go
  internal/httpapi/
    server.go
    entries.go
    entries_test.go
  internal/entries/
    entry.go
    service.go
    service_test.go
  internal/memory/
    store.go
    store_test.go
  api-contract.md
~~~

Your actual files may differ slightly, but the code must still make the request path obvious. Do not
reorganize files merely to match this tree.

### Required baseline behavior

- GET /healthz remains public and unchanged.
- Existing protected create, list, and member-read behavior remains authorized and tested.
- Existing validation and owner-assignment behavior remains unchanged.
- Existing response errors remain safe.
- GET /entries exists before pagination; it can begin as the current simple list shape.
- No database, new HTTP framework, OpenAPI runtime, or generated SDK appears.

### Code-reading exercise

Before running tests, trace the current GET /entries path and write answers:

1. Where is the route registered?
2. Which handler receives query text?
3. Where is the verified actor obtained?
4. Which function decides that the caller can list only their own Entries?
5. Does the store receive HTTP types, raw query strings, or a small application-level list request?
6. Where is response ordering currently determined? If it is map iteration, what contract bug exists?
7. Which response write commits status and headers?
8. What test currently proves user one cannot list user two’s Entry?
9. What would change if a raw array becomes a pagination object?
10. Which pieces must Session 8 preserve even though its store becomes PostgreSQL-backed?

### Prompt to generate the inspection baseline

> Create a new tracelog-session-07 directory by copying my Session 6 TraceLog code without changing its existing HTTP behavior. Preserve GET /healthz, protected Entry creation/list/read behavior, validation, ownership policy, package direction, and tests. Add only an empty api-contract.md template for me to complete; do not implement pagination, PATCH, DELETE, archive, versioned routes, OpenAPI generation, PostgreSQL, a router library, or a generic CRUD layer. First show the preserved-contract list, file tree, current request trace, and tests. Wait for my review before writing files.

---

## 6. Manual coding and test-first feature work

### Required vertical feature

Implement this yourself, with guidance:

> Make GET /entries return the bounded pagination contract using page and limit, while preserving
> owner scope and stable ID-ascending ordering.

This is not a cosmetic response change. It changes what clients receive and exercises the whole
existing path:

~~~text
query text
  -> handler parsing and HTTP error mapping
  -> small application list options
  -> owner-scoped service behavior
  -> deterministic memory-store page
  -> list response mapping and JSON encoding
  -> handler and service tests
~~~

### Test plan: write these before implementation

| Test | What it proves |
| --- | --- |
| Default list request | Omitting page/limit uses documented defaults and returns a pagination object. |
| First page | page=1 and limit=2 returns the first two entries in stable ID order. |
| Later page | page=2 and limit=2 returns the next entries, not duplicates from page one. |
| Empty page | A valid page beyond available results is 200 with empty data, not 404 or panic. |
| Invalid page | page=0, negative text, or non-number gets the documented client error. |
| Invalid limit | limit=0, negative text, non-number, or value above 100 gets the documented client error. |
| Owner scope | User one’s list never includes user two’s Entry, regardless of page/limit. |
| Metadata | total and totalPages match the ordered owner-scoped result set. |
| No mutation | Repeated GET requests and invalid query requests do not create, modify, or delete Entries. |

### Implementation order

1. Complete api-contract.md and obtain review of the contract, not code.
2. Seed a focused test fixture with at least five Entries for one owner and one Entry for another.
3. Write the default and page-one failing handler tests.
4. Decide where parsed page/limit become a small application-level options value. Keep raw query
   strings and HTTP errors in the handler.
5. Add the smallest service/store behavior needed to return a deterministic page for an explicit
   owner. Do not add a generic paginator package.
6. Write invalid-input tests before adding parsing branches.
7. Make ordering explicit. Do not rely on map iteration or a test’s accidental insertion order.
8. Calculate metadata from the same owner-scoped, filtered result set being paginated.
9. Map the result to the documented JSON response in the handler.
10. Run go fmt ./..., focused tests, then go test ./....
11. Make one manual API-client request and compare it with the test response.
12. Explain every changed line: input, output, mutation, error path, package reason, and contract
    effect.

### Tiny design choices you must make explicitly

| Decision | Your answer before code |
| --- | --- |
| Is page numbering 0-based or 1-based? | |
| What happens when page is omitted? | |
| What happens when limit is omitted? | |
| What maximum limit prevents accidental large work? | |
| What stable order is used? | |
| What happens for a valid page with no items? | |
| Does total count entries owned by every user or only the caller? | |
| Does an invalid query reach the store? | |
| What pre-existing JSON contract changes? | |

### Small, focused implementation prompt

> Guide me through adding Session 7 page-and-limit pagination to GET /entries without writing the feature for me. First require my route contract, response example, page/limit defaults, maximum, stable order, and test plan. Then ask me to implement one small step at a time: handler query parsing, a small application list-options value if justified, owner-scoped deterministic paging, response metadata, and tests. Preserve Session 4 authorization and Session 5 validation/error rules. Do not add a router, database, cache, pagination library, generic query parser, OpenAPI generator, or versioned API rewrite.

### Stretch design exercise: a real custom action

After the required feature passes, design but do not implement this route:

~~~text
POST /entries/{id}/archive
~~~

Answer:

1. Why is archive potentially more than PATCH {"status":"archived"}?
2. What authorization rule applies?
3. What if the Entry is already archived?
4. What effects might later belong in one transaction or a background job?
5. Does the action create a resource, update a resource, or both?
6. Which success status/representation best describes the outcome?
7. Which tests must exist before an implementation?

This is deliberately a design exercise. Session 9 will introduce background-job trade-offs; do not
invent queues or notifications today.

---

## 7. Intentional API-design bug lab

Choose one bug only after your baseline feature tests pass. I should provide the symptom and
reproduction, not the source location or solution.

| Bug | Faulty behavior | Observable consequence |
| --- | --- | --- |
| Unbounded limit | The handler accepts any positive limit. | A client can request excessive work/data; later SQL queries become dangerous. |
| Map-order pagination | The store pages values straight from a map. | Page results change or duplicate across calls; tests become flaky. |
| 404 empty list | A valid collection filter/page with no entries returns 404. | Client cannot distinguish no matches from a missing endpoint/member. |
| Pagination before owner scope | The code pages all Entries, then filters to the caller. | Metadata leaks other users’ counts or pages appear inexplicably sparse. |
| Default mismatch | Documentation says limit 20 while code defaults to 10. | Clients and tests disagree even when code appears to work. |
| Invalid-query fallthrough | The handler writes a 400 but still calls the service. | Bad input may produce a second response or unnecessary work. |
| Generic envelope drift | An agent wraps every response, including health and member endpoints, without a compatibility decision. | Existing clients break for no product reason. |
| Route-as-verb drift | New handlers use paths such as /entries/list or /entries/delete. | The API becomes inconsistent and harder to learn. |
| Action-as-field bug | Archive is implemented as an unguarded JSON status field change. | Required side effects, authorization, or transition rules can be bypassed. |
| Documentation drift | The OpenAPI/example contract accepts a different limit or response shape from code. | Client integration fails despite apparently valid documentation. |

### Debugging protocol

1. Reproduce the symptom with the smallest handler/service test. Use a manual client only to observe
   the external message after the test exists.
2. State the violated contract in one sentence.
3. Identify whether it is route design, parsing, ownership, ordering, response mapping, documentation,
   or a compatibility defect.
4. Trace client input through the handler, service, and store.
5. Find the earliest place the contract is violated.
6. Add a regression test that fails before the fix.
7. Make the smallest fix at the shared decision point.
8. Re-run formatting, focused tests, all tests, and one manual request.
9. Update api-contract.md or its OpenAPI fragment if the intended contract became clearer.
10. Explain why a framework, generic utility, or broad rewrite is not required.

### Bug-lab prompt

> I intentionally introduced a Session 7 API-contract bug. Be a debugging partner, not an auto-fixer. Ask me for the failing request/test, the relevant route/handler/service/store code, and my documented contract. Help me state the violated client-visible promise, trace the earliest mismatch, and write a regression test. Do not propose code until I explain the root cause and smallest safe correction.

---

## 8. Audit AI-generated API changes

Generated API code is especially risky because it can look tidy while silently changing what every
client receives. Treat route tables, JSON examples, and tests as part of the diff.

### Audit rubric

| Review question | Evidence of a good answer |
| --- | --- |
| Is the resource model clear? | Collection, member, and action routes describe client intent without Go/UI verbs. |
| Does each method fit its behavior? | GET reads, POST creates or names an explicit action, PATCH is partial update, DELETE has intentional no-content behavior. |
| Are status codes outcome-based? | Creation, ordinary success, deletion, bad input, missing member, empty collection, and authentication are distinguishable. |
| Is list behavior bounded? | Defaults, max limit, accepted parameters, order, empty result, and metadata are documented and tested. |
| Is ordering deterministic? | The API never exposes map iteration as an accidental order. |
| Is owner scope preserved? | Pagination and total count are computed inside the verified caller’s permitted data set. |
| Are errors safe and consistent? | No raw Go, map, future SQL, or credential data reaches clients. |
| Is compatibility considered? | Response-shape changes are explicit and tests/documentation update together. |
| Do docs match code? | Collection examples, OpenAPI fragment, and handler tests agree. |
| Is complexity proportionate? | Standard library, narrow existing boundaries, no generator/framework/library for one list endpoint. |

### Common AI-generated-code failure modes

- Adding /v1 to every route without a compatibility plan, then breaking all existing tests and clients.
- Turning every handler into POST /doThing because it is quicker to implement.
- Calling PUT for partial updates or treating PATCH as a blind map of arbitrary fields.
- Returning 404 for an empty collection or 200 with an error string for invalid query input.
- Using map iteration as a list order.
- Letting a client pass an unbounded limit, arbitrary sort field, arbitrary SQL-like filter, or
  duplicated query keys without a defined contract.
- Counting/paging every user’s Entries before applying ownership rules.
- Returning all Entry fields, including internal ownership or future storage fields, because JSON
  encoding makes it easy.
- Replacing a raw array with a generic envelope without calling it a breaking change.
- Adding a generic Pagination interface, query-builder package, reflection parser, or router library
  for one endpoint.
- Treating an OpenAPI file as correct because it parses, without comparing it to tests and examples.
- Generating a Swagger server/client SDK before the routes are stable.
- Implementing archive as a free-form status update that bypasses domain policy.
- Logging raw Authorization values while debugging API-client failures.

### API-diff audit prompt

> Review this generated Session 7 TraceLog API diff before I accept it. Start by extracting the old and new client contract: routes, methods, query parameters, auth/owner scope, defaults, response JSON, statuses, errors, ordering, and side effects. Then trace GET /entries through handler, service, and store. Check pagination bounds, stable order, empty-list behavior, authorization before paging/counting, compatibility, tests, API-client/OpenAPI documentation drift, security, and unnecessary dependencies or abstractions. Cite exact evidence. Classify findings as must-fix, design decision, or safe simplification. Do not write replacement code unless I ask.

### Contract-to-test comparison prompt

> Compare this api-contract.md or OpenAPI fragment against the actual TraceLog handler tests and code. List every mismatch in route, HTTP method, query default/bound, status, response field, error behavior, ordering, authorization, or side effect. Then give me one minimal test or documentation change for each real mismatch. Do not assume the documentation is correct just because it is formatted.

### Documentation check prompt

> I wrote this small OpenAPI fragment manually for GET /entries. Explain every field that affects a client and identify only the missing details that prevent a client from using the endpoint correctly. Do not generate a complete spec, server, client SDK, or new dependency. Ask me to compare the final result with an httptest case.

---

## 9. Checkpoint: design, trace, and defend

Complete this without notes first.

1. Define an API contract in one sentence, then list its minimum components.
2. Classify these paths as collection, member, relationship, or action. Explain any poor design.

| Path | Classification and reason |
| --- | --- |
| GET /entries | |
| GET /entries/42 | |
| GET /entries?tag=go | |
| POST /entries/42/archive | |
| POST /createEntry | |
| DELETE /entries/42 | |
| GET /entries/list | |

3. Fill in the route/method/status decisions:

| Client intent | Method/path | Expected successful result |
| --- | --- | --- |
| Create a new Entry | | |
| Fetch a caller’s second page of Entries | | |
| Fetch one Entry | | |
| Partially change permitted Entry fields | | |
| Remove one Entry | | |
| Archive one Entry through a domain action | | |

4. Explain why an empty list is 200 and a missing specific Entry is 404.
5. State the default page, default limit, maximum limit, ordering rule, and invalid-query behavior for
   TraceLog.
6. Draw the execution path for GET /entries?page=2&limit=2 from URL to JSON.
7. Explain why owner scope must be applied before the page and total count are calculated.
8. Write one test that catches map-order pagination.
9. Describe a breaking API change you made in this session and why it is acceptable only in this
   controlled learning project.
10. Review an agent proposal to add an OpenAPI code generator and generic pagination package. Accept,
    reject, or defer each with a concrete current need.
11. Explain why POST /entries/{id}/archive might return 200 while a clone action could return 201.
12. Open your API-client collection or curl log and explain how it complements, but does not replace,
    your tests.

### Passing standard

You are ready for Session 8 when you can:

- Design a resource route table before requesting handlers.
- Explain each route’s method, input, response, side effect, and error behavior.
- Distinguish an empty collection from a missing member.
- Add and audit bounded, deterministic, owner-scoped pagination.
- Recognize a custom domain action and avoid reducing it to a careless field update.
- Trace a public API contract through Go package boundaries.
- Use an API client, raw HTTP output, an executable handler test, and a small documentation artifact
  as complementary evidence.
- Reject a generated API change that is inconsistent, unsafe, breaking, or over-engineered.

If any answer is weak, repeat the relevant contract table and one focused test. Do not “solve” API
uncertainty by adding routes, versions, dependencies, or abstractions.

---

## 10. Update your progress log

~~~markdown
## Session 7 — REST API design

- Date and actual time:
- Diagnostic answer I changed:
- My TraceLog resource map:
- My page/limit contract:
- Stable ordering rule and why:
- API-client or curl evidence:
- OpenAPI/Swagger documentation artifact:
- Manual feature and tests:
- Chosen bug, violated contract, root cause, and regression check:
- AI-generated API finding:
- One dependency, route, or abstraction I rejected and why:
- Compatibility decision I recorded:
- Confidence (1–5):
- Question to carry into Session 8:
~~~

---

## Ready for Session 8

Session 8 introduces PostgreSQL. Bring:

- the Session 7 api-contract.md or OpenAPI fragment;
- passing list pagination and ownership tests;
- your stable ordering and total-count decisions; and
- one question about how the same contract should survive a SQL implementation.

The database may change how TraceLog stores and retrieves pages. It must not quietly change what the
API promises.
