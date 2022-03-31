package navigations

import (
	"net/url"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	log "github.com/sirupsen/logrus"
)

func RenderGuidesContent(c *fyne.Container) {
	log.Info("Helps Rendered")
	c.Objects = nil
	guidesContent := container.NewVBox(
		widget.NewRichTextFromMarkdown("## Disclaimer: Please backup your messages/contacts before testing"),
		widget.NewRichTextFromMarkdown("#	Connection"),
		widget.NewLabel("~ Use ifconfig(linux) or ipconfig(windows) to get your wi-fi ip address"),
		widget.NewLabel("~ Please setup port forwarding, if your pc/laptop not connected to KaiOS hotspot"),
		widget.NewRichTextFromMarkdown("#	Local Contacts"),
		widget.NewLabel("~ The origin of contact is KaiOS Device/VCF"),
		widget.NewLabel("~ Please use Restore, if you accidentally delete any contacts on your device"),
		widget.NewLabel("or when the KaiOS device is connected to Kai Suite for the first time"),
		widget.NewRichTextFromMarkdown("#	Google Contacts"),
		widget.NewLabel("~ The origin of contact is Google People API"),
		widget.NewLabel("~ Please use Restore, if you accidentally delete any contacts on yourdevice"),
		widget.NewLabel("or when the KaiOS device is connected to Kai Suite for the first time"),
		container.NewHBox(
			widget.NewRichTextFromMarkdown("# Setup Google API"),
			widget.NewButtonWithIcon("Open in browser", theme.LogoutIcon(), func(){
				url, _ := url.Parse("https://github.com/arma7x/kai-suite-app/blob/master/README.md#guides")
				if err := fyne.CurrentApp().OpenURL(url); err != nil {
					log.Info(err)
				}
			}),
		),
		widget.NewLabel("Video tutorial https://youtu.be/Wk6pk-uRUOE"),
		widget.NewLabel("1. Create new project, visit https://console.cloud.google.com/"),
		widget.NewLabel("2. Enable People API & Calendar API"),
		widget.NewLabel("3. Configure Consent Screen"),
		widget.NewLabel("4. Create Credentials"),
		widget.NewLabel("5. Download the credential json file and rename it as credentials.json"),
		widget.NewLabel("6. Open credentials.json, search for `http://localhost` and replace"),
		widget.NewLabel("it with `urn:ietf:wg:oauth:2.0:oob`"),
		widget.NewLabel("7. The credentials.json & Kai Suite(binary file) must reside in same folder/directory"),
	)
	contentScroller := container.NewVScroll(guidesContent)
	c.Add(contentScroller)
}
