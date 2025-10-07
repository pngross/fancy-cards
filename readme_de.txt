Compiling-Infos:
- Auf dem PC muss ein C-Compiler installiert sein (erforderlich für die fyne GUI library)
- Beim Aufbauen für Windows ist folgender Command zu verwenden (damit beim Öffnen kein Shell-Fenster aufgeht):
> go build -ldflags "-H=windowsgui" fancyCards

Vorbereitung (Ordner für die App richtig einrichten):

CSV FILES:
Alle CSV-Dateien, die das Programm einliest, müssen Semikolon (;) als Trennzeichen verwenden.

INI:
- benenne die Beispiel-Ini zu fancyCards.ini um (enthält Standard-Einstellungen und ein paar Sprachen)

DATEILISTE:
- die Liste der Dateien wird aus der Datei eingelesen, die unter "fileListConfigFile" in der ini eingelesen ist
- Kopiere example_files.csv (Namen anpassen und in Ordner config legen) und befülle sie mit den Infos zu den eigenen Dateien
- Die Datei muss denselben Aufbau wie die example_files.csv haben:
    SPALTE 1: AUSGANGSSPRACHE
    SPALTE 2: LERNSPRACHE
    SPALTE 3: DATEINAME (nicht relativ, sondern nur der Dateiname, siehe unten)
    SPALTE 4: GRUPPEN (können mehrere sein, durch Kommata getrennt)

DATEIEN: (KARTEIKARTEN sind hier)
- Diese müssen ebenfalls in CSV-Dateien abgelegt werden
- Das Programm erwartet, dass die Input-Dateien in einem bestimmten Ordner sind: "input_" + die Kürzel der Sprachen
    BEISPIEL: wenn man Englisch als Muttersprache spricht und Italienisch lernen will,
        gehören sie in diesen Ordner:
            input_en_it

        Es wird automatisch erkannt (anhand der Dateiliste) und alle Dateien für Englisch->Italienisch müssen hier abgelegt werden. 
        Es kann in der fancyCards.ini ein anderer Präfix als "input_" festgelegt werden, das ist aber nur zu Testzwecken empfohlen.

- Bei den Karteikarten-Dateien muss der folgende Aufbau eingehalten werden, damit keine Fehler entstehen:
    SPALTE 1: WORT IN AUSGANGSSPRACHE
    SPALTE 2: WORT IN LERNSPRACHE
    SPALTE 3: KOMMENTAR (in Ausgangssprache, Funktion wird noch nicht unterstützt)
    SPALTE 4: KOMMENTAR (in Lernsprache, Funktion wird noch nicht unterstützt)

Fehler werden im GUI angezeigt. Die App kann sowohl über den Dateien-Explorer als auch über die Kommandozeile geöffnet werden. 