package graph

import (
	"testing"
)

func TestIntNodeValueSerialization(t *testing.T) {
	tests := []struct {
		name     string
		value    IntNodeValue
		expected []byte
	}{
		{
			name:     "zero value",
			value:    0,
			expected: []byte{0, 0},
		},
		{
			name:     "small value",
			value:    42,
			expected: []byte{0, 42},
		},
		{
			name:     "max uint16 value",
			value:    0xFFFF,
			expected: []byte{255, 255},
		},
		{
			name:     "medium value",
			value:    1234,
			expected: []byte{4, 210},
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

func TestIntNodeValueDeserialization(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expected    IntNodeValue
		expectError bool
	}{
		{
			name:        "zero value",
			data:        []byte{0, 0},
			expected:    0,
			expectError: false,
		},
		{
			name:        "small value",
			data:        []byte{0, 42},
			expected:    42,
			expectError: false,
		},
		{
			name:        "max uint16 value",
			data:        []byte{255, 255},
			expected:    0xFFFF,
			expectError: false,
		},
		{
			name:        "medium value",
			data:        []byte{4, 210},
			expected:    1234,
			expectError: false,
		},
		{
			name:        "empty data",
			data:        []byte{},
			expected:    0,
			expectError: true,
		},
		{
			name:        "incomplete data",
			data:        []byte{0},
			expected:    0,
			expectError: true,
		},
		{
			name:        "extra data",
			data:        []byte{0, 42, 0},
			expected:    42,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := DeserializeIntNodeValue(tt.data)
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				if err != ErrDataTooShort {
					t.Errorf("expected ErrDataTooShort, got %v", err)
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

func TestIntNodeValueRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		value IntNodeValue
	}{
		{
			name:  "zero value",
			value: 0,
		},
		{
			name:  "small value",
			value: 42,
		},
		{
			name:  "max uint16 value",
			value: 0xFFFF,
		},
		{
			name:  "medium value",
			value: 1234,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize
			serialized := tt.value.Serialize()

			// Deserialize
			deserialized, err := DeserializeIntNodeValue(serialized)
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
