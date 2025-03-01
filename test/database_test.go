package test

import (
	"context"
	"testing"

	"github.com/sebastian-j-ibanez/flourish-backend/database"
)

func Test01GetUsernames(t *testing.T) {
	connection, err := database.ConnectToDatabase()
	if err != nil {
		panic("connection exploded")
	}

	rows, err := connection.Query(context.Background(), "SELECT Username FROM Users")
	if err != nil {
		panic("query exploded")
	}

	var usernames []string
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			panic("scan exploded")
		}
		usernames = append(usernames, name)
	}

	if len(usernames) <= 0 {
		t.Errorf("Received usernames: %d", len(usernames))
	}
}
