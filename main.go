package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Student struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

type OllamResponse struct {
	Response string `json:"response"`
}

var (
	students []Student
	mu       sync.Mutex
)

func CreateStudent(c *fiber.Ctx) error {
	var student Student
	c.BodyParser(&student)

	if student.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": errors.New("invalid name").Error(),
		})
	}

	if student.Age < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": errors.New("invalid age").Error(),
		})
	}

	if student.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": errors.New("invalid email").Error(),
		})
	}

	mu.Lock()
	student.ID = strconv.Itoa(len(students) + 1)
	students = append(students, student)
	mu.Unlock()

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"data": student,
	})
}

func GetStudents(c *fiber.Ctx) error {
	mu.Lock()
	defer mu.Unlock()
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": students,
	})
}

func GetStudentByID(c *fiber.Ctx) error {
	id := c.Params("id")
	mu.Lock()
	defer mu.Unlock()
	for _, item := range students {
		if item.ID == id {
			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"data": item,
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
		"error": errors.New("student not found").Error(),
	})
}

func UpdateStudent(c *fiber.Ctx) error {
	id := c.Params("id")
	var student Student
	c.BodyParser(&student)

	if student.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": errors.New("invalid name").Error(),
		})
	}

	if student.Age < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": errors.New("invalid age").Error(),
		})
	}

	if student.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": errors.New("invalid email").Error(),
		})
	}

	mu.Lock()
	defer mu.Unlock()
	for index, item := range students {
		if item.ID == id {
			students = append(students[:index], students[index+1:]...)
			student.ID = id
			students = append(students, student)
			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"data": student,
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
		"error": errors.New("student not found"),
	})
}

func DeleteStudent(c *fiber.Ctx) error {
	id := c.Params("id")

	mu.Lock()
	defer mu.Unlock()
	for index, item := range students {
		if item.ID == id {
			students = append(students[:index], students[index+1:]...)
			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"data": item,
			})
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
		"error": errors.New("student not found"),
	})
}

func GetStudentSummary(c *fiber.Ctx) error {
	id := c.Params("id")
	var student *Student

	mu.Lock()
	for _, item := range students {
		if item.ID == id {
			student = &item
			break
		}
	}
	mu.Unlock()

	if student == nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"error": errors.New("student not found"),
		})
	}

	summary, err := GenerateSummary(*student)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": errors.New("error generating summary"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"summary": &summary,
	})
}

func GenerateSummary(student Student) (*string, error) {
	agent := fiber.Post("http://localhost:11434/api/generate")

	type ReqBody struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		Stream bool   `json:"stream"`
	}

	byteStudentArray, err := json.Marshal(student)
	if err != nil {
		return nil, err
	}

	prompt := "generate summary for given student\n" + string(byteStudentArray)

	var reqBody = &ReqBody{
		Model:  "llama3",
		Prompt: prompt,
		Stream: false,
	}

	byteArray, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	agent.Body(byteArray)
	_, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return nil, err
	}

	var res OllamResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res.Response, nil
}

func main() {
	app := fiber.New()

	api := app.Group("/students")
	api.Post("/", CreateStudent)
	api.Get("/", GetStudents)
	api.Get("/:id", GetStudentByID)
	api.Put("/:id", UpdateStudent)
	api.Delete("/:id", DeleteStudent)
	api.Get("/:id/summary", GetStudentSummary)

	app.Listen(":3000")
}
