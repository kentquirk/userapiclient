package userapiclient

import (
	"bytes"
	"testing"
)

var ac *UserApiClient

func createClient() *UserApiClient {
	if ac == nil {
		ac = NewApiClient("gotest", "http://localhost:9107",
			"This needs to be the same secret everywhere. YaHut75NsK1f9UKUXuWqxNN0RUwHFBCy")
	}
	return ac
}

func TestServerLogin(t *testing.T) {
	// t.Skip("Skipping TestServerLogin")

	// ac.Init("http://localhost:9107",
	// "This needs to be the same secret everywhere. YaHut75NsK1f9UKUXuWqxNN0RUwHFBCy")
	ac := createClient()
	if err := ac.ServerLogin(); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestCheckTokenFail(t *testing.T) {
	// t.Skip("Skipping TestCheckToken")
	ac := createClient()
	a, err := ac.CheckToken("123") // this one should fail
	if err == nil {
		t.Log(a)
		t.Fail()
	}
	t.Log(err)
}

func TestLoginFail(t *testing.T) {
	// t.Skip("Skipping TestLogin")
	var w bytes.Buffer
	ac := createClient()
	ud, err := ac.Login("b@a.com", "ababb")
	if err == nil {
		t.Log(ud)
		t.Fail()
	}
	t.Log(w.String())
	t.Log(err)
}

func TestLoginSucceed(t *testing.T) {
	// t.Skip("Skipping TestLoginSucceed")
	ac := createClient()
	ud, err := ac.Login("a@a.com", "aaa")
	if err != nil {
		t.Fail()
	}
	t.Log(ud)
	t.Log(err)
}

func TestCheckTokenSucceed(t *testing.T) {
	// t.Skip("Skipping TestCheckToken")
	ac := createClient()
	td, err := ac.CheckToken(ac.ClientToken)
	if err != nil {
		t.Fail()
	}
	t.Log(td)
	t.Log(err)
}
