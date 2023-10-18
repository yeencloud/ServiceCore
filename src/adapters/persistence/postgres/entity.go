package postgres

import "github.com/google/uuid"

type ServiceDatabaseEntity struct {
	ID uuid.UUID

	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
}

func (entity *ServiceDatabaseEntity) GetID() uuid.UUID {
	return entity.ID
}
