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
		user := &User{
			Username: "Theo",
			Email:    "drumandbassbob@gmail.com",
		}

		err := user.PasswordHash.Set("Password")
		assert.NoError(t, err)

		tx, err := db.Begin()
		assert.NoError(t, err)
		defer tx.Rollback()

		err = userStore.CreateUser(tx, user)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID, "ID should be populated by RETURNING clause")
		assert.NotZero(t, user.CreatedAt)
		assert.NotZero(t, user.UpdatedAt)

		tx.Commit()

		query := `
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users 
		WHERE username = $1;`

		var dbUser User
		err = db.QueryRow(query, user.Username).Scan(&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.PasswordHash.hash, &dbUser.CreatedAt, &dbUser.UpdatedAt)
		assert.NoError(t, err)

		assert.Equal(t, user.ID, dbUser.ID)
		assert.Equal(t, user.Username, dbUser.Username)
		assert.Equal(t, user.Email, dbUser.Email)
		assert.Equal(t, user.PasswordHash.hash, dbUser.PasswordHash.hash)
		assert.Equal(t, user.CreatedAt, dbUser.CreatedAt)
		assert.Equal(t, user.UpdatedAt, dbUser.UpdatedAt)
	})

	t.Run("fails to create user with duplicate username", func(t *testing.T) {
		user := &User{
			Username: "Theo",
			Email:    "example@gmail.com",
		}

		err := user.PasswordHash.Set("Password")
		assert.NoError(t, err)

		tx, err := db.Begin()
		assert.NoError(t, err)
		defer tx.Rollback()

		err = userStore.CreateUser(tx, user)
		assert.Error(t, err, "Expected to fail creating duplicate user")
	})

	t.Run("fails to create user with duplicate email", func(t *testing.T) {
		user := &User{
			Username: "Bob",
			Email:    "drumandbassbob@gmail.com",
		}

		err := user.PasswordHash.Set("Password")
		assert.NoError(t, err)

		tx, err := db.Begin()
		assert.NoError(t, err)
		defer tx.Rollback()

		err = userStore.CreateUser(tx, user)
		assert.Error(t, err, "Expected to fail creating duplicate user")
	})
}

func TestGetUserByUsername(t *testing.T) {
	db := SetupTestDB(t)
	TruncateTables(t, db)
	userStore := NewPostgresUserStore(db)

	user := &User{
		Username: "Theo",
		Email:    "drumandbassbob@gmail.com",
	}

	err := user.PasswordHash.Set("Password")
	assert.NoError(t, err)

	tx, err := db.Begin()
	assert.NoError(t, err)
	defer tx.Rollback()

	err = userStore.CreateUser(tx, user)
	assert.NoError(t, err)
	tx.Commit()

	t.Run("get existing user by username", func(t *testing.T) {
		dbUser, err := userStore.GetUserByUsername(user.Username)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, dbUser.ID)
		assert.Equal(t, user.Username, dbUser.Username)
		assert.Equal(t, user.Email, dbUser.Email)
		assert.Equal(t, user.PasswordHash.hash, dbUser.PasswordHash.hash)
		assert.Equal(t, user.CreatedAt, dbUser.CreatedAt)
		assert.Equal(t, user.UpdatedAt, dbUser.UpdatedAt)
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

	user := &User{
		Username: "Theo",
		Email:    "drumandbassbob@gmail.com",
	}

	err := user.PasswordHash.Set("Password")
	assert.NoError(t, err)

	tx, err := db.Begin()
	assert.NoError(t, err)
	err = userStore.CreateUser(tx, user)
	assert.NoError(t, err)
	err = tx.Commit()
	assert.NoError(t, err)

	token, err := tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	assert.NoError(t, err)

	t.Run("ensure sign in with valid token", func(t *testing.T) {
		dbUser, err := userStore.GetUserToken(tokens.ScopeAuth, token.Plaintext)
		assert.NoError(t, err)

		assert.Equal(t, user.ID, dbUser.ID)
		assert.Equal(t, user.Username, dbUser.Username)
		assert.Equal(t, user.Email, dbUser.Email)
		assert.Equal(t, user.PasswordHash.hash, dbUser.PasswordHash.hash)
		assert.Equal(t, user.CreatedAt, dbUser.CreatedAt)
		assert.Equal(t, user.UpdatedAt, dbUser.UpdatedAt)
	})

	t.Run("ensure failed sign in when token is outdated", func(t *testing.T) {
		query := `
		UPDATE tokens 
		SET expiry = $1 
		WHERE user_id = $2 and scope = $3;
		`

		var expiry time.Time = time.Now().Add(-1 * time.Second)
		_, err := db.Exec(query, expiry, user.ID, tokens.ScopeAuth)
		assert.NoError(t, err)

		dbUser, err := userStore.GetUserToken(tokens.ScopeAuth, token.Plaintext)
		assert.NoError(t, err)
		assert.Nil(t, dbUser)
	})
}
