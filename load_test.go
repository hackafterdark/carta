package carta

import (
	"reflect"
	"testing"
	"time"

	"github.com/hackafterdark/carta/value"
)

func TestLoadRow(t *testing.T) {
	m, err := newMapper(reflect.TypeOf(&User{}))
	if err != nil {
		t.Fatalf("error creating new mapper: %s", err)
	}
	determineFieldsNames(m)

	columns := map[string]column{
		"ID": {
			name:        "ID",
			columnIndex: 0,
			i:           0,
		},
		"Name": {
			name:        "Name",
			columnIndex: 1,
			i:           1,
		},
	}
	allocateColumns(m, columns)

	idCell := value.NewCell("INT")
	idCell.SetInt64(1)
	nameCell := value.NewCell("VARCHAR")
	nameCell.SetString("John Doe")
	row := []interface{}{
		idCell,
		nameCell,
	}

	rsv := newResolver()
	err = loadRow(m, row, rsv)
	if err != nil {
		t.Fatalf("error loading row: %s", err)
	}

	if len(rsv.elements) != 1 {
		t.Fatalf("expected 1 element, got %d", len(rsv.elements))
	}

	for _, elem := range rsv.elements {
		user := elem.v.Interface().(User)
		if user.ID != 1 || user.Name != "John Doe" {
			t.Errorf("expected user to be {ID:1 Name:\"John Doe\"}, but got %+v", user)
		}
	}
}

func TestLoadRowNullValue(t *testing.T) {
	m, err := newMapper(reflect.TypeOf(&UserWithNullName{}))
	if err != nil {
		t.Fatalf("error creating new mapper: %s", err)
	}
	determineFieldsNames(m)

	columns := map[string]column{
		"ID": {
			name:        "ID",
			columnIndex: 0,
			i:           0,
		},
		"Name": {
			name:        "Name",
			columnIndex: 1,
			i:           1,
		},
	}
	allocateColumns(m, columns)

	idCell := value.NewCell("INT")
	idCell.SetInt64(1)
	nameCell := value.NewCell("VARCHAR")
	nameCell.SetNull()

	row := []interface{}{
		idCell,
		nameCell,
	}

	rsv := newResolver()
	err = loadRow(m, row, rsv)
	if err != nil {
		t.Fatalf("error loading row: %s", err)
	}

	if len(rsv.elements) != 1 {
		t.Fatalf("expected 1 element, got %d", len(rsv.elements))
	}

	for _, elem := range rsv.elements {
		user := elem.v.Interface().(UserWithNullName)
		if user.ID != 1 {
			t.Errorf("expected user ID to be 1, but got %d", user.ID)
		}
		if user.Name != nil {
			t.Errorf("expected user Name to be nil, but got %s", *user.Name)
		}
	}
}

func TestLoadRowDataTypes(t *testing.T) {
	m, err := newMapper(reflect.TypeOf(&AllTypes{}))
	if err != nil {
		t.Fatalf("error creating new mapper: %s", err)
	}
	determineFieldsNames(m)

	columns := map[string]column{
		"Bool": {
			name:        "Bool",
			columnIndex: 0,
			i:           0,
		},
		"Uint64": {
			name:        "Uint64",
			columnIndex: 1,
			i:           1,
		},
		"Float64": {
			name:        "Float64",
			columnIndex: 2,
			i:           2,
		},
		"Time": {
			name:        "Time",
			columnIndex: 3,
			i:           3,
		},
	}
	allocateColumns(m, columns)

	now := time.Now()
	boolCell := value.NewCell("BOOL")
	boolCell.SetBool(true)
	uint64Cell := value.NewCell("BIGINT")
	uint64Cell.SetInt64(12345) // SetInt64 is used for uints as well
	float64Cell := value.NewCell("FLOAT")
	float64Cell.SetFloat64(123.45)
	timeCell := value.NewCell("TIMESTAMP")
	timeCell.SetTime(now)

	row := []interface{}{
		boolCell,
		uint64Cell,
		float64Cell,
		timeCell,
	}

	rsv := newResolver()
	err = loadRow(m, row, rsv)
	if err != nil {
		t.Fatalf("error loading row: %s", err)
	}

	if len(rsv.elements) != 1 {
		t.Fatalf("expected 1 element, got %d", len(rsv.elements))
	}

	for _, elem := range rsv.elements {
		data := elem.v.Interface().(AllTypes)
		if data.Bool != true {
			t.Errorf("expected Bool to be true, but got %v", data.Bool)
		}
		if data.Uint64 != 12345 {
			t.Errorf("expected Uint64 to be 12345, but got %d", data.Uint64)
		}
		if data.Float64 != 123.45 {
			t.Errorf("expected Float64 to be 123.45, but got %f", data.Float64)
		}
		if !data.Time.Equal(now) {
			t.Errorf("expected Time to be %v, but got %v", now, data.Time)
		}
	}
}

func TestLoadRowConversionError(t *testing.T) {
	m, err := newMapper(reflect.TypeOf(&User{}))
	if err != nil {
		t.Fatalf("error creating new mapper: %s", err)
	}
	determineFieldsNames(m)

	columns := map[string]column{
		"ID": {
			name:        "ID",
			columnIndex: 0,
			i:           0,
		},
		"Name": {
			name:        "Name",
			columnIndex: 1,
			i:           1,
		},
	}
	allocateColumns(m, columns)

	idCell := value.NewCell("INT")
	idCell.SetString("not a number") // This should cause a conversion error
	nameCell := value.NewCell("VARCHAR")
	nameCell.SetString("John Doe")

	row := []interface{}{
		idCell,
		nameCell,
	}

	rsv := newResolver()
	err = loadRow(m, row, rsv)
	if err == nil {
		t.Fatalf("expected a conversion error, but got nil")
	}
}
