package main

import (
	prompt "github.com/c-bata/go-prompt"
)

var infoWindow *prompt.LinearInfoWindow
var p *prompt.Prompt

func printInfo() {
	infoWindow.AddLine("Test info")
	p.UpdateInfoWindow(infoWindow)
}

func executor(t string) {
	if t == "info" {
		go printInfo()
	}
	return
}

func completer(t prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "info"},
	}
}

func main() {
	infoWindow = prompt.NewLinearInfoWindow(100, true)
	p = prompt.New(
		executor,
		completer,
		prompt.OptionInfoWindowHeight(30),
	)
	p.UpdateInfoWindow(infoWindow)
	p.Run()
}
