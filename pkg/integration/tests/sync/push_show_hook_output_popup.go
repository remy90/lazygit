package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PushShowHookOutputPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show hook output in a live popup while pushing when enabled",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.ShowHookOutput = true
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")

		shell.CloneIntoRemote("origin")

		shell.SetBranchUpstream("master", "origin/master")

		shell.EmptyCommit("two")

		shell.CreateFile(
			".git/hooks/pre-push",
			"#!/bin/sh\necho 'Foo'\nsleep 0.5\necho 'Bar'\nsleep 0.5\necho 'Baz'\nexit 0\n",
		)
		shell.MakeExecutable(".git/hooks/pre-push")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().Content(Equals("↑1 repo → master"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Push)

		t.ExpectPopup().Alert().
			Title(Equals("Git output:")).
			Content(Contains("Foo"))

		t.Wait(900)

		t.ExpectPopup().Alert().
			Title(Equals("Git output:")).
			Content(Contains("Bar")).
			Confirm()

		assertSuccessfullyPushed(t)
	},
})
