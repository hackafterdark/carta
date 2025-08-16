# SQLX Support Plan

## 1. Introduction

This document outlines the plan to add support for the `github.com/jmoiron/sqlx` package to `carta`. The `sqlx` package is a popular extension to the standard `database/sql` package, and providing compatibility will make `carta` more versatile for Go developers.

Since `sqlx.Rows` is a wrapper around `sql.Rows`, we can add support by creating simple convenience functions that extract the underlying `sql.Rows` and pass it to our existing mapping logic.

## 2. Proposed Solution

We will introduce a new file, `sqlx.go`, to house the `sqlx`-specific compatibility functions. This approach keeps the new functionality isolated and maintains a clean separation of concerns between the core `database/sql` logic and the `sqlx` support.

### New Functions

We will create the following function in `sqlx.go`:

**`Mapx(rows *sqlx.Rows, dst interface{}) error`**

This function will serve as the entry point for mapping `sqlx.Rows`. It will be a lightweight wrapper around the existing `Map` function.

Example Implementation:

```go
package carta

import "github.com/jmoiron/sqlx"

// Mapx maps sqlx.Rows onto a struct or slice of structs.
// It is a convenience wrapper around the Map function.
func Mapx(rows *sqlx.Rows, dst interface{}) error {
	return Map(rows.Rows, dst)
}
```

## 3. Implementation Details

1.  **Create `sqlx.go`:**
    -   A new file named `sqlx.go` will be created in the root of the project.
    -   This file will contain the `Mapx` function.

2.  **Add `sqlx` dependency:**
    -   The `go.mod` file will be updated to include `github.com/jmoiron/sqlx` as a dependency.

## 4. Testing Strategy

A new test file, `sqlx_test.go`, will be created to validate the `sqlx` integration. The tests will ensure that `Mapx` works correctly for various mapping scenarios.

The testing strategy will include:
-   Setting up a test database (e.g., in-memory SQLite).
-   Using `sqlx` to execute queries and retrieve `*sqlx.Rows`.
-   Calling `Mapx` to map the results to structs.
-   Adding test cases for:
    -   Mapping to a single struct.
    -   Mapping to a slice of structs.
    -   Mapping structs with nested associations (`has-one` and `has-many`).
    -   Verifying that all fields are correctly populated.

## 5. Documentation

The `README.md` file will be updated to:
-   Mention the new support for `sqlx`.
-   Provide a clear example of how to use the `Mapx` function with `sqlx.Rows`.