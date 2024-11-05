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

func GetUserNameFromToken(uuid string) (string, error) {
	uname := ""
	err := db.QueryRow("SELECT uname FROM users WHERE token=?", uuid).Scan(&uname)
	if err != nil {
		return "", err
	}
	return uname, nil
}
