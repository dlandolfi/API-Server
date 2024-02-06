# Practice API server

This is a practice API server for my friends and I. Currently there is a Go backend and a Postgres database. There are plans to add more microservices are soon.

## Getting Started

- Install [Docker](https://docs.docker.com/engine/install/)
- Start everything up with `docker-compose up` while in the root dir of this project
  - If changes to the code are made you may need to add the `--build` to force a rebuild
- View logs with `docker-compose logs`
  - add the `-f` flag for realtime output
- Shutdown and remove containers with `docker-compose down`

## Endpoints

### Home Route

- Endpoint: /
- Description: Displays a simple "Hello World" message.

### Public Test Route

- Endpoint: /api/v1/testpublic
- Description: A public route for testing purposes.

### Protected Test Route

- Endpoint: /api/v1/testprivate
- Description: A protected route requiring authentication.

### Authenticated Test Route

- Endpoint: /api/v1/testauthenticated
- Description: An authenticated route making an external API request.
- Note: A `config.json` with valid credentials will be needed.

## Dependencies

This project relies on the following external libraries:

- [pq](https://github.com/lib/pq) v1.10.9
