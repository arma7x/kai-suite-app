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
	events, err := srv.Events.List("primary").ShowDeleted(false).SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Error("Unable to retrieve next ten of the user's events: ", err)
		return err
	}
	log.Info("Upcoming events:")
	if len(events.Items) == 0 { // type Event struct
		log.Info("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			log.Info("%v (%v)\n", item.Summary, date)
		}
	}
	return nil
}
