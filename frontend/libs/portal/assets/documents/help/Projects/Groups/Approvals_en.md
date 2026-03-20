# Approvals

## Overview

- In this view you find all approvals and reviews, which have been requested on the project
- Only users with the Owner role can request approvals and reviews
- Requesting an approval or review will automatically add a Viewer role for the requested users if new to the project
- Requesting an approval will automatically freeze and store the applying policy rules and policy check status next to the approval documentation
- Requesting an approval will automatically generate a German and English FOSS Disclosure Document for signing including SBOM and project metadata and pointers
- Before requesting an approval, you can mark the respective SBOM in the deliveries tab of a version with a star <i aria-hidden="true" class="v-icon notranslate mdi mdi-star"></i>
- When a project is child of a project group, then approvals are managed on project group level
- You can only approve one project version at a time
- In project groups one SBOM of each child project can be selected, e.g. to request approval for frontend and backend at once

## Internal Approvals

- Internal approvals require two different persons according to the four-eyes principle
- Internal approvals will create approval tasks within the Disclosure Portal
- At first the supplier approval tasks are created in parallel, after completion the customer approvals will be created
- An approval can be aborted by the requester on the task list
- Hash Chains: For an extra validation all files belonging to the approval process chain are stored and can be downloaded or requested on demand. For each file the SHA256 hash can be calculated and compared with the corresponding value available in data base. The hash values stored in data base can be accessed using the "copy reference" Action on the Disclosure Document.
