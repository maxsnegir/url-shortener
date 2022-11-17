package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type PostgresStorage struct {
	connection *pgx.Conn
}

func (ps *PostgresStorage) GetOriginalURL(shortURL string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (ps *PostgresStorage) SetShortURL(urlData URLData) error {
	//TODO implement me
	panic("implement me")
}

func (ps *PostgresStorage) GetUserURLs(userID string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (ps *PostgresStorage) SetUserURL(userID string, shortURL string) error {
	//TODO implement me
	panic("implement me")
}

func (ps *PostgresStorage) Ping(ctx context.Context) error {
	return ps.connection.Ping(ctx)
}

func (ps *PostgresStorage) Shutdown(ctx context.Context) error {
	return ps.connection.Close(ctx)
}

func NewPostgresStorage(ctx context.Context, dsn string) (ShortenerStorage, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	storage := &PostgresStorage{
		connection: conn,
	}
	if err := storage.Ping(ctx); err != nil {
		return nil, err
	}
	return storage, nil
}
