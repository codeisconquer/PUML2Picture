package main

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed assets/icon.png
var iconPng []byte

var resourceIconPng fyne.Resource = fyne.NewStaticResource("icon.png", iconPng)
