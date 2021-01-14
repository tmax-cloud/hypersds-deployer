package util

import(
	"testing"
	"os/user"
)

func TestRunSSHCmd(t *testing.T) {
	user, _ := user.Current()
	testCommand := []string{"ls", "-alh"}

	result, err := RunSSHCmd(user.Username, "localhost", testCommand...)

	t.Log(result.String())

	if err != nil {
		t.Fatal(err)
	}
}
