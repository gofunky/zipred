package gitignore

import (
	"testing"
)

func TestList(t *testing.T) {
	tests := []struct {
		name         string
		wantPatterns []string
		wantErr      bool
	}{
		{
			name:         "Check some patterns",
			wantPatterns: []string{"go", "java"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPatterns, err := List()
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, want := range tt.wantPatterns {
				if !contains(gotPatterns, want) {
					t.Errorf("List() = %v, want %v", gotPatterns, tt.wantPatterns)
				}
			}
		})
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
