package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	e := run()
	if e != nil {
		log.Fatal(e)
	}
}

func run() error {
	ima, bootErr := NewPhoenix()
	if bootErr != nil {
		log.Fatal(bootErr)
	}
	password := ima.NewPassword()
	if *ima.Persona.Bool(PersonaOutputJSON) {
		b, e := json.MarshalIndent(map[string]string{"password": password}, "", "  ")
		if e != nil {
			return e
		}
		fmt.Println(string(b))
	} else {
		fmt.Println(password)
	}
	return nil
}
