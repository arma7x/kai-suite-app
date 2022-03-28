package navigations

import (
	"net"
	"net/url"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"kai-suite/utils/global"
	"kai-suite/utils/websockethub"
	log "github.com/sirupsen/logrus"
	"fyne.io/fyne/v2/dialog"
)

var StatusText = "Disconnected"
var inputIp *widget.Select
var inputPort *widget.Entry
var DeviceStatus = make(chan bool)
var ConnectionStatus = make(chan bool)

var deviceLabel = widget.NewLabel("Device: -")
var ipPortLabel = widget.NewLabel("")
var buttonConnect = widget.NewButton("Connect", func() {
	if websockethub.Status == false {
		addr, err := global.CheckIPAddress(inputIp.Selected, inputPort.Text)
		if err != nil {
			log.Warn(err)
			ipPortLabel.SetText(err.Error())
			dialog.ShowError(err, global.WINDOW)
			return
		} else {
			ipPortLabel.SetText("Ip Address: " + addr)
			go websockethub.Start(addr)
		}
	} else {
		websockethub.Stop()
	}
})

func getNetworkCardIPAddresses() (ipList []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipList = append(ipList, ipnet.IP.String())
			}
		}
	}
	return
}

func RenderConnectionContent(c *fyne.Container) {
	websockethub.RegisterListener(ReloadThreads, ReloadMessages, RemoveContact, RefreshThreads)
	log.Info("Connection Rendered")
	go func() {
		for {
			select {
				case <- websockethub.GetClientConnectedChan():
					if websockethub.Client != nil {
						deviceLabel.SetText("Device: " + websockethub.Client.GetDevice() + " " + websockethub.Client.GetIMEI())
						DeviceStatus <- true
					} else {
						deviceLabel.SetText("Device: -")
						DeviceStatus <- false
					}
				case status := <- websockethub.GetConnectionChan():
					if (websockethub.Status) {
						StatusText = "Connected"
						buttonConnect.SetText("Disconnect")
						log.Info("Connected ", status)
						ConnectionStatus <- true
					} else {
						StatusText = "Disconnected"
						buttonConnect.SetText("Connect")
						log.Info("Disconnected ", status)
						ConnectionStatus <- false
					}
			}
		}
	}()
	c.Objects = nil
	inputPort = widget.NewEntry()
	inputPort.Text = "4444"
	inputIp = widget.NewSelect(getNetworkCardIPAddresses(), func(selected string) {
		ipPortLabel.SetText("Ip Address: " + inputIp.Selected + ":" + inputPort.Text)
	})
	inputIp.PlaceHolder = "(Select network card)"
	inputPort.OnChanged = func(val string) {
		ipPortLabel.SetText("Ip Address: " + inputIp.Selected + ":" + inputPort.Text)
	}
	ipPortLabel.SetText("Ip Address: " + inputIp.Selected + ":" + inputPort.Text)
	//if websockethub.Status == false {
	//	contentTitle.Set("Disconnected")
	//} else {
	//	contentTitle.Set("Connected")
	//}
	c.Add(
		container.NewVBox(
			deviceLabel,
			ipPortLabel,
			inputIp,
			inputPort,
			buttonConnect,
			widget.NewButton("Test Connection", func() {
				if addr, err := global.CheckIPAddress(inputIp.Selected, inputPort.Text); err == nil && websockethub.Status == true {
					testURL, _ := url.Parse("http://" + addr)
					if err := fyne.CurrentApp().OpenURL(testURL); err != nil {
						log.Warn(err)
					}
				}
			}),
		),
	)
}
