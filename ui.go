package main

import (
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type WordCardsApp struct {
	conf       CardsConfig
	app        fyne.App
	window     fyne.Window
	mainMenu   *fyne.Container
	rando      CardsRandomizer
	statistics map[string]Stats
	selectedLP LangPair
	reverse    bool
}

// ******************************************************
// INITIALIZE
// ******************************************************

func InitUI(conf CardsConfig) WordCardsApp {
	fyneApp := app.New()
	w := fyneApp.NewWindow("fancyCards")

	application := WordCardsApp{conf: conf, app: fyneApp, window: w}
	application.CreateMainMenu(conf)
	err := application.InitializeStatistics()
	application.ToMainMenu()
	application.HandleError(err)

	return application
}

func (a *WordCardsApp) CreateMainMenu(conf CardsConfig) {
	hello := widget.NewLabel("I <3 Wordcards")
	a.mainMenu = container.NewVBox(
		hello,
	)

	for _, lp := range conf.langPairs {
		lpstr := conf.GetLangPairAsString(lp)
		btn := widget.NewButton(lpstr, func() {
			a.OpenLangpairMenu(lp, false)
		})
		a.mainMenu.Add(btn)

		// Automatically adding reverse wordcards for each language pair
		// This skipped if there's a collision (language pair already exists in the original file)

		reversePair := lp.Flip()
		if !conf.LangPairExists(reversePair) {
			reverseLpstr := conf.GetLangPairAsString(reversePair)
			reverseBtn := widget.NewButton(reverseLpstr, func() {
				a.OpenLangpairMenu(lp, true)
			})
			a.mainMenu.Add(reverseBtn)
		}

	}

	a.mainMenu.Add(widget.NewButton("Anleitung", func() {
		a.HandleError(errors.New("Funktion ist noch nicht programmiert :P"))
	}))
}

// ******************************************************
// BASICS
// ******************************************************

func (a *WordCardsApp) ToMainMenu() {
	a.window.SetContent(a.mainMenu)
}

func (a *WordCardsApp) ReturnButton() *widget.Button {
	return widget.NewButton("Zurück zum Startmenü", func() {
		a.ToMainMenu()
	})
}

func (a *WordCardsApp) HandleError(err error) {
	if err == nil {
		return
	}
	errMsg := widget.NewLabel(err.Error())

	errWindow := container.NewVBox(
		errMsg,
		a.ReturnButton(),
	)

	a.window.SetContent(errWindow)
}

// ******************************************************
// WORDCARDS
// ******************************************************

func (a *WordCardsApp) OpenLangpairMenu(lp LangPair, reverse bool) {
	a.selectedLP = lp
	a.reverse = reverse

	hello := widget.NewLabel(a.conf.GetLangPairAsString(a.GetSelectedLangPair()))

	exerciseButton := widget.NewButton("Starten", func() {
		cards, err := ReadCards(a.conf, lp, reverse)
		a.HandleError(err)
		if err == nil {
			a.rando = NewRando(cards)
			a.LoadRandomCard()
		}
	})

	statsButton := widget.NewButton("Statistik ansehen", func() {
		a.ShowStatSummary()
	})

	lpMenu := container.NewVBox(hello, exerciseButton, statsButton, a.ReturnButton())
	a.window.SetContent(lpMenu)

}

func (a *WordCardsApp) GetSelectedLangPair() LangPair {
	if a.reverse {
		return a.selectedLP.Flip()
	} else {
		return a.selectedLP
	}
}

func (a *WordCardsApp) LoadRandomCard() {
	wc := a.rando.FetchRandomCard()

	textbox := widget.NewEntry()

	lang := a.conf.GetLangName(a.GetSelectedLangPair().sourceLang)
	inputWord := widget.NewLabel(fmt.Sprintf("%s: %s", lang, wc.sourceWord))

	checkBtn := widget.NewButton("Prüfen", func() {
		a.CheckCard(textbox.Text, wc)
	})

	cardsView := container.NewVBox(
		inputWord,
		textbox,
		checkBtn,
		a.ReturnButton(),
	)

	a.window.SetContent(cardsView)
}

func (a *WordCardsApp) CheckCard(word string, wc WordCard) {

	feedbackLabel := widget.NewLabel("")
	success := CheckInput(word, wc)
	switch success {
	case Wrong:
		feedbackLabel.SetText("Falsch!")
	case Similar:
		feedbackLabel.SetText("Ähnlich:")
	case Correct:
		feedbackLabel.SetText("Richtig!")
	case Skipped:
		feedbackLabel.SetText("Übersprungen...")
	}
	a.IncrementCount(success)

	continueBtn := widget.NewButton("Weiter", func() {
		a.LoadRandomCard()
	})

	correctSolution := widget.NewLabel(fmt.Sprintf("%s => %s", wc.sourceWord, wc.targetWord))

	resultView := container.NewVBox(
		feedbackLabel,
		correctSolution,
		continueBtn,
		a.ReturnButton(),
	)

	a.window.SetContent(resultView)
}

// ******************************************************
// STATISTICS
// ******************************************************

func (a *WordCardsApp) ShowStatSummary() {
	lp := a.GetSelectedLangPair()
	statEvals := a.GetStatEvals(lp)

	hello := widget.NewLabel(a.conf.GetLangPairAsString(lp))
	statPage := container.NewVBox(
		hello,
	)
	if len(statEvals) == 0 {
		statPage.Add(widget.NewLabel("Noch keine Statistik für diese Karten erfasst..."))
	}
	for key, eval := range statEvals {
		stringifyEval := fmt.Sprintf("%s - %d Karteikarten, %d richtig, %d falsch",
			key, eval.Count, eval.Successes, eval.Mistakes)
		statPage.Add(widget.NewLabel(stringifyEval))
	}

	statPage.Add(a.ReturnButton())
	a.window.SetContent(statPage)
}
