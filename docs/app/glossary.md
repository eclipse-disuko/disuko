# Glossary

Plain-language definitions for the terms used throughout these docs and inside DISUKO itself.

**SBOM (Software Bill of Materials)**
A full inventory of every third-party or open-source component used inside a piece of software — like an ingredients list for code, including each component's version and license.

**SPDX**
An industry-standard file format for expressing an SBOM's contents, maintained under the Linux Foundation. DISUKO's SBOM schemas are based on SPDX.

**Project**
A tracked piece of software in DISUKO — an application, service, or other deliverable whose open-source components and compliance status are managed over time.

**Project version**
A specific release or build of a project. Each version has its own SBOM delivery and its own independent compliance status.

**Group**
A parent that organizes several related projects together, mainly so access and API tokens can be shared across all of them at once.

**Schema**
The formal definition of what fields and structure an uploaded SBOM must contain to be considered valid. Schemas are versioned, so what's required can evolve without breaking older integrations.

**Policy rule**
A rule defining which licenses are acceptable ("allowed") or unacceptable ("denied") for a given project, used to automatically flag non-compliant components. A **calculated policy rule** applies to a whole category of similar licenses at once, rather than listing each one individually.

**Classification**
The specific legal duty or characteristic tied to a license — for example, an obligation to publish attribution, a right, a limitation, a prohibition, or an exception. Classifications are what drive the automatic policy checks.

**Classification matrix**
A single cross-reference view showing, for every combination of policy rule and classification, whether it's allowed, flagged for review, or denied.

**Notice (third-party notice)**
The attribution document — required by many open-source licenses — listing every component, its license, and its copyright notice. DISUKO can generate this automatically from an approved SBOM.

**Task / approval**
A review decision in progress. Submitting an SBOM for review creates a linked pair of tasks: the reviewer sees it as something to action, and the submitter sees it as pending, until a decision is made.

**Audit log**
The permanent record of who approved or rejected what, and when, kept at the project, version, and SBOM level.

**Personal access token / API token**
A credential that lets an automated system (a CI pipeline, another internal system) act on behalf of a user or a specific project without a human logging in.

**Public API / Internal API**
The public API is what external suppliers and CI pipelines use, scoped to a single project via a token. The internal API supports trusted automation within your own organization.
