package postgres

import (
	"database/sql"

	"github.com/dkotsyuruba/go-shortener/internal/model"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(dsn string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func (pr *PostgresRepository) Save(link *model.Link) (*model.Link, error) {
	query := `
	INSERT INTO links (id, original_url) VALUES ($1, $2)
	ON CONFLICT DO NOTHING
	`
	result, err := pr.db.Exec(query, link.ID, link.OriginalURL)
	if err != nil {
		return link, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return link, err
	}

	if rows == 0 {
		foundLink, ok := pr.FindByOriginalURL(link.OriginalURL)
		if ok {
			return foundLink, model.ErrDuplicatedURL
		}
	}

	return link, nil
}

func (pr *PostgresRepository) SaveAll(links []*model.Link) error {
	tx, err := pr.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
        INSERT INTO links (id, original_url) 
        VALUES ($1, $2)
        ON CONFLICT DO NOTHING
    `)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, link := range links {
		_, err := stmt.Exec(link.ID, link.OriginalURL)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (pr *PostgresRepository) FindByID(id string) (*model.Link, bool) {
	query := `
    SELECT id, original_url FROM links WHERE id=$1
    `
	var link model.Link
	err := pr.db.QueryRow(query, id).Scan(&link.ID, &link.OriginalURL)
	if err != nil {
		if err == sql.ErrNoRows {
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
	var link model.Link
	err := pr.db.QueryRow(query, url).Scan(&link.ID, &link.OriginalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		return nil, false
	}

	return &link, true
}

func (pr *PostgresRepository) Close() error {
	err := pr.db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (pr *PostgresRepository) Ping() error {
	err := pr.db.Ping()
	if err != nil {
		return err
	}

	return nil
}
