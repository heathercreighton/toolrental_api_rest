package database

import (
	"encoding/json"
	"errors"
	"sort"
)

// Tool defines the tools data model.
type Tool struct {
	ID       int     `json:"id"`
	Name     string  `json:"name,omitempty"`
	Desc     string  `json:"desc,omitempty"`
	Price    float64 `json:"price,omitempty"`
	Quantity int     `json:"quantity,omitempty"`
}

// All returns all records, sorted in ascending order for a given resource.
// An error is returned if the read process fails.
func (db *Tool) All() ([]Tool, error) {
	tools, err := db.fetchTools()
	if err != nil {
		return nil, err
	}
	return tools, nil
}

// FindByID locates the correct tool based on ID.
func (db *Tool) FindByID(id int) (*Tool, error) {
	tools, err := db.fetchTools()
	if err != nil {
		return nil, err
	}
	for _, t := range tools {
		if t.ID == id {
			return &t, nil
		}
	}
	return nil, errors.New("tool not found")
}

// Create takes in a tool object and adds the item to the db.
// An error is returned if the request body is empty,
// or if any error occurs when reading/writing data.
func (db *Tool) Create(tool *Tool) error {
	tools, err := db.fetchTools()
	if err != nil {
		return err
	}

	tool.ID = tools[len(tools)-1].ID + 1

	tools = append(tools, *tool)
	err = Save(toolData, tools)

	return err
}

// Update takes a Tool ID and tool object and handles updating the db.
// An error is returned if the process fails.
func (db *Tool) Update(id int, tool *Tool) error {
	tools, err := db.All()
	if err != nil {
		return err
	}
	for i, v := range tools {

		if v.ID == id {
			tool.ID = id
			tools[i] = *tool
			break
		}
	}

	err = Save(toolData, tools)
	return err
}

// Delete takes an ID and removes the item from the db.
// An error is returned if the process fails.
func (db *Tool) Delete(id int) error {
	tools, err := db.fetchTools()
	if err != nil {
		return err
	}
	for i, v := range tools {
		if id == v.ID {
			tools = append(tools[:i], tools[i+1:]...)
		}
	}

	if err = Save(toolData, tools); err != nil {
		return errors.New("problem making updates, please try again")
	}

	// if tool has an associated rental, delete rental
	r := Rental{}
	if err = r.cascade(id, 0); err != nil {
		return errors.New("rental not found")
	}

	return nil
}

// sortRecords handles sorting records by ID in ascending order.
func sortRecords(tools []Tool) []Tool {
	sort.SliceStable(tools, func(i, j int) bool {
		return tools[i].ID < tools[j].ID
	})
	return tools
}

// fetchTools retrieves all tools from the DB and parses the data.
// Upon success, a slice of tools is returned.
// On failure, an error is returned.
func (db *Tool) fetchTools() (tools []Tool, err error) {
	data, err := Load(toolData)
	if err != nil {
		if err != nil {
			return nil, err
		}
	}

	if err := json.Unmarshal(data, &tools); err != nil {
		return nil, err
	}

	return tools, nil
}

// UpdateQuantity handles updating the amount of available items.
func (db *Tool) UpdateQuantity(id int, action string) error {
	t, err := db.FindByID(id)
	if err != nil {
		return err
	}
	switch action {
	case "add":
		t.Quantity = t.Quantity + 1
	case "subtract":
		t.Quantity = t.Quantity - 1
	}
	t.Update(t.ID, t)
	return nil
}
