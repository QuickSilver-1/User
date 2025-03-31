package interfaces

import "user/internal/domain"

type ServerRepo interface {
	Start(domain.Port) error
	Shutdown()
}