package db

import (
	"context"
	"database/sql"
	"fmt"
	apierrors "git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/errors"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	client *sql.DB
}

func NewDatabase(config Config) (*Database, error) {
	connInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database,
	)
	client, err := sql.Open("postgres", connInfo)
	if err != nil {
		return nil, err
	}

	return &Database{client: client}, nil
}

func (db *Database) Init(ctx context.Context) error {
	_, err := db.client.ExecContext(ctx, initRequest)
	return err
}

// Drop drops all tables in database. USE ONLY FOR TESTS!
func (db *Database) Drop(ctx context.Context) error {
	_, err := db.client.ExecContext(ctx, dropRequest)
	return err
}

// Clean cleans all tables in database. USE ONLY FOR TESTS!
func (db *Database) Clean(ctx context.Context) error {
	_, err := db.client.ExecContext(ctx, cleanRequest)
	return err
}

func (db *Database) Save(ctx context.Context, shortSuffix, longLink, secretKey string, expirationDate time.Time, isDeleted bool) error {
	_, err := db.client.ExecContext(ctx, saveRequest, shortSuffix, longLink, secretKey, expirationDate, isDeleted)
	return err
}

func (db *Database) SelectBySuffix(ctx context.Context, shortSuffix string) (*Link, error) {
	row := db.client.QueryRowContext(ctx, selectBySuffixRequest, shortSuffix)

	var link Link
	var expirationTime time.Time
	err := row.Scan(&link.ShortSuffix, &link.Link, &link.SecretKey, &link.Clicks, &link.ExpirationDate, &link.Deleted)
	switch {
	case err == sql.ErrNoRows:
		return nil, SuffixNotFoundError
	case err != nil:
		return nil, err
	}

	// Проверка истечения срока действия ссылки
	// Округление до секунд
	// TODO: стоит перенести в отдельный метод, потому что этот метод вызывается не только во время редиректа
	expirationTime = link.ExpirationDate.Truncate(time.Second)
	currentTime := time.Now().Truncate(time.Second).UTC() // преобразование времени в UTC, для сравнения

	if expirationTime.Before(currentTime) {
		if _, err := db.client.ExecContext(ctx, updateDeletedBySuffixRequest, shortSuffix); err != nil {
			return nil, err
		}
		return nil, apierrors.NotFoundError{}
	}

	return &link, nil
}

func (db *Database) SelectByLink(ctx context.Context, longLink string) (*Link, error) {
	row := db.client.QueryRowContext(ctx, selectByLinkRequest, longLink)

	var link Link
	err := row.Scan(&link.ShortSuffix, &link.Link, &link.SecretKey, &link.Clicks, &link.ExpirationDate, &link.Deleted)
	switch {
	case err == sql.ErrNoRows:
		return nil, LinkNotFoundError
	case err != nil:
		return nil, err
	}

	return &link, nil
}

func (db *Database) SelectBySecretKey(ctx context.Context, secretKey string) (*Link, error) {
	row := db.client.QueryRowContext(ctx, selectBySecretKeyRequest, secretKey)

	var link Link
	err := row.Scan(&link.ShortSuffix, &link.Link, &link.SecretKey, &link.Clicks, &link.ExpirationDate, &link.Deleted)
	switch {
	case err == sql.ErrNoRows:
		return nil, SecretKeyNotFoundError
	case err != nil:
		return nil, err
	}

	return &link, nil
}

func (db *Database) DeleteBySecretKey(ctx context.Context, secretKey string) error {
	res, err := db.client.ExecContext(ctx, deleteBySecretKeyRequest, secretKey)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return SecretKeyNotFoundError
	}

	return nil
}

func (db *Database) IncrementClicksBySuffix(ctx context.Context, shortSuffix string) error {
	_, err := db.client.ExecContext(ctx, incrementClicksBySuffixRequest, shortSuffix)
	return err
}
