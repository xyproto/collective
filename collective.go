package collective

type Collective struct {
	Agents Agents
	TODO   []string
}

func NewCollective(todoList []string) (*Collective, error) {
	frank, err := NewAgent("Frank", "CEO", "llama3", false, "be just, fair, display excellence and performance")
	if err != nil {
		return nil, err
	}

	bob, err := NewAgent("Bob", "Artist", "tinyllama", true, "be creative and amusing")
	if err != nil {
		return nil, err
	}

	alice, err := NewAgent("Alice", "Developer", "llama3", false, "develop the best source code the world has ever seen")
	if err != nil {
		return nil, err
	}

	return &Collective{
		Agents: []*Agent{frank, bob, alice},
		TODO:   todoList,
	}, nil
}

func (co *Collective) CallUpon() {
	for _, agent := range co.Agents {
		agent.CallUpon(&co.TODO, co.Agents)
	}
}
