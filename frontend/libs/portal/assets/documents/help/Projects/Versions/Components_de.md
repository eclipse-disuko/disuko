# Komponenten Management

## <i class="v-icon notranslate material-icons mdi mdi-layers" style="color: var(--v-textColor-base); Caret-Farbe: var(--v-textColor-base);" ></i> Überblick:

- **Name & Version** aller Komponenten, die in der ausgewählten SBOM-Lieferung enthalten sind, sind in der Tabelle aufgeführt
- **Typ:** Eine Komponente ist in der Regel ein Paket, kann aber auch eine Datei sein
- **Lizenz gültig:** "Lizenz deklariert" Information der Komponente, es sei denn, es wird eine "Lizenz gefolgert"
  Information gemeldet
- **Policy-Regel:** Überprüfen Sie den Ergebnisstatus der Komponentenlizenz gemäß den für das Projekt geltenden
  Policy-Regeln

## <i class="v-icon notranslate material-icons mdi mdi-minus-circle text-error" style=""></i> Policy-Regel-Status "Verweigert":

- **"Verweigert"** ist generell nicht zu genehmigen, sollte daher vom Entwicklungsteam entfernt oder ersetzt werden.
- **"Verweigert (verboten)"** Lizenz ist ausnahmslos verboten für die Verwendung in selbst entwickelter Software (in jedem Distributionsmodell), Entfernen oder Ersetzen.

## <i class="v-icon notranslate material-icons policyStatusUnassertedColor--text mdi mdi-lightning-bolt-circle" style=""></i> Policy-Regel-Status "Nicht zugeordnet":

- **"Nicht zugeordnet"** bedeutet, dass Lizenzinformationen für die Komponente fehlen oder in der internen Wissensdatenbank unbekannt sind; Solche Fälle müssen gelöst werden.
- **Unbekannt:** Für eine Komponente sind überhaupt keine Lizenzinformationen vorhanden (wie "NONE", "" oder "NOASSERTION"), die Verwendung ist verboten; In diesem Fall muss das Entwicklungsteam die Lizenzinformationen untersuchen.
- **License Chart Status:** Es dürfen nur Lizenzen mit dem Status "License Chart" verwendet werden; Um den Status des fehlenden Lizenzdiagramms zu beheben, überprüfen Sie Ihre Optionen:
  - a. Die Lizenz ist möglicherweise unter einer anderen Kennung bekannt, daher können Sie dieses Problem in der Lizenzdatenbank überprüfen und durch Hinzufügen eines Alias zum bekannten Lizenzeintrag das Problem möglicherweise beheben.
  - b. die Entfernung oder den Ersatz durch das Entwicklungsteam verlangen.
  - c. Setzen eines License Chart für die jeweilige Lizenz (nach Abschluss von Option a).

## <i class="v-icon notranslate material-icons mdi mdi-alert text-warning" style=""></i> Policy-Regel-Status "Verwarnt":

- **"Verwarnt"** bedeutet, dass die Lizenz einer Komponente auf einer Verweigerungsliste steht, die eine Untersuchung und eine bewusste Entscheidung erfordert.
- Die Gründe, warum eine Lizenz mit einer Warnung versehen ist, können unterschiedlich sein. Zum Beispiel kann eine schwache Copyleft Lizenz wegen des Risikos gewarnt werden, dass das Copyleft unbeabsichtigt auf das Projekt angewendet werden könnte, z. B. wenn die Komponente statisch verknüpft ist.
- **Problembehebung:** Wenn eine Komponente mit einer Lizenz auf einer Warnliste enthalten ist, gibt es einige Schritte:
  - Bewerten Sie die potenziellen Risiken, die mit der Aufnahme der Komponente in Ihr Projekt verbunden sind, die Lizenzinformationen in der Lizenzdatenbank bieten eine Orientierungshilfe dafür
  - Prüfen Sie, ob die Komponente in der Mehrfachlizenzierung mit einer besseren Lizenzoption verfügbar ist
  - Prüfen Sie, ob der Lieferant die Komponente entfernen oder durch eine andere Komponente oder Funktionalität mit einer geeigneteren Lizenz ersetzen kann
  - Dokumentieren Sie die Entscheidung und die Begründung (diese Dokumentation ist intern aufzubewahren)

## <i class="v-icon notranslate material-icons mdi-help-circle-outline mdi v-theme--dark" style="color: var(--v-textColor-base); caret-color: var(--v-textColor-base);"></i> Policy-Regel-Status "Unentschieden":

- Aufgrund einer Mehrfachlizenzierung kann für eine Komponente mehr als eine Lizenz gelten.
- **"Unentschieden"** ist keine Aussage zum Policy-Regel-Status, bitte gehen Sie dazu in die Detailansicht der Komponente.
  Überprüfen Sie dort den Status der Policy-Regel, es gilt die "schlechteste" bzw. "problematischste" Lizenz.
  - **AND:** im Falle von "AND" (UND) gibt es keine Auswahl zu treffen, z.B. bei einer Dual-License-Komponente, die "MIT AND Apache-2.0" deklariert, müssen beide Lizenzen gleichzeitig eingehalten werden.
  - **OR:** im Falle einer "OR"-Aussage (ODER) können wir eine der angegebenen Lizenzen wählen und müssen dann nur eine Lizenz einhalten. Sofern ein "OR" nicht aufgelöst wird, wird die Komponente unter jeder der angegebenen Lizenzen aufgeführt. Um ein "OR" (je nach Anwendungsfall) aufzulösen, sollte die folgende Vorgehensweise angewendet werden:
    - Verwenden Sie die permissivste und gebräuchlichste Lizenz, es sei denn, Sie sind gezwungen, aus Konflikt- oder Kompatibilitätsgründen eine Copyleft-Lizenz zu verwenden (z. B. bei "MIT OR GPL-2.0" wählen Sie MIT)
    - Verwenden Sie nur genehmigte Lizenzen im Status "License Chart"
    - Die Verwendung alter Komponenten sollte nach Möglichkeit vermieden werden
    - Sofern ein "OR" nicht aufgelöst wird, muss die Komponente unter jeder der Lizenzen aufgeführt werden (zur Erinnerung: Nur genehmigte Lizenzen verwenden)
    - Dokumentieren Sie die Entscheidung und Begründung (interne Dokumentation)

## <i class="v-icon notranslate material-icons mdi mdi-check-circle text-green" style=""></i> Policy-Regel-Status "Zulässig":

- **Zulässig"** bedeutet, dass die überprüfte Lizenz für diesen Anwendungsfall zulässig ist.
- Dieser Status beschränkt sich auf die Lizenzbewertung, er spiegelt nicht die Technologie, Wartung, Codequalität oder Cybersicherheit der Komponente wider, diese müssen separat bewertet werden. Möglicherweise müssen Sie bei der Qualitätsprüfung der Lieferung weitere Erkenntnisse und Anmerkungen berücksichtigen.
