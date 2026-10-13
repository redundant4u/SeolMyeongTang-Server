package validator

import "testing"

func TestBoundaryValidation(t *testing.T) {
	v := New()
	for _, tc := range []struct {
		label, id string
		valid     bool
	}{
		{"", "abcdefgh", true}, {"Demo-1_label.2", "abc12345", true},
		{"-bad", "abcdefgh", false}, {"bad!", "abcdefgh", false}, {"ok", "ABC12345", false}, {"ok", "abc", false}, {"ok", "", false},
	} {
		t.Run(tc.label+tc.id, func(t *testing.T) {
			err := v.Validate(struct {
				Label string `validate:"k8slabel"`
				ID    string `validate:"sessionId"`
			}{tc.label, tc.id})
			if (err == nil) != tc.valid {
				t.Fatalf("validation %q %q: %v", tc.label, tc.id, err)
			}
		})
	}
}
