package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed assets/plantuml.jar
var plantumlJar []byte

// renderPUML nimmt PUML-Text und erzeugt daraus ein PNG oder SVG.
// Gibt das Bild als []byte zurück.
func renderPUML(code string, format string) ([]byte, error) {
	// Temporäres Verzeichnis anlegen
	tempDir, err := os.MkdirTemp("", "puml")
	if err != nil {
		return nil, fmt.Errorf("Temp dir Fehler: %w", err)
	}
	defer os.RemoveAll(tempDir) // Aufräumen

	// JAR-Datei speichern
	jarPath := filepath.Join(tempDir, "plantuml.jar")
	err = ioutil.WriteFile(jarPath, plantumlJar, 0644)
	if err != nil {
		return nil, fmt.Errorf("JAR schreiben Fehler: %w", err)
	}

	// PUML-Code speichern
	pumlPath := filepath.Join(tempDir, "diagram.puml")
	err = ioutil.WriteFile(pumlPath, []byte(code), 0644)
	if err != nil {
		return nil, fmt.Errorf("PUML schreiben Fehler: %w", err)
	}

	// Java-Aufruf vorbereiten
	outputFormat := "-t" + format
	cmd := exec.Command("java", "-jar", jarPath, outputFormat, pumlPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("PlantUML Fehler: %s", stderr.String())
	}

	// Ergebnisdatei laden
	outputFile := filepath.Join(tempDir, "diagram."+format)
	imgBytes, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, fmt.Errorf("Ergebnis lesen Fehler: %w", err)
	}

	return imgBytes, nil
}
