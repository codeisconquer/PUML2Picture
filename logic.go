package main

import (
	"bytes"
	_ "embed"
	"fmt"
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
		return nil, fmt.Errorf("temp dir fehler: %w", err)
	}
	defer os.RemoveAll(tempDir) // Aufräumen

	// JAR-Datei speichern
	jarPath := filepath.Join(tempDir, "plantuml.jar")
	err = os.WriteFile(jarPath, plantumlJar, 0644)
	if err != nil {
		return nil, fmt.Errorf("jar schreiben fehler: %w", err)
	}

	// PUML-Code speichern
	pumlPath := filepath.Join(tempDir, "diagram.puml")
	err = os.WriteFile(pumlPath, []byte(code), 0644)
	if err != nil {
		return nil, fmt.Errorf("puml schreiben fehler: %w", err)
	}

	// Java-Aufruf vorbereiten
	outputFormat := "-t" + format
	cmd := exec.Command("java", "-jar", jarPath, outputFormat, pumlPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("plantuml fehler: %s", stderr.String())
	}

	// Ergebnisdatei laden
	outputFile := filepath.Join(tempDir, "diagram."+format)
	imgBytes, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, fmt.Errorf("ergebnis lesen fehler: %w", err)
	}

	return imgBytes, nil
}

// isCommandAvailable prüft, ob ein Kommando im System verfügbar ist
func isCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// checkDependencies prüft, ob Java und Graphviz (dot) installiert sind
// Gibt eine Liste fehlender Komponenten zurück
func checkDependencies() []string {
	var missing []string

	if !isCommandAvailable("java") {
		missing = append(missing, "Java Runtime (JRE)")
	}
	if !isCommandAvailable("dot") {
		missing = append(missing, "Graphviz (dot)")
	}

	return missing
}
