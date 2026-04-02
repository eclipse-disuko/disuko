# Verwaltung von SBOM-Lieferungen

## Überblick:

- **Tabelle:** zeigt eine Liste aller SBOM-Lieferungen zusammen mit einigen Metadaten zum Ursprung; Eine SBOM kann entweder manuell in dieser Ansicht in Form einer SPDX-JSON-Datei hochgeladen oder über die API bereitgestellt werden. Nach dem Hochladen wird die SBOM anhand des Schemas validiert und nur akzeptiert, wenn sie gültig ist.
- **Aktionen:** Sie können hier eine SBOM für Freigaben oder Plausibilitätsprüfung markieren (dazu auf den <i aria-hidden="true" class="v-icon notranslate mdi mdi-star"></i> Stern klicken), die SBOM Referenzinformationen in Ihre Zwischenablage kopieren oder eine SBOM Datei herunterladen, indem Sie auf eine bestimmte Zeile klicken, Sie können die SBOM in der Ansicht KOMPONENTEN öffnen und dort den Inhalt überprüfen.
- **Validierung:** Beim Hochladen einer SBOM in das Disclosure Portal wird sie automatisch anhand eines JSON-Schemas validiert, das die Anforderungen an das Eingabeformat definiert. Wenn die Validierung fehlschlägt, wenden Sie sich bitte an den technischen Support. Die Validierung umfasst die folgenden Anmerkungen:
  - Der Abschnitt in der Stückliste, in dem ein problematisches Element gefunden wird. In der Regel handelt es sich dabei um "Pakete", in denen die Komponenten aufgelistet sind
  - Die Position des Elements in der Liste der Elemente. Zum Beispiel das erste Element in der Liste der Pakete (Index beginnt mit "0")
  - Das Attribut des Elements, das ein Problem verursacht hat - z. B. das Attribut "name"
  - Das Problem selbst, z. B. dass etwas benötigt, aber nicht gefunden wird, oder dass etwas unerwartet ist
- **Wichtig:**
  - Zur Löschung und Aufbewahrung von SBOM-Lieferungen: Die genehmigten bzw. geprüften SBOMs bleiben erhalten. Darüber hinaus werden die fünf zuletzt gelieferten SBOMs beibehalten. Sie können eine SBOM manuell sperren, um sie vom automatisierten Löschen auszuschließen. Das Löschen von SBOM-Lieferungen über dem Limit, die zudem nicht gesperrt wurden, wird durch einen erfolgreichen SBOM-Upload in der Projektversion angestoßen.

## Validierungsprobleme:

- **Erforderliche Attribute:** einige Attribute, die für die ordnungsgemäße Verarbeitung im Disclosure Portal erforderlich sind und im SBOM-Schema als "erforderlich" festgelegt wurden. Welche der Attribute obligatorisch sind, erfahren Sie in den "requiredAttributes"-Anweisungen. Falls für ein spezifisch benötigtes Attribut kein Wert bekannt ist (z. B. für 'versionInfo'), erfüllt ein leerer Wert ("" oder "NOASSERTION") die Schemavalidierung, solange das Attribut selbst enthalten ist.
- **Zusätzliche Attribute:** In einigen Fällen lässt das Schema zu, dass zusätzliche, nicht standardmäßige Attribute vorhanden sind. Der Schalter additionalProperties gibt an, ob zusätzliche Eigenschaften zulässig (true) oder nicht (false) sind.
- **Groß-/Kleinschreibung beachten:** Beim JSON-Schema wird zwischen Groß- und Kleinschreibung unterschieden, sodass möglicherweise ein Problem bei der Schemavalidierung auftritt, falls das zum Generieren der Stückliste verwendete Tool nicht die richtige Groß-/Kleinschreibung verwendet. Es muss genau übereinstimmen. SPDX-Lizenz-IDs werden im Schema-Validierungsschritt nicht validiert – bei Lizenz-IDs wird die Groß-/Kleinschreibung nicht beachtet.
- **Wo Probleme behoben werden können:** Die identifizierten Probleme müssen in dem Tool behoben werden, das die SBOM generiert (z. Black Duck) oder in einem Präprozessorschritt vor dem Einreichen einer SBOM an das Disclosure Portal.

## Grenzwerte:

- **Maximale Uploads** pro Projekt und Stunde: 30
- **Maximale Dateigröße** pro Upload: 200 MB
