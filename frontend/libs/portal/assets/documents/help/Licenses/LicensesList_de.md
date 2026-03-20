# <span class="material-icons material-symbols-outlined mdi mdi-shield-check-outline"></span> FOSS Lizenzdatenbank

## Lizenzen:

- **Tabelle:** zeigt eine Liste aller Lizenzen in der FOSS Lizenzdatenbank an, vorgefiltert, um nur überprüfte Lizenzen anzuzeigen
- **License Chart Status:** Nur Lizenzen mit dem Status "License Chart" wurden rechtlich geprüft und sind grundsätzlich zulässig, aber geltende Richtlinien wie Allow &amp; Deny-Listen müssen beachtet werden
- **License Chart fehlt:** Um den Status einer fehlenden Lizenz oder eines Lizenzdiagramms zu beheben, gibt es verschiedene Optionen:
- Möglicherweise ist die Lizenz unter einer anderen Kennung bekannt, so dass eine Untersuchung in der FOSS Lizenzdatenbank und Hinzufügen eines Alias des bekannten Lizenzeintrags dieses Problem beheben kann
- Setzen eines License Chart für die jeweilige Lizenz; Prüfen Sie zunächst, ob License Chart bereits unter einer anderen Kennung verfügbar ist
- <i aria-hidden="true" class="mdi-filter-variant mdi v-icon notranslate v-theme--dark v-icon--size-default text-blue"></i> Auf diese Liste wird ein Standardfilter angewendet, der nur Lizenzeinträge anzeigt, die den Status "License Chart" oder "Verboten" haben; Um alle Lizenzeinträge anzuzeigen, müssen Sie den Filter auf dieses Symbol in der Spaltenüberschrift "Bewertungsstatus" anpassen
- **Hinweis:** Verwenden Sie Strg+Klick, um eine Lizenz in einem neuen Browser-Tab zu öffnen

## Klassifizierungen:

- Insgesamt ist das Ziel, rechtliche Risiken transparent zu machen und Produktverantwortlichen zu helfen, ihr Produkt in einer Weise zu liefern, die rechtlich konform mit den FOSS-Lizenzanforderungen ist.
- Die Klassifizierungen unterscheiden zwischen drei Vorsichtsstufen. Sie reichen von "Alarm" über "Warnung" bis "Information":
  - <i aria-hidden="true" class="mdi-bell mdi v-icon notranslate v-theme--dark v-icon--size-default text-red"></i> **Alarm** zeigt an, dass die Klassifizierung potenziell ein hohes Risiko mit sich bringen kann.
  - <i aria-hidden="true" class="mdi-alert mdi v-icon notranslate v-theme--dark v-icon--size-default text-warning"></i> **Warnung** zeigt an, dass die Klassifizierung potenziell ein Risiko mit sich bringen kann. Das Risiko ist möglicherweise nicht so hoch oder so präsent wie bei „Alarm“, da es notwendig sein kann, zu bestimmen, ob das Ereignis, das die Klassifizierungszuweisung auslöst, tatsächlich eintritt.
  - <i aria-hidden="true" class="mdi-information-outline mdi v-icon notranslate v-theme--dark v-icon--size-default text-textColor"></i> **Information** ist eine Erinnerung daran, dass eine Verpflichtung möglicherweise erfüllt werden muss.
  - Anmerkung: In manchen Ansichten wird nur die höchste Risikostufe einer Lizenz angezeigt, während eine Lizenz typischerweise auch Klassifizierungen von niedrigeren Risikostufen enthält.

## Lizenztyp:

- **Nicht deklariert**: noch nicht bewertet und definiert
- **Open Source**: Open-Source-Software ist in der Regel kostenlos, wobei der Quellcode unter Bedingungen lizenziert ist, die jeder einsehen, ändern und verbessern kann
- **nicht-FOSS**: Der Oberbegriff von Freeware und proprietärer Software, die häufig ohne die Freiheiten aus Open Source definiert sind, wie Änderbarkeit und Verbesserbarkeit
- **Freeware**: Freeware ist Software, die in der Regel proprietär ist und kostenlos vertrieben wird, aber nicht die Freiheiten bietet, die Open Source bietet, wie z. B. Modifizieren und Verbessern
- **Proprietär**: Kommerzielle Software, die in der Regel gegen eine Gebühr zur Verfügung gestellt wird, ohne die Freiheiten, die Open Source bietet, wie z. B. Modifizieren und Verbessern
- **Public Domain**: Public Domain ist ein Sonderfall, bei dem Software kostenlos zur Verfügung steht, ohne lizenziert zu sein

## Bewertungsstatus:

- **Nicht gesetzt**
- **In rechtlicher Bewertung**
- **In Risikoprüfung**
- **In Policy-Zuordnung**
- **Bewertet**
- **Verboten**
- **Veraltet**

## Lizenzfamilie:

- **Nicht deklariert**
- **Permissiv**
- **Schwaches Copyleft**
- **Starkes Copyleft**
- **Netzwerk-Copyleft**