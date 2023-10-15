package telegram

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"tg-bot/clients/telegram"
	"tg-bot/lib/e"
	"tg-bot/storage"
	"tg-bot/storage/files"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	sendMessage := NewMessageSender(chatID, p.tg)
	slog.Info(fmt.Sprintf("got new command '%s' from '%s'", text, username))
	// add page: http://..
	// rnd page: /rnd
	// help: /help
	//start: /star: hi + help

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return sendMessage(msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't execute save page command", err) }()

	sendMessage := NewMessageSender(chatID, p.tg)
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}

	if isExists {
		return sendMessage(msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := sendMessage(msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't execute send random command", err) }()
	sendMessage := NewMessageSender(chatID, p.tg)

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, files.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, files.ErrNoSavedPages) {
		return sendMessage(msgNoSavedPages)
	}

	if err := sendMessage(page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	sendMessage := NewMessageSender(chatID, p.tg)
	if err := sendMessage(msgHelp); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendHello(chatID int) error {
	sendMessage := NewMessageSender(chatID, p.tg)
	if err := sendMessage(msgHello); err != nil {
		return err
	}
	return nil
}

func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
