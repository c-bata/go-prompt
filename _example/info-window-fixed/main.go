package main

import (
	prompt "github.com/c-bata/go-prompt"
)

var infoWindow *prompt.FixedInfoWindow
var p *prompt.Prompt

func printInfo() {
	infoWindow.Clear()
	l1 := infoWindow.RequestLine(0)
	*l1 = string("Hallo Line 1")
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
	infoWindow = prompt.NewFixedInfoWindow(20)
	p = prompt.New(
		executor,
		completer,
		prompt.OptionInfoWindowHeight(30),
	)
	p.UpdateInfoWindow(infoWindow)
	p.Run()
}
