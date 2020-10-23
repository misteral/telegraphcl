package cmd

import (
	"github.com/plant99/telegraphcl/pkg/user"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Operations related to Telegraph account management",
	Run: func(cmd *cobra.Command, args []string) {
		userNew := user.User{
			ShortName:   "shortname",
			AuthorName:  "authorname",
			AuthorUrl:   "https://shivashis.xyz",
			AccessToken: "",
			AuthUrl:     "",
			PageCount:   0,
		}
		user.CreateUser(userNew)
	},
}
