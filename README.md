# Order Service 

## Table of Contents
- [Getting Started](#getting-started)
- [Authentication](#authentication)
- [API Endpoints](#api-endpoints)


## Getting Started

### Installation

1. Clone the repository:
```bash
git clone https://github.com/shahriar-mohim007/orderservice.git
```

2. Start the application using Docker Compose:
```bash
docker-compose up
```

The service will be available at `http://localhost:8080`.

## Authentication

The API uses JWT (JSON Web Token) for authentication. You must first register and login to obtain an access token.

### User Registration

```bash
curl --location 'http://localhost:8080/api/v1/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "John Doe",
    "email": "01901901901@mailinator.com",
    "password": "321dsaf"
}'
```

### User Login

```bash
curl --location 'http://localhost:8080/api/v1/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "01901901901@mailinator.com",
    "password": "321dsaf"
}'
```

## API Endpoints

### Create Order
Creates a new delivery order in the system.

```bash
curl --location 'http://localhost:8080/api/v1/orders/' \
--header 'Authorization: Bearer JWT_ACCESS_TOKEN' \
--header 'Content-Type: application/json' \
--data '{
    "store_id": 131172,
    "merchant_order_id": "ORD-12345",
    "recipient_name": "John Doe",
    "recipient_phone": "01712345678",
    "recipient_address": "banani, gulshan 2, dhaka, bangladesh",
    "recipient_city": 1,
    "recipient_zone": 1,
    "recipient_area": 1,
    "delivery_type": 48,
    "item_type": 2,
    "special_instruction": "Handle with care",
    "item_quantity": 1,
    "item_weight": 13,
    "amount_to_collect": 1600.00,
    "item_description": "Electronics item"
}'
```

### Get All Orders
Retrieves a paginated list of orders with optional filters.

```bash
curl --location 'http://localhost:8080/api/v1/orders/all?transfer_status=1&archive=0&limit=1&page=2' \
--header 'Authorization: Bearer JWT_ACCESS_TOKEN'
```

Parameters:
- `transfer_status`: Filter by transfer status
- `archive`: Filter archived/non-archived orders
- `limit`: Number of records per page
- `page`: Page number

### Cancel Order
Cancels an existing order by its ID.

```bash
curl --location --request PUT 'http://localhost:8080/api/v1/orders/DA2501316CYUOG/cancel' \
--header 'Authorization: Bearer JWT_ACCESS_TOKEN'
```

### Logout
Invalidates the current access token.

```bash
curl --location --request POST 'http://localhost:8080/api/v1/logout' \
--header 'Authorization: Bearer JWT_ACCESS_TOKEN'
```

### Refresh Token
Obtains a new access token using a refresh token.

```bash
curl --location 'http://localhost:8080/api/v1/token/refresh' \
--header 'Content-Type: application/json' \
--data '{
    "refresh_token": "REFRESH_TOKEN"
}'
```



