package flags

import (
	"testing"
)

func TestPassDoubleDash(t *testing.T) {
	var opts = struct {
		Value bool `short:"v"`
	}{}

	p := NewParser(&opts, PassDoubleDash)
	ret, _, _, err := p.ParseCommandLine([]string{"-v", "--", "-v", "-g"})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
		return
	}

	if !opts.Value {
		t.Errorf("Expected Value to be true")
	}

	assertStringArray(t, ret, []string{"-v", "-g"})
}

func TestPassAfterNonOption(t *testing.T) {
	var opts = struct {
		Value bool `short:"v"`
	}{}

	p := NewParser(&opts, PassAfterNonOption)
	ret, _, _, err := p.ParseCommandLine([]string{"-v", "arg", "-v", "-g"})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
		return
	}

	if !opts.Value {
		t.Errorf("Expected Value to be true")
	}

	assertStringArray(t, ret, []string{"arg", "-v", "-g"})
}
