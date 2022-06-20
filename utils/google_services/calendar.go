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

func SyncCalendar(config *oauth2.Config, account *types.UserInfoAndToken, unsync_events []*calendar.Event) ([]*calendar.Event, []*calendar.Event, error) {
	var failEvents []*calendar.Event // type Event struct
	log.Info("Sync Calendars ", account.User.Id, ' ', len(unsync_events))
	ctx := context.Background()
	client := GetAuthClient(config, account.Token)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Error("Unable to retrieve Calendar client: ", err)
		return nil, nil, err
	}

	if len(unsync_events) > 0 {
		evtSrv := calendar.NewEventsService(srv);
		for _, evt := range unsync_events {
			evt.Id = ""
			if _, err := evtSrv.Insert("primary", evt).Do(); err != nil {
				log.Error("Err insert ", err)
				failEvents = append(failEvents, evt)
			}
		}
	}

	t := time.Now().Format(time.RFC3339)
	run := true
	var loopError error
	var events []*calendar.Event // type Event struct
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
	return failEvents, events, err
}
