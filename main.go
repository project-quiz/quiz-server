package main

import (
	"fmt"

	"guido.arkesteijn/quiz-server/database"
	"guido.arkesteijn/quiz-server/server"
)

var messageID int32 = 1
var stopServer bool = false

func main() {
	srv, err := database.Connect("192.168.2.18", "4600")

	if err != nil {
		fmt.Printf("error", err.Error())
	}

	questions, questionErr := srv.GetQuestions()

	if questionErr != nil {
		fmt.Println("Error : " + questionErr.Error())
	}

	for _, element := range questions {
		fmt.Println("Question:", element.Question)
		fmt.Println("Answers:")
		for _, answer := range element.Answers {
			fmt.Println(answer.Text)
		}
	}

	go server.StartServer(4500)

	for {
		if stopServer {
			break
		}
	}
}
