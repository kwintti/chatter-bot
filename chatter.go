package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
    openai "github.com/sashabaranov/go-openai"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s and %s", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p := loadPage(title)
    fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

type Page struct {
    Title string
    Body []byte
}

func loadPage(title string) *Page {
    filename := title + ".txt"
    body, _ := os.ReadFile(filename)
    return &Page{Title: title, Body: body}
}

func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/view/", viewHandler)
    fmt.Println("Listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func getClient() *openai.Client {
    client := openai.NewClient(os.Getenv("TOKEN"))
    client.BaseUrl = os.Getenv("BASE_URL")

    return client
}


    
