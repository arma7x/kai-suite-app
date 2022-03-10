package types

import(
	"golang.org/x/oauth2"
	google_oauth2 "google.golang.org/api/oauth2/v2"
)

type UserInfoAndToken struct {
	User	*google_oauth2.Userinfo	`json:"user"`
	Token	*oauth2.Token 					`json:"token"`
}

type Metadata struct {
	SyncID 				string	`json:"sync_id,omitempty"`
	SyncUpdated		string	`json:"sync_updated,omitempty"`
	Hash					string	`json:"hash,omitempty"`
	Deleted 			bool		`json:"deleted,omitempty"`
}
