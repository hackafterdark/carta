# Test Plan for Carta Library

## 1. Introduction and Goals

The primary goal of this test plan is to ensure the correctness and robustness of the `carta` library's functionality. The focus will be on writing unit tests that are easy to run and maintain, with minimal external dependencies. This plan outlines the strategy and specific test cases for testing the library's core features.

## 2. Testing Approach

We will primarily use **unit tests** to validate the functionality of the `carta` library. This approach allows us to test the mapping logic in isolation, without the need for a real database connection.

## 3. Tooling

As requested, we will use the `go-sqlmock` library for mocking database interactions. This library allows us to simulate the behavior of a SQL driver and `sql.Rows` in our tests, which is ideal for testing the `Map` function and its related components.

## 4. Test Coverage

The tests will cover the following files and their respective functionalities:

*   `mapper.go`: Core mapping logic, struct and field analysis.
*   `load.go`: Row loading and data conversion.
*   `setter.go`: Setting data to destination structs.
*   `column.go`: Column name to field mapping and snake case conversion.
*   `resolver.go`: Object resolution and deduplication.
*   `cache.go`: Caching of mappers.

## 5. Test Cases

### 5.1. `mapper.go`

*   **`TestMap` function:**
    *   Test mapping to a single struct.
    *   Test mapping to a slice of structs.
    *   Test mapping to a slice of pointers to structs.
    *   Test mapping with nested structs (one-to-one).
    *   Test mapping with nested slices of structs (one-to-many).
    *   Test mapping with `db` tags.
    *   Test mapping with different data types (int, string, bool, float, time.Time).
    *   Test mapping with `sql.Null` types.
    *   Test error handling for invalid destination types (e.g., not a pointer).
    *   Test caching by calling `Map` multiple times with the same query and destination type.

### 5.2. `load.go`

*   **`TestLoadRow` function:**
    *   Test loading a single row of data.
    *   Test loading multiple rows of data.
    *   Test data type conversions (e.g., string from DB to int in struct).
    *   Test handling of `NULL` values from the database.
    *   Test error handling for data conversion errors.

### 5.3. `column.go`

*   **`TestToSnakeCase` function:**
    *   Test conversion of `CamelCase` to `snake_case`.
    *   Test conversion of `camelCase` to `snake_case`.
    *   Test strings that are already in `snake_case`.
    *   Test strings with numbers.
*   **`TestAllocateColumns` function:**
    *   Test allocation of columns to struct fields.
    *   Test allocation with `db` tags.
    *   Test allocation with nested structs.

### 5.4. `cache.go`

*   **`TestCache`:**
    *   Test storing and loading a mapper from the cache.
    *   Test that the cache returns the correct mapper for a given set of columns and destination type.

## 6. Implementation Details

*   Each test file will be named `[filename]_test.go` (e.g., `mapper_test.go`).
*   Tests will be written using the standard Go `testing` package.
*   `go-sqlmock` will be used to create mock `*sql.DB` and `*sql.Rows` objects.
*   Assertions will be made to verify that the mapped data is correct and that errors are handled as expected.