package google_services

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"net/url"

	"fyne.io/fyne/v2"
	"kai-suite/types"
	"kai-suite/utils/global"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/people/v1"
	"google.golang.org/api/option"
	google_oauth2 "google.golang.org/api/oauth2/v2"
)

var(
	AuthInstance *http.Client
	TokenRepository = make(map[string]*types.UserInfoAndToken)
)

func init() {
	readTokensFromFile()
	log.Info("TokenRepository: ",len(TokenRepository))
}

func GetTokenFromWeb(config *oauth2.Config) error {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Info("Go to the following link in your browser then type the authorization code: ", authURL, "\n")
	url, _ := url.Parse(authURL)
	if err := fyne.CurrentApp().OpenURL(url); err != nil {
		return err
	}
	return nil
}

func SaveToken(config *oauth2.Config, authCode string) (*oauth2.Token, error) {
	var token *oauth2.Token
	var err error
	token, err = config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, err
	}

	var user *google_oauth2.Userinfo
	if user, err = FetchUserInfo(GetAuthClient(config, token)); err != nil {
		return nil, err
	}
	TokenRepository[user.Id] = &types.UserInfoAndToken{user, token}
	WriteTokensToFile()
	return token, nil
}

func WriteTokensToFile() {
	tokensFile := global.ResolvePath("tokens.json")
	var b []byte
	var err error
	b, err = json.Marshal(&TokenRepository)
	if err != nil {
		return
	}

	log.Info("Saving credential file to: ", tokensFile, "\n")
	f, err := os.OpenFile(tokensFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	if _, err := f.Write(b); err != nil {
		return
	}
	defer f.Close()
	global.InitDatabaseIndex(TokenRepository)
}

func readTokensFromFile() error {
	tokensFile := global.ResolvePath("tokens.json")
	b, err := os.ReadFile(tokensFile)
	if err != nil {
		log.Error(err)
		json := []byte("{}")
		if err := os.WriteFile(tokensFile, json, 0644); err != nil {
			return err
		} else {
			if _, err := os.ReadFile(tokensFile); err != nil {
				return err
			}
		}
	}
	json.Unmarshal(b, &TokenRepository)
	global.InitDatabaseIndex(TokenRepository)
	return nil
}

func GetConfig() (*oauth2.Config, error) {
	b, err := ioutil.ReadFile(global.ResolvePath("credentials.json"))
	if err != nil {
		return nil, err
	}
	return google.ConfigFromJSON(b, calendar.CalendarScope, people.ContactsScope, google_oauth2.UserinfoProfileScope, google_oauth2.UserinfoEmailScope)
}

func GetAuthClient(config *oauth2.Config, token *oauth2.Token) *http.Client {
	return config.Client(context.Background(), token)
}

func FetchUserInfo(client *http.Client) (*google_oauth2.Userinfo, error) {
	ctx := context.Background()
	srv, err := google_oauth2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Warn("Unable to create oauth2 Client: ", err)
	}
	return srv.Userinfo.V2.Me.Get().Do()
}

func RefreshToken(token *oauth2.Token) (*oauth2.Token, error) {
	config, err := GetConfig()
	if err != nil {
		log.Warn("Unable to parse client secret file to config: %v ", err)
	}
	ctx := context.Background()
	tokenSource := config.TokenSource(ctx, token)
	return tokenSource.Token()
}

func RemoveAccount(id string) {
	if _, found := TokenRepository[id]; found == true {
		tokensFile := global.ResolvePath("tokens.json")
		delete(TokenRepository, id)
		var b []byte
		b, err := json.Marshal(&TokenRepository)
		if err != nil {
			return
		}

		log.Info("Saving credential file to: ", tokensFile, "\n")
		f, err := os.OpenFile(tokensFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return
		}
		if _, err := f.Write(b); err != nil {
			return
		}
		defer f.Close()
		global.InitDatabaseIndex(TokenRepository)
	}
}
