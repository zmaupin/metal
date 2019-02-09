package mysql

import (
	"testing"
)

const sqlCommands = `
CREATE DATABASE IF NOT EXISTS rexecd;
CREATE TABLE IF NOT EXISTS migrations (
  id int NOT NULL AUTO_INCREMENT,
  succeeded tinyint,
  PRIMARY KEY (id));
`

func TestStatementGenerator(t *testing.T) {
	count := 0
	for got := range statementGenerator(sqlCommands) {
		if count == 0 {
			expected := "CREATE DATABASE IF NOT EXISTS rexecd;"
			if expected != got {
				t.Errorf("expected %s, got %s", expected, got)
			}
		}
		if count == 1 {
			expected := `
CREATE TABLE IF NOT EXISTS migrations (
  id int NOT NULL AUTO_INCREMENT,
  succeeded tinyint,
  PRIMARY KEY (id));`
			if expected != got {
				t.Errorf("expected %s, got %s", expected, got)
			}
		}
		count++
	}
}
