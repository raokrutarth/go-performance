package main

import (
	"encoding/json"
	"fmt"
)

/*
	verifying both marshal and unmarshal are needed in libovsdb when receiveing
	update() payload
*/

type TableUpdates struct {
	Updates map[string]TableUpdate `json:"updates,overflow"`
}

// TableUpdate represents a table update according to RFC7047
type TableUpdate struct {
	Rows map[string]RowUpdate `json:"rows,overflow"`
}

// RowUpdate represents a row update according to RFC7047
// It can also hold an Update2 notification as described in ovsdb-server(1)
type RowUpdate struct {
	UUID    UUID `json:"-,omitempty"`
	New     Row  `json:"new,omitempty"`
	Old     Row  `json:"old,omitempty"`
	Initial Row  `json:"initial,omitempty"`
	Insert  Row  `json:"insert,omitempty"`
	Delete  Row  `json:"delete,omitempty"`
	Modify  Row  `json:"modify,omitempty"`
}

// UUID is a UUID according to RFC7047
type UUID struct {
	GoUUID string `json:"uuid"`
}

// Row is a table Row according to RFC7047
type Row struct {
	Fields map[string]interface{}
}

func main() {
	mock_param := getMockDBPayload()
	params := []interface{}{mock_param}

	fmt.Printf("with marshal & unmarshal: %v", marshalAndUnmarshal(params))
	fmt.Println("\n")
	fmt.Printf("with cast: %v", cast(params))

}

func marshalAndUnmarshal(params []interface{}) map[string]map[string]RowUpdate {

	raw, ok := params[0].(map[string]interface{})
	if !ok {
		fmt.Println("1")
		return nil
	}

	b, err := json.Marshal(raw)

	if err != nil {
		fmt.Println("2")
		return nil
	}

	var rowUpdates map[string]map[string]RowUpdate
	if err := json.Unmarshal(b, &rowUpdates); err != nil {
		fmt.Println("3")
		return nil
	}
	return rowUpdates
}

func cast(params []interface{}) map[string]map[string]RowUpdate {
	// res := make(map[string]map[string]RowUpdate)

	tableUpdatesRaw, ok := params[0].(map[string]interface{})

	// tableUpdatesRaw -> ...

	if !ok {
		fmt.Println("9")
	}

	res := make(map[string]map[string]RowUpdate)

	for tableName, tableUpdateRaw := range tableUpdatesRaw {

		tableUpdate := tableUpdateRaw.(map[string]interface{})

		res[tableName] = make(map[string]RowUpdate)

		for rowUUID, rowUpdateRaw := range tableUpdate {
			fmt.Printf("%v\n", rowUpdateRaw)
			res[tableName][rowUUID] = rowUpdateRaw.(RowUpdate) // this fails.
			// i.e. need to cast until the very base primitive (i.e. map) to avoid
			// the marshal and unmarshal
		}

	}

	return res
}

/*
	Helper function
*/
func getMockDBPayload() interface{} {
	tu := TableUpdates{
		Updates: map[string]TableUpdate{
			"table1": TableUpdate{
				Rows: map[string]RowUpdate{
					"row1": RowUpdate{
						UUID: UUID{"10"},
						New: Row{
							Fields: map[string]interface{}{
								"col0": 42,
								"col1": map[string]int{
									"key1": 1,
									"key2": 2,
								},
							},
						},
					},
					"row2": RowUpdate{
						UUID: UUID{"20"},
						New: Row{
							Fields: map[string]interface{}{
								"col1": 43,
							},
						},
					},
				},
			},
			"table2": TableUpdate{
				Rows: map[string]RowUpdate{
					"row1.1": RowUpdate{
						UUID: UUID{"10.1"},
						New: Row{
							Fields: map[string]interface{}{
								"col0": 42.1,
							},
						},
					},
				},
			},
		},
	}

	// build the db payload from the actual struct
	temp := make(map[string]map[string]RowUpdate)
	for tableName, tableUpdate := range tu.Updates {
		temp[tableName] = map[string]RowUpdate{}
		for rowUUID, rowUpdate := range tableUpdate.Rows {
			temp[tableName][rowUUID] = rowUpdate
		}
	}

	// marshal it to a byte[]
	b, err := json.Marshal(temp)
	if err != nil {
		fmt.Println("2.m")
		return nil
	}

	// unmarshal it into an interface to mock the type of the object in the callback
	var res interface{}
	if err := json.Unmarshal(b, &res); err != nil {
		fmt.Println("3.m")
		return nil
	}

	return res
}
