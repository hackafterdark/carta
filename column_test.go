package carta

import (
	"reflect"
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "CamelCase to snake_case",
			input:    "CamelCase",
			expected: "camel_case",
		},
		{
			name:     "camelCase to snake_case",
			input:    "camelCase",
			expected: "camel_case",
		},
		{
			name:     "Already snake_case",
			input:    "snake_case",
			expected: "snake_case",
		},
		{
			name:     "String with numbers",
			input:    "UserID2",
			expected: "user_id2",
		},
		{
			name:     "Single word",
			input:    "user",
			expected: "user",
		},
		{
			name:     "String with space",
			input:    "user name",
			expected: "user_name",
		},
		{
			name:     "String with hyphen",
			input:    "user-name",
			expected: "user_name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := toSnakeCase(tc.input)
			if actual != tc.expected {
				t.Errorf("expected %s, but got %s", tc.expected, actual)
			}
		})
	}
}

func TestAllocateColumns(t *testing.T) {
	m, err := newMapper(reflect.TypeOf(&UserWithAddress{}))
	if err != nil {
		t.Fatalf("error creating new mapper: %s", err)
	}
	determineFieldsNames(m)

	columns := map[string]column{
		"ID": {
			name:        "ID",
			columnIndex: 0,
		},
		"Name": {
			name:        "Name",
			columnIndex: 1,
		},
		"Address_Street": {
			name:        "Address_Street",
			columnIndex: 2,
		},
		"Address_City": {
			name:        "Address_City",
			columnIndex: 3,
		},
	}
	err = allocateColumns(m, columns)
	if err != nil {
		t.Fatalf("error allocating columns: %s", err)
	}

	if len(m.PresentColumns) != 2 {
		t.Fatalf("expected 2 present columns for User, got %d", len(m.PresentColumns))
	}

	addressSubMap := m.SubMaps[2] // Address is the 3rd field (index 2)
	if len(addressSubMap.PresentColumns) != 2 {
		t.Fatalf("expected 2 present columns for Address, got %d", len(addressSubMap.PresentColumns))
	}

	if _, ok := addressSubMap.PresentColumns["Address_Street"]; !ok {
		t.Errorf("expected 'Address_Street' column to be present in submap")
	}
	if _, ok := addressSubMap.PresentColumns["Address_City"]; !ok {
		t.Errorf("expected 'Address_City' column to be present in submap")
	}
}

func TestAllocateColumnsBasic(t *testing.T) {
	m, err := newMapper(reflect.TypeOf(&[]string{}))
	if err != nil {
		t.Fatalf("error creating new mapper: %s", err)
	}
	determineFieldsNames(m)

	columns := map[string]column{
		"tag_name": {
			name:        "tag_name",
			columnIndex: 0,
		},
	}
	m.AncestorNames = []string{"tag_name"}
	err = allocateColumns(m, columns)
	if err != nil {
		t.Fatalf("error allocating columns: %s", err)
	}

	if len(m.PresentColumns) != 1 {
		t.Fatalf("expected 1 present column, got %d", len(m.PresentColumns))
	}
}

func TestGetColumnNameCandidatesWithCustomDelimiter(t *testing.T) {
	candidates := getColumnNameCandidates("field", []string{"parent"}, "->")
	expected := map[string]bool{
		"field":         true,
		"parent->field": true,
		"parent_field":  true,
	}

	for k := range expected {
		if _, ok := candidates[k]; !ok {
			t.Errorf("expected candidate %s to be present", k)
		}
	}
}

func TestAllocateColumnsDeeplyNested(t *testing.T) {
	m, err := newMapper(reflect.TypeOf(&BlogWithPosts{}))
	if err != nil {
		t.Fatalf("error creating new mapper: %s", err)
	}
	determineFieldsNames(m)

	columns := map[string]column{
		"id":                {name: "id", columnIndex: 0},
		"name":              {name: "name", columnIndex: 1},
		"posts_id":          {name: "posts_id", columnIndex: 2},
		"posts_title":       {name: "posts_title", columnIndex: 3},
		"posts_labels_id":   {name: "posts_labels_id", columnIndex: 4},
		"posts_labels_name": {name: "posts_labels_name", columnIndex: 5},
	}
	err = allocateColumns(m, columns)
	if err != nil {
		t.Fatalf("error allocating columns: %s", err)
	}

	if len(m.PresentColumns) != 2 {
		t.Fatalf("expected 2 present columns for Blog, got %d", len(m.PresentColumns))
	}

	postsSubMap := m.SubMaps[2] // Posts is the 3rd field (index 2)
	if len(postsSubMap.PresentColumns) != 2 {
		t.Fatalf("expected 2 present columns for Post, got %d", len(postsSubMap.PresentColumns))
	}

	labelsSubMap := postsSubMap.SubMaps[2] // Labels is the 3rd field (index 2)
	if len(labelsSubMap.PresentColumns) != 2 {
		t.Fatalf("expected 2 present columns for Label, got %d", len(labelsSubMap.PresentColumns))
	}

	if _, ok := labelsSubMap.PresentColumns["posts_labels_id"]; !ok {
		t.Errorf("expected 'posts_labels_id' column to be present in submap")
	}
	if _, ok := labelsSubMap.PresentColumns["posts_labels_name"]; !ok {
		t.Errorf("expected 'posts_labels_name' column to be present in submap")
	}
}
