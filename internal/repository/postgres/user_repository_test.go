package postgres

import (
	"testing"

	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	db.AutoMigrate(&entity.User{})

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &entity.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	err := repo.Create(user)

	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &entity.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	repo.Create(user)

	found, err := repo.FindByID(user.ID)

	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, user.Email, found.Email)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &entity.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	repo.Create(user)

	found, err := repo.FindByEmail("test@example.com")

	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, user.Email, found.Email)
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &entity.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	repo.Create(user)

	name := "Updated Name"
	user.Name = &name
	err := repo.Update(user)

	assert.NoError(t, err)

	updated, _ := repo.FindByID(user.ID)
	assert.NotNil(t, updated.Name)
	assert.Equal(t, "Updated Name", *updated.Name)
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &entity.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	repo.Create(user)

	err := repo.Delete(user.ID)
	assert.NoError(t, err)

	found, err := repo.FindByID(user.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}
