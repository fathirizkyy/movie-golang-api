# 📦 Gin Golang CRUD API with JWT Authentication and File Upload

This is a RESTful API built using [Gin](https://github.com/gin-gonic/gin) and [GORM](https://gorm.io/) in Golang. It provides full CRUD functionality for "Post" resources with support for `.jpg` and `.jpeg` image uploads, and JWT-based authentication for securing endpoints.

---

## 🧱 Features

- ✅ CRUD Post API
- ✅ Image upload (`.jpg`, `.jpeg`)
- ✅ User registration & login with hashed passwords
- ✅ JWT Authentication (middleware-protected routes)
- ✅ Log management with Logrus
- ✅ Password hashing with bcrypt

---

## 🛠️ Tech Stack

- [Gin](https://github.com/gin-gonic/gin) - HTTP Web Framework
- [GORM](https://gorm.io/) - ORM for Golang
- [MySQL](https://www.mysql.com/) - Database
- [JWT](https://github.com/golang-jwt/jwt) - Token-based authentication
- [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - Password hashing
- [Logrus](https://github.com/sirupsen/logrus) - Logging

---

## 🧩 Table Structure

### 📌 Post Table

| Column        | Type   | Description             |
| ------------- | ------ | ----------------------- |
| `id`          | uint   | Primary Key             |
| `title`       | string | Title of the post       |
| `image`       | string | File path of the image  |
| `description` | string | Description of the post |

### 📌 User Table

| Column     | Type   | Description      |
| ---------- | ------ | ---------------- |
| `id`       | uint   | Primary Key      |
| `name`     | string | Name of the user |
| `email`    | string | Must be unique   |
| `password` | string | Hashed password  |

---

## 🔐 Authentication

### Register

- **Endpoint**: `POST /register`
- **Payload**:

```json
{
  "name": "Your Name",
  "email": "you@example.com",
  "password": "yourpassword"
}
```

#### Validation

Name: required
Email: required, must be valid
Password: required, min 6 characters

### Login

- **Endpoint**: `POST /login`
- **Payload**:

```json
{
  "email": "you@example.com",
  "password": "yourpassword"
}
```

Returns JWT token if login is successful.

## 📘 Post API Endpoints

| Method | Endpoint        | Auth Required | Description                         |
| ------ | --------------- | ------------- | ----------------------------------- |
| GET    | `/api/post`     | ❌ No         | Get all posts                       |
| GET    | `/api/post/:id` | ✅ Yes        | Get post by ID                      |
| POST   | `/api/post`     | ✅ Yes        | Create new post (with image upload) |
| PUT    | `/api/post/:id` | ✅ Yes        | Update existing post                |
| DELETE | `/api/post/:id` | ✅ Yes        | Delete post                         |

## 🚀 Running the Project


## 📦 Installation

Install the required Go packages using the following commands:

```bash
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get github.com/golang-jwt/jwt/v5
go get github.com/sirupsen/logrus
go get golang.org/x/crypto/bcrypt
```

## 🔧 Configuration

- Make sure your MySQL server is running.
- Create an `uploads/` directory to store uploaded images.

## ▶️ Run the Server

Start the server with the following command:

```bash
go run main.go
```

http://localhost:8000

## 📁 Image Upload

- **Endpoint**: `POST /api/post`
- **Content-Type**: `multipart/form-data`
- **Form Fields**:
  - `title`: string
  - `description`: string
  - `image`: file (`.jpg` or `.jpeg`)

## 🔒 Authorization Middleware

- Only the route `GET /api/post` is public.
- All other routes require a valid JWT token in the `Authorization` header:

```http
Authorization: Bearer <token>

```

## 🙋‍♂️ Author

Developed by [fathirizkyy](https://github.com/fathirizkyy)
