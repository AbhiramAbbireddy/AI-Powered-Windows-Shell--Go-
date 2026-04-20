package parser

import "errors"

const (
	ActionListFiles       = "list_files"
	ActionListFolders     = "list_folders"
	ActionPrintWorkingDir = "print_working_dir"
	ActionShowIPAddress   = "show_ip_address"
	ActionDeleteFile      = "delete_file"
	ActionCreateFolder    = "create_folder"
	ActionRenameFile      = "rename_file"
	ActionUnknown         = "unknown"
)

var (
	ErrUnknownIntent   = errors.New("unknown input")
	ErrMissingTarget   = errors.New("missing target")
	ErrMissingNewName  = errors.New("missing new name")
	ErrUnsafeArguments = errors.New("input contains unsupported shell characters")
	ErrAIUnavailable   = errors.New("ai parser is not configured")
	ErrAIResponse      = errors.New("ai parser returned an invalid response")
)

type Intent struct {
	RawInput      string
	Normalized    string
	Action        string
	Target        string
	Source        string
	Destination   string
	RequiresInfo  bool
	Clarification string
	Explanation   string
}
