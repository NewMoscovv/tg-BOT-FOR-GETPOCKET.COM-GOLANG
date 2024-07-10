package boltDB

import (
	"First-TgBot-On-GO/pkg/repository"
	"errors"
	"github.com/boltdb/bolt"
	"strconv"
)

// TokenRepository - тип в котором хранится объект бд
type TokenRepository struct {
	db *bolt.DB
}

// NewTokenRepository - создание бд
func NewTokenRepository(db *bolt.DB) *TokenRepository {
	return &TokenRepository{db}
}

// Save - сохранение в бд чат айди и конкретный токен (аксес, реквест)
func (r *TokenRepository) Save(chatId int64, token string, bucket repository.Bucket) error {
	// обновляется база данных, внутри которой создается транзакция
	return r.db.Update(func(tx *bolt.Tx) error {
		// создаем ячейку для работы с ней
		b := tx.Bucket([]byte(bucket))
		// передаем информацию в ячейку и сохраняем ее в бд
		return b.Put(intToBytes(chatId), []byte(token))
	})
}

// Get - получение данных из бд
func (r *TokenRepository) Get(chatID int64, bucket repository.Bucket) (string, error) {
	var token string

	// проверяем базу данных, внутри которой создается транзакция
	err := r.db.View(func(tx *bolt.Tx) error {
		// создаем ячейку для работы с ней
		b := tx.Bucket([]byte(bucket))
		// в переменную data записываем значение из ячейки с ключом chatID который мы передали
		data := b.Get(intToBytes(chatID))
		// получаем токен
		token = string(data)
		return nil
	})
	if err != nil {
		return "", err
	}
	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil
}

// перевод из int в байты
func intToBytes(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}
