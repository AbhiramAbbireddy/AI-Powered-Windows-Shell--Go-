package safety

import "testing"

func TestIsDangerous(t *testing.T) {
	t.Parallel()

	tests := []struct {
		command       string
		wantDangerous bool
	}{
		{command: `dir`, wantDangerous: false},
		{command: `del "test.txt"`, wantDangerous: false},
		{command: `del *`, wantDangerous: true},
		{command: `format c:`, wantDangerous: true},
		{command: `shutdown /s /t 0`, wantDangerous: true},
		{command: `rd /s /q temp`, wantDangerous: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.command, func(t *testing.T) {
			t.Parallel()

			got, _ := IsDangerous(tt.command)
			if got != tt.wantDangerous {
				t.Fatalf("expected %v for %q, got %v", tt.wantDangerous, tt.command, got)
			}
		})
	}
}
