package carta

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/hackafterdark/carta/value"
)

type AllNullTypes struct {
	NullBool    sql.NullBool
	NullFloat64 sql.NullFloat64
	NullInt32   sql.NullInt32
	NullInt64   sql.NullInt64
	NullString  sql.NullString
	NullTime    sql.NullTime
}

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
	err = loadRow(m, row, rsv, 0)
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
	err = loadRow(m, row, rsv, 0)
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
	err = loadRow(m, row, rsv, 0)
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
	err = loadRow(m, row, rsv, 0)
	if err == nil {
		t.Fatalf("expected a conversion error, but got nil")
	}
}

func TestLoadRowNullToNonNull(t *testing.T) {
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
	idCell.SetNull()
	nameCell := value.NewCell("VARCHAR")
	nameCell.SetString("John Doe")

	row := []interface{}{
		idCell,
		nameCell,
	}

	rsv := newResolver()
	err = loadRow(m, row, rsv, 0)
	if err == nil {
		t.Fatalf("expected an error when loading a null value to a non-nullable field, but got nil")
	}
}

func TestLoadRowNullTypes(t *testing.T) {
	m, err := newMapper(reflect.TypeOf(&AllNullTypes{}))
	if err != nil {
		t.Fatalf("error creating new mapper: %s", err)
	}
	determineFieldsNames(m)

	columns := map[string]column{
		"NullBool":    {name: "NullBool", columnIndex: 0, i: 0},
		"NullFloat64": {name: "NullFloat64", columnIndex: 1, i: 1},
		"NullInt32":   {name: "NullInt32", columnIndex: 2, i: 2},
		"NullInt64":   {name: "NullInt64", columnIndex: 3, i: 3},
		"NullString":  {name: "NullString", columnIndex: 4, i: 4},
		"NullTime":    {name: "NullTime", columnIndex: 5, i: 5},
	}
	allocateColumns(m, columns)

	now := time.Now()
	row := []interface{}{
		value.NewCellWithData("BOOL", true),
		value.NewCellWithData("FLOAT", 123.45),
		value.NewCellWithData("INT", int32(123)),
		value.NewCellWithData("BIGINT", int64(12345)),
		value.NewCellWithData("VARCHAR", "hello"),
		value.NewCellWithData("TIMESTAMP", now),
	}

	rsv := newResolver()
	err = loadRow(m, row, rsv, 0)
	if err != nil {
		t.Fatalf("error loading row: %s", err)
	}

	if len(rsv.elements) != 1 {
		t.Fatalf("expected 1 element, got %d", len(rsv.elements))
	}

	for _, elem := range rsv.elements {
		data := elem.v.Interface().(AllNullTypes)
		if !data.NullBool.Valid || data.NullBool.Bool != true {
			t.Errorf("expected NullBool to be valid and true")
		}
		if !data.NullFloat64.Valid || data.NullFloat64.Float64 != 123.45 {
			t.Errorf("expected NullFloat64 to be valid and 123.45")
		}
		if !data.NullInt32.Valid || data.NullInt32.Int32 != 123 {
			t.Errorf("expected NullInt32 to be valid and 123")
		}
		if !data.NullInt64.Valid || data.NullInt64.Int64 != 12345 {
			t.Errorf("expected NullInt64 to be valid and 12345")
		}
		if !data.NullString.Valid || data.NullString.String != "hello" {
			t.Errorf("expected NullString to be valid and 'hello'")
		}
		if !data.NullTime.Valid || !data.NullTime.Time.Equal(now) {
			t.Errorf("expected NullTime to be valid and equal to now")
		}
	}
}
