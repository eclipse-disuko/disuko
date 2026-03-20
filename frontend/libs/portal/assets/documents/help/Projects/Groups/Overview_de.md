# Projektgruppen

## Übersicht

- <i aria-hidden="true" class="v-icon notranslate pr-1 material-icons labelIconColor--text mdi mdi-card-multiple-outline" style="font-size: 16px;" ></i> Eine **Projektgruppe** fasst eine Liste von Projekten (untergeordnete Projekte) zusammen
- <i aria-hidden="true" class="v-icon notranslate pr-1 material-icons labelIconColor--text mdi mdi-account-supervisor-circle" style="font-size: 16px;" ></i> Ein **Teilprojekt** hat eine übergeordnete Beziehung zur Gruppe
- Der Benutzer muss Eigentümer der Gruppe **und** Eigentümer des untergeordneten Elements sein, um eine untergeordnete Beziehung zu erstellen
- Ein Projekt kann nur einer einzelnen Gruppe untergeordnet sein
- Gelöschte Projekte können nicht als untergeordnetes Projekt zu einer Gruppe hinzugefügt werden

## Token

- Token der Projektgruppe werden an Teilprojekte weitergegeben
- Es ist daher möglich, einen technischen Token für alle Teilprojekte zu verwenden

## Benutzerverwaltung

- Die Benutzerverwaltung einer Gruppe wird **nicht** an Teilprojekte weitergegeben
- Ein Benutzer in einer Gruppe hat nur Rechte auf die Gruppe selbst (Ausnahme: Tokens)

## Löschen von Gruppen

- Wenn eine Gruppe gelöscht wird, wird der übergeordnete Eintrag für alle untergeordneten Elemente entfernt
- Wenn eine Gruppe gelöscht wird, bleiben die Teilprojekte als eigenständiges Projekt erhalten

## Löschen von Teilprojekten

- Wenn ein untergeordnetes Projekt gelöscht wird, wird das Projekt in der Projektgruppe als gelöscht markiert
- Die Gruppe zeigt dann das gelöschte Projekt mit einer Benachrichtigung an, dass es gelöscht wurde
- Der Dialog, in dem er untergeordnete Projekte bearbeitet, zeigt das gelöschte untergeordnete Projekt, aber **nicht** alle gelöschten Projekte.

## Entfernung von Teilprojekten

- Beim Entfernen eines untergeordneten Objekts wird das übergeordnete Attribut gelöscht und das Gruppenprojekt erhält einen neuen untergeordneten Listensatz
- Token auf Gruppenebene können nicht mehr für den Zugriff auf das entfernte untergeordnete Projekt verwendet werden
