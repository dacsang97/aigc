package main

import (
	"fmt"
	"os/exec"
)

func getGitChanges() (string, error) {
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

func commitChanges(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error committing changes: %v", err)
	}

	if push {
		cmd = exec.Command("git", "push", "origin")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error pushing changes: %v", err)
		}
		fmt.Println("Successfully pushed changes to origin")
	}

	return nil
}
