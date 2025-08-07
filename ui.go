package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// blankImage erzeugt ein leeres weißes 1x1 Bild
func blankImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	return img
}

// Bilddaten in die macOS-Zwischenablage legen (nur PNG)
func copyToClipboardPNG(data []byte) error {
	tmpfile, err := os.CreateTemp("", "*.png")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(data); err != nil {
		return err
	}
	tmpfile.Close()

	cmd := exec.Command("osascript", "-e", fmt.Sprintf(`set the clipboard to (read (POSIX file "%s") as «class PNGf»)`, tmpfile.Name()))
	return cmd.Run()
}

// Hauptfunktion zur Erstellung der UI
func buildUI(w fyne.Window) fyne.CanvasObject {
	// PUML-Code Eingabefeld
	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("PUML-Code hier eingeben oder einfügen...")

	// Formatwahl (png oder svg)
	formatSelect := widget.NewSelect([]string{"png", "svg"}, func(string) {})
	formatSelect.SetSelected("png")

	// Checkbox: Automatisch in Zwischenablage kopieren
	autoClipboard := widget.NewCheck("Auto to Clipboard", nil)

	// Dateiname-Eingabe (fixe Breite)
	filenameEntry := widget.NewEntry()
	filenameContainer := container.NewWithoutLayout(filenameEntry)

	filenameEntry.SetPlaceHolder("Dateiname")
	filenameEntry.Resize(fyne.NewSize(100, filenameEntry.MinSize().Height)) // Manuell setzen
	filenameContainer.Resize(fyne.NewSize(100, filenameEntry.MinSize().Height))

	// Bildvorschau
	imagePreview := canvas.NewImageFromImage(blankImage())
	imagePreview.FillMode = canvas.ImageFillContain
	imagePreview.SetMinSize(fyne.NewSize(400, 300))

	// Vorschau aktualisieren und speichern
	updatePreview := func() {
		format := formatSelect.Selected
		imgBytes, err := renderPUML(input.Text, format)
		if err != nil {
			fmt.Println("Fehler beim Rendern:", err)
			return
		}

		// Immer speichern
		saveCurrentImage(input.Text, format, filenameEntry.Text)

		// In Zwischenablage kopieren (nur PNG)
		if autoClipboard.Checked && format == "png" {
			err := copyToClipboardPNG(imgBytes)
			if err != nil {
				fmt.Println("Fehler beim Kopieren in Zwischenablage:", err)
			} else {
				fmt.Println("Bild erfolgreich in Zwischenablage kopiert.")
			}
		}

		// Vorschau aktualisieren
		if format == "png" {
			img, _, err := image.Decode(bytes.NewReader(imgBytes))
			if err != nil {
				fmt.Println("Fehler beim Decodieren:", err)
				return
			}
			imagePreview.Image = img
			imagePreview.Refresh()
			imagePreview.Show()
		} else {
			imagePreview.Image = nil
			imagePreview.Refresh()
		}
	}

	// Speichern-Button
	saveButton := widget.NewButton("Save", func() {
		updatePreview()
	})

	// Textänderung löst Vorschau aus
	input.OnChanged = func(string) {
		updatePreview()
	}

	// Layout zusammenbauen
	ui := container.NewVBox(
		container.NewBorder(nil, nil, nil, formatSelect, widget.NewLabel("PlantUML Code")),
		input,
		container.NewHBox(
			autoClipboard,
			// layout.NewSpacer(),
			saveButton,
			// layout.NewSpacer(),
			widget.NewLabel("Dateiname:"),
			filenameContainer,
		),
		widget.NewLabel("Vorschau (nur PNG):"),
		imagePreview,
	)

	return container.NewPadded(ui)
}

// Diagramm speichern im Downloads-Ordner
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
	} else {
		fmt.Println("Datei gespeichert:", fullPath)
	}
}
