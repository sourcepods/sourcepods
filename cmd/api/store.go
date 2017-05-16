package main

import (
	"database/sql"

	"github.com/gitpods/gitpods/handler"
	"github.com/gitpods/gitpods/store"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

func NewRouterStore(driver string, dsn string, secret []byte) (*handler.RouterStore, StoreCloser, error) {
	cookieStore := sessions.NewFilesystemStore("./dev/sessions/", secret)

	routerStore := handler.RouterStore{
		CookieStore: cookieStore,
	}

	var closer StoreCloser

	switch driver {
	case "memory":
		usersStore := store.NewUsersInMemory()
		repositoriesStore := store.NewRepositoriesInMemory(usersStore)
		usersRepositoriesStore := store.NewUsersRepositoriesInMemory(usersStore, repositoriesStore)

		closer = func() error { return nil }

		routerStore.UsersStore = usersStore
		routerStore.UsersRepositoriesStore = usersRepositoriesStore
		routerStore.AuthorizeStore = usersStore
	default:
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return nil, nil, err
		}

		if err := db.Ping(); err != nil {
			return nil, nil, err
		}

		closer = db.Close

		usersStore := store.NewUsersPostgres(db)
		routerStore.UsersStore = usersStore
		routerStore.AuthorizeStore = usersStore
	}

	return &routerStore, closer, nil
}
