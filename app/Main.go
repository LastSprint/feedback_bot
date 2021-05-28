package main

import (
	"github.com/LastSprint/feedback_bot/DB"
	"github.com/LastSprint/feedback_bot/Rest"
	"log"
	"net/http"
	"os"
)

func main() {

	path := os.Getenv("FEEDBACK_BOT_DB_FILE_PATH")

	db := DB.FileDB{FilePath: path}
	controller := Rest.SlackController{DB: &db}
	controller.Init()
	log.Fatal(http.ListenAndServe(":6654", nil))
}