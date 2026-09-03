package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestCreateGreetingDecodeIgnoresUnknownFieldsAndTrimsValues(t *testing.T) {
	var input createGreetingInput
	dec := json.NewDecoder(bytes.NewReader([]byte(`{"name":"  Ada  ","message":"  Hello  ","extra":"ignored"}`)))
	if err := dec.Decode(&input); err != nil {
		t.Fatalf("decode: %v", err)
	}

	name, message, details := validateGreetingInput(input)
	if len(details) != 0 {
		t.Fatalf("unexpected validation errors: %#v", details)
	}
	if name != "Ada" {
		t.Fatalf("name = %q, want %q", name, "Ada")
	}
	if message != "Hello" {
		t.Fatalf("message = %q, want %q", message, "Hello")
	}
}
