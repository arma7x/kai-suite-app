package navigations

import (
	"strconv"
	"kai-suite/types"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/layout"
	log "github.com/sirupsen/logrus"
	custom_widget "kai-suite/widgets"
)

type ThreadCardCached struct {
	Timestamp int
	Card fyne.CanvasObject
}

type MessageCardCached struct {
	Timestamp int
	Card fyne.CanvasObject
}

var (
	Threads	map[int]types.MozMobileMessageThread
	Messages	map[int]map[int]types.MozSmsMessage
	threadsCardCache map[int]*ThreadCardCached
	messagesCardCache map[int]map[int]*MessageCardCached
	threadsBox *fyne.Container
	threadsContainer *fyne.Container
	messagesBox *fyne.Container
	messagesContainer *fyne.Container
)

func threadsRepository() []types.MozMobileMessageThread {
	t1 := types.MozMobileMessageThread{
		Id: 1,
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

func messagesRepository(threadId int) []types.MozSmsMessage {
	messages := make(map[int][]types.MozSmsMessage)
	t1m1 := types.MozSmsMessage{Type: "sms", Id: 1, ThreadId: 1, Body: "This is received body of thread 1 msg 1", Delivery: "received", DeliveryStatus: "success", Read: true, Receiver: "", Sender: "20505", Timestamp: 1, MessageClass: "normal", }
	t1m2 := types.MozSmsMessage{Type: "sms", Id: 2, ThreadId: 1, Body: "This is sent body of thread 1 msg 2", Delivery: "sent", DeliveryStatus: "success", Read: true, Receiver: "20505", Sender: "", Timestamp: 1, MessageClass: "normal", }
	messages[1] = append(messages[1], t1m1)
	messages[1] = append(messages[1], t1m2)
	t2m1 := types.MozSmsMessage{Type: "sms", Id: 1, ThreadId: 2, Body: "This is received body of thread 2 msg 1", Delivery: "received", DeliveryStatus: "success", Read: true, Receiver: "", Sender: "15505", Timestamp: 1, MessageClass: "normal", }
	t2m2 := types.MozSmsMessage{Type: "sms", Id: 2, ThreadId: 2, Body: "This is sent body of thread 2 msg 2", Delivery: "sent", DeliveryStatus: "success", Read: true, Receiver: "15505", Sender: "", Timestamp: 1, MessageClass: "normal", }
	messages[2] = append(messages[2], t2m1)
	messages[2] = append(messages[2], t2m2)
	return messages[threadId]
}

func ViewMessagesThread(threadId int) {
	threadsBox.Hide()
	messagesBox.Show()
	messagesContainer.Objects = nil
	if _, exist := messagesCardCache[threadId]; exist == false {
		messagesCardCache[threadId] = make(map[int]*MessageCardCached)
	}
	if _, exist := Messages[threadId]; exist == true {
		for _, m := range Messages[threadId] {
			if item, exist := messagesCardCache[threadId][m.Id]; exist == false {
				messagesCardCache[threadId][m.Id] = &MessageCardCached{}
				messagesCardCache[threadId][m.Id].Timestamp = m.Timestamp
				card := &widget.Card{}
				card.SetSubTitle(m.Body)
				if m.Receiver != "" {
					messagesCardCache[threadId][m.Id].Card = container.NewHBox(
						layout.NewSpacer(),
						card,
					)
				} else {
					messagesCardCache[threadId][m.Id].Card = container.NewHBox(
						card,
						layout.NewSpacer(),
					)
				}
				messagesContainer.Add(messagesCardCache[threadId][m.Id].Card)
			} else {
				messagesContainer.Add(item.Card)
			}
		}
	}
}

func RefreshThreads() {
	threadsContainer.Objects = nil
	for _, t := range Threads {
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
						ViewMessagesThread(i)
					}
				}),
				))
				threadsCardCache[t.Id].Card = card
			}
			log.Info("Cached ", t.Id);
		} else {
			log.Info("Load ", t.Id);
			threadsCardCache[t.Id] = &ThreadCardCached{}
			threadsCardCache[t.Id].Timestamp = t.Timestamp
			card := &widget.Card{}
			card.SetTitle(t.Body)
			card.SetSubTitle(t.Participants[0])
			card.SetContent(container.NewHBox(
				custom_widget.NewButton(strconv.Itoa(t.Id), "View", func(scope string) {
					log.Info("Clicked view ", scope)
					if i, err := strconv.Atoi(scope); err == nil {
						ViewMessagesThread(i)
					}
				}),
			))
			threadsCardCache[t.Id].Card = card
		}
		threadsContainer.Add(threadsCardCache[t.Id].Card)
	}
}

func RenderMessagesContent(c *fyne.Container) {
	log.Info("Messages Rendered")
	c.Hide()
	threadsCardCache = make(map[int]*ThreadCardCached)
	messagesCardCache = make(map[int]map[int]*MessageCardCached)
	threadsContainer = container.NewVBox()
	RefreshThreads()
	messagesContainer = container.NewVBox()
	threadsBox = container.NewBorder(
		nil, nil, nil, nil,
		container.NewVScroll(container.NewVBox(threadsContainer)),
	)
	messagesBox = container.NewBorder(
		nil,
		container.NewBorder(
			nil, nil, nil,
			container.NewVBox(
				widget.NewButton("RETURN", func() {
					threadsBox.Show()
					messagesBox.Hide()
				}),
				widget.NewButton("SEND", func(){}),
			),
			widget.NewMultiLineEntry(),
		), nil, nil,
		container.NewVScroll(container.NewVBox(messagesContainer)),
	)
	messagesBox.Hide()
	c.Add(threadsBox)
	c.Add(messagesBox)
}
