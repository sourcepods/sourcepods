package user

//go:generate togo sql --dialect=postgres --package=user --input=store_postgres/*.sql --output=store_sql_gen.go
