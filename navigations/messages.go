package navigations

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/people/v1"
)

type ThreadCardCached struct {
	Hash string
	Card fyne.CanvasObject
}

var (
	threadsCache map[string]string
	threadsBox *fyne.Container
	threadsContainer *fyne.Container
	messagesBox *fyne.Container
	messagesContainer *fyne.Container
)

func ViewMessagesThread(namespace string, personsArr []*people.Person) {

}

func RenderMessagesContent(c *fyne.Container) {
	log.Info("Messages Rendered")
	c.Hide()
	threadsCache = make(map[string]string)
	threadsContainer = container.NewVBox()
	threadsContainer.Add(widget.NewButton("Thread1", func() {
		threadsBox.Hide()
		messagesBox.Show()
		messagesContainer.Objects = nil
		messagesContainer.Add(widget.NewLabel("T1 Message 1"))
		messagesContainer.Add(widget.NewLabel("T1 Message 2"))
		messagesContainer.Add(widget.NewLabel("T1 Message 3"))
		messagesContainer.Add(widget.NewLabel("T1 Message 4"))
		messagesContainer.Add(widget.NewLabel("T1 Message 5"))
		messagesContainer.Add(widget.NewLabel("T1 Message 6"))
	}))
	threadsContainer.Add(widget.NewButton("Thread2", func() {
		threadsBox.Hide()
		messagesBox.Show()
		messagesContainer.Objects = nil
		messagesContainer.Add(widget.NewLabel("T2 Message 1"))
		messagesContainer.Add(widget.NewLabel("T2 Message 2"))
		messagesContainer.Add(widget.NewLabel("T2 Message 3"))
		messagesContainer.Add(widget.NewLabel("T2 Message 4"))
		messagesContainer.Add(widget.NewLabel("T2 Message 5"))
		messagesContainer.Add(widget.NewLabel("T2 Message 6"))
	}))
	threadsContainer.Add(widget.NewButton("Thread3", func() {
		threadsBox.Hide()
		messagesBox.Show()
		messagesContainer.Objects = nil
		messagesContainer.Add(widget.NewLabel("T3 Message 1"))
		messagesContainer.Add(widget.NewLabel("T3 Message 2"))
		messagesContainer.Add(widget.NewLabel("T3 Message 3"))
		messagesContainer.Add(widget.NewLabel("T3 Message 4"))
		messagesContainer.Add(widget.NewLabel("T3 Message 5"))
		messagesContainer.Add(widget.NewLabel("T3 Message 6"))
	}))
	messagesContainer = container.NewVBox()
	threadsBox = container.NewBorder(
		nil, nil, nil, nil,
		container.NewVScroll(container.NewVBox(threadsContainer)),
	)
	messagesBox = container.NewBorder(
		widget.NewButton("Return", func() {
			threadsBox.Show()
			messagesBox.Hide()
		}),
		nil, nil, nil,
		container.NewVScroll(container.NewVBox(messagesContainer)),
	)
	messagesBox.Hide()
	c.Add(threadsBox)
	c.Add(messagesBox)
}
