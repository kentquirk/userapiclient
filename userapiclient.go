// This is a client module to support server-side use of the Tidepool
// service called user-api.
package userapiclient

import (
	"io"
	"net/http"
	"net/url"
)

type UserApiClient struct {
	serverName   string
	apiHost      string
	serverSecret string
	clientToken  string
	serverToken  string
	client       *http.Client
}

type UserData struct {
	username  string
	userEmail string
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
// secret that was passed in on the creation of the client object. It returns
// true if the server was successfully logged in.
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
		ac.serverToken = ""
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		ac.serverToken = ""
		return &StatusError{NewStatus(res.StatusCode, res.Status)}
	}
	ac.serverToken = res.Header.Get("x-tidepool-session-token")
	return nil
}

func (ac *UserApiClient) Login(username, password string, w io.Writer) (int64, error) {
	theurl, err := url.Parse(ac.apiHost)
	if err != nil {
		return 0, err
	}
	theurl.Path += "/login"

	req, _ := http.NewRequest("POST", theurl.String(), nil)
	req.SetBasicAuth(username, password)

	res, err := ac.client.Do(req)
	if err != nil {
		ac.serverToken = ""
		return 0, err
	}
	stat := NewStatus(res.StatusCode, res.Status)

	defer res.Body.Close()
	if res.StatusCode != 200 {
		ac.serverToken = ""
		return 0, &StatusError{stat}
	}
	return io.Copy(w, res.Body)
}

func (ac *UserApiClient) CheckToken(token string, w io.Writer) (int64, error) {
	theurl, err := url.Parse(ac.apiHost)
	if err != nil {
		return 0, err
	}
	theurl.Path += "/token/" + token

	req, _ := http.NewRequest("GET", theurl.String(), nil)
	req.Header.Add("x-tidepool-session-token", ac.serverToken)

	res, err := ac.client.Do(req)
	if err != nil {
		ac.serverToken = ""
		return 0, err
	}
	stat := NewStatus(res.StatusCode, res.Status)

	defer res.Body.Close()
	if res.StatusCode != 200 {
		ac.serverToken = ""
		return 0, &StatusError{stat}
	}
	return io.Copy(w, res.Body)
}
