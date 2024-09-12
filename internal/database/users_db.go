package database

type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email string, password string) (User, error) {
	_, err := db.GetUserByEmail(email)
	if err == nil {
		return User{}, ErrDuplicate
	}

	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	id := len(dbs.Users) + 1
	user := User{
		Id:          id,
		Email:       email,
		Password:    password,
		IsChirpyRed: false,
	}

	dbs.setUser(id, user)

	err = db.writeDB(dbs)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// UpgradeUser - Can updrade to Chirpy Red
func (db *DB) UpgradeUser(id int, isChirpyRed bool) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	user, ok := dbs.getUser(id)
	if !ok {
		return User{}, ErrNotExist
	}
	user.IsChirpyRed = isChirpyRed

	dbs.setUser(id, user)

	err = db.writeDB(dbs)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// UpdateUser - Update email and/or password
func (db *DB) UpdateUser(id int, email string, password string) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	user, ok := dbs.getUser(id)
	if !ok {
		return User{}, ErrNotExist
	}

	if email != user.Email {
		_, err = db.GetUserByEmail(email)
		if err == nil {
			return User{}, ErrDuplicate
		}
	}

	user.Email = email
	user.Password = password

	dbs.setUser(id, user)

	err = db.writeDB(dbs)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// GetUsers -  returns all chirps in the database
func (db *DB) GetUsers() ([]User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbs.Users))
	for _, user := range dbs.Users {
		users = append(users, user)
	}
	return users, nil
}

// GetUserById - Get one user
func (db *DB) GetUserById(id int) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbs.getUser(id)

	if ok {
		return user, nil
	}
	return User{}, ErrNotExist
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbs.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrNotExist
}
