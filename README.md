# eCommerce API

This is a backend eCommerce API built with Golang and the Gin framework. The API supports product management, user authentication, role-based authorization, and order handling. It also integrates with Amazon S3 for media storage and uses PostgreSQL as the primary database.

## Table of Contents
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Endpoints](#endpoints)
- [Authentication and Authorization](#authentication-and-authorization)
- [Contributing](#contributing)
- [License](#license)

---

### Features
- **Product Management**: CRUD operations for products, restricted to admin users.
- **User Authentication**: Register and login with JWT-based authentication.
- **Role-Based Access Control**: Only admin users can perform certain actions, such as product creation and deletion.
- **Order Management**: Allows users to place orders and view order history.
- **AWS S3 Integration**: Media files are uploaded to and stored in Amazon S3.
- **Swagger Documentation**: API documentation available through Swagger UI.

---

### Tech Stack
- **Language**: Go (Golang)
- **Framework**: Gin
- **Database**: PostgreSQL
- **Storage**: Amazon S3 for media files
- **Authentication**: JWT for secure token-based authentication

---

### Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-username/ecommerce-api.git
   cd ecommerce-api
| Endpoint               | Method | Description                                   | Access       |
|------------------------|--------|-----------------------------------------------|--------------|
| `/api/v1/auth/register` | POST   | Register a new user                          | Public       |
| `/api/v1/auth/login`    | POST   | Log in and receive a JWT                     | Public       |
| `/api/v1/products`      | POST   | Create a new product                         | Admin only   |
| `/api/v1/products`      | GET    | List all products                            | Public       |
| `/api/v1/products/:id`  | PUT    | Update a product by ID                       | Admin only   |
| `/api/v1/products/:id`  | DELETE | Delete a product by ID                       | Admin only   |
| `/api/v1/orders`        | POST   | Create a new order                           | User only    |
| `/api/v1/orders/:id`    | GET    | View order details                           | User only    |
