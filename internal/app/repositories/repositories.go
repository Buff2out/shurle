package repositories

import (
	"database/sql"

	"github.com/Buff2out/shurle/internal"
)

func SQLCreateTableURLs(urlsDB *sql.DB) error {
	q := `CREATE TABLE IF NOT EXISTS urls
    (
        id SERIAL PRIMARY KEY,
        short text NOT NULL,
        origin text NOT NULL,
		hashcode text NOT NULL
    )`
	_, errExec := urlsDB.Exec(q)
	return errExec
}

func SQLInsertURL(DB *sql.DB, infoURL *internal.InfoURL) error {
	_, errExec := DB.Exec("INSERT INTO urls(short, origin, hashcode) VALUES($1, $2, $3)", infoURL.ShortURL, infoURL.OriginalURL, infoURL.HashCode)
	return errExec
}

// новая задача. Сделать бизнес логику (методы) для структуры InfoURL
// Вот что значит пришёл с новыми знаниями. Рефакторинг будет постепенно.
// Ещё теперь лучше получается придумывать названия для функций. Более универсально.
// Изначально тут было ...ByHashCode
func SQLGetOriginURL(DB *sql.DB, hashCode string) *sql.Row {
	res := DB.QueryRow("SELECT origin FROM urls WHERE hashcode = $1", hashCode)
	return res
}

func SQLInsertTest(DB *sql.DB) error {
	b1, b2 := "bla1", "bla2"
	_, errExec := DB.Exec("INSERT INTO urls(short, origin) VALUES($1, $2)", b1, b2)
	return errExec
}

func SQLSelectAll(DB *sql.DB) (*sql.Rows, error) {
	q := `SELECT * FROM urls`
	result, errExec := DB.Query(q)
	return result, errExec
}
