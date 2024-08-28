package api

import (
	"encoding/json"
	"fmt"
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/modules"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

var (
	allowedTokens    = strings.Split(os.Getenv("BOT_TOKEN"), " ")
	lenAllowedTokens = len(allowedTokens)
)

const (
	statusCodeSuccess = 200
)

// Bot Handles all incoming traffic from webhooks.
func Bot(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	split := strings.Split(url, "/")
	if len(split) < 2 {
		fmt.Println(w, "url path too short")
		w.WriteHeader(statusCodeSuccess)

		return
	}

	botToken := split[len(split)-2]

	bot, _ := gotgbot.NewBot(botToken, &gotgbot.BotOpts{DisableTokenCheck: false})

	// Delete the webhook in case token is unauthorized.
	if lenAllowedTokens > 0 && allowedTokens[0] != "" && !findInStringSlice(allowedTokens, botToken) {
		_, _ = bot.DeleteWebhook(&gotgbot.DeleteWebhookOpts{DropPendingUpdates: true}) // It doesn't matter if it errors
		w.WriteHeader(statusCodeSuccess)
		return
	}

	var update gotgbot.Update

	body, err := io.ReadAll(r.Body)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error reading request body: %v", err)
		w.WriteHeader(statusCodeSuccess)
		return
	}

	err = json.Unmarshal(body, &update)
	if err != nil {
		fmt.Println("failed to unmarshal body ", err)
		w.WriteHeader(statusCodeSuccess)

		return
	}

	bot.Username = split[len(split)-1]

	err = modules.Dispatcher.ProcessUpdate(bot, &update, map[string]any{})
	if err != nil {
		fmt.Printf("error while processing update: %v", err)
	}

	w.WriteHeader(statusCodeSuccess)
}

func findInStringSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
