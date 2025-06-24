package terminal

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stxkxs/ok-cli/logger"
	"os/exec"
	"regexp"
	"time"
)

func ExecuteCommand(command string, timeout time.Duration) error {
	var output bytes.Buffer

	cmd := exec.Command("bash", "-c", command)

	cmd.Stdout = &output
	cmd.Stderr = &output

	if timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		cmd = exec.CommandContext(ctx, "bash", "-c", command)
		cmd.Stdout = &output
		cmd.Stderr = &output
	} else {
		logger.Logger.Warn().Msg("command reached timeout threshold")
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command execution failed: %s, error: %w, output: %s", command, err, output.String())
	}

	logger.Logger.Info().
		Str("command", command).
		Str("output", clean(output.String())).
		Msg("command output")

	return nil
}

func clean(input string) string {
	re := regexp.MustCompile(`[\n\t]+`)
	cleaned := re.ReplaceAllString(input, " ")
	cleaned = regexp.MustCompile(` +`).ReplaceAllString(cleaned, " ")
	return cleaned
}
