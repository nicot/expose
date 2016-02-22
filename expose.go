package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func serve(cmd *cobra.Command, args []string) {
}

var rootCommand = &cobra.Command{
	Use:   "expose",
	Short: "Expose a port on a public IP address",
	Run:   serve,
}
