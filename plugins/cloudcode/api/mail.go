package api

import (
	"github.com/pocketbase/pocketbase/tools/mailer"
	lua "github.com/yuin/gopher-lua"
	"net/mail"
)

func getMailModule(l *lua.LState) *lua.LTable {
	m := l.SetFuncs(l.NewTable(), mailExports)
	return m
}

var mailExports = map[string]lua.LGFunction{
	"send": sendMessage,
}

func sendMessage(l *lua.LState) int {
	// TODO: fully implement after https://github.com/pocketbase/pocketbase/issues/1727 is done.

	message := &mailer.Message{
		From: mail.Address{
			Address: getApp().Settings().Meta.SenderAddress,
			Name:    getApp().Settings().Meta.SenderName,
		},
		To:      mail.Address{Address: "test@example.com", Name: "Test Recipient"},
		Subject: "Test Subject",
		Text:    "This is the text of the message",
		HTML:    "<html><body><p><strong>This</strong> is the HTML text of the message</p></body></html>",
	}

	getApp().NewMailClient().Send(message)

	return 0
}
