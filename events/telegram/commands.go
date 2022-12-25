package telegram

import (
	"errors"
	"link_saver/clients/telegram"
	"link_saver/lib/e"
	"link_saver/storage"
	"link_saver/storage/files"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd = "/rnd"
	HelpCmd = "/help"
	StartCmd = "/start"
)

func (p *Processor) doComd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	
	log.Printf("got new command '%s' from %s",text,username)

	if isAddCmd(text){
		return p.savePage(chatID,text,username)
	}
	switch text {
	case RndCmd:
		return p.SendRandom(chatID,username)
	case HelpCmd:
		return p.SendHelp(chatID)
	case StartCmd:
		return p.SendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, text string, username string) (err error) {
	defer func() {err =e.Wrap("cant do command save page ",err )}()

	sendMsg := NewMessageSender(chatID,p.tg)

	page:= &storage.Page{
		URL:	pageURL
		UserName:	username,
	}

	isExists,err := p.storage.IsExists(page)
	if err!= nil {
		return err
	}
	if isExists{
		return sendMsg(msgAlredyExists)
	}
	if err:=p.storage.Save(page); err!= nil {
		return err
	}

	if err :=p.tg.SendMessage(chatID,msgSaved);err != nil{
		return err 
	}
	return nil
}


func (p* Processor) SendRandom(ChatID int, username string) (err error) { 
	defer func() {err=e.Wrap("cant do command: cant send random",err)}()

	page,err:=p.storage.PickRandom(username)
	if err != nil && !errors.Is(err,storage.ErrorNoSavedPages) {
		return err
	}
	if errors.Is(err,storage.ErrorNoSavedPages){
		return p.tg.SendMessage(ChatID, msgNoSavedPages)
	}

	if err:=p.tg.SendMessage(ChatID,page.URL); err!=nil {
		return err
	}
	return p.storage.Remove(page)
}

func (p *Processor) SendHelp(chatID int) error {
	return p.tg.SendMessage(chatID,msgHelp)
}

func (p *Processor) SendHello(chatID int) error {
	return p.tg.SendMessage(chatID,msgHello)
}

func NewMessageSender(chatID int, tg * telegram.Client) func(string) error {
	return func(msg string) error { 
		return tg.SendMessage(chatID,msg)

	}
}

func isAddCommand(text string) bool  { 
	return isUrl(text)
}

func isUrl(text string) bool{
	u,err:= url.Parse(text)
	return err ==nil && u.Host != ""
}