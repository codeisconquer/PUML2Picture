# 🖥 PUML Desktop App

Die **PUML Desktop App** ist ein schlanker, nativer PlantUML-Editor für macOS, mit dem du UML-Diagramme direkt auf deinem Rechner erstellen, anzeigen und speichern kannst – ohne externe Dienste, vollständig offline und mit lokal eingebettetem `plantuml.jar`.

![Screenshot](assets/icon.png)

---

## ✨ Features

- 📝 PUML-Code direkt in der App eingeben oder einfügen
- 📸 Live-Vorschau des Diagramms (PNG)
- 📂 Export als **PNG** oder **SVG**
- 💾 Speichern in den **Downloads-Ordner**
- 🧠 Automatisches Speichern bei Eingabe (Auto-Save)
- 🔤 Benutzerdefinierbarer Dateiname
- 🍏 Native macOS `.app` mit eigenem Icon
- ☁️ Offline-fähig – kein Internet notwendig

---

## 🛠 Voraussetzungen

- [Go](https://golang.org/dl/) **≥ 1.21**
- [Java Runtime (JRE)](https://adoptium.net/) (z. B. via `brew install openjdk`)
- `plantuml.jar` in `assets/` (einmalig herunterladen von [plantuml.com/download](https://plantuml.com/download))

---

## 🚀 Starten der App

```bash
go run .
