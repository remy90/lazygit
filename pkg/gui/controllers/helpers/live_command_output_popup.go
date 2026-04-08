package helpers

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

const maxLiveCommandOutputLength = 20000

type liveCommandOutputPopupGUI interface {
	CreatePopupPanel(opts types.CreatePopupPanelOpts)
	Context() types.IContextMgr
	OnUIThread(f func() error)
	SetViewContent(view *gocui.View, content string)
	Views() types.Views
	State() types.IStateAccessor
}

type LiveCommandOutputPopup struct {
	gui   liveCommandOutputPopupGUI
	title string

	mu      sync.Mutex
	content []byte
	opened  bool
	key     string

	lastRenderedContent string
}

var nextLiveCommandOutputPopupID uint64

func NewLiveCommandOutputPopup(gui liveCommandOutputPopupGUI, title string) *LiveCommandOutputPopup {
	id := atomic.AddUint64(&nextLiveCommandOutputPopupID, 1)

	return &LiveCommandOutputPopup{
		gui:   gui,
		title: title,
		key:   fmt.Sprintf("live-command-output-%d", id),
	}
}

func (self *LiveCommandOutputPopup) OnOutput(chunk string) {
	if chunk == "" {
		return
	}

	self.mu.Lock()
	self.content = append(self.content, chunk...)
	if len(self.content) > maxLiveCommandOutputLength {
		self.content = self.content[len(self.content)-maxLiveCommandOutputLength:]
	}
	self.mu.Unlock()

	self.gui.OnUIThread(func() error {
		self.mu.Lock()
		latestContent := strings.TrimSpace(string(self.content))
		if !self.opened {
			self.gui.CreatePopupPanel(types.CreatePopupPanelOpts{
				HasLoader:   true,
				InstanceKey: self.key,
				Title:       self.title,
				Prompt:      latestContent,
			})

			self.opened = true
			self.lastRenderedContent = latestContent
		}

		if latestContent == self.lastRenderedContent {
			self.mu.Unlock()
			return nil
		}

		self.lastRenderedContent = latestContent
		self.mu.Unlock()

		if self.isCurrentPopup() {
			self.gui.SetViewContent(self.gui.Views().Confirmation, style.AttrBold.Sprint(latestContent))
		}

		return nil
	})
}

func (self *LiveCommandOutputPopup) Close() {
	self.mu.Lock()
	opened := self.opened
	self.mu.Unlock()

	if !opened {
		return
	}

	self.gui.OnUIThread(func() error {
		if self.isCurrentPopup() {
			self.gui.Context().Pop()
		}

		return nil
	})
}

func (self *LiveCommandOutputPopup) isCurrentPopup() bool {
	currentPopupOpts := self.gui.State().GetRepoState().GetCurrentPopupOpts()
	return currentPopupOpts != nil && currentPopupOpts.InstanceKey == self.key
}
