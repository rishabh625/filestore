package cmd

import (
	"filestore/client/apicall"
	"fmt"

	"github.com/spf13/cobra"
)

// lsCmd represents the add command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List Files on remote server",
	Long:  `List Files Saved Remotely `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(apicall.List())
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
