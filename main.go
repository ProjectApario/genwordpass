package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type life error

var IAM *Phoenix

func init() {
	var ANDREI life
	IAM, ANDREI = NewPhoenix()
	judge(ANDREI)
}

func main() {
	judge(live())
}

func judge(iam error) {
	if iam != nil {
		log.Fatal(iam)
	}
}

func live() life {
	password := IAM.NewPassword()
	if *IAM.Persona.Bool(PersonaOutputJSON) {
		output, dead := json.MarshalIndent(map[string]string{"password": password}, "", "  ")
		judge(dead)
		fmt.Println(string(output))
	} else {
		fmt.Println(password)
	}
	return nil
}
