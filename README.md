# JAT - Job Application Tracker

JAT is a self-hosted job application tracker designed to help users manage their job search pipeline, from initial applications to final offers.
The end goal is to build a **monolithic Go binary** that bundles both a web-based user interface and a terminal user interface (TUI), making it easy to deploy as a single executable or container image.

This application was designed to be hosted by users in their own local networks/infrastructures. Because of this I didn't go too deep regarding security of the platform, however normal users aren't able to view or modify data that doesn't belong to them.

Great for people that simply want to track their job applications without using an Excel sheet.

Recommended for individuals or vocational schools/universities that want to provide a useful tool to help their students keep track of job applications, within the facilities local network.  

I learnt Go doing [boot.dev](https://boot.dev)'s Go Backend Developer Path, and doing so you are required to develop 2 personal projects, a starter one and a capstone one.
Instead of building 2 projects I decided it would be a better approach to simply build something that I've tried to build in the past. But at the same time spice it up a little, I have some experience building web-based applications but I've always been curious about TUI applications so I decided to combine both into a single project.
The idea is for this project to always be open source and open for anyone no matter their IT skills, as time goes on this file will have better documentation to help less technical people (check the roadmap bellow).


> **Status:** This project is under heavy development. Features, APIs, and data models may change without notice until a stable `v1.0.0` release is reached.

## Goals

- Provide a simple web UI for managing job applications and their changes over time.
- Provide a TUI interface for users who prefer working from the terminal.
- Ship everything in a single, self-contained binary with no separate frontend server required.
- Provide 3 operation modes for the binary: Standalone, Server and Client.
  - Server: A web server with a REST API that allows external clients to connect to it and also serves the react web ui, obviously this requires a database to connect to.
  - Client: A TUI application that connects to the REST API and provides most, if not all, of the features of the web ui, no need to run a database on the computer that runs it.
  - Standalone: The same TUI application as the Client mode, but without the requirement for a server running somewhere else, however it will still require a PostgreSQL database to connect to.

## Installation

JAT can be installed with `go install`:

```bash
go install github.com/filz0r/jat@latest
```

> **Note:** Installation via `go install` is currently **not recommended**. The project is still under heavy development, and breaking changes are expected. For now, please build and run the project using the provided `Makefile` from a local clone.

## Development

To build and run the server locally:

```bash
# Build the binary
make build_server

# Build and start the server in development mode
make dev_server
```

The server expects environment variables to be configured. A sample `.env` file is provided in the repository for reference.

## Roadmap

The following roadmap is a rough outline of what is planned before the first stable release. It is subject to change as the project evolves.

- [x] **v0.1.0**
  - Initial project scaffold, API server, database schema, and core job application CRUD endpoints.
- [ ] **v0.2.0** 
  - Standardize error handling (right now the messages are rather generic and not that much descriptive) with custom interfaces 
  - Implementation of the react based web ui that works with the existing REST endpoints
  - Adding more admin features to the REST API (check the comments on ./internal/api/api.go)
  - Adding OpenAPI spec (idk if this is possible due to me using normal go/http std library, but if so it will help a lot with the web ui implementation)
- [ ] **v0.3.0**
  - Adding a build pipeline for both server and web ui
  - Implementing the Client TUI interface
  - Initial docker container builds
- [ ] **v0.4.0**
  - Implementing standalone TUI interface (there's a possibility this wont be implemented)
  - Adding a statistics dashboard to the users homepage
- [ ] **v0.5.0**
  - Adding a statistics dashboard for the admin interface (web ui only)
  - Adding a notification system (probably with websockets, if possible)
  - Adding version checking from within the web server (basically triggers a notification warning all admin users that a new version was released)
- [ ] **v0.6.0**
  - Adding a way to import .csv/.xsl(x) (idk if the latter is possible but that's the idea) files into the web ui
  - Adding time based notifications such as after 2 months of a job application not having changes suggesting changing it to a ghosted state
- [ ] **v0.7.0**
  - Improving the Application State logic for all users (some of the state_kind should only have 1 record per user)
  - Improving the company logic (right now it will probably break if you delete a company that has existing job applications, still not sure on how GORM handles the soft deleted values, when they have connections to other tables)
- [ ] **v0.8.0**
  - Adding proper CLI option support (as in right now manually parsing arguments is a pain and this can easily be improved) 
  - Documentation for development
  - Documentation for deployments
  - Container image public release
- [ ] **v0.9.0**
  - Bug fixes, polish, and feature freeze for v1.0.0.
- [ ] **v1.0.0**
  - Stable public release.
- [ ] **beyond v1.0.0**
  - Open up feature requests
  - Improve GitHub dependant stuff like enforcing a standard issue format, tags etc.
  - Open up external contributions.

## Contributing

External contributions are currently **closed**. JAT is a personal project until `v1.0.0` is released, at which point contribution guidelines may be opened. Feel free to open issues for bugs or suggestions, but pull requests will not be accepted before the first stable release.
