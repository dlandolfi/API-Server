# Practice API server

This is a practice API server for my friends and I. Currently there is a Go backend and a Postgres database. There are plans to add more microservices soon.

## Getting Started

- Install [Docker](https://docs.docker.com/engine/install/)
- Start everything up with `docker-compose up` while in the root dir of this project
  - If changes to the code are made you may need to add the `--build` to force a rebuild
- View logs with `docker-compose logs`
  - add the `-f` flag for realtime output
- Shutdown and remove containers with `docker-compose down`

## Endpoints

### Home Route

- **Endpoint**: `/`
- **Description**: Displays a simple "Hello World" message.

### Get User Route

- **Endpoint**: `/api/v1/getuser`
- **Description**: A route for querying the database for a user by id
- **Usage**: `curl -X GET "/api/v1/getuser?id=123"`

### Protected Test Route

- **Endpoint**: `/api/v1/testprivate`
- **Description**: A protected route requiring authentication.

### Authenticated Test Route

- **Endpoint**: `/api/v1/testauthenticated`
- **Description**: An authenticated route making an external API request.
- **Note**: A `config.json` with valid credentials will be needed.
