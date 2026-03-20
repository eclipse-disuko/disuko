# SBOM Vergleichen

## Überblick:

- **Auswahl:** Wählen Sie die aktuelle SBOM, die als Grundlage für den Vergleich verwendet werden soll, und eine frühere SBOM-Lieferung aus demselben Projekt aus. Sie können eine SBOM-Lieferung aus anderen Projektversionen auswählen. Die SBOM-Lieferungen werden auf gleiche, geänderte, hinzugefügte und entfernte Komponenten verglichen. In der resultierenden Tabelle ist der Vereinigungssatz beider Stücklistenlieferungen aufgeführt.
- **Tabelle:** Die Liste hat den gleichen Aufbau wie die Komponentenliste. Ein Klick auf eine Komponente zeigt die Attribute der Komponente und identifizierte Unterschiede.
- **Unterschiede:** Der Differenzstatus wird für jede Komponente angezeigt:
  - <i class="v-icon notranslate material-icons mdi mdi-plus"></i> bedeutet, dass die Komponente nur in der aktuellen Stückliste vorhanden ist
  - <i class="v-icon notranslate mr-2 material-icons mdi mdi-minus"></i> bedeutet, dass die Komponente nur in der vorherigen SBOM vorhanden ist
  - <i class="v-icon notranslate material-icons mdi mdi-compare-horizontal"></i> bedeutet, dass sich die Komponenteninformationen in der aktuellen und vorherigen SBOM unterscheiden. Verwenden Sie die Detailansicht, um die Attribute anzuzeigen und die identifizierten Unterschiede hervorzuheben
