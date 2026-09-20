# Session 8 — PostgreSQL Persistence, Queries, and Transactions

[← Course overview](../study-plan.md) · [← Session 7](session-07-rest-api-design.md)

## Session contract

**Time:** 115–120 minutes  
**Lecture sequence:** L12 — Mastering Databases with PostgreSQL  
**Project:** TraceLog version 2.3 — PostgreSQL-backed Entries and tags with the Session 7 API
contract preserved  
**Goal:** Move TraceLog's storage boundary from memory to PostgreSQL without changing its public
behavior. Model the existing invariants in a migration, run parameterized and cancellation-aware
queries through Go's database/sql package, and use one transaction to prove that an Entry and its
tags are stored together or not at all.

Session 7 defined the client promise. This session changes only how durable state is stored:

~~~text
authenticated request
  -> handler preserves the API contract
  -> service preserves validation and ownership policy
  -> PostgreSQL store executes parameterized SQL with request context
  -> schema constraints defend durable data
  -> transaction commits Entry + tags together or rolls both back
  -> handler returns the same safe status and JSON contract
~~~

The database is not a replacement for the handler or service. It is a durable boundary with its
own guarantees, failure modes, resources, and concurrency behavior.

### What we are not building yet

Do not add:

- Redis, an application cache, a task queue, background jobs, Elasticsearch, replicas, sharding, or
  a cloud database service.
- An ORM, query builder, code generator, generic repository framework, or a second database access
  library.
- A custom migration framework. Keep ordinary SQL migration files and use one existing migration
  command or psql.
- A database connection per request. One long-lived sql.DB handle already manages a pool.
- Dynamic SQL for arbitrary client-selected columns, sort fields, table names, or operators.
- Triggers merely because L12 demonstrates them. TraceLog has no update-timestamp behavior that
  needs one in this session.
- Indexes for imagined queries, database tuning flags, custom isolation levels, prepared-statement
  caches, retry loops, or premature pool tuning.
- A rewrite of the Session 7 handlers, JSON shapes, authentication stub, validation rules, or
  service boundary.

We will add the one unavoidable application dependency: a PostgreSQL driver compatible with
database/sql. Everything else remains SQL, the Go standard library, and the existing TraceLog code.

---

## Source material

Read L12 as the primary conceptual sequence. Use the Go and PostgreSQL documentation to verify the
exact behavior of the APIs and SQL you exercise.

| Primary material | What to focus on |
| --- | --- |
| [L12 — Mastering Databases with PostgreSQL](../resources/backend-from-first-principles/notes/12-12-mastering-databases-with-postgres.md) | Persistence and DBMS responsibilities; relational versus non-relational storage; schemas, types, keys, relationships, migrations, CRUD SQL, joins, parameterized queries, pagination, indexes, and triggers. |
| [L12 transcript](../resources/backend-from-first-principles/transcripts/12-F7Vwp2Xo5Do.en.vtt) | The full lecture examples: PostgreSQL modeling, migration workflow, TablePlus inspection, relationship tables, ordered queries, page/limit translated to offset/limit, parameterized values, index trade-offs, and trigger demonstration. |
| *Let's Go* §§4.1–4.9 | Supporting Go detail for drivers, database/sql, one long-lived pool, model methods, Exec/QueryRow/Query, Scan, rows cleanup, error handling, and transactions. The book uses MySQL; translate placeholders and types deliberately for PostgreSQL. |
| [Go database access documentation](https://go.dev/doc/database/) | Current standard-library lifecycle: opening a handle, pooling, querying, cancellation, resource cleanup, and transactions. |
| [pgx database/sql compatibility layer](https://pkg.go.dev/github.com/jackc/pgx/v5/stdlib) | The narrow PostgreSQL driver bridge used by this course; verify the current module documentation before adding it. |
| [PostgreSQL documentation](https://www.postgresql.org/docs/current/) | Exact SQL behavior for constraints, transactions, EXPLAIN, data types, and indexes when the lecture summary is not precise enough. |

### Lecture concepts to carry forward

L12 moves from why persistence exists to how a backend models and queries durable relational data:

- A DBMS organizes persistent data, provides access operations, enforces integrity, coordinates
  concurrent work, and controls database-level access.
- A relational database uses tables, rows, typed columns, keys, and relationships. A flexible
  document model is not automatically better or worse; the current data and access patterns decide.
- A schema is executable policy. Types, NOT NULL, CHECK, UNIQUE, primary-key, and foreign-key
  constraints reject states the database must never persist.
- Migrations make schema changes explicit and ordered. They are source-controlled changes, not
  undocumented clicks in a GUI.
- SQL expresses create, read, update, delete, joins, filters, sorting, and pagination.
- Client values belong in parameters. They do not belong in SQL strings built with concatenation.
- List queries need explicit ORDER BY. Relational results have no useful implied order.
- Indexes trade write/storage cost for selected read paths. A field appearing in a WHERE, JOIN, or
  ORDER BY is a candidate to investigate, not an automatic command to add an index.
- Triggers can enforce database-side behavior, but hidden behavior also costs clarity. Use one only
  for an invariant or behavior that genuinely belongs in the database.

The course adds Go-specific depth required by the study plan: sql.DB pooling, request-context
propagation, rows cleanup, stable error classification, and a transaction around a multi-write
invariant.

---

## Outcomes

By the end of this session, you should be able to:

- Explain persistence, DBMS, relational schema, table, row, column, primary key, foreign key,
  constraint, migration, query, transaction, and index using concrete TraceLog examples.
- Explain why PostgreSQL fits TraceLog's related users, entries, and tags without claiming that one
  database type is universally best.
- Translate existing application invariants into narrow database constraints.
- Model one-to-many and many-to-many relationships and explain why Entry tags need an association
  table when tags are reusable values.
- Apply and inspect an ordered SQL migration instead of changing the schema only through a GUI.
- Explain why application validation improves client feedback while database constraints protect
  durable truth from every writer.
- Use PostgreSQL placeholders such as $1 and pass values separately.
- Explain why parameterization protects values but cannot safely parameterize arbitrary identifiers
  such as a column name or sort direction.
- Distinguish ExecContext, QueryRowContext, and QueryContext.
- Explain why *sql.DB is a concurrency-safe pool handle rather than one physical connection.
- Open one database handle at startup, confirm connectivity with PingContext, inject it into the
  PostgreSQL store, and close it during application shutdown.
- Propagate the request context to every database operation that belongs to that request.
- Iterate sql.Rows safely with defer rows.Close, Scan, and rows.Err.
- Translate sql.ErrNoRows into the established application not-found category without leaking raw
  database errors to HTTP clients.
- Preserve Session 7's owner-scoped, ID-ascending, bounded page-and-limit contract in SQL.
- Calculate OFFSET from a validated 1-based page and prove ownership is applied before count and
  pagination.
- Use a transaction for Entry plus tag inserts and explain commit, rollback, and the invariant it
  protects.
- Write one PostgreSQL integration test that proves a failed tag insert leaves no partial Entry.
- Use EXPLAIN to inspect a query plan without turning one tiny test plan into a performance claim.
- Audit generated schema and storage code for unsafe SQL, context loss, resource leaks, partial
  writes, broken ownership, fragile scans, leaked errors, unsupported indexes, and extra layers.

You do not need to memorize every PostgreSQL type, SQL command, isolation level, index structure,
query-plan node, pool setting, or driver-specific API. You must understand the data invariant and
resource lifecycle, be able to trace one query, and know what evidence to inspect next.

### Understand deeply versus recognize and look up

| Learn deeply now | Recognize and look up when needed |
| --- | --- |
| Schema constraints; primary/foreign keys; explicit ordering; parameterized values; sql.DB as a pool; QueryRow versus Query; Scan; rows Close/Err; context propagation; sql.ErrNoRows; transaction begin/rollback/commit; owner-scoped pagination; integration-test isolation. | Normal forms beyond the current model; every PostgreSQL type; custom isolation levels; deadlock diagnosis; lock modes; B-tree internals; partial/expression indexes; materialized views; triggers; stored procedures; replication; partitioning; sharding; connection proxies; ORM and query-generator APIs. |

---

## Active practice and tool lab

### Source-aware tool register

| Tool or technology | Source and role | Depth for this session | Proof of competence |
| --- | --- | --- | --- |
| PostgreSQL | The relational DBMS used throughout L12. | Hands-on. | Apply the TraceLog migration to a dedicated local test database, run valid and invalid writes, and inspect resulting rows and constraints. |
| TablePlus | Demonstrated in L12 for viewing tables, rows, and changes. | Hands-on if already available; otherwise use psql. | Inspect the same schema and data you can explain from raw SQL. Do not let GUI clicks become the only migration history. |
| dbmate or the existing migration command | L12 demonstrates a small command-line migration workflow; the transcript may render the tool name imprecisely. | Hands-on if already available. | Apply the ordered migration and show its recorded state. If no migration tool is installed, apply the SQL with psql for this exercise rather than adding a custom runner. |
| psql | Course-supporting PostgreSQL client. | Hands-on. | Inspect table definitions, execute one parameterized-style prepared statement, provoke one constraint failure, and run EXPLAIN on the list query. |
| database/sql | Go standard-library database boundary used by the study plan and supporting book. | Hands-on. | Trace one pool handle from startup to QueryContext/Scan/Close and one transaction from BeginTx to Commit or deferred Rollback. |
| pgx/v5/stdlib | Narrow driver required to connect database/sql to PostgreSQL. | Working use. | Explain why it exists, inspect the added go.mod/go.sum diff, and use PostgreSQL $1 placeholders. Do not use both native pgx and database/sql APIs in this session. |
| go test | Course-supporting proof tool. | Hands-on. | Run unit tests plus a focused PostgreSQL integration test against a dedicated test DSN. |
| EXPLAIN | L12 motivates indexes from query access paths; PostgreSQL exposes the chosen plan. | Bounded inspection. | Read the plan for the owner-scoped ordered list query and state what the tiny fixture cannot prove. |

A GUI is optional. A real PostgreSQL server and observable SQL behavior are not. Use an existing
local installation or container setup if one already exists; do not introduce orchestration files
merely for this lesson.

### Manual-first rule

Before an agent writes a store method, you must:

1. Write the durable invariants for users, Entries, and tags in plain language.
2. Draw the tables, primary keys, foreign keys, and relationship cardinalities.
3. Write the parameterized SQL for one owner-scoped page by hand.
4. State why creating an Entry with its tags is one transaction.
5. Write the failing integration test for partial transaction failure.
6. Explain how request cancellation reaches the driver and PostgreSQL operation.

You will manually implement one query and one transaction boundary. An agent may review the schema,
point to documentation, or inspect your failure, but it may not replace your attempt with an ORM or
generate the whole persistence layer first.

### Tool drill

1. Connect to a dedicated TraceLog test database with psql. Print the current database and server
   version so you know which system you are changing.
2. Apply the first migration. Use psql's table-description commands or PostgreSQL catalog queries to
   inspect columns, defaults, constraints, and indexes.
3. Insert the two test principals used by the authentication stub. Insert one valid Entry and tags
   with SQL, then query them back in explicit ID order.
4. Attempt one invalid row, such as a duration outside 1–480 or a tag longer than 30 characters.
   Capture the constraint failure and verify that no valid application path depends on its raw text.
5. Prepare and execute one owner lookup with values supplied separately. Compare it with unsafe
   string concatenation and explain why only the first is acceptable.
6. Run the Session 7 page query with WHERE owner_id, ORDER BY, LIMIT, and OFFSET. Confirm another
   user's row never appears.
7. Run EXPLAIN on that query. Identify scan, filter/index condition, and sort nodes if present. Do not
   infer production performance from a few rows.
8. Run the focused Go integration test, cancel one context before a query, and inspect the returned
   error category without exposing it through the HTTP response.

### Completion evidence

Finish the tool lab with:

- the migration filename and successful apply command;
- an inspected schema showing keys and constraints;
- one valid query result and one deliberate constraint failure;
- the parameterized owner-scoped list SQL and its EXPLAIN output;
- the exact integration-test name proving transaction rollback;
- a passing focused integration test and go test ./... result; and
- a sentence explaining why sql.DB, sql.Rows, and sql.Tx have different lifecycles.

Never run destructive cleanup against an unverified database. Use a dedicated test DSN and test
rows with unique owner IDs; delete only those rows if cleanup is necessary.

---

## Agenda

| Time | Activity | Evidence you produce |
| --- | --- | --- |
| 0–12 min | Diagnostic and Session 7 recall | Predictions about persistence, constraints, and query flow |
| 12–28 min | L12 relational model | TraceLog table/relationship diagram |
| 28–43 min | Schema and migration design | Durable invariant table and reviewed migration |
| 43–58 min | database/sql lifecycle | Pool, context, query, rows, and error trace |
| 58–82 min | Manual query and transaction | Parameterized page query plus atomic Entry/tag creation |
| 82–98 min | PostgreSQL integration tests | Rollback, ownership, not-found, and cancellation evidence |
| 98–108 min | EXPLAIN and intentional bug | Query-plan observation and root-cause diagnosis |
| 108–118 min | AI-generated SQL audit | Evidence-based findings and smallest safe changes |
| 118–120 min | Checkpoint and progress log | Readiness decision for caching and background work |

If you have only 90 minutes, complete the schema, one parameterized query, the transaction, and its
rollback integration test. Do the EXPLAIN and AI-audit exercises before Session 9. Do not skip the
transaction proof.

---

## 1. Start with a diagnostic

Answer without searching first. Keep each answer tied to TraceLog rather than reciting a definition.

1. What problem does PostgreSQL solve that the in-memory store cannot solve after the process exits?
2. Which TraceLog rules belong in application validation, which belong in database constraints, and
   which belong in both?
3. Why does a foreign key from entries.owner_id to users.id matter if the service already assigns
   the authenticated principal ID?
4. Is sql.DB one open network connection? What happens when concurrent requests use it?
5. What is the difference between QueryRowContext and QueryContext?
6. Where does an error from QueryRowContext appear when no call to Scan has happened yet?
7. Why is this unsafe: "WHERE owner_id = '" + ownerID + "'"?
8. Can $1 safely stand for a client-supplied column name in ORDER BY $1? Why or why not?
9. What order will SELECT return if the query omits ORDER BY?
10. For page=3 and limit=20, what OFFSET should the SQL receive?
11. Why should Entry insertion and all of its tag insertions share one transaction?
12. What must happen if the second tag insert fails after the Entry row was inserted?
13. Why must rows.Close and rows.Err both appear in a multi-row query path?
14. Which response should a client see for a raw PostgreSQL constraint or connection error?
15. Does adding an index guarantee this list query will become faster? What evidence is missing?

Save your answers under “Session 8 diagnostic” in the progress log. Do not correct them yet.

### Start-of-session prompt for Codex

> Teach Session 8 from sessions/session-08-postgresql-persistence.md. Use L12 as the primary source and the official Go/PostgreSQL documentation for exact behavior. Ask the diagnostic questions one at a time. Require me to define the schema invariants, draw relationships, write one parameterized owner-scoped page query, and write the failing rollback integration test before implementation. Make me implement one QueryContext path and one Entry-plus-tags transaction manually. Preserve the Session 7 API contract. Use the hint ladder; do not add an ORM, query generator, cache, queue, search service, generic repository, custom migration framework, or speculative tuning.

---

## 2. Persistence and relational integrity

Persistence means TraceLog state survives beyond the lifetime of one Go process. PostgreSQL does
more than write bytes to disk: it gives multiple clients a structured query interface, coordinates
concurrent work, checks declared constraints, and commits or rejects changes with transaction
semantics.

### The three-layer integrity model

The same rule can appear at more than one boundary for different reasons:

| Layer | Job | TraceLog example |
| --- | --- | --- |
| HTTP/input preparation | Give precise client feedback and normalize untrusted input. | Trim title, lowercase/deduplicate tags, return stable field errors. |
| Service/application behavior | Enforce use-case and identity policy. | Assign owner from the verified principal; never from JSON. |
| PostgreSQL schema | Reject invalid durable state from every writer. | owner_id NOT NULL with a foreign key; duration CHECK; unique tag membership. |

Application validation does not replace constraints because tests, scripts, future services, manual
SQL, and defects may write through another path. Constraints do not replace application validation
because a raw constraint name is poor client feedback and arrives too late to express the whole API
contract cleanly.

### Relational choices for TraceLog

The current data has stable fields and explicit relationships:

~~~text
users 1 ───────< entries
                    |
                    | 1
                    v
              entry_tags >────── 1 tags
~~~

- One user can own many Entries; each Entry has one owner.
- One Entry can have many tags, and the same normalized tag can label many Entries.
- The entries.owner_id foreign key represents the ownership relationship.
- entry_tags represents many-to-many membership and stores position so the Session 5 first-seen tag
  order remains reproducible.

Do not put a comma-separated tag string in entries. It makes membership, uniqueness, filtering, and
constraints application parsing problems again. Do not add a universal entity table or polymorphic
relationship for three concrete tables.

### Schema decisions before SQL

Fill the last column before writing a migration.

| Durable rule | Database mechanism | Your exact decision |
| --- | --- | --- |
| Auth stub principals referenced by Entries must exist | users primary key plus entries foreign key | |
| Entry IDs are unique and generated by the database | bigint identity primary key | |
| Every Entry has one owner | owner_id NOT NULL REFERENCES users(id) | |
| Title is present and bounded after normalization | NOT NULL plus a CHECK consistent with the application contract | |
| Duration is 1–480 minutes | integer NOT NULL CHECK | |
| Creation time is always recorded | timestamptz NOT NULL DEFAULT current timestamp | |
| Normalized tag name is present and bounded | tags primary key plus CHECK | |
| An Entry cannot contain one tag twice | entry_tags composite primary key or UNIQUE | |
| Tag response order preserves first meaningful occurrence | position column plus per-entry uniqueness | |
| Deleting an Entry cannot leave memberships behind | entry_tags foreign key with deliberate delete behavior | |
| Deleting a user with Entries has a declared policy | choose RESTRICT or CASCADE and explain it | |

Use text for the teaching-stub user IDs because earlier sessions use values such as user-1. Do not
invent UUID conversion work during the persistence lesson. Use timestamptz for real instants; avoid
losing timezone meaning with a timestamp choice you cannot explain.

---

## 3. Migrations, constraints, and query-shaped indexes

A migration is a reviewed, ordered schema change. It should make a blank test database reach the
same structure every time.

### Minimum migration set

Keep the first version small:

~~~text
migrations/
  001_create_tracelog.sql
~~~

Whether your migration tool uses separate up/down files or marked sections is a tool choice. The
durable content is ordinary SQL that creates users, entries, tags, and entry_tags with their keys and
constraints.

Do not change a table manually in TablePlus and then attempt to remember the click sequence. The GUI
can inspect the result; the migration is the reproducible source of truth.

### A schema shape to review, not paste blindly

~~~sql
CREATE TABLE users (
    id text PRIMARY KEY
);

CREATE TABLE entries (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    owner_id text NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    title text NOT NULL CHECK (char_length(title) BETWEEN 1 AND 200),
    duration_minutes integer NOT NULL CHECK (duration_minutes BETWEEN 1 AND 480),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE tags (
    name text PRIMARY KEY CHECK (char_length(name) BETWEEN 1 AND 30)
);

CREATE TABLE entry_tags (
    entry_id bigint NOT NULL REFERENCES entries (id) ON DELETE CASCADE,
    tag_name text NOT NULL REFERENCES tags (name) ON DELETE RESTRICT,
    position smallint NOT NULL CHECK (position BETWEEN 0 AND 4),
    PRIMARY KEY (entry_id, tag_name),
    UNIQUE (entry_id, position)
);
~~~

Compare every number and delete action with the earlier API contract. If your established title
limit differs, preserve it instead of accepting this example. Schema and application rules must not
drift.

### The first justified index

The Session 7 hot path is already known:

~~~sql
SELECT id, owner_id, title, duration_minutes, created_at
FROM entries
WHERE owner_id = $1
ORDER BY id ASC
LIMIT $2 OFFSET $3;
~~~

The primary key already indexes id, but this query first scopes by owner and then orders that
owner's rows. One candidate is:

~~~sql
CREATE INDEX entries_owner_id_id_idx ON entries (owner_id, id);
~~~

Treat it as a hypothesis tied to a real query. Inspect EXPLAIN before and after with a meaningful
fixture later. A sequential scan on a tiny table is not evidence that PostgreSQL is broken, and an
index appearing in the schema is not evidence that production latency improved.

Defer tag-filter indexes until the tag filter exists and its query is known. Do not add a trigger:
there is no current updated_at contract to automate.

### Migration safety questions

Before applying a migration, answer:

1. Which database and schema will change?
2. Is it a dedicated local/test database?
3. What data or table can the down path remove?
4. Can the up path run on the expected prior schema?
5. Do application and migration deployment order need coordination?
6. How will a blank database prove the full migration history works?
7. Which constraints can reject existing rows in a later migration?

For this session, apply to a disposable dedicated test database. Do not practise destructive down
migrations on a database containing anything you care about.

---

## 4. Go's database boundary

### One long-lived pool handle

The lifecycle is:

~~~text
main startup
  -> read dedicated DATABASE_URL
  -> sql.Open("pgx", dsn) creates a database handle
  -> PingContext confirms a usable connection now
  -> inject *sql.DB into postgres.Store
  -> concurrent requests borrow connections through the handle
  -> server shutdown completes
  -> db.Close releases pool resources
~~~

sql.Open validates driver arguments and returns a handle; it does not prove the database is
reachable. PingContext performs that startup check. The sql.DB handle is safe for concurrent use and
manages a pool. Do not call sql.Open in a handler and do not place *sql.DB in request context.

Pool limits are operational settings, not constants to guess during this lesson. Keep defaults
unless you have a concrete database limit or measurement to configure against.

### Match the operation to the API

| Need | database/sql operation | Required handling |
| --- | --- | --- |
| INSERT/UPDATE/DELETE without returned rows | ExecContext | Check returned error and, when behavior depends on it, affected-row count. |
| One row | QueryRowContext | Call Scan; classify sql.ErrNoRows there. |
| Many rows | QueryContext | Check query error, defer rows.Close, loop Next/Scan, then check rows.Err. |
| Several atomic operations | BeginTx plus Tx methods | Defer Rollback, use tx for every operation, check Commit. |

Do not use SELECT * in store code. An explicit column list makes the Scan order reviewable and keeps
an unrelated schema addition from silently changing the mapping.

### Parameterization and identifiers

Safe value flow:

~~~go
row := db.QueryRowContext(ctx, `
    SELECT id, owner_id, title, duration_minutes, created_at
    FROM entries
    WHERE id = $1
`, id)
~~~

The query text is fixed. id is sent as a value for $1. The driver and PostgreSQL do not reinterpret
that value as SQL syntax.

Placeholders are not a general template system. A client-selected sort column cannot become ORDER
BY $1 with the intended identifier semantics. If a future API allows sorting, validate against a
small server-owned allowlist and choose fixed query text. Session 8 keeps ID ASC fixed.

### Multi-row lifecycle

Every QueryContext path must make these steps visible:

~~~text
query with ctx and parameters
  -> check immediate query error
  -> defer rows.Close()
  -> for rows.Next(): Scan explicit columns
  -> check rows.Err() after the loop
  -> return mapped application values
~~~

Close releases the result's resources and allows its connection to return cleanly to the pool.
rows.Err catches errors that occur during iteration after the query began. One does not replace the
other.

### Error translation

The PostgreSQL store owns database-specific interpretation; HTTP remains unaware of PostgreSQL.

| Database observation | Store/application meaning | HTTP behavior |
| --- | --- | --- |
| sql.ErrNoRows for an accessible member query | Established not-found category | Existing 404/access policy |
| Context canceled/deadline exceeded | Request stopped or timed out | Existing safe cancellation/server behavior |
| Known unique/constraint conflict the API exposes | Stable conflict or validation category if deliberately mapped | Existing safe client error chosen by contract |
| Connection, scan, syntax, or unknown constraint failure | Wrapped internal error for logs/caller | Safe 5xx; no SQL, DSN, driver, or constraint text |

Wrap errors with operation context for developers, but preserve errors.Is behavior where the caller
needs classification. Never return raw PostgreSQL errors in JSON.

---

## 5. Inspect TraceLog version 2.3

Start from the completed Session 7 project. Preserve its handler/service boundaries and add the
smallest database-specific edge:

~~~text
tracelog-session-08/
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
  internal/postgres/
    store.go
    store_integration_test.go
  migrations/
    001_create_tracelog.sql
  api-contract.md
~~~

Your actual files may differ. Keep the current layout if its dependency direction is already clear.
The memory store may remain for fast service/handler tests; it is now a real substitution, not a
reason to grow a generic repository hierarchy.

### Required preserved behavior

- GET /healthz remains public and unchanged.
- Authentication remains the clearly labeled teaching stub from Session 4.
- Owner ID still comes from the verified principal, never JSON or SQL input supplied by a client.
- Title, duration, and tag normalization/validation remain the Session 5 contract.
- Handlers still translate HTTP; the service still owns application behavior; the store owns SQL.
- GET /entries keeps page=1, limit=20, maximum limit 100, ID-ascending order, owner scope, total, and
  totalPages from Session 7.
- A valid empty page remains 200 with an empty data array.
- Existing client-visible error shapes remain safe and stable.
- Restarting the Go process no longer erases successfully committed Entries.

### Code-reading exercise

Before editing, trace the Session 7 code and mark the exact replacement boundary:

1. Where is the Store interface consumed and why does its current shape fit or fail the SQL query?
2. Which memory-store behavior is accidental, such as map ordering, and which is contractual?
3. Where does the handler parse and bound page/limit?
4. Where is owner identity established and passed explicitly to the service/store?
5. Which function must receive context.Context for cancellation to reach PostgreSQL?
6. Where will the one *sql.DB be created and closed?
7. Which tests should remain unit tests using memory or a narrow fake?
8. Which behavior can only be proved by a PostgreSQL integration test?
9. Does the current Create store method provide enough information to insert tags in one
   transaction without importing HTTP types?
10. How will a database-generated ID return to the existing response contract?
11. How will sql.ErrNoRows become the application's existing not-found error?
12. Which SQL query preserves Session 7 owner scope before count, order, limit, and offset?
13. Can tags be loaded without one query per Entry? What bounded strategy will you use?
14. Which current test would catch a database implementation that changes the API response shape?

### Prompt to generate the inspection baseline

> Create a new tracelog-session-08 directory by copying my completed Session 7 TraceLog project. Preserve every route, JSON shape, status, validation rule, ownership rule, pagination default/bound, ordering rule, and test. Add only: an internal/postgres package skeleton, one SQL migration location, a PostgreSQL integration-test file with no completed store behavior, and the pgx/v5/stdlib database/sql driver dependency after showing me the exact go.mod impact. Keep the memory store for existing fast tests. First show the preserved-contract list, dependency direction, schema invariants, proposed migration, query inventory, and test split. Wait for my review before writing files. Do not implement queries or transactions, add an ORM/query generator, change the API, or add cache/queue/search infrastructure.

---

## 6. Manual coding and test-first persistence work

### Required vertical change

Implement this yourself, with guidance:

> Persist one authenticated user's Entry and normalized ordered tags atomically, then return an
> owner-scoped page from PostgreSQL with the exact Session 7 contract.

This exercises the full path without creating a second application design:

~~~text
POST /entries
  -> existing decode/prepare/auth/service path
  -> BeginTx
  -> INSERT entry ... RETURNING id, created_at
  -> INSERT/locate normalized tags through the same tx
  -> INSERT entry_tags with position through the same tx
  -> Commit
  -> existing response mapping

GET /entries?page=&limit=
  -> existing bounded parsing and verified owner
  -> COUNT(*) WHERE owner_id = $1
  -> SELECT explicit columns WHERE owner_id = $1
       ORDER BY id ASC LIMIT $2 OFFSET $3
  -> bounded tag retrieval without per-entry queries
  -> existing pagination response
~~~

### Test plan: write these before implementation

| Test | What it proves |
| --- | --- |
| Migration from blank database | All required tables, keys, and constraints can be reproduced. |
| Valid create round trip | Entry, generated ID/time, normalized tags, positions, and owner persist and read back. |
| Failed tag insert rolls back Entry | A lower-level tag constraint failure leaves neither partial Entry nor memberships. |
| Missing principal row | Foreign key rejects an orphan owner; application returns no partial Entry. |
| Get missing member | sql.ErrNoRows becomes the established not-found category. |
| Owner-scoped first page | Only the caller's Entries appear, in ID-ascending order, with correct limit. |
| Later and empty page | OFFSET math is correct and a valid beyond-end page remains an empty success. |
| Owner-scoped total | Another user's rows do not affect total or totalPages. |
| Canceled context | A canceled context reaches the database call and no successful mutation is reported. |
| Constraint defense | An invalid direct store value cannot bypass the durable duration/tag/title rules. |
| Existing handler suite | PostgreSQL work did not change client-visible Session 7 behavior. |

Run integration tests only against a dedicated test database. Give each test unique user IDs and
data, and clean up only what that test owns. Parallel database tests are optional; deterministic
isolation is more important than shaving seconds from this session.

### First failing transaction test

Use a deliberate store-level value that passes far enough to insert the Entry but makes a later tag
write violate a database constraint. Then assert:

~~~text
Create returned an error
AND no matching Entry row exists
AND no membership row exists
~~~

This test intentionally exercises the database boundary below normal input preparation. Its purpose
is not to accept invalid client data; it proves that an unexpected later write cannot leave durable
partial state.

### Transaction shape to understand

~~~go
tx, err := db.BeginTx(ctx, nil)
if err != nil {
    return Entry{}, fmt.Errorf("begin create entry: %w", err)
}
defer tx.Rollback()

// Every Entry and tag write uses tx, never db.

if err := tx.Commit(); err != nil {
    return Entry{}, fmt.Errorf("commit create entry: %w", err)
}
~~~

Deferred Rollback is safe after a successful Commit; it returns sql.ErrTxDone and has no effect.
Keep the transaction short. Do not perform HTTP calls, wait for user input, or launch background
work inside it.

A frequent atomicity bug is starting tx, inserting the Entry with tx, then inserting tags with db.
Those tag writes are outside the transaction even though the function visually contains a
transaction. Review the receiver of every Exec/Query call.

### Parameterized list query

Write the values separately:

~~~sql
SELECT id, owner_id, title, duration_minutes, created_at
FROM entries
WHERE owner_id = $1
ORDER BY id ASC
LIMIT $2 OFFSET $3;
~~~

For validated 1-based inputs:

~~~text
offset = (page - 1) * limit
~~~

The handler already validates the public values, but the application/store boundary should receive
a bounded list request it can trust by construction or validate defensively. Never interpolate
owner, limit, offset, filter, or sort text into SQL.

Count from the same owner-scoped set:

~~~sql
SELECT count(*)
FROM entries
WHERE owner_id = $1;
~~~

Two statements are acceptable for this course. Under concurrent writes, total and page data can be
observed at slightly different moments. Add snapshot coordination only if the product requires that
stronger promise; do not hide the trade-off.

### Avoid the N+1 tag query

Do not run one tag SELECT for every Entry on a page. Use one bounded second query for all Entry IDs
on the page, or one reviewed join/aggregation query. Prefer the option you can explain and test.

For a bounded second query, preserve these properties:

- it runs only when the page contains Entries;
- its ID set comes from the already owner-scoped page, not raw client text;
- it returns entry_id, tag_name, and position in deterministic order;
- Go attaches tags to the correct Entry without changing Entry order; and
- it remains a fixed number of queries for a page, not one query per row.

Do not add a generic data loader to solve one bounded relationship.

### Implementation order

1. Copy the completed Session 7 project and run its existing tests unchanged.
2. Write the invariant table and migration. Review every constraint and delete action.
3. Apply the migration to a blank dedicated test database and inspect it with psql or TablePlus.
4. Add the one PostgreSQL driver dependency and explain its purpose from go.mod.
5. Open and PingContext one sql.DB at startup; inject it at the composition root.
6. Write the failing integration test for rollback after a later tag failure.
7. Implement Entry-plus-tags creation using BeginTx, deferred Rollback, only tx operations, and
   checked Commit.
8. Write failing owner-scoped list tests for page one, a later page, empty page, total, and other-user
   exclusion.
9. Implement the fixed count and page queries with context and PostgreSQL placeholders.
10. Load tags in one bounded additional query or one reviewed join, preserving position.
11. Implement member lookup and classify sql.ErrNoRows at the storage boundary.
12. Add the canceled-context integration test.
13. Run go fmt ./..., focused unit tests, focused integration tests, then go test ./....
14. Restart the server and prove one committed Entry remains.
15. Run EXPLAIN on the exact list query and record only what the plan actually shows.
16. Inspect git diff and explain every schema, dependency, query, error, and lifecycle change.

### Decisions you must make explicitly

| Decision | Your answer before code |
| --- | --- |
| Dedicated test database name/DSN guard | |
| users-to-entries delete behavior | |
| Entry title database bound | |
| Tag membership and position constraints | |
| Migration apply and rollback command | |
| Where sql.DB is created and closed | |
| Where database-generated ID/time are scanned | |
| Store error representing not found | |
| List count and page query order | |
| Tag-loading strategy without N+1 | |
| Transaction invariant and exact statements inside it | |
| Integration-test cleanup/isolation strategy | |
| Candidate index and the query that justifies it | |

### Focused implementation prompt

> Guide me through Session 8 PostgreSQL persistence without writing the whole store. First require my schema invariants, relationship diagram, migration review, query inventory, test-database safety check, parameterized list SQL, and failing rollback test. Then help one step at a time: open/Ping one database/sql pool with pgx/v5/stdlib, implement the Entry-plus-tags transaction using only tx calls, implement the owner-scoped count/page QueryContext path with rows cleanup, classify sql.ErrNoRows, retrieve tags without N+1 queries, and run focused integration tests. Preserve every Session 7 API behavior. Do not add an ORM, query builder, generator, generic repository, cache, queue, search service, retry layer, trigger, or speculative pool/index tuning.

### Feature addition: tag filtering with pagination

Only after the required persistence tests pass, add the study-plan feature:

~~~text
GET /entries?tag=go&page=1&limit=20
~~~

Define and test before SQL:

1. tag uses the same normalization as create input.
2. owner scope still applies before count and page selection.
3. total counts only Entries matching both owner and normalized tag.
4. each Entry appears once even if joins could multiply rows.
5. result order remains Entry ID ascending.
6. an unknown tag returns 200 with an empty data array.
7. the SQL still uses value parameters.

An EXISTS predicate is often clearer than joining only to test membership. Write and explain the
query yourself. Inspect whether entry_tags already has a useful key order for this lookup before
adding an index. Do not add a cache; Session 9 begins by asking whether measured repeated work
justifies one.

---

## 7. Intentional database bug lab

Choose one bug only after the baseline integration tests pass. I should provide the symptom and
reproduction, not the source location or solution.

| Bug | Faulty behavior | Observable consequence |
| --- | --- | --- |
| String-built SQL | owner or tag is concatenated into query text. | Injection risk, quoting defects, and unreviewable query shapes. |
| New pool per request | Handler calls sql.Open for each request. | Connection growth, lifecycle confusion, and avoidable failures under concurrency. |
| sql.Open treated as connectivity proof | Startup never calls PingContext. | The app appears healthy until its first real query fails. |
| Context discarded | Store uses context.Background instead of the received context. | Canceled clients leave database work running. |
| rows not closed | QueryContext result is returned from without Close. | Connections remain occupied and later requests wait or fail. |
| rows.Err ignored | Iteration failure after some rows is treated as success. | A partial page can be returned as complete. |
| SELECT * scan drift | A migration adds/reorders a selected column. | Scan fails or fields map incorrectly. |
| Mixed tx/db writes | Entry uses tx but tags use db. | A tag failure can leave a committed or independently visible partial Entry. |
| Rollback omitted | An early return leaves the transaction open. | A connection remains tied up and no clear atomic outcome exists. |
| Commit error ignored | Function returns success after a failed commit. | Client believes data persisted when it did not. |
| Count lacks owner filter | Page rows are scoped, total is global. | Cross-user metadata leaks and totalPages is wrong. |
| Pagination lacks ORDER BY | LIMIT/OFFSET operate on unspecified order. | Entries repeat or move between pages unpredictably. |
| Wrong offset formula | page * limit is used for 1-based pages. | Page one skips the first page of data. |
| N+1 tag loading | One tag query runs per listed Entry. | Query count grows with page size. |
| Constraint text leaked | Raw PostgreSQL error is encoded to JSON. | Internal schema and implementation details reach clients. |
| Test points at shared database | Cleanup truncates or drops unverified state. | Data outside the test can be destroyed. |
| Speculative index | An index is added without its query or plan. | Writes/storage cost increases with no demonstrated benefit. |

### Debugging protocol

1. Reproduce the symptom with the smallest PostgreSQL integration test or pool statistic. Do not
   start by rewriting the store.
2. State the violated invariant: safety, atomicity, ownership, cancellation, completeness, ordering,
   resource release, or error secrecy.
3. Trace the request context, parameters, DB/Tx receiver, result rows, error mapping, and cleanup.
4. Inspect PostgreSQL state directly after the failure. Do not infer rollback from the Go return
   value alone.
5. Write a regression test that fails before the fix.
6. Find the earliest shared decision point where the invariant is lost.
7. Make the smallest fix there.
8. Re-run the focused test, all integration tests, all Go tests, and one manual SQL observation.
9. Inspect the diff for a new leak, hidden retry, broad abstraction, or changed public contract.

### Bug-lab prompt

> I intentionally introduced a Session 8 PostgreSQL bug. Be a debugging partner, not an auto-fixer. Ask for the failing request/test, exact query, schema constraint, DB/Tx receiver, context path, and database state after failure. Make me state the invariant and write a focused integration test. Help me narrow the earliest root cause. Do not propose code until I explain why the failure occurs and what the smallest shared fix is.

---

## 8. Audit AI-generated schema and storage code

Generated persistence code can pass happy-path handler tests while corrupting invariants, leaking
connections, or weakening authorization. Review the migration, query text, Go lifecycle, tests, and
client contract together.

### Audit rubric

| Review question | Evidence of a good answer |
| --- | --- |
| Does the schema encode real invariants? | Types, nullability, checks, unique keys, and foreign keys match established TraceLog rules. |
| Is ownership structural and query-scoped? | owner_id is constrained and every member/list/count query uses the verified owner policy. |
| Are values parameterized? | SQL text is fixed and client/application values are separate arguments. |
| Is ordering explicit and stable? | Page SQL has ORDER BY id ASC before LIMIT/OFFSET. |
| Is database/sql used with the right lifecycle? | One startup pool, PingContext, injection, shutdown Close, no per-request sql.Open. |
| Does cancellation reach I/O? | Received context flows to BeginTx/ExecContext/QueryRowContext/QueryContext. |
| Are multi-row resources released? | Query error, defer rows.Close, Scan errors, and rows.Err are all handled. |
| Is the transaction truly atomic? | Every related write uses tx; Rollback is deferred; Commit error is checked. |
| Are errors classified safely? | Not found/conflict/cancellation are deliberate; raw driver/schema details never reach JSON. |
| Does pagination preserve Session 7? | Bounds, owner scope, order, empty page, total, and totalPages match the documented contract. |
| Are tags loaded with bounded queries? | No per-Entry tag query; tag order and empty tags remain stable. |
| Are indexes query-driven? | Each non-constraint index names the real WHERE/JOIN/ORDER BY path and has plan/measurement evidence appropriate to the stage. |
| Are integration tests isolated? | Dedicated test DSN, unique fixtures, narrow cleanup, no shared database destruction. |
| Is complexity proportionate? | One driver and existing boundaries; no ORM/generator/generic repository/retry/cache added. |

### Common AI-generated-code failure modes

- Replacing database/sql with native pgx throughout the app even though the course and current
  boundary already chose database/sql.
- Adding GORM, sqlx, sqlc, squirrel, a repository base class, or a unit-of-work abstraction before
  one hand-written query is understood.
- Opening or pinging a database inside every handler.
- Storing a DSN, password, or test credential in committed Go code, logs, or error responses.
- Using fmt.Sprintf or concatenation for owner, tag, limit, offset, or filter values.
- Treating placeholders as safe dynamic identifiers.
- Using context.Background in storage because it is convenient.
- Forgetting rows.Close, rows.Err, Scan errors, deferred Rollback, or Commit errors.
- Calling db.ExecContext from inside a transaction function that should call tx.ExecContext.
- Returning success after a partial Entry/tag write.
- Relying only on application validation and omitting durable constraints.
- Returning raw unique/foreign-key/check-constraint errors to clients.
- Letting SQL pagination count all users or filter ownership after LIMIT/OFFSET.
- Omitting ORDER BY or changing ID-ascending order because a query happened to look sorted.
- Loading tags with one query per Entry.
- Using SELECT * and positional Scan across an evolving table.
- Adding indexes to every foreign key and filter candidate without inspecting actual queries.
- Using triggers for ordinary visible application behavior that is easier to trace in Go.
- Running integration tests against the developer's default database and using TRUNCATE CASCADE for
  cleanup.
- Changing the public API response while claiming the work was only a storage swap.

### Persistence-diff audit prompt

> Review this generated Session 8 TraceLog persistence diff before I accept it. First list the preserved Session 7 API behavior and the durable schema invariants. Then trace startup pool creation, one create transaction, member lookup, owner-scoped count/page query, tag loading, cancellation, rows cleanup, error translation, and shutdown. Inspect every SQL statement for fixed text, parameters, explicit columns, ordering, ownership, cardinality, and index assumptions. Verify all related writes use tx and Commit is checked. Review migration safety and integration-test isolation. Cite exact evidence and classify findings as must-fix, design decision, or safe simplification. Flag ORMs, generators, generic repositories, retries, caches, triggers, or tuning that lack a current requirement. Do not rewrite the implementation unless I ask.

### Query inventory prompt

> Extract every SQL statement in this Session 8 diff into a table with caller, purpose, expected row count, parameters, owner/auth scope, ordering, transaction membership, context source, returned columns, error mapping, cleanup, supporting constraint/index, and integration test. Flag any row where the code or test provides no evidence. Do not suggest a new library.

### Transaction review prompt

> Trace this Entry-plus-tags transaction line by line. Identify the invariant, when a connection is acquired, every statement receiver (db or tx), every possible early return, deferred rollback behavior, commit handling, context cancellation, and database state after each failure. Then name the smallest missing regression test. Do not propose nested transactions, retries, or a unit-of-work abstraction.

### Migration review prompt

> Compare this PostgreSQL migration with the Session 4–7 TraceLog contracts. Check type choices, nullability, primary/foreign/unique/check constraints, delete actions, normalized tag rules, first-seen tag order, defaults, migration ordering, destructive down behavior, and query-driven indexes. List mismatches and unsupported additions. Do not add columns or indexes for hypothetical future features.

---

## 9. Checkpoint: model, query, transact, and defend

Complete this without notes first.

1. Draw the users, entries, tags, and entry_tags relationships. Label primary keys, foreign keys,
   and cardinalities.
2. For each rule, state whether it belongs in input preparation, service behavior, database schema,
   or more than one. Explain why.

| Rule | Boundary and reason |
| --- | --- |
| Trim title whitespace | |
| Duration is 1–480 | |
| Owner comes from verified principal | |
| owner_id references an existing user | |
| One normalized tag appears once per Entry | |
| Tag response order preserves first occurrence | |
| Client sees no raw PostgreSQL error | |

3. Explain why sql.DB is a pool handle and where TraceLog creates, shares, and closes it.
4. Explain why sql.Open and PingContext are separate steps.
5. Choose ExecContext, QueryRowContext, or QueryContext for each operation:

| Operation | Choice and required cleanup/error step |
| --- | --- |
| INSERT ... RETURNING id | |
| SELECT one Entry by ID | |
| SELECT a page of Entries | |
| DELETE one Entry without returned columns | |

6. Write the owner-scoped page query with PostgreSQL placeholders and compute the offset for page 3,
   limit 20.
7. Trace a canceled HTTP request through handler, service, store, driver, and database call.
8. Explain why QueryRowContext reports not found at Scan and how the store preserves that category.
9. Write the full sql.Rows lifecycle in order and explain the separate jobs of Close and Err.
10. Draw the Entry-plus-tags transaction. Mark every failure point and expected database state.
11. Explain how a function can begin a transaction correctly but accidentally execute writes
    outside it.
12. State which integration test proves there is no partial Entry after a tag failure.
13. Explain why owner scope must appear in both count and page SQL.
14. Describe one bounded tag-loading strategy that avoids N+1 queries.
15. Given an EXPLAIN plan on five rows, state what you can observe and what performance conclusion
    you cannot support.
16. Review a proposal to add an ORM, automatic retries, eight indexes, and Redis during this storage
    change. Accept, reject, or defer each from current evidence.
17. Explain why a migration and a passing unit test are both insufficient without a PostgreSQL
    integration test.
18. Prove the API contract survived the storage swap using one handler test and one manual request.

### Passing standard

You are ready for Session 9 when you can:

- Model the current TraceLog data and justify every constraint.
- Apply and inspect a migration against a dedicated test database.
- Explain one sql.DB pool lifecycle from startup through shutdown.
- Write and trace a parameterized, context-aware PostgreSQL query.
- Handle QueryRow and Query/Rows resources and errors correctly.
- Preserve owner scope, stable order, and bounded pagination in SQL.
- Implement and prove the Entry-plus-tags atomic transaction.
- Classify not-found and internal database errors without leaking implementation detail.
- Use one integration test to prove rollback and one to prove owner-scoped pagination.
- Read a basic EXPLAIN plan without inventing a performance conclusion.
- Reject persistence abstractions and infrastructure that do not solve a current measured problem.

If any answer is weak, repeat one narrow SQL path and its integration test. Do not hide uncertainty
behind an ORM, generator, retry loop, cache, or broad rewrite.

---

## 10. Update your progress log

~~~markdown
## Session 8 — PostgreSQL persistence

- Date and actual time:
- Diagnostic answer I changed:
- TraceLog relationship diagram:
- Migration applied and inspected:
- Durable constraints I can justify:
- Parameterized query I wrote manually:
- sql.DB and context lifecycle explanation:
- Transaction invariant:
- Rollback integration-test evidence:
- Owner-scoped pagination test evidence:
- Tag-loading strategy and query count:
- EXPLAIN observation and its limit:
- Chosen bug, root cause, and regression check:
- AI-generated persistence finding:
- Dependency, index, trigger, or abstraction I rejected and why:
- Confidence (1–5):
- Question to carry into Session 9:
~~~

---

## Ready for Session 9

Session 9 introduces caching and background jobs. Bring:

- the Session 7 API contract and passing handler tests;
- the Session 8 migration and schema diagram;
- the exact owner-scoped count, page, and tag queries;
- the transaction rollback integration test;
- one EXPLAIN plan or measured query observation; and
- one candidate piece of repeated or deferrable work, with evidence rather than intuition.

A cache is meaningful only after you can name the query it avoids, the staleness it introduces, and
the invalidation event. A background job is meaningful only after you can name the work, ownership,
retry/idempotency behavior, and failure policy. The PostgreSQL implementation you can now trace is
the baseline for both decisions.
