package total_amount_db

import (
	"database/sql"
	"errors"

	"github.com/ibilalkayy/flow/entities"
	"github.com/ibilalkayy/flow/handler"
	"github.com/jedib0t/go-pretty/v6/table"
)

type MyTotalAmountDB struct {
	*handler.Handler
}

func (h MyTotalAmountDB) InsertTotalAmountCategory(tv *entities.TotalAmountVariables) error {
	data, err := h.Deps.Connect.Table("framework/db/migrations/003_create_total_amount_table.sql", 1)
	if err != nil {
		return err
	}

	query := "INSERT INTO TotalAmountCategories(included_categories, labels) VALUES($1, $2)"
	insert, err := data.Prepare(query)
	if err != nil {
		return err
	}

	defer insert.Close()

	if len(tv.Included) != 0 {
		_, err = insert.Exec(tv.Included, tv.Label)
		if err != nil {
			return err
		}
	} else {
		return errors.New("write total amount and category. see 'flow total-amount set -h'")
	}
	return nil
}

func (h MyTotalAmountDB) ViewTotalAmountCategories() (string, [][2]string, error) {
	tv := new(entities.TotalAmountVariables)
	var values [][2]string

	db, err := h.Deps.Connect.Connection()
	if err != nil {
		return "", [][2]string{}, err
	}
	defer db.Close()

	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"Included Categories", "Labels"})

	var rows *sql.Rows
	query := "SELECT included_categories, labels FROM TotalAmountCategories"
	rows, err = db.Query(query)
	if err != nil {
		return "", [][2]string{}, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&tv.Included, &tv.Label); err != nil {
			return "", [][2]string{}, err
		}
		values = append(values, [2]string{tv.Included, tv.Label})
	}

	if err := rows.Err(); err != nil {
		return "", [][2]string{}, err
	}

	for i := 0; i < len(values); i++ {
		tw.AppendRow([]interface{}{values[i][0], values[i][1]})
	}
	tableRender := "Categories\n" + tw.Render()

	return tableRender, values, nil
}
