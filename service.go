package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

)

// Service is a simple CRUD interface for users
type Service interface {
	PostUser  (ctx context.Context, u User) error
	GetUser   (ctx context.Context, id string) (User, error)
	PutUser   (ctx context.Context, id string, u User) error
	PatchUser (ctx context.Context, id string, u User) error
	DeleteUser(ctx context.Context, id string) error
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

// UserModel represents the model of a user
type UserModel struct {
	gorm.Model
	FirstName    string
	LastName     string
	Username	 string
	Email        string  `gorm:"type:varchar(100);unique_index"`
	Password 	 string
	Role         string  `gorm:"size:255"`
}

// errors
var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type service struct {
	db gorm.DB
}

func NewService(db gorm.DB) Service {
	return &service{db}
}

/**
 * SETUP STORE
 */
func (s *service) PostUser(ctx context.Context, u User) error {
	// POST = create, don't overwrite
	s.db.Create(&User{u.FirstName, u.LastName, u.Username, u.Password, u.Email, u.Role})

	return nil
}

func (s *service) GetUser(ctx context.Context, id string) (User, error) {
	// GET = if found, return user
	var user User
	s.db.First(&user, id)

	return user, nil
}

func (s *service) PutUser(ctx context.Context, id string, u User) error {
	// PUT = create or update

	// find in orm
	var user User
	s.db.First(&user, id)

	fmt.Println("USER IN DB",user)

	// if exists update
	// else create

	return nil
}

func (s *service) PatchUser(ctx context.Context, id string, u User) error {
	// PATCH = update existing, don't create
	// find in orm
	var user User
	s.db.First(&user, id)

	fmt.Println("USER IN DB",user)

	return nil
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	// DELETE = if found, delete user
	var user User
	s.db.Delete(&user)

	return nil
}

//func (s *inmemService) PutUser(ctx context.Context, username string, u User) error {
//	if username != u.Username {
//		return ErrInconsistentIDs
//	}
//	s.mtx.Lock()
//	defer s.mtx.Unlock()
//	s.m[username] = u // PUT = create or update
//	return nil
//}

//func (s *inmemService) PatchUser(ctx context.Context, username string, u User) error {
//	if u.Username != "" && username != u.Username {
//		return ErrInconsistentIDs
//	}
//
//	s.mtx.Lock()
//	defer s.mtx.Unlock()
//
//	existing, ok := s.m[username]
//	if !ok {
//		return ErrNotFound // PATCH = update existing, don't create
//	}
//
//	// fields that can be modified
//	if u.FirstName != "" {
//		existing.FirstName = u.FirstName
//	}
//
//	if u.LastName != "" {
//		existing.LastName = u.LastName
//	}
//
//	if u.Email != "" {
//		existing.Email= u.Email
//	}
//
//	if u.Role!= "" {
//		existing.Role= u.Role
//	}
//
//	s.m[username] = existing
//	return nil
//}

//func (s *inmemService) DeleteUser(ctx context.Context, username string) error {
//	s.mtx.Lock()
//	defer s.mtx.Unlock()
//	if _, ok := s.m[username]; !ok {
//		return ErrNotFound
//	}
//	delete(s.m, username)
//	return nil
//}