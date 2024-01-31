/*
Copyright © 2024 Naveen <imnaveenbharath@gmail.com>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/iam-naveen/magic/utils"
	"github.com/spf13/cobra"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var (
	fileMode   bool
	folderMode bool
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the files in your Google Drive",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := utils.GetTokenFromFile("token.json")
		if err != nil {
			log.Fatalf("Authenticate using - `magic auth` command before using the cli")
		}
		config := utils.GetConfigFromFile("credentials.json")

		if token.Expiry.Before(time.Now()) {
			token = utils.RefreshToken(config, token)
			utils.SaveToken("token.json", token)
		}

		client := config.Client(context.Background(), token)

		service, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
		if err != nil {
			log.Fatalf("Unable to retrieve Drive client: %v", err)
		}
		res, err := service.Files.List().Do()
		if err != nil {
			log.Fatalf("Unable to retrieve files: %v", err)
		}

		for _, file := range res.Files {
			if fileMode && file.MimeType != "application/vnd.google-apps.folder" {
				fmt.Printf("%s\n", file.Name)
			}
			if folderMode && file.MimeType == "application/vnd.google-apps.folder" {
				fmt.Printf("%s\n", file.Name)
			}
		}
	},
}

func init() {
	showCmd.Flags().BoolVarP(&fileMode, "files", "f", true, "Show the files in your Google Drive")
	showCmd.Flags().BoolVarP(&folderMode, "folders", "F", true, "Show the folders in your Google Drive")
	showCmd.MarkFlagsMutuallyExclusive("files", "folders")
	rootCmd.AddCommand(showCmd)
}
