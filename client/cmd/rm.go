package cmd

import (
	"filestore/client/apicall"
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Delete Files from remote server",
	Long:  `Delete Files from remote server Passsing n number of files as agruments`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Remove(args))
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}

// Remove Calls Remote Server to delete files
func Remove(file []string) string {
	var finalstr string
	for _, f := range file {
		str := apicall.Remove(f)
		finalstr = fmt.Sprintf("%s\n%s", finalstr, str)
	}
	return finalstr
}
