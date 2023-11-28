package database

import (
	"github.com/google/uuid"
	"log"
)

func (a *DatabaseAdapter) Save(data *DummyOrm) (uuid.UUID, error) {
	if err := a.db.Create(data).Error; err != nil {
		log.Println("Can't create data : ", err)
		return uuid.Nil, err
	}

	return data.UserId, nil
}

func (a *DatabaseAdapter) GetByUUID(uuid *uuid.UUID) (DummyOrm, error) {
	var res DummyOrm
	if err := a.db.First(&res, "user_id = ?", uuid).Error; err != nil {
		log.Printf("Can't get data : %v", err)
		return res, err
	}

	return res, nil
}
