package cmd

import (
  "github.com/spf13/cobra"
  "github.com/joho/godotenv"
  "strings"
  "log"
  "os"
)

func init() {
	rootCmd.AddCommand(registerCmd)
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Registers a new user",
	Run: func(cmd *cobra.Command, args []string) {
		godotenv.Load()
		passphrase := os.Getenv("Passphrase")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Username: ")
		username, _ := reader.ReadString('\n')
		if _, err := os.Stat(".registeredUsers"); err == nil {
			data, _ := ioutil.ReadFile(".registeredUsers")
			decryptedData = decrypt(data, passphrase)
		
			if decryptedData.Contains(scanner.Text(), username) {
				log.Fatal("User with this username already exists")
			}

			newData := decryptedData + "\n" + username + ":" + password + "\n"

			if _, err := file.Write(encrypt(newData, passphrase)); err != nil {
				log.Fatal(err)
			}

		} else {
			newUser := username + ":" + password + "\n"
			f, _ := os.Create(".registeredUsers")
			defer f.Close()
			f.Write(encrypt(newUser, password))
		}
	},
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

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}