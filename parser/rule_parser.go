package parser

import (
	"regexp"
	"slices"
	"strings"

	"ai-shell-windows/utils"
)

var (
	createFolderPattern = regexp.MustCompile(`^(create|make)\s+folder\s+(.+)$`)
	deleteFilePattern   = regexp.MustCompile(`^(delete|remove)\s+file(?:\s+(.+))?$`)
	renameFilePattern   = regexp.MustCompile(`^(rename)\s+file\s+(.+?)\s+to\s+(.+)$`)
)

type RuleParser struct{}

func NewRuleParser() RuleParser {
	return RuleParser{}
}

func (RuleParser) Parse(input string) (Intent, error) {
	normalized := utils.NormalizeText(input)
	intent := Intent{
		RawInput:   input,
		Normalized: normalized,
		Action:     ActionUnknown,
	}

	switch {
	case normalized == "":
		return intent, ErrUnknownIntent
	case isListFolders(normalized):
		intent.Action = ActionListFolders
		intent.Target = "current_directory"
		intent.Explanation = "Shows only folders in the current directory."
		return intent, nil
	case isListFiles(normalized):
		intent.Action = ActionListFiles
		intent.Target = "current_directory"
		intent.Explanation = "Lists files in the current directory."
		return intent, nil
	case isPrintWorkingDir(normalized):
		intent.Action = ActionPrintWorkingDir
		intent.Target = "current_directory"
		intent.Explanation = "Shows the current directory path."
		return intent, nil
	case createFolderPattern.MatchString(normalized):
		matches := createFolderPattern.FindStringSubmatch(normalized)
		intent.Action = ActionCreateFolder
		intent.Target = strings.TrimSpace(matches[2])
		intent.Explanation = "Creates a folder with the requested name."
		if intent.Target == "" {
			intent.RequiresInfo = true
			intent.Clarification = "Which folder name should I create?"
			return intent, ErrMissingTarget
		}
		return intent, nil
	case deleteFilePattern.MatchString(normalized):
		matches := deleteFilePattern.FindStringSubmatch(normalized)
		intent.Action = ActionDeleteFile
		intent.Target = strings.TrimSpace(matches[2])
		intent.Explanation = "Deletes the requested file."
		if intent.Target == "" {
			intent.RequiresInfo = true
			intent.Clarification = "Which file should I delete?"
			return intent, ErrMissingTarget
		}
		return intent, nil
	case renameFilePattern.MatchString(normalized):
		matches := renameFilePattern.FindStringSubmatch(normalized)
		intent.Action = ActionRenameFile
		intent.Source = strings.TrimSpace(matches[2])
		intent.Destination = strings.TrimSpace(matches[3])
		intent.Explanation = "Renames the requested file."
		if intent.Source == "" {
			intent.RequiresInfo = true
			intent.Clarification = "Which file should I rename?"
			return intent, ErrMissingTarget
		}
		if intent.Destination == "" {
			intent.RequiresInfo = true
			intent.Clarification = "What should the new name be?"
			return intent, ErrMissingNewName
		}
		return intent, nil
	default:
		return intent, ErrUnknownIntent
	}
}

func isListFolders(input string) bool {
	tokens := strings.Fields(input)

	if hasAnyToken(tokens, "folder", "folders", "directory", "directories") &&
		hasAnyToken(tokens, "show", "list", "display") &&
		!hasAnyToken(tokens, "file", "files") {
		return true
	}

	return containsPhrase(input,
		"what are the folders in this directory",
		"list folders",
		"show folders",
		"show only folders",
	)
}

func isListFiles(input string) bool {
	tokens := strings.Fields(input)

	if hasAnyToken(tokens, "file", "files", "folder", "folders", "directory", "directories") &&
		hasAnyToken(tokens, "show", "list", "display") {
		return true
	}

	return containsPhrase(input,
		"show files",
		"list files",
		"show files in this folder",
		"show directory",
		"list directory",
		"directory contents",
	)
}

func isPrintWorkingDir(input string) bool {
	tokens := strings.Fields(input)

	if containsPhrase(input,
		"print working dir",
		"print working directory",
		"where am i",
		"show current path",
		"what folder am i in",
	) {
		return true
	}

	if hasAllTokens(tokens, "current", "directory") && !hasAnyToken(tokens, "file", "files", "folder", "folders") {
		return true
	}

	return false
}

func containsPhrase(input string, phrases ...string) bool {
	for _, phrase := range phrases {
		if strings.Contains(input, phrase) {
			return true
		}
	}

	return false
}

func hasAnyToken(tokens []string, values ...string) bool {
	for _, token := range tokens {
		if slices.Contains(values, token) {
			return true
		}
	}

	return false
}

func hasAllTokens(tokens []string, values ...string) bool {
	for _, value := range values {
		if !slices.Contains(tokens, value) {
			return false
		}
	}

	return true
}
