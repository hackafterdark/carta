# Test Plan for Deeply Nested Structs in Carta

## 1. Introduction and Goals

This document outlines the test plan for verifying the functionality of the `carta` library when dealing with deeply nested structs. Building upon the existing `TESTS_PLAN.md`, this plan specifically targets the recursive capabilities of the mapper to handle "n+1" relationships, where a struct contains a slice of another struct, which in turn contains a slice of a third struct (e.g., `Blog` -> `Posts` -> `Labels`).

The primary goal is to ensure that `carta` can correctly map data from a flat `sql.Rows` result set to a complex, multi-level object graph.

## 2. Testing Approach

We will continue to use **unit tests** with the `go-sqlmock` library. This allows us to simulate database results without requiring a live database connection, ensuring our tests are fast, reliable, and focused on the mapping logic itself.

The tests will simulate a JOIN query that flattens a one-to-many-to-many relationship into a single result set. For example, a query that joins `blogs`, `posts`, and `labels` tables.

## 3. Test Scenarios

### 3.1. Test Structs

The following structs will be used to model the deeply nested relationship:

```go
type Label struct {
	ID   int
	Name string
}

type Post struct {
	ID     int
	Title  string
	Labels []*Label
}

type Blog struct {
	ID    int
	Name  string
	Posts []Post
}
```

### 3.2. `mapper.go` Test Cases

*   **`TestMapDeeplyNested` function:**
    *   Test mapping to a single `Blog` struct with multiple `Posts`, where each `Post` has multiple `Labels`.
    *   Test mapping to a slice of `Blog` structs (`[]*Blog`).
    *   Test mapping where a `Post` has no associated `Labels`. The `Labels` slice should be empty or nil.
    *   Verify that the object graph is correctly reconstructed, with the correct number of posts for each blog and the correct number of labels for each post.
    *   Test with `db` tags to ensure custom column names are respected at all levels of nesting.
    *   Test with `carta` tags to handle different delimiters for nested struct fields.
    *   Test error handling when the destination is not a pointer.

### 3.3. `column.go` Test Cases

*   **`TestAllocateColumnsDeeplyNested` function:**
    *   Test the `allocateColumns` function with column names that represent a deeply nested structure (e.g., `id`, `name`, `posts_id`, `posts_title`, `posts_labels_id`, `posts_labels_name`).
    *   Verify that columns are correctly allocated to the fields of the nested structs (`Blog`, `Post`, and `Label`).
    *   Test with different naming conventions (e.g., `CamelCase`, `snake_case`) for struct fields and ensure they are correctly mapped to snake_cased column names.
    *   Test the use of custom delimiters in `carta` tags to correctly parse column names for nested structs.

## 4. Implementation Details

*   Tests will be added to `mapper_test.go` and `column_test.go`.
*   `go-sqlmock` will be used to create a mock `*sql.Rows` object with the flattened data from the simulated JOIN query.
*   Assertions will be used to verify the structure and content of the resulting mapped structs. For example, we will check `len(blog.Posts)` and `len(blog.Posts[0].Labels)`.