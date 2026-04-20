package commands

import (
	"testing"

	"ai-shell-windows/parser"
)

func TestMapIntent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		intent      parser.Intent
		wantCommand string
		wantErr     error
	}{
		{
			name: "list files",
			intent: parser.Intent{
				Action: parser.ActionListFiles,
			},
			wantCommand: "dir",
		},
		{
			name: "current directory",
			intent: parser.Intent{
				Action: parser.ActionPrintWorkingDir,
			},
			wantCommand: "cd",
		},
		{
			name: "list folders",
			intent: parser.Intent{
				Action: parser.ActionListFolders,
			},
			wantCommand: "dir /ad",
		},
		{
			name: "create folder",
			intent: parser.Intent{
				Action: parser.ActionCreateFolder,
				Target: "demo folder",
			},
			wantCommand: `mkdir "demo folder"`,
		},
		{
			name: "delete file",
			intent: parser.Intent{
				Action: parser.ActionDeleteFile,
				Target: "test.txt",
			},
			wantCommand: `del "test.txt"`,
		},
		{
			name: "rename file",
			intent: parser.Intent{
				Action:      parser.ActionRenameFile,
				Source:      "old.txt",
				Destination: "new.txt",
			},
			wantCommand: `ren "old.txt" "new.txt"`,
		},
		{
			name: "block shell metacharacters",
			intent: parser.Intent{
				Action: parser.ActionDeleteFile,
				Target: `a.txt & del *`,
			},
			wantErr: parser.ErrUnsafeArguments,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			command, _, err := MapIntent(tt.intent)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if command != tt.wantCommand {
				t.Fatalf("expected command %q, got %q", tt.wantCommand, command)
			}
		})
	}
}
