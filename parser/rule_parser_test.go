package parser

import "testing"

func TestRuleParserParse(t *testing.T) {
	t.Parallel()

	parser := NewRuleParser()

	tests := []struct {
		name       string
		input      string
		wantAction string
		wantTarget string
		wantSource string
		wantDest   string
		wantErr    error
	}{
		{
			name:       "list files",
			input:      "show files in this folder",
			wantAction: ActionListFiles,
			wantTarget: "current_directory",
		},
		{
			name:       "list all files in current directory",
			input:      "list all files in current directory",
			wantAction: ActionListFiles,
			wantTarget: "current_directory",
		},
		{
			name:       "list folders in directory",
			input:      "what are the folders in this directory",
			wantAction: ActionListFolders,
			wantTarget: "current_directory",
		},
		{
			name:       "show directory contents",
			input:      "show directory contents",
			wantAction: ActionListFiles,
			wantTarget: "current_directory",
		},
		{
			name:       "print working directory",
			input:      "current directory",
			wantAction: ActionPrintWorkingDir,
			wantTarget: "current_directory",
		},
		{
			name:       "show ip address",
			input:      "what is my ip address",
			wantAction: ActionShowIPAddress,
		},
		{
			name:       "create folder",
			input:      "create folder test",
			wantAction: ActionCreateFolder,
			wantTarget: "test",
		},
		{
			name:       "delete file",
			input:      "delete file report.txt",
			wantAction: ActionDeleteFile,
			wantTarget: "report.txt",
		},
		{
			name:       "rename file",
			input:      "rename file old.txt to new.txt",
			wantAction: ActionRenameFile,
			wantSource: "old.txt",
			wantDest:   "new.txt",
		},
		{
			name:    "ambiguous delete",
			input:   "delete file",
			wantErr: ErrMissingTarget,
		},
		{
			name:    "unknown input",
			input:   "dance for me",
			wantErr: ErrUnknownIntent,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			intent, err := parser.Parse(tt.input)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if intent.Action != tt.wantAction {
				t.Fatalf("expected action %q, got %q", tt.wantAction, intent.Action)
			}
			if intent.Target != tt.wantTarget {
				t.Fatalf("expected target %q, got %q", tt.wantTarget, intent.Target)
			}
			if intent.Source != tt.wantSource {
				t.Fatalf("expected source %q, got %q", tt.wantSource, intent.Source)
			}
			if intent.Destination != tt.wantDest {
				t.Fatalf("expected destination %q, got %q", tt.wantDest, intent.Destination)
			}
		})
	}
}
