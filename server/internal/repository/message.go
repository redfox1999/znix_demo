package repository

import (
	"server/internal/model"

	"github.com/jmoiron/sqlx"
)

type MessageRepository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (m *MessageRepository) Create(msg *model.Message) (int64, error) {
	query := "INSERT INTO message (msg) VALUES (?)"
	result, err := m.db.Exec(query, msg.Msg)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *MessageRepository) GetByID(id int) (*model.Message, error) {
	var msg model.Message
	query := "SELECT id, msg FROM message WHERE id = ?"
	err := m.db.Get(&msg, query, id)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (m *MessageRepository) GetAll() ([]model.Message, error) {
	var messages []model.Message
	query := "SELECT id, msg FROM message"
	err := m.db.Select(&messages, query)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (m *MessageRepository) Update(msg *model.Message) (int64, error) {
	query := "UPDATE message SET msg = ? WHERE id = ?"
	result, err := m.db.Exec(query, msg.Msg, msg.ID)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}

func (m *MessageRepository) Delete(id int) (int64, error) {
	query := "DELETE FROM message WHERE id = ?"
	result, err := m.db.Exec(query, id)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}
