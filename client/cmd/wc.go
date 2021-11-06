package cmd

import (
	"filestore/client/apicall"
	"fmt"

	"github.com/spf13/cobra"
)

// wcCmd represents the add command
var wcCmd = &cobra.Command{
	Use:   "wc",
	Short: "Retrievs word count",
	Long:  `Gives Word Count of files present in filestore`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(apicall.WC())
	},
}

func init() {
	rootCmd.AddCommand(wcCmd)
}
