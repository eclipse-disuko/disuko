# SBOM Compare

## Overview:

- **Selection:** select the current SBOM to be taken as basis for the comparison, and a previous SBOM delivery from the same project. You can select an SBOM delivery from other project versions. The SBOM deliveries are compared for equal, changed, added and removed components. The resulting table lists the union set of both SBOM deliveries.
- **Table:** the list has the same setup as the components list-. A click on a component shows the component's attributes and identified differences.
- **Differences:** the difference status is shown for each component:
  - <i class="v-icon notranslate material-icons mdi mdi-plus"></i> means the component is only present in the current SBOM
  - <i class="v-icon notranslate mr-2 material-icons mdi mdi-minus"></i> means the component is only present in the previous SBOM
  - <i class="v-icon notranslate material-icons mdi mdi-compare-horizontal"></i> means the component information differs in current and previous SBOM- Use the details view to display the attributes and highlights the identified differences
