# What is DISUKO?

Almost every piece of modern software is built on open-source code. A typical application can easily depend on hundreds — or thousands — of third-party components, each with its own license and its own strings attached: some require you to publish attribution, some require you to offer your own source code in return, and some can't legally be used in certain products at all.

Keeping track of that by hand — in spreadsheets, emails, and tribal knowledge — doesn't scale. Components get missed, obligations get forgotten, and by the time anyone notices, the software has already shipped.

**DISUKO is the system that keeps track of it for you.**

It's a central place where:

- Every project declares exactly what open-source components it uses, via a standard **Software Bill of Materials (SBOM)**.
- Those components are automatically checked against the license rules your organization has agreed to.
- The people responsible for a project can see, at a glance, whether it's compliant — and approve or reject a release accordingly.
- The paperwork that compliance requires (attribution files, audit trails) gets generated for you instead of assembled by hand.

DISUKO doesn't replace the tools that scan your code or generate an SBOM in the first place (a tool like [ORT](https://oss-review-toolkit.org/) does that job well already). Instead, DISUKO is what happens *after* the scan: it consumes the SBOM, checks it, routes it for a decision, and keeps the record of that decision for as long as you need it.

## Who it's for

DISUKO is built around the way software actually gets made — by more than one team, often across organizational boundaries:

- **Project or product owners** create a project in DISUKO, decide who can contribute data to it, and make the final call on whether a release is compliant.
- **Development teams and external suppliers** generate the SBOM for the software they build and submit it — through the web app or directly from a CI pipeline via the API.
- **Compliance, legal, or license reviewers** define the policies (which licenses are acceptable, what obligations they carry) that every project gets checked against, and step in when a decision needs expert judgment.
- **Administrators** manage the underlying reference data — license definitions, SBOM schemas, user access — that keeps the rest of the system consistent.

Because access and roles are scoped per project (and per project group), the same instance of DISUKO can support many teams and many suppliers at once, without everyone seeing everyone else's data.

## Why it exists

DISUKO grew out of years of real-world experience running open-source license compliance for a large software organization, where the same manual, paper-based process was repeated for every product, every release. It's now developed as an open-source project under the [Eclipse Foundation](https://projects.eclipse.org/projects/technology.disuko), so any organization facing the same problem can use it, adapt it, and contribute back.

Ready to see what it actually does? Continue to [Features](/features), or jump straight to [how it works](/workflow) end-to-end.
