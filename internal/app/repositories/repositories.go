package repositories

import "database/sql"

func SQLCreateTableURLs(urlsDB *sql.DB) error {
	q := `CREATE TABLE IF NOT EXISTS public.urls
(
    id uuid NOT NULL,
    short_url text COLLATE pg_catalog."default" NOT NULL,
    origin_url text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT urls_pkey PRIMARY KEY (id)
)`
	_, errExec := urlsDB.Exec(q)
	return errExec
}
