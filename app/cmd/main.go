package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func main(){
	fmt.Println("Hello!")
	logrus.Errorf("Error")
}