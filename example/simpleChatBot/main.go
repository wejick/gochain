package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/wejick/gochain/chain/conversation"
	"github.com/wejick/gochain/model"
	_openai "github.com/wejick/gochain/model/openAI"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Type .quit to exit")

	var authToken = os.Getenv("OPENAI_API_KEY")
	chatModel := _openai.NewOpenAIChatModel(authToken, "", _openai.GPT3Dot5Turbo0301)
	memory := []model.ChatMessage{}
	convoChain := conversation.NewConversationChain(chatModel, memory, "You're helpful chatbot that answer human question very concisely")
	convoChain.AppendToMemory(model.ChatMessage{Role: model.ChatMessageRoleAssistant, Content: "Hi, My name is GioAI"})

	for {
		fmt.Print("User : ")
		chat, _ := reader.ReadString('\n')

		// Remove newline character from the command string
		chat = chat[:len(chat)-1]

		if chat == ".quit" {
			break
		}

		output, err := convoChain.SimpleRun(context.Background(), chat)
		if err != nil {
			fmt.Println("error " + err.Error())
			break
		}

		fmt.Println("AI :", output)
	}

	fmt.Println("Program exited.")
}