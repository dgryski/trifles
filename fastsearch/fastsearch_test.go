package fastsearch

import "testing"

func TestSearch(t *testing.T) {

	var tests = []struct {
		needle   string
		haystack string
		i        int
	}{
		{"hello", "hello", 0},
		{"hello", " hello", 1},
		{"hello", " hell", -1},
		{"ell", "hello", 1},
		{"hellp", " hellohellohello", -1},
	}

	for _, tt := range tests {
		if i := fastsearch([]byte(tt.haystack), []byte(tt.needle)); i != tt.i {
			t.Errorf("fastsearch(%q,%q)=%v, want %v", tt.haystack, tt.needle, i, tt.i)
		}
	}
}
