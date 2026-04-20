package shell

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"ai-shell-windows/commands"
	"ai-shell-windows/config"
	"ai-shell-windows/parser"
	"ai-shell-windows/safety"
)

func StartShell(cfg config.Config) error {
	ruleParser := parser.NewRuleParser()
	aiParser := parser.NewAIParser(cfg.GroqAPIKey, cfg.GroqBaseURL, cfg.GroqModel, cfg.AIRequestTimeout)
	scanner := bufio.NewScanner(os.Stdin)

	PrintWelcome()

	for {
		PrintPrompt()
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return err
			}
			fmt.Println()
			return nil
		}

		rawInput := scanner.Text()
		normalized := strings.TrimSpace(strings.ToLower(rawInput))
		if normalized == "" {
			continue
		}

		if normalized == "exit" || normalized == "quit" {
			fmt.Println("Bye.")
			return nil
		}

		intent, err := ruleParser.Parse(rawInput)
		if err != nil {
			if errors.Is(err, parser.ErrUnknownIntent) && cfg.EnableAI {
				intent, err = aiParser.Parse(rawInput)
			}
			if err != nil {
				switch {
				case errors.Is(err, parser.ErrMissingTarget):
					RenderMissingInfo("Which file or folder should I use?")
				case errors.Is(err, parser.ErrMissingNewName):
					RenderMissingInfo("What should the new name be?")
				case errors.Is(err, parser.ErrUnknownIntent):
					RenderUnknownInput()
				case errors.Is(err, parser.ErrAIUnavailable):
					RenderMissingInfo("Groq AI parsing is not configured. Set GROQ_API_KEY and try again.")
				default:
					RenderMissingInfo(fmt.Sprintf("AI parsing failed: %v", err))
				}
				continue
			}
		}

		if intent.RequiresInfo {
			if intent.Clarification != "" {
				RenderMissingInfo(intent.Clarification)
			} else {
				RenderMissingInfo("Which file or folder should I use?")
			}
			continue
		}

		command, explanation, err := commands.MapIntent(intent)
		if err != nil {
			switch {
			case errors.Is(err, parser.ErrUnsafeArguments):
				RenderMissingInfo("That input contains unsupported shell characters, so I blocked it.")
			default:
				RenderUnknownInput()
			}
			continue
		}
		if intent.Explanation != "" {
			explanation = intent.Explanation
		}

		if cfg.PreviewCommands {
			RenderCommand(command, explanation)
		}

		if cfg.DangerousPrompts {
			if dangerous, reason := safety.IsDangerous(command); dangerous {
				confirmed, confirmErr := AskConfirmation(command, reason)
				if confirmErr != nil {
					return confirmErr
				}
				if !confirmed {
					fmt.Println("Command canceled.")
					continue
				}
			}
		}

		result, execErr := ExecuteCommand(cfg.Shell, command)
		RenderOutput(result, execErr)

		if historyErr := AppendHistory(cfg.HistoryPath, HistoryEntry{
			UserInput: rawInput,
			Command:   command,
		}); historyErr != nil {
			fmt.Printf("History warning: %v\n", historyErr)
		}
	}
}
