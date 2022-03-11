package types

import(
	"golang.org/x/oauth2"
	google_oauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/people/v1"
)

type UserInfoAndToken struct {
	User	*google_oauth2.Userinfo	`json:"user"`
	Token	*oauth2.Token 					`json:"token"`
}

type WebsocketRxMessageFlag struct {
	Flag int				`json:"flag"`
	Message string	`json:"message"`
}

type RxClientFlag0 string

type Metadata struct {
	SyncID 				string	`json:"sync_id,omitempty"`			//KaiContact.id
	SyncUpdated		string	`json:"sync_updated,omitempty"`	//KaiContact.updated
	Hash					string	`json:"hash,omitempty"`
	Deleted 			bool		`json:"deleted,omitempty"`
}


// QUEUE []TxSyncContact

type KaiContact struct {}

type TxSyncContact1 struct {
	Namespace string				`json:"namespace"`	//account:people:id
	Metadata								`json:"metadata"`
	Person *people.Person		`json:"person"`
}

// On Rx, pop QUEUE, next QUEUE
// successfully add or update contact data on kaios
type RxSyncContactFlag2 struct {
	Namespace			string	`json:"namespace"`	//account:people:id
	SyncID				string	`json:"sync_id"`
	SyncUpdated		string	`json:"sync_updated"`
}

// On Rx, add KaiContact to desktop local contacts, push QUEUE, next QUEUE
// received kaicontact from kaios then TxSync for comformation
type RxSyncContactFlag4 struct {
	Namespace			string	`json:"namespace"`	//local:people:KaiContact.id
	KaiContact						`json:"kai_contact"`
}

// On Rx, pop QUEUE, next QUEUE
// both desktop & kaios remove this contact
type RxSyncContactFlag6 struct {
	Namespace			string	`json:"namespace"`	//account:people:id
}
