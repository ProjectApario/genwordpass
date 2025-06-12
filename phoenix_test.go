package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// BenchmarkPhoenix determines the ns/op per NewPassword() invocation
func BenchmarkPhoenix(b *testing.B) {
	for n := 0; n < b.N; n++ {
		i, e := NewPhoenix()
		assert.NoError(b, e)
		assert.NotNil(b, i)
		assert.NotEmpty(b, i.NewPassword())
	}
}

// FuzzNewPhoenix will use Go's Fuzz testing functionality to challenge the CLI arguments that can be passed into the app
func FuzzNewPhoenix(f *testing.F) {
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	testCases := []struct {
		languages   string
		words       string
		separators  string
		expectError bool
		description string
	}{
		{"", "3", "1", false, "empty languages should use defaults"},
		{"en", "3", "1", false, "valid case"},
		{"en,es,fr", "5", "2", false, "multiple languages"},
		{"invalid", "3", "1", false, "invalid language should not fail parsing"},
		{"en", "-1", "1", true, "negative words should fail"},
		{"en", "abc", "1", true, "non-numeric words should fail"},
		{"en", "1000", "1", true, "words too large should fail"},
		{"", "", "", true, "all empty should fail"},
		{"en", "3", "abc", true, "non-numeric separators should fail"},
		{"en,", "3", "1", false, "trailing comma should be handled"},
		{",en", "3", "1", false, "leading comma should be handled"},
		{"en,,es", "3", "1", false, "double comma should be handled"},
		{"EN", "3", "1", false, "uppercase language should work"},
		{"en-US", "3", "1", false, "language with region should work"},
		{"123", "3", "1", false, "numeric language should not fail parsing"},
		{"en", "0", "1", true, "zero words should fail"},
		{"en", "3", "-999", true, "very negative separators should fail"},
		{"en", "999999", "1", true, "very large words should fail"},
		{"a", "3", "1", false, "single char language should not fail parsing"},
		{"en", "1.5", "1", true, "decimal words should fail"},
		{"en", "3", "1.5", true, "decimal separators should fail"},
		{"en", "1", "1", true, "words below minimum (3) should fail"},
		{"en", "2", "1", true, "words below minimum (3) should fail"},
		{"en", "70", "1", true, "words above maximum (69) should fail"},
		{"en", "3", "0", true, "separators below minimum (1) should fail"},
		{"en", "3", "8", true, "separators above maximum (7) should fail"},
		{"en", "69", "7", false, "maximum valid values should work"},
	}

	for _, tc := range testCases {
		f.Add(tc.languages, tc.words, tc.separators, tc.expectError)
	}

	f.Fuzz(func(t *testing.T, languages, words, separators string, expectError bool) {
		os.Args = []string{
			originalArgs[0],
			"-languages", languages,
			"-words", words,
			"-separators", separators,
		}

		ima, bootErr := NewPhoenix()

		if expectError {
			assert.Error(t, bootErr, "Expected error for args: languages=%q, words=%q, separators=%q",
				languages, words, separators)
		} else {
			assert.NoError(t, bootErr, "Expected no error for args: languages=%q, words=%q, separators=%q",
				languages, words, separators)
			if bootErr == nil {
				assert.NotNil(t, ima, "Expected non-nil result when no error occurred")
			}
		}
	})
}

// TestPhoenix consumes the runtime of the application to verify full coverage of package
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
