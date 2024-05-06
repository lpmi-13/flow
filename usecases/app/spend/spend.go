package usecases_spending

import (
	"errors"
	"fmt"

	"github.com/ibilalkayy/flow/entities"
)

func (m MySpending) SpendMoney(category string, spending_amount int) error {
	values, err := m.ViewBudget(category)
	if err != nil {
		return err
	}

	details, err := ExtractBudgetValues(values)
	if err != nil {
		return err
	}

	included_categories_in_total_amount, value, err := m.TotalAmountValues()
	if err != nil {
		return err
	}

	categoryName, ok1 := details[0].(string)
	budget_category_amount, ok2 := details[1].(int)
	budget_category_spent_amount, ok3 := details[2].(int)
	budget_category_remaining_amount, ok4 := details[3].(int)

	total_amount, ok5 := value[0].(int)
	total_amount_spent, ok6 := value[1].(int)
	total_amount_status, ok7 := value[2].(string)

	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 {
		return errors.New("unable to return to int or string")
	}

	sv := entities.SpendingVariables{
		Category:                      category,
		CategoryName:                  categoryName,
		TotalAmountStatus:             total_amount_status,
		IncludedCatogeries:            included_categories_in_total_amount,
		TotalAmount:                   total_amount,
		SpendingAmount:                spending_amount,
		TotalAmountSpent:              total_amount_spent,
		BudgetCategoryAmount:          budget_category_amount,
		BudgetCategorySpentAmount:     budget_category_spent_amount,
		BudgetCategoryRemainingAmount: budget_category_remaining_amount,
	}

	fmt.Println(included_categories_in_total_amount[1][0])
	fmt.Println(included_categories_in_total_amount[0][0])

	err = m.ValidBudgetValues(&sv)
	if err != nil {
		return err
	}

	return nil
}

func ExtractBudgetValues(values [5]interface{}) ([4]interface{}, error) {
	categoryName, ok1 := values[1].(string)
	budget_category_amount, ok2 := values[2].(int)
	budget_category_spent_amount, ok3 := values[3].(int)
	budget_category_remaining_amount, ok4 := values[4].(int)

	if !ok1 || !ok2 || !ok3 || !ok4 {
		return [4]interface{}{}, errors.New("unable to convert budget amount to int or string")
	}

	details := [4]interface{}{categoryName, budget_category_amount, budget_category_spent_amount, budget_category_remaining_amount}
	return details, nil
}

func (m MySpending) ValidBudgetValues(sv *entities.SpendingVariables) error {
	if sv.TotalAmountStatus != "active" {
		return errors.New("make your total amount status active. see 'flow total-amount -h'")
	}

	if sv.TotalAmountSpent+sv.SpendingAmount > sv.TotalAmount {
		return errors.New("you have exceeded the total amount")
	}

	foundCategory := false
	for i := 0; i < len(sv.IncludedCatogeries); i++ {
		if sv.Category == sv.CategoryName && sv.Category == sv.IncludedCatogeries[i][0] {
			foundCategory = true
			budget_category_total_spending_amount := sv.SpendingAmount + sv.BudgetCategorySpentAmount

			if budget_category_total_spending_amount <= sv.BudgetCategoryAmount {
				err := m.UpdateBudgetAndTotalAmount(sv)
				if err != nil {
					return err
				}
				break
			} else if sv.SpendingAmount <= sv.BudgetCategoryRemainingAmount {
				err := m.UpdateBudgetAndTotalAmount(sv)
				if err != nil {
					return err
				}
				break
			} else if sv.SpendingAmount > sv.BudgetCategoryRemainingAmount && sv.SpendingAmount <= sv.TotalAmount && sv.BudgetCategorySpentAmount <= sv.TotalAmount && budget_category_total_spending_amount <= sv.TotalAmount {
				err := m.HandleExceededBudget(sv)
				if err != nil {
					return err
				}
				break
			} else {
				return errors.New("you have exceeded the total amount")
			}
		}
	}

	if !foundCategory {
		return errors.New("category is not found. setup the alert 'flow budget alert setup -h' or include the category in your total amount 'flow total-amount set -h'")
	}
	return nil
}

func (m MySpending) UpdateBudgetAndTotalAmount(sv *entities.SpendingVariables) error {
	err := m.AddBudgetExpenditure(sv.SpendingAmount, sv.Category)
	if err != nil {
		return err
	}
	err = m.CalculateRemaining(sv.Category)
	if err != nil {
		return err
	}

	err = m.StoreHistory(sv.Category, sv.SpendingAmount)
	if err != nil {
		return err
	}

	fmt.Println("Enjoy your spending!")
	return nil
}

func (m MySpending) HandleExceededBudget(sv *entities.SpendingVariables) error {
	var answer string
	fmt.Printf("You have spent %d and your remaining balance is %d but your budget is %d\n", sv.BudgetCategorySpentAmount, sv.BudgetCategoryRemainingAmount, sv.BudgetCategoryAmount)
	fmt.Printf("Do you still want to spend? [yes/no]: ")
	fmt.Scanln(&answer)

	switch answer {
	case "yes", "y":
		m.SendAlertEmail(sv.Category)
		err := m.UpdateBudgetAndTotalAmount(sv)
		if err != nil {
			return err
		}
	case "no", "n":
		fmt.Println("Alright")
	default:
		return errors.New("select the right option")
	}
	return nil
}
