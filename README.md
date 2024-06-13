# Student API

This is a simple REST API built in Go that performs basic CRUD (Create, Read, Update, Delete) operations on a list of students. Each student has the following attributes: `ID`, `Name`, `Age`, and `Email`. Additionally, the API integrates with the Ollama API to generate a summary of a student's profile.

## Requirements

1. **Initialize a Go module:** Ensure your project is a Go module by initializing it with `go mod init`.
2. **REST API Endpoints:**
   - Create a new student: `POST /students`
   - Get all students: `GET /students`
   - Get a student by ID: `GET /students/{id}`
   - Update a student by ID: `PUT /students/{id}`
   - Delete a student by ID: `DELETE /students/{id}`
   - Generate a summary of a student by ID using Ollama: `GET /students/{id}/summary`
3. **Data Storage:** Used an in-memory data structure (e.g., a slice or map) to store student information. No need for a database.
4. **Ollama Integration:** Use the Ollama API to generate a summary of a student's profile.
5. **Error Handling:** Handle errors appropriately, such as when a student with a specified ID does not exist.
6. **Input Validation:** Ensure that the input data for creating and updating students is valid.
7. **Concurrency:** Ensure your API can handle concurrent requests safely.

## Getting Started

### Prerequisites

- Go 1.18 or higher
- Internet connection to use the Ollama API

### Installing

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/student-api.git
   cd student-api
   ```

2. Install Dependency

    ```sh
   go get github.com/gofiber/fiber/v2
   ```

> **_NOTE:_** Make sure llama3 model running on your system. For install [Ollama](https://github.com/ollama/ollama/blob/main/README.md)


### Running API

Run the Go program to start the server:

```sh
go run main.go
```

### API Endpoints

#### 1. Create Student

```sh
POST /students
```

- Request Body

    ```sh
    {
        "name": "test",
        "age": 23,
        "email": "test@example.com"
    }
    ```

#### 2. Get All Students

```sh
GET /students
```

#### 3. Get Student By ID

```sh
GET /students/{id}
```

#### 4. Update Student

```sh
POST /students/{id}
```

- Request Body

    ```sh
    {
        "name": "test updated",
        "age": 25,
        "email": "testupdate@example.com"
    }
    ```

#### 5. Delete Student By ID

```sh
GET /students/{id}
```

#### 6. Generate Student Summary

```sh
GET /students/{id}/summary
```


### Concurrency

The API uses a sync.Mutex to ensure safe concurrent access to the in-memory data structure storing student information.

### Error Handling

The API returns appropriate error messages and status codes when a student is not found or when input validation fails.

### Input Validation

The API validates the input data for creating and updating students. It checks for:

- Non-empty Name
- Positive Age
- Non-empty Email