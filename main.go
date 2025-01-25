package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

// Path to the directory where docker-compose.yml is located
const dockerComposeDir = "/home/techops/musicbrainz-docker/"

// Slack Webhook URL (Replace with your actual URL)
const slackWebhookURL = "YOUR KEY GOES HERE"

// sendToSlack sends a message to Slack
func sendToSlack(message string) error {

	payload := map[string]string{"text": message}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(slackWebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Slack webhook failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// fetchLast10Lines gets the last 10 lines from the log file, running docker compose from the correct directory
func fetchLast10Lines() ([]string, error) {
	cmd := exec.Command("docker", "compose", "exec", "musicbrainz", "/usr/bin/tail", "-n", "10", "mirror.log")
	cmd.Dir = dockerComposeDir // Set working directory where docker-compose.yml exists

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return nil, fmt.Errorf("Error creating StdoutPipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("Error starting command: %v", err)
	}

	var lines []string
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("Command execution error: %v", err)
	}

	return lines, nil
}

func main() {
	lines, err := fetchLast10Lines()

	if err != nil {
		fmt.Println("Error fetching logs:", err)
	} else if len(lines) > 0 {
		message := strings.Join(lines, "\n")
		//fmt.Println("Sending logs to Slack...")
		if err := sendToSlack(message); err != nil {
			fmt.Println("Error sending to Slack:", err)
		} else {
			fmt.Println("Sent logs to Slack successfully.")
		}
	} else {
		fmt.Println("No logs found.")
	}

}
