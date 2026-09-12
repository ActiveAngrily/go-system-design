# Session 1 — Backend Foundations and Go Program Flow

[← Course overview](../study-plan.md)

## Session contract

**Time:** 90–120 minutes  
**Lecture sequence:** L01–L04  
**Project:** TraceLog, version 0 — a tiny in-memory program  
**Goal:** Build the mental model required to read a small Go program confidently before we add HTTP, databases, or frameworks.

This is not a passive watching session. You will explain code, predict behavior, write a small amount of Go, diagnose one bug, and review deliberately small AI-generated code.

### What we are not building yet

No HTTP server, router, database, authentication, framework, interface hierarchy, or repository layer. Those would hide the foundations this session is meant to establish.

---

## Source material

Use the lectures in this order. Read the note first; use the transcript only when a point feels unclear.

| Order | Primary lecture material | Use it for |
| --- | --- | --- |
| 1 | [L01 — Roadmap](../resources/backend-from-first-principles/notes/01-1-roadmap-for-backend-from-first-principles.md) · [transcript](../resources/backend-from-first-principles/transcripts/01-0Rwb4Xmlcwc.en.vtt) | See the whole backend curriculum and where this first session fits. |
| 2 | [L02 — Backend engineer path](../resources/backend-from-first-principles/notes/02-2-walk-the-path-of-a-true-backend-engineer.md) · [transcript](../resources/backend-from-first-principles/transcripts/02-3qFjZbFRSAU.en.vtt) | Adopt an engineering workflow: understand systems, test assumptions, and trace failures. |
| 3 | [L03 — What is a backend?](../resources/backend-from-first-principles/notes/03-3-what-is-a-backend-how-do-they-work-and-why-do-we-need-them.md) · [transcript](../resources/backend-from-first-principles/transcripts/03-6Ss4dJD9Kzg.en.vtt) | Form the request-to-data-to-response mental model. |
| 4 | [L04 — First principles](../resources/backend-from-first-principles/notes/04-4-benefits-of-learning-backend-engineering-from-first-principles.md) · [transcript](../resources/backend-from-first-principles/transcripts/04-6fqZs5Z3k9A.en.vtt) | Learn why we build up from execution and data flow instead of memorizing tools. |

**Supporting reference only:** the markdown edition of *Let’s Go*, sections 1.1, 2.1, and 2.2. Use it to clarify Go modules, packages, and a minimal program; do not reorganize this lesson around book chapters.

---

## Outcomes

By the end of the session, you should be able to:

- Describe a backend as code that receives input, applies rules, reads or changes state, and produces an output.
- Trace a tiny Go executable from module to package to main function to helper call.
- Recognize variables, basic types, structs, slices, maps, pointers, functions, methods, and errors in a small program.
- Explain the difference between a Go module and a package.
- Explain why a nil map can be read but cannot be assigned to.
- Identify a sensible first place to investigate when a small backend program misbehaves.
- Refuse unnecessary architecture in a one-file, in-memory program.

You are not expected to memorize every Go built-in, method signature, package, or compiler error. You should know how to read the code, state what it is doing, and use Go documentation or a focused experiment when you are unsure.

---

## Active practice and tool lab

### Source-aware tool register

| Tool or technology | Source and role | Depth for this session | Proof of competence |
| --- | --- | --- | --- |
| Go toolchain: go run, go test, and go fmt | Course-supporting. These make the program’s execution, tests, and formatting observable. | Hands-on. | Run all three, predict what each will do, and explain the result. |
| go doc and go list | Course-supporting. They expose documentation and the module/package view behind the lecture’s first-principles approach. | Hands-on. | Use go doc to answer one concrete question and use go list to identify the package being built. |
| Node.js, Rust, Axum, and ORMs | L02/L04 use these as comparative examples, not as prerequisites for this Go course. | Recognition only. | State why an ORM or framework can hide execution flow; do not install any of them. |

There is no external Go library to add today. The useful tool skill is making the compiler, tests,
documentation, and package list answer a specific question.

### Manual-first rule

You may inspect the deliberately tiny starter program, but you must then manually implement or
repair one behavior: blank-title rejection, tag-map initialization, duplicate prevention, or the
summary output. Write its expected behavior and a test first. Before asking for a patch, explain the
call path from main to that behavior and name the line that mutates state.

Run this small tool drill from the session directory:

~~~text
go list
go doc strings.TrimSpace
go fmt ./...
go test ./...
go run .
~~~

Record what each command revealed. For the deliberate nil-map bug, make a prediction, reproduce it
with the smallest test or run, then use the panic and the relevant map documentation to identify why
the write fails. Do not ask an agent to replace the file.

### Completion evidence

You finish only after you can explain every production line in your manual change: values in,
values out, mutation, possible error, and why its package belongs where it does. Save the command
output or a short note in the progress tracker under both Evidence and Tool/library proof.

---

## Agenda

| Time | Activity | Evidence you produce |
| --- | --- | --- |
| 0–10 min | Diagnostic and session setup | Written answers, even if uncertain |
| 10–30 min | L01–L04 synthesis | A one-page backend component map |
| 30–45 min | Go execution walkthrough | A spoken or written trace of a tiny program |
| 45–70 min | Build TraceLog version 0 | A formatted program and a passing test |
| 70–85 min | Intentional bug investigation | A failing test, root cause, and smallest fix |
| 85–100 min | AI-generated code audit | A short review with evidence |
| 100–120 min | Checkpoint and progress log | Answers and next-session readiness |

If you have only 90 minutes, complete through the bug lab and defer the AI audit until the next session.

---

## 1. Start with a diagnostic

Answer these before reading anything. Short answers are enough. Do not look them up yet.

1. What is the first function called when a normal Go executable starts?
2. What is the difference between a function and a method?
3. What is the difference between a module and a package?
4. What happens if a program assigns a value into a nil map?
5. When should a program return an error rather than panic?
6. You are dropped into an unfamiliar backend repository. What three files or locations would you inspect first to understand how the program starts?
7. In broad terms, what stages happen between a client request and a client response?

Save your answers under a “Session 1 diagnostic” heading in your progress log. Incorrect answers are useful baseline data, not failure.

### Start-of-session prompt for Codex

> Teach Session 1 from sessions/session-01-backend-foundations.md. Begin by asking me the diagnostic questions one at a time. Do not give answers immediately. Ask follow-up questions that expose my reasoning, then explain the correct model with a small Go example. Keep the work scoped to L01–L04 and do not introduce HTTP implementation yet.

---

## 2. Lecture synthesis: a backend from first principles

After reviewing L01–L04, write this in your own words:

> A backend accepts input at a boundary, validates and transforms it, applies business rules, uses dependencies or state when needed, and returns an output. Every part can fail, so I need to trace data and errors across those steps.

Now make a component map for a familiar product, such as a notes app, online shop, or a task list.

| Question | Your answer |
| --- | --- |
| What input enters the system? | |
| What business rule decides whether it is valid? | |
| What state must be read or changed? | |
| What output should the caller receive? | |
| What can fail at each step? | |
| What information must not be trusted from the caller? | |

Do not design classes or tables yet. The point is to see the flow before selecting tools.

### The execution-flow model

For the TraceLog program you will build today, use this model:

~~~text
go run .
  -> loads the module
  -> compiles the main package
  -> calls main()
  -> main calls newEntry(...)
  -> main calls methods on Entry
  -> output is printed or an error is handled
~~~

Later, HTTP will become another input boundary. The tracing habit will not change:

~~~text
request -> handler -> validation -> application rule -> dependency/state -> response
~~~

---

## 3. Go concepts to understand today

| Concept | Understand deeply enough to explain | Look up when needed |
| --- | --- | --- |
| Module | A versioned collection of packages named by go.mod. | Exact go command options and module-path conventions. |
| Package | A group of Go source files compiled together; executable programs use package main. | The full standard-library package catalog. |
| main function | Program execution begins at main in the main package. | Build tags and alternate build modes. |
| Function | Named behavior that takes inputs and may return values or an error. | Less common parameter and return forms. |
| Method | A function attached to a named type through a receiver. | Method sets and advanced receiver edge cases. |
| Variable and type | A variable holds a value of a type; types constrain operations. | Every conversion rule and numeric limit. |
| Struct | A named collection of related fields. | Reflection and field-tag details. |
| Slice | A dynamic view over ordered values. | Capacity tuning and aliasing edge cases. |
| Map | A key-to-value collection; writes require an initialized map. | Performance internals and all iteration details. |
| Pointer | A value that refers to another value; useful when a method must mutate the original or copying is undesirable. | Unsafe pointers and low-level memory details. |
| Error | An ordinary returned value representing an expected failure path. | Every error-wrapping helper and pattern. |

### Recognition exercise

For each line below, state what it does before asking Codex to confirm.

~~~go
type Entry struct {
    ID    int
    Title string
    Tags  map[string]bool
}

func newEntry(id int, title string) (*Entry, error) {
    if title == "" {
        return nil, errors.New("title is required")
    }
    return &Entry{ID: id, Title: title, Tags: make(map[string]bool)}, nil
}

func (e *Entry) AddTag(tag string) {
    e.Tags[tag] = true
}
~~~

Questions:

1. Which identifier names a type?
2. Which values are inputs to newEntry?
3. Why does newEntry return a pointer?
4. What is initialized by make?
5. Why does AddTag use a pointer receiver?
6. What would happen if Tags were nil when AddTag runs?

---

## 4. Inspect the starter codebase

Ask me to generate the following only after you finish the lecture synthesis.

### TraceLog version 0

Create a temporary learning directory named tracelog-session-01 with only:

~~~text
tracelog-session-01/
  go.mod
  main.go
  main_test.go
~~~

The program should contain:

- An Entry struct with ID, Title, and Tags fields.
- A newEntry function that rejects a blank title.
- An AddTag method that avoids duplicate tags.
- A Summary method that returns a short string.
- A main function that creates one entry, adds a tag, and prints a summary.
- One small test for successful construction and one test for blank-title rejection.

Constraints:

- Use only the Go standard library.
- Keep all production code in main.go for this session.
- Do not add a database, web server, router, interface, controller, service, repository, configuration system, logger, or external dependency.
- Do not hide behavior behind helpers you cannot explain.

### Code-reading task

Before running the program, inspect it and answer:

1. Which package will the Go tool compile as the executable?
2. In exact call order, what happens after main starts?
3. Where is the blank-title rule enforced?
4. Which operation mutates state?
5. What result tells the caller that construction failed?
6. What invariant should always be true for a successful Entry?
7. What test would catch a nil-map bug?
8. Which field can later become a security or authorization concern, and why?

Then run these checks:

~~~text
go fmt ./...
go test ./...
go run .
~~~

Do not accept an explanation such as “it works.” Say what each command verifies:

- gofmt normalizes Go source formatting.
- go test compiles and executes tests.
- go run compiles and runs the executable path.

### Prompt to generate the starter code

> Create the smallest possible Go learning program in a new tracelog-session-01 directory. Use only go.mod, main.go, and main_test.go. It must contain an Entry struct; newEntry with blank-title validation; AddTag without duplicates; Summary; a printable main path; and two focused tests. Show the proposed file tree and code before writing files. Do not add HTTP, a database, interfaces, layers, third-party dependencies, or future-proof abstractions. After generating it, wait for my review.

---

## 5. Manual coding exercise

Now write or complete the implementation yourself. Keep it small.

### Required behavior

1. Construct an Entry with an integer ID and non-blank title.
2. If the title is blank after trimming surrounding spaces, return an error.
3. Start every Entry with a usable tag map.
4. Add a tag only once.
5. Produce a readable summary that includes the ID and title.
6. Test both success and the rejected blank-title path.

### Think before coding

Answer these aloud or in writing:

- Why is returning an error preferable to letting an invalid Entry travel through the program?
- Why is a pointer receiver appropriate for AddTag?
- What makes a map write different from reading a missing map key?
- Why is a map a reasonable data structure for duplicate prevention today?
- Why would a generic Store interface be needless in this program?

### Completion check

Your implementation is ready only if all of these are true:

- It is formatted.
- Tests pass.
- The normal program output matches your prediction.
- A blank title has a tested error path.
- Adding the same tag twice does not create duplicate state.
- You can explain every production-code line without asking an agent to translate it.

If you cannot explain a line, mark it with a question and ask for a narrow explanation. Do not replace the whole file.

### Narrow explanation prompt

> Explain only this Go line or function in the context of my TraceLog program: [paste it]. First describe the values entering and leaving it, then describe any mutation or error path. Do not simplify away details, rewrite the program, or introduce concepts not used here.

---

## 6. Intentional bug lab

Pick one bug below. I should inject it after your baseline tests pass; you should diagnose it before I offer a solution.

| Bug | Minimal faulty change | Observable symptom |
| --- | --- | --- |
| Nil map write | Return an Entry without initializing Tags. | Adding the first tag panics. |
| Value receiver | Change AddTag from a pointer receiver to a value receiver and mutate a non-reference field in the method. | The caller does not see the expected mutation. |
| Ignored construction error | Continue after newEntry returns an error. | A later nil dereference or invalid behavior occurs. |
| Missing trim | Check title before trimming whitespace. | A title containing only spaces is accepted. |
| Broken duplicate rule | Append tags to a slice without checking existing values. | Duplicate tags appear in the summary. |

### Debugging protocol

Use this order. It is the habit we are training.

1. Reproduce the failure with the smallest command or test.
2. State the actual symptom precisely: panic, wrong output, unexpected success, or compile failure.
3. Trace the relevant caller and callee instead of guessing.
4. State the invariant that should have held.
5. Add or improve one test that fails for the right reason.
6. Make the smallest root-cause fix.
7. Run formatting, the focused test, all tests, and the program path again.
8. Explain why the fix belongs where it does.

Do not ask an agent to rewrite the code blindly. A rewrite may hide the cause and create a different bug.

### Bug-lab prompt

> I intentionally broke my TraceLog Session 1 program. Act as a debugging partner, not an auto-fixer. Ask me for the failure output and the smallest relevant code first. Help me form hypotheses, trace the call flow, and write a failing test. Do not propose a patch until I state the invariant and probable root cause.

---

## 7. Audit AI-generated Go code

Have Codex generate its own solution to the starter-code prompt, but do not accept it without review.

Review it using this rubric.

| Question | What good evidence looks like |
| --- | --- |
| Does it meet the requested behavior? | Tests prove success, blank-title rejection, and duplicate handling. |
| Are mutation boundaries clear? | AddTag changes the original Entry intentionally and safely. |
| Is error behavior explicit? | Invalid construction returns an error; main handles it sensibly. |
| Is data initialized? | Tags is usable before its first write. |
| Is the design idiomatic and small? | Clear names, ordinary Go constructs, no needless abstraction. |
| Did the agent add unrequested architecture? | No interface with one implementation; no storage layers without storage. |
| Can you prove the code works? | Commands and tests cover the key paths. |

Write review comments in this form:

~~~text
Finding:
Evidence:
Why it matters:
Smallest acceptable change:
Test or check that proves it:
~~~

### Audit prompt

> Review this generated Go diff as if it will become the first commit of a backend service. Check behavior, error paths, mutation, nil handling, tests, Go idioms, and unnecessary abstractions. Cite the exact code responsible for each issue. Separate must-fix issues from optional simplifications. Do not write a replacement implementation unless I ask.

### Documentation comparison prompt

> Compare this small Go pattern with the relevant official Go documentation: [paste code or name the pattern]. Tell me what the language guarantees, what this code assumes, and one tiny experiment or test that would verify the assumption locally.

---

## 8. Checkpoint: explain, predict, modify

Complete this without notes first.

1. Draw or describe the execution path from go run . to the first printed output.
2. Explain module, package, function, and method in one sentence each.
3. Explain why returning an error is part of normal program control flow in Go.
4. Predict the outcome of reading from a nil map and writing to a nil map.
5. Given a failing AddTag test, name the first two locations you would inspect and why.
6. Add one safe behavior change: reject an empty tag after trimming it. Write the test before the implementation.
7. Review an agent suggestion to add an EntryRepository interface. Accept, reject, or defer it, with a reason tied to current requirements.

### Passing standard

You are ready for Session 2 when you can:

- Trace the small program accurately without hand-waving.
- Make the empty-tag change with a passing test.
- Diagnose the selected bug using evidence rather than a blind rewrite.
- Explain why the current program is intentionally simple.
- Identify one thing you would look up instead of pretending to memorize it.

If a checkpoint answer is weak, repeat the smallest relevant exercise. Do not rewatch all four lectures automatically.

---

## 9. Update your progress log

Add this after the session:

~~~markdown
## Session 1 — Backend foundations

- Date and actual time:
- Diagnostic answers I changed:
- Program-flow model in my own words:
- TraceLog commit or folder:
- Bug selected, symptom, root cause, and test:
- AI review finding:
- One Go concept I understand:
- One Go concept I will look up next time:
- Confidence (1–5):
- Question to carry into Session 2:
~~~

---

## Ready for Session 2

The next session adds the HTTP request and response boundary. Bring your TraceLog version 0 program and your answers to these questions:

- Where should untrusted request input enter the program?
- How can the program turn a returned error into a useful client response?
- Which parts of today’s program should remain unchanged when HTTP is introduced?

> Teach Session 2 only after I complete this checkpoint. Start by reviewing my Session 1 progress log and one piece of TraceLog code I can explain least confidently.
