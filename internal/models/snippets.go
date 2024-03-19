package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a Snippet type to hold the data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets
// table?
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type BookDB struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *BookDB) Insert(title string, content string, expires int) (int, error) {
	stmt := `
	INSERT INTO snippets (title, content, created, expires)
	VALUES(?,?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))
	`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *BookDB) Get(id int) (*Snippet, error) {
	stmt := `
	SELECT id, title, content, created, expires FROM snippets where expires > UTC_TIMESTAMP() AND id = ?
	`

	row := m.DB.QueryRow(stmt, id)

	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *BookDB) Latest() ([]*Snippet, error) {
	stmt := `
	SELECT id, title, content, created, expires FROM snippets WHERE
	expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10
	`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// Important not to forget to close connection!
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil

}

func OpenDB(connStr string) (*BookDB, error) {
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &BookDB{
		DB: db,
	}, nil
}

func (s *BookDB) SeedDB() error {

	query := `CREATE TABLE IF NOT EXISTS snippets (
		id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
		title VARCHAR(100) NOT NULL,
		content TEXT NOT NULL,
		created DATETIME NOT NULL,
		expires DATETIME NOT NULL)`

	_, err := s.DB.Exec(query)

	if err != nil {
		return err
	}

	query = `
	INSERT INTO snippets (title, content, created, expires) VALUES ("An old silent pond","An old silent pon..\nA frog jumps into the pond,\nsplash! Silence again.\n\n- Matsuo Basho!",UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(),INTERVAL 365 DAY))
	`
	_, err = s.DB.Exec(query)

	if err != nil {
		return err
	}

	query = `
	INSERT INTO snippets (title, content, created, expires) VALUES (
		"Over the wintry forest",
		"Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n- Natsume soseki",
		UTC_TIMESTAMP(),
		DATE_ADD(UTC_TIMESTAMP(),INTERVAL 365 DAY)
	)
	`
	_, err = s.DB.Exec(query)

	if err != nil {
		return err
	}

	query = `
	INSERT INTO snippets (title, content, created, expires) VALUES (
		"First autumn morning",
		"First autumn morning\nthe mirror I stare into\nshows my fathers face. \n\n - Murakami kijo",
		UTC_TIMESTAMP(),
		DATE_ADD(UTC_TIMESTAMP(),INTERVAL 365 DAY)
	)
	`
	_, err = s.DB.Exec(query)

	return err
}
