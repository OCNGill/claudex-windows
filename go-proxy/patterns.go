package main

import (
	"fmt"
	"io"
	"os"
)

const (
	colorReset  = "\x1b[0m"
	colorCyan   = "\x1b[36m"
	colorGreen  = "\x1b[32m"
	colorYellow = "\x1b[33m"
)

// SetupPatterns configures all the pattern matching rules
func SetupPatterns(interceptor *Interceptor) error {
	// INPUT RULES - Only checked when user presses ENTER

	// Rule 1: Append custom text to "hello" pattern in user input
	err := interceptor.AddInputRule(`(?i)hello`, func(input string, writer io.Writer) bool {
		// Clear the typed "hello" by sending backspaces
		inputLen := len(input)
		for i := 0; i < inputLen; i++ {
			writer.Write([]byte{0x08}) // Send backspace (Ctrl-H)
		}

		// Now type the complete message with appended text
		completeMessage := input + ", this is a custom text"
		writer.Write([]byte(completeMessage))

		// Return false to let the original ENTER go through
		return false
	})
	if err != nil {
		return err
	}

	// Rule 2: Replace "goodbye" with test message in user input
	err = interceptor.AddInputRule(`(?i)goodbye`, func(input string, writer io.Writer) bool {
		customMessage := fmt.Sprintf("\n%s[Claudex]%s %sIntercepted \"goodbye\" - sending different message to Claude...%s\n",
			colorYellow, colorReset, colorGreen, colorReset)
		// Write to stderr for user notification
		os.Stderr.WriteString(customMessage)

		// Send the replacement message to Claude via PTY
		replacementMessage := "the test we are doing is working\r"
		writer.Write([]byte(replacementMessage))

		return true // Block original, we sent replacement
	})
	if err != nil {
		return err
	}

	// OUTPUT RULES - Checked continuously on Claude's output

	// Rule 3: Detect "hello world" in output
	err = interceptor.AddOutputRule(`(?i)hello world`, func(input string, writer io.Writer) bool {
		customMessage := fmt.Sprintf("\n%s[Claudex]%s %s🎉 Hello World detected in OUTPUT!%s\n",
			colorYellow, colorReset, colorCyan, colorReset)
		os.Stderr.WriteString(customMessage)
		return false // Don't block, just notify
	})
	if err != nil {
		return err
	}

	return nil
}
