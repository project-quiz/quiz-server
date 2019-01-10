package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	Question "github.com/project-quiz/quiz-go-model/Question"

	//Use _ because it is needed for mysql driver to be imported.
	_ "github.com/go-sql-driver/mysql"
)

type DatabaseService struct {
	database *sql.DB
}

//Connect connect with the questions database.
func Connect() (*DatabaseService, error) {
	ip, ipSuccess := os.LookupEnv("DATABASE_IP")
	port, portSuccess := os.LookupEnv("DATABASE_PORT")
	username, usernameSuccess := os.LookupEnv("DATABASE_USERNAME")
	password, passwordSuccess := os.LookupEnv("DATABASE_PASSWORD")

	success := ipSuccess && portSuccess && usernameSuccess && passwordSuccess

	if !success {
		return nil, errors.New("not enough database values available")
	}

	db, err := sql.Open("mysql", username+":"+password+"@tcp("+ip+":"+port+")/quiz")

	if err != nil {
		fmt.Println("err" + err.Error())
		return nil, err
	}

	service := DatabaseService{db}

	return &service, err
}

func (service *DatabaseService) GetQuestion(guid string) {
	Result, errDB := service.database.Query("SELECT guid,text FROM questions WHERE guid='" + guid + "'")

	if errDB != nil {
		fmt.Println("Error" + errDB.Error())
	}

	for Result.Next() {
		var (
			guid string
			text string
		)
		if err := Result.Scan(&guid, &text); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("guid %s question is %s\n", guid, text)
	}
}

//GetQuestions Get all questions from the database.
func (service *DatabaseService) GetQuestions() (questions []Question.Question, err error) {
	Result, err := service.database.Query("SELECT guid,text FROM questions")

	if err != nil {
		fmt.Println("Question error: " + err.Error())
	}

	q := []Question.Question{}

	for Result.Next() {
		var (
			guid string
			text string
		)
		if err := Result.Scan(&guid, &text); err != nil {
			log.Fatal(err)
		}

		a, answerErr := service.GetAnswers(guid)

		if answerErr != nil {
			fmt.Println("Answer err for question " + guid + ": " + answerErr.Error())
		}

		question := Question.Question{Guid: guid, Question: text, Answers: a}
		q = append(q, question)
	}

	return q, err
}

func (service *DatabaseService) GetAnswers(questionGuid string) (answers []*Question.Answer, err error) {
	Result, err := service.database.Query("SELECT guid,answer FROM answers WHERE question='" + questionGuid + "'")

	a := []*Question.Answer{}

	for Result.Next() {
		var (
			guid   string
			answer string
		)
		if err := Result.Scan(&guid, &answer); err != nil {
			log.Fatal(err)
		}
		obj := Question.Answer{Guid: guid, Text: answer}

		a = append(a, &obj)
	}
	return a, err
}
