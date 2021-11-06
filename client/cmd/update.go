package cmd

import (
	"filestore/client/apicall"
	"fmt"
	"sync"

	"github.com/spf13/cobra"
)

// updateCmd represents the add command
var updateCmd = &cobra.Command{
	Use:   "update file.txt ...",
	Short: "update contents of local file to remote",
	Long: `update contents of file.txt(passed in arguments) in
	server with the local file.txt(passed in arguments) or create a new file.txt(passed in arguments) in server if it is
	absent.`,
	Run: func(cmd *cobra.Command, args []string) {
		Updatefile(args)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

// Addfile ... Async Reads all input files , Calls Server for Repeated contents if server returns true calls copy API or calls Add API to store file prints string as error or operation results
func Updatefile(files []string) {
	wg := &sync.WaitGroup{}
	datamap := make(map[string]([]byte))
	errormap := make(map[string]error)
	mutx := sync.Mutex{}
	for _, v := range files {
		wg.Add(1)
		go func(local string) {
			content, err := readFile(local)
			if err != nil {
				errormap[v] = err
			} else {
				mutx.Lock()
				datamap[local] = content
				mutx.Unlock()
			}
			wg.Done()
		}(v)
	}
	wg.Wait()
	mutx.Lock()
	for k, v := range datamap {
		hashval := generateHash(v)
		hashexiststatus, _ := apicall.Hexists(hashval, k)
		if hashexiststatus {
			status := apicall.CopyCall(k, hashval)
			if !status {
				fmt.Println("Failed to Update file ", k)
			} else {
				fmt.Println("File Copied ", k)
			}
		} else {
			status := apicall.AddCall(k, hashval)
			if !status {
				fmt.Println("Failed to Update file ", k)
			} else {
				fmt.Println("New File Created as File was not present on server", k)
			}
		}
	}
	mutx.Unlock()
	for k, v := range errormap {
		fmt.Println("Failed to Add file ", k, " Reason: ", v)
	}
	fmt.Println("All Files Processed")
}
