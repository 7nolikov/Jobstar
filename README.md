# Jobstar

A Golang web application that allows posting job vacancies, managing candidate pipelines with a state machine, and visualizing statistics with graphs. Built with PostgreSQL and HTMX.

## Table of Contents

1. [Project Overview](#project-overview)
2. [Setting up the Project](#setting-up-the-project)
3. [Creating the Git Repository](#creating-the-git-repository)
4. [Implementing HTMX for Interactivity](#implementing-htmx-for-interactivity)
5. [Database Schema and Migrations](#database-schema-and-migrations)
6. [Writing Tests and Using Best Practices](#writing-tests-and-using-best-practices)
7. [Setting Up Continuous Integration (CI)](#setting-up-continuous-integration-ci)
8. [Deployment to Google Cloud Free Tier](#deployment-to-google-cloud-free-tier)
9. [Releasing the First Version](#releasing-the-first-version)
10. [Best Practices Recap](#best-practices-recap)
11. [Congratulations](#congratulations)

## Project Overview

### Functionality of the Application

Job Vacancies Management: Create, update, delete, and list job vacancies.
Candidate Pipeline Management: Track candidates through different stages using a state machine (e.g., Applied, Interviewing, Offered, Hired, Rejected).
Statistics Visualization: Display graphs showing statistics like the number of applicants at each stage.

### Prerequisites

- Go 1.19 or later
- PostgreSQL
- Docker (for containerization)
- Google Cloud Account (for deployment)

### Use of HTMX for Server-Driven UI Interactions

Dynamic Content Loading: HTMX allows partial page updates without full reloads.
Server-Driven Interactions: Simplifies the implementation of AJAX requests handled by the server.
Enhanced User Experience: Provides a more responsive interface with less JavaScript code.

### Use of PostgreSQL as the Database

Data Storage: Stores job postings, candidate information, and pipeline data.
Reliability: PostgreSQL is a robust, open-source relational database.
Advanced Features: Supports complex queries and transactions.

## Setting up the Project

### Best Practices for Structuring a Golang Project

- Project layout

```go
your-project/
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── handlers/
│   ├── models/
│   └── db/
├── pkg/
├── configs/
├── scripts/
├── test/
├── go.mod
├── go.sum
└── README.md

```

- Package Organization:

  - cmd/: Entry points of the application.
  - internal/: Private application code.
  - pkg/: Reusable packages.
  - configs/: Configuration files.
  - scripts/: Build and deploy scripts.
  - test/: Test data and utilities.

- Configuration Management:

  - Use environment variables.
  - Utilize packages like github.com/spf13/viper or github.com/kelseyhightower/envconfig.

### Setting Up the Development Environment

Install Go

- Download and Install:
  - Go to golang.org/dl and download the latest version.
  - Follow the installation instructions for your operating system.
- Verify Installation:

```bash
go version
```

Install PostgreSQL
- MacOS:

```bash
brew install postgresql
```

- Ubuntu:

```bash
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib
```

- Windows: Download the installer from postgresql.org.

Set Up HTMX

- No installation required; include it via CDN in your HTML templates:

```html
<script src="https://unpkg.com/htmx.org@1.7.0"></script>
```

### Using Go Modules for Dependency Management

- Initialize Go Modules:

```bash
go mod init github.com/yourusername/yourproject
```

- Add Dependencies:
  - Import packages in your code, and Go will automatically add them to go.mod.
  - Alternatively, use go get:

```bash
go get github.com/gorilla/mux
```

## Creating the Git Repository

### Initialize the Repository

1. Create a New Directory:

```bash
mkdir your-project
cd your-project
```

2. Initialize Git:

```bash
git init
```

3. Create a .gitignore File:

```bash
touch .gitignore
```

Add the following to .gitignore:

```bash
# Binaries
/bin/
/build/
*.exe

# Vendor directory
/vendor/

# Logs
*.log

# Environment variables
.env

# IDE files
.idea/
.vscode/
```

4. Initial Commit:

```bash
git add .
git commit -m "Initial commit"
```

### Link to a Remote Repository

1. Create a Repository on GitHub/GitLab (without README or .gitignore).
2. Add Remote Origin:

```bash
git remote add origin https://github.com/yourusername/your-project.git
```

3. Push to Remote:

```bash
git push -u origin master
```

## Implementing HTMX for Interactivity

### Include HTMX in Templates

#### In your base HTML template:

```html
<head>
    <!-- Other head elements -->
    <script src="https://unpkg.com/htmx.org@1.7.0"></script>
</head>
```

#### Use HTMX for Form Submissions

Example: Adding a Candidate

```html
<form hx-post="/candidates" hx-target="#candidate-list" hx-swap="beforeend">
    <!-- Form fields -->
    <button type="submit">Add Candidate</button>
</form>
```

#### Handling the Request in Go

```go
func CreateCandidate(w http.ResponseWriter, r *http.Request) {
    // Process form data
    // Save to database
    // Return a partial HTML snippet for the new candidate
}
```

## Database Schema and Migrations

### Define the Database Schema

#### Create Tables

```sql
-- Vacancies Table
CREATE TABLE vacancies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    location VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Candidates Table
CREATE TABLE candidates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    vacancy_id INTEGER REFERENCES vacancies(id),
    state VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

#### Manage Migrations with `golang-migrate`

Install `golang-migrate`

- MacOS:

```bash
brew install golang-migrate
```

- Others: Download from golang-migrate GitHub.

#### Create Migrations Directory

```bash
mkdir -p db/migrations
```

#### Create a Migration File

```bash
migrate create -ext sql -dir db/migrations -seq create_tables
```

#### Write SQL in Migration Files

- xxxx_create_tables.up.sql: Place CREATE TABLE statements.
- xxxx_create_tables.down.sql: Place DROP TABLE statements.

#### Run Migrations

```bash
migrate -database "postgres://user:password@localhost:5432/yourdb?sslmode=disable" -path db/migrations up
```

## Writing Tests and Using Best Practices

### Best Practices

- Write Tests Early: Implement tests as you develop features.
- Use the `testing` Package: Standard Go testing tools.
- Organize Tests: Place test files alongside code files with `_test.go` suffix.
- Mock External Dependencies: Use interfaces to mock database interactions.

### Example Unit Test

```go
package handlers_test

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "your-project/internal/handlers"
)

func TestListVacancies(t *testing.T) {
    req, err := http.NewRequest("GET", "/vacancies", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(handlers.ListVacancies)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Additional assertions
}
```

### Run Tests

```bash
go test ./...
```

## Setting Up Continuous Integration (CI)

### Using GitHub Actions

Create Workflow File

Create `.github/workflows/ci.yml`:

```yaml
name: CI

on:
  push:
    branches: [master]
  pull_request:
    branches: [master]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: '^1.19'
    - name: Install dependencies
      run: go mod download
    - name: Run tests
      run: go test ./...
```

## Deployment to Google Cloud Free Tier

### Set Up Google Cloud Account and Project

1. Sign Up: Create an account at cloud.google.com.
2. Create a New Project: In the console, click on the project dropdown and select "New Project".

### Configure PostgreSQL Instance

1. Navigate to Cloud SQL: In the Google Cloud Console.
2. Create Instance:
  2.1 Choose PostgreSQL.
  2.2 Set up instance ID, root password, and region.
3. Configure Connectivity:
  3.1 Enable public IP.
  3.2 Set authorized networks to your IP for development.

### Deploy Go Application with Cloud Run

#### Containerize the Application

1. Create a Dockerfile:

```dockerfile
# Use official Golang image as base
FROM golang:1.19-alpine

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main ./cmd/app

# Expose port
EXPOSE 8080

# Start the application
CMD ["./main"]
```

2. Build the Docker Image:

```bash
docker build -t gcr.io/your-project-id/your-app .
```

#### Deploy to Cloud Run

1. Enable Cloud Run API:

```bash
gcloud services enable run.googleapis.com
```

2. Deploy the Service:

```bash
gcloud run deploy your-app \
  --image gcr.io/your-project-id/your-app \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --add-cloudsql-instances your-instance-connection-name \
  --set-env-vars DB_CONNECTION_STRING="postgres://user:password@/dbname?host=/cloudsql/your-instance-connection-name&sslmode=disable"
```

### Set Up Environment Variables

- DB_CONNECTION_STRING: Use the Cloud SQL instance connection name.
- Other Configurations: Set any other necessary environment variables.

### Configure CI/CD for Deployment

Update GitHub Actions Workflow
Add deployment step to `ci.yml`:

```yaml
- name: Authenticate to Google Cloud
  uses: google-github-actions/auth@v0
  with:
    credentials_json: ${{ secrets.GCP_CREDENTIALS }}

- name: Set up Cloud SDK
  uses: google-github-actions/setup-gcloud@v0

- name: Deploy to Cloud Run
  run: |
    gcloud run deploy your-app \
      --image gcr.io/your-project-id/your-app \
      --region us-central1 \
      --platform managed \
      --allow-unauthenticated \
      --add-cloudsql-instances your-instance-connection-name \
      --set-env-vars DB_CONNECTION_STRING="postgres://user:password@/dbname?host=/cloudsql/your-instance-connection-name&sslmode=disable"
```

Ensure you have added `GCP_CREDENTIALS` as a secret in your GitHub repository.

## Releasing the First Version

### Tagging the Release

1. Create a Git Tag:

```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

2. Update Release Notes: In GitHub/GitLab, add release notes describing the features.

### Monitoring the Application

- Cloud Monitoring: Use to track performance metrics.
- Cloud Logging: Access logs to troubleshoot issues.
- Error Reporting: Receive alerts for application errors.

## Best Practices Recap

- Project Structure: Maintain a clean, standard Go project layout.
- Dependency Management: Use Go modules and pin dependency versions.
- Configuration Management: Use environment variables and configuration packages.
- Database Migrations: Use tools like golang-migrate to manage schema changes.
- Testing: Write unit and integration tests; aim for comprehensive coverage.
- CI/CD: Automate builds, testing, and deployments with tools like GitHub Actions.
- Security: Secure sensitive data and use best practices for authentication.
- Code Quality: Follow clean code principles; write readable and maintainable code.
- Documentation: Document your code and provide clear README files.
- Monitoring and Logging: Implement logging and monitoring for production environments.

## Congratulations

You've built and deployed a Golang web application with PostgreSQL and HTMX, following industry best practices. Continue to iterate on your project, adding features, improving performance, and refining the user experience.
