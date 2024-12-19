package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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
    err := godotenv.Load()
    if err != nil {
        log.Println("Couldn't load .env file", err)
    }
    // http.HandleFunc("/", handler)
    // http.HandleFunc("/view/", viewHandler)
    // fmt.Println("Listening on :8080")
    // log.Fatal(http.ListenAndServe(":8080", nil))
    fmt.Println("Miten voin auttaa?:")
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    err = scanner.Err()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf(chatter(scanner.Text()))

}

func getClient() *openai.Client {
    config := openai.DefaultConfig(os.Getenv("TOKEN"))
    config.BaseURL = os.Getenv("BASE_URL") 

    client := openai.NewClientWithConfig(config)

    return client
}

func chatter(msg string) string {
    client := getClient()
    resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
            Model: "meta-llama/Meta-Llama-3.1-8B-Instruct-fast",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: msg,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	return resp.Choices[0].Message.Content

}

    
