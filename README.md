# 🧵 Forum Web Application Project

## 📚 Overview

This project consists of creating a web-based forum with the following functionalities:

- Communication between users via **posts** and **comments**.
- Association of **categories** to posts.
- Ability to **like** or **dislike** posts and comments.
- Ability to **filter posts** by different criteria.

---

## 🗄️ Database: SQLite

All forum data (users, posts, comments, etc.) will be stored using **SQLite**, a lightweight and efficient embedded SQL database engine.

### 🔧 Requirements:
- Use at least:
  - One `SELECT` query
  - One `CREATE` query
  - One `INSERT` query

> 💡 **Tip:** Design an **Entity-Relationship Diagram (ERD)** before implementation for better structure and performance.

📖 [Learn more about SQLite](https://www.sqlite.org/index.html)

---

## 🔐 Authentication System

Users must register and log in to post or comment.

### ✅ Registration Requirements:
- **Email**
  - Must be unique.
  - Return an error if already taken.
- **Username**
- **Password**
  - 🔐 (Bonus) Encrypt password before storing it using **bcrypt**.

### 🔐 Session Management:
- Use **cookies** to track login sessions.
- Each session must have an **expiration time**.
- 🏅 (Bonus) Use **UUIDs** for session identification.

---

## 💬 User Communication

Registered users will be able to create posts and comments.

### 📝 Posts:
- Must be associated with one or more **categories** (user-defined).
- Visible to all users (including guests).

### 💭 Comments:
- Only registered users can comment.
- Visible to all users.

> 🔒 Non-registered users have read-only access to posts and comments.

---

## 👍 Likes & Dislikes

- Only **registered users** can like/dislike posts and comments.
- The number of likes/dislikes must be visible to **all users**.

---

## 🔍 Filtering Mechanism

Allow filtering of posts based on:

- **Categories** (acts like subforums)
- **Created posts** (by the logged-in user)
- **Liked posts** (by the logged-in user)

> 📝 The last two filters are only available to authenticated users.

---

## 🐳 Docker Integration

You **must** dockerize the project using **Docker**.

📘 Refer to: _ascii-art-web-dockerize_ subject for basic Docker usage.

### Docker Concepts to Apply:
- Containerization of the full application.
- Dependency management and compatibility.
- Docker image creation and running containers.

---

## ⚙️ General Instructions

- Use **SQLite** as the database.
- Handle **HTTP status** and **website errors** gracefully.
- Handle **technical errors** (invalid queries, sessions, cookies, etc.).
- Follow **good programming practices** (code quality, modularity, documentation).
- (Recommended) Write **unit tests** for critical components.

---

## 📦 Allowed Go Packages

You may use the following Go packages:

- Standard library packages.
- `github.com/mattn/go-sqlite3`
- `golang.org/x/crypto/bcrypt`
- `github.com/gofrs/uuid` **or** `github.com/google/uuid`

> ❌ **Frontend frameworks/libraries** such as React, Angular, or Vue are **not allowed**.

---

## 🎯 Learning Objectives

This project will help you understand and practice:

### 🌐 Web Fundamentals:
- HTML
- HTTP protocol
- Sessions and cookies

### 🐳 Docker Basics:
- Application containerization
- Dependency isolation
- Image creation

### 🛢️ SQL Language:
- Creating and manipulating databases
- Writing SQL queries (CRUD operations)

### 🔒 Security:
- User authentication
- Password encryption (Bonus)

---

## ✅ Final Notes

Make sure your application is:

- Functional and user-friendly
- Secure
- Efficient
- Follows best practices
- Fully containerized with Docker

> 🚀 Build, Test, and Ship it with confidence!
