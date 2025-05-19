# Cooperative Product Sales and Inventory Management System

A web-based management system developed for cooperatives, capable of product sales and inventory tracking.

## Features

- **User Management**: Role-based authorization system (Admin, Member, Sales, Warehouse)
- **Product Management**: Product listing, adding, updating, deleting
- **Inventory Tracking**: Input-output operations, location-based inventory management
- **Sales Management**: Sales forms, invoicing, return operations
- **Customer Management**: Customer registration, debt and payment tracking
- **Supplier Management**: Supplier registration, purchase records
- **Reporting**: Sales, inventory, and financial reports
- **User-Friendly Interface**: Responsive design, modern UI

## Technology

- **Backend**: Go (Golang) and Gin Framework
- **Frontend**: HTML5, TailwindCSS, Alpine.js
- **Database**: SQLite (can be migrated to PostgreSQL)
- **Authentication**: JWT token based

## Installation

### Requirements

- Go 1.16 or higher
- SQLite

### Steps

1. Clone the repository:
```bash
git clone https://github.com/kullanici/koopsatis.git
cd koopsatis
```

2. Install required dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o koopsatis ./cmd/server
```

4. Run the application:
```bash
./koopsatis
```

5. Navigate to the following address in your browser:
```
http://localhost:8080
```

## Development Environment

For development, you can use the following command:

```bash
go run ./cmd/server/main.go
```

## Default User

When the application is first run, an admin user is automatically created:

- **Username**: admin
- **Password**: admin123

It is recommended to change this password after the first login for security purposes.

## API Documentation

RESTful APIs can be accessed through the following base URL:

```
http://localhost:8080/api
```

### Authentication Endpoints

- `POST /api/auth/login`: User login
- `POST /api/auth/register`: New user registration

### User Endpoints

- `GET /api/users`: List all users (admin only)
- `GET /api/users/:id`: Get a specific user
- `PUT /api/users/:id`: Update user information
- `DELETE /api/users/:id`: Delete a user (admin only)

### Product Endpoints

- `POST /api/products`: Add new product
- `GET /api/products`: List all products
- `GET /api/products/:id`: Get a specific product
- `PUT /api/products/:id`: Update product information
- `DELETE /api/products/:id`: Delete a product

