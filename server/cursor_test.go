package server

import (
	"testing"

	"github.com/sivchari/go-mcp/apis"
	"github.com/sivchari/go-mcp/server/internal/testhelper"
)

func Test_nextCursor(t *testing.T) {
	type args struct {
		cursor *string
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		{
			name: "cursor is nil",
			args: args{
				cursor: nil,
			},
			want: testhelper.To("2"),
		},
		{
			name: "cursor is empty",
			args: args{
				cursor: testhelper.To(""),
			},
			want: testhelper.To("2"),
		},
		{
			name: "cursor is not empty",
			args: args{
				cursor: testhelper.To("1"),
			},
			want: testhelper.To("2"),
		},
		{
			name: "cursor is 2",
			args: args{
				cursor: testhelper.To("2"),
			},
			want: testhelper.To("3"),
		},
		{
			name: "cursor is last",
			args: args{
				cursor: testhelper.To("3"),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &Server{
				prompts: prompts{
					lists: map[string][]apis.Prompt{
						"1": {},
						"2": {},
						"3": {},
					},
				},
			}
			got := s.nextCursor(tt.args.cursor)
			if tt.want == nil && got != nil {
				t.Fatalf("nextCursor() = %v, want nil", *got)
			}
			if tt.want != nil && got == nil {
				t.Fatalf("nextCursor() = nil, want %v", *tt.want)
			}
			if tt.want == nil && got == nil {
				return
			}
			if *got != *tt.want {
				t.Fatalf("nextCursor() = %v, want %v", *got, *tt.want)
			}
		})
	}
}
