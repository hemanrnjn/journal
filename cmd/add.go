package cmd

import (
  "github.com/spf13/cobra"
)

// var Username string
// var Password string

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new entry to the journal",
	Run: func(cmd *cobra.Command, args []string) {
		if isLoggedIn() {

		}
	},
}

func isLoggedIn() bool {
	if _, err := os.Stat(".loggedUser"); err == nil {
		content, err := ioutil.ReadFile("testdata/hello")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("File contents: %s", content)
	}
	return false
}