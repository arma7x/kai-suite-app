package google

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

	r, err := srv.People.Connections.List("people/me").PageSize(10).
		PersonFields("names,emailAddresses").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve people. %v", err)
	}
	if len(r.Connections) > 0 {
		fmt.Print("List 10 connection names:\n")
		for _, c := range r.Connections {
			names := c.Names
			if len(names) > 0 {
				name := names[0].DisplayName
				fmt.Printf("%s\n", name)
			}
		}
	} else {
		fmt.Print("No connections found.")
	}
}
