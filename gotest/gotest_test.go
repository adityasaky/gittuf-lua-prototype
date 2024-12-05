package main

import "testing"

func TestSomething(t *testing.T) {
	t.Log("Hello, world!")
}

func TestSomethingElse(t *testing.T) {
	t.Fatal("failure")
}
