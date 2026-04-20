package commands

import (
	"fmt"

	"ai-shell-windows/parser"
	"ai-shell-windows/utils"
)

func MapIntent(intent parser.Intent) (string, string, error) {
	switch intent.Action {
	case parser.ActionListFiles:
		return CommandDir, "Lists files in the current directory.", nil
	case parser.ActionListFolders:
		return CommandDirAD, "Lists folders in the current directory.", nil
	case parser.ActionPrintWorkingDir:
		return CommandCD, "Shows the current directory path.", nil
	case parser.ActionShowIPAddress:
		return CommandShowLocalIP, "Shows Windows network configuration, including IPv4 addresses.", nil
	case parser.ActionCreateFolder:
		if err := validateArgument(intent.Target); err != nil {
			return "", "", err
		}
		return fmt.Sprintf("%s %s", CommandMkdir, utils.QuoteCMDArg(intent.Target)), "Creates a folder with the requested name.", nil
	case parser.ActionDeleteFile:
		if err := validateArgument(intent.Target); err != nil {
			return "", "", err
		}
		return fmt.Sprintf("%s %s", CommandDel, utils.QuoteCMDArg(intent.Target)), "Deletes the requested file.", nil
	case parser.ActionRenameFile:
		if err := validateArgument(intent.Source); err != nil {
			return "", "", err
		}
		if err := validateArgument(intent.Destination); err != nil {
			return "", "", err
		}
		return fmt.Sprintf("%s %s %s", CommandRen, utils.QuoteCMDArg(intent.Source), utils.QuoteCMDArg(intent.Destination)), "Renames the requested file.", nil
	default:
		return "", "", parser.ErrUnknownIntent
	}
}

func validateArgument(value string) error {
	if value == "" {
		return parser.ErrMissingTarget
	}
	if utils.ContainsShellMetacharacters(value) {
		return parser.ErrUnsafeArguments
	}
	return nil
}
