package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	// App starten
	a := app.NewWithID("com.example.puml")
	w := a.NewWindow("PUML")
	w.Resize(fyne.NewSize(800, 600))

	// Benutzeroberfläche aufbauen
	content := buildUI(w)
	w.SetContent(content)

	// App-Icon setzen
	w.SetIcon(resourceIconPng)

	w.ShowAndRun()
}
