package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
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

type Conversation struct {
    Messages []openai.ChatCompletionMessage
}

func newConversation() *Conversation {
    return &Conversation{
        Messages: []openai.ChatCompletionMessage{
            {
                Role: openai.ChatMessageRoleSystem,
                Content: "Olet ride-hailing yhtiön neuvokas asiakaspalvelija. Vastaat rennosti.",
            },
        },
    }
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Couldn't load .env file", err)
    }
    component := hello("JOhn")
    http.Handle("/", templ.Handler(component))
    // http.HandleFunc("/", handler)
    // http.HandleFunc("/view/", viewHandler)
    fmt.Println("Listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
    fmt.Println("Miten voin auttaa?:")
    for {
        scanner := bufio.NewScanner(os.Stdin)
        scanner.Scan()
        err = scanner.Err()
        if err != nil {
            log.Fatal(err)
        }
        client := getClient() 
        conv := newConversation()
        conv.addMessage(openai.ChatMessageRoleUser, scanner.Text())

        resp, err := conv.getCompletion(client)
        if err != nil {
            log.Fatalf("Couldn't get a response from llm: %v", err)
        }
        fmt.Println(resp)
    }
}

func getClient() *openai.Client {
    config := openai.DefaultConfig(os.Getenv("TOKEN"))
    config.BaseURL = os.Getenv("BASE_URL") 

    client := openai.NewClientWithConfig(config)

    return client
}

func (c *Conversation) getCompletion(client *openai.Client) (string,error) {
    client = getClient()
    resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
            Model: "meta-llama/Llama-3.3-70B-Instruct",
			Messages: c.Messages,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil

}

func (c *Conversation) addMessage(role string, content string){
    c.Messages = append(c.Messages, openai.ChatCompletionMessage{
        Role: role,
        Content: content,
    })
}
