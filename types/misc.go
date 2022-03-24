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

type WebsocketMessageFlag struct {
	Flag int		`json:"flag"`
	Data string	`json:"data"`
}

type RxClientFlag0 string

type Metadata struct {
	SyncID 				string	`json:"sync_id,omitempty"`			//KaiContact.id
	SyncUpdated		string	`json:"sync_updated,omitempty"`	//KaiContact.updated
	Hash					string	`json:"hash,omitempty"`
	Deleted 			bool		`json:"deleted"`
}

type KaiEmail struct {
	Type 	[]string	`json:"type,omitempty"`
	Value string		`json:"value,omitempty"`
}

type KaiTel struct {
	Type 	[]string	`json:"type,omitempty"`
	Value string		`json:"value,omitempty"`
}

type KaiContact struct {
	Id 					string			`json:"id"`
	Published		string			`json:"published"`
	Updated			string			`json:"updated"`
	Email				[]KaiEmail	`json:"email,omitempty"`
	Tel					[]KaiTel		`json:"tel,omitempty"`
	Name 				[]string		`json:"name"`
	GivenName		[]string		`json:"givenName"`
	FamilyName	[]string		`json:"familyName"`
	Key					[]string		`json:"key"`
}

type LocalContactSync struct {
	KaiContact			`json:"kai_contact"`
	Metadata				`json:"metadata"`
}

type LocalContactMergedSync struct {
	Person 				people.Person		`json:"person"`
	KaiContact										`json:"kai_contact"`
	Metadata											`json:"metadata"`
}

//

type TxSyncGoogleContact struct {
	Namespace string				`json:"namespace"`	//account:people:id
	Metadata								`json:"metadata"`
	Person *people.Person		`json:"person"`
}

type TxRestoreGoogleContact3 struct {}

type TxSyncLocalContact5 struct {
	Metadata		map[string]Metadata					`json:"metadata"`
	Persons			map[string]people.Person		`json:"persons"`
}

type TxRestoreLocalContact7 struct {
	Metadata		map[string]Metadata					`json:"metadata"`
	Persons			map[string]people.Person		`json:"persons"`
}

type TxSyncSMS9 struct {}

type TxSendSMS11 struct {
	Receivers		[]string	`json:"receivers,omitempty"`
	Message			string		`json:"message,omitempty"`
	IccId				string		`json:"iccId,omitempty"`
}

type TxSyncSMSRead13 struct {
	Id []int	`json:"id"`
}

// On Rx, pop QUEUE, next QUEUE
// successfully add or update contact data on kaios
type RxSyncDevice0 struct {
	Device			string	`json:"device"`
	IMEI				string	`json:"imei"`
}

// On Rx, pop QUEUE, next QUEUE
// successfully add or update contact data on kaios
type RxSyncContactFlag2 struct {
	Namespace			string	`json:"namespace"`	//account:people:id
	SyncID				string	`json:"sync_id,omitempty"`
	SyncUpdated		string	`json:"sync_updated,omitempty"`
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

type RxSyncLocalContactFlag8 struct {
	PushList		[]LocalContactSync				`json:"push_list"`
	SyncList		[]LocalContactSync				`json:"sync_list"`
	MergedList	[]LocalContactMergedSync	`json:"merged_list"`
	DeleteList	[]Metadata								`json:"delete_list"`
}

type RxSyncSMSFlag10 struct {
	Threads			map[int]*MozMobileMessageThread	`json:"threads"`
	Messages		map[int][]*MozSmsMessage		`json:"messages"`
}

type RxRestoreContactFlag12 struct {
	Namespace			string	`json:"namespace"`	//account:people:id
	SyncID				string	`json:"sync_id,omitempty"`
	SyncUpdated		string	`json:"sync_updated,omitempty"`
}
