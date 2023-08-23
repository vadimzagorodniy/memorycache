package main

import (
	"testing"
	"time"
)

func TestProfile_Set(t *testing.T) {
	p := New()
	expectedVal := "test"
	p.Insert(expectedVal)

	if p.Orders[0].Value != expectedVal {
		t.Errorf("Expected %q to be a value of order", expectedVal)
	}
	newExpectedVal := "changed"
	p.Set(p.Orders[0].UUID, newExpectedVal)

	if p.Orders[0].Value != "changed" {
		t.Errorf("Expected %q to be a new value of order", newExpectedVal)
	}
}

func TestProfile_Insert(t *testing.T) {
	p := New()
	expectedVal := "test"
	p.Insert(expectedVal)

	if p.Orders[0].Value != "test" {
		t.Errorf("Expected %q to be a value of order", expectedVal)
	}
}

func TestProfile_Delete(t *testing.T) {
	p := New()
	expectedVal := "test"
	forDelete := "delete"
	id := p.Insert(expectedVal)
	forDeleteId := p.Insert(forDelete)
	if p.Orders[0].UUID != *id {
		t.Errorf("Expected %q to be equal to order UUID %q", *id, p.Orders[0].UUID)
	}
	if p.Orders[1].UUID != *forDeleteId {
		t.Errorf("Expected %q to be equal to order UUID %q", *forDeleteId, p.Orders[0].UUID)
	}

	answer := *p.Delete(*id)
	if answer != true {
		t.Errorf("Expected answer to be true")
	}
}

func TestProfile_IsExpiredTTL(t *testing.T) {
	p := New()
	id := *p.Insert("test")

	time.Sleep(TTL + TTL)
	answer := p.Insert("test2")
	if answer != nil {
		t.Errorf("Expected return of the function to be nil")
	}
	answer2 := p.Set(id, "test3")
	if answer2 != nil {
		t.Errorf("Expected return of the function to be nil")
	}
	answer3 := p.Delete(id)
	if answer3 != nil {
		t.Errorf("Expected return of the function to be nil")
	}
}
