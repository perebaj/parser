package gpt

import "testing"

func TestCreateUserPrompt(t *testing.T) {
	usrPromp, err := createUserPrompt("hello world")
	if err != nil {
		t.Error(err)
	}

	gotPrompt := `
	TEXT
		hello world
	END TEXT`

	if usrPromp != gotPrompt {
		t.Errorf("got %s, want %s", usrPromp, gotPrompt)
	}
}
