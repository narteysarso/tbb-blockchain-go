package main

import (
	"fmt"
	"os"

	"github.com/narteysarso/tbb-bockchain-go/database"
	"github.com/spf13/cobra"
)

func balancesCmd() *cobra.Command {
	var balancesCmd = &cobra.Command{
		Use:   "balances",
		Short: "Interact with balances (list...)",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s.%s.%s-beta %s ", Major, Minor, Fix, Verbal)
		},
	}

	balancesCmd.AddCommand(balancesListCmd)

	return balancesCmd
}

var balancesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all balances.",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := database.NewStateFromDisk()

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		defer state.Close()

		fmt.Println("Accounts balances:")
		fmt.Println("__________________")
		fmt.Println("")
		for account, balance := range state.Balances {
			fmt.Println(fmt.Sprintf("%s %d", account, balance))
		}
	},
}
