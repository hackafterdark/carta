package carta

import (
	"errors"
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

type Blog struct {
	ID     int    `db:"id"`
	Title  string `db:"title"`
	Author Author `carta:"author"`
}

type Author struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func TestCartaTagMap(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title", "author_id", "author_name"}).
		AddRow(1, "My First Blog", 101, "John Doe")

	mock.ExpectQuery("SELECT (.+) FROM blogs").WillReturnRows(rows)

	sqlRows, err := db.Query("SELECT * FROM blogs")
	if err != nil {
		t.Fatalf("error '%s' was not expected when querying rows", err)
	}

	var blogs []Blog
	err = Map(sqlRows, &blogs)
	if err != nil {
		t.Errorf("error was not expected while mapping rows: %s", err)
	}

	if len(blogs) != 1 {
		t.Fatalf("expected 1 blog, got %d", len(blogs))
	}

	expectedBlog := Blog{
		ID:    1,
		Title: "My First Blog",
		Author: Author{
			ID:   101,
			Name: "John Doe",
		},
	}

	if !reflect.DeepEqual(blogs[0], expectedBlog) {
		t.Errorf("expected blog to be %+v, but got %+v", expectedBlog, blogs[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMapRowsScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Name"}).
		AddRow(1, "John Doe").
		RowError(0, errors.New("scan error"))

	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)

	sqlRows, err := db.Query("SELECT * FROM users")
	if err != nil {
		t.Fatalf("error '%s' was not expected when querying rows", err)
	}

	var users []User
	err = Map(sqlRows, &users)
	if err == nil {
		t.Errorf("expected an error when mapping rows with a scan error, but got nil")
	}
}

type BlogWithCustomDelimiter struct {
	ID     int    `db:"id"`
	Title  string `db:"title"`
	Author Author `carta:"author,delimiter=_"`
}

func TestCartaTagCustomDelimiterMap(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title", "author_id", "author_name"}).
		AddRow(1, "My First Blog", 101, "John Doe")

	mock.ExpectQuery("SELECT (.+) FROM blogs").WillReturnRows(rows)

	sqlRows, err := db.Query("SELECT * FROM blogs")
	if err != nil {
		t.Fatalf("error '%s' was not expected when querying rows", err)
	}

	var blogs []BlogWithCustomDelimiter
	err = Map(sqlRows, &blogs)
	if err != nil {
		t.Errorf("error was not expected while mapping rows: %s", err)
	}

	if len(blogs) != 1 {
		t.Fatalf("expected 1 blog, got %d", len(blogs))
	}

	expectedBlog := BlogWithCustomDelimiter{
		ID:    1,
		Title: "My First Blog",
		Author: Author{
			ID:   101,
			Name: "John Doe",
		},
	}

	if !reflect.DeepEqual(blogs[0], expectedBlog) {
		t.Errorf("expected blog to be %+v, but got %+v", expectedBlog, blogs[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMapDeeplyNested(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	label1, label2, label3 := "Label 1", "Label 2", "Label 3"
	labelID1, labelID2, labelID3 := 1001, 1002, 1003
	rows := sqlmock.NewRows([]string{"id", "name", "posts_id", "posts_title", "posts_labels_id", "posts_labels_name"}).
		AddRow(1, "Blog 1", 101, "Post 1", &labelID1, &label1).
		AddRow(1, "Blog 1", 101, "Post 1", &labelID2, &label2).
		AddRow(1, "Blog 1", 102, "Post 2", &labelID3, &label3)

	mock.ExpectQuery("SELECT (.+) FROM blogs").WillReturnRows(rows)

	sqlRows, err := db.Query("SELECT * FROM blogs")
	if err != nil {
		t.Fatalf("error '%s' was not expected when querying rows", err)
	}

	var blogs []BlogWithPosts
	err = Map(sqlRows, &blogs)
	if err != nil {
		t.Errorf("error was not expected while mapping rows: %s", err)
	}

	if len(blogs) != 1 {
		t.Fatalf("expected 1 blog, got %d", len(blogs))
	}

	if len(blogs[0].Posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(blogs[0].Posts))
	}

	if len(blogs[0].Posts[0].Labels) != 2 {
		t.Fatalf("expected 2 labels for the first post, got %d", len(blogs[0].Posts[0].Labels))
	}

	if len(blogs[0].Posts[1].Labels) != 1 {
		t.Fatalf("expected 1 label for the second post, got %d", len(blogs[0].Posts[1].Labels))
	}

	if *blogs[0].Posts[0].Labels[0].Name != "Label 1" {
		t.Errorf("incorrect label name")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMapDeeplyNestedNoLabels(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	label3 := "Label 3"
	labelID3 := 1003
	rows := sqlmock.NewRows([]string{"id", "name", "posts_id", "posts_title", "posts_labels_id", "posts_labels_name"}).
		AddRow(1, "Blog 1", 101, "Post 1", nil, nil).
		AddRow(1, "Blog 1", 102, "Post 2", &labelID3, &label3)

	mock.ExpectQuery("SELECT (.+) FROM blogs").WillReturnRows(rows)

	sqlRows, err := db.Query("SELECT * FROM blogs")
	if err != nil {
		t.Fatalf("error '%s' was not expected when querying rows", err)
	}

	var blogs []BlogWithPosts
	err = Map(sqlRows, &blogs)
	if err != nil {
		t.Errorf("error was not expected while mapping rows: %s", err)
	}

	if len(blogs) != 1 {
		t.Fatalf("expected 1 blog, got %d", len(blogs))
	}

	if len(blogs[0].Posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(blogs[0].Posts))
	}

	if len(blogs[0].Posts[0].Labels) != 0 {
		t.Fatalf("expected 0 labels for the first post, got %d", len(blogs[0].Posts[0].Labels))
	}

	if len(blogs[0].Posts[1].Labels) != 1 {
		t.Fatalf("expected 1 label for the second post, got %d", len(blogs[0].Posts[1].Labels))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMapToBasicSlice(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"tag"}).
		AddRow("tag1").
		AddRow("tag2").
		AddRow("tag1") // Duplicate tag

	mock.ExpectQuery("SELECT (.+) FROM tags").WillReturnRows(rows)

	sqlRows, err := db.Query("SELECT * FROM tags")
	if err != nil {
		t.Fatalf("error '%s' was not expected when querying rows", err)
	}

	var tags []string
	err = Map(sqlRows, &tags)
	if err != nil {
		t.Errorf("error was not expected while mapping rows: %s", err)
	}

	if len(tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(tags))
	}

	expectedTags := []string{"tag1", "tag2", "tag1"}
	if !reflect.DeepEqual(tags, expectedTags) {
		t.Errorf("expected tags to be %+v, but got %+v", expectedTags, tags)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMapToBasicSlice_MultipleColumnsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"tag", "extra"}).
		AddRow("tag1", "x")

	mock.ExpectQuery("SELECT (.+) FROM tags").WillReturnRows(rows)

	sqlRows, err := db.Query("SELECT * FROM tags")
	if err != nil {
		t.Fatal(err)
	}

	var tags []string
	err = Map(sqlRows, &tags)
	if err == nil {
		t.Fatalf("expected error when mapping to []string with multiple columns, got nil")
	}
}
