package cmd

import "testing"

func TestRootCommands(t *testing.T) {
	if command, _, err := rootCmd.Find([]string{"get"}); err != nil || command != getCmd {
		t.Fatalf("get command not registered: command=%v err=%v", command, err)
	}
	if command, _, err := rootCmd.Find([]string{"set"}); err != nil || command != setCmd {
		t.Fatalf("set command not registered: command=%v err=%v", command, err)
	}
	if command, _, err := rootCmd.Find([]string{"display"}); err != nil || command != displayCmd {
		t.Fatalf("display command not registered: command=%v err=%v", command, err)
	}

	for _, oldName := range []string{"list", "show"} {
		command, _, err := rootCmd.Find([]string{oldName})
		if err == nil && command != rootCmd {
			t.Fatalf("old %q command is still registered", oldName)
		}
	}
}

func TestGetAndSetResources(t *testing.T) {
	for _, name := range []string{"image", "gif", "dashboard", "animation"} {
		getResource, _, getErr := getCmd.Find([]string{name})
		if getErr != nil || getResource == getCmd {
			t.Errorf("get %s command not registered: command=%v err=%v", name, getResource, getErr)
		}

		setResource, _, setErr := setCmd.Find([]string{name})
		if setErr != nil || setResource == setCmd {
			t.Errorf("set %s command not registered: command=%v err=%v", name, setResource, setErr)
		}
	}
}
