package git

import (
	"fmt"
	"os/exec"
)

type Git struct {
	shouldPush bool
}

func New(shouldPush bool) *Git {
	return &Git{
		shouldPush: shouldPush,
	}
}

func (g *Git) GetStagedChanges() (string, error) {
	// Stage all changes
	stageCmd := exec.Command("git", "add", ".")
	if err := stageCmd.Run(); err != nil {
		return "", fmt.Errorf("error staging changes: %v", err)
	}

	// Get staged changes
	cmd := exec.Command("git", "diff", "--cached", "--name-status")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	changes := string(output)
	if changes == "" {
		return "", fmt.Errorf("no staged changes found")
	}

	return changes, nil
}

func (g *Git) Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error committing changes: %v", err)
	}

	if g.shouldPush {
		if err := g.Push(); err != nil {
			return err
		}
	}

	return nil
}

func (g *Git) Push() error {
	cmd := exec.Command("git", "push")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to push changes: %v\n%s", err, output)
	}
	return nil
}
