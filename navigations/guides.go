package navigations

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"
)

func RenderGuidesContent(c *fyne.Container) {
	log.Info("Helps Rendered")
	c.Hide()
	guidesContent := container.NewVBox(
		widget.NewRichTextFromMarkdown("#	Connections"),
		widget.NewLabel("~ Use ifconfig(linux) or ipconfig(windows) to get your wi-fi ip address"),
		widget.NewLabel("~ Please setup port forwarding, if your pc not using/connected to KaiOS hotspot/tethering"),
		widget.NewRichTextFromMarkdown("#	Local Contacts"),
		widget.NewLabel("~ The origin of contact is KaiOS Device/VCF"),
		widget.NewLabel("~ Please use Restore, if you accidentally delete any contacts on your device"),
		widget.NewLabel("or when the KaiOS device is connected to Kai Suite for the first time"),
		widget.NewRichTextFromMarkdown("#	Google Contacts"),
		widget.NewLabel("~ The origin of contact is Google People API"),
		widget.NewLabel("~ Please use Restore, if you accidentally delete any contacts on yourdevice"),
		widget.NewLabel("or when the KaiOS device is connected to Kai Suite for the first time"),
		widget.NewRichTextFromMarkdown("#	Setup Google API"),
		widget.NewLabel("~ TBD"),
	)
	contentScroller := container.NewVScroll(guidesContent)
	c.Add(contentScroller)
}
