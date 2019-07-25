package cmd

import (
  	"fmt"
  	"github.com/spf13/cobra"
  	"bufio"
	"os"
	"strings"
	"golang.org/x/crypto/ssh/terminal"
	"crypto/md5"
)

var LoggedInUser string

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Logs in an existing user",
	Run: func(cmd *cobra.Command, args []string) {
		username, password := credentials()
		if _, err := os.Stat(".registeredUsers"); err == nil {
			content, err := ioutil.ReadFile(".registeredUsers")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("File contents: %s", content)
		}
	},
}

func credentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(0)
	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password)
}