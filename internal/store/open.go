package store

import (
	"database/sql"
	"net/url"

	_ "github.com/lib/pq"
)

func Open(dbURL string) (*sql.DB, error) {
	return sql.Open("postgres", dbURL)
}

// sanitizeDSN masks password before printing DSN.
func sanitizeDSN(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil {
		return dsn // fallback to raw
	}
	if u.User != nil {
		username := u.User.Username()
		if _, hasPwd := u.User.Password(); hasPwd {
			u.User = url.UserPassword(username, "*****")
		}
	}
	// Remove query noise like sslmode=disable for cleaner logs if you want
	//u.RawQuery = ""
	return u.String()
}
