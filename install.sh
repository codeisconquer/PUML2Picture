#!/bin/bash

# Homebrew prüfen
if ! command -v brew >/dev/null 2>&1; then
  echo "❌ Homebrew nicht gefunden."
  echo "👉 Bitte installiere Homebrew zuerst: https://brew.sh"
  exit 1
fi

# Java prüfen und ggf. installieren
if ! command -v java >/dev/null 2>&1; then
  echo "🔧 Java wird installiert..."
  brew install openjdk
else
  echo "✅ Java ist bereits installiert."
fi

# Graphviz prüfen und ggf. installieren
if ! command -v dot >/dev/null 2>&1; then
  echo "🔧 Graphviz wird installiert..."
  brew install graphviz
else
  echo "✅ Graphviz ist bereits installiert."
fi
