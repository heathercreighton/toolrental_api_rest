package database

import (
	"encoding/json"
	"errors"
	"fmt"
)

var tool = &Tool{}
var user = &User{}
var rental = &Rental{}

// Rental defines the rental data model.
type Rental struct {
	ID     int  `json:"id"`
	Active bool `json:"active"`
	UserID int  `json:"user_id"`
	ToolID int  `json:"tool_id"`
}

// RentalData contains User & Tool info.
type RentalData struct {
	*Rental `json:"rental"`
	*User   `json:"user"`
	*Tool   `json:"tool"`
}

// All handles gathering rental info and returns a RentalData object.
// An error is returned if any record is missing.
func (db *Rental) All() ([]RentalData, error) {
	rentals, err := db.fetchRentals()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	records, err := generateRentalData(rentals)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return records, nil
}

func generateRentalData(r []Rental) (records []RentalData, err error) {
	for i, v := range r {
		t, err := tool.FindByID(v.ToolID)
		u, err := user.FindByID(v.UserID)
		if err != nil {
			return nil, err
		}

		records = append(records, RentalData{&r[i], u, t})
	}
	return
}

// fetchRentals retrieves all rentals from the DB and parses the data.
// Upon success, a slice of rentals is returned.
// On failure, an error is returned.
func (db *Rental) fetchRentals() ([]Rental, error) {
	data, err := Load(rentalData)
	if err != nil {
		return nil, err
	}
	var rentals []Rental
	var activeRentals []Rental
	if err := json.Unmarshal(data, &rentals); err != nil {
		return nil, err
	}

	for _, v := range rentals {
		if v.Active {
			activeRentals = append(activeRentals, v)
		}
	}

	return activeRentals, nil
}

// FindByID locates the correct rental based on ID.
func (db *Rental) FindByID(id int) (*RentalData, error) {
	rentals, err := db.fetchRentals()
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("rental not found")
	}

	var records []Rental
	for _, v := range rentals {
		if v.ID == id {
			records = append(records, v)
			break
		}
	}

	if len(records) == 0 {
		fmt.Println(err)
		return nil, errors.New("rental not found")
	}

	data, err := generateRentalData(records)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("rental not found")
	}

	return &data[0], nil
}

// cascade is responsible for deleting associated rentals
// when a tool or user is deleted.
func (db *Rental) cascade(toolID, userID int) (err error) {
	rentals, err := db.fetchRentals()
	if err != nil {
		return err
	}

	for i, v := range rentals {
		if v.ToolID == toolID || v.UserID == userID {
			rentals = append(rentals[:i], rentals[i+1:]...)
		}
	}
	err = Save(rentalData, rentals)

	return err
}

// Create takes in a rental object and adds the item to the db.
// An error is returned if the request body is empty,
// or if any error occurs when reading/writing data.
func (db *Rental) Create(rental *Rental) error {
	rentals, err := db.fetchRentals()
	if err != nil {
		return err
	}

	rental.ID = rentals[len(rentals)-1].ID + 1
	rental.Active = true

	rentals = append(rentals, *rental)

	err = Save(rentalData, rentals)
	// find tool and update quantity
	err = tool.UpdateQuantity(rental.ToolID, "subtract")
	return err
}

// Update takes a rental ID and rental object and handles updating the db.
// An error is returned if the process fails.
func (db *Rental) Update(id int, rental *Rental) error {
	rentals, err := db.fetchRentals()
	if err != nil {
		fmt.Println(err)
		return err
	}
	for i, v := range rentals {
		if v.ID == id {
			rental.ID = id
			rentals[i] = *rental
			break
		}
	}

	err = Save(rentalData, rentals)
	return err
}

// Delete takes an ID and removes the item from the db.
// An error is returned if the process fails.
func (db *Rental) Delete(id int) error {
	rentals, err := db.fetchRentals()
	if err != nil {
		fmt.Println(err)
		return err
	}

	for i, v := range rentals {
		if v.ID == id {
			if err := tool.UpdateQuantity(v.ToolID, "add"); err != nil {
				fmt.Println("error: ", err)
				return err
			}
			rentals[i].Active = false
			break
		}
	}

	err = Save(rentalData, rentals)
	return err
}
