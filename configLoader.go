package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	ini "gopkg.in/ini.v1"
)

type InputFile struct {
	fileName         string
	groups           []string
	sourceWordCol    int
	targetWordCol    int
	sourceCommentCol int
	targetCommentCol int
	skipHeaderLine   bool
}

func defaultInputFile() InputFile {
	k := InputFile{
		sourceWordCol:    0,
		targetWordCol:    1,
		sourceCommentCol: 2,
		skipHeaderLine:   true,
		targetCommentCol: 3,
	}
	return k
}

type LangPair struct {
	sourceLang string
	targetLang string
}

func (l LangPair) ToString() string {
	return l.sourceLang + "_" + l.targetLang
}

type CardsConfig struct {
	languageNames       map[string]string
	langPairs           []LangPair
	savDir              string
	inputDirPrefix      string
	languagesConfigFile string
	files               map[string][]InputFile
}

func (c CardsConfig) Init() CardsConfig {
	c.languageNames = map[string]string{}
	c.inputDirPrefix = ""
	c.languagesConfigFile = ""
	c.langPairs = []LangPair{}
	c.files = map[string][]InputFile{}
	return c
}

func (c *CardsConfig) AddFile(fn string, lp LangPair) {
	km := defaultInputFile()
	km.fileName = fn
	c.files[lp.ToString()] = append(c.files[lp.ToString()], km)
}

func (c *CardsConfig) ValidateAndAddFile(file InputFile, lp LangPair) error {
	if c.GetLangName(lp.sourceLang) == "" {
		return fmt.Errorf("Ungültige Ausgangssprache %s", lp.sourceLang)
	} else if c.GetLangName(lp.targetLang) == "" {
		return fmt.Errorf("Ungültige Lernsprache %s", lp.targetLang)
	}

	c.files[lp.ToString()] = append(c.files[lp.ToString()], file)
	if !c.LangPairExists(lp) {
		c.langPairs = append(c.langPairs, lp)
	}
	return nil
}

func (c CardsConfig) GetInputFiles(lpst string) []InputFile {
	return c.files[lpst]
}

func (c CardsConfig) GetLangName(id string) string {
	return string(c.languageNames[id])
}

func (c CardsConfig) GetLangPairAsString(lp LangPair) string {
	sl := c.GetLangName(lp.sourceLang)
	tl := c.GetLangName(lp.targetLang)
	lpstr := fmt.Sprintf("%s -> %s", sl, tl)
	return lpstr
}

func (c CardsConfig) LangPairExists(lp LangPair) bool {
	for _, existingLP := range c.langPairs {
		if lp.ToString() == existingLP.ToString() {
			return true
		}
	}
	return false
}

func processLanguageFileLine(input []string, i int) (InputFile, LangPair, error) {
	f := defaultInputFile()
	f.groups = make([]string, 0)
	lp := LangPair{}
	var err error

	if len(input) < 3 {
		err = fmt.Errorf("Zeile %d ist zu kurz!\n", i+1)
	} else if input[0] == "" || input[1] == "" || input[2] == "" {
		err = fmt.Errorf("Zeile %d ist ungültig - in Spalten 1 bis 3 dürfen keine leeren Einträge sein!\n", i+1)
	} else {
		lp.sourceLang = input[0]
		lp.targetLang = input[1]
		f.fileName = input[2]

		if len(input) >= 4 {
			f.groups = strings.Split(input[3], ",")
		}
	}
	return f, lp, err
}

func (c *CardsConfig) ReadLanguagesFile() error {

	f, err := os.Open(c.languagesConfigFile)
	if err != nil {
		return err
	}
	// Lese Sprachen aus CSV-Datei ein
	reader := csv.NewReader(f)
	reader.Comma = ';'

	inputData, err := reader.ReadAll()
	if err != nil {
		return err
	}
	for i, ds := range inputData {
		if i == 0 {
			continue
		}
		file, lp, err := processLanguageFileLine(ds, i)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		err = c.ValidateAndAddFile(file, lp)
		if err != nil {
			fmt.Printf("Zeile %d: %s\n", i+1, err.Error())
		}

	}

	return nil
}

func loadConfigsIni(inipath string) (CardsConfig, error) {
	c := CardsConfig{}.Init()

	CardsIniReader, err := ini.Load(inipath)
	if err != nil {
		return c, err
	}

	langSection := CardsIniReader.Section("LANGUAGES")
	for _, k := range langSection.Keys() {
		langKey := k.Name()
		c.languageNames[langKey] = k.String()
	}

	configfilesSection := CardsIniReader.Section("CONFIGFILES")
	for _, k := range configfilesSection.Keys() {
		val := k.String()
		switch k.Name() {
		case "inputDirPrefix":
			c.inputDirPrefix = val
		case "languagesConfigFile":
			c.languagesConfigFile = val
		case "savDir":
			c.savDir = val
		}
	}

	if len(configfilesSection.Keys()) == 0 {
		return c, fmt.Errorf("Fehler beim Einlesen der %s: Bereich [CONFIGFILES] fehlt", inipath)
	}
	if c.inputDirPrefix == "" || c.languagesConfigFile == "" {
		return c, fmt.Errorf("Fehler beim Einlesen der %s:\n Im Bereich [CONFIGFILES] fehlt inputDirPrefix und/oder languagesConfigFile", inipath)
	}
	if len(c.languageNames) == 0 {
		return c, fmt.Errorf("Fehler beim Einlesen der %s: Keine Sprachen definiert", inipath)
	}

	c.ReadLanguagesFile()

	return c, nil
}
