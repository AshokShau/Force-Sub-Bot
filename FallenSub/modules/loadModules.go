package modules

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

// LoadModules loads all the modules
func LoadModules(d *ext.Dispatcher) {
	d.AddHandler(handlers.NewCommand("start", start))
	d.AddHandler(handlers.NewCommand("ping", ping))
	d.AddHandler(handlers.NewCommand("fsub", setFSub))

	d.AddHandlerToGroup(handlers.NewMessage(message.All, fSubWatcher), 0)
	d.AddHandler(handlers.NewCallback(callbackquery.Prefix("unmuteMe_"), unMuteMe))
}
