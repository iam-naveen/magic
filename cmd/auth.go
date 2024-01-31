/*
Copyright © 2024
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/iam-naveen/magic/utils"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

var refreshMode bool

var authCodeChannel = make(chan string)

var wg = &sync.WaitGroup{}

func listenForCallback() {

	http.HandleFunc("/callback", func(res http.ResponseWriter, req *http.Request) {
		code := req.URL.Query()["code"][0]
		authCodeChannel <- code
		wg.Done()
	})

	wg.Add(1)
	go func() {
		http.ListenAndServe(":1234", nil)
	}()
}

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Command to authenticate with Google Drive",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.ReadFile("token.json")
		if err != nil || refreshMode {
			cred, credErr := os.ReadFile("credentials.json")
			if credErr != nil {
				log.Fatalf("Unable to read client secret file: %v", err)
				return
			}

			config, configErr := google.ConfigFromJSON(cred, drive.DriveScope)
			if configErr != nil {
				log.Fatalf("Unable to parse client secret file to config: %v", err)
				return
			}

			listenForCallback()

			authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
			browserErr := utils.OpenBrowser(authURL)
			if browserErr != nil {
				fmt.Printf("Open this URL in your Browser: %s", authURL)
			}
			code := <-authCodeChannel

			fmt.Println("Waiting for Authentication Code...")
			wg.Wait()
			
			fmt.Println("Authentication Code received...")
			fmt.Println("Generating Token...")

			token, tokenErr := config.Exchange(context.Background(), code)
			if tokenErr != nil {
				log.Fatalf("Unable to generate token %v", tokenErr)
				return
			}

			utils.SaveToken("token.json", token)
			fmt.Printf("Authentication Successful...")
			return
		}
		fmt.Println("Authentication is already Done...")

	},
}

func init() {
	authCmd.Flags().BoolVarP(&refreshMode, "refresh", "r", false, "Refresh the token")
	rootCmd.AddCommand(authCmd)
}
