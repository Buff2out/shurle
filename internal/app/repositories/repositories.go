package repositories

import "database/sql"

func SQLCreateTableURLs(urlsDB *sql.DB) error {
	q := `CREATE TABLE IF NOT EXISTS urls
    (
        id SERIAL PRIMARY KEY,
        short text NOT NULL,
        origin text NOT NULL
    )`
	_, errExec := urlsDB.Exec(q)
	return errExec
}

func SQLInsert(urlsDB *sql.DB) error {
	b1, b2 := "bla1", "bla2"
	_, errExec := urlsDB.Exec("INSERT INTO urls(short, origin) VALUES($1, $2)", b1, b2)
	return errExec
}

func SQLSelectAll(urlsDB *sql.DB) (*sql.Rows, error) {
	q := `SELECT * FROM urls`
	result, errExec := urlsDB.Query(q)
	return result, errExec
}
