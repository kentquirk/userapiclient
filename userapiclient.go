// This is a client module to support server-side use of the Tidepool
// service called user-api.
package userapiclient

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

// UserApiClient manages the local data for a client. A client is intended to be shared among multiple
// goroutines so it's OK to treat it as a singleton (and probably a good idea).
type UserApiClient struct {
	ClientToken string // stores the most recently received client token
	ServerToken string // stores the most recently received server token

	serverName   string       // the name of the server as passed to tidepool
	apiHost      string       // the host address of the user-api
	serverSecret string       // the shared secret for this configuration
	client       *http.Client // store a reference to the http client so we can reuse it
}

// UserData is the data structure returned from a successful Login query.
type UserData struct {
	UserID   string   // the tidepool-assigned user ID
	UserName string   // the user-assigned name for the login (usually an email address)
	Emails   []string // the array of email addresses associated with this account
}

// TokenData is the data structure returned from a successful CheckToken query.
type TokenData struct {
	UserID   string // the UserID stored in the token
	IsServer bool   // true or false depending on whether the token was a servertoken
}

// NewApiClient constructs an api client object.
func NewApiClient(name, host, serversecret string) *UserApiClient {
	ac := new(UserApiClient)
	ac.serverName = name
	ac.apiHost = host
	ac.serverSecret = serversecret
	ac.client = &http.Client{}
	return ac
}

// ServerLogin issues a request to the server for a login, using the stored
// secret that was passed in on the creation of the client object. If
// successful, it stores the returned token in ServerToken.
func (ac *UserApiClient) ServerLogin() error {
	theurl, err := url.Parse(ac.apiHost)
	if err != nil {
		return err
	}
	theurl.Path += "/serverlogin"

	req, _ := http.NewRequest("POST", theurl.String(), nil)
	req.Header.Add("x-tidepool-server-name", ac.serverName)
	req.Header.Add("x-tidepool-server-secret", ac.serverSecret)

	res, err := ac.client.Do(req)
	if err != nil {
		ac.ServerToken = ""
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		ac.ServerToken = ""
		return &StatusError{NewStatus(res.StatusCode, res.Status)}
	}
	ac.ServerToken = res.Header.Get("x-tidepool-session-token")
	return nil
}

func extractUserData(r io.Reader) *UserData {
	var buf bytes.Buffer
	buf.ReadFrom(r)
	// var m map[string]interface{}
	var ud UserData
	err := json.Unmarshal(buf.Bytes(), &ud)
	log.Printf("%v, %v", err, ud)
	return &ud
}

// Login logs in a user with a username and password. Returns a UserData object if successful
// and also stores the returned login token into ClientToken.
func (ac *UserApiClient) Login(username, password string) (*UserData, error) {
	theurl, err := url.Parse(ac.apiHost)
	if err != nil {
		return nil, err
	}
	theurl.Path += "/login"

	req, _ := http.NewRequest("POST", theurl.String(), nil)
	req.SetBasicAuth(username, password)

	res, err := ac.client.Do(req)
	if err != nil {
		ac.ClientToken = ""
		return nil, err
	}
	stat := NewStatus(res.StatusCode, res.Status)
	ac.ClientToken = res.Header.Get("x-tidepool-session-token")
	defer res.Body.Close()
	if res.StatusCode != 200 {
		ac.ClientToken = ""
		return nil, &StatusError{stat}
	}
	ud := extractUserData(res.Body)
	return ud, nil
}

// CheckToken tests a token with the user-api to make sure it's current;
// if so, it returns the data encoded in the token.
func (ac *UserApiClient) CheckToken(token string) (*TokenData, error) {
	theurl, err := url.Parse(ac.apiHost)
	if err != nil {
		return nil, err
	}
	theurl.Path += "/token/" + token

	req, _ := http.NewRequest("GET", theurl.String(), nil)
	log.Printf("URL: %s", theurl.String())
	req.Header.Add("x-tidepool-session-token", ac.ServerToken)

	res, err := ac.client.Do(req)
	if err != nil {
		return nil, err
	}
	stat := NewStatus(res.StatusCode, res.Status)

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, &StatusError{stat}
	}

	// we have it so convert it
	var buf bytes.Buffer
	var td TokenData
	buf.ReadFrom(res.Body)
	if err = json.Unmarshal(buf.Bytes(), &td); err != nil {
		return nil, err
	}
	return &td, nil
}
