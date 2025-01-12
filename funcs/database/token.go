package forum

import (
	"fmt"
	types "forum/funcs/types"
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
	err := Db.QueryRow("SELECT user_id,created_at FROM tokens WHERE token=?", uuid).Scan(&id, &created_time)
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
func GetUserInfoByLoginInfo(identifier string) (*types.User, error) {
	query := `SELECT id, password FROM users WHERE email = ? OR uname = ?`
	user := &types.User{}
	err := Db.QueryRow(query, identifier, identifier).Scan(&user.ID, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func SetToken(token string, id int) error {
	query := "UPDATE tokens SET token=?,created_at=CURRENT_TIMESTAMP WHERE user_id=?"
	_, err := Db.Exec(query, token, id)
	if err != nil {
		return err
	}
	return nil
}
