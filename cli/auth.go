package cli

import (
	"github.com/arelate/gog_auth"
	"github.com/arelate/gog_urls"
	"github.com/arelate/vangogh_api/cli/input"
	"github.com/boggydigital/coost"
	"net/url"
	"os"
	"path/filepath"
)

const cookiesFilename = "cookies.txt"

func AuthHandler(u *url.URL) error {
	q := u.Query()
	username := q.Get("username")
	password := q.Get("password")
	tempDir := q.Get("temp-directory")

	return Auth(username, password, tempDir)
}

func Auth(username, password, tempDir string) error {

	cookieFile, err := os.Open(filepath.Join(tempDir, cookiesFilename))
	defer cookieFile.Close()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	cj, err := coost.NewJar(cookieFile, gog_urls.GogHost)
	if err != nil {
		return err
	}

	hc := cj.NewHttpClient()

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

	return cj.Store(cookieFile)
}
