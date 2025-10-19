package postgres

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dkotsyuruba/go-shortener/internal/model"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(dsn string) (*PostgresRepository, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	repo := &PostgresRepository{pool: pool}

	err = repo.Init(dsn)
	if err != nil {
		pool.Close()
		return nil, err
	}

	return repo, nil
}

func (pr *PostgresRepository) Init(dsn string) error {
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return err
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.Database)

	m, err := migrate.New("file://migrations", connString)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (pr *PostgresRepository) Save(link *model.Link) (*model.Link, error) {
	query := `
        INSERT INTO links (id, original_url, user_id) VALUES ($1, $2, $3)
        ON CONFLICT DO NOTHING
    `

	ctx := context.Background()

	result, err := pr.pool.Exec(ctx, query, link.ID, link.OriginalURL, link.UUID)
	if err != nil {
		return link, err
	}

	rows := result.RowsAffected()

	if rows == 0 {
		foundLink, ok := pr.FindByOriginalURL(link.OriginalURL)
		if ok {
			return foundLink, model.ErrDuplicatedURL
		}
	}

	return link, nil
}

func (pr *PostgresRepository) SaveAll(links []*model.Link) error {
	ctx := context.Background()

	tx, err := pr.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	for _, link := range links {
		_, err := tx.Exec(ctx, `
            INSERT INTO links (id, original_url, user_id) 
            VALUES ($1, $2, $3)
            ON CONFLICT DO NOTHING
        `, link.ID, link.OriginalURL, link.UUID)

		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (pr *PostgresRepository) FindByID(id string) (*model.Link, bool) {
	query := `
        SELECT id, original_url FROM links WHERE id=$1
    `

	ctx := context.Background()
	var link model.Link

	err := pr.pool.QueryRow(ctx, query, id).Scan(&link.ID, &link.OriginalURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, false
		}
		return nil, false
	}

	return &link, true
}

func (pr *PostgresRepository) FindByOriginalURL(url string) (*model.Link, bool) {
	query := `
        SELECT id, original_url FROM links WHERE original_url=$1
    `

	ctx := context.Background()
	var link model.Link

	err := pr.pool.QueryRow(ctx, query, url).Scan(&link.ID, &link.OriginalURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, false
		}
		return nil, false
	}

	return &link, true
}

func (pr *PostgresRepository) FindAllByUserID(id string) ([]*model.Link, error) {
	query := `
        SELECT id, original_url FROM links WHERE user_id = $1
    `

	ctx := context.Background()

	rows, err := pr.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*model.Link

	for rows.Next() {
		var link model.Link
		err := rows.Scan(&link.ID, &link.OriginalURL)
		if err != nil {
			return nil, err
		}
		links = append(links, &link)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(links) == 0 {
		return []*model.Link{}, nil
	}

	return links, nil
}

func (pr *PostgresRepository) Close() error {
	pr.pool.Close()
	return nil
}

func (pr *PostgresRepository) Ping() error {
	ctx := context.Background()

	err := pr.pool.Ping(ctx)
	if err != nil {
		return err
	}

	return nil
}
