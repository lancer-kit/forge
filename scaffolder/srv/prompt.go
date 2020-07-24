package srv

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

func askPromptInput(message, defaultValue string) (string, error) {
	var promptAnswer string

	promptQuestion := &survey.Input{
		Message: message,
		Default: defaultValue,
	}

	err := survey.AskOne(promptQuestion, &promptAnswer, survey.WithValidator(survey.Required))
	if err != nil {
		return "", fmt.Errorf("failed to get the srv answer: %s", err)
	}
	return promptAnswer, nil
}

func askPromptSelect(message, defaultValue string, options []string) (string, error) {
	var promptAnswer string

	promptQuestion := &survey.Select{
		Message: message,
		Options: options,
		Default: defaultValue,
	}

	err := survey.AskOne(promptQuestion, &promptAnswer, survey.WithValidator(survey.Required))
	if err != nil {
		return "", fmt.Errorf("failed to get the srv answer: %s", err)
	}
	return promptAnswer, nil
}
