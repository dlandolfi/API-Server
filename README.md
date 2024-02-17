# Practice API server

This is a practice API server for my friends and I. Currently there is a Go backend and a Postgres database. There are plans to add more microservices soon.

## Getting Started

- Install [Docker](https://docs.docker.com/engine/install/)
- Start everything up with `docker-compose up` while in the root dir of this project
  - If changes to the code are made you may need to add the `--build` to force a rebuild
- View logs with `docker-compose logs`
  - add the `-f` flag for realtime output
- Shutdown and remove containers with `docker-compose down`

## Routes

The routes package contains the `setupRoutes(r *mux.Router)` function, which is responsible for defining and configuring HTTP routes.

### Route Definitions

- Home Route
  - Path: /
  - Handler Function: homeHandler
  - HTTP Methods: GET
- Public Test Route
  - Path: /api/v1/testpublic
  - Handler Function: testPublic
  - HTTP Methods: GET, OPTIONS
- Private Test Route
  - Path: /api/v1/testprivate
  - Handler Function: authMiddleware(http.HandlerFunc(testPrivate))
  - Middleware: authMiddleware
  - HTTP Methods: All (wrapped with authentication middleware)
- Authenticated Test Route
  - Path: /api/v1/testauthenticated
  - Handler Function: authMiddleware(http.HandlerFunc(testAuthenticated))
  - Middleware: authMiddleware
  - HTTP Methods: All (wrapped with authentication middleware)
- Get User Route
  - Path: /api/v1/getuser
  - Handler Function: getUserHandler
  - HTTP Methods: GET, OPTIONS
  - Query Parameter: id (e.g., /getuser?id=n)
- Insert User Route
  - Path: /api/v1/insertuser
  - Handler Function: createUserInDb
  - HTTP Methods: POST, OPTIONS
  - Parameters
    - `first_name`: The first name of the user.
    - `last_name`: The last name of the user.
    - `email`: The email address of the user.
- Get All Users Route
  - Path: /api/v1/getallusers
  - Handler Function: getAllUsersHandler
  - HTTP Methods: GET, OPTIONS

## Middlewares

- CORS Method Middleware: Applied globally to handle Cross-Origin Resource Sharing (CORS) for HTTP methods.
- Security Headers Middleware: Applied globally to enhance security by setting appropriate HTTP headers.
- No Cache Header Middleware: Applied globally to prevent caching of responses.
