package google_services

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

func People(client *http.Client) {
	ctx := context.Background()

	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create people Client %v", err)
	}

	run := true;
	fields := "names,phoneNumbers,emailAddresses,addresses,birthdays,metadata"
	var connections []*people.Person // type Person struct
	var r *people.ListConnectionsResponse
	var rErr error
	r, rErr = srv.People.Connections.List("people/me").PageSize(20).PersonFields(fields).Do()
	for (run) {
		if rErr != nil {
			log.Fatalf("Unable to retrieve people. %v", err)
			run = false
		} else {
			fmt.Print(r.NextPageToken, "\n")
			connections = append(connections, r.Connections...)
			if r.NextPageToken == "" {
				run = false
			} else {
				r, rErr = srv.People.Connections.List("people/me").PageSize(20).PersonFields(fields).PageToken(r.NextPageToken).Do()
			}
		}
	}
	if len(connections) > 0 {
		fmt.Print("List 10 connection names:\n")
		for i, c := range connections {
			// fmt.Print(i, " ", c.Metadata.Sources[0].UpdateTime, " ", c.Names[0].DisplayName, "\n\n")
			b, _ := c.MarshalJSON();
			fmt.Print(i, string(b), "\n\n")
		}
	} else {
		fmt.Print("No connections found.")
	}
}
