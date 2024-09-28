# 1. Overview

The Indego App is designed to collect real-time data from the Indego bike-sharing system and weather services and store this data in a PostgreSQL database. The application uses a layered architecture to ensure separation of concerns, allowing for better maintainability, scalability, and testing. The app is containerized using Docker, enabling easy setup and deployment.

# 2. Key Technologies
- Golang: The backend is built with Go for high performance and scalability.
- Docker: Both the app and the PostgreSQL database are containerized for ease of deployment and isolation.
- PostgreSQL: The relational database stores snapshots of Indego bike data and weather data.
- Auth0: Authentication is managed via Auth0, using OAuth for secure API access.

# 3. Architecture and Design
## 3.1 Layered Architecture

The app is structured into distinct layers to ensure separation of concerns. Each layer access another using an interface. The main advantage of doing it this way is that it allows us to switch out multiple components without affecting preceeding layers. 

For example, we can stop using Postgres and use some other DB and we wouldnt have to make much changes. Our new DB simply need to implement DB interface and things will be just fine.

API Layer (HTTP Handlers):
- Routes incoming requests to the appropriate services.
- Includes basic request validation and response formatting.
- Handlers are implemented using Gin.
  
Service Layer:
- Contains the business logic of the application.
- Fetches and processes data from external APIs (Indego & Weather) and interacts with the database via repository functions.
- Services are abstracted via interfaces, allowing for easier mocking in tests.

Data Access Layer (Repository):
- Contains database-related operations (fetch, insert, update).
- SQL queries are written to handle data storage and retrieval in an efficient manner.
  
Configuration Layer:
- Manages configuration loading (from config.yaml and .env).
- Allows environment variable overrides for key configuration parameters (like database connection details and API keys).

## 3.2 Configuration Management
The application configuration is managed using a combination of Viper (for loading config from YAML) and environment variables (for overriding sensitive details like API keys and database credentials).

- config.yaml: The default configuration is stored here and includes server settings, database connection details, API URLs, and authentication configurations.
- Environment variables: Sensitive and environment-specific details like DATABASE_HOST, DATABASE_USER, AUTH0_CLIENT_SECRET, etc., are overridden by environment variables defined in the .env file.
  
The configuration file is loaded at the start of the app, and any corresponding environment variables are used to override values where applicable.

# 4. Docker Setup
The app utilizes multi-stage Docker builds to create a minimal, production-ready container. Two Docker images are created:

Builder Image: Builds the Go binaries for the main application and the migration utility.
Final Image: Contains only the necessary binaries and configuration files for running the app.

# 5. API Documentation and OpenAPI Spec
The API endpoints are documented using Swagger via swaggo. The OpenAPI spec is automatically generated, making it easier for clients to interact with the app.

# 6. Hitting the API Endpoints
The application provides two ways to authenticate and hit the API endpoints:
1. Auth0 JWT Token: You can obtain a JWT token using the Auth0 credentials, which will be used for all API requests. This method takes priority if both session and token authentication are available.
2. Session-Based Authentication: Alternatively, users can log in via a web browser, which creates a session that can be used for subsequent requests. If a valid session exists, it will automatically be used for authentication.

## 6.1 Using Auth0 JWT Token
To authenticate using Auth0, you must first obtain a JWT token by hitting the Auth0 token endpoint. The obtained token should then be included in the Authorization header for subsequent API requests.

## 6.2 Using Session-Based Authentication
If you prefer session-based authentication, you can log in via the browser, which establishes a session for the user. Once authenticated, you can hit the API endpoints without needing to manually pass an access token.
