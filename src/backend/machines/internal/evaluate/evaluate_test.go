package evaluate

import (
	"testing"
)

func TestSolve(t *testing.T) {
	ch := make(chan string)

	go Solve("2 + 3 * 4 - 6 / 2", ch)

	result := <-ch

	if result != "11" {
		t.Errorf("Expected result 11, but got %s", result)
	}
	go Solve("2 + 2 * 2", ch)
	result = <-ch
	if result != "6" {
		t.Errorf("Expected result 6, but got %s", result)
	}
}

func TestSubSolve(t *testing.T) {
	// Test multiplication
	result, err := subSolve("2", "*", "3")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if result != "6" {
		t.Errorf("Expected result 6 for 2 * 3, but got %s", result)
	}

	// Test division
	result, err = subSolve("10", "/", "5")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if result != "2" {
		t.Errorf("Expected result 2 for 10 / 5, but got %s", result)
	}

	// Test addition
	result, err = subSolve("5", "+", "3")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if result != "8" {
		t.Errorf("Expected result 8 for 5 + 3, but got %s", result)
	}

	// Test subtraction
	result, err = subSolve("10", "-", "3")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if result != "7" {
		t.Errorf("Expected result 7 for 10 - 3, but got %s", result)
	}

	// Test invalid operator
	_, err = subSolve("10", "%", "3")
	if err == nil {
		t.Error("Expected an error for an invalid operator, but got none")
	}
}
