# indego-app

This application fetches and stores real-time data from the Indego bike share system and weather services. The app stores this data in a PostgreSQL database and exposes several API endpoints to retrieve snapshots of station data and weather conditions at specific times.

Refer to [design_documentation/design_documentation.md](design_documentation/design_documentation.md) for implementation details.

## Prerequisites

Ensure you have the following prerequisites installed on your machine:

1. **OS**: Ubuntu 22.04
2. **Go**: 1.23.0 or later
3. **Docker**: for running the PostgreSQL database
4. **Postman or curl**: for hitting the API endpoints
5. **Swag CLI**: for generating Swagger documentation (optional)
   - Install Swag CLI with: 
     ```bash
     go install github.com/swaggo/swag/cmd/swag@latest
     ```
   - **Note**: You only need to install Swag CLI if you plan to modify or regenerate the Swagger documentation. If you're just running the app, the documentation is already generated, and you can view it at `http://localhost:3000/swagger/index.html#/`.


## Setting Up the Environment

### 1. Clone the repository
```bash
git clone https://github.com/nyeinsoe26/indego-app.git
cd indego-app
```

### 2. Set Up Environment Variables

Ensure you have the following in your .env file:
```bash
CLIENT_ID=your-client-id
CLIENT_SECRET=your-client-secret
AUTH_DOMAIN=https://your-auth-domain
```
For convenience, the client ID and secret are already provided in the `.env` file. These will be used for generating access tokens to authenticate with the API.

### 3 Running the app via script:

There are 2 shell scripts to run the app.

1. To run the app in local machine(I say local but it run postgres docker container but runs the go code locally), run the following script:
```bash
./scripts/run_local.sh
```

2. To dockerize the app and run it as container, run the following script:
```bash
./scripts/run_docker.sh
```

Essentially, both the scripts will run the postgres container, do migration, then launch the app.

### 4. Steps to launching the app manually
If you want to try launching things manually, here is what you need to do.
1. Launch postgres container.
```bash
docker run -d --rm --name=indego_db -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=pass123 -e POSTGRES_DB=indego_db postgres:17.0-alpine3.20
```
2. Migrate the Database
```bash
go run migration/migrate.go --up --config config.yaml
```

3. Launch the Application
```bash
go run cmd/main.go --config config.yaml
```

### 5. Access the API Documentation
You can view the auto-generated API documentation at:
```bash
http://localhost:3000/swagger/index.html#/
```

## Hitting the Endpoints

### Via Auth0 JWT (Programmatic Access)
Use this method if you are accessing the API programmatically (e.g., via curl, Postman, or any client that needs to authenticate via a JWT token).

You can obtain a JWT access token from your Auth0 tenant:
```bash
curl --request POST \
  --url https://<YOUR-AUTH0-DOMAIN>/oauth/token \
  --header 'content-type: application/json' \
  --data '{
    "client_id": "<YOUR_CLIENT_ID>",
    "client_secret": "<YOUR_CLIENT_SECRET>",
    "audience": "https://your-api-identifier",
    "grant_type": "client_credentials"
  }'
```

Replace <YOUR-AUTH0-DOMAIN>, <YOUR_CLIENT_ID>, and <YOUR_CLIENT_SECRET> with actual values provided in `.env`.
The response should include an access token:
```bash
{
  "access_token": "YOUR_ACCESS_TOKEN",
  "token_type": "Bearer"
}
```

In Postman or curl, add the following header to your requests:
```bash
Authorization: Bearer <YOUR_ACCESS_TOKEN>
```

### Via Browser (Session-based Authentication)
To access the API via your browser, you will first need to log in:
```
http://localhost:3000/login
```

You can use your Google account to log in. After successful login, you can hit the following endpoints directly in your browser without manually generating a JWT token:
```
http://localhost:3000/api/v1/stations?at=2019-09-01T10:00:00Z

http://localhost:3000/api/v1/stations/3a06?at=2019-09-01T10:00:00Z
```
Session-based authentication will be handled automatically after login.



## Testing
### 1. Unit Tests
Run the unit tests with code coverage (this excludes functional tests):

```bash
go test -v -coverprofile=coverage.out -run='^Test[^Functional]' ./...
```

### 2. Functional Tests
```bash
go test ./internal/app/api/functional_test.go
```