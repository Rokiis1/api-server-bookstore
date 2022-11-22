package main

// Go mod tidy
import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"golang.org/x/mod/sumdb/storage"
	"gorm.io/gorm"
)

type Book struct{
	Author      string     `json:"author"`
	Title       string     `json:"title"`
	Publisher   string     `json:"publisher"`
}

type Repository struct{
	DB *gorm.DB
}

func (r *Repository)CreateBook(context *fiber.Ctx) error {
	book := Book{}

	err := context.BodyParser(&book)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
			return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not create book"})
		    return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message":"book created successfully"})
    return nil

}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	err := r.DB.Find(booksModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not find books"})
			return err		
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"books fetched successfully",
		"data":bookModels,
	})
	return nil

} 

func(r *Repository) SetupRoutes(app *fiber.App){
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_books/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main(){
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	db, err := storage.NewConnection(config)

    if err != nil {
		log.Fatal("could not load the database", err)
	}

    r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":4000")
}