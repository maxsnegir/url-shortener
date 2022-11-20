package storage

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func (ps *PostgresStorage) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	var originalURL string
	err := ps.db.GetContext(ctx, &originalURL, "SELECT original_url FROM url_data ud WHERE ud.short_url=$1", shortURL)
	return originalURL, err
}

func (ps *PostgresStorage) saveURLData(ctx context.Context, tx *sql.Tx, urlData URLData) (int, error) {
	const query = "INSERT INTO url_data(short_url, original_url)  VALUES ($1, $2) RETURNING url_data_id;"
	var urlDataID int
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return urlDataID, err
	}
	defer stmt.Close()
	if err := stmt.QueryRow(urlData.ShortURL, urlData.OriginalURL).Scan(&urlDataID); err != nil {
		return urlDataID, err
	}
	return urlDataID, nil
}

func (ps *PostgresStorage) saveUserURL(ctx context.Context, tx *sql.Tx, userID string, urlDataID int) error {
	const query = "INSERT INTO user_url VALUES ($1, $2);"
	_, err := tx.ExecContext(ctx, query, userID, urlDataID)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStorage) SetShortURL(urlData URLData) error {
	return nil
}

func (ps *PostgresStorage) SaveData(ctx context.Context, userID string, urlData URLData) (err error) {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	urlDataID, err := ps.saveURLData(ctx, tx, urlData)
	if err != nil {
		return err
	}
	if err := ps.saveUserURL(ctx, tx, userID, urlDataID); err != nil {
		return err
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
		);
`
	var userURLs []URLData
	err := ps.db.SelectContext(ctx, &userURLs, query, userID)
	return userURLs, err

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
		CREATE INDEX IF NOT EXISTS user_token ON user_url (user_token);
	`
	ps.db.MustExec(schema)
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
