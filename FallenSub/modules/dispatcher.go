package modules

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var Dispatcher = newDispatcher()

// newDispatcher creates a new dispatcher and loads modules.
func newDispatcher() *ext.Dispatcher {
	dispatcher := ext.NewDispatcher(nil)
	LoadModules(dispatcher)

	return dispatcher
}
