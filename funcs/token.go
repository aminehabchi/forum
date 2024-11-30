package forum

import (
	"log"

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
	err := db.QueryRow("SELECT id FROM users WHERE token=?", uuid).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
