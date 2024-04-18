package total_amount_subhandler

import (
	"fmt"
	"log"

	"github.com/ibilalkayy/flow/db/total_amount_db"
	"github.com/spf13/cobra"
)

// CheckCmd represents the check command
var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		values, err := total_amount_db.ViewTotalAmount()
		if err != nil {
			log.Fatal(err)
		}

		status, ok := values[2].(string)
		if !ok {
			fmt.Println("unable to convert to string")
		}

		fmt.Printf("This is your total amount status: %s\n", status)
	},
}