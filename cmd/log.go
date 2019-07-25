package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var entry []string
var Add string
var View string

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.Flags().StringVarP(&Add, "add", "a", "", "Add a new entry to the journal")
	logCmd.Flags().StringVarP(&View, "view", "v", "", "View all your entries")
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Adds a new entry to the journal or view your entries",
	Run: func(cmd *cobra.Command, args []string) {
		if isLoggedIn() {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				entry = append(entry, scanner.Text())
			}

			if scanner.Err() != nil {
				log.Fatal("Error reading line!")
			}

			for _, val := range entry {
				fmt.Println(val)
			}
		} else {
			fmt.Println("Login first!")
			login()
		}
	},
}

func isLoggedIn() bool {
	if LoggedInUser != "" {
		return true
	}
	return false
}
