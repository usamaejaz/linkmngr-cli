package client

import (
	"reflect"
	"testing"
)

func TestExtractDomains(t *testing.T) {
	tests := []struct {
		name   string
		in     any
		want   []string
		wantOK bool
	}{
		{
			name:   "string list",
			in:     []any{"a.example", "b.example"},
			want:   []string{"a.example", "b.example"},
			wantOK: true,
		},
		{
			name: "object list",
			in: []any{
				map[string]any{"domain": "a.example"},
				map[string]any{"domain": "b.example"},
			},
			want:   []string{"a.example", "b.example"},
			wantOK: true,
		},
		{
			name: "wrapped domains field",
			in: map[string]any{
				"domains": []any{"a.example", "b.example"},
			},
			want:   []string{"a.example", "b.example"},
			wantOK: true,
		},
		{
			name: "wrapped items field",
			in: map[string]any{
				"items": []any{
					map[string]any{"domain": "a.example"},
					map[string]any{"domain": "b.example"},
				},
			},
			want:   []string{"a.example", "b.example"},
			wantOK: true,
		},
		{
			name:   "invalid shape",
			in:     map[string]any{"foo": "bar"},
			wantOK: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, ok := extractDomains(tc.in)
			if ok != tc.wantOK {
				t.Fatalf("ok mismatch: got %v want %v", ok, tc.wantOK)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("domains mismatch: got %#v want %#v", got, tc.want)
			}
		})
	}
}
