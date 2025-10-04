package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Stats struct {
	Count     int `json:"count"`
	Successes int `json:"successes"`
	Mistakes  int `json:"mistakes"`
}

// ******************************************************
// READ/WRITE BASICS
// ******************************************************

func (a *WordCardsApp) IncrementCount(success SuccessLevel) {
	s := a.LoadCurrentStats()
	s.Count++
	if success == Wrong {
		s.Mistakes++
	} else if success != Skipped {
		s.Successes++
	}
	a.UpdateCurrentStats(s)
}

func (a *WordCardsApp) LoadStats(lp LangPair, month, year int) Stats {
	s, ok := a.statistics[fmt.Sprintf("%s_%d_%d", lp.ToString(), month, year)]
	if ok {
		return s
	}
	return Stats{}
}

func (a *WordCardsApp) LoadCurrentStats() Stats {
	lp := a.selectedLP
	if a.reverse {
		lp = LangPair{sourceLang: lp.targetLang, targetLang: lp.sourceLang}
	}
	d := time.Now()
	month := int(d.Month())
	return a.LoadStats(lp, month, d.Year())
}

func (a *WordCardsApp) UpdateCurrentStats(s Stats) {
	lp := a.selectedLP
	if a.reverse {
		lp = LangPair{sourceLang: lp.targetLang, targetLang: lp.sourceLang}
	}
	d := time.Now()
	month := int(d.Month())
	a.statistics[fmt.Sprintf("%s_%d_%d", lp.ToString(), month, d.Year())] = s
}

// ******************************************************
// LISTING
// ******************************************************

func (a *WordCardsApp) GetStatEvals(lp LangPair) map[string]Stats {
	prefix := lp.ToString()
	mp := map[string]Stats{}
	for key, value := range a.statistics {
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		outkey := strings.ReplaceAll(strings.ReplaceAll(key, prefix+"_", ""), "_", "/")
		mp[outkey] = value
	}
	return mp
}

// ******************************************************
// HANDLING FILES
// ******************************************************

func (a *WordCardsApp) SaveStatistics() {
	if a.conf.savDir == "" {
		return
	}
	savFile := filepath.Join(a.conf.savDir, "statistik.json")
	backupFile := filepath.Join(a.conf.savDir, "_statistik.json")

	if FileExists(backupFile) {
		os.Remove(backupFile)
	}
	if FileExists(savFile) {
		os.Rename(savFile, backupFile)
	}

	data, err := json.MarshalIndent(a.statistics, "", "  ") // compact JSON
	if err != nil {
		os.Rename(backupFile, savFile)
		return
	}

	if err := os.WriteFile(savFile, data, 0644); err != nil {
		os.Rename(backupFile, savFile)
		return
	}
}

func (a *WordCardsApp) InitializeStatistics() error {
	a.statistics = map[string]Stats{}
	if a.conf.savDir == "" {
		return errors.New("savDir fehlt in fancyCards.ini")
	}

	savFile := filepath.Join(a.conf.savDir, "statistik.json")
	backupFile := filepath.Join(a.conf.savDir, "_statistik.json")
	backupFile2 := filepath.Join(a.conf.savDir, "__statistik.json")

	f, err := os.Open(savFile)

	if err != nil {
		if FileExists(backupFile) && !FileExists(backupFile2) {
			os.Rename(backupFile, backupFile2)
		}
		return nil
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	var st map[string]Stats
	err = json.Unmarshal(data, &st)
	if err != nil {
		return err
	}

	a.statistics = st
	return nil
}
