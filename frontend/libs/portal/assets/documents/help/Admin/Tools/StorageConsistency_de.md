## How to fix storage inconsistency

The feature to correct storage inconsistency is off by default. To enable it, set the environment variable "
ENABLE_FIX_DATA_INTEGRITY=true" in any of the values-XXX.yaml files located in the infra/charts/disuko-backend directory of
the disuko-backend helm charts.

This feature cannot be activated through a simple frontend button, but must be initiated by a REST call. Corresponding
call templates are located in the http/analyse.http file.
