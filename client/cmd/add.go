package cmd

import (
	"crypto/sha256"
	"errors"
	"filestore/client/apicall"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add file.txt ...",
	Short: "Transfer local Files to remote server",
	Long:  `Save Local Files Remotely Passsing n number of files as agrument to add`,
	Run: func(cmd *cobra.Command, args []string) {
		Addfile(args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

// Addfile ... Async Reads all input files , Calls Server for Repeated contents if server returns true Calls copy API or Calls Add API to store file prints string as error or operation results
func Addfile(files []string) {
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
		hashexiststatus := apicall.Hexists(hashval)
		duplicatestatus := apicall.Fexists(k)
		if hashexiststatus && !duplicatestatus {
			status := apicall.CopyCall(k, hashval)
			if !status {
				fmt.Println("Failed to Add file ", k)
			} else {
				fmt.Println("File Copied ", k)
			}
		} else {
			if duplicatestatus {
				fmt.Println("File already Exists ", k)
			} else {
				status := apicall.AddCall(k, hashval)
				if !status {
					fmt.Println("Failed to Add file ", k)
				} else {
					fmt.Println("New File Created", k)
				}
			}
		}
	}
	mutx.Unlock()
	for k, v := range errormap {
		fmt.Println("Failed to Add file ", k, " Reason: ", v)
	}
	fmt.Println("All Files Processed")
}

// readFile... Checks File is present and able to read it returns error if file does not exists ans bytes of data (file content) if file exists
func readFile(filepath string) ([]byte, error) {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// generateHash ... Generates SHA256 Hash and returns it as string value
func generateHash(data []byte) string {
	sha256 := sha256.Sum256(data)
	strval := fmt.Sprintf("%x", sha256)
	return strval
}
