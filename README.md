# Carta
Dead simple SQL data mapper for complex Go structs. 

Load SQL data onto Go structs while keeping track of has-one and has-many relationships

## Examples 
Using carta is very simple. All you need to do is: 
```
// 1) Run your query
if rows, err = sqlDB.Query(blogQuery); err != nil {
	// error
}

// 2) Instantiate a slice(or struct) which you want to populate 
blogs := []Blog{}

// 3) Map the SQL rows to your slice
carta.Map(rows, &blogs)
```

Assume that in above exmple, we are using a schema containing has-one and has-many relationships:

![schema](https://i.ibb.co/SPH3zhQ/Schema.png)

And here is our SQL query along with the corresponging Go struct:
```
select
       b.id,
       b.title,
       p.id        as  posts_id,         
       p.name      as  posts_name,
       a.id        as  author_id,      
       a.username  as  author_username
from blog b
       left outer join author a    on  b.author_id = a.id
       left outer join post p      on  b.id = p.blog_id
```

```
type Blog struct {
        Id     int    `db:"id"`
        Title  string `db:"title"`
        Posts  []Post `carta:"posts"`
        Author Author `carta:"author"`
}
type Post struct {
        Id   int    `db:"id"`
        Name string `db:"name"`
}
type Author struct {
        Id       int    `db:"id"`
        Username string `db:"username"`
}
```
Carta will map the SQL rows while keeping track of those relationships. 

Results: 
```
rows:
id | title | posts_id | posts_name | author_id | author_username
1  | Foo   | 1        | Bar        | 1         | John
1  | Foo   | 2        | Baz        | 1         | John
2  | Egg   | 3        | Beacon     | 2         | Ed

blogs:
[{
	"id": 1,
	"title": "Foo",
	"author": {
		"id": 1,
		"username": "John"
	},
	"posts": [{
			"id": 1,
			"name": "Bar"
		}, {
			"id": 2,
			"name": "Baz"
		}]
}, {
	"id": 2,
	"title": "Egg",
	"author": {
		"id": 2,
		"username": "Ed"
	},
	"posts": [{
			"id": 3,
			"name": "Beacon"
		}]
}]
```


## Comparison to Related Projects

#### GORM
Carta is NOT an an object-relational mapper(ORM). Read more in [Approach](#Approach)

#### sqlx
Sqlx does not track has-many relationships when mapping SQL data. This works fine when all your relationships are at most has-one (Blog has one Author) ie, each SQL row corresponds to one struct. However, handling has-many relationships (Blog has many Posts), requires  running many queries or running manual post-processing of the result. Carta handles these complexities automatically.

## Guide

### Column and Field Names

Carta matches your SQL columns with the corresponding struct fields.

#### Basic Fields
For basic types (int, string, etc.), use the `db` tag to specify the column name. If no tag is provided, Carta will use the snake_case version of the field name.

```go
type User struct {
	// Tag is specified, so "user_id" is the expected column name.
	Id int `db:"user_id"`

	// No tag, so "user_name" is the expected column name.
	UserName string
}
```

#### Associations (Nested Structs)
For nested structs (has-one or has-many relationships), use the `carta` tag to define a prefix for the nested struct's columns. This allows you to reuse struct definitions. The SQL query must then use aliases for the columns of the joined table.

**Example:**
```go
type Blog struct {
    Id     int    `db:"id"`
    Title  string `db:"title"`
    Author Author `carta:"author"` // "author" is the prefix
}

type Author struct {
    Id       int    `db:"id"`       // Maps to "author_id"
    Username string `db:"username"` // Maps to "author_username"
}
```

**Corresponding SQL Query:**
```sql
select
    b.id,
    b.title,
    a.id as author_id,
    a.username as author_username
from blogs b
left join authors a on b.author_id = a.id
```

This design promotes struct reusability. The `Author` struct can be used on its own to map to a query like `select id, username from authors` or as a nested struct within `Blog` as shown above.

### Data Types and Relationships

Any primative types, time.Time, protobuf Timestamp, and sql.NullX can be loaded with Carta.
These types are one-to-one mapped with your SQL columns

To define more complex SQL relationships use slices and structs as in example below:

```
type Blog struct {
	BlogId int  // Will map directly with "blog_id" column 

	// If your SQL data can be "null", use pointers or sql.NullX
	AuthorId  *int
	CreatedOn *timestamp.Timestamp // protobuf timestamp
	UpdatedOn *time.Time
	SonsorId  sql.NullInt64

	// To define has-one relationship, use nested structs 
	// or pointer to a struct
	Author *Author `carta:"author"`

	// To define has-many relationship, use slices
	// options include: *[]*Post, []*Post, *[]Post, []Post
	Posts []*Post `carta:"posts"`

	// If your has-many relationship corresponds to one column,
	// you can use a slice of a settable type
	TagIds     []int           `db:"tag_id"`
	CommentIds []sql.NullInt64 `db:"comment_id"`
}
```

### Drivers 

Recommended driver for Postgres is [lib/pg](https://github.com/lib/pq), for MySql use [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql).

When using MySql, carta expects time data to arrive in time.Time format. Therefore, make sure to add "parseTime=true" in your connection string, when using DATE and DATETIME types.

Other types, such as TIME, will will be converted from plain text in future versions of Carta.

## Installation 
```
go get -u github.com/hackafterdark/carta
```


## Important Notes 

Carta removes any duplicate rows. This is a side effect of the data mapping as it is unclear which object to instantiate if the same data arrives more than once.
If this is not a desired outcome, you should include a uniquely identifiable columns in your query and the corresponding fields in your structs.
 
To prevent relatively expensive reflect operations, carta caches the structure of your struct using the column mames of your query response as well as the type of your struct. 

## Approach
Carta adopts the "database mapping" approach (described in Martin Fowler's [book](https://books.google.com/books?id=FyWZt5DdvFkC&lpg=PA1&dq=Patterns%20of%20Enterprise%20Application%20Architecture%20by%20Martin%20Fowler&pg=PT187#v=onepage&q=active%20record&f=false)) which is useful among organizations with strict code review processes.

Carta is not an object-relational mapper(ORM). With large and complex datasets, using ORMs becomes restrictive and reduces performance when working with complex queries. 

### License
Apache License
