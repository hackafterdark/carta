# Bug Fix: `carta.Map` Incorrectly Handles Empty Slices and Nil Associations

## Summary of the Issue

When using `carta.Map` with a `LEFT JOIN` where the "many" side of a `has-many` relationship had no corresponding entries, the mapper would incorrectly produce a slice containing a single, zero-valued struct instead of an empty slice. For example, a user with no posts would be mapped as:

```json
{
  "user_id": 1,
  "posts": [
    {
      "post_id": null,
      "content": null
    }
  ]
}
```

Instead of the expected:

```json
{
  "user_id": 1,
  "posts": []
}
```

A related issue was discovered where `has-one` pointer associations were not being set to `nil` when the joined data was `NULL`, leading to an empty struct instead of a `null` value in the final JSON output.

The root cause was traced to the column allocation logic in `column.go`, where the sub-mapper for the nested struct was greedily and incorrectly claiming a non-`NULL` column from its parent (e.g., `user_id`), causing the `isNil` check to fail and an empty struct to be erroneously created.

## The Plan

1.  **Analyze Core Logic:** Review `mapper.go`, `load.go`, and `setter.go` to understand the mapping and data loading process.
2.  **Isolate the Flaw:** Investigate `column.go` to pinpoint the exact flaw in the `allocateColumns` function that led to incorrect column claims by sub-mappers.
3.  **Refactor Column Allocation:** Modify the `allocateColumns` function in `column.go` to iterate through a struct's fields first, then search for a matching column. This makes column claims unambiguous and prevents sub-mappers from claiming columns belonging to their parents.
4.  **Address Nil Associations:** Update the `setDst` function in `setter.go` to check if the resolver for a `has-one` pointer association is empty. If it is, the field is left as `nil`, ensuring it marshals to `null`.
5.  **Add Verification Tests:**
    *   Create a test case to verify that a `has-many` relationship with no results from a `LEFT JOIN` correctly produces an empty slice (`[]`).
    *   Create a test case to verify that a `has-one` pointer relationship with no results from a `LEFT JOIN` correctly produces a `nil` value.
6.  **Implement and Verify:** Execute the plan by applying the code changes and running the new tests to confirm the fixes.
7.  **Cleanup:** Remove the temporary test cases to keep the test suite clean.

## The Result

The plan was executed successfully.

*   The column allocation logic in `column.go` was refactored to be more precise, resolving the primary issue of incorrectly generated empty structs in slices.
*   The setter logic in `setter.go` was updated to correctly handle `nil` pointer associations.
*   All tests passed, confirming that both the empty slice and `nil` association issues have been fixed. The `carta` package now behaves as expected in these scenarios.