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

func TestSetDstCollection(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	m, err := newMapper(reflect.TypeOf(&[]User{}))
	if err != nil {
		t.Fatalf("error creating new mapper: %s", err)
	}

	rsv := newResolver()
	// User 1
	elem1 := reflect.New(m.Typ).Elem()
	elem1.FieldByName("ID").SetInt(1)
	elem1.FieldByName("Name").SetString("John Doe")
	rsv.elements["1"] = &element{v: elem1}
	rsv.elementOrder = append(rsv.elementOrder, "1")
	// User 2
	elem2 := reflect.New(m.Typ).Elem()
	elem2.FieldByName("ID").SetInt(2)
	elem2.FieldByName("Name").SetString("Jane Doe")
	rsv.elements["2"] = &element{v: elem2}
	rsv.elementOrder = append(rsv.elementOrder, "2")

	var users []User
	dst := reflect.ValueOf(&users)

	err = setDst(m, dst, rsv)
	if err != nil {
		t.Fatalf("error setting destination: %s", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected users slice to have 2 elements, but got %d", len(users))
	}

	if users[0].ID != 1 || users[0].Name != "John Doe" {
		t.Errorf("expected user 1 to be {ID:1 Name:\"John Doe\"}, but got %+v", users[0])
	}
	if users[1].ID != 2 || users[1].Name != "Jane Doe" {
		t.Errorf("expected user 2 to be {ID:2 Name:\"Jane Doe\"}, but got %+v", users[1])
	}
}

func TestSetDstCollectionOfPointers(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	m, err := newMapper(reflect.TypeOf(&[]*User{}))
	if err != nil {
		t.Fatalf("error creating new mapper: %s", err)
	}

	rsv := newResolver()
	// User 1
	elem1 := reflect.New(m.Typ).Elem()
	elem1.FieldByName("ID").SetInt(1)
	elem1.FieldByName("Name").SetString("John Doe")
	rsv.elements["1"] = &element{v: elem1}
	rsv.elementOrder = append(rsv.elementOrder, "1")
	// User 2
	elem2 := reflect.New(m.Typ).Elem()
	elem2.FieldByName("ID").SetInt(2)
	elem2.FieldByName("Name").SetString("Jane Doe")
	rsv.elements["2"] = &element{v: elem2}
	rsv.elementOrder = append(rsv.elementOrder, "2")

	var users []*User
	dst := reflect.ValueOf(&users)

	err = setDst(m, dst, rsv)
	if err != nil {
		t.Fatalf("error setting destination: %s", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected users slice to have 2 elements, but got %d", len(users))
	}

	if users[0].ID != 1 || users[0].Name != "John Doe" {
		t.Errorf("expected user 1 to be {ID:1 Name:\"John Doe\"}, but got %+v", users[0])
	}
	if users[1].ID != 2 || users[1].Name != "Jane Doe" {
		t.Errorf("expected user 2 to be {ID:2 Name:\"Jane Doe\"}, but got %+v", users[1])
	}
}

func TestSetDstNested(t *testing.T) {
	type Profile struct {
		ID    int
		Email string
	}
	type User struct {
		ID      int
		Name    string
		Profile Profile
	}

	m, err := newMapper(reflect.TypeOf(&User{}))
	if err != nil {
		t.Fatalf("error creating new mapper: %s", err)
	}

	rsv := newResolver()
	elem := reflect.New(m.Typ).Elem()
	elem.FieldByName("ID").SetInt(1)
	elem.FieldByName("Name").SetString("John Doe")
	rsv.elements["1"] = &element{v: elem, subMaps: make(map[fieldIndex]*resolver)}
	rsv.elementOrder = append(rsv.elementOrder, "1")

	// submap for Profile
	profileRsv := newResolver()
	profileElem := reflect.New(m.SubMaps[2].Typ).Elem()
	profileElem.FieldByName("ID").SetInt(101)
	profileElem.FieldByName("Email").SetString("john.doe@example.com")
	profileRsv.elements["101"] = &element{v: profileElem}
	profileRsv.elementOrder = append(profileRsv.elementOrder, "101")

	rsv.elements["1"].subMaps[fieldIndex(2)] = profileRsv

	var user User
	dst := reflect.ValueOf(&user)

	err = setDst(m, dst, rsv)
	if err != nil {
		t.Fatalf("error setting destination: %s", err)
	}

	if user.ID != 1 || user.Name != "John Doe" {
		t.Errorf("expected user to be {ID:1 Name:\"John Doe\" ...}, but got %+v", user)
	}
	if user.Profile.ID != 101 || user.Profile.Email != "john.doe@example.com" {
		t.Errorf("expected profile to be {ID:101 Email:\"john.doe@example.com\"}, but got %+v", user.Profile)
	}
}

func TestSetDstNestedPointer(t *testing.T) {
	type Profile struct {
		ID    int
		Email string
	}
	type User struct {
		ID      int
		Name    string
		Profile *Profile
	}

	m, err := newMapper(reflect.TypeOf(&User{}))
	if err != nil {
		t.Fatalf("error creating new mapper: %s", err)
	}

	rsv := newResolver()
	elem := reflect.New(m.Typ).Elem()
	elem.FieldByName("ID").SetInt(1)
	elem.FieldByName("Name").SetString("John Doe")
	rsv.elements["1"] = &element{v: elem, subMaps: make(map[fieldIndex]*resolver)}
	rsv.elementOrder = append(rsv.elementOrder, "1")

	// submap for Profile
	profileRsv := newResolver()
	profileElem := reflect.New(m.SubMaps[2].Typ).Elem()
	profileElem.FieldByName("ID").SetInt(101)
	profileElem.FieldByName("Email").SetString("john.doe@example.com")
	profileRsv.elements["101"] = &element{v: profileElem}
	profileRsv.elementOrder = append(profileRsv.elementOrder, "101")

	rsv.elements["1"].subMaps[fieldIndex(2)] = profileRsv

	var user User
	dst := reflect.ValueOf(&user)

	err = setDst(m, dst, rsv)
	if err != nil {
		t.Fatalf("error setting destination: %s", err)
	}

	if user.ID != 1 || user.Name != "John Doe" {
		t.Errorf("expected user to be {ID:1 Name:\"John Doe\" ...}, but got %+v", user)
	}
	if user.Profile == nil {
		t.Fatal("expected user profile to be non-nil")
	}
	if user.Profile.ID != 101 || user.Profile.Email != "john.doe@example.com" {
		t.Errorf("expected profile to be {ID:101 Email:\"john.doe@example.com\"}, but got %+v", user.Profile)
	}
}

func TestSetDstNestedCollection(t *testing.T) {
	type Post struct {
		ID    int
		Title string
	}
	type User struct {
		ID    int
		Name  string
		Posts []Post
	}

	m, err := newMapper(reflect.TypeOf(&User{}))
	if err != nil {
		t.Fatalf("error creating new mapper: %s", err)
	}

	rsv := newResolver()
	elem := reflect.New(m.Typ).Elem()
	elem.FieldByName("ID").SetInt(1)
	elem.FieldByName("Name").SetString("John Doe")
	rsv.elements["1"] = &element{v: elem, subMaps: make(map[fieldIndex]*resolver)}
	rsv.elementOrder = append(rsv.elementOrder, "1")

	// submap for Posts
	postsRsv := newResolver()
	// Post 1
	postElem1 := reflect.New(m.SubMaps[2].Typ).Elem()
	postElem1.FieldByName("ID").SetInt(101)
	postElem1.FieldByName("Title").SetString("First Post")
	postsRsv.elements["101"] = &element{v: postElem1}
	postsRsv.elementOrder = append(postsRsv.elementOrder, "101")
	// Post 2
	postElem2 := reflect.New(m.SubMaps[2].Typ).Elem()
	postElem2.FieldByName("ID").SetInt(102)
	postElem2.FieldByName("Title").SetString("Second Post")
	postsRsv.elements["102"] = &element{v: postElem2}
	postsRsv.elementOrder = append(postsRsv.elementOrder, "102")

	rsv.elements["1"].subMaps[fieldIndex(2)] = postsRsv

	var user User
	dst := reflect.ValueOf(&user)

	err = setDst(m, dst, rsv)
	if err != nil {
		t.Fatalf("error setting destination: %s", err)
	}

	if user.ID != 1 || user.Name != "John Doe" {
		t.Errorf("expected user to be {ID:1 Name:\"John Doe\" ...}, but got %+v", user)
	}
	if len(user.Posts) != 2 {
		t.Fatalf("expected user to have 2 posts, but got %d", len(user.Posts))
	}
	if user.Posts[0].ID != 101 || user.Posts[0].Title != "First Post" {
		t.Errorf("expected post 1 to be {ID:101 Title:\"First Post\"}, but got %+v", user.Posts[0])
	}
	if user.Posts[1].ID != 102 || user.Posts[1].Title != "Second Post" {
		t.Errorf("expected post 2 to be {ID:102 Title:\"Second Post\"}, but got %+v", user.Posts[1])
	}
}

func TestSetDstNestedCollectionWithPointerSlice(t *testing.T) {
	type Post struct {
		ID    int
		Title string
	}
	type User struct {
		ID    int
		Name  string
		Posts *[]Post
	}

	m, err := newMapper(reflect.TypeOf(&User{}))
	if err != nil {
		t.Fatalf("error creating new mapper: %s", err)
	}

	rsv := newResolver()
	elem := reflect.New(m.Typ).Elem()
	elem.FieldByName("ID").SetInt(1)
	elem.FieldByName("Name").SetString("John Doe")
	rsv.elements["1"] = &element{v: elem, subMaps: make(map[fieldIndex]*resolver)}
	rsv.elementOrder = append(rsv.elementOrder, "1")

	// submap for Posts
	postsRsv := newResolver()
	// Post 1
	postElem1 := reflect.New(m.SubMaps[2].Typ).Elem()
	postElem1.FieldByName("ID").SetInt(101)
	postElem1.FieldByName("Title").SetString("First Post")
	postsRsv.elements["101"] = &element{v: postElem1}
	postsRsv.elementOrder = append(postsRsv.elementOrder, "101")

	rsv.elements["1"].subMaps[fieldIndex(2)] = postsRsv

	var user User
	dst := reflect.ValueOf(&user)

	err = setDst(m, dst, rsv)
	if err != nil {
		t.Fatalf("error setting destination: %s", err)
	}

	if user.ID != 1 || user.Name != "John Doe" {
		t.Errorf("expected user to be {ID:1 Name:\"John Doe\" ...}, but got %+v", user)
	}
	if user.Posts == nil {
		t.Fatal("expected user posts to be non-nil")
	}
	if len(*user.Posts) != 1 {
		t.Fatalf("expected user to have 1 post, but got %d", len(*user.Posts))
	}
	if (*user.Posts)[0].ID != 101 || (*user.Posts)[0].Title != "First Post" {
		t.Errorf("expected post 1 to be {ID:101 Title:\"First Post\"}, but got %+v", (*user.Posts)[0])
	}
}

func TestSetDstNestedCollectionWithPointerElements(t *testing.T) {
	type Post struct {
		ID    int
		Title string
	}
	type User struct {
		ID    int
		Name  string
		Posts []*Post
	}

	m, err := newMapper(reflect.TypeOf(&User{}))
	if err != nil {
		t.Fatalf("error creating new mapper: %s", err)
	}

	rsv := newResolver()
	elem := reflect.New(m.Typ).Elem()
	elem.FieldByName("ID").SetInt(1)
	elem.FieldByName("Name").SetString("John Doe")
	rsv.elements["1"] = &element{v: elem, subMaps: make(map[fieldIndex]*resolver)}
	rsv.elementOrder = append(rsv.elementOrder, "1")

	// submap for Posts
	postsRsv := newResolver()
	// Post 1
	postElem1 := reflect.New(m.SubMaps[2].Typ).Elem()
	postElem1.FieldByName("ID").SetInt(101)
	postElem1.FieldByName("Title").SetString("First Post")
	postsRsv.elements["101"] = &element{v: postElem1}
	postsRsv.elementOrder = append(postsRsv.elementOrder, "101")

	rsv.elements["1"].subMaps[fieldIndex(2)] = postsRsv

	var user User
	dst := reflect.ValueOf(&user)

	err = setDst(m, dst, rsv)
	if err != nil {
		t.Fatalf("error setting destination: %s", err)
	}

	if user.ID != 1 || user.Name != "John Doe" {
		t.Errorf("expected user to be {ID:1 Name:\"John Doe\" ...}, but got %+v", user)
	}
	if len(user.Posts) != 1 {
		t.Fatalf("expected user to have 1 post, but got %d", len(user.Posts))
	}
	if user.Posts[0].ID != 101 || user.Posts[0].Title != "First Post" {
		t.Errorf("expected post 1 to be {ID:101 Title:\"First Post\"}, but got %+v", user.Posts[0])
	}
}
