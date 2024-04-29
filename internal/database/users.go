package database

func (db *DB) CreateUser(email string) (User, error) {

	dbContent, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	newUser := User{}

	for i := 1; ; i++ {
		_, ok := dbContent.Users[i]
		if !ok {
			newUser = User{
				Id:    i,
				Email: email,
			}
			dbContent.Users[i] = newUser
			break
		}
	}

	err = db.writeDB(dbContent)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}
