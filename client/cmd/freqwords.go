package cmd

import (
	"encoding/json"
	"filestore/client/apicall"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type Wc struct {
	Word  string `json:"Word"`
	Count int    `json:"Count"`
}

// freqCmd represents the add command
var freqCmd = &cobra.Command{
	Use:   "freqwords [--limit n] [--order asc|desc]",
	Short: "Gives n number of words in order asc|desc",
	Long:  `Lists all the file in filestore and returns most frequently used words based on options given.`,
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetString("limit")
		order, _ := cmd.Flags().GetString("order")
		m := GetFreqWords(limit, order)
		for _, v := range m {
			fmt.Printf("Word : Count \n ############## : ##############\n")
			for _, val := range v {
				fmt.Printf("%s : %d\n", val.Word, val.Count)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(freqCmd)
	freqCmd.PersistentFlags().String("limit", "10", "Gives n number of frequently used words")
	freqCmd.PersistentFlags().String("order", "desc", "Gives n number of frequently used words in given order asc or desc")
}

// ValidateFlags ... Checks limit is number and order is asc or desc if not returns default value 10,asc
func ValidateFlags(x, y string) (limit, order string) {
	if strings.EqualFold(y, "asc") {
		order = "asc"
	} else {
		order = "desc"
	}
	_, err := strconv.Atoi(x)
	if err != nil {
		limit = "10"
	}
	return
}

// GetFreqWords ... Does API call to get Frequent words in store
func GetFreqWords(limit, order string) map[string][]Wc {
	bytesd := apicall.FreqWords(limit, order)
	if bytesd == nil {
		return nil
	}
	obj := make(map[string][]Wc)
	err := json.Unmarshal(bytesd, &obj)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return obj
}
