package storage

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var DuplicateErr = errors.New("DuplicateErr")

type PostgresStorage struct {
	db *sqlx.DB
}

func (ps *PostgresStorage) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	const query = "SELECT original_url FROM url_data ud WHERE ud.short_url=$1;"
	var originalURL string
	err := ps.db.GetContext(ctx, &originalURL, query, shortURL)
	return originalURL, err
}

func (ps *PostgresStorage) SaveData(ctx context.Context, userID string, urlData URLData) (err error) {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	urlDataID, err := ps.createURLData(ctx, urlData)
	if err != nil {
		if ps.IsDuplicateError(err) {
			return DuplicateErr
		}
		return err
	}
	if err := ps.getOrCreateUserURL(ctx, userID, urlDataID); err != nil {
		return err
	}
	return tx.Commit()
}

func (ps *PostgresStorage) createURLData(ctx context.Context, urlData URLData) (int, error) {
	const query = `
		INSERT INTO url_data(short_url, original_url)  
		VALUES (:short_url, :original_url) 
		RETURNING url_data_id;`
	var urlDataID int
	stmt, err := ps.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return urlDataID, err
	}
	err = stmt.Get(&urlDataID, urlData)
	return urlDataID, err
}

func (ps *PostgresStorage) getOrCreateUserURL(ctx context.Context, userID string, urlDataID int) error {
	const insertQuery = `INSERT INTO user_url VALUES ($1, $2);`
	const selectQuery = `SELECT count(*) FROM user_url uu WHERE uu.user_token=$1 AND uu.url_data_id=$2`
	var count int
	if err := ps.db.GetContext(ctx, &count, selectQuery, userID, urlDataID); err != nil {
		return err
	}
	if count != 0 {
		return nil
	}
	_, err := ps.db.ExecContext(ctx, insertQuery, userID, urlDataID)
	return err
}

func (ps *PostgresStorage) SetShortURL(urlData URLData) error {
	return nil
}

func (ps *PostgresStorage) SaveDataBatch(ctx context.Context, userID string, urlData []URLData) (err error) {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, url := range urlData {
		urlDataID, err := ps.createURLData(ctx, url)
		if err != nil {
			if ps.IsDuplicateError(err) {
				return DuplicateErr
			}
			return err
		}
		if err := ps.getOrCreateUserURL(ctx, userID, urlDataID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (ps *PostgresStorage) GetUserURLs(ctx context.Context, userID string) ([]URLData, error) {
	const query = `
		SELECT ud.short_url, ud.original_url 
		FROM url_data ud 
		WHERE  ud.url_data_id IN (
		    SELECT uu.url_data_id
		    	FROM user_url uu 
		    WHERE uu.user_token = $1
		);`
	var userURLs []URLData
	err := ps.db.SelectContext(ctx, &userURLs, query, userID)
	return userURLs, err
}

func (ps *PostgresStorage) IsDuplicateError(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == "23505" {
			return true
		}
	}
	return false
}

func (ps *PostgresStorage) SetUserURL(userID string, shortURL string) error {
	//TODO implement me
	panic("implement me")
}

func (ps *PostgresStorage) Ping(ctx context.Context) error {
	return ps.db.PingContext(ctx)
}

func (ps *PostgresStorage) Shutdown(ctx context.Context) error {
	return ps.db.Close()
}

func (ps *PostgresStorage) initPostgresStorage(ctx context.Context) {
	const schema = `
		CREATE TABLE IF NOT EXISTS url_data (
		    url_data_id SERIAL PRIMARY KEY,
		    short_url VARCHAR(255) UNIQUE,
		    original_url VARCHAR(255) NOT NULL
		);
		CREATE TABLE IF NOT EXISTS user_url (
		    user_token VARCHAR(36) NOT NULL,
		    url_data_id INTEGER NOT NULL,
		    CONSTRAINT user_url_data FOREIGN KEY(url_data_id) REFERENCES url_data(url_data_id)
		);
		CREATE UNIQUE INDEX IF NOT EXISTS user_url_data ON user_url (user_token, url_data_id);
		CREATE INDEX IF NOT EXISTS user_token ON user_url (user_token);
	`
	ps.db.MustExecContext(ctx, schema)
}

func NewPostgresStorage(ctx context.Context, dsn string) (ShortenerStorage, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	storage := &PostgresStorage{
		db: db,
	}
	if err := storage.Ping(ctx); err != nil {
		return nil, err
	}
	storage.initPostgresStorage(ctx)
	return storage, nil
}
