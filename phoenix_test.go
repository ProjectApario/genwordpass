package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkPhoenix(b *testing.B) {
	for n := 0; n < b.N; n++ {
		i, e := NewPhoenix()
		assert.NoError(b, e)
		assert.NotNil(b, i)
		assert.NotEmpty(b, i.NewPassword())
	}
}

func TestPhoenix(t *testing.T) {
	langs := []string{"en", "de", "fr", "ro", "es", "ru"}
	for _, lang := range langs {
		os.Args = []string{os.Args[0], "-languages", lang, "-words", "7", "-separators", "3"}
		ima, bootErr := NewPhoenix()
		assert.NoError(t, bootErr)
		assert.NotNil(t, ima)
		t.Run("Rise", func(t *testing.T) {
			assert.NoError(t, ima.Rise())
			assert.Equal(t, "Bool", string(ima.Persona.MutagenesisOfFig(PersonaVerbose)))
			assert.NotEmpty(t, *ima.Persona.List(PersonaLanguages))
		})
		t.Run("NewPassword", func(t *testing.T) {
			password := ima.NewPassword()
			assert.NotEmpty(t, password)
		})
		t.Run("LoadWords", func(t *testing.T) {
			err := ima.loadWords()
			assert.NoError(t, err)
		})
		t.Run("RandomInteger", func(t *testing.T) {
			randomInteger := ima.randomInt(1000)
			assert.NotEmpty(t, randomInteger)
			assert.LessOrEqual(t, randomInteger, 1000)
			assert.GreaterOrEqual(t, randomInteger, 0)
		})
		t.Run("GenerateWordPassword", func(t *testing.T) {
			for i := 1; i <= 1000; i++ {
				wordPassword, err := ima.generateWordPassword()
				assert.NoError(t, err)
				assert.NotEmpty(t, wordPassword)
			}
		})
		t.Run("ShouldCapitalize", func(t *testing.T) {
			trues, falses := 0, 0
			for i := 1; i <= 10000; i++ {
				b, e := ima.ShouldCapitalize()
				assert.NoError(t, e)
				if b {
					trues++
				} else {
					falses++
				}
			}
			if trues == 0 {
				t.Error("should capitalize at least one word")
			}
		})
	}
}
