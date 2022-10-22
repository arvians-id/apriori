package setup

import (
	"database/sql"
)

func TearDownTest(db *sql.DB) error {
	_, err := db.Exec(`DELETE FROM apriories;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM products;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM transactions;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM payloads;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM user_orders;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM users;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM password_resets;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM categories;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM comments;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM notifications;`)
	if err != nil {
		return err
	}

	return nil
}
