# DISUKO – focuses on consuming SBOMs and the resulting actions based on their assessment.

DISUKO is an open-source project under the umbrella of the [Eclipse Foundation](https://projects.eclipse.org/projects/technology.disuko).  
It provides a modular and lightweight base to quickly start working with **Disuko** functionalities.  
The goal is to deliver a ready-to-run entry point with minimal setup effort.

---

## Features

- Docker-based setup (via `docker-compose`)
- Ready-to-run demo environment
- Includes example users and credentials
- Extendable for custom requirements
- Supports SBOM (Software Bill of Materials) integration

---

## Quickstart

Run the following command in the project root directory:

```bash
./setup-dev.sh   # Windows: setup-dev.ps1
```

```bash
docker-compose up --build -d
```

Check if all containers are running:

```bash
docker-compose ps --format "{{.Service}} {{.State}}"
```

### Open in browser

[https://localhost:3009/](https://localhost:3009/)

### Credentials

```
Username: CUSTOMER1
Password: CUSTOMER1
```
```
Username: CUSTOMER2
Password: CUSTOMER2
```

### Troubleshooting

- If something goes wrong (e.g., login issues), try logging out first:  
  [Logout User](https://localhost:3009/api/v1/oauth/logout)

- For the setup wizard, if an owner or company name is required, you may use "dummy" as value.

---

## SBOM Support

DISUKO supports uploading Software Bill of Materials (SBOMs) after successfully creating a project.  
Before uploading an SBOM, you must first upload an SBOM schema under **Admin** with the label `common standard`.

The official SPDX schema can be downloaded here:  
[SPDX 2.3 Schema (JSON)](https://github.com/spdx/spdx-spec/blob/support/2.3/schemas/spdx-schema.json)

---

## Next Steps

- Integrate your own configurations and data sources
- Enable additional modules and extensions
- Experiment with SBOM uploads for project transparency and compliance

---

## Contributing

Contributions are welcome and appreciated.

This project follows the Eclipse Foundation development and contribution processes.
Before contributing, please make sure you are familiar with the following resources:

- Eclipse Foundation Contributing Guide  
  https://www.eclipse.org/contribute/

- Eclipse Contributor Agreement (ECA)  
  https://www.eclipse.org/legal/ECA.php

By submitting a pull request, you confirm that you have the right to contribute the code
and that you agree to the terms of the Eclipse Contributor Agreement.

No additional project specific contribution guidelines are required at this time.

### Security

This project provides a Gitleaks configuration file to help contributors
detect accidental secret commits. Usage is optional and can be integrated
locally or in CI environments.

## Code of Conduct

This project follows the Eclipse Foundation Code of Conduct to ensure a respectful,
inclusive, and harassment free environment for everyone involved.

All participants are expected to adhere to the rules defined here:  
https://www.eclipse.org/org/documents/Community_Code_of_Conduct.php

By participating in this project, you agree to uphold this Code of Conduct in all project related spaces.

## License

This project is licensed under the [Apache-2.0](LICENSE).

## Note

The installation variants provided serve exclusively as templates for test environments. Although they are ready for immediate use, they must be adapted to the specific requirements of the target environment before going live. This includes, in particular, additional hardening and security measures.
