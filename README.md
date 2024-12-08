# **Connect App Platform - REST API Service**  

![Go Badge](https://img.shields.io/badge/Go-1.x-blue) ![Docker Badge](https://img.shields.io/badge/Docker-Enabled-blue) ![Postgres Badge](https://img.shields.io/badge/Postgres-Database-green) ![Redis Badge](https://img.shields.io/badge/Redis-Caching-red) ![SendGrid Badge](https://img.shields.io/badge/Email-SendGrid-blue) ![CI/CD Badge](https://img.shields.io/badge/CI%2FCD-GitHub%20Actions-blue)

---

## **Overview**  

Connect App Platform is a scalable REST API service built with Go, featuring secure user management, role-based access control, and persistent data storage using PostgreSQL. It includes Redis caching, email notifications via SendGrid, Docker containerization, and deployment to Google Cloud Run. The API also supports CI/CD pipelines configured with GitHub Actions for automated testing, building, and deployment.

---

## [Documentation Link](https://connect-app-platform-16033377029.europe-west10.run.app/v1/swagger/index.html)
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

- **`.github/`**: CI/CD GitHub Actions workflows  
- **`bin/`**: Compiled binaries (optional or for tooling)  
- **`cmd/`**: Main application entry points  
- **`docs/`**: Project documentation  
- **`internal/`**: Core business logic  
- **`scripts/`**: Scripts for automation and maintenance  
- **`tmp/`**: Temporary files  
- **`web/`**: Web frontend source code  
- **`CHANGELOG.md`**: Project changelog  
- **`docker-compose.yml`**: Docker services configuration  
- **`Dockerfile`**: Docker container build file  
- **`Makefile`**: Automation commands  

---

## Technologies Used

- **Language**: Go  
- **Database**: PostgreSQL
- **Caching**: Redis  
- **Email Service**: SendGrid  
- **Containerization**: Docker, Docker Compose  
- **Deployment**: Google Cloud Run  
- **Web Frontend**: Vue.js/Nuxt.js  
- **CI/CD**: GitHub Actions  
- **Documentation**: Swagger UI  
