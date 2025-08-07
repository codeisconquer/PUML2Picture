if ! command -v brew >/dev/null 2>&1; then
  echo "Bitte installiere Homebrew zuerst: https://brew.sh"
  exit 1
fi

brew install openjdk graphviz
