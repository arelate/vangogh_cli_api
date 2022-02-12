package cli

import (
	"bufio"
	"fmt"
	"github.com/arelate/gog_atu"
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

	cj, err := coost.NewJar(cookieFile, gog_atu.GogHost)
	if err != nil {
		return err
	}

	hc := cj.NewHttpClient()

	li, err := gog_atu.LoggedIn(hc)
	if err != nil {
		return err
	}

	if li {
		return nil
	}

	if err := gog_atu.Login(hc, username, password, requestText); err != nil {
		return err
	}

	return cj.Store(cookieFile)
}

func requestText(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		return scanner.Text()
	}
	return ""
}
