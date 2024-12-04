package store

import (
	"context"
	"database/sql"
	"time"
)


func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {}

func (m *MockUserStore) Create(ctx context.Context, tx *sql.Tx, u *Users) error {
	return nil
}

func (m *MockUserStore) GetUser(ctx context.Context, userID int64) (*Users, error) {
	return &Users{ID: userID}, nil
}

func (m *MockUserStore) GetByEmail(context.Context, string) (*Users, error) {
	return &Users{}, nil
}

func (m *MockUserStore) CreateAndInvite(ctx context.Context, user *Users, token string, exp time.Duration) error {
	return nil
}

func (m *MockUserStore) Activation(ctx context.Context, t string) error {
	return nil
}

func (m *MockUserStore) Delete(ctx context.Context, id int64) error {
	return nil
}