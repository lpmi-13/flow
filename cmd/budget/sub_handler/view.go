package budget_subhandler

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// ViewCmd represents the view command
var ViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View the alert notification values",
	Run: func(cmd *cobra.Command, args []string) {
		category, _ := cmd.Flags().GetString("category")

		h := TakeHandler()
		table, err := h.Deps.AlertDB.ViewAlert(category)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(table[0])
	},
}

func init() {
	ViewCmd.Flags().StringP("category", "c", "", "Write the category to see the alert notification values")
}
