package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func (ps *PostgresStorage) GetOriginalURL(ctx context.Context, shortURL string) (URLData, error) {
	const query = "SELECT original_url, deleted FROM url_data ud WHERE ud.short_url=$1;"
	urlData := URLData{ShortURL: shortURL}
	err := ps.db.GetContext(ctx, &urlData, query, shortURL)
	return urlData, err
}

func (ps *PostgresStorage) SaveData(ctx context.Context, userToken string, urlData URLData) (err error) {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	urlDataID, err := ps.CreateURLData(ctx, urlData)
	if err != nil {
		if isDuplicateErr(err) {
			return NewDuplicateError(urlData.ShortURL)
		}
		return err
	}
	if err := ps.CreateUserURL(ctx, userToken, urlDataID); err != nil {
		return err
	}
	return tx.Commit()
}

func (ps *PostgresStorage) CreateURLData(ctx context.Context, urlData URLData) (int, error) {
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

func (ps *PostgresStorage) CreateUserURL(ctx context.Context, userToken string, urlDataID int) error {
	const query = `INSERT INTO user_url VALUES ($1, $2);`
	_, err := ps.db.ExecContext(ctx, query, userToken, urlDataID)
	return err
}

func (ps *PostgresStorage) SaveDataBatch(ctx context.Context, userToken string, urlData []URLData) (err error) {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, url := range urlData {
		urlDataID, err := ps.CreateURLData(ctx, url)
		if err != nil {
			if isDuplicateErr(err) {
				return NewDuplicateError(url.ShortURL)
			}
			return err
		}
		if err := ps.CreateUserURL(ctx, userToken, urlDataID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (ps *PostgresStorage) GetUserURLs(ctx context.Context, userToken string) ([]URLData, error) {
	const query = `
		SELECT ud.short_url, ud.original_url 
		FROM url_data ud 
		WHERE  ud.url_data_id IN (
		    SELECT uu.url_data_id
		    	FROM user_url uu 
		    WHERE uu.user_token = $1
		);`
	var userURLs []URLData
	err := ps.db.SelectContext(ctx, &userURLs, query, userToken)
	return userURLs, err
}

func (ps *PostgresStorage) Ping(ctx context.Context) error {
	return ps.db.PingContext(ctx)
}

func (ps *PostgresStorage) Shutdown(ctx context.Context) error {
	return ps.db.Close()
}

func (ps *PostgresStorage) DeleteURLs(ctx context.Context, urlsToDelete []string) error {
	const query = `UPDATE url_data SET deleted = True  WHERE short_url = any($1);`
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = ps.db.Exec(query, pq.Array(urlsToDelete))
	if err != nil {
		return err
	}
	return tx.Commit()

}
func (ps *PostgresStorage) initPostgresStorage(ctx context.Context) {
	const schema = `
		CREATE TABLE IF NOT EXISTS url_data (
		    url_data_id SERIAL PRIMARY KEY,
		    short_url VARCHAR(255) UNIQUE,
		    original_url VARCHAR(255) NOT NULL,
		    deleted BOOLEAN NOT NULL DEFAULT FALSE
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
