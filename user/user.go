package user

import (
	"errors"

	"github.com/asdine/storm/v3"
	"gopkg.in/mgo.v2/bson"
)

// Move DB, add init (create if not exist)

type User struct {
	ID   bson.ObjectId `json:"id" storm:"id"`
	Name string        `json:"name"`
	Role string        `json:"role"`
}

const (
	dbPath = "users.db"
)

// Errors - TODO: add more here
var (
	ErrRecordInvalid = errors.New("record is invalid")
)

// Get all users from the DB
func GetAll() ([]User, error) {
	db, err := storm.Open(dbPath)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	users := []User{}
	err = db.All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Get a single user from the DB by ID
func GetOne(id bson.ObjectId) (*User, error) {
	db, err := storm.Open(dbPath)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	u := new(User)
	err = db.One("ID", id, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Delete a single user from the DB by ID
func DeleteOne(id bson.ObjectId) error {
	db, err := storm.Open(dbPath)
	if err != nil {
		return err
	}

	defer db.Close()

	u := new(User)
	err = db.One("ID", id, u)
	if err != nil {
		return err
	}

	return db.DeleteStruct(u)
}

// Save updates to the DB
func (u *User) Save() error {
	if err := u.validateRecordData(); err != nil {
		return err
	}
	db, err := storm.Open(dbPath)
	if err != nil {
		return err
	}

	defer db.Close()

	return db.Save(u)
}

// Validate that the record contains valid data
func (u *User) validateRecordData() error {
	if u.Name == "" {
		return ErrRecordInvalid
	}
	return nil
}
