# Development Guide

## Requirements

- Go 1.26+
  - [installation docs](https://go.dev/doc/install) or use your OS package manager if applicable.
- Node.js and npm (for the React web UI).
  - [installation docs](https://nodejs.org/en/download) or use your OS package manager if applicable.
- [swag](https://github.com/swaggo/swag) CLI.
  - after installing go run `go install github.com/swaggo/swag/cmd/swag@latest` in the terminal.
- [air](https://github.com/air-verse/air) CLI.
  - after installing go run `go install github.com/air-verse/air@latest` in the terminal.
- [GNU make](https://www.gnu.org/software/make/) 
  - Install it using your OS package manager (eg: `sudo pacman -Syu make` or `sudo apt install make`)
- Docker and Docker Compose (required for the PostgreSQL DB but you can also install it manually if you prefer so, but then the make commands won't work as they were implemented for using the `docker-compose.db.yml` file)
  - Arch Linux users can install this using `sudo pacman -Syu docker docker-compose` other Linux distros/OS's should refer to the [docker docks](https://docs.docker.com/engine/install/).
  - You can also install the Docker Desktop if you don't like using a terminal, but this won't be covered in this guide.

## Using the Makefile

There are a bunch of commands and settings in the Makefile that were added to simplify both development and deployment of newer versions, in this guide we will be focusing on the development side of things

```bash
# setup the web ui environment (basically installs the Node.js dependencies for the Web UI)
$ make dev_setup_web

# create the PostgreSQL development DB
$ make dev_db_up

# start backend server
$ make dev_backend

# start the frontend server
$ make dev_frontend

# regenerate API spec files
$ make generate_api

# reset database to a blank state
$ make reset_db

# stop the development database
$ make dev_db_down

# connect directly to the database container
$ make connect_db
```
Note: 
If you want to develop both the frontend and backend at the same time, you need to open 2 different terminals (one for the backend the other for the frontend)

## Developing
1. Run `cp .env.sample .env`
2. Change the `SECRET_JWT` variable in the **.env** file and uncomment the `JAT_DEV=1` line
3. If you intend to make changes to the Database models you can also uncomment the `DEV_DB=1` line to get debug logs from GORM
4. In one terminal run `make dev_backend`
5. In a second terminal run `make dev_frontend`
6. Whenever a new API endpoint is added you need to run `make generate_api`, this will automatically regenerate the Swagger/OpenAPI spec files and update the `web/src/api/gen-spec.ts` files to allow the frontend to have access to a type safe API
7. The Vite dev server proxies `/api` to `http://localhost:4200`, keeping the browser origin consistent with the backend's `Secure; SameSite=Strict` cookies.

## Contributions
1. If this point isn't crossed out, external contributions are still not accepted by anyone and any Pull requests will be rejected.
2. While code contributions aren't currently being accepted, bug reports and opening issues requesting new features will **always** be welcome, just be mindful of AI slop.
3. As of **v0.2.0+** all features/improvements will have a GitHub issue associated with them and must first be merged into the **development** branch via a Pull Request, opening a Pull request to the main branch will result in it being closed with a comment referring to this file.
4. I don't mind the usage of AI, but I'm a firm believer that AI is a crutch and not a substitute for your legs.
   - **You must understand what the code you are contributing does, and if asked you must be able to give out explanations about how and why you chose that approach, this is *non negotiable*.**
   - In React its quite easy to determine if something is AI generated or not, as they tend to over complicate and ignore existing design patterns and implement their own.
   - As I implement new features into this app I'll also start to understand the difference between AI generated code regarding the Go backend as well.
   - First you will be warned and asked to make changes, the second time and those after your PR's will be closed with a comment referring to this file.
5. Eventually I'll be adding a doc file regarding design patterns and how to structure new features and link it in this point, until then, first read the code and how its structured and try to replicate the same patterns (this is something that AI can help you with).