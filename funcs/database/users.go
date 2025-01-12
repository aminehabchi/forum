package forum

func InsertUserInfo(email, password, uname string) error {
	selector := `INSERT INTO users(password,uname,email) VALUES (?,?,?)`
	result, err := Db.Exec(selector, password, uname, email)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	selector = `INSERT INTO tokens(user_id) VALUES (?)`
	_, err = Db.Exec(selector, int(id))
	if err != nil {
		return err
	}
	return nil
}
