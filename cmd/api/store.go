package main

import (
	"database/sql"

	"github.com/gitpods/gitpods/handler"
	"github.com/gitpods/gitpods/store"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

func NewRouterStore(driver string, dsn string, secret []byte) (*handler.RouterStore, error) {
	cookieStore := sessions.NewFilesystemStore("./dev/sessions/", secret)

	routerStore := handler.RouterStore{
		CookieStore: cookieStore,
	}

	if driver == "memory" {
		usersStore := store.NewUsersInMemory()
		repositoriesStore := store.NewRepositoriesInMemory(usersStore)
		usersRepositoriesStore := store.NewUsersRepositoriesInMemory(usersStore, repositoriesStore)

		routerStore.UsersStore = usersStore
		routerStore.UsersRepositoriesStore = usersRepositoriesStore
		routerStore.AuthorizeStore = usersStore
	}

	if driver == "postgres" {
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return nil, err
		}

		if err := db.Ping(); err != nil {
			return nil, err
		}

		usersStore := store.NewUsersPostgres(db)
		routerStore.UsersStore = usersStore
		routerStore.AuthorizeStore = usersStore
	}

	return &routerStore, nil
}
