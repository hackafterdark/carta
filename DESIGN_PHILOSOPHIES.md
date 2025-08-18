# Design Philosophies

This document includes information around design philosophies and decisions made to help document and illustrate scenarios one may encounter when using this package.

## Approach
Carta adopts the "database mapping" approach (described in Martin Fowler's [book](https://books.google.com/books?id=FyWZt5DdvFkC&lpg=PA1&dq=Patterns%20of%20Enterprise%20Application%20Architecture%20by%20Martin%20Fowler&pg=PT187#v=onepage&q=active%20record&f=false)) which is useful among organizations with strict code review processes.

## Comparison to Related Projects

#### GORM
Carta is NOT an object-relational mapper (ORM).

#### sqlx
Sqlx does not track has-many relationships when mapping SQL data. This works fine when all your relationships are at most has-one (Blog has one Author) ie, each SQL row corresponds to one struct. However, handling has-many relationships (Blog has many Posts), requires  running many queries or running manual post-processing of the result. Carta handles these complexities automatically.

## Protection vs. Graceful Handling

A core design principle of the `carta` mapper is to prioritize **user protection and clarity** over attempting a "graceful" but potentially incorrect guess. The library's guiding philosophy is to only proceed if the user's intent is perfectly clear. If there is any ambiguity in the mapping operation, `carta` will **fail fast** by returning an error, forcing the developer to be more explicit.

Making a guess might seem helpful, but it can hide serious, silent bugs. The following scenarios illustrate the balance between failing on ambiguous operations (Protection) and handling well-defined transformations (Graceful Handling).

---

### Scenario 1: Multi-column Query to a Basic Slice (Protection)

-   **Query:** `SELECT name, email FROM users`
-   **Destination:** `var data []string`
-   **Behavior:** `carta.Map` **returns an error immediately**: `carta: when mapping to a slice of a basic type, the query must return exactly one column (got 2)`.
-   **Why this is Protection:** The library has no way of knowing if the user intended to map the `name` or the `email` column. A "graceful" solution might be to pick the first column arbitrarily, but this could lead to the wrong data being silently loaded into the slice. By failing fast, `carta` forces the developer to write an unambiguous query (e.g., `SELECT name FROM users`), ensuring the result is guaranteed to be correct.

---

### Scenario 2: SQL `NULL` to a Non-nullable Go Field (Protection)

-   **Query:** `SELECT id, NULL AS name FROM users`
-   **Destination:** `var users []User` (where `User.Name` is a `string`)
-   **Behavior:** `carta.Map` **returns an error during scanning** (e.g., `carta: cannot load NULL into non-nullable type string for column name`).
-   **Why this is Protection:** A standard Go `string` cannot represent a `NULL` value. A "graceful" but incorrect solution would be to use the zero value (`""`), which is valid data and semantically different from "no data". This can cause subtle bugs in application logic. By failing, `carta` forces the developer to explicitly handle nullability in their Go struct by using a pointer (`*string`) or a nullable type (`sql.NullString`), making the code more robust and correct.

---

### Scenario 3: Merging `JOIN`ed Rows into Structs (Graceful Handling)

-   **Query:** `SELECT b.id, p.id FROM blogs b JOIN posts p ON b.id = p.blog_id`
-   **Destination:** `var blogs []BlogWithPosts`
-   **Behavior:** `carta` **gracefully handles** the fact that the same blog ID appears in multiple rows. It creates one `Blog` object and appends each unique `Post` to its `Posts` slice.
-   **Why this is Graceful:** This is the core purpose of the library. There is no ambiguity. The library uses the unique ID of the `Blog` (the `b.id` column) to understand that these rows all describe the same parent entity. This is a well-defined transformation, not a guess.

---

### Scenario 4: Default `_` Delimiter Support (Graceful Handling)

-   **Query:** `SELECT b.id, a.id AS "author_id" FROM blogs b JOIN authors a ON b.author_id = a.id`
-   **Destination:** `var blogs []BlogWithAuthor`
-   **Behavior:** `carta` **gracefully handles** the `author_id` column, correctly mapping it to the `Author` struct's `id` field, even though the default delimiter is `->`.
-   **Why this is Graceful:** This convenience was a deliberate design choice. Since SQL `SELECT` statements must have unambiguous column names (which you control with aliases), there is no risk of conflict with actual database field names that contain underscores. This allows for more natural-looking column names in queries without requiring an explicit `delimiter=_` option in the `carta` tag. If another delimiter is desired, it must be set explicitly.
