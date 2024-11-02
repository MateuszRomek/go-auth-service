# Go Authentication Service

A simple authentication service built with Go, designed for learning purposes. This project demonstrates basic authentication and session management implementation in Go.

## Overview

This service provides essential authentication functionality including:

- User registration and login
- Session management
- Password hashing and security
- Database migrations

## Features

- User registration with secure password hashing
- User login with session management
- Database-backed session storage
- Clean architecture with separation of concerns
- SQL migrations for database schema management

## Getting Started

### Prerequisites

- Go 1.23 or higher
- PostgreSQL database

### Installation

1. Clone the repository

2. Install dependencies

```bash
go mod download
```

3. Run database migrations

```bash
make migrate-up
```

4. Start the server

```bash
air
```

## API Endpoints

- `POST /register` - Register a new user
- `POST /login` - Login user
- `POST /logout` - Logout user
- `GET /me` - Get current user information
