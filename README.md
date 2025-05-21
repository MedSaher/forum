# 🧩 Web Forum Project

A web-based discussion platform built using **Go** and **SQLite**. This forum allows users to engage through posts and comments, react with likes and dislikes, and filter posts by different criteria. It is fully containerized using Docker.

---

## 📌 Objectives

This project aims to deliver a functional web forum with the following capabilities:

- 🗣️ User communication via **posts and comments**
- 🗂️ Categorization of posts
- 👍👎 Like and dislike functionality for posts and comments
- 🔍 Filtering of posts by:
  - Category
  - Posts created by the logged-in user
  - Posts liked by the user
- 🔐 Full authentication system using cookies and sessions
- 🔐 (Bonus) Password encryption with bcrypt
- 📦 Containerization with Docker

---

## 🛠️ Technologies Used

- **Language**: Go (Golang)
- **Database**: SQLite
- **Authentication**: HTTP cookies, bcrypt (bonus)
- **Session Management**: Cookie-based, with expiration
- **UUID** (Bonus): `google/uuid` or `gofrs/uuid`
- **Containerization**: Docker
- **Frontend**: Pure HTML and CSS (No JS frameworks)

---

## 🔐 Authentication Features

- User **registration** with:
  - Email (must be unique)
  - Username
  - Password
- Encrypted password storage (bcrypt - bonus)
- Session handling with cookie expiration
- Only one active session per user allowed

---

## ✍️ Communication

- Only **authenticated users** can create posts and comments
- Posts can be assigned to one or more **categories**
- Comments are linked to individual posts
- Posts and comments are **publicly viewable** to all users

---

## 👍👎 Likes and Dislikes

- Available only for **registered users**
- Users can like/dislike both **posts and comments**
- Number of likes/dislikes is visible to **everyone**

---

## 🔎 Filtering Mechanism

- Users can filter posts by:
  - **Category** (like subforums)
  - **Created posts** (by logged-in user)
  - **Liked posts** (by logged-in user)

---

## 🗄️ SQLite Requirements

The forum uses SQLite to persist data. Your database must implement:

- At least one `CREATE` query
- At least one `INSERT` query
- At least one `SELECT` query

Entities include:
- Users
- Sessions
- Posts
- Comments
- Categories
- Likes/Dislikes

💡 It is highly recommended to design an **ER diagram** before implementing your schema.

---

## 🐳 Docker Integration

The project is fully containerized using Docker.

### Build and Run:

```bash
docker build -t forum-app .
docker run -p 8080:8080 forum-app
