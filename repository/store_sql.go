package repository

//go:generate togo sql --dialect=postgres --package=repository --input=store_postgres/*.sql --output=store_sql_gen.go
