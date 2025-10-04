package main

func main() {

	conf, err := loadConfigsIni("fancyCards.ini")

	application := InitUI(conf)
	application.HandleError(err)
	application.window.ShowAndRun()
	application.SaveStatistics()
}
