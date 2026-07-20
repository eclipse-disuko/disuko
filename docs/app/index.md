---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

hero:
  name: "DISUKO"
  text: "Open-source license compliance, without the paperwork"
  tagline: Track every open-source component in your software, check it against your policies, and prove compliance — automatically, not on a spreadsheet.
  actions:
    - theme: brand
      text: What is DISUKO?
      link: /introduction
    - theme: alt
      text: Explore the features
      link: /features
    - theme: alt
      text: Run it locally
      link: /local-setup

features:
  - icon: 📋
    title: SBOM intake, built in
    details: Upload a Software Bill of Materials (based on the SPDX standard) for any project and DISUKO validates it against a schema automatically — no more guessing whether a file is usable.
  - icon: ⚖️
    title: Automatic license policy checks
    details: Every component's license is checked against the allow/deny rules you define, so problems surface long before a release ships.
  - icon: ✅
    title: A clear approval workflow
    details: Reviewers get a to-do list, suppliers get status updates, and every approval or rejection is written to an audit log.
  - icon: 📄
    title: One-click attribution files
    details: Generate the third-party NOTICE file your legal team needs — in text, HTML, or JSON — straight from an approved SBOM.
  - icon: 🗂️
    title: Modeled on how software is really built
    details: Projects, versions, groups, and role-based permissions reflect the reality of internal teams and external suppliers working together.
  - icon: 🔌
    title: API-first
    details: Everything the UI can do is backed by an API, so CI pipelines and other systems can upload, check, and query compliance status automatically.
---
