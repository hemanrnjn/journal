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

func init() {
	rootCmd.AddCommand(registerCmd)
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Registers a new user",
	Run: func(cmd *cobra.Command, args []string) {
		register()
	},
}

func register() {
	if _, err := os.Stat(".registeredUsers"); err == nil {
		decryptedData := string(decryptFile(".registeredUsers", passphrase))
		fmt.Println("Number of new lines: ", strings.Count(decryptedData, "\n"))
		if strings.Count(decryptedData, "\n") >= 10 {
			log.Fatal("Max users limit [10 users] reached!")
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Username: ")
		username, _ := reader.ReadString('\n')
		username = strings.Trim(username, "\n")

		if strings.Contains(decryptedData, username) {
			log.Info("User with this username already exists")
		}

		fmt.Print("Enter Password: ")
		bytePassword, _ := terminal.ReadPassword(0)
		password := string(bytePassword)
		fmt.Println("")

		newData := decryptedData + username + ":" + password + "\n"

		encryptFile(".registeredUsers", []byte(newData), passphrase)
		if _, err := os.Stat("." + username); os.IsNotExist(err) {
			os.Mkdir("."+username, os.ModePerm)
		}
		log.Info("User Registered")
		LoggedInUser = username
		journal()

	} else {
		username, password := credentials()
		fmt.Println("")
		username = strings.Trim(username, "\n")
		newUser := username + ":" + password + "\n"
		encryptFile(".registeredUsers", []byte(newUser), passphrase)
		if _, err := os.Stat("." + username); os.IsNotExist(err) {
			os.Mkdir("."+username, 0777)
		}
		log.Info("User Registered")
		LoggedInUser = username
		journal()
	}
}
