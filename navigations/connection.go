package navigations

import (
	"net"
	"strconv"
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

var websocketBtnTxtChan = make(chan string)
var websocketClientConnectedChan = make(chan bool)

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
			websockethub.Init(addr, websocketClientConnectedChan, ReloadThreads, ReloadMessages, RemoveContact, RefreshThreads)
			go websockethub.Start(onStatusChange)
		}
	} else {
		websockethub.Stop(onStatusChange)
	}
})

func onStatusChange(status bool, err error) {
	if (status) {
		log.Info("Connected")
		StatusText = "Connected"
		websocketBtnTxtChan <- "Disconnect"
	} else {
		log.Info("Disconnected")
		StatusText = "Disconnected"
		websocketBtnTxtChan <- "Connect"
	}
	if err != nil {
		log.Warn(strconv.FormatBool(status), err.Error())
	}
}

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
	log.Info("Connection Rendered")
	go func() {
		for {
			select {
				case txt := <- websocketBtnTxtChan:
					buttonConnect.SetText(txt)
				case present := <- websocketClientConnectedChan:
					if present {
						deviceLabel.SetText("Device: " + websockethub.Client.GetDevice() + " " + websockethub.Client.GetIMEI())
					} else {
						deviceLabel.SetText("Device: -")
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
		),
	)
}
