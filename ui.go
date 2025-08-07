package main

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func buildUI(w fyne.Window) fyne.CanvasObject {
	// Eingabefeld für PUML-Code
	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("PUML-Code hier eingeben oder einfügen...")

	// Format-Auswahl
	formatSelect := widget.NewSelect([]string{"png", "svg"}, func(string) {})
	formatSelect.SetSelected("png")

	// Auto-Save Checkbox
	autoSave := widget.NewCheck("Auto-Save", nil)

	// Dateiname-Eingabe mit fixer Breite
	filenameEntry := widget.NewEntry()
	filenameEntry.SetPlaceHolder("Dateiname")
	filenameEntry.Resize(fyne.NewSize(100, filenameEntry.MinSize().Height))

	// Bildvorschau
	imagePreview := canvas.NewImageFromImage(nil)
	imagePreview.FillMode = canvas.ImageFillContain

	// Vorschau-Update-Logik
	updatePreview := func() {
		format := formatSelect.Selected
		imgBytes, err := renderPUML(input.Text, format)
		if err != nil {
			fmt.Println("Fehler beim Rendern:", err)
			return
		}

		if autoSave.Checked {
			saveCurrentImage(input.Text, format, filenameEntry.Text)
		}

		if format == "png" {
			img, _, err := image.Decode(bytes.NewReader(imgBytes))
			if err != nil {
				fmt.Println("Fehler beim Decodieren:", err)
				return
			}
			imagePreview.Image = img
			imagePreview.Refresh()
		} else {
			imagePreview.Image = nil
			imagePreview.Refresh()
		}
	}

	// Save-Button
	saveButton := widget.NewButton("Save", func() {
		saveCurrentImage(input.Text, formatSelect.Selected, filenameEntry.Text)
		updatePreview()
	})

	// Textänderung = Vorschau aktualisieren
	input.OnChanged = func(string) {
		updatePreview()
	}

	// GUI zusammenbauen
	ui := container.NewVBox(
		container.NewBorder(nil, nil, nil, formatSelect, widget.NewLabel("PlantUML Code")),
		input,
		container.NewHBox(
			autoSave,
			container.NewMax(
				container.NewWithoutLayout(filenameEntry),
			),
			layout.NewSpacer(),
			saveButton,
		),
		widget.NewLabel("Vorschau (nur PNG):"),
		imagePreview,
	)

	return container.NewPadded(ui)
}

// Speichert das aktuelle Diagramm im Downloads-Ordner
func saveCurrentImage(code, format, filename string) {
	imgBytes, err := renderPUML(code, format)
	if err != nil {
		fmt.Println("Fehler beim Speichern:", err)
		return
	}

	user, _ := user.Current()
	downloads := filepath.Join(user.HomeDir, "Downloads")
	os.MkdirAll(downloads, os.ModePerm)

	if filename == "" {
		filename = fmt.Sprintf("puml_%d", time.Now().Unix())
	}
	fullPath := filepath.Join(downloads, filename+"."+format)

	err = ioutil.WriteFile(fullPath, imgBytes, 0644)
	if err != nil {
		fmt.Println("Fehler beim Schreiben der Datei:", err)
	}
}
