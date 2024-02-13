package config

import "testing"

func TestConfig(t *testing.T) {
	testCases := []struct{
		name string
		expected string
	}{
		{"database.mysql.host", "43.136.232.116"},
		{"database.mysql.port", "3307"},
		{"database.mysql.user", "xiaofei"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if actual := Get(tc.name); actual != tc.expected {
				t.Errorf("Config(%v) = %v; expected %v", tc.name, actual, tc.expected)
			}
		})
	}
}