package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var SetCmd = &cobra.Command{
	Use:   "set",
	Short: "Re-write default config parameters",
	Long:  `Re-write default config parameters`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use `meteo set coords` to re-write default corrdinates")
	},
}

func init() {
	rootCmd.AddCommand(SetCmd)

}
