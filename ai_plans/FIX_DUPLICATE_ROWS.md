# Problem: Incorrect De-duplication When Mapping to Basic Slices

## Summary
The `carta` library is designed to de-duplicate entities when mapping SQL rows to slices of structs (e.g., `[]User`). This is achieved by generating a unique ID for each entity based on the content of its primary key columns. This behavior is correct for handling `JOIN`s where a single entity might appear across multiple rows.

However, this same logic is incorrectly applied when the mapping destination is a slice of a basic type (e.g., `[]string`, `[]int`). In this scenario, rows with duplicate values are treated as the same entity and are de-duplicated, which is incorrect. The desired behavior is to preserve every row from the result set, including duplicates.

This issue is the root cause for the following problems:
1.  The `if m.IsBasic` code path in `load.go` lacks test coverage because no tests exist for mapping to basic slices.
2.  Attempts to write such tests lead to infinite loops and incorrect behavior because the column allocation and unique ID generation logic are not designed to handle this case.

## Proposed Solution
The solution is to create a distinct execution path for "basic mappers" (`m.IsBasic == true`) that ensures every row is treated as a unique element.

This will be accomplished in two main steps:

### 1. Fix Column Allocation (`allocateColumns`)
The logic will be modified to enforce a clear rule for basic slices: the source SQL query must return **exactly one column**.

-   If `m.IsBasic` is true, the function will bypass the existing name-matching logic.
-   It will validate that only one column is present in the query result.
-   This single column will be assigned as the `PresentColumn` for the mapper.
-   If more than one column is found, the function will return an error to prevent ambiguity.

### 2. Fix Unique ID Generation (`loadRow`)
The logic will be modified to generate a unique ID based on the row's position rather than its content.

-   If `m.IsBasic` is true, the call to `getUniqueId(row, m)` will be bypassed.
-   A new, position-based unique ID will be generated for each row (e.g., using a simple counter that increments with each row processed).
-   This ensures that every row, regardless of its content, is treated as a distinct element to be added to the destination slice.

This approach preserves the existing, correct behavior for struct mapping while introducing a new, robust path for handling basic slices correctly.

## Plan
1.  **Modify `column.go`**: Update the `allocateColumns` function to implement the single-column rule for basic mappers.
2.  **Modify `load.go`**: Update the `loadRow` function to use a position-based counter for unique ID generation when `m.IsBasic` is true.
3.  **Add Tests**: Create a new test case in `mapper_test.go` that maps a query result to a slice of a basic type (e.g., `[]string`) to validate the fix and provide coverage for the `m.IsBasic` code path.