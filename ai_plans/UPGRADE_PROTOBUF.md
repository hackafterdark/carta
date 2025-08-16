# Protobuf Upgrade Summary

This document outlines the changes made to upgrade the protobuf library and clarifies the use of `timestamppb.Timestamp` within the `carta` package.

## Summary of Changes

The primary goal was to resolve the use of a deprecated protobuf package, `github.com/golang/protobuf/ptypes`.

1.  **Updated Deprecated Import**: In `value/cell.go`, the deprecated package `github.com/golang/protobuf/ptypes` was replaced with the current, official package `google.golang.org/protobuf/types/known/timestamppb`.

2.  **Corrected Timestamp Conversion**: The `Timestamp()` method in `value/cell.go` was updated to manually convert a `time.Time` object to a `timestamppb.Timestamp`. This was necessary because the project's version of the `protobuf` library did not include the newer helper functions like `timestamppb.New()`.

3.  **Resolved Type Mismatch**: A panic occurred during testing because the `Timestamp()` method was returning a pointer (`*timestamppb.Timestamp`) while the mapping logic in `load.go` expected a value (`timestamppb.Timestamp`). The method was corrected to return a value, resolving the panic.

4.  **Updated Test Cases**: The test suite, specifically in `load_test.go`, was updated to use the new `timestamppb` package and reflect the changes in the `Timestamp()` method, ensuring all tests pass.

5.  **Dependency Management**: The `go.mod` and `go.sum` files were updated by running `go mod tidy` to ensure the new dependency was correctly managed and the project builds cleanly.

## Use of `timestamppb.Timestamp` in `carta`

The `carta` package uses `timestamppb.Timestamp` to ensure compatibility with Protocol Buffers (protobuf), which is a language-neutral, platform-neutral mechanism for serializing structured data.

### Why Not Just `time.Time`?

While `time.Time` is the standard for handling time within Go applications, `timestamppb.Timestamp` is the official, standardized protobuf representation for timestamps. This design choice provides several key benefits:

*   **Interoperability**: It allows for seamless data exchange with other services that use protobuf, such as gRPC APIs, without requiring custom conversion logic.
*   **Standardization**: It ensures that timestamp data is handled consistently across different programming languages and platforms.

### Compatibility with `database/sql`

`timestamppb.Timestamp` is **not** directly compatible with Go's `database/sql` package, which expects `time.Time` for SQL timestamp/datetime types.

The `carta` library bridges this gap by acting as a mapping layer:
1.  The SQL driver scans a timestamp from the database into a `time.Time` object.
2.  `carta`'s internal `value.Cell` stores this `time.Time` object.
3.  When mapping to a struct field, `carta` checks the field's type:
    *   If the field is `time.Time`, it sets the value directly.
    *   If the field is `timestamppb.Timestamp`, it calls the `cell.Timestamp()` method to perform the conversion before setting the value.

This allows developers the flexibility to use either `time.Time` for standard Go applications or `timestamppb.Timestamp` when interoperability with protobuf systems is required.