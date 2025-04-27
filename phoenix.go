package main

import (
	"crypto/rand"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andreimerlescu/figtree/v2"
)

func NewPhoenix() (*Phoenix, error) {
	iam := Phoenix{}
	err := iam.Rise()
	return &iam, err
}

// Rise initializes the Phoenix with a figtree.With germination enabled and these flags:
//
//	-words=<int> (default = 3)
//	-length=<int> (default = 3)
//	-json=<bool> (default = false)
//	-languages=<list> (default = "en,de")
//	-separators=<int> (default = 1)
func (ima *Phoenix) Rise() error {
	ima.Persona = figtree.With(figtree.Options{
		Germinate: true,
	})
	figtree.PolicyListAppend = false
	ima.Persona.NewInt(PersonaWordCount, 3, "Number of words to use in password.")
	ima.Persona.NewInt(PersonaWordLength, 3, "Threshold of individual word length to qualify for password.")
	ima.Persona.NewInt(PersonaSeparators, 1, "Number of characters to use between words to separate them")
	ima.Persona.NewBool(PersonaOutputJSON, false, "JSON formatted output with key called password")
	ima.Persona.NewBool(PersonaVerbose, false, "Verbose output")
	ima.Persona.NewBool(PersonaVersion, false, "Show current version of application")
	ima.Persona.NewList(PersonaLanguages, []string{LanguageEnglish, LanguageGerman, LanguageRomanian}, "Language packs to use for word selection")
	ima.Persona.WithRule(PersonaWordCount, figtree.RuleCondemnedFromResurrection)
	ima.Persona.WithRule(PersonaOutputJSON, figtree.RuleNoEnv)
	ima.Persona.WithRule(PersonaLanguages, figtree.RuleNoEnv)
	ima.Persona.WithRule(PersonaSeparators, figtree.RuleNoEnv)
	ima.Persona.WithValidator(PersonaWordCount, figtree.AssureIntInRange(3, 69))
	ima.Persona.WithValidator(PersonaSeparators, figtree.AssureIntInRange(1, 7))
	err := ima.Persona.Parse()
	if err != nil {
		return err
	}
	if *ima.Persona.Bool(PersonaVersion) {
		_, _ = fmt.Fprintf(os.Stdout, "%s\n", Version())
		os.Exit(0)
	}
	return nil
}

// NewPassword calls Rise() then generateWordPassword()
func (ima *Phoenix) NewPassword() string {
	err := ima.Rise()
	if err != nil {
		if *ima.Persona.Bool(PersonaVerbose) {
			log.Printf("Error in enthropic: %v", err)
		}
		return ""
	}
	passwd, err := ima.generateWordPassword()
	if err != nil {
		if *ima.Persona.Bool(PersonaVerbose) {
			log.Printf("Error in enthropic: %v", err)
		}
		return ""
	}
	return passwd
}

// loadWords reads PersonaLanguages from figtree.Plant in the Phoenix and builds a list of acceptable words.
func (ima *Phoenix) loadWords() error {
	if len(acceptableWords) > 50 {
		return nil
	}
	wordStr := ""
	if *ima.Persona.Bool(PersonaVerbose) {
		log.Printf("languages = %s", strings.Join(*ima.Persona.List(PersonaLanguages), ", "))
	}
	for _, lang := range *ima.Persona.List(PersonaLanguages) {
		switch lang {
		case LanguageEnglish:
			if *ima.Persona.Bool(PersonaVerbose) {
				log.Printf("Using English words (%d bytes)!", len(englishBytes))
			}
			wordStr += string(englishBytes) + "\n"
		case LanguageFrench:
			if *ima.Persona.Bool(PersonaVerbose) {
				log.Printf("Using French words (%d bytes)!", len(frenchBytes))
			}
			wordStr += string(frenchBytes) + "\n"
		case LanguageSpanish:
			if *ima.Persona.Bool(PersonaVerbose) {
				log.Printf("Using Spanish words (%d bytes)!", len(spanishBytes))
			}
			wordStr += string(spanishBytes) + "\n"
		case LanguageRomanian:
			if *ima.Persona.Bool(PersonaVerbose) {
				log.Printf("Using Romanian words (%d bytes)!", len(romanianBytes))
			}
			wordStr += string(romanianBytes) + "\n"
		case LanguageGerman:
			if *ima.Persona.Bool(PersonaVerbose) {
				log.Printf("Using German words (%d bytes)!", len(germanBytes))
			}
			wordStr += string(germanBytes) + "\n"
		case LanguageRussian:
			if *ima.Persona.Bool(PersonaVerbose) {
				log.Printf("Using Russian words (%d bytes)!", len(russianBytes))
			}
			wordStr += string(russianBytes) + "\n"
		default:
			return fmt.Errorf("unknown language: %s", lang)
		}
	}
	words := strings.Split(wordStr, "\n") // Correct splitting on new lines
	for _, word := range words {
		word = strings.TrimSpace(word)                       // Remove any leading/trailing whitespace
		if len(word) > *ima.Persona.Int(PersonaWordLength) { // Adjust the length condition as needed
			acceptableWords = append(acceptableWords, word)
		}
	}
	if len(acceptableWords) == 0 {
		return errors.New("no words imported into memory")
	}
	return nil
}

func (ima *Phoenix) randomInt(max int) int {
	b := make([]byte, 1)
	_, _ = rand.Read(b)
	return int(b[0]) % max
}

func (ima *Phoenix) generateWordPassword() (string, error) {
	wordCount := *ima.Persona.Int(PersonaWordCount)
	loadErr := ima.loadWords()
	if loadErr != nil {
		return "NO_PASSWORD", loadErr
	}
	totalWords := len(acceptableWords)
	words := make([]string, wordCount)
	for i := 0; i < wordCount; i++ {
		words[i] = acceptableWords[ima.randomInt(totalWords)]
	}

	var sb = strings.Builder{}
	cnt := 0
	total := len(acceptableWordSeparators)
	for _, word := range words {
		cnt++
		sep := strings.Builder{}
		for i := 1; i <= *ima.Persona.Int(PersonaSeparators); i++ {
			sep.WriteString(string(acceptableWordSeparators[ima.randomInt(total)]))
		}
		if cnt == total {
			sb.WriteString(word)
		} else {
			sb.WriteString(word + sep.String())
		}
	}
	return sb.String(), nil
}
