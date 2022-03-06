package google_services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"net/url"

	"fyne.io/fyne/v2"
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
)

func init() {
	_, err_cre := ioutil.ReadFile(global.ResolvePath("credentials.json"))
	_, err_tok := ioutil.ReadFile(global.ResolvePath("token.json"))
	if err_cre == nil && err_tok == nil {
		AuthInstance = GetAuth()
	}
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := global.ResolvePath("token.json")
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		return nil
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func GetTokenFromWeb(config *oauth2.Config) error {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)
	url, _ := url.Parse(authURL)
	if err := fyne.CurrentApp().OpenURL(url); err != nil {
		return err
	}
	return nil
}

// Saves a token to a file path.
func SaveToken(config *oauth2.Config, path string, authCode string) (*oauth2.Token, error) {
	var token *oauth2.Token
	var err error
	token, err = config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return token, nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func GetConfig() (*oauth2.Config, error) {
	b, err := ioutil.ReadFile(global.ResolvePath("credentials.json"))
	if err != nil {
		return nil, err
	}
	return google.ConfigFromJSON(b, calendar.CalendarScope, people.ContactsScope, google_oauth2.UserinfoProfileScope)
}

func GetAuth() *http.Client {
	config, err := GetConfig()
	if err != nil {
		log.Warn("Unable to parse client secret file to config: %v ", err)
	}
	return getClient(config)
}

func FetchUserInfo(client *http.Client) (*google_oauth2.Userinfo, error) {
	config, _ := GetConfig()
	restoredToken, _ := tokenFromFile(global.ResolvePath("token.json"))
	tokenSource := config.TokenSource(oauth2.NoContext, restoredToken)
	newclient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	ctx := context.Background()
	srv, err := google_oauth2.NewService(ctx, option.WithHTTPClient(newclient))
	if err != nil {
		log.Warn("Unable to create oauth2 Client: ", err)
	}
	return srv.Userinfo.V2.Me.Get().Do()
}
