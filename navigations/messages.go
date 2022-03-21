package navigations

import (
	"sort"
	"time"
	"strings"
	"strconv"
	"kai-suite/types"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/data/binding"
	log "github.com/sirupsen/logrus"
	custom_widget "kai-suite/widgets"
	"kai-suite/utils/global"
	"fyne.io/fyne/v2/dialog"
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
	FocusedThread int
	Threads	map[int]*types.MozMobileMessageThread
	Messages	map[int][]*types.MozSmsMessage
	threadsCardCache map[int]*ThreadCardCached
	messagesCardCache map[int]map[int]*MessageCardCached
	threadsBox *fyne.Container
	threadsContainer *fyne.Container
	messagesBox *fyne.Container
	messagesContainer *fyne.Container
	messagesScroller *container.Scroll
	textMessageEntry = widget.NewMultiLineEntry()
	textMessage = binding.NewString()
)

func ViewMessagesThread(threadId int) {
	log.Info("View thread: ", threadId)
	FocusedThread = threadId
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
				var richText string
				words := strings.Split(m.Body, " ")
				for i, word := range words {
					if i != 0 && i %10 == 0 {
						richText += "\n\n " + word + " "
					} else {
						richText += word + " "
					}
				}
				card.SetContent(widget.NewRichTextFromMarkdown(richText))
				if m.Receiver != "" {
					messagesCardCache[threadId][m.Id].Card = container.NewBorder(nil,nil,nil,card)
				} else {
					messagesCardCache[threadId][m.Id].Card = container.NewBorder(nil,nil,card,nil)
				}
				log.Info("Load Message ", threadId, ": ", m.Id)
				messagesContainer.Add(messagesCardCache[threadId][m.Id].Card)
			} else {
				log.Info("Cached Message ", threadId, ": ", m.Id)
				messagesContainer.Add(item.Card)
			}
		}
		messagesContainer.Refresh()
	}
	time.AfterFunc(time.Second / 2, messagesScroller.ScrollToBottom)
}

func RefreshThreads() {
	log.Info("Refresh Threads ", len(Threads))
	threadsContainer.Objects = nil
	var sortedThreads []*types.MozMobileMessageThread
	for _, t := range Threads {
		log.Info("Threads ", t.Id, " ", len(sortedThreads))
		sortedThreads = append(sortedThreads, t)
	}
	if len(sortedThreads) > 1 {
		sort.Slice(sortedThreads, func(i, j int) bool {
			return sortedThreads[i].Timestamp > sortedThreads[j].Timestamp
		})
	}
	for _, t := range sortedThreads {
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
					if i, err := strconv.Atoi(scope); err == nil {
						ViewMessagesThread(i)
					}
				}),
				))
				threadsCardCache[t.Id].Card = card
			}
			log.Info("Cached Thread ", t.Id)
		} else {
			log.Info("Load Thread ", t.Id)
			threadsCardCache[t.Id] = &ThreadCardCached{}
			threadsCardCache[t.Id].Timestamp = t.Timestamp
			card := &widget.Card{}
			card.SetTitle(t.Body)
			card.SetSubTitle(t.Participants[0])
			card.SetContent(container.NewHBox(
				custom_widget.NewButton(strconv.Itoa(t.Id), "View", func(scope string) {
					if i, err := strconv.Atoi(scope); err == nil {
						ViewMessagesThread(i)
					}
				}),
			))
			threadsCardCache[t.Id].Card = card
		}
		threadsContainer.Add(threadsCardCache[t.Id].Card)
	}
	threadsContainer.Refresh()
	if FocusedThread != 0 {
		ViewMessagesThread(FocusedThread)
	}
}

func RenderMessagesContent(c *fyne.Container, sendSMSCb func(receivers []string, message string, iccId string)) {
	var newDialog dialog.Dialog
	to := widget.NewEntry()
	body := widget.NewMultiLineEntry()
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "To", Widget: to},
			{Text: "Body", Widget: body},
		},
		SubmitText: "Send",
		OnSubmit: func() {
			log.Info(to.Text, " ", body.Text)
			if to.Text != "" && body.Text != "" {
				sendSMSCb([]string{to.Text}, body.Text, "")
			}
			newDialog.Hide()
		},
	}
	newDialog = dialog.NewCustom("New Message", "Cancel", container.NewMax(form), global.WINDOW);
	log.Info("Messages Rendered")
	c.Hide()
	textMessageEntry.Bind(textMessage)
	threadsCardCache = make(map[int]*ThreadCardCached)
	messagesCardCache = make(map[int]map[int]*MessageCardCached)
	threadsContainer = container.NewVBox()
	messagesContainer = container.NewVBox()
	threadsBox = container.NewBorder(
		container.NewHBox(
			widget.NewButton("New Message", func() {
				newDialog.Show()
				sz := newDialog.MinSize()
				sz.Width = 400
				newDialog.Resize(sz)
			}),
		),
		nil, nil, nil,
		container.NewVScroll(container.NewVBox(threadsContainer)),
	)
	messagesScroller = container.NewVScroll(container.NewVBox(messagesContainer))
	messagesBox = container.NewBorder(
		nil,
		container.NewBorder(
			nil, nil, nil,
			container.NewVBox(
				widget.NewButton("RETURN", func() {
					FocusedThread = 0
					threadsBox.Show()
					messagesBox.Hide()
					RefreshThreads()
				}),
				widget.NewButton("SEND", func(){
					if FocusedThread != 0 {
						text, _ := textMessage.Get()
						if text != "" {
							sendSMSCb(Threads[FocusedThread].Participants, text, Messages[FocusedThread][0].IccId)
							textMessage.Set("")
						}
					}
				}),
			),
			textMessageEntry,
		), nil, nil,
		messagesScroller,
	)
	messagesBox.Hide()
	c.Add(threadsBox)
	c.Add(messagesBox)
}
