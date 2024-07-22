package dispatcher

import (
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/modules"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var Dispatcher = NewDispatcher()

// NewDispatcher creates a new dispatcher and loads modules.
func NewDispatcher() *ext.Dispatcher {
	dispatcher := ext.NewDispatcher(nil)
	modules.LoadModules(dispatcher)

	return dispatcher
}
