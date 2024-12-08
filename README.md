# **Connect App Platform - REST API Service**  

![Go Badge](https://img.shields.io/badge/Go-1.x-blue) ![Docker Badge](https://img.shields.io/badge/Docker-Enabled-blue) ![Postgres Badge](https://img.shields.io/badge/Postgres-Database-green) ![Redis Badge](https://img.shields.io/badge/Redis-Caching-red) ![SendGrid Badge](https://img.shields.io/badge/Email-SendGrid-blue) ![CI/CD Badge](https://img.shields.io/badge/CI%2FCD-GitHub%20Actions-blue)

---

## **Overview**  

Connect App Platform is a scalable REST API service built with Go, featuring secure user management, role-based access control, and persistent data storage using PostgreSQL. It includes Redis caching, email notifications via SendGrid, Docker containerization, and deployment to Google Cloud Run. The API also supports CI/CD pipelines configured with GitHub Actions for automated testing, building, and deployment.

---

## Key Features

- **Authentication & Authorization**: JWT-based authentication with role-based access control for Admins, Moderators, and Users.
- **User Management**: Full CRUD operations for user accounts, including registration, login, and role assignment.
- **Data Persistence**: PostgreSQL for scalable and secure relational data storage.
- **Caching**: Redis caching for frequently requested data and reduced server load.
- **Email Notifications**: SendGrid-powered email notifications for account verification, password resets, and more.
- **Containerization**: Docker-based deployment for scalability and portability.
- **Cloud Deployment**: Google Cloud Run for high availability and performance.
- **Web Frontend Integration**: A connected frontend application is included in the `web` directory.
- **API Documentation**: Comprehensive API documentation using Swagger UI.
- **CI/CD Pipelines**: Automated workflows for testing, building, and deployment using GitHub Actions.

---

## Project Structure

```plaintext
.github/            # CI/CD GitHub Actions workflows
bin/                # Compiled binaries (optional or for tooling)
cmd/                # Main application entry points
docs/               # Project documentation
internal/           # Core business logic
scripts/            # Scripts for automation and maintenance
tmp/                # Temporary files
web/                # Web frontend source code
CHANGELOG.md        # Project changelog
docker-compose.yml  # Docker services configuration
Dockerfile          # Docker container build file
Makefile            # Automation commands

Technologies Used
Language: Go (Golang)
Database: PostgreSQL
Caching: Redis
Email Service: SendGrid
Containerization: Docker, Docker Compose
Deployment: Google Cloud Run
Web Frontend: Vue.js/Nuxt.js
CI/CD: GitHub Actions
Documentation: Swagger UI
Installation & Setup
1. Clone the Repository
bash
Copy code
git clone https://github.com/Shakhboz06/your-repo.git
cd your-repo
2. Setup Environment Variables
Create a .env file with the following content:

env
Copy code
DB_HOST=postgres
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_db
REDIS_HOST=redis
REDIS_PORT=6379
JWT_SECRET=your_secret
SENDGRID_API_KEY=your_sendgrid_api_key
EMAIL_FROM=your_email_address
3. Run Docker Services
bash
Copy code
docker-compose up --build
4. Access the API Documentation
Visit the Swagger UI Documentation.

API Endpoints Overview
Authentication & Authorization
Method	Endpoint	Description
POST	/auth/register	Register a new user
POST	/auth/login	Login and get JWT token
GET	/auth/me	Get current user info
User Management
Method	Endpoint	Description
GET	/users	Get all users (Admin)
POST	/users	Create a new user
PUT	/users/{id}	Update user information
DELETE	/users/{id}	Delete a user (Admin)
Email Notifications
Method	Endpoint	Description
POST	/auth/forgot-password	Send password reset email
POST	/auth/verify-email	Send email verification
CI/CD Pipeline with GitHub Actions
Automated workflows are located in the .github/workflows directory and handle:

Build & Test: Automatically run tests on every push or pull request.
Docker Build: Build and push Docker images.
Deployment: Deploy to Google Cloud Run upon successful builds and tests.🚀
