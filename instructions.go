package main

import (
	"fmt"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (a *WordCardsApp) OpenInstructions() {

	header := NewViewHeader("Anleitung")

	instructionsText := "DATEILISTE:"
	instructionsText += fmt.Sprintf("\nIn der Datei %s alle Dateien auflisten, aus denen Karteikarten gelesen werden sollen.", a.conf.fileListConfigFile)
	instructionsText += "\nZEILE 1: ÜBERSCHRIFT (wird übersprungen)"
	instructionsText += "\nSPALTE 1: Kürzel der Ausgangssprache (z. B. de)\nSPALTE 2: Kürzel der Lernsprache (z. B. fr)"
	instructionsText += "\nSPALTE 3: Dateiname\nPro Sprache können beliebig viele Karteikarten-Dateien hier eingetragen werden."

	instructionsText += "\n\nINPUT-DATEIEN: (die Wörter für die Karteikarten)"
	instructionsText += "\nFür jedes Sprachpaar sollen die Input-Dateien in einem eigenen Ordner liegen."
	instructionsText += fmt.Sprintf("\nBeispiel: Deutsch-Französisch => Ordner %sde_fr", a.conf.inputDirPrefix)
	instructionsText += "\n\nDiese Dateien müssen CSV-Tabellen sein, mit dem folgenden Aufbau:"
	instructionsText += "\nZEILE 1: ÜBERSCHRIFT (wird übersprungen)"
	instructionsText += "\nSPALTE 1: Wort in der Ausgangssprache\nSPALTE 2: Wort in der Lernsprache"
	instructions := widget.NewLabel(instructionsText)

	lpMenu := container.NewVBox(header, instructions, a.ReturnButton())

	a.window.SetContent(lpMenu)
}
