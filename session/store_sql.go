package session

//go:generate togo sql --dialect=postgres --package=session --input=store_postgres/*.sql --output=store_sql_gen.go
