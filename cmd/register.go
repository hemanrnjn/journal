package cmd

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
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
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.Trim(username, "\n")
	passphrase := "ambitionjournal123"
	if _, err := os.Stat(".registeredUsers"); err == nil {
		decryptedData := string(decryptFile(".registeredUsers", passphrase))

		if strings.Contains(decryptedData, username) {
			log.Fatal("User with this username already exists")
		}

		fmt.Print("Enter Password: ")
		bytePassword, _ := terminal.ReadPassword(0)
		password := string(bytePassword)

		newData := decryptedData + username + ":" + password + "\n"
		fmt.Println(newData)

		encryptFile(".registeredUsers", []byte(newData), passphrase)
		if _, err := os.Stat("." + username); os.IsNotExist(err) {
			os.Mkdir("."+username, os.ModeDir)
		}
		log.Info("User Registered")

	} else {
		fmt.Print("Enter Password: ")
		bytePassword, _ := terminal.ReadPassword(0)
		password := string(bytePassword)
		fmt.Print("\n")
		newUser := username + ":" + password + "\n"
		encryptFile(".registeredUsers", []byte(newUser), passphrase)
		if _, err := os.Stat("." + username); os.IsNotExist(err) {
			os.Mkdir("."+username, os.ModeDir)
		}
		log.Info("User Registered")
	}
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func encryptFile(filename string, data []byte, passphrase string) {
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(encrypt(data, passphrase))
}

func decryptFile(filename string, passphrase string) []byte {
	data, _ := ioutil.ReadFile(filename)
	return decrypt(data, passphrase)
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
