package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PushNoHookOutputPopupByDefault = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Do not show hook output popup after a successful push by default",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")

		shell.CloneIntoRemote("origin")

		shell.SetBranchUpstream("master", "origin/master")

		shell.EmptyCommit("two")

		shell.CreateFile(
			".git/hooks/pre-push",
			"#!/bin/sh\necho 'husky - pre-push checks passed'\nexit 0\n",
		)
		shell.MakeExecutable(".git/hooks/pre-push")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().Content(Equals("↑1 repo → master"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Push)

		assertSuccessfullyPushed(t)
	},
})
