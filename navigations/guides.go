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
		widget.NewSeparator(),
		widget.NewRichTextFromMarkdown("#	Connection"),
		widget.NewLabel("1. Use ifconfig(linux) or ipconfig(windows) to get your Wi-Fi ip address"),
		widget.NewLabel("2. Your computer and KaiOS device must connect to the same network. Please setup"),
		widget.NewLabel("port forwarding, if your computer not connected to KaiOS Wi-Fi hotspot"),
		widget.NewSeparator(),
		widget.NewRichTextFromMarkdown("#	Local Contacts"),
		widget.NewLabel("1. The origin of contact is from KaiOS Device/VCF"),
		widget.NewLabel("2. Please use Restore, if you accidentally delete any contacts on your device"),
		widget.NewLabel("or when your KaiOS device is connected to Kai Suite for the first time"),
		widget.NewSeparator(),
		widget.NewRichTextFromMarkdown("#	Google Contacts"),
		widget.NewLabel("1. The origin of contact is from Google People"),
		widget.NewLabel("2. Please use Restore, if you accidentally delete any contacts on your device"),
		widget.NewLabel("or when your KaiOS device is connected to Kai Suite for the first time"),
		widget.NewSeparator(),
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
