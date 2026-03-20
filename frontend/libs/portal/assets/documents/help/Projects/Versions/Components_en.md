# Components Management

## <i class="v-icon notranslate material-icons mdi mdi-layers" style="color: var(--v-textColor-base); caret-color: var(--v-textColor-base);"></i> Overview:

- **Name & Version** of all components included in the selected SBOM delivery are listed in the table
- **Type:** A component is typically a package, but can also be a file
- **License Effective:** "License Declared" information of the component, unless a "License Concluded" information is reported
- **Policy:** check result status of component's license according to the policy rules applicable for the project

## <i class="v-icon notranslate material-icons mdi mdi-minus-circle text-error" style=""></i> Policy Status "Denied":

- **"Denied"** is not to be approved in general, therefore should be removed or replaced by the development team.
- **"Denied (forbidden)"** license is forbidden for use in custom developed software without exception (in any distribution model), remove or replace.

## <i class="v-icon notranslate material-icons policyStatusUnassertedColor--text mdi mdi-lightning-bolt-circle" style=""></i> Policy Status "Unasserted":

- **"Unasserted"** means license information for the component is missing or it is unknown in the internal knowledge base; such cases need to be resolved.
- **Unknown:** no license information at all is present for a component (like "NONE", "" or "NOASSERTION"), the use is forbidden; in this case the development team needs to investigate the license information.
- **License Chart Status:** only licenses with status "license chart" may be used; to resolve missing license chart status, check your options:
  - a. license may be known under another identifier, so check the license database and add an alias to the known license entry may resolve this issue.
  - b. request removal or replacement by the development team.
  - c. Set License Chart for the respective license (after finishing option a.).

## <i class="v-icon notranslate material-icons mdi mdi-alert text-warning" style=""></i> Policy Status "Warned":

- **"Warned"** means that the license of a component is on a deny list requiring investigation and a conscious decision.
- Reasons why a license has a warning may vary. For example, a weak copyleft license may be warned because of the risk that the copyleft could unintentionally apply to the project, for instance when the component is statically linked.
- **Resolving:** If a component with a license on a warn list is included, there are some steps:
  - Assess the potential risks of including the component in your project, the license information in the license database will provide guidance for that
  - Check if the component is available in multi licensing with a better license option
  - Check if the supplier can remove or replace the component with another component or functionality with a more suitable license
  - Document the decision and the reason (this documentation is to be kept internal)

## <i class="v-icon notranslate material-icons mdi-help-circle-outline mdi v-theme--dark" style="color: var(--v-textColor-base); caret-color: var(--v-textColor-base);"></i> Policy Status "Questioned":

- More than one license applies to that component due to Multi-Licensing.
- **"Questioned"** is no statement on the policy status, please go to the component's detail view for this.
  Check the policy rules status there, the ”worst case” or “most problematic” license applies.
  - **AND:** in the case of "AND" there is no choice to make, for example on a Dual-Licensed component declaring "MIT AND Apache-2.0", both licenses have to be complied with concurrently.
  - **OR:** in the case of an "OR" statement, we may choose one of these licenses and so we only need to comply with one license. Unless an "OR" is resolved, the component is listed in the license concluded attribute under each of the licenses. To resolve an "OR" (depending on the use case), the following procedure should be applied:
    - Use the most permissive and common license, unless you are forced to use a copyleft license for conflict reasons or compatibility reasons (e.g. for MIT OR GPL-2.0-or-later, choose MIT)
    - Use only licenses with license charts
    - Usage of old components should be avoided where possible
    - Unless an "OR" is not resolved in the license concluded attribute, the component is to be listed under each of the licenses (remember: use only approved licenses)
    - Document the decision and the reason (this documentation is to be kept internal)

## <i class="v-icon notranslate material-icons mdi mdi-check-circle text-green" style=""></i> Policy Status "Allowed":

- **"Allowed"** means, that the reviewed license is allowed for this use case.
- This status is limited to the license assessment, it does not reflect the component's technology, maintenance, code quality or cyber security, these need to be assessed separately. You may need to consider further findings and remarks in the quality check of the delivery.
