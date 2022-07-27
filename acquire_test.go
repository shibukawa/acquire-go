package acquire

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAcquireFromUnder(t *testing.T) {
	current, err := filepath.Abs(".")
	if err != nil {
		panic(err.Error())
	}

	type args struct {
		t        TargetType
		folder   string
		under    string
		patterns []string
	}
	tests := []struct {
		name        string
		args        args
		wantMatches []string
		wantErr     error
	}{
		{
			name: "search in same folder, abs path",
			args: args{
				t:        All,
				folder:   filepath.Join(current, "testdata"),
				patterns: []string{"target1.txt"},
			},
			wantMatches: []string{
				filepath.Join(current, "testdata", "target1.txt"),
			},
			wantErr: nil,
		},
		{
			name: "search in same folder, relative path",
			args: args{
				t:        All,
				folder:   "testdata",
				patterns: []string{"target1.txt"},
			},
			wantMatches: []string{
				filepath.Join(current, "testdata", "target1.txt"),
			},
			wantErr: nil,
		},
		{
			name: "search in same folder, match only file (OK)",
			args: args{
				t:        File,
				folder:   "testdata",
				patterns: []string{"target1.txt"},
			},
			wantMatches: []string{
				filepath.Join(current, "testdata", "target1.txt"),
			},
			wantErr: nil,
		},
		{
			name: "search in same folder, match only file (NG)",
			args: args{
				t:        File,
				folder:   "testdata",
				patterns: []string{"target"},
			},
			wantMatches: nil,
			wantErr:     ErrNotFound,
		},
		{
			name: "search in same folder, match only dir (OK)",
			args: args{
				t:        Dir,
				folder:   "testdata",
				patterns: []string{"target"},
			},
			wantMatches: []string{
				filepath.Join(current, "testdata", "target"),
			},
			wantErr: nil,
		},
		{
			name: "search in same folder, match only dir (NG)",
			args: args{
				t:        Dir,
				folder:   "testdata",
				patterns: []string{"target1.txt"},
			},
			wantMatches: nil,
			wantErr:     ErrNotFound,
		},
		{
			name: "search in same folder, multiple patterns",
			args: args{
				t:        File,
				folder:   "testdata",
				patterns: []string{"target1.txt", "target2.txt"},
			},
			wantMatches: []string{
				filepath.Join(current, "testdata", "target1.txt"),
				filepath.Join(current, "testdata", "target2.txt"),
			},
			wantErr: nil,
		},
		{
			name: "search in same folder, glob patterns 1",
			args: args{
				t:        File,
				folder:   "testdata",
				patterns: []string{"*.txt"},
			},
			wantMatches: []string{
				filepath.Join(current, "testdata", "submatch.txt"),
				filepath.Join(current, "testdata", "target1.txt"),
				filepath.Join(current, "testdata", "target2.txt"),
			},
			wantErr: nil,
		},
		{
			name: "search in same folder, glob patterns 2",
			args: args{
				t:        File,
				folder:   "testdata",
				patterns: []string{"target*.txt"},
			},
			wantMatches: []string{
				filepath.Join(current, "testdata", "target1.txt"),
				filepath.Join(current, "testdata", "target2.txt"),
			},
			wantErr: nil,
		},
		{
			name: "search from deep folder, match in parent folder",
			args: args{
				t:        File,
				folder:   "testdata/sub1",
				patterns: []string{"target*.txt"},
			},
			wantMatches: []string{
				filepath.Join(current, "testdata", "target1.txt"),
				filepath.Join(current, "testdata", "target2.txt"),
			},
			wantErr: nil,
		},
		{
			name: "search from deep folder, match in ancestor folder",
			args: args{
				t:        File,
				folder:   "testdata/sub1/sub2/sub3",
				patterns: []string{"target*.txt"},
			},
			wantMatches: []string{
				filepath.Join(current, "testdata", "target1.txt"),
				filepath.Join(current, "testdata", "target2.txt"),
			},
			wantErr: nil,
		},
		{
			name: "search from deep folder, match in uncle folder",
			args: args{
				t:        File,
				folder:   "testdata/sub1/sub2",
				patterns: []string{"target/*.txt"},
			},
			wantMatches: []string{
				filepath.Join(current, "testdata", "target", "target1.txt"),
				filepath.Join(current, "testdata", "target", "target2.txt"),
			},
			wantErr: nil,
		},
		{
			name: "search from deep folder, match in the middle folder",
			args: args{
				t:        File,
				folder:   "testdata/sub1/sub2/sub3",
				patterns: []string{"*.txt"},
			},
			wantMatches: []string{
				filepath.Join(current, "testdata", "sub1", "sub2", "submatch.txt"),
			},
			wantErr: nil,
		},
		{
			name: "search from deep folder, match in ancestor folder under limit",
			args: args{
				t:        File,
				folder:   "testdata/sub1/sub2/sub3",
				under:    "testdata/sub1/sub2",
				patterns: []string{"*.txt"},
			},
			wantMatches: []string{
				filepath.Join(current, "testdata", "sub1", "sub2", "submatch.txt"),
			},
			wantErr: nil,
		},
		{
			name: "search from deep folder, not match in ancestor folder over limit",
			args: args{
				t:        File,
				folder:   "testdata/sub1/sub2/sub3",
				under:    "testdata/sub1/sub2",
				patterns: []string{"target*.txt"},
			},
			wantMatches: nil,
			wantErr:     ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatches, err := AcquireFromUnder(tt.args.t, tt.args.folder, tt.args.under, tt.args.patterns...)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.wantErr))
			}
			assert.Equal(t, tt.wantMatches, gotMatches)
		})
	}
}
