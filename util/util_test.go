package util

import (
    "testing"
)

func TestIsEmpty(t *testing.T) {
    testCases := []struct {
        name     string
        value    interface{}
        expected bool
    }{
        {"Empty int", 0, true},
        {"Non-empty int", 100, false},
        {"Empty string", "", true},
        {"Non-empty string", "hello", false},
        {"Empty slice", []int{}, true},
        {"Non-empty slice", []int{1}, false},
        {"Empty map", map[string]int{}, true},
        {"Non-empty map", map[string]int{"a": 1}, false},
        {"False bool", false, true},
        {"True bool", true, false},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            if actual := IsEmpty(tc.value); actual != tc.expected {
                t.Errorf("IsEmpty(%v) = %v; expected %v", tc.value, actual, tc.expected)
            }
        })
    }
}
