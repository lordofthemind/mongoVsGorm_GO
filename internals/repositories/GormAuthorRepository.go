package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/mongoVsGorm_GO/internals/types"
	"gorm.io/gorm"
)

type GORMAuthorRepository struct {
	db *gorm.DB
}

func NewGORMAuthorRepository(db *gorm.DB) *GORMAuthorRepository {
	return &GORMAuthorRepository{db: db}
}

func (repo *GORMAuthorRepository) CreateAuthor(ctx context.Context, name string, bio string, email string, dateOfBirth *time.Time) (uuid.UUID, error) {
	id := uuid.New()
	author := types.Author{
		ID:          id,
		Name:        name,
		Bio:         bio,
		Email:       email,
		DateOfBirth: dateOfBirth,
	}
	if err := repo.db.WithContext(ctx).Create(&author).Error; err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (repo *GORMAuthorRepository) GetAuthor(ctx context.Context, id uuid.UUID) (types.Author, error) {
	var author types.Author
	if err := repo.db.WithContext(ctx).First(&author, "id = ?", id).Error; err != nil {
		return types.Author{}, err
	}
	return author, nil
}

func (repo *GORMAuthorRepository) ListAuthors(ctx context.Context) ([]types.Author, error) {
	var authors []types.Author
	if err := repo.db.WithContext(ctx).Find(&authors).Error; err != nil {
		return nil, err
	}
	return authors, nil
}

func (repo *GORMAuthorRepository) DeleteAuthor(ctx context.Context, id uuid.UUID) error {
	if err := repo.db.WithContext(ctx).Delete(&types.Author{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (repo *GORMAuthorRepository) UpdateAuthor(ctx context.Context, id uuid.UUID, name string, bio string, email string, dateOfBirth *time.Time) error {
	if err := repo.db.WithContext(ctx).Model(&types.Author{}).Where("id = ?", id).Updates(types.Author{
		Name:        name,
		Bio:         bio,
		Email:       email,
		DateOfBirth: dateOfBirth,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *GORMAuthorRepository) GetAuthorsByBirthdateRange(ctx context.Context, startDate, endDate time.Time) ([]types.Author, error) {
	var authors []types.Author
	if err := repo.db.WithContext(ctx).Where("date_of_birth BETWEEN ? AND ?", startDate, endDate).Find(&authors).Error; err != nil {
		return nil, err
	}
	return authors, nil
}
