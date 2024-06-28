package collective

import (
	"fmt"
	"strings"

	"github.com/xyproto/ollamaclient/v2"
)

type Agents []*Agent

type Agent struct {
	Name        string
	Description string
	Purpose     string
	Brain       *ollamaclient.Config
	Workers     Agents
	Memory      []string
}

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

// YesOrNoWithoutContext accepts a "yes/no" question and returns true if the Agent believes that the answer is "yes".
// If there are errors, or the answer is no or something else, the function returns false.
func (a *Agent) YesOrNoWithoutContext(question string) bool {
	answer, err := a.Brain.GetOutput(question + "\nAnswer either YES or NO. Only answer very briefly: YES or NO.")
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(answer), "yes")
}

// YesOrNo accepts a "yes/no" question and returns true if the Agent believes that the answer is "yes".
// If there are errors, or the answer is no or something else, the function returns false.
func (a *Agent) YesOrNo(question string) bool {
	lowercaseAnswer := strings.ToLower(a.Ask(question + "\nAnswer either YES or NO. Only answer very briefly: YES or NO."))
	return strings.Contains(lowercaseAnswer, "yes")
}

// AskWithoutContext this agent a question without context.
// Returns the error as a string if something went wrong in the "brain".
func (a *Agent) AskWithoutContext(question string) string {
	answer, err := a.Brain.GetOutput(question)
	if err != nil {
		return err.Error()
	}
	return answer
}

// Ask this agent a question, with context/memory.
// Returns the error as a string if something went wrong in the "brain".
func (a *Agent) Ask(question string) string {
	var answer string
	var err error
	if len(a.Memory) > 0 {
		const memoryEntriesUsedInPrompt = 100
		firstIndex := 0
		lastIndex := len(a.Memory) - 1
		if lastIndex > memoryEntriesUsedInPrompt {
			firstIndex = lastIndex - memoryEntriesUsedInPrompt
		}
		context := strings.Join(a.Memory[firstIndex:lastIndex], "\n")
		answer, err = a.Brain.GetOutput("This is the conversation up until now:\n```\n" + context + "\n```\n\n" + question)
	} else {
		answer, err = a.Brain.GetOutput(question)
	}
	if err != nil {
		return err.Error()
	}
	a.Memory = append(a.Memory, question)
	a.Memory = append(a.Memory, answer)
	return answer
}

func (a *Agent) Do(taskInstructions, taskContext string) error {
	var sb strings.Builder
	sb.WriteString(a.Name + ", it is time to work!\n")
	sb.WriteString("This describes you: " + a.Description + "\n")
	sb.WriteString("The context for the that you will be solving is: " + taskContext + "\n")
	sb.WriteString("A Linux bash prompt and a Python REPL are at your disposal. Just prefix commands with LINUX: or PYTHON: and you will be given the output. This can be useful for all sorts of things.\n")
	sb.WriteString("The task you will be solving is this: " + taskInstructions + "\n")
	sb.WriteString("If the task is a question and you know the answer, just say the answer, prefixed with ANSWER:\n")

	prompt1 := sb.String()

	fmt.Printf("Prompt for %s: %s.\n", a.Name, prompt1)

	action1 := a.Ask(prompt1)

	fmt.Printf("Action from %s: %s.\n", a.Name, action1)

	toolResult := a.AskWithoutContext("You are a tool. You are being called with this code: " + action1 + "\nWhat is your reply?")

	prompt2 := "The response from the tools are: " + toolResult + ". Do you want to call another function, or declare the task as complete? If it is complete, what is your conclusion or summary of the completed task?"

	fmt.Printf("Prompt for %s: %s.\n", a.Name, prompt2)

	conclusion := a.Ask(prompt2)
	fmt.Printf("Conclusion from %s: %s.\n", a.Name, conclusion)

	for {
		callAnother := a.YesOrNoWithoutContext("Does the AI want to call another tool, given this text? Answer either YES or NO. Only answer YES or NO. " + conclusion)
		if callAnother {
			fmt.Println("TO IMPLEMENT: CALL ANOTHER")
		} else {
			fmt.Println("No more calling of tools. Ending this call.")
			break
		}
	}

	fmt.Println("call upon " + a.Name + " complete")

	return nil
}
