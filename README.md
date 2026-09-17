# JAT - Job Application Tracker

<video src="./docs/assets/demo-clip.mp4" poster="./docs/assets/demo-clip.jpg" controls muted preload="metadata"></video>

JAT is a self-hosted job application tracker designed to help users manage their job search pipeline, from initial applications to final offers.
**Great for people that simply want to track their job applications without using an Excel sheet.**

I initially started to build this application as my final project for the [boot.dev](https://boot.dev) **Backend Go Developer** course.
However I quickly realized that I was trying to bundle too many features into what was supposed to be a simple skill showcase project.
Because of that I decided to build the course showcase project and link it to **v0.2.0** and from then on I would add the entire feature set that are described [here](./ROADMAP.md) as time goes on.

The end goal is to build a **monolithic Go binary** that bundles both a web-based user interface and a terminal user interface (TUI), making it easy to deploy as a single executable or container image.

Recommended for individuals or vocational schools/universities that want to provide a useful tool to help their students keep track of job applications, within the facilities local network.  
This application was designed to be hosted by users in their own local networks/infrastructures. 

> **Status:** This project is under heavy development. Features, APIs, and data models may change without notice until this message is removed.

## Demo

1. Clone this repository using `git clone https://github.com/filz0r/jat`
2. Install the dependencies outlined in [here](./docs/development.md#requirements)
3. Run `make demo` inside a terminal opened with the downloaded repository
4. Go to `http://localhost:4200` in your browser
5. Login using `test@test.com` and `123123` to use a normal user
6. Login using `a@a.com` and `123123` to use a admin user
7. Test out the application

## Running locally

1. Clone this repository using `git clone https://github.com/filz0r/jat`.
2. Copy `.env.sample` to `.env`.
3. Install the dependencies and make changes to the `.env` file [(check this doc for further reference)](./docs/development.md#developing).
4. Then run the build step `make`.
5. Launch the database container using `make dev_db_up`.
6. Launch the server using `./jat server`.
> **Note:** Eventually this will change to a single docker-compose.yml file that will essentially do all these steps for you, but I'm leaving that to a different release.

## Development

Check the development [docs](./docs/development.md).

## Roadmap

Moved to [here](./ROADMAP.md).

## Contributing

Moved to [here](./docs/development.md#contributions).