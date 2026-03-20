# SBOM Deliveries Management

## Overview:

- **Table:** shows a list of all SBOM deliveries along with some metadata on the origin; an SBOM can either be uploaded manually in this view in form of an SPDX JSON file, or provided via API. After upload, the SBOM is validated against the schema and only accepted if valid.
- **Actions:** you may mark an SBOM for approvals or checks here (just click on the <i aria-hidden="true" class="v-icon notranslate mdi mdi-star"></i> star), copy the SBOM reference information into your clipboard, or download an SBOM file by clicking on a specific row, you can open the SBOM in COMPONENTS view and inspect the contents there.
- **Validation:** when uploading an SBOM to Disclosure Portal, it is automatically validated against a JSON schema, which defines the input format requirements. If validation fails, please contact the technical support. Validation includes the following remarks:
  - The section in the SBOM in which a problem item is found. Typically this is "packages", where the components are listed
  - The position of the item in the list of items. For example, the first item in the list of packages (index starts with "0")
  - The attribute of the item which caused a problem - for example the "name" attribute
  - The problem itself such as something is required but not found, or something is unexpected
- **Important:**
  - Regarding SBOM deletion and retention: The SBOMs which were subject to an approval or review will be retained. Furthermore, the five most recently delivered SBOMs will be kept. You can manually lock an SBOM to exclude it from automated deletion. The deletion of limit-reached, unlocked SBOMs will be triggered through a successful SBOM upload to the project version.

## Validation Problems:

- **Required Attributes:** some attributes, which are necessary for proper processing in Disclosure Portal and have been set as "required" in the SBOM schema. You can find out which of the attributes are mandatory in the "requiredAttributes" statements. In case there is no value known for a specifically required attribute (like for 'versionInfo'), an empty value ("" or "NOASSERTION") will satisfy the schema validation, as long as the attribute itself is contained.
- **Additional Attributes:** for some cases, the schema allows additional, non-standard attributes to be present. The switch additionalProperties indicates if additional properties are allowed (true) or not (false).
- **Case Sensitive:** JSON Schema is case sensitive, so you might run into a schema validation problem in case the tool used for generating the SBOM does not use the right upper/lower case spelling. It needs to match exactly. SPDX License IDs are not validated in the schema validation step - License IDs are treated as case insensitive.
- **Where to fix problems:** The identified problems need to be addressed in the tool which generates the SBOM (e.g. Black Duck), or in a preprocessor step before submitting an SBOM to Disclosure Portal.

## Limits:

- **Maximum Uploads** per project per hour: 30
- **Maximum File Size** per upload: 200 MB
