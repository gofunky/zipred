package zipred

import (
	"io"
	"reflect"
	"testing"

	"bytes"
	"io/ioutil"
)

func Test_downloadFiles(t *testing.T) {
	type args struct {
		URL   string
		usage func(data io.ReadCloser) (useErr error)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Download empty",
			args: args{
				URL: "invalid",
				usage: func(data io.ReadCloser) error {
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "Download custom",
			args: args{
				URL: "https://github.com/gofunky/zipred/archive/master.zip",
				usage: func(data io.ReadCloser) error {
					_, err := ioutil.ReadAll(data)
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
		{
			name: "Check download fail",
			args: args{
				URL: "invalid",
				usage: func(data io.ReadCloser) error {
					return nil
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := downloadFiles(tt.args.URL, tt.args.usage); (err != nil) != tt.wantErr {
				t.Errorf("downloadFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_readNext(t *testing.T) {
	firstKey := "first"
	firstContent := "wins\n"
	secondKey := "second"
	secondContent := "loses\t"
	type fields struct {
		URL               string
		SelectedTemplates map[string]bool
		AdditionalRules   []string
		fetchedTemplates  map[string][]byte
	}
	type args struct {
		reader io.Reader
		key    string
	}
	tests := []struct {
		name       string
		fields     fields
		args       []args
		wantErr    bool
		wantTarget map[string][]byte
	}{
		{
			name:       "Single parse",
			args:       []args{{bytes.NewBufferString(firstContent), firstKey}},
			wantTarget: map[string][]byte{firstKey: []byte(firstContent)},
		},
		{
			name: "Dual parse",
			args: []args{
				{bytes.NewBufferString(firstContent), firstKey},
				{bytes.NewBufferString(firstContent), firstKey},
			},
			wantTarget: map[string][]byte{firstKey: []byte(firstContent)},
		},
		{
			name: "Dual parse different keys",
			args: []args{
				{bytes.NewBufferString(firstContent), firstKey},
				{bytes.NewBufferString(secondContent), secondKey},
			},
			wantTarget: map[string][]byte{firstKey: []byte(firstContent), secondKey: []byte(secondContent)},
		},
		{
			name:    "Without a valid key",
			args:    []args{{bytes.NewBufferString(firstContent), ""}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := make(map[string][]byte, len(tt.args))
			for _, arg := range tt.args {
				if err := readNext(arg.key, arg.reader, target); (err != nil) != tt.wantErr {
					t.Errorf("readNext() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if !reflect.DeepEqual(target, tt.wantTarget) && !tt.wantErr {
				t.Errorf("readNext() target = %v, wantTarget %v", target, tt.wantTarget)
			}
		})
	}
}
