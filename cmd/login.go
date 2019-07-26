package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var LoggedInUser string

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Logs in an existing user",
	Run: func(cmd *cobra.Command, args []string) {
		login()
	},
}

func login() {
	if _, err := os.Stat(".registeredUsers"); err == nil {
		passphrase := "ambitionjournal123"
		decryptedData := string(decryptFile(".registeredUsers", passphrase))
		scanner := bufio.NewScanner(strings.NewReader(decryptedData))
	outer:
		for scanner.Scan() {
			text := scanner.Text()
			for {
				username, password := credentials()
				if strings.Index(text, username) != -1 {
					if strings.Split(text, ":")[1] == password {
						LoggedInUser = username
						log.Info("Login Successful! Welcome ", LoggedInUser)
						break outer
					} else {
						log.Info("Incorrect Password, Try again")
					}
				}
			}
		}
		journal()
	} else {
		log.Info("User does not exist. Register First!")
		register()
	}
}

func credentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, _ := terminal.ReadPassword(0)
	fmt.Print("\n")
	password := string(bytePassword)

	return strings.Trim(username, " \n"), strings.TrimSpace(password)
}
