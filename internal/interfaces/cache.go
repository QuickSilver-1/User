package interfaces

import "user/internal/domain"

// CacheRepo представляет интерфейс для работы с кэшом
type CacheRepo interface {
	// CreateKey создает ключ в кэше
	CreateKey(domain.Id, domain.User) error

	// GetText получает пользователя из кэша по идентификатору
	GetByKey(domain.Id) (*domain.User, error)
	
	// DelKey удаляет ключ из кэша
	DelKey(domain.Id) error
	
	// Close закрывает подключение
	Close() error
}
