package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "journal",
	Short: "Personal Journal CLI App",
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("Command works!")
	// },
}
  
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}