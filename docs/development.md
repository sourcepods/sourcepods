# Development

SourcePods is written in [Go](https://golang.org/) and [Dart](https://www.dartlang.org/).
All backend server-side components are in Go and the web UI is written with AngularDart.

Please make sure to have both Go and Dart installed.  
For development you also need to have [Docker](https://docs.docker.com/install/) running.

## Setting up

First clone SourcePods to `$GOPATH/src/github.com/sourcepods/sourcepods` and then change into the directory.
This is not 100% necessary anymore, as we now use [Go modules](https://github.com/golang/go/wiki/Modules)
but should make life easier for everyone when debugging a problem with the maintainers.

### sourcepods-dev

SourcePods ships with a development binary called `sourcepods-dev`.
During development this binary will help you run all components concurrently
(api, storage, ui and a [Caddy](http://caddyserver.com) as reverse proxy).
It also helps to enable a feature with one command line flag with all
components at once (it forwards these flags).

```bash
make dev/sourcepods-dev
```

Now we can use that binary to setup all external dependencies,
like our database (Cockroach) which will be run as a Docker container.
It will also pull down all Dart dependencies by running `pub get` and
download a Caddy binary into `./dev/caddy` to be used a proxy during development:

After starting Cockroach as database the command will run all migrations.
You can check if the tables were successfully created in the Cockroach Console on [localhost:8080](http://localhost:8080/).

```bash
./dev/sourcepods-dev setup
```

### Creating user

For now you need to create the users manually and the api binary has this functionality built-in.
(This will be removed in the future, as we will transition to OIDC login)

```bash
GITPODS_DATABASE_DSN=postgres://root@localhost:26257/sourcepods?sslmode=disable \
    ./dev/api users create --email admin@localhost.com --username admin --name Admin --password password
```

## Begin with development

### Compiling the binaries

If you simply want to compile all binaries you can use:

```bash
make build
```

After running server application you can sign to UI via entered email and password.

### During Development

```bash
./dev/sourcepods-dev
```

To make life easier when developing on an application that is composed of multiple components,
we have created a wrapper to run all components of SourcePods at once.
This will start `dev/api`, `dev/storage`, `dev/ui` (by default it start the UI in docker though)
and Caddy as reverse proxy in front of it all.
Once you Ctrl+C and quite the command, it will gracefully shut down all components.

Check [localhost:3000](http://localhost:3000) and you show find a running instance of SourcePods running on your machine.

It will proxy all requests to the UI component (which can be in a container, the `./dev/ui` binary, or Dart's development server)
and all requests starting with `/api` will be proxied to the API component.
Additionally the storage component is running to serves the API via [gRPC](https://grpc.io/).

#### Live reloading Go binaries

The `sourcepods-dev` command has a flag to enable watching all Go files in the project's folders.
Once a change on the filesystem has been detected it will recompile the binaries and
only on success restart the components with the new binaries, making development as easy and quick as possible.
Essentially, it's an endless loop running `make build`, once a change is detected.

```bash
sourcepods-dev --watch
```

This can be really helpful when working on the API, for example. Just save your change,
wait a second and hit the endpoint again with the newest version of it running.

#### Tracing

SourcePods has built-in support for OpenTracing and [Jeager](https://www.jaegertracing.io/) is the default tracing backend.
You can start a Jaeger all-in-one container for local development by following their guide:  
https://www.jaegertracing.io/docs/1.8/getting-started/#all-in-one

Once Jaeger is running, simply add the `--tracing` flag and all traces will be sent to Jaeger.

```bash
sourcepods-dev --tracing
```

All log lines of the API and other components should contain a `request=2ae7cb76-fcec-4143-8d1d-ff4bf4913aea` key-value pair.
You can copy this key-value-pair and paste it at Jaeger's UI into the _Tags_ input field. 
After pressing enter you'll see the exact trace for the request you have the logs to.
See [#31](https://github.com/sourcepods/sourcepods/pull/31#issue-244309605) for more information.

### Dart (UI) development

Our UI is written in Dart with the [AngularDart](https://webdev.dartlang.org/angular) framework.
The source code can be found inside this repository under `/ui`.
As written above you need to have [Dart installed](https://webdev.dartlang.org/tools/sdk#install).

AngularDart has one dependency for development.
You can install [webdev](https://webdev.dartlang.org/tools/webdev) by running `pub global activate webdev`.

After that it will be easiest to use the integrated `sourcepods-dev` command that starts the reverse proxy and fowards
the requests correctly to the `webdev serve` command it starts in the background.

It is the easiest to simply run
```bash
sourcepods-dev --ui=dart
```

_Make sure to stop the Docker container running the UI!_ - `docker stop sourcepods-ui`.
