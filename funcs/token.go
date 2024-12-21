package forum

import (
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
)

func GenereteTocken() (string, error) {
	// Create a Version 4 UUID.
	u2, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
		return "", err
	}
	return u2.String(), nil
}

func GetUserIDFromToken(uuid string) (int, error) {
	var id int
	var created_time string
	err := db.QueryRow("SELECT user_id,created_at FROM tokens WHERE token=?", uuid).Scan(&id, &created_time)
	if err != nil {
		return 0, err
	}
	t, err := time.Parse(time.RFC3339, created_time)
	if err != nil {
		return 0, err
	}
	t = t.Add(1 * time.Hour)
	if time.Now().After(t) {
		err = fmt.Errorf("expired token")
		return 0, err
	}
	return id, nil
}
