package cli

import (
	"github.com/arelate/gog_auth"
	"github.com/arelate/gog_urls"
	"github.com/arelate/vangogh_api/cli/input"
	"github.com/boggydigital/cooja"
	"net/url"
)

const tempDirectory = "/var/tmp/vangogh"

var gogHosts = []string{gog_urls.GogHost}

func AuthHandler(u *url.URL) error {
	q := u.Query()
	username := q.Get("username")
	password := q.Get("password")

	return Auth(username, password)
}

func Auth(username, password string) error {

	cj, err := cooja.NewJar(gogHosts, tempDirectory)
	if err != nil {
		return err
	}

	hc := cj.GetClient()

	li, err := gog_auth.LoggedIn(hc)
	if err != nil {
		return err
	}

	if li {
		return nil
	}

	if err := gog_auth.Login(hc, username, password, input.RequestText); err != nil {
		return err
	}

	return cj.Save()
}
