package carta

import (
	"reflect"
	"testing"
)

func TestSetDstAssociation(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	m, err := newMapper(reflect.TypeOf(&User{}))
	if err != nil {
		t.Fatalf("error creating new mapper: %s", err)
	}

	rsv := newResolver()
	elem := reflect.New(m.Typ).Elem()
	elem.FieldByName("ID").SetInt(1)
	elem.FieldByName("Name").SetString("John Doe")
	rsv.elements["1"] = &element{v: elem}
	rsv.elementOrder = append(rsv.elementOrder, "1")

	var user User
	dst := reflect.ValueOf(&user)

	err = setDst(m, dst, rsv)
	if err != nil {
		t.Fatalf("error setting destination: %s", err)
	}

	if user.ID != 1 || user.Name != "John Doe" {
		t.Errorf("expected user to be {ID:1 Name:\"John Doe\"}, but got %+v", user)
	}
}
