# Session 4 — Authentication, Authorization, and Trust Boundaries

[← Course overview](../study-plan.md) · [← Session 3](session-03-routing-and-json.md)

## Session contract

**Time:** 100–120 minutes  
**Lecture sequence:** L08 — Authentication and authorization  
**Project:** TraceLog version 1.5 — protected entries with ownership rules  
**Goal:** Learn to trace a caller’s identity from an HTTP request to an authorization decision, and recognize the difference between a teaching stub and a production authentication system.

Session 3 established route selection and JSON boundaries. This session adds a verified caller identity and a small policy: a user can access their own entries, while an administrator may access any entry. The important learning outcome is not building a login system; it is being able to audit where identity comes from, what is trusted, and whether permission is checked for the actual resource.

### Safety and scope

The code exercise uses test-only opaque tokens mapped to preconfigured principals. It is an identity-verification stub, not a deployment-ready authentication system.

Do not build:

- A password registration or login endpoint.
- Password storage, hashing parameters, password reset, MFA, JWT issuance, OAuth, OIDC, cookie sessions, or SSO.
- A real token secret in source code.
- A custom cryptographic algorithm.
- A broad authorization framework, policy language, or role hierarchy.
- An authentication bypass such as trusting an X-User-ID header supplied by the client.

Those systems need explicit product requirements, reviewed libraries or a mature identity provider, secure configuration, threat modeling, and ongoing maintenance. Today we study the execution flow and test an authorization policy safely.

---

## Source material

Read L08 in order. The supporting book material is useful for recognizing implementation patterns, but it is not the authority for current production security choices.

| Primary material | What to focus on |
| --- | --- |
| [L08 — Authentication and authorization](../resources/backend-from-first-principles/notes/08-8-authentication-and-authorization-for-backend-engineers.md) | Authentication versus authorization; stateful sessions and stateless tokens; password handling; OAuth/OIDC; roles and permissions. |
| [L08 transcript](../resources/backend-from-first-principles/transcripts/08-A95rliroC8Q.en.vtt) | Revisit the session, JWT, OAuth/OIDC, RBAC, and password-storage explanations when a term remains unclear. |
| *Let’s Go* §§9.1–9.3 and §§11.1–11.6 | Supporting reference for recognizing Go web-session and authorization patterns only. Verify unfamiliar security behavior against current documentation and a reviewed design. |

---

## Outcomes

By the end of this session, you should be able to:

- State precisely that authentication identifies a caller and authorization decides whether that caller may perform an action on a resource.
- Trace a protected request from an Authorization header to a verified principal in request context to a resource-specific authorization decision.
- Recognize that headers, JSON bodies, path values, and query values are client-controlled until verification or validation establishes otherwise.
- Explain why an entry owner must come from the verified principal, not from a client-supplied owner field.
- Distinguish 401 for missing or invalid authentication from 403 for an authenticated caller lacking permission.
- Recognize when a system may intentionally return 404 instead of 403 to avoid revealing resource existence, and why that is an explicit policy decision.
- Explain the high-level trade-offs between server-side sessions, bearer tokens, JWTs, cookies, OAuth, and OpenID Connect.
- Recognize that passwords require deliberately designed password hashing and must never be stored or logged in plaintext.
- Test authentication and authorization paths separately.
- Audit AI-generated Go code for identity spoofing, missing ownership checks, credential leaks, insecure token handling, and overbuilt security plumbing.

You are not expected to memorize JWT claim rules, OAuth grant details, cookie attribute matrices, cryptographic algorithms, password-hashing parameters, or identity-provider APIs. You must know which questions to ask, which behavior to test, and when to stop an agent from improvising security code.

---

## Active practice and tool lab

### Source-aware tool register

| Tool or technology | Source and role | Depth for this session | Proof of competence |
| --- | --- | --- | --- |
| Authorization header plus curl or Postman | Course-supporting way to exercise L08’s protected-request flow. | Hands-on. | Send a missing, malformed, known, and unknown test token; explain the 401 or success result without exposing token text. |
| JWT | L08 discusses JWTs as one bearer-token design. | Safe inspection only. | Decode a fabricated local token-shaped sample and explain that decoding claims is not signature verification. Never paste real tokens into online tools. |
| Cookies and browser storage inspection | L08 discusses session and cookie-based identity. | Working familiarity. | Identify Secure, HttpOnly, and SameSite as attributes to inspect on a disposable local test cookie; do not build a login system. |
| Redis | L08 presents it as a possible session store. | Recognition now; hands-on caching work is deferred to the caching lecture. | Draw where session lookup would occur and explain its availability and invalidation trade-off. |
| OAuth and OpenID Connect | L08 presents external identity/provider flows. | Recognition only. | Trace who authenticates the user, who receives an assertion, and why this API must still authorize resources. Do not create a provider account or app. |

No production authentication library, key store, password hasher, or identity-provider SDK belongs in
this session. The test-only opaque token stub exists so you can learn trust boundaries without
pretending to implement real security.

### Manual-first rule

Write the authentication and authorization tests before the wrapper or policy code. Manually add one
ownership test or policy branch, then explain every line that reads the header, establishes a
Principal, stores or retrieves context data, assigns owner ID, or makes an allow/deny decision.
Security code is never accepted because it “looks standard.”

### Tool drill

1. Exercise the protected routes with a non-secret test token through curl or Postman. Capture only
   status and safe response data in your notes.
2. Compare 401 for missing/invalid identity with the chosen 403-or-404 policy for a verified
   non-owner. State the attacker-controlled value and the protected property for each test.
3. Decode only a fabricated local JWT-shaped string with base64 or a small local Go snippet. Confirm
   that readable JSON is not authenticated data until a verifier checks its signature and claims.
4. Sketch the session, JWT, and OAuth/OIDC flows from the lecture; mark where TraceLog would still
   perform resource authorization.

### Completion evidence

You complete the session only with focused tests for missing credentials, invalid credentials,
owner access, non-owner access, admin access, and spoofed owner input. Your tool evidence must not
contain raw authorization values, secrets, or copied real tokens.

---

## Agenda

| Time | Activity | Evidence you produce |
| --- | --- | --- |
| 0–10 min | Diagnostic and Session 3 recall | Short written answers |
| 10–30 min | L08 model and terminology | Identity and access decision map |
| 30–45 min | Read the protected-request flow | Annotated request trace |
| 45–65 min | Inspect TraceLog version 1.5 | Tests for identity and ownership |
| 65–82 min | Test-first authorization feature | A resource-policy change |
| 82–97 min | Intentional security bug lab | Root cause and regression test |
| 97–112 min | AI-generated-code audit | Evidence-based security review |
| 112–120 min | Checkpoint and progress log | Session 5 readiness |

---

## 1. Start with a diagnostic

Answer without searching first.

1. What is the difference between authentication and authorization?
2. Is a valid bearer token proof that a caller may read any resource? Why or why not?
3. Where should an Entry owner identifier come from when a user creates an entry?
4. Which response category fits a missing or invalid credential: 401, 403, or 404?
5. Which response category fits a valid caller who tries to read another user’s private entry?
6. Why is an X-User-ID request header not a trustworthy identity mechanism by itself?
7. What information should never appear in logs, error messages, test snapshots, or source control?
8. What is a session at a high level? What is a token at a high level?
9. Why is “we use JWT” not enough information to decide that an authentication design is safe?

Record your answers. A weak answer is a map for the rest of the session.

### Start-of-session prompt for Codex

> Teach Session 4 from sessions/session-04-authentication-and-authorization.md. Use L08 as the primary source. Ask my diagnostic questions one at a time and make me separate authentication, authorization, ownership, and input validation. Use only the clearly labeled test-only token stub described here. Do not build login, password handling, JWT issuance, OAuth, sessions, cookies, or custom cryptography.

---

## 2. The identity-and-permission model

### The two questions

~~~text
authentication: “Who is making this request?”
authorization: “May that verified caller do this action on this resource?”
~~~

They are consecutive but independent decisions.

A request can be:

| Caller state | Resource policy result | Correct broad outcome |
| --- | --- | --- |
| No usable identity | Not evaluated | 401 authentication failure |
| Identity is verified | Policy allows the action | Success path |
| Identity is verified | Policy denies the action | 403, or a deliberate 404 policy |
| Identity is verified | Resource does not exist | 404 |
| JSON body is malformed | Authentication may not even reach behavior | 400 input failure |

### TraceLog protected-request flow

~~~text
HTTP request
  -> route selects a protected handler
  -> requirePrincipal reads Authorization input
  -> verifier maps a known test token to a Principal
  -> verified Principal is attached to this request context
  -> handler reads path/body values
  -> handler loads the requested Entry
  -> policy compares principal and entry owner or role
  -> handler returns success, 401, 403, or 404
~~~

The key point is where trust changes:

~~~text
Authorization header text: untrusted
token after successful verifier check: evidence for a principal
Principal in request context: trusted only for this request
owner ID in client JSON: untrusted and must not choose ownership
entry loaded from the server’s store: server-controlled state
policy decision: explicit server behavior
~~~

### Authentication is not authorization

A request with a valid identity can still be forbidden.

For example:

~~~text
Principal: user-2
Requested entry owner: user-1
Requested action: read entry

authentication: succeeds
authorization: fails unless an explicit administrator policy allows it
~~~

Never collapse this into “token present means allowed.” That is the beginning of horizontal privilege escalation.

### The minimal policy for this session

Use a small, explainable policy:

~~~text
A normal user may list and read only their own entries.
A normal user creates entries owned by themselves.
An administrator may read any entry.
The client never chooses or overrides owner ID.
~~~

A policy can be expressed as a small function or clearly named branch. Do not create an authorization engine to evaluate one ownership rule.

---

## 3. L08 concepts: understand, recognize, and look up

### Understand deeply

| Topic | What you should be able to explain |
| --- | --- |
| Authentication | Verifies an identity claim before application behavior trusts it. |
| Authorization | Uses verified identity, action, resource, and policy to decide access. |
| Ownership | A resource has an owner controlled by server-side behavior, not client input. |
| Trust boundary | Anything from the network is attacker-controlled until it is verified or validated. |
| 401 versus 403 | 401 means no valid authenticated caller; 403 means valid caller but policy denial. |
| Context propagation | A verified principal can travel with one request through request context; raw credentials should not. |
| Resource-level checks | Permission must be tested against the specific entry, not only at route registration. |
| Security tests | Test absent credential, invalid credential, permitted access, forbidden access, and anti-spoofing behavior. |

### Recognize and reason about

| Topic | What to recognize |
| --- | --- |
| Stateful sessions | The server stores or can look up session state; a browser often presents an opaque session identifier. Revocation and server-side control are easier, but storage is required. |
| Bearer tokens | Possession is enough to use them, so leakage is serious. Their format does not itself decide safety. |
| JWT | A signed token format, not an authorization policy and not automatically encrypted. Correct validation includes more than merely decoding it. |
| Cookies | Browser-managed credential transport with security attributes and CSRF considerations. A cookie is not inherently authentication. |
| OAuth 2.0 | Delegated authorization between a client, resource owner, authorization server, and resource server. |
| OpenID Connect | An identity layer built on OAuth-style flows; often used for sign-in. |
| RBAC | Role-based access control, such as user and admin roles. Roles still must be applied to concrete actions and resources. |
| ABAC or policy rules | Attribute-based decisions can be useful when simple roles and ownership stop fitting. |
| Password hashing | Passwords require one-way, password-specific hashing through a vetted approach; plaintext and fast general-purpose hashes are not acceptable designs. |
| MFA, recovery, revocation, rate limiting | Real identity systems require lifecycle and abuse controls beyond a token check. |

### Look up immediately before implementation

- The current recommended password-hashing approach and its parameters.
- Exact cookie attributes such as Secure, HttpOnly, and SameSite.
- JWT issuer, audience, expiration, key rotation, signature, and algorithm-validation requirements.
- OAuth/OIDC flow choice, redirect URI validation, scopes, provider configuration, and token storage.
- CSRF protections for cookie-based browser sessions.
- Identity-provider library documentation, secret-management rules, and production threat-model requirements.
- The project-specific choice between 403 and resource-hiding 404 behavior.

Do not guess any of these because you saw a similar implementation in generated code.

---

## 4. Security mental models before code

### Passwords

A password is a secret supplied by a user. The server should not need to recover the original password later.

~~~text
never: store plaintext password
never: log plaintext password
never: return password data in an API response
never: invent password hashing or crypto
instead: use a reviewed password-specific hashing design when a real identity system is in scope
~~~

This session deliberately does not process passwords. Treat any agent that introduces them unasked as expanding scope and security risk.

### Sessions, tokens, and cookies

| Mechanism | Useful mental model | Important caution |
| --- | --- | --- |
| Server-side session | Client presents an opaque identifier; server looks up session state. | Session fixation, storage, revocation, cookie attributes, and CSRF need design. |
| Opaque bearer token | Client presents a secret; server or gateway verifies it against state. | Anyone who obtains it can use it until the system rejects it. |
| JWT | Client presents signed claims; verifier checks signature and required claims. | Decoding is not verification; revocation and key lifecycle require deliberate design. |
| Cookie transport | Browser sends a cookie under matching rules. | Browser convenience creates cross-site and theft considerations; set attributes intentionally. |
| OAuth/OIDC | A specialist authorization or identity system supplies verified results. | Do not copy a redirect/token flow without provider requirements and reviewed docs. |

### CORS is not auth

CORS tells browsers when JavaScript from one origin may read a response. It does not prevent a non-browser client from sending a request, identify a caller, or decide whether they may access an Entry.

### A security review question set

For every protected endpoint, ask:

1. What is the identity source?
2. Who verifies it?
3. Is the raw credential ever logged, returned, or stored?
4. Which user-controlled values are still untrusted?
5. Which resource is being accessed?
6. What rule allows or denies that exact action?
7. What status and body are revealed on failure?
8. Which test proves a caller cannot cross a tenant or ownership boundary?

---

## 5. Inspect TraceLog version 1.5

Create a new directory named tracelog-session-04 by evolving the Session 3 API.

~~~text
tracelog-session-04/
  go.mod
  main.go
  server.go
  entry.go
  auth.go
  server_test.go
~~~

### Required behavior

- GET /healthz remains public.
- GET /entries, POST /entries, and GET /entries/{id} require a verified principal.
- A missing or malformed Authorization credential returns 401 and does not run protected behavior.
- The test-only verifier maps known opaque test tokens to principals with an ID and role.
- POST /entries assigns owner ID from the verified principal. The request JSON has no owner field.
- A normal user lists and reads only their own entries.
- An administrator can read entries across ownership boundaries.
- A normal user attempting to read another user’s known entry receives the policy response chosen and documented for this exercise.
- Tests exercise the actual returned handler through httptest.

### Teaching-only identity stub

The verifier may use an injected map of test token strings to Principals. The test values are fixtures, not deployment secrets.

~~~text
Authorization: Bearer test-user-one
  -> verifier
  -> Principal{ID: "user-1", Role: "user"}

Authorization: Bearer test-admin
  -> verifier
  -> Principal{ID: "admin-1", Role: "admin"}
~~~

The only job of this stub is to create a trustworthy Principal after a successful lookup. It must not be copied to a production service.

### Minimal code shape to recognize

~~~text
requirePrincipal(next handler)
  -> reads Authorization header
  -> rejects missing or invalid credential
  -> obtains verified Principal from verifier
  -> creates a request context that carries Principal
  -> calls next with the updated request

protected handler
  -> reads Principal from request context
  -> parses route/body input
  -> loads Entry
  -> evaluates ownership or role policy
  -> writes response
~~~

### Context rule

A request context is appropriate for request-scoped identity that has already been verified. Use a private typed key and a small accessor function. Do not use ordinary strings as context keys, do not put raw bearer tokens into context, and do not use context as a replacement for explicit function parameters or persistent storage.

This is an early, narrow exposure to middleware-shaped code and request context. Session 6 will teach package boundaries, middleware, and context in greater depth.

### Code-reading exercise

Before running any tests, trace three requests.

| Request | Questions to answer |
| --- | --- |
| GET /entries without Authorization | Where does authentication stop the request? What must not happen? |
| POST /entries with a user token | Where does owner ID originate? What must the JSON decoder not control? |
| GET /entries/7 with another user’s token | Where is the Entry loaded? Where is ownership compared? What response policy applies? |

Then identify:

1. The type that represents a verified caller.
2. The exact function that verifies a raw credential.
3. The exact function that attaches identity to the request.
4. The exact function that reads identity back out.
5. The exact policy that grants or denies a resource action.
6. Which values should never appear in a response or log.
7. Which routes remain public and why.
8. Whether an invalid token can cause an Entry lookup or write.
9. Whether client JSON can set owner ID.
10. Whether the administrator rule is a role rule, an ownership rule, or both.

### Prompt to generate the inspection code

> Create a new tracelog-session-04 directory by evolving the Session 3 TraceLog API. Add a clearly labeled test-only opaque-token verifier that maps injected test tokens to a Principal with ID and role. Protect GET /entries, POST /entries, and GET /entries/{id}; keep GET /healthz public. Store Entry ownership server-side from the verified Principal, never from JSON. Normal users may access only their own entries; an admin may read any entry. Use one small handler wrapper and typed request-context accessors, concrete in-memory storage, and httptest tests. Before writing code, show the trust-boundary map, route policy table, status decisions, file tree, and tests. Do not build login, passwords, JWT, OAuth/OIDC, cookies, sessions, database layers, a router framework, policy engine, or real secrets.

---

## 6. Manual coding and test-first feature work

Start with authentication tests before writing or changing the policy.

### Authentication contract

| Request condition | Expected behavior |
| --- | --- |
| No Authorization header | 401; protected handler must not list, read, or create data. |
| Malformed Authorization scheme | 401; generic safe error behavior. |
| Unknown test token | 401; generic safe error behavior. |
| Known test token | Principal is available only to the protected request path. |
| GET /healthz without token | 200; health remains deliberately public. |

### Authorization contract

For this lesson, make the policy explicit:

| Verified principal | Requested Entry | Expected behavior |
| --- | --- | --- |
| owner user-1 | Entry owned by user-1 | Allow |
| ordinary user-2 | Entry owned by user-1 | Deny with the documented policy response |
| admin | Entry owned by user-1 | Allow |
| any user | New entry creation | Server assigns that user as owner |

### Your implementation order

1. Write a test for missing authentication on GET /entries.
2. Write a test for invalid authentication on POST /entries.
3. Write a test proving a user-created Entry belongs to the verified principal, even though JSON cannot supply owner ID.
4. Write a test for an owner reading their Entry.
5. Write a test for a different ordinary user reading it.
6. Write a test for the administrator rule.
7. Implement the smallest verifier and handler wrapper that make the authentication tests pass.
8. Implement the smallest ownership policy that makes the authorization tests pass.
9. Format and run all tests.
10. Explain why identity verification and permission checking belong in different places.

### Feature addition: resource ownership

Implement the following change yourself:

> A caller must not be able to create an Entry on behalf of another user.

Start by writing a request body that tries to include an owner field. Your test should prove that this field is either rejected by the defined JSON contract or ignored because server-side ownership always comes from Principal. Then confirm the stored Entry owner matches Principal.ID.

### Decision exercise: 403 or 404?

For an authenticated ordinary user requesting another user’s entry, choose a policy:

- **403:** clearly communicates that identity is valid but access is prohibited.
- **404:** hides whether the resource exists, which can be useful for sensitive resources.

For TraceLog, choose one, record it, and test it consistently. Explain the trade-off. Do not let an agent silently make the choice.

### Focused implementation prompt

> Guide me through adding ownership authorization to my Session 4 TraceLog API. First ask me to write tests for missing credentials, invalid credentials, owner access, non-owner access, admin access, and spoofed owner input. Then help me make the smallest changes: verified Principal, request-context accessor, one auth wrapper, server-assigned owner, and one resource policy. Do not introduce login, JWT, passwords, cookies, database interfaces, or a generic policy framework.

---

## 7. Intentional security bug lab

Choose one bug only after baseline tests pass. Your diagnosis should name the violated security property, not just the bad line.

| Bug | Faulty behavior | Security consequence |
| --- | --- | --- |
| Client-owned owner ID | POST body includes ownerId and storage trusts it. | A caller can create or impersonate ownership for another user. |
| Authentication without authorization | Any verified token may read any Entry by ID. | Horizontal privilege escalation. |
| Trusting X-User-ID | Handler treats a client header as a verified principal. | Anyone can impersonate a user. |
| Missing return after 401 | Wrapper sends an authentication failure but still calls the protected handler. | Unauthorized behavior can still execute. |
| Raw credential leak | Error path, log, or test output includes Authorization value. | Bearer credential disclosure. |
| Role from JSON | Client body or query decides whether caller is admin. | Privilege escalation. |
| Context key collision | A shared string context key is used inconsistently. | Wrong data may be interpreted as identity. |
| 200 on denial | Handler returns a success status with an error-like body. | Clients and monitoring cannot reliably see policy failure. |

### Debugging and security-audit protocol

1. Reproduce with a focused httptest case.
2. State the attacker capability: what can they control?
3. State the protected property: identity, ownership, secrecy, or policy enforcement.
4. Trace the data from HTTP boundary to the first trust decision.
5. Find the earliest place the program treats untrusted data as trusted.
6. Add a regression test that fails before the fix.
7. Make the smallest root-cause fix at the shared decision point.
8. Verify no state change occurred on denied requests.
9. Run formatting and all tests.
10. Write one sentence describing what the test now prevents.

### Questions I will ask before helping

- What raw values can the caller send?
- What source is supposed to establish identity?
- What resource and action need authorization?
- Is the problem authentication, authorization, validation, or response disclosure?
- Which code path makes the first trust decision?
- What test proves that a denied request cannot change or reveal state?

### Bug-lab prompt

> I intentionally introduced a Session 4 auth or authorization bug. Be a security-focused debugging partner, not an auto-fixer. Ask me for the failing request/test and the smallest relevant auth, context, and handler code. Help me identify the attacker-controlled value, the violated security property, and the earliest incorrect trust decision. Require a regression test and do not propose a patch until I explain the root cause.

---

## 8. Audit AI-generated identity and access code

Generated security code needs more skepticism than ordinary plumbing. A clean-looking diff can still allow cross-user access or mishandle credentials.

### Audit rubric

| Review question | Evidence of a good answer |
| --- | --- |
| Is authentication distinct from authorization? | Token verification happens before the handler; resource policy is still checked after lookup. |
| Is identity sourced safely? | No client-controlled user or role header/body field is trusted as identity. |
| Is ownership server-assigned? | Creation uses Principal.ID; request JSON cannot choose another owner. |
| Are resource checks specific? | The policy compares the principal against the exact Entry or a deliberate admin rule. |
| Are credentials protected? | Authorization values do not appear in logs, responses, fixtures intended for sharing, or source configuration. |
| Are failures truthful and safe? | 401, 403 or deliberate 404, and 500 are not collapsed into 200. |
| Is context used narrowly? | A typed key and accessor carry a verified Principal for one request only. |
| Are tests adversarial? | Tests prove missing, invalid, non-owner, spoofed input, and allowed admin behavior. |
| Is the solution proportionate? | A test stub and small wrapper, not an improvised identity platform. |
| Does it defer high-risk work? | Passwords, JWT, sessions, OAuth/OIDC, cookies, and cryptography are explicitly out of scope. |

### Common AI-generated-code failure modes

- Decoding a JWT but never verifying its signature or critical claims.
- Treating a bearer token as encrypted just because it looks opaque.
- Trusting X-User-ID, X-Role, ownerId, or admin flags sent by the caller.
- Checking only that a user is logged in, not that they own the requested resource.
- Returning detailed credential errors that reveal whether an account exists.
- Logging Authorization headers, cookies, passwords, reset links, or tokens.
- Storing passwords as plaintext, reversible encryption, or a fast general-purpose hash.
- Creating a custom token format or custom cryptographic helper.
- Adding permissive CORS headers as a substitute for authorization.
- Returning 200 for authentication or authorization failure.
- Adding a session/JWT/OAuth library without requirements, configuration, key lifecycle, or a security review.
- Burying policy in a generic helper nobody can trace or test.

### Audit prompt

> Review this generated Session 4 Go diff as a security-sensitive change. Trace identity from request header through verification and request context to resource-level authorization. Find impersonation paths, cross-user access, role spoofing, secret leakage, incorrect status behavior, response disclosure, missing adversarial tests, unsafe crypto or token assumptions, and unnecessary auth infrastructure. Cite exact code and classify each finding as must-fix, design decision requiring product input, or optional simplification. Do not write a replacement implementation unless I request it.

### Assumption-challenge prompt

> List every security assumption in this TraceLog auth change: token source, verifier trust, token lifecycle, transport security, role assignment, owner assignment, context propagation, resource policy, error disclosure, logging, and tests. Mark each as enforced, tested, documented, or merely assumed. Identify which assumptions make this a teaching stub rather than a production authentication design.

---

## 9. Checkpoint: trace, decide, and defend

Complete this without notes first.

1. State authentication and authorization in one precise sentence each.
2. Trace GET /entries/7 with a valid ordinary-user token from HTTP header to final response.
3. Explain why a verified user can still receive a forbidden or hidden result.
4. Explain why owner ID must come from Principal, not JSON.
5. Place these in order:

~~~text
route selection
identity verification
JSON decoding
resource lookup
authorization policy
response encoding
~~~

Then explain where some ordering can vary safely and where it cannot.

6. Predict the expected category of result:

| Request | Expected category |
| --- | --- |
| GET /healthz without Authorization | |
| GET /entries without Authorization | |
| GET /entries with an unknown token | |
| POST /entries with user-1 token and ownerId set to user-2 | |
| GET /entries/7 as entry owner | |
| GET /entries/7 as another ordinary user | |
| GET /entries/7 as admin | |

7. Write a test that proves a client cannot choose an Entry owner.
8. Review an agent proposal to add JWT login. Reject, defer, or accept it with a security and scope reason.
9. Explain one fact about sessions, bearer tokens, JWTs, cookies, OAuth, or password hashing that you would verify in current documentation before implementation.

### Passing standard

You are ready for Session 5 when you can:

- Trace identity and permission checks as distinct steps.
- Prove with a test that an ordinary user cannot cross an ownership boundary.
- Explain exactly why the teaching token stub is not production authentication.
- Identify at least three generated-code security failure modes.
- Choose and document a response policy for unauthorized resource access.
- State which security details must be looked up or reviewed instead of improvised.

If you cannot pass one item, repeat the matching test or audit exercise. Do not add a more sophisticated identity system to compensate for unclear fundamentals.

---

## 10. Update your progress log

~~~markdown
## Session 4 — Authentication and authorization

- Date and actual time:
- Diagnostic answer I changed:
- My authentication versus authorization explanation:
- Trust-boundary map:
- TraceLog folder or commit:
- Chosen non-owner response policy and rationale:
- Chosen bug, attacker capability, root cause, and regression test:
- AI audit finding:
- Security detail I will verify before real implementation:
- Confidence (1–5):
- Question to carry into Session 5:
~~~

---

## Ready for Session 5

Session 5 separates structural decoding from semantic validation and transformation. Keep this request flow in mind:

~~~text
route -> authenticate -> deserialize -> validate and transform -> authorize action -> apply behavior -> serialize response
~~~

The exact order can depend on the operation, but do not use validation as authentication or authentication as authorization.

Bring your Session 4 code and prepare to answer:

- Which input values still need validation after authentication succeeds?
- Where should trimming, lowercasing, deduplication, and length limits happen?
- How will you prove invalid input cannot become stored state?

> Teach Session 5 only after I show you the Session 4 checkpoint answers and one ownership test I can explain line by line.
