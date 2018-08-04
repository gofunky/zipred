package gitignore

import (
	"testing"
)

func TestGet(t *testing.T) {
	type args struct {
		patterns []string
	}
	tests := []struct {
		name      string
		args      args
		wantFiles []string
		wantErr   bool
	}{
		{
			name:      "Check some patterns",
			args:      args{[]string{"go", "java"}},
			wantFiles: []string{"go", "java"},
		},
		{
			name:    "Invalid patterns",
			args:    args{[]string{"invalid"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFiles, err := Get(tt.args.patterns)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				for _, pattern := range tt.wantFiles {
					if val, ok := gotFiles[pattern]; !ok || len(val) == 0 {
						t.Errorf("Get() is missing key %s or it is empty", pattern)
					}
				}
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	tests := []struct {
		name      string
		wantFiles []string
		wantErr   bool
	}{
		{
			name:      "Check some patterns",
			wantFiles: []string{"go", "java"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFiles, err := GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				for _, pattern := range tt.wantFiles {
					if val, ok := gotFiles[pattern]; !ok || len(val) == 0 {
						t.Errorf("GetAll() is missing key %s or it is empty", pattern)
					}
				}
			}
		})
	}
}
