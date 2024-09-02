package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type PostChatGroqType struct {
	Promt string `json:"promt"`
}

func PostChatGroq(c *gin.Context) {

	err := handleCheckHost(c)
	if(err != nil){
		c.JSON(http.StatusOK,gin.H{
				"result": err.Error(),
				"message": "Error",
			})
			return
	}

	var reqBody PostChatGroqType
	
	if err := c.BindJSON(&reqBody); err != nil {
		c.IndentedJSON(http.StatusInternalServerError,"error")
	}

	apiKey := os.Getenv("GROQ_API_KEY")
	base_url := os.Getenv("GROQ_BASE_URL")
	model := os.Getenv("GROQ_MODEL")

	llm, err := openai.New(
        openai.WithModel(model),
        openai.WithBaseURL(base_url),
        openai.WithToken(apiKey),
    )
    if err != nil {
        log.Fatal(err)
    }

	// var completion  string

    ctx := context.Background()
    completion, err := llms.GenerateFromSinglePrompt(ctx,
        llm,
        reqBody.Promt,
        llms.WithTemperature(0.8),
        llms.WithMaxTokens(4096),
        // llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
        //     fmt.Print(string(chunk))
        //     return nil
        // }),
    )
    fmt.Println()
    if err != nil {
        log.Fatal(err)
    }

	c.JSON(http.StatusOK,gin.H{
		"result": completion,
		"message": "Success",
	})

}

func main() {
	err := godotenv.Load()
	if err != nil {
        log.Fatalf("Error loading .env file")
    }
	router := gin.Default()
	router.POST("/api/chat-groq",PostChatGroq)
	router.Run()
}

func handleCheckHost(c *gin.Context)(error){
	host := c.Request.Host

	fmt.Printf("host in func : %v",host)

	if (host != "localhost:8001"){
		return errors.New("Not a valid host")
	}

	return nil
	
	
	
}
