package handler

import "testing"

func TestSanitizePath(t *testing.T) {
	tests := []struct {
		input    string
		wantOK   bool
		wantPath string
	}{
		{"foo/bar.txt", true, "foo/bar.txt"},
		{"/foo/bar", true, "/foo/bar"},
		{"foo.txt", true, "foo.txt"},
		{"", true, "."},
		{"../etc/passwd", false, ""},
		{"foo/../../bar", false, ""},
		{"../../etc/shadow", false, ""},
		{"foo/..bar", true, "foo/..bar"}, // ".." in filename is OK, only ".." as segment is rejected
		{"foo/bar/../../../baz", false, ""},
		{"a/b/c/d.txt", true, "a/b/c/d.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			gotPath, gotOK := sanitizePath(tt.input)
			if gotOK != tt.wantOK {
				t.Errorf("sanitizePath(%q) ok = %v, want %v", tt.input, gotOK, tt.wantOK)
			}
			if gotOK && gotPath != tt.wantPath {
				t.Errorf("sanitizePath(%q) path = %q, want %q", tt.input, gotPath, tt.wantPath)
			}
		})
	}
}
