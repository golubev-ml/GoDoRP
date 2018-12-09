package main

import (
	"encoding/json"
	"fmt"
	"godorpdb"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

// CRUD Route Handlers
func createPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setCors(w)
	decoder := json.NewDecoder(r.Body)
	var newPost godorpdb.Post
	if err := decoder.Decode(&newPost); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	godorpdb.DB.Create(&newPost)
	res, err := json.Marshal(newPost)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(res)
}

func deletePostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setCors(w)
	var deletedPost godorpdb.Post
	godorpdb.DB.Where("ID = ?", ps.ByName("postId")).Delete(&deletedPost) // write now this returns a blank item not the deleted item
	res, err := json.Marshal(deletedPost)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.Write(res)
}

func updatePostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setCors(w)
	type body struct {
		Author  string
		Message string
	}
	var updates body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updates); err != nil {
		http.Error(w, err.Error(), 400)
	}

	var updatedPost godorpdb.Post
	godorpdb.DB.Where("ID = ?", ps.ByName("postId")).First(&updatedPost)
	updatedPost.Author = updates.Author
	updatedPost.Message = updates.Message
	godorpdb.DB.Save(&updatedPost)
	res, err := json.Marshal(updatedPost)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.Write(res)
}

func showPostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setCors(w)
	var post godorpdb.Post
	godorpdb.DB.Where("ID = ?", ps.ByName("postId")).First(&post)
	res, err := json.Marshal(post)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(res)
}

func indexPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setCors(w)
	var posts []godorpdb.Post
	godorpdb.DB.Find(&posts)
	res, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(res)
}

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setCors(w)
	fmt.Fprintf(w, "This is the RESTful api")
}

// used for COR preflight checks
func corsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setCors(w)
}

// util
func getFrontendURL() string {
	if os.Getenv("APP_ENV") == "production" {
		return "http://localhost:3000" // change this to production domain
	}
	return "http://localhost:3000"
}

func setCors(w http.ResponseWriter) {
	frontendURL := getFrontendURL()
	w.Header().Set("Access-Control-Allow-Origin", frontendURL)
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

//Canary test to make sure Travis-CI is working.
func Canary(word string) string {
	return word
}

func main() {
	defer godorpdb.DB.Close()

	// add router and routes
	router := httprouter.New()
	router.GET("/", indexHandler)
	router.POST("/posts", createPostHandler)
	router.GET("/posts/:postId", showPostHandler)
	router.DELETE("/posts/:postId", deletePostHandler)
	router.PUT("/posts/:postId", updatePostHandler)
	router.GET("/posts", indexPostHandler)
	router.OPTIONS("/*any", corsHandler)

	// add database
	_, err := godorpdb.InitDb()
	if err != nil {
		log.Println("connection to DB failed, aborting...")
		log.Fatal(err)
	}

	log.Println("connected to DB")

	// print env
	env := os.Getenv("APP_ENV")
	if env == "production" {
		log.Println("Running api server in production mode")
	} else {
		log.Println("Running api server in dev mode")
	}

	http.ListenAndServe(":8080", router)
}
