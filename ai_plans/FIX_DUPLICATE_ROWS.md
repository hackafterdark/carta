# Plan: Fix Incorrect De-duplication for Basic Slices

## 1. Problem Summary
The `carta` library was incorrectly de-duplicating rows when mapping to a slice of a basic type (e.g., `[]string`). The logic, designed to merge `JOIN`ed rows for slices of structs, was misapplied, causing data loss. This also meant the `m.IsBasic` code path was entirely untested.

The goal was to modify the library to correctly preserve all rows, including duplicates, when mapping to a basic slice, and to add the necessary test coverage.

## 2. Evolution of the Solution

The final solution was reached through an iterative process of implementation and refinement based on code review feedback.

### Initial Implementation
The first version of the fix introduced two key changes:
1.  **Position-Based Unique IDs:** In `load.go`, the `loadRow` function was modified. When `m.IsBasic` is true, it now generates a unique ID based on the row's position in the result set (e.g., "row-0", "row-1") instead of its content. This ensures every row is treated as a unique element.
2.  **Single-Column Rule:** In `column.go`, the `allocateColumns` function was updated to enforce a strict rule: if the destination is a basic slice, the SQL query must return **exactly one column**. This prevents ambiguity.

### Refinements from Code Review
Feedback from a code review (via Coderabbit) prompted several improvements:
-   **Performance:** In `load.go`, `fmt.Sprintf` was replaced with the more performant `strconv.Itoa` for generating the position-based unique ID.
-   **Idiomatic Go:** Error creation was changed from `errors.New(fmt.Sprintf(...))` to the more idiomatic `fmt.Errorf`.
-   **Clearer Errors:** The error message for the single-column rule was improved to include the actual number of columns found, aiding debugging.
-   **Test Coverage:** A negative test case was added to `mapper_test.go` to ensure the single-column rule correctly returns an error.

### Final Fix: Handling Nested Basic Mappers
The most critical refinement came from identifying a flaw in the single-column rule: it did not correctly handle **nested** basic slices (e.g., a struct field like `Tags []string`). The initial logic would have incorrectly failed if other columns for the parent struct were present.

The final patch corrected this by making the logic in `allocateColumns` more nuanced:
-   **For top-level basic slices** (`len(m.AncestorNames) == 0`), the query must still contain exactly one column overall.
-   **For nested basic slices**, the function now searches the remaining columns for exactly one that matches the ancestor-qualified name (e.g., `tags`). It returns an error if zero or more than one match is found.

This final change ensures the logic is robust for both top-level and nested use cases.

## 3. Summary of Changes Executed
1.  **Modified `load.go`**:
    -   Updated `loadRow` to accept a `rowCount` parameter.
    -   Implemented logic to generate a unique ID from `rowCount` when `m.IsBasic` is true.
    -   Refactored error handling and string formatting based on code review feedback.
2.  **Modified `column.go`**:
    -   Updated `allocateColumns` to differentiate between top-level and nested basic mappers, enforcing the correct single-column matching rule for each.
    -   Improved the error message to be more descriptive.
3.  **Modified `mapper.go`**:
    -   Corrected the logic in `determineFieldsNames` to properly handle casing in `carta` tags, ensuring ancestor names are generated correctly.
4.  **Added Tests to `mapper_test.go`**:
    -   Added a test for a top-level basic slice (`[]string`) to verify that duplicates are preserved.
    -   Added a negative test to ensure an error is returned for a multi-column query to a top-level basic slice.
    -   Added a test for a nested basic slice (`PostWithTags.Tags []string`) to verify correct mapping.
    -   Added negative tests to ensure errors are returned for nested basic slices with zero or multiple matching columns.
5.  **Updated Documentation**:
    -   Updated `README.md` to clarify the difference in de-duplication behavior.
    -   Created `DESIGN_PHILOSOPHIES.md` to document the "fail-fast" error handling approach.