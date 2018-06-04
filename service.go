package users

import (
	"context"
	"errors"
	"sync"
)

// Service is a simple CRUD interface for users
type Service interface {
	PostUser(ctx context.Context, u User) error
	GetUser(ctx context.Context, username string) (User, error)
	PutUser(ctx context.Context, username string, u User) error
	PatchUser(ctx context.Context, username string, u User) error
	DeleteUser(ctx context.Context, username string) error
}

// User represents a single user
type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type inmemService struct {
	mtx sync.RWMutex
	m   map[string]User
}

func NewInmemService() Service {
	return &inmemService{
		m: map[string]User{},
	}
}

func (s *inmemService) PostUser(ctx context.Context, u User) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	// use username for key
	if _, ok := s.m[u.Username]; ok {
		return ErrAlreadyExists // POST = create, don't overwrite
	}
	s.m[u.Username] = u
	return nil
}

func (s *inmemService) GetUser(ctx context.Context, username string) (User, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	p, ok := s.m[username]
	if !ok {
		return User{}, ErrNotFound
	}
	return p, nil
}

func (s *inmemService) PutUser(ctx context.Context, username string, u User) error {
	if username != u.Username {
		return ErrInconsistentIDs
	}
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.m[username] = u // PUT = create or update
	return nil
}

func (s *inmemService) PatchUser(ctx context.Context, username string, u User) error {
	if u.Username != "" && username != u.Username {
		return ErrInconsistentIDs
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	existing, ok := s.m[username]
	if !ok {
		return ErrNotFound // PATCH = update existing, don't create
	}

	if u.FirstName != "" {
		existing.FirstName = u.FirstName
	}

	s.m[username] = existing
	return nil
}

func (s *inmemService) DeleteUser(ctx context.Context, username string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[username]; !ok {
		return ErrNotFound
	}
	delete(s.m, username)
	return nil
}