package google_services

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

var (
	fields = "names,phoneNumbers,emailAddresses,addresses,birthdays,metadata"
	updateFields = "names,phoneNumbers,emailAddresses,addresses,birthdays"
)

func GetContacts(client *http.Client) []*people.Person {
	ctx := context.Background()

	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Warn("Unable to create people Client: ", err)
	}

	run := true;
	var connections []*people.Person // type Person struct
	var r *people.ListConnectionsResponse
	var rErr error
	r, rErr = srv.People.Connections.List("people/me").PageSize(1000).PersonFields(fields).Do()
	for (run) {
		if rErr != nil {
			log.Warn("Unable to retrieve people: ", err)
			run = false
		} else {
			if r.NextPageToken != "" {
				log.Info(r.NextPageToken, "\n")
			}
			connections = append(connections, r.Connections...)
			if r.NextPageToken == "" {
				run = false
			} else {
				r, rErr = srv.People.Connections.List("people/me").PageSize(20).PersonFields(fields).PageToken(r.NextPageToken).Do()
			}
		}
	}
	return connections
}

func CreateContacts() {}

func UpdateContacts(client *http.Client, contacts map[string]people.Person) ([]*people.Person, []*people.Person) {
	ctx := context.Background()
	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Warn("Unable to create people Client: ", err)
	}
	var success []*people.Person
	var fail		[]*people.Person
	for key, pr := range contacts {
		if p, err := srv.People.UpdateContact(key, &pr).UpdatePersonFields(updateFields).Do(); err == nil {
			log.Info("SUCCESS Person: ", p.Metadata.Sources[0].UpdateTime)
			success = append(success, p)
		} else {
			log.Info("FAIL Person: ", err)
			fail = append(fail, p)
		}
	}
	return success, fail
	//batch := people.BatchUpdateContactsRequest{
		//Contacts: contacts,
		//UpdateMask: updateFields,
	//}
	//if response, err := srv.People.BatchUpdateContacts(&batch).Do(); err == nil {
		//log.Warn("Batch update success: ", len(response.UpdateResult), " ", len(contacts))
		//for key, pr := range response.UpdateResult {
			//b, _ := pr.Person.MarshalJSON()
			//log.Info("Updated Person: ", key, ": ", string(b))
		//}
	//} else {
		//log.Warn("Batch update fails: ", err)
	//}
	// PeopleBatchUpdateContactsCall
	// PeopleUpdateContactCall
}

func DeleteContacts() {}

func SearchContacts() {}
