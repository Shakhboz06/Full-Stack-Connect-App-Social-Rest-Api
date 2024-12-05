package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type UserStore struct {
	db *sql.DB
}

type Users struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  Password `json:"-"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
	RoleID    int64    `json:"role_id"`
	Role      Roles    `json:"role"`
}

type Password struct {
	text *string
	hash []byte
}

func (pass *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	pass.text = &text
	pass.hash = hash
	return nil
}
func (u *UserStore) Create(ctx context.Context, tx *sql.Tx, user *Users) error {
	query := ` 
		INSERT INTO Users (username, email, password, role_id)
		VALUES ($1, $2, $3, (SELECT id FROM roles WHERE name = $4)) RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	role := user.Role.Name
	if role == ""{
		role = "user"
	}

	err := tx.QueryRowContext(ctx, query, user.Username, user.Email, user.Password.hash, role).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}
	return nil
}

func (s *UserStore) GetUser(ctx context.Context, userId int64) (*Users, error) {
	query := `SELECT users.id, username, email, created_at, roles.* FROM users 
	JOIN roles ON (users.role_id = roles.id)
	WHERE users.id = $1 AND is_active = true
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user Users
	err := s.db.QueryRowContext(ctx, query, userId).
	Scan(&user.ID, &user.Username, 
		&user.Email, &user.CreatedAt, &user.Role.ID, &user.Role.Name, &user.Role.Level, &user.Role.Description)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}

	}

	return &user, nil
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *Users, token string, invitationExp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := s.createAndInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) createAndInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userID int64) error {
	query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Activation(ctx context.Context, token string) error {

	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		user, err := s.getUserFromInvitation(ctx, tx, token)

		if err != nil {
			return err
		}

		user.IsActive = true

		if err := s.getUserUpdate(ctx, tx, user); err != nil {
			return err
		}

		if err := s.deleteUserInvitation(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*Users, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.is_active
		FROM users u
		JOIN user_invitations ui ON u.id = ui.user_id
		WHERE ui.token = $1 AND ui.expiry > $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])

	user := &Users{}
	err := tx.QueryRowContext(ctx, query, hashedToken, time.Now()).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.IsActive)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStore) getUserUpdate(ctx context.Context, tx *sql.Tx, user *Users) error {
	query := `UPDATE users SET username = $1, email = $2, is_active = $3
	WHERE id = $4 `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM user_invitations WHERE user_id = $1`

	_, err := tx.ExecContext(ctx, query, userID)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Delete(ctx context.Context, userID int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.deleteUser(ctx, tx, userID); err != nil {
			return err
		}

		if err := s.deleteUserInvitation(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) deleteUser(ctx context.Context, tx *sql.Tx, userId int64) error {
	query := `DELETE FROM users u WHERE u.id = $1`

	_, err := tx.ExecContext(ctx, query, userId)


	_, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*Users, error) {
	query := `
		SELECT id, username, email, password, created_at, is_active
	 	FROM users WHERE email = $1 AND is_active = true
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	users := &Users{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(&users.ID, &users.Username, &users.Email, &users.Password.hash, &users.CreatedAt, &users.IsActive)

	if err != nil {
		return nil, err
	}

	return users, nil

}
