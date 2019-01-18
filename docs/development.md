# Development

GitPods is written in Go and Dart.
All backend, server-side components are in Go, and all the web UI is written with AngularDart.

Please make sure to have both Go and Dart installed, as well as Docker running.

## Setting up

First clone gitpods to `$GOPATH/src/github.com/gitpods/gitpods` and then change into the directory.
Next you should compile the gitpods binary itself, which will help you setup all other components.

```bash
go install -v ./cmd/gitpods
```

_Note: Make sure to have `$GOPATH/bin` in your `$PATH`_.

Now we can use that binary to setup our database (Cockroach) with docker and
pull down all dart dependencies and Caddy, as a development proxy:

```bash
gitpods dev setup
```

## Compiling the binaries

If you simply want to compile all binaries you can use:

```bash
gitpods build
```

## Installing migrations

```bash
./dev/api db migrate --migrations-path ./schema/cockroach/ --database-dsn=postgres://root@localhost:26257/gitpods?sslmode=disable
```

You can check migrated data in Cockroach Console on [localhost:8080](http://localhost:8080/).

## Creating user

```bash
./dev/api users create --email admin@localhost.com --username admin --name Admin --password password
```

After running server application you can sign to UI via entered email and password.

## During Development

We have created a wrapper to run all components of GitPods at once, and also shut them down at once.

```bash
gitpods dev
```

This will start Caddy as a proxy on [localhost:3000](http://localhost:3000).
It will proxy all requests to the UI component (or the dart development server if enabled)
and those requests starting with `/api` will be proxied to the API.
Additionally the storage component will be run to serve the API.

Killing this program will also kill all components at once.

You can enable tracing for OpenTracing for all components by running;

#### Live reloading Go binaries

The command also has a flag to enable watching all Go files.
Once a change on the filesystem has been detected it will recompile the binaries and
only on success restart the components with the new binaries, making development as easy and quick as possible.

```bash
gitpods dev --watch
```


#### Tracing

```bash
gitpods dev --tracing
```

### Dart (UI) development

It is the easiest to simply run
```bash
gitpods dev --dart
```

Instead of compiling and running the UI Go binary this will start the Dart development server
`webdev serve` and compile the UI every time a change to `/ui` has been detected.
