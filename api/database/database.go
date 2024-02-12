package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)


type User struct {
	ID    int
	Name  string
	Email string
}
type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

const (
	userContextKey = "user"
)

var (
	ErrUserNotFound = errors.New("user not found in context")
)

var db *sql.DB

var ErrCommentNotFound = errors.New("comment not found")

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func InitDB(dataSourceName string) error {
	var err error
	DB, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	return nil
}

func AddComment(ctx context.Context, comment Comment) (int, error) {
	query := `
		INSERT INTO comments (post_id, user_id, content, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var commentID int
	err := db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
		comment.CreatedAt,
	).Scan(&commentID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert comment: %v", err)
	}

	return commentID, nil
}

func GetUserFromContext(ctx context.Context) (*User, error) {
	user, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func CreatePost(post Post) (*Post, error) {
	query := `
		INSERT INTO posts (user_id, title, content, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, title, content, created_at
	`

	var createdPost Post
	err := db.QueryRow(
		query,
		post.UserID,
		post.Title,
		post.Content,
		time.Now(),
	).Scan(
		&createdPost.ID,
		&createdPost.UserID,
		&createdPost.Title,
		&createdPost.Content,
		&createdPost.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create post: %v", err)
	}

	return &createdPost, nil
}

func GetTags(ctx context.Context) ([]Tag, error) {
	query := `
		SELECT id, name
		FROM tags
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %v", err)
	}

	return tags, nil
}

func AddTag(ctx context.Context, tag Tag) (int, error) {
	query := `
		INSERT INTO tags (name)
		VALUES ($1)
		RETURNING id
	`

	var tagID int
	err := db.QueryRowContext(
		ctx,
		query,
		tag.Name,
	).Scan(&tagID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert tag: %v", err)
	}

	return tagID, nil
}

func GetCommentsForPost(ctx context.Context, postID int) ([]Comment, error) {
	query := `
		SELECT id, post_id, user_id, content, created_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at DESC
	`

	rows, err := db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %v", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %v", err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func DeleteComment(ctx context.Context, commentID int) error {
	query := `
		DELETE FROM comments
		WHERE id = $1
	`

	result, err := db.ExecContext(ctx, query, commentID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return ErrCommentNotFound
	}

	return nil
}

func GetPaginatedPosts(ctx context.Context, page, pageSize, postID int) ([]Post, error) {
	query := `
        SELECT id, user_id, title, content, created_at
        FROM posts
        ORDER BY created_at DESC
        OFFSET $1 LIMIT $2
    `

	offset := (page - 1) * pageSize

	rows, err := db.QueryContext(ctx, query, offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error in result rows: %v", err)
	}

	return posts, nil
}

func GetPaginatedComments(ctx context.Context, postID, offset, limit int) ([]Comment, error) {
	query := `
		SELECT id, post_id, user_id, content
		FROM comments
		WHERE post_id = $1
		ORDER BY id
		OFFSET $2 LIMIT $3
	`

	rows, err := db.QueryContext(ctx, query, postID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %v", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content); err != nil {
			return nil, fmt.Errorf("failed to scan comment row: %v", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading comments rows: %v", err)
	}

	return comments, nil
}
