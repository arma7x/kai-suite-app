package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ProgressInfiniteDialog struct {
	dialog Dialog

	bar *widget.ProgressBarInfinite
}

func NewProgressInfinite(title, message string, parent fyne.Window) *ProgressInfiniteDialog {
	bar := widget.NewProgressBarInfinite()
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(200, 0))
	d := NewCustom(title, message, container.NewMax(rect, bar), parent)
	bar.Show()
	d.Show()
	return &ProgressInfiniteDialog{d, bar}
}

func (d *ProgressInfiniteDialog) Show() {
	d.bar.Show()
	d.dialog.Show()
}

func (d *ProgressInfiniteDialog) Hide() {
	d.bar.Hide()
	d.dialog.Hide()
}
