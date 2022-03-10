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

type Metadata struct {
	SyncID 				string	`json:"sync_id,omitempty"`			//KaiContact.id
	SyncUpdated		string	`json:"sync_updated,omitempty"`	//KaiContact.updated
	Hash					string	`json:"hash,omitempty"`
	Deleted 			bool		`json:"deleted,omitempty"`
}

type KaiContact struct {}

type TxSync struct {
	Namespace string				`json:"namespace"`	//account:people:id
	Metadata								`json:"metadata"`
	Person people.Person		`json:"person"`
}

// successfully add or update contact data on kaios
type RxSyncFlag2 struct {
	Namespace			string	`json:"namespace"`	//account:people:id
	SyncID				string	`json:"sync_id"`
	SyncUpdated		string	`json:"sync_updated"`
}

// received kaicontact from kaios then TxSync for comformation
type RxSyncFlag4 struct {
	Namespace			string	`json:"namespace"`	//local:people:KaiContact.id
	KaiContact						`json:"kai_contact"`
}

// both desktop & kaios remove this contact
type RxSyncFlag6 struct {
	Namespace			string	`json:"namespace"`	//account:people:id
}
