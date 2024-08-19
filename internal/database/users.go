package database

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) CreateUser(email string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStruct.Users) + 1
	user := User{id, email}

	dbStruct.Users[id] = user

	err = db.writeDB(dbStruct)
	if err != nil {
		return user, err
	}

	return user, nil
}
