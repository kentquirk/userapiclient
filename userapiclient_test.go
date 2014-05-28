package userapiclient

import (
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

func TestCheckToken(t *testing.T) {
	// t.Skip("Skipping TestCheckToken")
	ac := createClient()
	a, err := ac.CheckToken("123", nil) // this one should fail
	if err == nil {
		t.Log(a)
		t.Fail()
	}
	t.Log(err)
}

func TestLogin(t *testing.T) {
	// t.Skip("Skipping TestLogin")
	ac := createClient()
	a, err := ac.Login("b@b.com", "bbb", nil)
	if err == nil {
		t.Log(a)
		t.Fail()
	}
	t.Log(err)
}
