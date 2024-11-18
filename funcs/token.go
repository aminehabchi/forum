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

func GetUserNameFromToken(uuid string) (int, error) {
	id := 0
	err := db.QueryRow("SELECT id FROM users WHERE token=?", uuid).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}
