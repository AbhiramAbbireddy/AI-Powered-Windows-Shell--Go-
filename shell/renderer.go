package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func PrintWelcome() {
	fmt.Println("AI Shell for Windows")
	fmt.Println(`Type natural language commands, or "exit" to quit.`)
}

func PrintPrompt() {
	fmt.Print(">> ")
}

func RenderCommand(command, explanation string) {
	fmt.Printf("AI Command: %s\n", command)
	if explanation != "" {
		fmt.Printf("Why: %s\n", explanation)
	}
}

func RenderOutput(result ExecutionResult, err error) {
	if err != nil {
		fmt.Printf("Execution error: %v\n", err)
	}

	output := strings.TrimSpace(result.Stdout)
	errors := strings.TrimSpace(result.Stderr)

	if output != "" {
		fmt.Println("Output:")
		fmt.Println(output)
	}

	if errors != "" {
		fmt.Println("Error:")
		fmt.Println(errors)
	}

	if output == "" && errors == "" && err == nil {
		fmt.Println("Command completed successfully.")
	}
}

func RenderUnknownInput() {
	fmt.Println("I don't understand. Try again.")
}

func RenderMissingInfo(message string) {
	fmt.Println(message)
}

func AskConfirmation(command, reason string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Confirm execution of %q? %s [y/N]: ", command, reason)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}
