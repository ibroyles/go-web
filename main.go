package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// note represents data about a record note.
type note struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Author      string `json:"author"`
	LastUpdated string `json:"lastUpdated"`
}

type Error struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// notes slice to seed record note data.
var notesData []note
var errors []Error
var users []User
var currentUser User

func init() {
	file, err := os.Open("notes.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&notesData)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
}

func IndentedJSON(c *gin.Context, code int, obj interface{}) {
	c.IndentedJSON(code, obj)
}

func preflightCheck(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	c.JSON(http.StatusOK, gin.H{"message": "set headers"})
}

func main() {
	router := gin.Default()
	router.SetHTMLTemplate(template.Must(template.ParseGlob("static/templates/*.html")))
	router.GET("/notes", getNotes)
	router.GET("/notes/:id", getNoteByID)
	router.POST("/notes", postNotes)
	router.DELETE("/notes/:id", deleteNote)
	router.PATCH("/notes/:id", patchNote)
	//router.GET("/", getUser)
	//router.POST("/", postUser)

	//router.GET("/home", viewNotes)
	//router.GET("/errors", getErrors)
	//router.GET("/errors/:id", deleteError)
	//router.GET("/edit", editNote)
	//router.OPTIONS("/notes", preflightCheck)

	log.Fatal(router.Run("localhost:8080"))
}

// getNotes responds with the list of all notes as JSON.
func getNotes(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Access-Control-Allow-Methods, Authorization")
	//c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	IndentedJSON(c, http.StatusOK, notesData)
}

// postNotes adds an note from JSON received in the request body.
func postNotes(c *gin.Context) {
	var newNote note
	var newError Error

	// Call BindJSON to bind the received JSON to
	// newNote.
	if err := c.BindJSON(&newNote); err != nil {
		IndentedJSON(c, http.StatusBadRequest, "Bad Request")
		return
	}

	if newNote.Title == "" || newNote.Body == "" || newNote.Author == "" {
		newError.Message = "One or more required fields are missing"
		newError.ID = strconv.Itoa(len(errors) + 1)
		errors = append(errors, newError)
		IndentedJSON(c, http.StatusBadRequest, gin.H{
			"Note": note{Title: "bad request", Body: "Bad request body", Author: "Author"},
		})
		return
	}

	newNote.ID = strconv.Itoa(len(notesData) + 1)
	newNote.LastUpdated = time.Now().String()
	//Add the new note to the slice.
	notesData = append(notesData, newNote)
	IndentedJSON(c, http.StatusOK, gin.H{"message": "ok"})
	IndentedJSON(c, http.StatusCreated, newNote)
	//TODO write notesData to file
	file, _ := os.Create("notes.json")
	if err := json.NewEncoder(file).Encode(notesData); err != nil {
		panic("Error opening json file to write to")
	}
}

// getNoteByID locates the note whose ID value matches the id
// parameter sent by the client, then returns that note as a response.
func getNoteByID(c *gin.Context) {
	id := c.Param("id")

	// Loop through the list of notes, looking for
	// an note whose ID value matches the parameter.
	for _, a := range notesData {
		if a.ID == id {
			IndentedJSON(c, http.StatusOK, gin.H{
				"id":          a.ID,
				"title":       a.Title,
				"body":        a.Body,
				"author":      a.Author,
				"lastUpdated": a.LastUpdated,
			})
			return
		}
	}
	IndentedJSON(c, http.StatusNotFound, gin.H{"message": "note not found"})
}

func deleteNote(c *gin.Context) {
	id := c.Param("id")

	for i, a := range notesData {
		if a.ID == id {
			notesData = append(notesData[:i], notesData[i+1:]...)
		}
	}
}

func patchNote(c *gin.Context) {
	id := c.Param("id")

	var newNote note

	for i, a := range notesData {
		if a.ID == id {
			//newNote.LastUpdated = time.Now().String()
			//newNote.ID = a.ID
			//newNote.Author = a.Author
			newNote = a
			notesData = append(notesData[:i], notesData[i+1:]...)
		}
	}

	err := c.BindJSON(&newNote)
	if err != nil {
		panic("Bad Request")
	}

	notesData = append(notesData, newNote)

	IndentedJSON(c, http.StatusOK, gin.H{"message": "ok", "note": newNote})

}

//func viewNotes(c *gin.Context) {
//
//	var usersNotes []note
//	var user User
//
//	if err := c.Bind(&user); err != nil {
//		c.String(http.StatusBadRequest, "Bad Request")
//		return
//	}
//
//	for _, a := range users {
//		if a.Username == user.Username {
//			currentUser = user
//		}
//	}
//
//	for _, a := range notesData {
//		if a.Author == currentUser.Username {
//			usersNotes = append(usersNotes, a)
//		}
//	}
//
//	c.HTML(http.StatusOK, "home.html", gin.H{
//		"Notes": usersNotes,
//		"User":  currentUser,
//	})
//}

func getErrors(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"Errors": errors,
	})
}

func deleteError(c *gin.Context) {
	id := c.Param("id")

	for i, a := range errors {
		if a.ID == id {
			errors = append(errors[:i], errors[i+1:]...)
			IndentedJSON(c, http.StatusOK, errors)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error not found"})

}

func editNote(c *gin.Context) {
	//TODO: Implement authentication here
	c.HTML(http.StatusOK, "edit.html", gin.H{
		"Author": currentUser.Username,
	})
}

//func getUser(c *gin.Context) {
//	c.HTML(http.StatusOK, "login.html", gin.H{})
//}

//func postUser(c *gin.Context) {
//
//	var user User
//	var usersNotes []note
//
//	if err := c.Bind(&user); err != nil {
//		c.String(http.StatusBadRequest, "Bad Request")
//		return
//	}
//
//	if user.Username == "" {
//		c.String(http.StatusBadRequest, "Username is missing")
//	}
//
//	for _, a := range users {
//		if a.Username == user.Username {
//			currentUser = user
//		}
//	}
//
//	user.ID = strconv.Itoa(len(users) + 1)
//	currentUser = user
//	users = append(users, user)
//
//	for _, a := range notesData {
//		if a.Author == currentUser.Username {
//			usersNotes = append(usersNotes, a)
//		}
//	}
//	c.HTML(http.StatusCreated, "home.html", gin.H{
//		"User":  user,
//		"Notes": usersNotes,
//	})
//
//}
