package interfaces

import "user/internal/domain"

type UserRepo interface {
	Create(domain.User) (*domain.Id, error)
	Get(domain.Id) (*domain.User, error)
	Update(domain.User) error
}