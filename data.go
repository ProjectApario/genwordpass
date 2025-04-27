package main

import (
	"embed"
	"strings"

	"github.com/andreimerlescu/figtree/v2"
)

//go:embed VERSION
var versionBytes embed.FS

var currentVersion string

func Version() string {
	if len(currentVersion) == 0 {
		versionBytes, err := versionBytes.ReadFile("VERSION")
		if err != nil {
			return ""
		}
		currentVersion = strings.TrimSpace(string(versionBytes))
	}
	return currentVersion
}

//go:embed data/en.txt
var englishBytes []byte

//go:embed data/fr.txt
var frenchBytes []byte

//go:embed data/es.txt
var spanishBytes []byte

//go:embed data/de.txt
var germanBytes []byte

//go:embed data/ro.txt
var romanianBytes []byte

//go:embed data/ru.txt
var russianBytes []byte

var (
	acceptableWordSeparators = "!@#$%^&*()_+1234567890-=,.></?;:[]|"
	acceptableWords          []string
)

type Phoenix struct {
	Persona figtree.Plant
}

const (
	LanguageGerman   string = "de"
	LanguageEnglish  string = "en"
	LanguageSpanish  string = "es"
	LanguageFrench   string = "fr"
	LanguageRomanian string = "ro"
	LanguageRussian  string = "ru"

	PersonaOutputJSON string = "json"
	PersonaWordCount  string = "words"
	PersonaWordLength string = "length"
	PersonaVerbose    string = "verbose"
	PersonaVersion    string = "version"
	PersonaLanguages  string = "languages"
	PersonaSeparators string = "separators"
)
