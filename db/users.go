package database

import (
	"encoding/json"
	"errors"
)

// User defines the user model.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// All returns all records, sorted in ascending order for a given resource.
// An error is returned if the read process fails.
func (db *User) All() ([]User, error) {
	users, err := db.fetchUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

// FindByID locates the correct user based on ID.
func (db *User) FindByID(id int) (*User, error) {
	users, err := db.fetchUsers()
	if err != nil {
		return nil, err
	}
	for _, t := range users {
		if t.ID == id {
			return &t, nil
		}
	}
	return nil, errors.New("user not found")
}

// Create takes in a user object and adds the item to the db.
// An error is returned if the request body is empty,
// or if any error occurs when reading/writing data.
func (db *User) Create(user *User) error {
	users, err := db.fetchUsers()
	if err != nil {
		return err
	}

	user.ID = users[len(users)-1].ID + 1

	users = append(users, *user)
	err = Save(userData, users)

	return err
}

// Update takes a user ID and user object and handles updating the db.
// An error is returned if the process fails.
func (db *User) Update(id int, user *User) error {
	users, err := db.All()
	if err != nil {
		return err
	}
	for i, v := range users {

		if v.ID == id {
			user.ID = id
			users[i] = *user
			break
		}
	}
	// fmt.Printf("users after: %+v", users)

	err = Save(userData, users)
	return err
}

// Delete takes an ID and removes the item from the db.
// An error is returned if the process fails.
func (db *User) Delete(id int) error {
	users, err := db.fetchUsers()
	if err != nil {
		return err
	}
	for i, v := range users {
		if id == v.ID {
			users = append(users[:i], users[i+1:]...)
		}
	}

	if err = Save(userData, users); err != nil {
		return errors.New("problem making updates, please try again")
	}

	// if user has an associated rental, delete rental
	r := Rental{}
	if err = r.cascade(0, id); err != nil {
		return errors.New("rental not found")
	}

	return nil
}

// fetchUsers retrieves all users from the DB and parses the data.
// Upon success, a slice of users is returned.
// On failure, an error is returned.
func (db *User) fetchUsers() (users []User, err error) {
	data, err := Load(userData)
	if err != nil {
		if err != nil {
			return nil, err
		}
	}

	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}

	return users, nil
}
