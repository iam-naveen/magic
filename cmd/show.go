/*
Copyright © 2024 Naveen <imnaveenbharath@gmail.com>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/iam-naveen/magic/utils"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the files in your Google Drive",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfig("credentials.json")
		client := getClient(config)

		service, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
		if err != nil {
			log.Fatalf("Unable to retrieve Drive client: %v", err)
		}

		result, err := service.Files.List().PageSize(50).Fields("nextPageToken, files(id, name)").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve files: %v", err)
		}

		for _, i := range result.Files {
			fmt.Printf("%s\n", i.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}

/*
Read the Credentials file and
parse the JSON to oauth2.Config
*/
func getConfig(credentialsPath string) *oauth2.Config {
	jsonKey, err := os.ReadFile(credentialsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(jsonKey, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return config
}

/*
Checks if the token is present already,
else initiates the oauth2 flow to get the token
and save it.
*/
func getClient(config *oauth2.Config) *http.Client {
	tokenPath := "token.json"
	token, err := utils.GetTokenFromFile(tokenPath)
	if err != nil {
		startCallbackServer()
		token = utils.GetTokenFromWeb(config)
		utils.SaveToken(tokenPath, token)
	} else if token.Expiry.Before(time.Now()) {
		token, err = utils.RefreshToken(config, token)
		if err != nil {
			log.Fatalf("Unable to refresh token: %v", err)
		}
		utils.SaveToken(tokenPath, token)
	}
	return config.Client(context.Background(), token)
}

/*
Starts a server to listen for the oauth2 callback
and writes the response to a channel.
*/
func startCallbackServer() {
	http.HandleFunc("/callback", utils.HandleAuthCallback)
	go func() {
		if err := http.ListenAndServe(":1234", nil); err != nil {
			log.Fatal(err)
		}
	}()
}
