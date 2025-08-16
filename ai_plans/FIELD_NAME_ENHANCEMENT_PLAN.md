# Field Name Enhancement Plan

## 1. Introduction

The current implementation of `carta` requires developers to create specialized structs for associations by prefixing `db` tags to ensure uniqueness (e.g., `db:"author_id"`). This approach prevents struct reusability. An `Author` struct designed for a `Blog` association cannot be used for a standalone `authors` table query without modification.

This plan details the introduction of a new `carta` struct tag. This tag will define the prefix for an association, allowing the SQL query to handle column aliasing. This makes nested structs highly reusable, as their `db` tags can remain generic (e.g., `db:"id"`).

## 2. Proposed Solution

We will introduce a new struct tag, `carta`, to be used on fields that represent nested struct associations (has-one or has-many). The value of the `carta` tag will specify the prefix that the mapper should expect for the columns related to that association. The actual column aliasing will be handled in the SQL query.

This decouples the Go struct definition from the specific query structure, promoting reusability.

### Example:

**Corrected Struct Definitions:**

```go
type Blog struct {
    Id     int    `db:"id"`
    Title  string `db:"title"`
    Author Author `carta:"author"` // 'author' is the key for the association
}

type Author struct {
    Id       int    `db:"id"`
    Username string `db:"username"`
}
```

With this structure, the `Author` struct is now generic. It can be used in two scenarios:

1.  **Standalone Query:** The `Author` struct can be used directly.
    ```sql
    select id, username from authors;
    ```

2.  **Joined Query:** The `Author` struct can be used as a nested field within `Blog`. The SQL query must alias the columns for the author using the `carta` tag value as a prefix.
    ```sql
    select
        b.id,
        b.title,
        a.id        as author_id,
        a.username  as author_username
    from blog b
    left join author a on b.author_id = a.id;
    ```
The `carta` mapper will see the `carta:"author"` tag on the `Blog.Author` field and will look for columns named `author_id` and `author_username` to map to the `Author` struct's `id` and `username` fields.

## 3. Implementation Details

The primary changes will be confined to `mapper.go` to correctly process the new `carta` tag.

### `mapper.go` Modifications

1.  **Update `determineFieldsNames` function:**
    - This function inspects struct fields to prepare for mapping.
    - The logic will be updated: When processing a field that is a sub-map (a nested struct or slice of structs), it will first check for a `carta` tag.
    - If a `carta:"..."` tag is found, its value will be used as the `Field.Name` for that association. This name is then added to the `AncestorNames` slice for the sub-mapper.
    - If no `carta` tag is present, the existing logic (using the `db` tag or the field's name) will be used as a fallback to provide some backward compatibility.

2.  **Update `nameFromTag` function:**
    - A new `nameFromTag` function will be created to specifically look for the `carta` tag on struct fields that are associations. The existing `nameFromTag` (which looks for `db`) will be used for basic fields. We need to differentiate the logic for associations vs. basic fields.

### `column.go` Modifications

-   The `getColumnNameCandidates` function will not require changes. It already constructs candidate column names by prepending `AncestorNames` to a field's name (e.g., `ancestor_field`). Since `mapper.go` will now be populating `AncestorNames` with the value from the `carta` tag, `getColumnNameCandidates` will correctly generate the expected aliased column names (e.g., `author_id`, `author_username`).

## 4. Testing Strategy

New unit tests will be added to `mapper_test.go` to validate the new functionality. The tests will cover:
-   Mapping a joined query where column aliases match the `carta` tag prefix.
-   Confirming that a single struct definition can be used for both a joined query and a standalone query without modification.
-   Correctly mapping deeply nested structs using multiple `carta` tags.
-   Verifying the fallback behavior when no `carta` tag is present on an association.

## 5. Documentation

The `README.md` file will be updated to:
-   Clearly explain the new recommended approach using the `carta` tag for all associations.
-   Provide examples of the reusable struct definitions and the corresponding SQL queries with aliasing.
-   Update the "Column and Field Names" section to reflect the new convention.