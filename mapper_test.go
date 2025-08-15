package carta

import (
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type User struct {
	ID   int
	Name string
}

type UserWithAddress struct {
	ID      int
	Name    string
	Address Address
}

type Address struct {
	Street string
	City   string
}

type UserWithPosts struct {
	ID    int
	Name  string
	Posts []Post
}

type Post struct {
	Title   string
	Content string
}

type UserWithTags struct {
	ID   int    `db:"user_id"`
	Name string `db:"user_name"`
}

type UserWithNullName struct {
	ID   int
	Name *string
}

type AllTypes struct {
	Bool    bool
	Uint64  uint64
	Float64 float64
	Time    time.Time
}

func TestSimpleMap(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Name"}).
		AddRow(1, "John Doe")

	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)

	sqlRows, err := db.Query("SELECT * FROM users")
	if err != nil {
		t.Fatalf("error '%s' was not expected when querying rows", err)
	}

	var users []User
	err = Map(sqlRows, &users)
	if err != nil {
		t.Errorf("error was not expected while mapping rows: %s", err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}

	if users[0].ID != 1 || users[0].Name != "John Doe" {
		t.Errorf("expected user to be {ID:1 Name:\"John Doe\"}, but got {ID:%d Name:\"%s\"}", users[0].ID, users[0].Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPointerMap(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Name"}).
		AddRow(1, "John Doe")

	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)

	sqlRows, err := db.Query("SELECT * FROM users")
	if err != nil {
		t.Fatalf("error '%s' was not expected when querying rows", err)
	}

	var users []*User
	err = Map(sqlRows, &users)
	if err != nil {
		t.Errorf("error was not expected while mapping rows: %s", err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}

	if users[0].ID != 1 || users[0].Name != "John Doe" {
		t.Errorf("expected user to be {ID:1 Name:\"John Doe\"}, but got {ID:%d Name:\"%s\"}", users[0].ID, users[0].Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestNestedMap(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Name", "Address_Street", "Address_City"}).
		AddRow(1, "John Doe", "123 Main St", "Anytown")

	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)

	sqlRows, err := db.Query("SELECT * FROM users")
	if err != nil {
		t.Fatalf("error '%s' was not expected when querying rows", err)
	}

	var users []UserWithAddress
	err = Map(sqlRows, &users)
	if err != nil {
		t.Errorf("error was not expected while mapping rows: %s", err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}

	expectedUser := UserWithAddress{
		ID:   1,
		Name: "John Doe",
		Address: Address{
			Street: "123 Main St",
			City:   "Anytown",
		},
	}

	if users[0].ID != expectedUser.ID || users[0].Name != expectedUser.Name || users[0].Address.Street != expectedUser.Address.Street || users[0].Address.City != expectedUser.Address.City {
		t.Errorf("expected user to be %+v, but got %+v", expectedUser, users[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestNestedCollectionMap(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Name", "Posts_Title", "Posts_Content"}).
		AddRow(1, "John Doe", "First Post", "Hello World").
		AddRow(1, "John Doe", "Second Post", "Another post")

	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)

	sqlRows, err := db.Query("SELECT * FROM users")
	if err != nil {
		t.Fatalf("error '%s' was not expected when querying rows", err)
	}

	var users []UserWithPosts
	err = Map(sqlRows, &users)
	if err != nil {
		t.Errorf("error was not expected while mapping rows: %s", err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}

	if len(users[0].Posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(users[0].Posts))
	}

	if users[0].Posts[0].Title != "First Post" || users[0].Posts[1].Title != "Second Post" {
		t.Errorf("posts not mapped correctly")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTagMap(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"user_id", "user_name"}).
		AddRow(1, "John Doe")

	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)

	sqlRows, err := db.Query("SELECT * FROM users")
	if err != nil {
		t.Fatalf("error '%s' was not expected when querying rows", err)
	}

	var users []UserWithTags
	err = Map(sqlRows, &users)
	if err != nil {
		t.Errorf("error was not expected while mapping rows: %s", err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}

	if users[0].ID != 1 || users[0].Name != "John Doe" {
		t.Errorf("expected user to be {ID:1 Name:\"John Doe\"}, but got {ID:%d Name:\"%s\"}", users[0].ID, users[0].Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMapToNonPointer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Name"}).
		AddRow(1, "John Doe")

	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)

	sqlRows, err := db.Query("SELECT * FROM users")
	if err != nil {
		t.Fatalf("error '%s' was not expected when querying rows", err)
	}

	var users []User
	err = Map(sqlRows, users) // Pass by value instead of by pointer
	if err == nil {
		t.Errorf("expected an error when mapping to a non-pointer, but got nil")
	}
}

func TestNewMapperError(t *testing.T) {
	_, err := newMapper(reflect.TypeOf(1)) // Pass an invalid type
	if err == nil {
		t.Errorf("expected an error when creating a new mapper with an invalid type, but got nil")
	}
}
