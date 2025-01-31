# Order Service API Documentation

## Table of Contents
- [Getting Started](#getting-started)
- [Authentication](#authentication)
- [API Endpoints](#api-endpoints)
- [Error Handling](#error-handling)

## Getting Started

### Prerequisites
- Docker and Docker Compose installed on your system
- Git for cloning the repository

### Installation

1. Clone the repository:
```bash
git clone https://github.com/shahriar-mohim007/orderservice.git

Start the application using Docker Compose:

bash
Navigate to the project directory and run the following command to start the services:

```bash
docker-compose up

The service will be available at http://localhost:8080.
Authentication
The API uses JWT (JSON Web Token) for authentication. You must first register and login to obtain an access token.
User Registration
curl --location 'http://localhost:8080/api/v1/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "John Doe",
    "email": "01901901901@mailinator.com",
    "password": "321dsaf"
}'
User Login
