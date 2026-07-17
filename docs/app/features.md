# Features

A tour of what DISUKO actually does, grouped by the job each part handles. If a term is unfamiliar, check the [glossary](/glossary).

## Projects, versions, and groups

Everything in DISUKO is organized around **projects** — the applications, services, or products your organization ships — and **versions** of those projects, representing a specific release or build.

- A project moves through a simple lifecycle (new → active → archived) as it gets its first delivery and, eventually, is retired.
- Each version has its own compliance status, so an older release and a newer one can be at completely different stages of review at the same time.
- **Project groups** let you organize related projects under one parent — useful when one supplier or one platform produces several related products, since access and API tokens can be shared across the whole group.

## SBOM intake and validation

Suppliers and teams don't have to guess what format to submit. DISUKO defines the **schema** a valid SBOM must match (based on the SPDX standard), and any project can request the current schema before generating its file.

- SBOMs are uploaded — through the UI or the API — against a specific project version.
- Every upload is automatically validated against the schema; malformed files are flagged immediately instead of silently causing problems later.
- Because the intake is schema-driven, DISUKO can evolve what it requires over time without breaking older integrations, since older schema versions stay available until retired.

## License and policy engine

This is the core of what makes DISUKO more than a filing cabinet: it doesn't just store what licenses a project uses — it judges them.

- Administrators maintain a central **license database** and a set of **policy rules**: which licenses are acceptable for which kind of project, and which are not.
- **Calculated policy rules** let you define policy at the level of a whole category (a "bucket" of similar licenses) instead of listing every license by name one at a time.
- Every component in an SBOM is checked against the policy that applies to its project, and anything unknown, unapproved, or outright forbidden is surfaced clearly — no digging required.
- A **classification matrix** gives reviewers a single cross-reference view of policy rules against classifications, so it's obvious at a glance what's allowed, what needs a closer look, and what's denied.

## Classifications (obligations)

Every license can carry real legal duties — DISUKO calls these **classifications** (for example: an obligation to publish attribution, a right, a limitation, or an outright prohibition). Classifications are managed centrally by administrators and are what drives the automatic policy checks described above, so the rules stay consistent across every project instead of being reinvented per team.

## Review and approval workflow

Compliance isn't just a green checkmark — it's a decision someone has to own. DISUKO models that explicitly:

- Reviewers can leave **remarks** on scans, licenses, or general findings, and track them through to resolution — individually or in bulk.
- An **approval workflow** turns a review into a real task: the approver gets it on their to-do list, the requester sees it as pending, and both are linked so nobody has to chase the other down over email.
- Every approval or rejection — along with any comment attached to it — is written to an **audit log** at the project, version, and SBOM level, so there's always a record of who decided what, and when.
- A dedicated **Tasks** view gives every reviewer a single place to see everything waiting on their decision, and tasks can be delegated when someone is unavailable.

## Notices and attribution

Once a version is approved, DISUKO can generate the **third-party notice** file — the attribution document that lists every component, its license, and its copyright notice — in plain text, HTML, or JSON, ready to ship alongside the product or hand to legal.

## Roles, permissions, and access

Access in DISUKO is deliberately fine-grained rather than all-or-nothing:

- Project-level roles (owner, supplier, viewer, and similar) control who can do what on a given project.
- System-wide administrative roles (for licenses, policies, classifications, schemas, and more) let you split responsibility for reference data across the people who actually own that domain, instead of handing every admin the same master key.
- Both human users (via single sign-on) and automated systems (via personal access tokens or project-scoped API tokens) can be granted exactly the access they need — nothing more.

## Admin and governance tools

Behind the scenes, an administration area gives the people running a DISUKO instance the controls they need to keep it healthy and up to date:

- Manage SBOM schemas, license data, and policy rules as living reference data, versioned and auditable.
- Configure mail templates and send test notifications, manage in-app announcements and a site-wide notification bar.
- Toggle features on or off, manage translations, and review scheduled background jobs.
- Keep the Terms of Use current, and manage data retention for user accounts.

## API access

Nothing in DISUKO is UI-only. A public, versioned API lets external systems and CI/CD pipelines authenticate with a project-scoped token and:

- Upload an SBOM for a specific project version as part of an automated build.
- Check the current compliance status of a project or version.
- Retrieve the generated notice file or policy rules for a project.

A separate internal API supports trusted automation within your own organization. This means DISUKO can sit quietly inside a release pipeline, gating or reporting on compliance without anyone having to open the web app at all.

---

Curious how these pieces fit together in practice? Read [how it works](/workflow) for a step-by-step walkthrough.
