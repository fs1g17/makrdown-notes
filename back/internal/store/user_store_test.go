package store

import (
	"markdown-notes/internal/tokens"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db := SetupTestDB(t)
	TruncateTables(t, db)
	userStore := NewPostgresUserStore(db)

	t.Run("creates user successfully", func(t *testing.T) {
		user := CreateTestUser(t, db, userStore, "Theo", "drumandbassbob@gmail.com", "Password")
		assert.NotZero(t, user.ID, "ID should be populated by RETURNING clause")
		assert.NotZero(t, user.CreatedAt)
		assert.NotZero(t, user.UpdatedAt)

		query := `
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users 
		WHERE username = $1;`

		var dbUser User
		err := db.QueryRow(query, user.Username).Scan(&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.PasswordHash.hash, &dbUser.CreatedAt, &dbUser.UpdatedAt)
		assert.NoError(t, err)

		CompareUsers(t, user, &dbUser)
	})

	t.Run("fails to create user with duplicate username", func(t *testing.T) {
		user := &User{
			Username: "Theo",
			Email:    "other@gmail.com",
		}
		user.PasswordHash.Set("Password")
		tx, _ := db.Begin()
		err := userStore.CreateUser(tx, user)
		assert.Error(t, err, "Expected to fail creating duplicate user")
	})

	t.Run("fails to create user with duplicate email", func(t *testing.T) {
		user := &User{
			Username: "Other",
			Email:    "drumandbassbob@gmail.com",
		}
		user.PasswordHash.Set("Password")
		tx, _ := db.Begin()
		err := userStore.CreateUser(tx, user)
		assert.Error(t, err, "Expected to fail creating duplicate user")
	})
}

func TestGetUserByUsername(t *testing.T) {
	db := SetupTestDB(t)
	TruncateTables(t, db)
	userStore := NewPostgresUserStore(db)

	user := CreateTestUser(t, db, userStore, "Theo", "drumandbassbob@gmail.com", "Password")

	t.Run("get existing user by username", func(t *testing.T) {
		dbUser, err := userStore.GetUserByUsername(user.Username)
		assert.NoError(t, err)
		CompareUsers(t, user, dbUser)
	})

	t.Run("returns nil when not found", func(t *testing.T) {
		dbUser, err := userStore.GetUserByUsername("nonExistent")
		assert.NoError(t, err)
		assert.Nil(t, dbUser)
	})
}

func TestGetUserToken(t *testing.T) {
	db := SetupTestDB(t)
	TruncateTables(t, db)
	userStore := NewPostgresUserStore(db)
	tokenStore := NewPostgresTokenStore(db)

	user := CreateTestUser(t, db, userStore, "Theo", "drumandbassbob@gmail.com", "Password")
	otherUser := CreateTestUser(t, db, userStore, "Theo2", "example@gmail.com", "Password")

	token, err := tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	assert.NoError(t, err)
	outdatedToken, err := tokenStore.CreateNewToken(otherUser.ID, -1*time.Second, tokens.ScopeAuth)
	assert.NoError(t, err)

	t.Run("ensure sign in with valid token", func(t *testing.T) {
		dbUser, err := userStore.GetUserToken(tokens.ScopeAuth, token.Plaintext)
		assert.NoError(t, err)
		CompareUsers(t, user, dbUser)
	})

	t.Run("ensure failed sign in when token is outdated", func(t *testing.T) {
		dbUser, err := userStore.GetUserToken(tokens.ScopeAuth, outdatedToken.Plaintext)
		assert.NoError(t, err)
		assert.Nil(t, dbUser)
	})
}
