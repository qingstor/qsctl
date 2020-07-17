package main

import "testing"

func TestShell(t *testing.T) {
	args := []string{"cp", "-r", "--abc"}
	CpCommand.SetArgs(args)
	if err := CpCommand.Flags().Parse(args); err != nil {
		t.Fatal(err)
	}
	t.Log(cpInput)
}
