package coloringgraph

import (
	"testing"
)

func TestColorNodeValueSerialization(t *testing.T) {
	tests := []struct {
		name     string
		value    ColorNodeValue
		expected []byte
	}{
		{
			name:     "empty color",
			value:    "",
			expected: []byte{},
		},
		{
			name:     "simple color",
			value:    "red",
			expected: []byte("red"),
		},
		{
			name:     "color with spaces",
			value:    "light blue",
			expected: []byte("light blue"),
		},
		{
			name:     "color with special characters",
			value:    "red#FF0000",
			expected: []byte("red#FF0000"),
		},
		{
			name:     "color with numbers",
			value:    "color123",
			expected: []byte("color123"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serialized := tt.value.Serialize()
			if len(serialized) != len(tt.expected) {
				t.Errorf("serialized length = %v, want %v", len(serialized), len(tt.expected))
			}
			for i := range serialized {
				if serialized[i] != tt.expected[i] {
					t.Errorf("serialized[%d] = %v, want %v", i, serialized[i], tt.expected[i])
				}
			}
		})
	}
}

func TestColorNodeValueDeserialization(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expected    ColorNodeValue
		expectError bool
	}{
		{
			name:        "empty data",
			data:        []byte{},
			expected:    "",
			expectError: false,
		},
		{
			name:        "simple color",
			data:        []byte("red"),
			expected:    "red",
			expectError: false,
		},
		{
			name:        "color with spaces",
			data:        []byte("light blue"),
			expected:    "light blue",
			expectError: false,
		},
		{
			name:        "color with special characters",
			data:        []byte("red#FF0000"),
			expected:    "red#FF0000",
			expectError: false,
		},
		{
			name:        "color with numbers",
			data:        []byte("color123"),
			expected:    "color123",
			expectError: false,
		},
		{
			name:        "color with unicode characters",
			data:        []byte("couleur"),
			expected:    "couleur",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := DeserializeColorNodeValue(tt.data)
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if value != tt.expected {
				t.Errorf("deserialized value = %v, want %v", value, tt.expected)
			}
		})
	}
}

func TestColorNodeValueRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		value ColorNodeValue
	}{
		{
			name:  "empty color",
			value: "",
		},
		{
			name:  "simple color",
			value: "red",
		},
		{
			name:  "color with spaces",
			value: "light blue",
		},
		{
			name:  "color with special characters",
			value: "red#FF0000",
		},
		{
			name:  "color with numbers",
			value: "color123",
		},
		{
			name:  "color with unicode characters",
			value: "couleur",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize
			serialized := tt.value.Serialize()

			// Deserialize
			deserialized, err := DeserializeColorNodeValue(serialized)
			if err != nil {
				t.Errorf("deserialization failed: %v", err)
			}

			// Compare
			if deserialized != tt.value {
				t.Errorf("round trip value = %v, want %v", deserialized, tt.value)
			}
		})
	}
}
