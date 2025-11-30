package sample

import "context"

// User represents a user in the system.
type User struct {
	ID    int64
	Name  string
	Email string
}

// UserRepository defines the interface for user storage.
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	Create(ctx context.Context, user *User) error
}

// GetUserByID retrieves a user by their ID.
func GetUserByID(ctx context.Context, repo UserRepository, id int64) (*User, error) {
	return repo.GetByID(ctx, id)
}

// userService is an internal service for user operations.
type userService struct {
	repo UserRepository
}

// NewUserService creates a new user service.
func NewUserService(repo UserRepository) *userService {
	return &userService{repo: repo}
}

// GetUser retrieves a user using the service.
func (s *userService) GetUser(ctx context.Context, id int64) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

// MaxUsers is the maximum number of users allowed.
const MaxUsers = 1000

// DefaultPageSize is the default pagination size.
var DefaultPageSize = 20
