package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "journal",
	Short: "Personal Journal CLI App",
	Run: func(cmd *cobra.Command, args []string) {
		landing()
	},
}

func landing() {
	fmt.Println("Welcome to your personal journal! \n\nEnter from the options below: \n\n1. Login \n2. Register \n3. Exit\n")
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	choice = strings.Trim(choice, "\n")
	if choice == "1" {
		login()
	} else if choice == "2" {
		register()
	} else if choice == "3" {
		os.Exit(0)
	} else {
		fmt.Println("Invalid Choice")
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
