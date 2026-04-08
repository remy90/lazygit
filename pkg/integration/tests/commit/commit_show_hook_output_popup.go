package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitShowHookOutputPopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show hook output in a live popup while committing when enabled",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.ShowHookOutput = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("myfile", "myfile content")
		shell.CreateFile(
			".git/hooks/pre-commit",
			"#!/bin/sh\necho 'Foo'\nsleep 0.5\necho 'Bar'\nsleep 0.5\necho 'Baz'\nexit 0\n",
		)
		shell.MakeExecutable(".git/hooks/pre-commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ?? myfile"),
			).
			SelectNextItem().
			PressPrimaryAction().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().Type("my commit message").Confirm()

		t.ExpectPopup().Alert().
			Title(Equals("Git output:")).
			Content(Contains("Foo"))

		t.Wait(900)

		t.ExpectPopup().Alert().
			Title(Equals("Git output:")).
			Content(Contains("Bar")).
			Confirm()

		t.Views().Commits().
			Focus().
			Lines(
				Contains("my commit message").IsSelected(),
			)
	},
})
