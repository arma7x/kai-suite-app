package google_services

import (
	"context"
	"time"
	"kai-suite/types"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/calendar/v3"
)

func SyncCalendar(config *oauth2.Config, account *types.UserInfoAndToken) error {
	ctx := context.Background()
	client := GetAuthClient(config, account.Token)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Error("Unable to retrieve Calendar client: ", err)
		return err
	}

	t := time.Now().Format(time.RFC3339)
	run := true
	var loopError error
	var events []*calendar.Event // type Person struct
	var r *calendar.Events

	r, loopError = srv.Events.List("primary").ShowDeleted(false).SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	for (run) {
		if loopError != nil {
			log.Error("Unable to retrieve next ten of the user's events: ", loopError)
			run = false
		} else {
			events = append(events, r.Items...)
			if r.NextPageToken == "" {
				run = false
			} else {
				log.Info(r.NextPageToken)
				r, loopError = srv.Events.List("primary").ShowDeleted(false).SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").PageToken(r.NextPageToken).Do()
			}
		}
	}
	for _, item := range events {
		date_start := item.Start.DateTime
		if date_start == "" {
			date_start = item.Start.Date
		}
		date_end := item.End.DateTime
		if date_end == "" {
			date_end = item.End.Date
		}
		log.Printf("%v, %v, %v, %v, %v\n", item.Id, item.Summary, item.Description, date_start, date_end)
	}
	return nil
}
