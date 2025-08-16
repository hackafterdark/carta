# Delimiter Enhancement Plan

## 1. Introduction

The current convention for aliasing columns in a SQL query relies on an underscore (`_`) as a delimiter between the association prefix and the field name (e.g., `author_id`). This can lead to ambiguity when field names themselves contain underscores (e.g., `user_name`), making it difficult for a developer reading the SQL query to distinguish between a nested field and a simple field with an underscore in its name.

To improve clarity and provide more flexibility, this plan outlines a change to the default delimiter and introduces a mechanism for making the delimiter configurable.

## 2. Proposed Solution

We will change the default delimiter from `_` to `->`. This will make nested fields in SQL aliases more explicit (e.g., `author->id`).

Furthermore, we will enhance the `carta` tag to allow the delimiter to be configured on a per-association basis.

### New `carta` Tag Syntax:

The `carta` tag will now support a comma-separated format for options:

```go
`carta:"<prefix>[,<option1>=<value1>][,<option2>=<value2>]..."`
```

For this enhancement, we will introduce the `delimiter` option:

```go
type Blog struct {
  Id     int    `db:"id"`
  Title  string `db:"title"`
  // Uses the new default delimiter '->'
  Author Author `carta:"author"`
  // Overrides the default delimiter to be '_'
  Editor Author `carta:"editor,delimiter=_"`
}
```

-   If no `delimiter` is specified, the default `->` will be used.
-   If `delimiter` is specified, its value will be used.

### Example Queries:

**Using the default delimiter (`->`):**

```sql
select
    b.id,
    b.title,
    a.id as "author->id",
    a.username as "author->username"
from blogs b
left join authors a on b.author_id = a.id
```
*(Note: Quoting the alias may be necessary depending on the SQL dialect)*

**Using a custom delimiter (`_`):**

```sql
select
    b.id,
    b.title,
    e.id as "editor_id",
    e.username as "editor_username"
from blogs b
left join authors e on b.editor_id = e.id
```

## 3. Implementation Details

### `mapper.go` Modifications

1.  **Add `Delimiter` to `Mapper` struct:**
    -   A new field, `Delimiter string`, will be added to the `Mapper` struct. This will hold the delimiter for the sub-map.

2.  **Update `determineFieldsNames` function:**
    -   This function will be updated to parse the `carta` tag.
    -   It will split the tag string by `,`. The first part will be the prefix (`Field.Name`).
    -   Subsequent parts will be parsed as key-value pairs (e.g., `delimiter=_`).
    -   The parsed delimiter will be stored in the `subMap.Delimiter` field. If no delimiter is specified, it will default to `->`.

### `column.go` Modifications

1.  **Update `allocateColumns` function:**
    -   This function will need to be aware of the new `Delimiter` field. When it recursively calls itself for sub-maps, it will use the `subMap.Delimiter` when constructing the `AncestorNames`.

2.  **Update `getColumnNameCandidates` function:**
    -   This function will be modified to accept the delimiter as an argument.
    -   Instead of hardcoding `_`, it will use the provided delimiter to join the `ancestorNames` and the `fieldName`.

## 4. Testing Strategy

New unit tests will be added to `mapper_test.go` and `column_test.go` to validate the new functionality. The tests will cover:
-   Mapping a query using the default `->` delimiter.
-   Mapping a query using a custom delimiter specified in the `carta` tag (e.g., `_`, `-`).
-   Ensuring the old underscore behavior works if specified explicitly.
-   Correctly mapping deeply nested structs with different delimiters at each level.

## 5. Documentation

The `README.md` file will be updated to:
-   Explain the new default delimiter (`->`).
-   Document the new `carta` tag syntax for configuring the delimiter.
-   Provide clear examples of struct definitions and the corresponding SQL queries with the new delimiter syntax.