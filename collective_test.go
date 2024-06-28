package collective

import "testing"

func TestColletive(t *testing.T) {
	todo := []string{
		"Figure out what 2+2 is.",
		"Figure out if elephants are green.",
	}

	_, err := NewCollective(todo)
	if err != nil {
		t.Error(err)
	}

	//co.CallUpon()
}

func TestAsk(t *testing.T) {
	mallory, err := NewAgent("Mallory", "CEO", "gemma2", false, "just, fair, displays excellence and performance")
	if err != nil {
		t.Error(err)
	}
	if mallory.YesOrNo("cows are green") {
		t.Fail()
	}
	if !mallory.YesOrNo("some frogs are green") {
		t.Fail()
	}

}
