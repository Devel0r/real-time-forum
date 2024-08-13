package sqlite

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Pruel/real-time-forum/pkg/cstructs"
	"github.com/Pruel/real-time-forum/pkg/serror"
)

type Database struct {
	SQLite *sql.DB
	MySQL  *sql.DB
	Log    *log.Logger
}

func New() *Database {
	return &Database{
		SQLite: &sql.DB{},
	}
}

// InitDatabase
func InitDatabase(cfg *cstructs.Config) (*Database, error) {
	// Validate
	if cfg == nil {
		return nil, serror.ErrNilConfigStruct
	}

	// Connect
	sqlite, err := sql.Open("sqlite3", cfg.DatabaseFilePath)
	if err != nil {
		return nil, err
	}

	// Assign sqlite
	db := New()
	db.SQLite = sqlite

	return db, nil
}

// Close connection
func (db *Database) Close() error {
	if err := db.SQLite.Close(); err != nil {
		return err
	}
	return nil
}
