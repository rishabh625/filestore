package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Transfer local Files to remote server",
	Long:  `Save Local Files Remotely Passsing n number of files as agrument to add`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called", args[0])
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
