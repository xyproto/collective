package collective

import "testing"

func TestColletive(t *testing.T) {
	todo := []string{
		"Figure out what 2+2 is.",
		"Figure out if elephants are green.",
	}

	co, err := NewCollective(todo)
	if err != nil {
		t.Error(err)
	}

	co.CallUpon()
}
