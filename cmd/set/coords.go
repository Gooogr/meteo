package set

import (
	"fmt"
	cmd "meteo/cmd"

	"github.com/spf13/cobra"
)

// coordsCmd represents the coords command
var coordsCmd = &cobra.Command{
	Use:   "coords",
	Short: "Re-write default latitude and longitude",
	Long:  `Re-write default latitude and longitude`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("coords called")
	},
}

func init() {
	cmd.SetCmd.AddCommand(coordsCmd)
}
