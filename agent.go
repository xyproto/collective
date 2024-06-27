package collective

import (
	"errors"
	"fmt"
	"strings"

	"github.com/xyproto/ask"
	"github.com/xyproto/ollamaclient/v2"
)

type Agent struct {
	Name        string
	Description string
	Purpose     string
	Brain       *ollamaclient.Config
}

type Agents []*Agent

func NewAgent(name, description, modelName string, creative bool, purpose string) (*Agent, error) {
	oc := ollamaclient.New()
	oc.ModelName = modelName

	err := oc.PullIfNeeded(true)
	if err != nil {
		return nil, err
	}

	if !creative {
		oc.SetReproducible()
	}

	return &Agent{
		Name:        name,
		Description: description,
		Purpose:     purpose,
		Brain:       oc,
	}, nil
}

func (a *Agent) CallUpon(TODO *[]string, coWorkers Agents) error {
	var currentTask string
	if len(*TODO) > 0 {
		firstItem := (*TODO)[0]
		currentTask = firstItem
	} else {
		return errors.New("no tasks to complete, ending call")
	}

	fmt.Println("Current task: " + currentTask)

	currentTools := `
	// This function searches Wikipedia and returns the first result, as a string
	searchWikipedia(string) -> string
	// This function searches Google and returns the first result, as a string
	searchGoogle(string) -> string
	// This function sends an e-mail, given an address, a subject and a body
	sendEmail(string, string, string)
	// This function returns the last received non-archived e-mail as a string, or an empty string if no e-mail is available
	lastEmail() -> string
	// This function archives the last received non-archived e-mail, or does nothing if no e-mail is available
	archiveLastEmail()
	`

	prompt1 := fmt.Sprintf("%s, it is time to work. This describes you: %s Your current task is: %s The tools at your disposal that you can call as if they were Lua functions are: %s. \nWrite the Lua code that calls one of these tools, and I will return the response to you. The goal is to complete the task, or at least try to complete the task in up to three different ways. If you know the answer and do not need the tools, then just output the answer.", a.Name, a.Description, currentTask, currentTools)

	fmt.Printf("Prompt for %s: %s.\n", a.Name, prompt1)

	action1, err := a.Brain.GetOutput(prompt1)
	if err != nil {
		return err
	}

	fmt.Printf("Action from %s: %s.\n", a.Name, action1)

	toolResult := ask.Ask("If the AI is attempting to call a tool, what should the response from the tool be? ")

	prompt2 := "The response from the tools are: " + toolResult + ". Do you want to call another function, or declare the task as complete? If it is complete, what is your conclusion or summary of the completed task?"

	fmt.Printf("Prompt for %s: %s.\n", a.Name, prompt2)

	conclusion, err := a.Brain.GetOutput(prompt2)
	if err != nil {
		return err
	}

	fmt.Printf("Conclusion from %s: %s.\n", a.Name, conclusion)

	for {
		callAnother := ask.Ask("Does the AI want to call another tool, given this text? Answer either YES or NO. Only answer YES or NO. " + conclusion)
		if strings.Contains(strings.ToLower(callAnother), "yes") {
			fmt.Println("TO IMPLEMENT")
		} else {
			fmt.Println("No more calling of tools. Ending this call.")
			break
		}
	}

	fmt.Println("call upon " + a.Name + " complete")

	return nil
}
