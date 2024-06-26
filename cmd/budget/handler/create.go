package budget_handler

import (
	"log"

	"github.com/ibilalkayy/flow/entities"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create the budget of different categories",
	Run: func(cmd *cobra.Command, args []string) {
		category, _ := cmd.Flags().GetString("category")
		amount, _ := cmd.Flags().GetString("amount")

		h := TakeHandler()
		amountInt := h.Deps.Common.StringToInt(amount)
		bv := entities.BudgetVariables{Category: category, Amount: amountInt}

		err := h.Deps.ManageBudget.CreateBudget(&bv)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	CreateCmd.Flags().StringP("category", "c", "", "Write the category like groceries, utilities, etc")
	CreateCmd.Flags().StringP("amount", "a", "", "Write the total amount for that category")
}
