package widget

import(
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/widget"
)

type Button struct {
  Namespace string
  OnClick func(string)
  widget.Button
}

// NewButton creates a new button widget with the set label and tap handler
func NewButton(namespace, label string, tapped func(string)) *Button {
	button := &Button{}
	button.Namespace = namespace
	button.OnClick = tapped

	button.ExtendBaseWidget(button)
	button.Text = label
	return button
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler
func (b *Button) Tapped(evt *fyne.PointEvent) {
	if b.Disabled() {
		return
	}
	b.OnClick(b.Namespace)
	b.Button.Tapped(evt)
}
