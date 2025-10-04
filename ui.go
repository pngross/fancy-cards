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
	selectedLP LangPair
}

func InitUI(conf CardsConfig) WordCardsApp {
	fyneApp := app.New()
	w := fyneApp.NewWindow("fancyCards")

	application := WordCardsApp{conf: conf, app: fyneApp, window: w}
	application.CreateMainMenu(conf)
	application.ToMainMenu()
	return application
}

func (a *WordCardsApp) HandleError(err error) {
	if err == nil {
		return
	}
	errMsg := widget.NewLabel(err.Error())
	recoverBtn := widget.NewButton("Zurück zum Startmenü", func() {
		a.ToMainMenu()
	})

	errWindow := container.NewVBox(
		errMsg,
		recoverBtn,
	)

	a.window.SetContent(errWindow)
}

func (a *WordCardsApp) LoadRandomCard() {
	wc := a.rando.FetchRandomCard()

	textbox := widget.NewEntry()
	returnBtn := widget.NewButton("Zurück zum Startmenü", func() {
		a.ToMainMenu()
	})

	inputWord := widget.NewLabel(fmt.Sprintf("%s: %s",
		a.conf.GetLangName(a.selectedLP.sourceLang), wc.sourceWord))

	checkBtn := widget.NewButton("Prüfen", func() {
		a.CheckCard(textbox.Text, wc)
	})

	cardsView := container.NewVBox(
		inputWord,
		textbox,
		checkBtn,
		returnBtn,
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

	continueBtn := widget.NewButton("Weiter", func() {
		a.LoadRandomCard()
	})
	returnBtn := widget.NewButton("Zurück zum Startmenü", func() {
		a.ToMainMenu()
	})

	correctSolution := widget.NewLabel(fmt.Sprintf("%s => %s", wc.sourceWord, wc.targetWord))

	resultView := container.NewVBox(
		feedbackLabel,
		correctSolution,
		continueBtn,
		returnBtn,
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
			a.SelectLangpair(conf, lp)
		})
		a.mainMenu.Add(btn)
	}

	a.mainMenu.Add(widget.NewButton("Anleitung", func() {
		a.HandleError(errors.New("Funktion ist noch nicht programmiert :P"))
	}))

	a.mainMenu.Add(widget.NewButton("Statistik", func() {
		a.HandleError(errors.New("Funktion ist noch nicht programmiert :P"))
	}))
}

func (a *WordCardsApp) SelectLangpair(conf CardsConfig, lp LangPair) {
	a.selectedLP = lp
	cards, err := ReadCards(conf, lp)
	a.HandleError(err)
	if err == nil {
		a.rando = NewRando(cards)
		a.LoadRandomCard()
	}
}

func (a *WordCardsApp) ToMainMenu() {
	a.window.SetContent(a.mainMenu)
}
