package navigations

import (
	"strconv"
	"kai-suite/types"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"
	custom_widget "kai-suite/widgets"
)

type ThreadCardCached struct {
	Timestamp int
	Card fyne.CanvasObject
}

var (
	threadsCardCache map[int]*ThreadCardCached
	threadsBox *fyne.Container
	threadsContainer *fyne.Container
	messagesBox *fyne.Container
	messagesContainer *fyne.Container
)

func threadsRepository() []types.MozMobileMessageThread {
	t1 := types.MozMobileMessageThread{
		Id: 3,
		Body: "Help",
		UnreadCount: 0,
		Participants: []string{"Hotlink@"}, 
		Timestamp: 1647463907980,
		LastMessageSubject: "",
		LastMessageType: "sms",
	}
	t2 := types.MozMobileMessageThread{
		Id: 2,
		Body: "To check balance, reply CHECK. To e…",
		UnreadCount: 0,
		Participants: []string{"20505"}, 
		Timestamp: 1647463907980,
		LastMessageSubject: "",
		LastMessageType: "1647463583560",
	}
	return []types.MozMobileMessageThread{t1, t2}
}

func ViewMessagesThread(threadId int, threads []types.MozMobileMessageThread) {
	
}

func viewThread(threadId int) {
	threadsBox.Hide()
	messagesBox.Show()
	messagesContainer.Objects = nil
	messagesContainer.Add(widget.NewLabel("T3 Message 1"))
	messagesContainer.Add(widget.NewLabel("T3 Message 2"))
	messagesContainer.Add(widget.NewLabel("T3 Message 3"))
	messagesContainer.Add(widget.NewLabel("T3 Message 4"))
	messagesContainer.Add(widget.NewLabel("T3 Message 5"))
	messagesContainer.Add(widget.NewLabel("T3 Message 6"))
}

func RenderMessagesContent(c *fyne.Container) {
	log.Info("Messages Rendered")
	c.Hide()
	threadsCardCache = make(map[int]*ThreadCardCached)
	threadsContainer = container.NewVBox()
	for _, t := range threadsRepository() {
		if _, exist := threadsCardCache[t.Id]; exist == true {
			if threadsCardCache[t.Id].Timestamp != t.Timestamp {
				threadsCardCache[t.Id].Timestamp = t.Timestamp
				threadsCardCache[t.Id] = &ThreadCardCached{}
				threadsCardCache[t.Id].Timestamp = t.Timestamp
				card := &widget.Card{}
				card.SetTitle(t.Body)
				card.SetSubTitle(t.Participants[0])
				card.SetContent(container.NewHBox(
					custom_widget.NewButton(strconv.Itoa(t.Id), "View", func(scope string) {
					log.Info("Clicked view ", scope)
					if i, err := strconv.Atoi(scope); err == nil {
						viewThread(i)
					}
				}),
				))
				threadsCardCache[t.Id].Card = card
			}
		} else {
			threadsCardCache[t.Id] = &ThreadCardCached{}
			threadsCardCache[t.Id].Timestamp = t.Timestamp
			card := &widget.Card{}
			card.SetTitle(t.Body)
			card.SetSubTitle(t.Participants[0])
			card.SetContent(container.NewHBox(
				custom_widget.NewButton(strconv.Itoa(t.Id), "View", func(scope string) {
					log.Info("Clicked view ", scope)
					if i, err := strconv.Atoi(scope); err == nil {
						viewThread(i)
					}
				}),
			))
			threadsCardCache[t.Id].Card = card
		}
		threadsContainer.Add(threadsCardCache[t.Id].Card)
	}
	messagesContainer = container.NewVBox()
	threadsBox = container.NewBorder(
		nil, nil, nil, nil,
		container.NewVScroll(container.NewVBox(threadsContainer)),
	)
	messagesBox = container.NewBorder(
		nil,
		container.NewBorder(
			nil, nil,
			widget.NewButton("Return", func() {
				threadsBox.Show()
				messagesBox.Hide()
			}),
			widget.NewButton("SEND", func(){}),
			widget.NewMultiLineEntry(),
		), nil, nil,
		container.NewVScroll(container.NewVBox(messagesContainer)),
	)
	messagesBox.Hide()
	c.Add(threadsBox)
	c.Add(messagesBox)
}
