package store

import (
	"context"
	"database/sql"
)

// Kind of a repository for multiple entities
type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, string) (Post, error)
		Delete(context.Context, string) error
		Update(context.Context, *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
	Comments interface {
		GetByPostID(context.Context, int64) ([]Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}
