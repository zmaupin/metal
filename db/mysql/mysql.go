package mysql

import (
	"bufio"
	"context"
	"database/sql"
	"strings"
)

// ExecuteSQL executes each
func ExecuteSQL(ctx context.Context, db *sql.DB, s string) error {
	for statement := range statementGenerator(s) {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}

func statementGenerator(s string) chan string {
	ch := make(chan string)
	go func() {
		reader := bufio.NewReader(strings.NewReader(strings.TrimSpace(s)))
		for {
			currentString, err := reader.ReadString(';')
			if err != nil {
				close(ch)
				break
			}
			ch <- currentString
		}
	}()
	return ch
}
