package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

/*
	verifying both marshal and unmarshal are needed in libovsdb when receiveing
	update() payload with a simpler workload

	converting an interface{} -> map[outerKey]InnerStruct
*/

type Outer struct {
	Ins map[string]Inner `json:"ins,overflow"`
}

type Inner struct {
	ColumnValue interface{} `json:"value,overflow"`
}

// main function
func BenchmarkCast(b *testing.B) {
	mockPayload := getMockPayload()
	params := []interface{}{mockPayload}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		castResult := cast(params[0])
		if castResult != nil {
			// do nothing
		}
	}
}

func BenchmarkMarshalUnmarshal(b *testing.B) {
	mockPayload := getMockPayload()
	params := []interface{}{mockPayload}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		marshalUnmarshalResult := marshalAndUnmarshal(params[0])
		if marshalUnmarshalResult != nil {
			// do nothing
		}
	}

}

func TestEqual(t *testing.T) {
	mockPayload := getMockPayload()
	params := []interface{}{mockPayload}

	marshalUnmarshalResult := fmt.Sprintf("%v", marshalAndUnmarshal(params[0]))
	castResult := fmt.Sprintf("%v", cast(params[0]))

	if castResult != marshalUnmarshalResult {
		t.Fatalf("sprintf() results not equal")
	}

	t.Logf("\nresult during equality test: \n%s\n", castResult)
}

func marshalAndUnmarshal(param interface{}) map[string]Inner {

	raw := param
	b, err := json.Marshal(raw)

	if err != nil {
		fmt.Println("2")
		return nil
	}

	var rowUpdates map[string]Inner

	if err := json.Unmarshal(b, &rowUpdates); err != nil {
		fmt.Println("3")
		return nil
	}

	return rowUpdates
}

func cast(param interface{}) map[string]Inner {

	res := make(map[string]Inner)
	outerRaw := param.(map[string]interface{})

	for outerKey, insRaw := range outerRaw {
		inner := insRaw.(map[string]interface{})
		res[outerKey] = Inner{ColumnValue: inner["value"]}
	}

	return res
}

/*
	Helper function
*/
func getMockPayload() interface{} {
	outer := Outer{
		Ins: map[string]Inner{
			"out_key1": Inner{ColumnValue: 42},
			"out_key2": Inner{ColumnValue: []string{"a", "b", "c"}},
			"out_key3": Inner{ColumnValue: "xyz"},
			"out_key4": Inner{ColumnValue: map[string]int{"lol": 9}},
			"out_key5": Inner{ColumnValue: 42},
			"out_key6": Inner{ColumnValue: []string{"a", "b", "c"}},
			"out_key7": Inner{ColumnValue: "xyz"},
			"out_key8": Inner{ColumnValue: map[string]int{"lol": 9}},
		},
	}

	// build the db payload from the actual struct
	temp := make(map[string]Inner)

	for outKey, inner := range outer.Ins {
		temp[outKey] = inner
	}

	// marshal it to a byte[]
	b, err := json.Marshal(temp)
	if err != nil {
		fmt.Println("unable to marshal struct payload")
		return nil
	}

	// unmarshal it into an interface to mock the type of the object in the callback
	var res interface{}
	if err := json.Unmarshal(b, &res); err != nil {
		fmt.Println("unable to unmarshal payload bytes into interface")
		return nil
	}

	return res
}
