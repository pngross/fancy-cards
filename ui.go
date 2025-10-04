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

func (a *WordCardsApp) LoadRandomCard() {
	wc := a.rando.FetchRandomCard()

	textbox := widget.NewEntry()

	var lang string
	if a.reverse {
		lang = a.conf.GetLangName(a.selectedLP.targetLang)
	} else {
		lang = a.conf.GetLangName(a.selectedLP.sourceLang)
	}

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

	var sw, tw string
	if a.reverse {
		tw = wc.sourceWord
		sw = wc.targetWord
	} else {
		sw = wc.sourceWord
		tw = wc.targetWord
	}

	correctSolution := widget.NewLabel(fmt.Sprintf("%s => %s", sw, tw))

	resultView := container.NewVBox(
		feedbackLabel,
		correctSolution,
		continueBtn,
		a.ReturnButton(),
	)

	a.window.SetContent(resultView)
}

func (a *WordCardsApp) CreateMainMenu(conf CardsConfig) {
	hello := widget.NewLabel("I <3 Wordcards")
	a.mainMenu = container.NewVBox(
		hello,
	)

	for _, lp := range conf.langPairs {
		lpstr := conf.GetLangPairAsString(lp)
		btn := widget.NewButton(lpstr, func() {
			a.SelectLangpair(lp, false)
		})
		a.mainMenu.Add(btn)

		// Automatically adding reverse wordcards for each language pair
		// This skipped if there's a collision (language pair already exists in the original file)

		reversePair := LangPair{sourceLang: lp.targetLang, targetLang: lp.sourceLang}
		if !conf.LangPairExists(reversePair) {
			reverseLpstr := conf.GetLangPairAsString(reversePair)
			reverseBtn := widget.NewButton(reverseLpstr, func() {
				a.SelectLangpair(lp, true)
			})
			a.mainMenu.Add(reverseBtn)
		}

	}

	a.mainMenu.Add(widget.NewButton("Anleitung", func() {
		a.HandleError(errors.New("Funktion ist noch nicht programmiert :P"))
	}))

	a.mainMenu.Add(widget.NewButton("Statistik", func() {
		a.DisplayStatsMenu()
	}))
}

func (a *WordCardsApp) SelectLangpair(lp LangPair, reverse bool) {
	a.selectedLP = lp
	a.reverse = reverse
	cards, err := ReadCards(a.conf, lp, reverse)
	a.HandleError(err)
	if err == nil {
		a.rando = NewRando(cards)
		a.LoadRandomCard()
	}
}

// ******************************************************
// STATISTICS
// ******************************************************

func (a *WordCardsApp) DisplayStatsMenu() {
	hello := widget.NewLabel("Statistik-Übersicht")
	statMenu := container.NewVBox(
		hello,
	)

	for _, lp := range a.conf.langPairs {
		lpstr := a.conf.GetLangPairAsString(lp)
		btn := widget.NewButton(lpstr, func() {
			a.ShowStatSummary(lp)
		})
		statMenu.Add(btn)

		// Automatically adding reverse wordcards for each language pair
		// This skipped if there's a collision (language pair already exists in the original file)

		reversePair := LangPair{sourceLang: lp.targetLang, targetLang: lp.sourceLang}
		if !a.conf.LangPairExists(reversePair) {
			reverseLpstr := a.conf.GetLangPairAsString(reversePair)
			reverseBtn := widget.NewButton(reverseLpstr, func() {
				a.ShowStatSummary(reversePair)
			})
			statMenu.Add(reverseBtn)
		}
	}

	statMenu.Add(a.ReturnButton())
	a.window.SetContent(statMenu)
}

func (a *WordCardsApp) ShowStatSummary(lp LangPair) {
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
