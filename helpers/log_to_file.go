package helpers

import (
	"fmt"
	"log"
	"os"
)

func LogToFile(err error) {
	log.Println(err)
	file, err := os.OpenFile("error.logs", os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = file.Write([]byte(fmt.Sprint(err) + "\n"))
	if err != nil {
		log.Println(err)
		return
	}
}
