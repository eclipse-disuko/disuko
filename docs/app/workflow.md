# How it works

A typical project's journey through DISUKO, from first setup to a compliant release.

## 1. Set up the project

A project owner creates a **project** in DISUKO to represent the piece of software being tracked. If several related products share a supplier or a platform, they can be organized under a **project group** instead of managed one by one.

## 2. Grant access

The project owner decides who can contribute data: an internal team, an external supplier, or both. Each gets either a login to the web app or an API token, scoped to that project (or the whole group).

## 3. Create a version

Every release or build gets its own **project version** in DISUKO. This keeps each release's compliance status independent — an older version can stay in whatever state it finished in while a new one goes through review from scratch.

## 4. Generate the SBOM

The team or supplier uses whatever tooling they already rely on (a scanner, a build-integrated generator, a tool like [ORT](https://oss-review-toolkit.org/)) to produce a Software Bill of Materials — a full list of every third-party component in the software, and its license. DISUKO doesn't dictate the tool, only the schema the output has to match (based on the SPDX standard), which any project can request at any time.

## 5. Upload it

The SBOM is submitted against the specific project version, either by a person through the web app or automatically by a CI pipeline through the API.

## 6. Automatic validation

DISUKO checks the file against the required schema the moment it arrives. A malformed file is flagged immediately; a valid one moves the version into review.

## 7. Automatic policy check

Every component's license is checked against the policy rules that apply to this project — the allow/deny list, and the classifications (obligations, prohibitions, and so on) tied to each license. Anything unknown, unapproved, or forbidden is surfaced right away, instead of being discovered during an audit months later.

## 8. Review

The project owner or an assigned reviewer opens the project's dashboard and sees everything they need: the component list, validation results, any policy flags, and reviewer remarks raised so far. Remarks can be worked through and resolved individually or in bulk before a final decision is made.

## 9. Approve or reject

The reviewer makes the call. This creates a linked approval task: the reviewer sees it as something to do, the person who submitted the SBOM sees it as pending, and both are tied to the same decision so there's no ambiguity about where things stand. A decision can carry a comment explaining it, and once made, it's locked and written to the audit log at the project, version, and SBOM level.

- **Rejected** → the reason is visible to whoever submitted it, who fixes the issue and resubmits, going back to step 5.
- **Approved** → the version's compliance status is now official.

## 10. Generate the paperwork

Once a version is approved, DISUKO can generate the **third-party notice** file directly from it — the attribution document required by many open-source licenses — ready to ship with the product or hand to legal, with no manual assembly required.

## 11. Everything stays available

The approved status, the SBOM, and the notice file can all be queried later — by a person checking history, or by another system via the API — long after the release has shipped. Nothing has to be re-derived from scratch the next time someone asks "what's in this product, and is it compliant?"

## 12. Repeat for the next release

As the software evolves, steps 3–10 repeat for every new version. Because DISUKO keeps the full history, later deliveries can be compared against earlier ones, making it easy to see exactly what changed from one release to the next.

---

Want the vocabulary behind these steps spelled out? See the [glossary](/glossary).
