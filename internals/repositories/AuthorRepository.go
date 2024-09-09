package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/mongoVsGorm_GO/internals/types"
)

type AuthorRepository interface {
	CreateAuthor(ctx context.Context, name string, bio string, email string, dateOfBirth *time.Time) (uuid.UUID, error)
	GetAuthor(ctx context.Context, id uuid.UUID) (types.Author, error)
	ListAuthors(ctx context.Context) ([]types.Author, error)
	DeleteAuthor(ctx context.Context, id uuid.UUID) error
	UpdateAuthor(ctx context.Context, id uuid.UUID, name string, bio string, email string, dateOfBirth *time.Time) error
	GetAuthorsByBirthdateRange(ctx context.Context, startDate, endDate time.Time) ([]types.Author, error)
}
