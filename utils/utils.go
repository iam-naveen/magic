package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
)

var (
	authCodeChan = make(chan string)
)

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["code"]
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'code' is missing")
		return
	}
	authCode := keys[0]
	authCodeChan <- authCode
	fmt.Fprintf(w, "Authorization successful.\n\nYou can close this tab now. Go back to the terminal to continue.")
}

func OpenBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported platform")
	}
	return cmd.Run()
}

func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-Token", oauth2.AccessTypeOffline)

	if err := OpenBrowser(authURL); err != nil {
		fmt.Printf("Go to the following link in your browser then type the "+
			"authorization code: \n%v\n", authURL)
	}

	authCode := <-authCodeChan
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

func GetTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func RefreshToken(config *oauth2.Config, token *oauth2.Token) (*oauth2.Token, error) {
	newSource := config.TokenSource(context.Background(), token)
	newToken, err := newSource.Token()
	if err != nil {
		return nil, err
	}
	return newToken, nil
}

func SaveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential...")
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

