package service

import (
	"server/internal/model"
	"server/internal/repository"

	"github.com/jmoiron/sqlx"
)

type MessageService struct {
	repo *repository.MessageRepository
}

func NewMessageService(db *sqlx.DB) *MessageService {
	return &MessageService{
		repo: repository.NewMessageRepository(db),
	}
}

func (s *MessageService) CreateMessage(msg string) (int64, error) {
	message := &model.Message{Msg: msg}
	return s.repo.Create(message)
}

func (s *MessageService) GetMessageByID(id int) (*model.Message, error) {
	return s.repo.GetByID(id)
}

func (s *MessageService) GetAllMessages() ([]model.Message, error) {
	return s.repo.GetAll()
}

func (s *MessageService) UpdateMessage(id int, msg string) (int64, error) {
	message := &model.Message{ID: id, Msg: msg}
	return s.repo.Update(message)
}

func (s *MessageService) DeleteMessage(id int) (int64, error) {
	return s.repo.Delete(id)
}
