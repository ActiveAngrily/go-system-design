# Session 5 — Validation, Transformation, and Data Integrity

[← Course overview](../study-plan.md) · [← Session 4](session-04-authentication-and-authorization.md)

## Session contract

**Time:** 100–120 minutes  
**Lecture sequence:** L09 — Validations and transformations  
**Project:** TraceLog version 2 — normalized, validated entry creation  
**Goal:** Turn a decoded but untrusted request into a small, well-defined internal value before it can change server state.

Session 3 taught that JSON can be decoded into Go values. Session 4 taught that an authenticated caller is not automatically authorized to do everything. This session adds the next important boundary: a decoded request is still not acceptable data until the server normalizes it and checks its rules.

### What we are not building yet

Do not add a validation framework, reflection-based struct-tag validator, schema generator, generic middleware pipeline, database constraint system, frontend, package-layer architecture, or a rules engine.

We will write one small, explicit preparation function and a few focused tests. Later sessions will place behavior into clearer application boundaries and databases will provide a second integrity backstop.

---

## Source material

Read L09 first. The book is supporting reference only; the course continues in lecture order.

| Primary material | What to focus on |
| --- | --- |
| [L09 — Validations and transformations](../resources/backend-from-first-principles/notes/09-9-validations-and-transformations-for-backend-engineers.md) | Server-side validation at the input boundary; structural, type, semantic, and complex checks; transformations; integrity and security; client validation versus server validation. |
| [L09 transcript](../resources/backend-from-first-principles/transcripts/09-qedj_JjjL-U.en.vtt) | Revisit the examples of syntactic versus semantic validation, normalization, invalid data reaching storage, and why frontend validation is only a user-experience aid. |
| *Let’s Go* §§8.1–8.6 | Supporting examples for form and validation patterns. Use them to recognize Go techniques, not to replace this lecture-led plan. |

---

## Outcomes

By the end of this session, you should be able to:

- Explain the difference between JSON decoding, syntactic validation, semantic validation, transformation, and authorization.
- Trace POST /entries from untrusted bytes to a normalized Entry that is safe to store.
- State why backend validation is mandatory even when a frontend validates input.
- Define clear data invariants for title, durationMinutes, and tags.
- Normalize input deliberately instead of letting formatting differences create duplicate or inconsistent state.
- Return a consistent client-error response without leaking internal details.
- Prove with tests that invalid input cannot mutate TraceLog state.
- Explain why a small pure prepare function is easier to read, test, and audit than a generic validation framework at this stage.
- Recognize which validation rules belong at the API boundary, which belong in business behavior, and which later require database constraints.
- Audit AI-generated Go code for missing validation, unsafe transformation, validation after writes, nondeterministic output, and invented abstractions.

You do not need to memorize every regular expression, validation-library API, Unicode corner case, email rule, phone-number format, schema language, or database constraint syntax. You should know the input contract, flow of trust, and how to look up a rule before you enforce it.

---

## Active practice and tool lab

### Source-aware tool register

| Tool or technology | Source and role | Depth for this session | Proof of competence |
| --- | --- | --- | --- |
| Insomnia | Demonstrated in L09 as a client for exercising an API. | Hands-on if available. | Build four local requests: valid, malformed JSON, semantically invalid JSON, and spoofed owner input; inspect status, body, and unchanged state. |
| curl or Postman | Course-supporting alternatives if Insomnia is unavailable. | Hands-on. | Reproduce the same contract matrix without relying on a frontend. |
| Frontend/form validation | L09 contrasts client validation with mandatory server validation. | Working familiarity. | Explain why it improves user experience but cannot establish backend integrity. |
| Validation libraries | Not demonstrated as a required solution in L09. | Deliberately deferred. | Reject an unnecessary library proposal and explain why this small pure function plus tests is clearer today. |

Do not add a frontend or a reflection/tag-based validator to make the exercise look more realistic.
The server contract is the source of truth.

### Manual-first rule

You must manually write the prepareEntry contract and its table-driven tests before generating any
implementation. Then implement one transformation yourself: title trimming, duration validation, or
stable tag normalization/deduplication. For every line, explain whether it decodes input, transforms
it, checks an invariant, writes state, or returns a client error.

### Tool drill

Create a local environment/base URL in Insomnia if you use it; never save tokens or secrets there.
Send each contract case and record the expected versus actual status, JSON response, and entry count.
Use handler tests to prove the count does not change on rejected requests. A client tool can reveal a
bad response; it cannot prove no hidden mutation occurred without a test or explicit state check.

### Completion evidence

Finish only after you can distinguish malformed JSON (400 in this course) from semantic field
errors (422), demonstrate one rejected request in a client tool, and show the focused test that
proves invalid input did not mutate TraceLog state.

---

## Agenda

| Time | Activity | Evidence you produce |
| --- | --- | --- |
| 0–10 min | Diagnostic and Session 4 recall | Short written answers |
| 10–30 min | L09 validation model | A validation matrix and request flow |
| 30–45 min | Inspect normalization behavior | Predicted input/output cases |
| 45–65 min | Inspect TraceLog version 2 | An end-to-end POST trace |
| 65–82 min | Test-first validation feature | A small pure function and handler tests |
| 82–97 min | Intentional bug lab | Root cause and regression test |
| 97–112 min | AI-generated-code audit | Evidence-based review |
| 112–120 min | Checkpoint and progress log | Session 6 readiness |

---

## 1. Start with a diagnostic

Answer without searching first.

1. Is valid JSON automatically valid business input? Give one example.
2. What is the difference between trimming a title and validating a title?
3. Why must a server validate input even if the web form already does?
4. What should happen to state when a request contains malformed JSON?
5. What should happen to state when JSON is well-formed but title is only spaces?
6. Why could tags Go, go, and space-go-space need normalization?
7. Why might a map be useful while normalizing tags but unsuitable as the final response order by itself?
8. What is the difference between 400 and the course’s chosen 422 response?
9. Does authentication make title, durationMinutes, or tags trustworthy?

Save your answers. We will use the ones you change as progress evidence.

### Start-of-session prompt for Codex

> Teach Session 5 from sessions/session-05-validation-and-transformation.md. Use L09 as the primary source. Ask my diagnostic questions one at a time, then make me distinguish decoding, normalization, syntactic validation, semantic validation, authentication, and authorization. Help me make one small pure preparation function and tests; do not add a validation library, framework, database, or broad middleware system.

---

## 2. The validation and transformation model

### The request pipeline

For authenticated entry creation, use this working model:

~~~text
route match
  -> authentication establishes Principal
  -> bounded body is decoded from JSON
  -> raw input is normalized
  -> normalized input is validated
  -> authorization permits creation for Principal
  -> server assigns ownership and stores Entry
  -> response is serialized as JSON
~~~

The exact ordering of small checks can vary. For example, request-body size must be limited before substantial processing, and some cheap structural checks happen during decoding. The invariant that matters is this:

> No invalid or unnormalized client input may reach the state-changing operation.

### Separate the jobs

| Job | Question it answers | TraceLog example |
| --- | --- | --- |
| Routing | Which handler gets this method and path? | Does POST /entries reach entry creation? |
| Authentication | Who is the caller? | Does the test token identify user-1? |
| Deserialization | Can bytes become a Go input value? | Is the body valid JSON with expected types? |
| Transformation | What canonical form should this value use? | Trim title; trim and lowercase tags; deduplicate tags. |
| Syntactic validation | Does the input have the expected basic form? | Is the decoded field a string, number, or array as required? |
| Semantic validation | Does the value make sense for this API? | Is title non-blank and duration within the defined range? |
| Authorization | May this principal perform this action? | May user-1 create an entry owned by user-1? |
| Storage integrity | Can the server safely persist this state? | Later: database constraints protect against invalid or conflicting data. |

Do not rename all of these as “validation.” Their different jobs tell you where to debug and what a test must prove.

### Client validation versus server validation

~~~text
client validation: fast feedback and usability
server validation: mandatory protection for integrity and security
~~~

A browser, mobile application, curl, a script, and an attacker can all call the same API. The backend must enforce its rules even if a frontend is perfect.

### Input sources are all untrusted

The same validation habit applies to:

- JSON request bodies.
- URL path parameters.
- Query parameters.
- Headers.
- File metadata.
- Form fields.
- Message-queue payloads.
- Data from another service.

Authentication changes who is making the request. It does not make their title or tags automatically valid.

---

## 3. Define TraceLog’s contract before coding

The feature for this session extends POST /entries.

### Raw create request

~~~json
{
  "title": "  Review Go routes  ",
  "durationMinutes": 45,
  "tags": [" Go ", "backend", "go"]
}
~~~

### Normalized stored and returned result

~~~json
{
  "id": 1,
  "title": "Review Go routes",
  "durationMinutes": 45,
  "tags": ["go", "backend"]
}
~~~

The API should not silently change the meaning of user input. We will make a few deliberate, documented transformations only.

### Rules for this learning API

| Field | Transformation | Validation rule | Why |
| --- | --- | --- | --- |
| title | Trim surrounding whitespace; preserve the user’s internal capitalization. | Required after trimming; 1–120 user-visible characters. | Prevent blank or accidentally padded titles while avoiding surprising title rewriting. |
| durationMinutes | No formatting transformation. | Integer from 1 through 480. | Gives a small, meaningful range for a learning session. |
| tags | Trim surrounding whitespace, lowercase, deduplicate while preserving first occurrence order. | Optional; each normalized tag is 1–30 user-visible characters; at most 5 supplied tags. | Makes filter-like labels consistent and bounded. |
| owner | No client input is accepted. | Assigned from verified Principal only. | Preserves the Session 4 ownership boundary. |
| unknown JSON fields | No transformation. | Reject under the strict request contract. | Prevents callers from assuming unsupported data is accepted or silently ignored. |

### Design decisions to state explicitly

- A malformed JSON document or wrong JSON type is a 400 request problem.
- A syntactically valid JSON body that violates title, duration, or tag rules is a 422 request problem for this course.
- Some APIs use 400 for both categories. Consistency and documentation matter more than imitating a status code without a contract.
- Validation errors should identify a field and a safe reason, but must not expose internal errors, credentials, stack traces, or database details.
- Tags are optional. Decide and test whether the response represents no tags as an empty JSON array rather than null.
- An integer field cannot distinguish a missing value from JSON zero without using a different input representation. For this rule, both are invalid, so no pointer is needed yet. Add a pointer only when product behavior genuinely needs that distinction.

### Field-error response shape

Use a small, stable error shape rather than raw Go errors.

~~~json
{
  "errors": [
    {
      "field": "title",
      "message": "must not be blank"
    },
    {
      "field": "durationMinutes",
      "message": "must be between 1 and 480"
    }
  ]
}
~~~

A slice is a useful representation because error order can be made stable. Do not collect errors in a map and then rely on map iteration order for a client contract or a test.

### Write a validation matrix

Before opening code, fill in expected transformation, status, and state effect.

| Input condition | Normalize? | Expected status | Entry count changes? | Expected field error |
| --- | --- | --- | --- | --- |
| Valid title, duration, and distinct tags | | | | |
| Title contains only spaces | | | | |
| Title exceeds the limit after trimming | | | | |
| durationMinutes is 0 | | | | |
| durationMinutes is 481 | | | | |
| tags contain Go, go, and spaces | | | | |
| One tag contains only spaces | | | | |
| More than 5 tags are sent | | | | |
| tags is omitted | | | | |
| JSON has an unknown owner field | | | | |
| JSON syntax is malformed | | | | |

---

## 4. Go concepts to understand today

| Concept | Understand deeply enough to explain | Look up when needed |
| --- | --- | --- |
| Pure function | Given the same raw input, it returns the same prepared value or field errors without touching HTTP or storage. |
| Slice | Holds ordered tags and ordered field errors; preserve predictable output order. | Allocation tuning and low-level backing-array behavior. |
| Map | Can record tags already seen during deduplication. | Map internals and performance details. |
| Zero value | A decoded int defaults to 0 when omitted; decide whether that distinction matters. | Custom unmarshaling and presence tracking. |
| Multiple return values | A preparation function can return a prepared value plus a list of errors. | Generic result types and clever error abstractions. |
| errors versus field errors | An internal Go error differs from safe client-facing validation feedback. | Full production error-taxonomy designs. |
| strings functions | Trim whitespace and normalize tag case deliberately. | Locale-specific text normalization and every Unicode edge case. |
| UTF-8 rune counting | Use an appropriate measure when a product rule says user-visible characters. | Grapheme clusters and international text policy. |
| JSON input type | Represents raw client input, not stored domain truth. | Reflection-based validation tags and schema generation. |

### Small code shape to recognize

~~~go
type createEntryInput struct {
    Title           string
    DurationMinutes int
    Tags            []string
}

type entryDraft struct {
    Title           string
    DurationMinutes int
    Tags            []string
}

func prepareEntry(input createEntryInput) (entryDraft, []fieldError) {
    // normalize and validate only
}
~~~

The raw input type reflects what arrived. The draft represents values that have passed the rules. This is a small boundary, not an invitation to build a type hierarchy.

### Why a map alone is not enough

For input tags:

~~~text
input:       [" Go ", "backend", "go"]
normalized:  ["go", "backend", "go"]
seen map:    {"go": true, "backend": true}
result:      ["go", "backend"]
~~~

The map helps answer “have I seen this tag?” The result slice preserves the meaningful first-seen order. If you rebuild the result directly from a map, ordering becomes unpredictable.

---

## 5. Inspect TraceLog version 2

Create a new directory named tracelog-session-05 by evolving the Session 4 API.

~~~text
tracelog-session-05/
  go.mod
  main.go
  server.go
  entry.go
  auth.go
  validation.go
  validation_test.go
  server_test.go
~~~

### Required behavior

- Existing authentication and ownership rules remain intact.
- POST /entries accepts title, durationMinutes, and optional tags.
- The server bounds the body before decoding it.
- Malformed JSON and type mismatch return a safe request error and do not mutate state.
- The prepare function trims and validates title.
- The prepare function validates durationMinutes.
- The prepare function normalizes, deduplicates, validates, and bounds tags.
- Unknown JSON fields are rejected under the chosen strict contract.
- Semantic validation failures return the course’s consistent field-error response and do not mutate state.
- A successful response exposes normalized values, not raw values.
- Handler tests prove state remains unchanged for every rejected create request.
- The code uses no validation framework, generic validator interface, or reflection-based tags.

### The end-to-end trace to annotate

Trace this request before running it:

~~~text
POST /entries
Authorization: Bearer test-user-one
Content-Type: application/json

{"title":"  Build a validator  ","durationMinutes":45,"tags":[" Go ","backend","go"]}
~~~

Fill in this table.

| Step | Value or decision | Trusted yet? | Possible failure |
| --- | --- | --- | --- |
| Authentication | | | |
| JSON decode | | | |
| Title normalization | | | |
| Duration validation | | | |
| Tag normalization and deduplication | | | |
| Field-error decision | | | |
| Ownership assignment | | | |
| Store write | | | |
| JSON response | | | |

### Code-reading task

Before executing the code, answer:

1. Where is the request body size limited?
2. Where do decoding errors stop the request?
3. Which function owns normalization and semantic validation?
4. Which function is allowed to mutate the store?
5. What guarantees that validation happens before a create operation?
6. How are field errors represented and ordered?
7. Where does owner ID come from?
8. Which exact transformation turns two tag spellings into one normalized tag?
9. Does an error response contain raw input or internal Go errors?
10. Which tests prove rejected requests did not increase the entry count?
11. What rule would need a database constraint later even if handler validation exists?

### Prompt to generate the inspection code

> Create a new tracelog-session-05 directory by evolving my Session 4 TraceLog API. Keep its test-only authentication and ownership policy. Extend POST /entries with title, durationMinutes, and optional tags. Add one small pure prepareEntry-style function that trims title, validates title length and duration 1–480, trims/lowercases/deduplicates tags in first-seen order, validates non-blank normalized tags, and limits supplied tags to five. Return stable field errors for semantic failures, use 400 for malformed JSON/type errors and 422 for semantic field errors, and prove rejected requests do not mutate state. Use only the standard library and focused tests. Do not add a validation library, tags/reflection validator, generic pipeline, database, or architectural layers. Show the invariants, error contract, file tree, and tests before writing code.

---

## 6. Manual coding and test-first feature work

### Start with pure preparation tests

Write tests for the preparation function before attaching it to the HTTP handler.

| Case | Expected prepared value or error |
| --- | --- |
| Title with outer spaces | Trimmed title, no title error |
| Title only spaces | Title field error |
| Valid duration 45 | Preserved value, no duration error |
| Duration 0 | Duration field error |
| Tags Go, backend, go | Tags become go, backend |
| Tag with only spaces | Tag field error |
| Six supplied tags | Tags field error |
| Omitted tags | Empty, consistent tag result and no tag error |

Then write handler tests for:

1. A valid authenticated create response contains normalized values.
2. A malformed JSON body returns 400 and does not add an entry.
3. A syntactically valid invalid body returns 422 and does not add an entry.
4. A client cannot include owner in JSON to change ownership.
5. A later GET by the owner shows the normalized stored Entry.
6. An unauthorized caller still cannot access the Entry.

### Implementation order

1. Make the data invariants visible in tests or a short contract comment.
2. Create the raw createEntryInput type.
3. Write a pure prepareEntry-style function.
4. Normalize title and tags before their semantic checks.
5. Preserve first-seen tag order while deduplicating.
6. Return a stable slice of field errors.
7. Call preparation only after bounded JSON decode succeeds.
8. Return immediately on field errors.
9. Create and store an Entry only from the prepared value plus verified Principal.
10. Format, run focused tests, run all tests, then make one manual request.

### Add one behavior yourself

Implement this feature with guidance:

> Tags should be case-insensitive and whitespace-insensitive, but response order should preserve the first meaningful occurrence.

Your tests should show:

~~~text
input tags:  [" Go ", "backend", "go", " BACKEND "]
stored tags: ["go", "backend"]
~~~

Do not sort tags unless the product contract asks for sorting. Deduplication and sorting are different transformations.

### Reasoning questions

- Why should title be trimmed before checking whether it is blank?
- Why should a user-visible title limit not quietly use a byte count unless the contract says bytes?
- Why is a body-size limit a separate protection from field-length validation?
- Why does a pure preparation function make tests faster and bugs easier to isolate?
- Why should an invalid request return before the store call?
- Why do all clients need server-side validation even if a browser form already validates?
- Why is a tag map useful internally but not necessarily correct as JSON output?

### Focused implementation prompt

> Guide me through implementing the Session 5 TraceLog input rules without writing the whole feature at once. First ask me to define the raw input, normalized result, invariants, field-error shape, and tests. Then help me make one small change at a time: title normalization, duration validation, tag deduplication with stable order, handler integration, and no-state-change tests. Do not introduce a validation library, generic validator, database, middleware framework, or new layers.

---

## 7. Intentional bug lab

Choose one bug after baseline tests pass. Your diagnosis must include the broken invariant and whether invalid state could be created.

| Bug | Faulty behavior | Consequence |
| --- | --- | --- |
| Validate after storage | Handler writes Entry before calling preparation or checking errors. | Invalid state persists even when client receives an error. |
| Blank-after-trim bug | Code checks title before trimming it. | A title containing only spaces is accepted. |
| Deduplicate-before-normalize bug | Code treats Go and space-go-space as different before trimming/lowercasing. | Semantically duplicate tags remain. |
| Map-order bug | Code builds output tags or errors by iterating only over a map. | Response ordering is nondeterministic and tests can be flaky. |
| Client-only validation | Server assumes a UI rejected bad input. | Scripts or malicious callers can create invalid state. |
| Owner-field regression | New input type adds ownerId and storage trusts it. | Session 4 ownership protection is bypassed. |
| Error-but-continue bug | Handler writes field errors and then creates the Entry. | Client sees failure but data changes. |
| Bounds-after-work bug | Code accepts an excessive input collection before it applies a meaningful bound. | Unnecessary resource use or policy bypass. |
| Raw-error leak | Response exposes decoder, internal, or storage error details. | API leaks implementation information. |

### Debugging protocol

1. Reproduce with a focused pure-function or handler test.
2. State the exact raw input and expected normalized output or error.
3. State the invariant: no invalid state, canonical tags, bounded input, or ownership preservation.
4. Trace body decode, preparation, error response, and store call in order.
5. Identify the earliest point the invariant could have been enforced.
6. Add a regression test that fails before the fix.
7. Make the smallest root-cause fix.
8. Verify entry count and stored values after the failure path.
9. Format and run all tests.
10. Explain why a framework would not make this one rule clearer.

### Questions I will ask before helping

- Is the problem structural decoding, normalization, semantic validation, authorization, or storage?
- What did the caller send?
- What should the canonical stored value be?
- Did the store change despite an error?
- Which function had enough information to stop the invalid input earliest?
- Does the failing test verify status only, or also verify no state changed?

### Bug-lab prompt

> I intentionally introduced a Session 5 validation or transformation bug. Act as a debugging partner, not an auto-fixer. Ask for the smallest failing test, raw input, actual stored/output value, and relevant preparation and handler code. Help me state the broken invariant, trace whether state mutated, and add a regression test. Do not propose a patch until I identify the earliest correct enforcement point.

---

## 8. Audit AI-generated validation code

Generated validation code often looks thorough while missing the actual contract or adding difficult-to-audit machinery.

### Audit rubric

| Review question | Evidence of a good answer |
| --- | --- |
| Are decoding and semantic validation separated? | Malformed JSON/type errors differ from rule violations in code and tests. |
| Is normalization deliberate? | Title and tags are transformed in the defined order before rule checks. |
| Are values bounded? | Body and collection/field limits exist before state changes. |
| Is ordering stable? | Output tags and field errors do not depend on map iteration. |
| Can a rejected request mutate state? | Tests prove entry count and stored data are unchanged. |
| Does ownership remain server-controlled? | No input field can choose owner. |
| Are response errors safe and useful? | Field names and safe messages appear; internals and raw secrets do not. |
| Are all client paths covered? | Server-side tests do not depend on frontend behavior. |
| Is the design small? | One explicit preparation function, no generic validation framework or reflection machinery. |
| Are business rules named? | Limits and transformations are visible, documented, and tested rather than accidental. |

### Common AI-generated-code failure modes

- Treating successful JSON decoding as complete validation.
- Using an unbounded body read and only checking field values later.
- Checking an untrimmed string for emptiness.
- Lowercasing a human title when only tags should be canonicalized.
- Deduplicating with a map and losing intended order.
- Validating after a store write or returning error without returning from the handler.
- Accepting unknown fields silently when the documented API is strict.
- Depending on browser validation for security or integrity.
- Returning raw decoder, database, or internal errors to clients.
- Adding owner or role fields to client input and trusting them.
- Adding a generic validation interface, reflection tags, or an external library before the concrete rules are understood.
- Adding arbitrary regular expressions without explaining their limits or test cases.
- Over-validating a domain field based on invented assumptions rather than a product rule.

### Audit prompt

> Review this generated Session 5 Go diff before I accept it. Trace untrusted input from body through decoding, normalization, semantic validation, authorization, storage, and response. Verify body and field bounds, transformation order, tag ordering, field-error stability, ownership protection, status decisions, and no-state-change behavior on errors. Find invented requirements, internal error leaks, missing tests, and unnecessary validation abstractions. Cite exact code and separate must-fix defects from optional simplifications. Do not replace the implementation unless I ask.

### Assumption-challenge prompt

> List every validation and transformation rule this TraceLog change enforces or assumes. For each, state its input source, canonical form, failure response, state effect, test, and product rationale. Mark rules that are explicit requirements versus agent-invented policy. Recommend deleting any rule that cannot be justified.

---

## 9. Checkpoint: validate, transform, and prove

Complete this without notes first.

1. Explain decoding, transformation, syntactic validation, semantic validation, authentication, and authorization in one sentence each.
2. Trace this request from raw JSON to outcome:

~~~json
{
  "title": "   ",
  "durationMinutes": 45,
  "tags": ["Go", " go "]
}
~~~

State the normalized values, field errors, status category, and whether an Entry may be stored.

3. Explain why frontend validation cannot protect TraceLog data integrity.
4. Explain why Go and go should result in one tag but title capitalization should not be changed automatically.
5. Predict the expected result for each request:

| Request condition | Expected category | State change? |
| --- | --- | --- |
| Valid authenticated create | | |
| Malformed JSON | | |
| Wrong JSON type for durationMinutes | | |
| Empty title after trimming | | |
| durationMinutes is 0 | | |
| Sixth supplied tag | | |
| Duplicate normalized tags | | |
| Unknown owner field | | |
| Missing authentication | | |

6. Write a test proving a semantically invalid POST /entries does not increase entry count.
7. Write a test proving tag normalization preserves first-seen order.
8. Review an agent proposal to install a validation library. Accept, reject, or defer it using the current rules and codebase size.
9. Name one field rule you would look up or obtain from product requirements rather than invent.

### Passing standard

You are ready for Session 6 when you can:

- Explain where every input rule belongs in the request flow.
- Write a pure normalization/validation test before handler code.
- Prove invalid input cannot mutate state.
- Explain why server validation persists even with a polished client.
- Audit a generated diff for transformation-order and state-integrity bugs.
- Defend each implemented rule with a product, integrity, security, or contract reason.

If one answer is weak, repeat the focused test or request trace. Do not add a more complex framework to compensate for unclear fundamentals.

---

## 10. Update your progress log

~~~markdown
## Session 5 — Validation and transformation

- Date and actual time:
- Diagnostic answer I changed:
- My request pipeline:
- TraceLog folder or commit:
- Input rules and product rationale:
- Chosen bug, broken invariant, root cause, and regression test:
- One generated-code finding:
- Validation detail I will look up when needed:
- Confidence (1–5):
- Question to carry into Session 6:
~~~

---

## Ready for Session 6

Week 2 begins with L10: controllers, services, repositories, middleware, and request context. Your current code has deliberately stayed flat so every boundary is visible.

Bring your TraceLog version 2 code and prepare to answer:

- Which code is HTTP-specific?
- Which code is pure behavior that could run outside an HTTP handler?
- Which code owns in-memory state today?
- Which boundaries would make the next feature easier to change or test?
- Which abstractions would still be premature?

> Teach Session 6 only after I show you the Session 5 checkpoint answers and can trace a rejected POST from body through a no-state-change assertion.
