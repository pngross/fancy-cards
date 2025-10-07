package main

import (
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
	viewHeader := NewViewHeader("I <3 Wordcards")
	a.mainMenu = container.NewVBox(
		viewHeader,
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
		a.OpenInstructions()
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

func NewViewHeader(str string) *widget.Label {
	vh := widget.NewLabel(str)
	vh.Alignment = fyne.TextAlignCenter
	vh.TextStyle = fyne.TextStyle{Bold: true}
	return vh
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

	viewHeader := NewViewHeader(a.conf.GetLangPairAsString(a.GetSelectedLangPair()))

	exerciseButton := widget.NewButton("Starten", func() {
		cards, err := ReadCards(a.conf, lp, reverse, []string{})
		a.HandleError(err)
		if err == nil {
			a.rando = NewRando(cards)
			a.LoadRandomCard()
		}
	})

	toGroupsButton := widget.NewButton("Wörter-Gruppen auswählen", func() {
		a.GroupSelection()
	})

	statsButton := widget.NewButton("Statistik ansehen", func() {
		a.ShowStatSummary()
	})

	lpMenu := container.NewVBox(viewHeader, exerciseButton, toGroupsButton, statsButton, a.ReturnButton())
	a.window.SetContent(lpMenu)

}

func (a *WordCardsApp) GetSelectedLangPair() LangPair {
	if a.reverse {
		return a.selectedLP.Flip()
	} else {
		return a.selectedLP
	}
}

func (a *WordCardsApp) GroupSelection() {

	viewHeader := NewViewHeader(a.conf.GetLangPairAsString(a.GetSelectedLangPair()) + " - Wörter-Gruppen auswählen")

	checkboxGroup := widget.NewCheckGroup(a.conf.GetGroups(a.selectedLP), func(strs []string) {})

	exerciseButton := widget.NewButton("Starten", func() {
		cards, err := ReadCards(a.conf, a.selectedLP, a.reverse, checkboxGroup.Selected)
		a.HandleError(err)
		if err == nil {
			a.rando = NewRando(cards)
			a.LoadRandomCard()
		}
	})

	backButton := widget.NewButton("Zurück", func() {
		a.OpenLangpairMenu(a.selectedLP, a.reverse)
	})

	groupsMenu := container.NewVBox(viewHeader, checkboxGroup, exerciseButton, backButton, a.ReturnButton())
	a.window.SetContent(groupsMenu)
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

	viewHeader := NewViewHeader(a.conf.GetLangPairAsString(lp))
	statPage := container.NewVBox(
		viewHeader,
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
