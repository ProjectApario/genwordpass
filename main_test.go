package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"syscall"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	if code != 0 {
		os.Exit(code)
	}
}

func TestRun(t *testing.T) {
	t.Run("Version", func(t *testing.T) {
		os.Args = []string{os.Args[0], "-version"}
		assert.Panics(t, func() { // this is expected with NoError
			assert.NoError(t, run()) // this calls os.Exit(0) which causes panic in tests
		})
	})
	t.Run("Help", func(t *testing.T) {
		os.Args = []string{os.Args[0], "-help"}
		err := exec.Command(os.Args[0], "-help").Run()
		if exitError, ok := err.(*exec.ExitError); !ok {
			assert.NoError(t, err)
		} else {
			t.Log(string(exitError.Stderr))
		}
	})
	os.Stdout = os.NewFile(uintptr(syscall.Stdin), os.DevNull) // protect from printing passwords to stdout
	t.Run("JSON", func(t *testing.T) {
		os.Args = []string{os.Args[0], "-json"}
		assert.NoError(t, run())
	})
	t.Run("TEXT", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		assert.NoError(t, run())
	})
	langs := []string{"en", "de", "fr", "es", "ro", "ru"}
	for _, lang := range langs {
		t.Run(lang, func(t *testing.T) {
			os.Args = []string{os.Args[0], "-languages=" + lang}
			assert.NoError(t, run())
		})
	}
}
