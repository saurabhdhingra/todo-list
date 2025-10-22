# Go RESTful Todo List API with JWT Authentication

<img width="1833" height="780" alt="todo-list-api-bsrdd" src="https://github.com/user-attachments/assets/b238b0b4-b61f-40ae-b006-5b7b08c523e9" />

This project implements a secure, robust, and fully functional RESTful API for managing a user's to-do list, built with **Go (Golang)**, the **Gin framework**, and **PostgreSQL** via the **GORM** ORM.

It emphasizes best practices, including user authentication with JWT, secure password storage using bcrypt, and robust error handling with proper HTTP status codes.

## Features

* **User Authentication:** Secure registration and login using JWT (JSON Web Tokens).
* **Secure Passwords:** All passwords are hashed using bcrypt before storage.
* **CRUD Operations:** Full functionality for creating, reading, updating, and deleting todo items.
* **Authorization:** Ensures users can only access or modify their own todo items.
* **Pagination & Filtering:** Supports fetching todo lists with specified limits, page numbers, and optional filtering by status (e.g., done/not\_done).
* **RESTful Design:** Clean, resource-oriented endpoints.

***

## Prerequisites

Before running the project, you must have the following installed:

1.  **Go:** Version 1.21 or higher.
2.  **PostgreSQL:** A running instance of a PostgreSQL server.
3.  **Git:** For cloning the repository.

***

## Getting Started

### 1. Clone the Repository

```bash
git clone [https://github.com/saurabhdhingra/todo-list.git](https://github.com/saurabhdhingra/todo-list.git)
cd todo-list
```

2. **Configure Database**

Edit the database connection string in config/database.go to match your local PostgreSQL setup.

config/database.go Snippet:

```
// Replace with your actual database credentials
dsn := "host=localhost user=user password=password dbname=todoapi port=5432 sslmode=disable TimeZone=Asia/Shanghai" 
```

3. **Install Dependencies**

Run the following command to download all required Go packages:

Bash
```
go mod tidy
```

4. **Run the Application**

Start the API server:

Bash
```
go run main.go
```
The server will start on port 8080.

Server listening on :8080
Database connection successful and migrations complete.

## API Endpoints
The base URL for all endpoints is http://localhost:8080.

### Authentication

```
Endpoint	        Method	    Body Parameters	        Description
/register	        POST	    name, email, password	Create a new user. Returns a JWT token.
/login	            POST	    email, password	        Authenticate and return a JWT token.
```

### To-Do Management (Requires Authorization: Bearer <TOKEN>)

```
Endpoint	        Method	    Body/Query Parameters	        Description
/todos	            POST	    title, description	            Create a new to-do item.
/todos	            GET	        ?page=1&limit=10&status=done	Retrieve paginated list of todos.
/todos/:id	        PUT	        title, description, done	    Update an existing todo item. Must own the item.
/todos/:id	        DELETE	    (None)	                        Delete a todo item. Must own the item.
```

## Security & Authorization
All /todos endpoints are protected by the AuthMiddleware.

### How to Authenticate Requests

Call the /login or /register endpoint to receive a JWT string.

For subsequent requests to protected routes, include the JWT in the HTTP header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Authorization (Ownership)

The API enforces ownership. A user can only PUT or DELETE a todo item if the user_id stored in the JWT matches the UserID associated with the todo item in the database.

Attempting to access a todo that doesn't belong to the user results in a 403 Forbidden response.

## Project Structure
The project follows a modular structure for better separation of concerns:

```
todo-list/
├── main.go               # Application entry point, router setup
├── config/
│   └── database.go       # Database connection (PostgreSQL/GORM)
├── models/
│   ├── user.go           # GORM User model
│   └── todo.go           # GORM Todo model
├── handlers/
│   ├── auth.go           # Handlers for Register and Login
│   └── todo.go           # Handlers for Todo CRUD (Create, Get, Update, Delete)
├── middleware/
│   └── auth.go           # JWT validation middleware
└── utils/
    └── auth.go           # JWT token generation and validation logic
```

## Acknowledgement
https://roadmap.sh/projects/todo-list-api
