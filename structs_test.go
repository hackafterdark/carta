package carta

// Label represents a label associated with a blog post.
type Label struct {
	ID   *int    `db:"id"`
	Name *string `db:"name"`
}

// PostWithLabels represents a blog post that has multiple labels.
type PostWithLabels struct {
	ID     int     `db:"id"`
	Title  string  `db:"title"`
	Labels []Label `carta:"labels"`
}

// BlogWithPosts represents a blog that has multiple posts.
type BlogWithPosts struct {
	ID    int              `db:"id"`
	Name  string           `db:"name"`
	Posts []PostWithLabels `carta:"posts"`
}
