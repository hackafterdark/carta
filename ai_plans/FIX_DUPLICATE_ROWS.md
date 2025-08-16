# Plan: Fix Incorrect De-duplication for Basic Slices

## 1. Problem Summary
The `carta` library was incorrectly de-duplicating rows when mapping to a slice of a basic type (e.g., `[]string`). The logic, designed to merge `JOIN`ed rows for slices of structs, was misapplied, causing data loss. This also meant the `m.IsBasic` code path was entirely untested.

The goal was to modify the library to correctly preserve all rows, including duplicates, when mapping to a basic slice, and to add the necessary test coverage.

## 2. Evolution of the Solution

The final solution was reached through an iterative process of implementation and refinement based on code review feedback.

### Initial Implementation
The first version of the fix introduced two key changes:
1.  **Position-Based Unique IDs:** In `load.go`, when `m.IsBasic` is true, `loadRow` generates a per-result-set unique ID (e.g., "row-0", "row-1") from the zero-based row index, ensuring every row is treated as unique within that mapping operation.
2.  **Single-Column Rule:** In `column.go`, the `allocateColumns` function was updated to enforce a strict rule: if the destination is a basic slice, the SQL query must return **exactly one column**. This prevents ambiguity.

### Refinements from Code Review
Feedback from code review prompted several improvements:
-   **Performance:** In `load.go`, `fmt.Sprintf` was replaced with `strconv.Itoa` for generating the position-based unique ID.
-   **Idiomatic Go:** Error creation now uses `fmt.Errorf` instead of `errors.New(fmt.Sprintf(...))`.
-   **Clearer errors:** The single-column rule error message includes the actual number of columns found.
-   **Test coverage:** A negative test was added to `mapper_test.go` to ensure the single-column rule correctly returns an error.

### Final Fix: Handling Nested Basic Mappers
The most critical refinement came from identifying a flaw in the single-column rule: it did not correctly handle **nested** basic slices (e.g., a struct field like `Tags []string`). The initial logic would have incorrectly failed if other columns for the parent struct were present.

The final patch corrected this by making the logic in `allocateColumns` more nuanced:
-   **For top-level basic slices** (`len(m.AncestorNames) == 0`), the result set must contain exactly one projected column (as labeled by the driver after alias resolution). Expressions are allowed if aliased to a single column.
-   **For nested basic slices**, the function now searches the remaining columns for exactly one that matches the ancestor-qualified name (e.g., `tags`). It returns an error if zero or more than one match is found.

This final change ensures the logic is robust for both top-level and nested use cases.

## 3. Summary of Changes Executed
1.  **Modified `load.go`**:
    -   Updated `loadRow` to accept a `rowCount` parameter and propagate it to nested mappers.
    -   For `m.IsBasic`, generate a per-row unique ID from `rowCount` to preserve duplicates (applies to nested basics as well).
    -   Refactored error handling and string formatting based on code review feedback.
2.  **Modified `column.go`**:
    -   Updated `allocateColumns` to differentiate between top-level and nested basic mappers, enforcing the correct single-column matching rule for each.
    -   Improved the error message to be more descriptive.
3.  **Modified `mapper.go`**:
    -   Corrected the logic in `determineFieldsNames` to properly handle casing in `carta` tags, ensuring ancestor names are generated correctly.
4.  **Added Tests to `mapper_test.go`**:
    -   Top-level `[]string`: verifies duplicates are preserved.
    -   Top-level `[]string` (negative): multi-column queries produce an error.
    -   Nested `PostWithTags.Tags []string`: verifies correct column matching and mapping.
    -   Nested (negative): zero or multiple matching columns produce an error.
5.  **Updated Documentation**:
    -   Updated `README.md` to clarify the difference in de-duplication behavior.
    -   Created `DESIGN_PHILOSOPHIES.md` to document the "fail-fast" error handling approach.