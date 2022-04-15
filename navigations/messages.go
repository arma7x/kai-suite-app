package navigations

import (
	"errors"
	"sort"
	"time"
	"strings"
	"strconv"
	"encoding/json"
	"kai-suite/types"
	"golang.org/x/exp/utf8string"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/data/binding"
	log "github.com/sirupsen/logrus"
	custom_widget "kai-suite/widgets"
	"kai-suite/utils/global"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/canvas"
	
)

type ThreadCardCached struct {
	Timestamp int
	UnreadCount int
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
	newMessageButton *widget.Button
	newMessageDialog dialog.Dialog
	messageRecipient *widget.Entry
	messageBody *widget.Entry
	threadsCardCache map[int]*ThreadCardCached
	messagesCardCache map[int]map[int]*MessageCardCached
	threadsBox *fyne.Container
	threadsContainer *fyne.Container
	messagesBox *fyne.Container
	messagesContainer *fyne.Container
	messagesScroller *container.Scroll
	textMessageEntry = widget.NewMultiLineEntry()
	textMessage = binding.NewString()
	smsReadIdChan = make(chan []int)
	deleteMessagesFn = func(id []int){}
)

func ReloadThreads(threads map[int]*types.MozMobileMessageThread) {
	Threads = threads
	if _, focusExist := Threads[FocusedThread]; focusExist == false {
		FocusedThread = 0
		threadsBox.Show()
		messagesBox.Hide()
	}
	for _, t := range Threads {
		if t.Id != FocusedThread { //  || global.VISIBILITY == false 
			if t.UnreadCount > 0 {
				global.APP.SendNotification(fyne.NewNotification(t.Participants[0], t.Body))
			}
		}
	}
	//if global.VISIBILITY == true {
	//	global.WINDOW.RequestFocus()
	//}
}

func ReloadMessages(messages map[int][]*types.MozSmsMessage) {
	Messages = messages
}

func renderMessageMenuItem(m *types.MozSmsMessage) *custom_widget.ContextMenuButton {
	exportMenu := fyne.NewMenuItem("Copy", func() {
		global.WINDOW.Clipboard().SetContent(m.Body)
	})
	deleteMenu := fyne.NewMenuItem("Delete", func() {
		dialog.ShowCustomConfirm("Confirm", "Yes", "No", container.NewMax(widget.NewLabel("Are sure to delete this message ?")), func(action bool) {
			if action == true {
				ids := []int{m.Id}
				deleteMessagesFn(ids)
			}
		}, global.WINDOW)
	})
	menu := fyne.NewMenu("", exportMenu, deleteMenu)
	return custom_widget.NewContextMenu(theme.MoreVerticalIcon(), menu)
}

func ViewThreadMessages(threadId int) {
	log.Info("View thread: ", threadId)
	var recvSMSId []int
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
				tm := time.Unix(int64(m.Timestamp)/1000, (int64(m.Timestamp)%1000)*1000*1000).Local().Format("Mon, 02 Jan 2006 03:04 PM")
				if m.Receiver != "" {
					if m.Delivery == "error" {
						card.SetContent(
							container.NewVBox(
								container.NewHBox(
									widget.NewIcon(theme.WarningIcon()),
									layout.NewSpacer(),
									widget.NewRichTextFromMarkdown(richText),
								),
								container.NewBorder(
									nil, nil,
									container.New(layout.NewHBoxLayout(), &canvas.Text{ Text: tm, TextSize: 11}),
									renderMessageMenuItem(m),
								),
							),
						)
					} else {
						card.SetContent(
							container.NewVBox(
								container.NewHBox(
									layout.NewSpacer(),
									widget.NewRichTextFromMarkdown(richText),
								),
								container.NewBorder(
									nil, nil,
									container.New(layout.NewHBoxLayout(), &canvas.Text{ Text: tm, TextSize: 11}),
									renderMessageMenuItem(m),
								),
							),
						)
					}
					messagesCardCache[threadId][m.Id].Card = container.NewBorder(nil,nil,nil,card)
				} else {
					card.SetContent(
						container.NewVBox(
							widget.NewRichTextFromMarkdown(richText),
							container.NewBorder(
								nil, nil,
								renderMessageMenuItem(m),
								container.New(layout.NewHBoxLayout(), &canvas.Text{ Text: tm, TextSize: 11}),
							),
						),
					)
					messagesCardCache[threadId][m.Id].Card = container.NewBorder(nil,nil,card,nil)
					recvSMSId = append(recvSMSId, m.Id)
				}
				// log.Info("Load Message ", threadId, ": ", m.Id)
				messagesContainer.Add(messagesCardCache[threadId][m.Id].Card)
			} else {
				// log.Info("Cached Message ", threadId, ": ", m.Id)
				messagesContainer.Add(item.Card)
			}
		}
		smsReadIdChan <-recvSMSId
		messagesContainer.Refresh()
	}
	time.AfterFunc(time.Second / 2, messagesScroller.ScrollToBottom)
}

func renderThreadMenuItem(id int) *custom_widget.ContextMenuButton {
	exportMenu := fyne.NewMenuItem("Export", func() {
		dialog.ShowFileSave(func(f fyne.URIWriteCloser, err error) {
			if err == nil && f != nil {
				defer f.Close()
				if thread, exist := Messages[id]; exist == true {
					if b, err := json.Marshal(thread); err == nil {
						_, err := f.Write(b[:])
						if err != nil {
							dialog.ShowError(err, global.WINDOW)
						}
						err = f.Close()
						if err != nil {
							dialog.ShowError(err, global.WINDOW)
						}
						log.Info("Saved to...", f.URI())
						dialog.ShowInformation("Success", f.URI().Path(), global.WINDOW);
					} else {
						dialog.ShowError(errors.New("Unknown Error"), global.WINDOW)
					}
				} else {
					dialog.ShowError(errors.New("Thread not exist"), global.WINDOW)
				}
			} else {
				dialog.ShowError(errors.New("Unknown Error"), global.WINDOW)
			}
		}, global.WINDOW)
	})
	deleteMenu := fyne.NewMenuItem("Delete", func() {
		if messages, exist := Messages[id]; exist == true {
			dialog.ShowCustomConfirm("Confirm", "Yes", "No", container.NewMax(widget.NewLabel("Are sure to delete this thread ?")), func(action bool) {
				if action == true {
					var ids []int
					for _, m := range messages{
						ids = append(ids, m.Id)
					}
					deleteMessagesFn(ids)
				}
			}, global.WINDOW)
		}
	})
	menu := fyne.NewMenu("", exportMenu, deleteMenu)
	return custom_widget.NewContextMenu(theme.MoreVerticalIcon(), menu)
}

func RefreshThreads() {
	log.Info("Threads Length: ", len(Threads))
	threadsContainer.Objects = nil
	var sortedThreads []*types.MozMobileMessageThread
	for _, t := range Threads {
		// log.Info("Threads ", t.Id, " ", len(sortedThreads))
		sortedThreads = append(sortedThreads, t)
	}
	if len(sortedThreads) > 1 {
		sort.Slice(sortedThreads, func(i, j int) bool {
			return sortedThreads[i].Timestamp > sortedThreads[j].Timestamp
		})
	}
	for _, t := range sortedThreads {
		if _, exist := threadsCardCache[t.Id]; exist == true {
			if threadsCardCache[t.Id].Timestamp != t.Timestamp || threadsCardCache[t.Id].UnreadCount != t.UnreadCount {
				threadsCardCache[t.Id].Timestamp = t.Timestamp
				threadsCardCache[t.Id].UnreadCount = t.UnreadCount
				card := &widget.Card{}
				tBody := utf8string.NewString(t.Body)
				if tBody.RuneCount() > 50 {
					card.SetTitle(tBody.Slice(0, 50) + "...")
				} else {
					card.SetTitle(tBody.String())
				}
				card.SetSubTitle(t.Participants[0])
				if t.UnreadCount > 0 {
					card.SetSubTitle(t.Participants[0] + "(" + strconv.Itoa(t.UnreadCount) + ")")
				}
				tm := time.Unix(int64(t.Timestamp)/1000, (int64(t.Timestamp)%1000)*1000*1000).Local().Format("Mon, 02 Jan 2006 03:04 PM")
				card.SetContent(container.NewBorder(
					nil, nil,
					container.New(layout.NewHBoxLayout(), &canvas.Text{ Text: tm, TextSize: 11}),
					container.NewHBox(
						custom_widget.NewButton(strconv.Itoa(t.Id), "View", func(scope string) {
							if i, err := strconv.Atoi(scope); err == nil {
								ViewThreadMessages(i)
							}
						}),
						renderThreadMenuItem(t.Id),
					),
				))
				threadsCardCache[t.Id].Card = card
			}
			// log.Info("Cached Thread ", t.Id)
		} else {
			// log.Info("Load Thread ", t.Id)
			threadsCardCache[t.Id] = &ThreadCardCached{}
			threadsCardCache[t.Id].Timestamp = t.Timestamp
			threadsCardCache[t.Id].UnreadCount = t.UnreadCount
			card := &widget.Card{}
			tBody := utf8string.NewString(t.Body)
			if tBody.RuneCount() > 50 {
				card.SetTitle(tBody.Slice(0, 50) + "...")
			} else {
				card.SetTitle(tBody.String())
			}
			card.SetSubTitle(t.Participants[0])
			if t.UnreadCount > 0 {
				card.SetSubTitle(t.Participants[0] + "(" + strconv.Itoa(t.UnreadCount) + ")")
			}
			tm := time.Unix(int64(t.Timestamp)/1000, (int64(t.Timestamp)%1000)*1000*1000).Local().Format("Mon, 02 Jan 2006 03:04 PM")
			card.SetContent(container.NewBorder(
				nil, nil,
				container.New(layout.NewHBoxLayout(), &canvas.Text{ Text: tm, TextSize: 11}),
				container.NewHBox(
					custom_widget.NewButton(strconv.Itoa(t.Id), "View", func(scope string) {
						if i, err := strconv.Atoi(scope); err == nil {
							ViewThreadMessages(i)
						}
					}),
					renderThreadMenuItem(t.Id),
				),
			))
			threadsCardCache[t.Id].Card = card
		}
		threadsContainer.Add(threadsCardCache[t.Id].Card)
	}
	threadsContainer.Refresh()
	if FocusedThread != 0 {
		ViewThreadMessages(FocusedThread)
	}
}

func SendNewMessage(recipient, body string) {
	messageRecipient.Text = recipient
	messageBody.Text = body
	ev := &fyne.PointEvent{}
	ev.AbsolutePosition = fyne.CurrentApp().Driver().AbsolutePositionForObject(newMessageButton)
	newMessageButton.Tapped(ev)
}

func RenderMessagesContent(c *fyne.Container, syncSMSCb func(), sendSMSCb func([]string, string, string), syncSMSReadCb func([]int), deleteMessagesCb func([]int)) {
	log.Info("Messages Rendered")
	go func() {
		for {
			select {
				case ids := <- smsReadIdChan:
					syncSMSReadCb(ids)
			}
		}
	}()
	deleteMessagesFn = deleteMessagesCb
	messageRecipient = widget.NewEntry()
	messageBody = widget.NewMultiLineEntry()
	messageBody.Wrapping = fyne.TextWrapWord
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Recipient", Widget: messageRecipient},
			{Text: "Body", Widget: messageBody},
		},
		SubmitText: "Send",
		OnSubmit: func() {
			// log.Info(messageRecipient.Text, " ", messageBody.Text)
			if messageRecipient.Text != "" && messageBody.Text != "" {
				sendSMSCb([]string{messageRecipient.Text}, messageBody.Text, "")
				messageRecipient.Text = ""
				messageBody.Text = ""
			}
			newMessageDialog.Hide()
		},
	}
	newMessageDialog = dialog.NewCustom("New Message", "Cancel", container.NewMax(form), global.WINDOW);
	newMessageDialog.SetOnClosed(func() {
		messageRecipient.Text = ""
		messageBody.Text = ""
	})
	newMessageButton = widget.NewButton("New Message", func() {
		newMessageDialog.Show()
		sz := newMessageDialog.MinSize()
		sz.Width = 400
		newMessageDialog.Resize(sz)
	})
	c.Hide()
	textMessageEntry.Wrapping = fyne.TextWrapWord
	textMessageEntry.Bind(textMessage)
	threadsCardCache = make(map[int]*ThreadCardCached)
	messagesCardCache = make(map[int]map[int]*MessageCardCached)
	threadsContainer = container.NewVBox()
	messagesContainer = container.NewVBox()
	threadsBox = container.NewBorder(
		container.NewHBox(
			newMessageButton,
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
					syncSMSCb()
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
