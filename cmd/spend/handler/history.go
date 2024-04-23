package spend_handler

import (
	"log"

	"github.com/ibilalkayy/flow/db/budget_db"
	"github.com/spf13/cobra"
)

// HistoryCmd represents the history command
var HistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Show the transaction history",
	Run: func(cmd *cobra.Command, args []string) {
		category, _ := cmd.Flags().GetString("category")
		err := budget_db.ViewHistory(category)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	HistoryCmd.Flags().StringP("category", "c", "", "Write the category to show it's history")
}
