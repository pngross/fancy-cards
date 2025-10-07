package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	random "math/rand/v2"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type WordCard struct {
	sourceWord    string
	sourceComment string
	targetWord    string
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
		if len(rando.cards) == 1 {
			break
		}
	}
	rando.prevpos = pos
	return rando.cards[pos]
}

func readCardsFromCsv(mapp InputFile, inputdir string, reverse bool) ([]WordCard, error) {
	karten := []WordCard{}
	path := filepath.Join(inputdir, mapp.fileName)
	f, err := os.Open(path)
	if err != nil {
		return karten, fmt.Errorf("Datei '%s' konnte nicht geöffnet werden", path)
	}

	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	csvReader.LazyQuotes = false
	inputData, err := csvReader.ReadAll()
	if err != nil {
		return karten, fmt.Errorf("Aus Datei '%s' konnten keine Karteikarten gelesen werden!\nBitte die Datei prüfen.", path)
	}

	for i, ds := range inputData {

		if i == 0 && mapp.skipHeaderLine {
			continue
		}
		if len(ds) <= mapp.sourceWordCol || len(ds) <= mapp.targetWordCol {
			continue
		}

		var karte WordCard
		if reverse {
			karte = WordCard{sourceWord: ds[mapp.targetWordCol],
				targetWord: ds[mapp.sourceWordCol]}
		} else {
			karte = WordCard{sourceWord: ds[mapp.sourceWordCol],
				targetWord: ds[mapp.targetWordCol]}
		}

		if len(ds) > mapp.targetCommentCol {
			karte.sourceComment = ds[mapp.targetCommentCol]
		}
		karten = append(karten, karte)
	}

	return karten, nil
}

func ReadCards(conf CardsConfig, lp LangPair, reverse bool, groups []string) ([]WordCard, error) {
	inputfiles := conf.GetInputFiles(lp.ToString())

	allCards := []WordCard{}
	for _, file := range inputfiles {

		// Skip files whose groups don't match at least one of the groups provided in args
		found := false
		for _, group := range groups {
			if slices.Contains(file.groups, group) {
				found = true
				break
			}
		}
		if !found && len(groups) > 0 {
			continue
		}

		karten, err := readCardsFromCsv(file, conf.inputDirPrefix+lp.ToString(), reverse)
		if err != nil {
			return allCards, err
		}
		allCards = append(allCards, karten...)
	}
	if len(allCards) == 0 {
		return allCards, errors.New("Es wurden keine Karteikarten gefunden!")
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
