# Freigaben

## Übersicht

- In dieser Ansicht finden Sie alle Freigaben und Prüfungen, die für das Projekt angefordert wurden
- Nur Benutzer mit der Rolle "Eigentümer" können Freigaben und Überprüfungen anfordern
- Wenn Sie eine Freigabe oder Prüfung anfordern, wird automatisch eine Betrachter-Rolle für die angeforderten Benutzer hinzugefügt, wenn sie neu im Projekt sind
- Wenn Sie eine Freigabe anfordern, werden die angewendeten Policy-Regeln und der Status der Policy-Regel-Prüfung automatisch eingefroren und neben der Freigabedokumentation gespeichert
- Durch die angeforderte Freigabe wird automatisch ein deutsches und englisches FOSS-Offenlegungsdokument zur Unterzeichnung generiert, welches SBOM- und Projektmetadaten sowie entsprechende Zeiger enthält
- Bevor Sie eine Freigabe anfordern, können Sie die jeweilige SBOM im SBOM-Lieferungen-Tab einer Version mit einem Stern markieren <i aria-hidden="true" class="v-icon notranslate mdi mdi-star"></i>
- Wenn ein Projekt einer Projektgruppe zugeordnet ist, werden Freigaben auf Projektgruppenebene verwaltet
- Sie können jeweils nur eine Projektversion freigeben oder prüfen lassen
- In Projektgruppen kann jeweils eine SBOM je enthaltenem Projekt ausgewählt werden, um z.B. die Freigabe für Frontend und Backend gleichzeitig anzufordern

## Interne Freigaben

- Für interne Freigaben ist das 4 Augen Prinzip zu beachten
- Interne Freigaben erzeugen zu Freigabeaufgaben innerhalb des Disclosure Portals
- Zunächst werden parallel die Lieferantenfreigabeaufgaben erstellt, nach Abschluss werden die Kundenfreigabeaufgaben erstellt
- Eine Freigabe kann vom Anforderer in der Aufgabenliste abgebrochen werden
- Hash-Ketten: Für eine zusätzliche Validierung werden alle zur Genehmigungsprozess-Kette gehörenden Dateien gespeichert und können bei Bedarf heruntergeladen oder angefordert werden. Für jede Datei kann der SHA256-Hash berechnet und mit dem entsprechenden in der Datenbank gespeicherten Wert verglichen werden. Die in der Datenbank gespeicherten Hash-Werte können über die Aktion "Referenz kopieren" auf dem Disclosure Dokument abgerufen werden.
