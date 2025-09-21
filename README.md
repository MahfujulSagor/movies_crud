# 🎬 Movies CRUD

[![Go](https://img.shields.io/badge/Go-1.25.1-blue?logo=go\&logoColor=white)](https://golang.org/)
[![SQLite](https://img.shields.io/badge/SQLite-3.42-orange?logo=sqlite\&logoColor=white)](https://www.sqlite.org/index.html)
[![MIT License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

A **CRUD application in Go with SQLite** for managing movies, directors, and casts.
Supports insertion, retrieval, updating, and deletion of movies along with director and cast details.

---

## Features

* Insert new movies with directors and casts
* Update existing movies using full JSON objects
* Retrieve movies by ID or list with pagination
* Delete movies by ID
* Prevent duplicate directors while allowing multiple casts
* Atomic operations using SQLite transactions

---

## Schema

**Directors table**

```sql
CREATE TABLE directors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    age INTEGER NOT NULL
);
```

**Casts table**

```sql
CREATE TABLE casts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    actor TEXT NOT NULL,
    actress TEXT NOT NULL
);
```

**Movies table**

```sql
CREATE TABLE movies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    rating INTEGER NOT NULL,
    director_id INTEGER,
    cast_id INTEGER,
    FOREIGN KEY(director_id) REFERENCES directors(id),
    FOREIGN KEY(cast_id) REFERENCES casts(id)
);
```

---

## Types

```go
type Movie struct {
    ID       int64     `json:"id"`
    Title    string    `json:"name" validate:"required"`
    Rating   int       `json:"rating" validate:"required,gte=0,lte=10"`
    Director *Director `json:"director"`
    Cast     *Cast     `json:"cast"`
}

type Director struct {
    ID   int64  `json:"id"`
    Name string `json:"name" validate:"required"`
    Age  int    `json:"age" validate:"required,gte=0,lte=110"`
}

type Cast struct {
    ID      int64  `json:"id"`
    Actor   string `json:"actor" validate:"required"`
    Actress string `json:"actress" validate:"required"`
}
```

---

## Setup Instructions

1. **Clone the repository**

```bash
git clone https://github.com/MahfujulSagor/movies_crud.git
cd movies_crud
```

2. **Install Go dependencies**

```bash
go mod tidy
```

3. **Run the project**

```bash
go run cmd/movies/main.go
```

Server runs at:
👉 `http://localhost:8080`

## If you dont setup `.env`

### Run the server like this

```bash
go run cmd/movies/main.go -config config/config.yaml
```

Server runs at:
👉 `http://localhost:8080`

## 📂 Project Structure

```markdown
movies_crud/
├── cmd/               # Application entry point
├── config/            # Configuration file (ignored in Git)
├── db/                # Database (ignored in Git)
├── internal/
│   ├── config/        # Configuration loading (env, YAML)
│   ├── db/            # SQLite database logic
│   ├── logger/        # Centralized logging
│   ├── response/      # JSON response helpers
│   ├── student/       # Student handlers
│   └── types/         # Domain models
├── logs/              # Log output (ignored in Git)
├── .env               # Environment variables (ignored in Git)
└── go.mod
```

## ⚙️ Configuration

The app is configured via a YAML configuration file (e.g., config.yaml).

Example config.yaml:

```yaml
env: "development"
db_path: "db/movies.db"
http:
  host: "localhost"
  port: 8080
logging:
  level: "debug"
  file: "logs/app.log"
```

Example `.env` file:

```env
CONFIG_PATH="config/config.yaml"
```

---

## 📡 API Endpoints

### Health Check

```bash
GET /
```

### Students

| Method   | Endpoint             | Description                     |
| -------- | -------------------- | ------------------------------- |
| `POST`   | `/api/v1/movies`      | Create a new movie            |
| `GET`    | `/api/v1/movies`      | List movies (with pagination) |
| `GET`    | `/api/v1/movies/{id}` | Get movie by ID               |
| `PUT`    | `/api/v1/movies/{id}` | Update movie by ID            |
| `DELETE` | `/api/v1/movies/{id}` | Delete movie by ID            |

---

## 📖 Example Request / Response

### Create Movie

```http
POST /api/v1/movies
Content-Type: application/json

{
    "id": 1,
    "name": "Interstellar",
    "rating": 9,
    "director": {
        "name": "Christopher Nolan",
        "age": 54
    },
    "cast": {
        "actor": "Matthew McConaughey",
        "actress": "Anne Hathaway"
    }
}
```

Response:

```json
{
  "success": "OK",
  "message": "Movie created with ID: 1"
}
```

## Example Movie JSON

```json
{
    "id": 1,
    "name": "Interstellar",
    "rating": 9,
    "director": {
        "name": "Christopher Nolan",
        "age": 54
    },
    "cast": {
        "actor": "Matthew McConaughey",
        "actress": "Anne Hathaway"
    }
}
```

---

## Seeding Examples

```go
movies := []types.Movie{
    {
        Title:  "Interstellar",
        Rating: 9,
        Director: &types.Director{
            Name: "Christopher Nolan",
            Age:  54,
        },
        Cast: &types.Cast{
            Actor:   "Matthew McConaughey",
            Actress: "Anne Hathaway",
        },
    },
    {
        Title:  "Inception",
        Rating: 9,
        Director: &types.Director{
            Name: "Christopher Nolan",
            Age:  54,
        },
        Cast: &types.Cast{
            Actor:   "Leonardo DiCaprio",
            Actress: "Elliot Page",
        },
    },
    {
        Title:  "The Matrix",
        Rating: 9,
        Director: &types.Director{
            Name: "Lana Wachowski",
            Age:  56,
        },
        Cast: &types.Cast{
            Actor:   "Keanu Reeves",
            Actress: "Carrie-Anne Moss",
        },
    },
}
```

Insert them using:

```go
for _, m := range movies {
    _, err := db.CreateMovie(m.Title, m.Rating, m.Director, m.Cast)
    if err != nil {
        log.Fatalf("Error seeding movie: %v", err)
    }
}
```

## 🛠 Development Notes

- Logs are stored in `logs/app.log`
- In **development**, logs also print to console
- SQLite DB file defaults to `movies.db` in the project root
- Graceful shutdown ensures ongoing requests complete within 10s

---

## 📜 License

This project is licensed under the **MIT License**. See [LICENSE](LICENSE) for details.
