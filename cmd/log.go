package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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
		journal()
	},
}

func journal() {
	fmt.Println("Enter from the options below:\n\n 1. Add an Entry\n 2. View Journal\n 3. Logout\n 4. Exit\n")
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	choice = strings.Trim(choice, "\n")
	if choice == "1" {
		addJournal()
	} else if choice == "2" {
		viewJournal()
	} else if choice == "3" {
		LoggedInUser = ""
		log.Info("Logged Out Successfully!")
		login()
	} else if choice == "4" {
		os.Exit(0)
	} else {
		fmt.Println("Invalid Choice")
	}
}

func addJournal() {
	if isLoggedIn() {
		if _, err := os.Stat("." + LoggedInUser + "/journal"); err == nil {
			decryptedData := string(decryptFile("."+LoggedInUser+"/journal", passphrase))
			findEntries := regexp.MustCompile(">>")
			matches := findEntries.FindAllStringIndex(decryptedData, -1)
			if len(matches) >= 50 {
				log.Info("Max Entries exceed 50. Deleting oldest one to write new one.")
				i := strings.LastIndex(decryptedData, ">>")
				decryptedData = decryptedData[0:i]
			}

			entry := getNewEntry()

			newData := entry + "\n" + decryptedData
			encryptFile("."+LoggedInUser+"/journal", []byte(newData), passphrase)
			log.Info("Entry Successfully Added!")
			journal()
		} else {
			file, _ := filepath.Abs("." + LoggedInUser + "/journal")
			f, err := os.Create(file)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			entry := getNewEntry()

			f.Write(encrypt([]byte(entry+"\n"), passphrase))
			log.Info("Entry Successfully Added!")
			journal()
		}

	} else {
		fmt.Println("Login first!")
		login()
	}
}

func getNewEntry() string {
	var entries []string
	entry := ""
	fmt.Println("\nAdd your new entry. You can write in multiple lines. Press Ctrl+D on a new line when done. \n")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		entries = append(entries, scanner.Text())
	}

	if scanner.Err() != nil {
		log.Fatal("Error reading line!")
	}

	for _, val := range entries {
		entry = entry + val + "\n"
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	entry = ">> [" + timestamp + "] : " + entry
	return entry
}

func viewJournal() {
	if isLoggedIn() {
		if _, err := os.Stat("." + LoggedInUser + "/journal"); err == nil {
			decryptedData := string(decryptFile("."+LoggedInUser+"/journal", passphrase))
			fmt.Println("Your Journal Entries:\n", decryptedData, "\n")
		} else {
			log.Info("Journal Doesn't exist. Add a new entry first!\n")
		}
		journal()
	}
}

func isLoggedIn() bool {
	if LoggedInUser != "" {
		return true
	}
	return false
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
