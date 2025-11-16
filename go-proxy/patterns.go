package main

import (
	"fmt"
	"io"
	"os"
	// "time" // Temporarily commented out while BMad output rule is disabled
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
		// Simply append the custom text (no clearing needed)
		appendedText := ", this is a custom text"

		if interceptor.GetLogFile() != nil {
			fmt.Fprintf(interceptor.GetLogFile(), "[CLAUDEX] Appending: '%s'\n", appendedText)
		}

		writer.Write([]byte(appendedText))

		// Return false to let the original ENTER go through
		return false
	})
	if err != nil {
		return err
	}

	// Rule 2: Capture any substring that would autocomplete to /BMad:agents:dev
	// This matches: /d, /de, /dev, /b, /bm, /bmad, /a, /ag, /agents, etc.
	err = interceptor.AddInputRule(`(?i)^/(d|de|dev|b|bm|bma|bmad|a|ag|age|agen|agent|agents)$`, func(input string, writer io.Writer) bool {
		// Append custom text when this pattern is detected
		appendedText := " - BMad agents development mode activated"

		if interceptor.GetLogFile() != nil {
			fmt.Fprintf(interceptor.GetLogFile(), "[CLAUDEX] BMad command pattern '%s' matched, appending: '%s'\n", input, appendedText)
		}

		writer.Write([]byte(appendedText))

		// Return false to let the original ENTER go through
		return false
	})
	if err != nil {
		return err
	}

	// Rule 3: Replace "goodbye" with test message in user input
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

	// Rule 4: Detect when /BMad:agents:dev is running and interrupt it
	// TEMPORARILY DISABLED FOR DEBUGGING
	/*
		err = interceptor.AddOutputRule(`/BMad:agents:dev is running`, func(input string, writer io.Writer) bool {
			if interceptor.GetLogFile() != nil {
				fmt.Fprintf(interceptor.GetLogFile(), "\n========== BMAD OUTPUT INTERCEPTION START ==========\n")
				fmt.Fprintf(interceptor.GetLogFile(), "[CLAUDEX] Detected '/BMad:agents:dev is running' in output\n")
			}

			ptyWriter := interceptor.GetPtyWriter()
			if ptyWriter == nil {
				if interceptor.GetLogFile() != nil {
					fmt.Fprintf(interceptor.GetLogFile(), "[CLAUDEX ERROR] PTY writer not available!\n")
					fmt.Fprintf(interceptor.GetLogFile(), "========== BMAD OUTPUT INTERCEPTION FAILED ==========\n\n")
				}
				return false
			}

			// Send ESC to interrupt Claude (0x1B is the ESC byte)
			if interceptor.GetLogFile() != nil {
				fmt.Fprintf(interceptor.GetLogFile(), "[CLAUDEX] Step 1: Sending ESC (byte 0x1B) to interrupt Claude\n")
			}
			n, err := ptyWriter.Write([]byte{0x1B})
			if interceptor.GetLogFile() != nil {
				fmt.Fprintf(interceptor.GetLogFile(), "[CLAUDEX] ESC sent: wrote %d bytes, error: %v\n", n, err)
			}

			// Wait for Claude to process the interrupt and return to prompt
			if interceptor.GetLogFile() != nil {
				fmt.Fprintf(interceptor.GetLogFile(), "[CLAUDEX] Waiting 200ms for interrupt to take effect...\n")
			}
			time.Sleep(200 * time.Millisecond)

			// Send the modified command character by character (simulating typing)
			modifiedCommand := "/BMad:agents:dev with custom text"

			if interceptor.GetLogFile() != nil {
				fmt.Fprintf(interceptor.GetLogFile(), "[CLAUDEX] Step 2: Typing modified command character by character: '%s'\n", modifiedCommand)
			}

			for i, char := range modifiedCommand {
				ptyWriter.Write([]byte{byte(char)})
				// Small delay between characters to simulate human typing
				if i < len(modifiedCommand)-1 {
					time.Sleep(5 * time.Millisecond)
				}
			}

			if interceptor.GetLogFile() != nil {
				fmt.Fprintf(interceptor.GetLogFile(), "[CLAUDEX] Command text typed character by character\n")
			}

			// Wait a bit longer for the text to be processed
			if interceptor.GetLogFile() != nil {
				fmt.Fprintf(interceptor.GetLogFile(), "[CLAUDEX] Waiting 300ms before sending Enter...\n")
			}
			time.Sleep(300 * time.Millisecond)

			// Send Enter once
			if interceptor.GetLogFile() != nil {
				fmt.Fprintf(interceptor.GetLogFile(), "[CLAUDEX] Step 3: Sending Enter byte 13 (0x0D = \\r)\n")
			}
			n, err = ptyWriter.Write([]byte{13}) // byte 13 = 0x0D = \r
			if interceptor.GetLogFile() != nil {
				fmt.Fprintf(interceptor.GetLogFile(), "[CLAUDEX] Enter sent: wrote %d bytes, error: %v\n", n, err)
			}

			if interceptor.GetLogFile() != nil {
				fmt.Fprintf(interceptor.GetLogFile(), "[CLAUDEX] Interception complete!\n")
				fmt.Fprintf(interceptor.GetLogFile(), "========== BMAD OUTPUT INTERCEPTION END ==========\n\n")
			}

			return false // Don't block output, we want to see what happens
		})
		if err != nil {
			return err
		}
	*/

	// Rule 5: Detect "hello world" in output
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
