package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Posts struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comment   []Comment `json:"comments"`
	User      Users      `json:"user"`
}

type PostswithMetadata struct {
	Posts
	CommentCount int `json:"comment_count"`
}
type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) GetUserFeed(ctx context.Context, UserID int64, pg PaginatedFeed)([]PostswithMetadata, error){
	query :=`
	SELECT 
		p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
		u.username,
		COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
		WHERE 
			f.user_id = $1 AND
			(
				COALESCE(p.tags, '{}') @> $2 OR
				(p.title ILIKE '%' || $3 || '%' OR p.content ILIKE '%' || $3 || '%')
			)
		GROUP BY p.id, u.username
		ORDER BY p.created_at `  + pg.Sort + ` 
		LIMIT $4 OFFSET $5 
		`
	

	
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, UserID, pq.Array(pg.Tags), pg.Search, pg.Limit, pg.Offset)

	if err != nil{
		return nil, err
	}

	defer rows.Close()

	var feed []PostswithMetadata

	for rows.Next(){
		var post PostswithMetadata
		err := rows.Scan(
			&post.ID, 
			&post.UserID, 
			&post.Title,
			&post.Content, 
			&post.CreatedAt, 
			&post.Version, 
			pq.Array(&post.Tags), 
			&post.User.Username, 
			&post.CommentCount)

		if err != nil{
			return nil, err
		}
		feed = append(feed, post)
	
	}
	return feed, nil
}	
func (s *PostsStore) Create(ctx context.Context, post *Posts) error {

	query := ` 
		INSERT INTO Posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserID, pq.Array(post.Tags)).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostsStore) GetbyID(ctx context.Context, postID int64) (*Posts, error) {
	query := `SELECT id, user_id, title, content, created_at, updated_at, tags, version
		FROM Posts 
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post Posts
	err := s.db.QueryRowContext(ctx, query, postID).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt, pq.Array(&post.Tags), &post.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}

	}

	return &post, nil
}

func (s *PostsStore) DeletebyID(ctx context.Context, postID int64) error {
	query := "DELETE FROM posts WHERE id = $1"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, postID)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostsStore) UpdatebyID(ctx context.Context, post *Posts) error {

	query := `
		UPDATE posts 
		SET 
		title = COALESCE($1, title), 
		content = COALESCE($2, content),
		version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version).Scan(&post.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}
