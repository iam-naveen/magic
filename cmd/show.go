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
			utils.SaveToken("token.js", token)
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
		fmt.Println(len(res.Files))

		for _, file := range res.Files {
			fmt.Printf("%s\n", file.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
