package main

func main() {

	if !FileExists("fancyCards.ini") {
		CreateDefaultIni("fancyCards.ini")
	}
	conf, err := loadConfigsIni("fancyCards.ini")

	application := InitUI(conf)
	application.HandleErrorList(err)
	application.window.ShowAndRun()
	application.SaveStatistics()
}
