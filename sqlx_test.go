package carta

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setupSqlxDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	createAuthorTable := `
	CREATE TABLE author (
		id INTEGER PRIMARY KEY,
		username TEXT
	);`
	_, err = db.Exec(createAuthorTable)
	assert.NoError(t, err)

	createBlogTable := `
	CREATE TABLE blog (
		id INTEGER PRIMARY KEY,
		title TEXT,
		author_id INTEGER
	);`
	_, err = db.Exec(createBlogTable)
	assert.NoError(t, err)

	_, err = db.Exec("INSERT INTO author (id, username) VALUES (?, ?)", 1, "johndoe")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO blog (id, title, author_id) VALUES (?, ?, ?)", 1, "My First Post", 1)
	assert.NoError(t, err)

	return db
}

type SqlxAuthor struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
}

type SqlxBlog struct {
	Id     int        `db:"id"`
	Title  string     `db:"title"`
	Author SqlxAuthor `carta:"author"`
}

func TestMapx(t *testing.T) {
	db := setupSqlxDB(t)
	defer db.Close()

	query := `
	SELECT
		b.id,
		b.title,
		a.id AS "author->id",
		a.username AS "author->username"
	FROM blog b
	LEFT JOIN author a ON b.author_id = a.id`

	rows, err := db.Queryx(query)
	assert.NoError(t, err)

	var blogs []*SqlxBlog
	err = Mapx(rows, &blogs)
	assert.NoError(t, err)

	assert.Len(t, blogs, 1)
	assert.Equal(t, 1, blogs[0].Id)
	assert.Equal(t, "My First Post", blogs[0].Title)
	assert.Equal(t, 1, blogs[0].Author.Id)
	assert.Equal(t, "johndoe", blogs[0].Author.Username)
}
