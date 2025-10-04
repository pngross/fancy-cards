package main

import (
	"encoding/csv"
	"fmt"
	random "math/rand/v2"
	"os"
	"path/filepath"
	"strings"
)

type WordCard struct {
	sourceWord    string
	sourceComment string
	targetWord    string
	group         string
}

type CardsRandomizer struct {
	cards    []WordCard
	prevpos  int
	randoSum int
}

type SuccessLevel int

const (
	Wrong SuccessLevel = iota
	Similar
	Correct
	Skipped
)

func NewRando(list []WordCard) CardsRandomizer {
	rando := CardsRandomizer{prevpos: -1, cards: list, randoSum: len(list)}
	return rando
}

func (rando *CardsRandomizer) FetchRandomCard() WordCard {
	pos := rando.prevpos
	for pos == rando.prevpos {
		pos = random.IntN(rando.randoSum)
	}
	rando.prevpos = pos
	return rando.cards[pos]
}

func readCardsFromCsv(mapp InputFile, inputdir, group string) ([]WordCard, error) {
	karten := []WordCard{}
	f, err := os.Open(filepath.Join(inputdir, mapp.fileName))
	if err != nil {
		return karten, err
	}
	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	csvReader.LazyQuotes = false
	inputData, err := csvReader.ReadAll()
	if err != nil {
		return karten, err
	}
	for i, ds := range inputData {

		if i == 0 && mapp.skipHeaderLine {
			continue
		}
		if len(ds) <= mapp.sourceWordCol || len(ds) <= mapp.targetWordCol {
			continue
		}

		karte := WordCard{sourceWord: ds[mapp.sourceWordCol],
			targetWord: ds[mapp.targetWordCol],
			group:      group}

		if len(ds) > mapp.sourceCommentCol {
			karte.sourceComment = ds[mapp.sourceCommentCol]
		}
		karten = append(karten, karte)
	}

	return karten, nil
}

func ReadCards(conf CardsConfig, lp LangPair) ([]WordCard, error) {
	mappings := conf.GetInputFiles(lp.ToString())

	allCards := []WordCard{}
	for _, mapp := range mappings {
		karten, err := readCardsFromCsv(mapp, conf.inputDirPrefix+lp.ToString(), "test")
		if err != nil {
			fmt.Println(err.Error())
			return allCards, err
		}
		allCards = append(allCards, karten...)
	}
	return allCards, nil
}

func CheckInput(word string, wc WordCard) SuccessLevel {
	in := strings.ToLower(word)
	targ := strings.ToLower(wc.targetWord)
	if word == "" {
		return Skipped
	} else if targ == in {
		return Correct
	} else if strings.Contains(targ, in) || strings.Contains(in, targ) {
		return Similar
	}
	return Wrong
}
