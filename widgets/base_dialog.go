// Package dialog defines standard dialog windows for application GUIs.
package widget // import "fyne.io/fyne/v2/dialog"

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	padWidth  = 32
	padHeight = 16
)

// Dialog is the common API for any dialog window with a single dismiss button
type Dialog interface {
	Show()
	Hide()
	SetDismissText(label string)
	SetOnClosed(closed func())
	Refresh()
	Resize(size fyne.Size)

	// Since: 2.1
	MinSize() fyne.Size
}

// Declare conformity to Dialog interface
var _ Dialog = (*dialog)(nil)

type dialog struct {
	callback    func(bool)
	title       string
	icon        fyne.Resource
	desiredSize fyne.Size

	win            *widget.PopUp
	bg             *themedBackground
	content, label fyne.CanvasObject
	dismiss        *widget.Button
	parent         fyne.Window
	layout         *dialogLayout
}

// NewCustom creates and returns a dialog over the specified application using custom
// content. The button will have the dismiss text set.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func NewCustom(title, dismiss string, content fyne.CanvasObject, parent fyne.Window) Dialog {
	d := &dialog{content: container.NewHBox(layout.NewSpacer(), widget.NewLabel(dismiss), layout.NewSpacer()), title: title, icon: nil, parent: parent}
	d.layout = &dialogLayout{d: d}

	d.dismiss = &widget.Button{Text: dismiss,
		OnTapped: func(){},
	}
	d.setButtons(content)

	return d
}

// NewCustomConfirm creates and returns a dialog over the specified application using
// custom content. The cancel button will have the dismiss text set and the "OK" will
// use the confirm text. The response callback is called on user action.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func NewCustomConfirm(title, confirm, dismiss string, content fyne.CanvasObject,
	callback func(bool), parent fyne.Window) Dialog {
	d := &dialog{content: content, title: title, icon: nil, parent: parent}
	d.layout = &dialogLayout{d: d}
	d.callback = callback

	d.dismiss = &widget.Button{Text: dismiss, Icon: theme.CancelIcon(),
		OnTapped: d.Hide,
	}
	ok := &widget.Button{Text: confirm, Icon: theme.ConfirmIcon(), Importance: widget.HighImportance,
		OnTapped: func() {
			d.hideWithResponse(true)
		},
	}
	d.setButtons(container.NewHBox(layout.NewSpacer(), d.dismiss, ok, layout.NewSpacer()))

	return d
}

// ShowCustom shows a dialog over the specified application using custom
// content. The button will have the dismiss text set.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func ShowCustom(title, dismiss string, content fyne.CanvasObject, parent fyne.Window) {
	NewCustom(title, dismiss, content, parent).Show()
}

// ShowCustomConfirm shows a dialog over the specified application using custom
// content. The cancel button will have the dismiss text set and the "OK" will use
// the confirm text. The response callback is called on user action.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func ShowCustomConfirm(title, confirm, dismiss string, content fyne.CanvasObject,
	callback func(bool), parent fyne.Window) {
	NewCustomConfirm(title, confirm, dismiss, content, callback, parent).Show()
}

func (d *dialog) Hide() {
	d.hideWithResponse(false)
}

// MinSize returns the size that this dialog should not shrink below
//
// Since: 2.1
func (d *dialog) MinSize() fyne.Size {
	return d.win.MinSize()
}

func (d *dialog) Show() {
	if !d.desiredSize.IsZero() {
		d.win.Resize(d.desiredSize)
	}
	d.win.Show()
}

func (d *dialog) Refresh() {
	d.win.Refresh()
}

// Resize dialog, call this function after dialog show
func (d *dialog) Resize(size fyne.Size) {
	d.desiredSize = size
	d.win.Resize(size)
}

// SetDismissText allows custom text to be set in the confirmation button
func (d *dialog) SetDismissText(label string) {
	d.dismiss.SetText(label)
	d.win.Refresh()
}

// SetOnClosed allows to set a callback function that is called when
// the dialog is closed
func (d *dialog) SetOnClosed(closed func()) {
	// if there is already a callback set, remember it and call both
	originalCallback := d.callback

	d.callback = func(response bool) {
		closed()
		if originalCallback != nil {
			originalCallback(response)
		}
	}
}

func (d *dialog) hideWithResponse(resp bool) {
	d.win.Hide()
	if d.callback != nil {
		d.callback(resp)
	}
}

func (d *dialog) setButtons(buttons fyne.CanvasObject) {
	d.bg = newThemedBackground()
	d.label = widget.NewLabelWithStyle(d.title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	var content fyne.CanvasObject
	if d.icon == nil {
		content = container.New(d.layout,
			&canvas.Image{},
			d.bg,
			d.content,
			buttons,
			d.label,
		)
	} else {
		bgIcon := canvas.NewImageFromResource(d.icon)
		content = container.New(d.layout,
			bgIcon,
			d.bg,
			d.content,
			buttons,
			d.label,
		)
	}

	d.win = widget.NewModalPopUp(content, d.parent.Canvas())
	d.Refresh()
}

func newDialog(title, message string, icon fyne.Resource, callback func(bool), parent fyne.Window) *dialog {
	d := &dialog{content: newLabel(message), title: title, icon: icon, parent: parent}
	d.layout = &dialogLayout{d: d}

	d.callback = callback

	return d
}

func newLabel(message string) fyne.CanvasObject {
	return widget.NewLabelWithStyle(message, fyne.TextAlignCenter, fyne.TextStyle{})
}

func newButtonList(buttons ...*widget.Button) fyne.CanvasObject {
	list := container.New(layout.NewGridLayout(len(buttons)))
	for _, button := range buttons {
		list.Add(button)
	}

	return list
}

// ===============================================================
// ThemedBackground
// ===============================================================

type themedBackground struct {
	widget.BaseWidget
}

func newThemedBackground() *themedBackground {
	t := &themedBackground{}
	t.ExtendBaseWidget(t)
	return t
}

func (t *themedBackground) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	rect := canvas.NewRectangle(theme.BackgroundColor())
	return &themedBackgroundRenderer{rect, []fyne.CanvasObject{rect}}
}

type themedBackgroundRenderer struct {
	rect    *canvas.Rectangle
	objects []fyne.CanvasObject
}

func (renderer *themedBackgroundRenderer) Destroy() {
}

func (renderer *themedBackgroundRenderer) Layout(size fyne.Size) {
	renderer.rect.Resize(size)
}

func (renderer *themedBackgroundRenderer) MinSize() fyne.Size {
	return renderer.rect.MinSize()
}

func (renderer *themedBackgroundRenderer) Objects() []fyne.CanvasObject {
	return renderer.objects
}

func (renderer *themedBackgroundRenderer) Refresh() {
	r, g, b, _ := ToNRGBA(theme.BackgroundColor())
	bg := &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 230}
	renderer.rect.FillColor = bg
}

// ===============================================================
// DialogLayout
// ===============================================================

type dialogLayout struct {
	d *dialog
}

func (l *dialogLayout) Layout(obj []fyne.CanvasObject, size fyne.Size) {
	l.d.bg.Move(fyne.NewPos(0, 0))
	l.d.bg.Resize(size)

	btnMin := obj[3].MinSize()

	// icon
	iconHeight := padHeight*2 + l.d.label.MinSize().Height*2 - theme.Padding()
	obj[0].Resize(fyne.NewSize(iconHeight, iconHeight))
	obj[0].Move(fyne.NewPos(size.Width-iconHeight+theme.Padding(), -theme.Padding()))

	// buttons
	obj[3].Resize(btnMin)
	obj[3].Move(fyne.NewPos(size.Width/2-(btnMin.Width/2), size.Height-padHeight-btnMin.Height))

	// content
	contentStart := l.d.label.Position().Y + l.d.label.MinSize().Height + padHeight
	contentEnd := obj[3].Position().Y - theme.Padding()
	obj[2].Move(fyne.NewPos(padWidth/2, l.d.label.MinSize().Height+padHeight))
	obj[2].Resize(fyne.NewSize(size.Width-padWidth, contentEnd-contentStart))
}

func (l *dialogLayout) MinSize(obj []fyne.CanvasObject) fyne.Size {
	contentMin := obj[2].MinSize()
	btnMin := obj[3].MinSize()

	width := fyne.Max(fyne.Max(contentMin.Width, btnMin.Width), obj[4].MinSize().Width) + padWidth
	height := contentMin.Height + btnMin.Height + l.d.label.MinSize().Height + theme.Padding() + padHeight*2

	return fyne.NewSize(width, height)
}

func ToNRGBA(c color.Color) (r, g, b, a int) {
	// We use UnmultiplyAlpha with RGBA, RGBA64, and unrecognized implementations of Color.
	// It works for all Colors whose RGBA() method is implemented according to spec, but is only necessary for those.
	// Only RGBA and RGBA64 have components which are already premultiplied.
	switch col := c.(type) {
	// NRGBA and NRGBA64 are not premultiplied
	case color.NRGBA:
		r = int(col.R)
		g = int(col.G)
		b = int(col.B)
		a = int(col.A)
	case *color.NRGBA:
		r = int(col.R)
		g = int(col.G)
		b = int(col.B)
		a = int(col.A)
	case color.NRGBA64:
		r = int(col.R) >> 8
		g = int(col.G) >> 8
		b = int(col.B) >> 8
		a = int(col.A) >> 8
	case *color.NRGBA64:
		r = int(col.R) >> 8
		g = int(col.G) >> 8
		b = int(col.B) >> 8
		a = int(col.A) >> 8
	// Gray and Gray16 have no alpha component
	case *color.Gray:
		r = int(col.Y)
		g = int(col.Y)
		b = int(col.Y)
		a = 0xff
	case color.Gray:
		r = int(col.Y)
		g = int(col.Y)
		b = int(col.Y)
		a = 0xff
	case *color.Gray16:
		r = int(col.Y) >> 8
		g = int(col.Y) >> 8
		b = int(col.Y) >> 8
		a = 0xff
	case color.Gray16:
		r = int(col.Y) >> 8
		g = int(col.Y) >> 8
		b = int(col.Y) >> 8
		a = 0xff
	// Alpha and Alpha16 contain only an alpha component.
	case color.Alpha:
		r = 0xff
		g = 0xff
		b = 0xff
		a = int(col.A)
	case *color.Alpha:
		r = 0xff
		g = 0xff
		b = 0xff
		a = int(col.A)
	case color.Alpha16:
		r = 0xff
		g = 0xff
		b = 0xff
		a = int(col.A) >> 8
	case *color.Alpha16:
		r = 0xff
		g = 0xff
		b = 0xff
		a = int(col.A) >> 8
	default: // RGBA, RGBA64, and unknown implementations of Color
		r, g, b, a = unmultiplyAlpha(c)
	}
	return
}

func unmultiplyAlpha(c color.Color) (r, g, b, a int) {
	red, green, blue, alpha := c.RGBA()
	if alpha != 0 && alpha != 0xffff {
		red = (red * 0xffff) / alpha
		green = (green * 0xffff) / alpha
		blue = (blue * 0xffff) / alpha
	}
	// Convert from range 0-65535 to range 0-255
	r = int(red >> 8)
	g = int(green >> 8)
	b = int(blue >> 8)
	a = int(alpha >> 8)
	return
}
