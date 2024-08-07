package service

import (
	"sync"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/sets/treeset"
)

type dispatcher interface {
	patch(message RawMessage)
}

type notifyManager interface {
	DelNotify(Sender)
	AddNotify(Sender)
}

type subscribe struct {
	topic       map[string]notifier
	subscribers treeset.Set
	mutex       sync.Mutex
}

func (s *subscribe) patch(msg RawMessage) {

}

func (s *subscribe) DelNotify(sender Sender) {

}
func (s *subscribe) AddNotify(sender Sender) {
	sub := s.subscribers
	if sub.Contains(sender) {

	}
}

/*manage a topic distribute */
type notifier struct {
	subscribe arraylist.List
}

func (n *notifier) notify(message RawMessage) {
	subs := n.subscribe
	subs.Each(func(index int, value interface{}) {
		sender := value.(Sender)
		if !sender.IsDelete() {
			sender.send(message.Content)
		} else {
			subs.Remove(index)
		}
	})
}

type Sender interface {
	send(*ArticleView) error
	IsDelete() bool
	setDelete(bool)
}

type TelegramSender struct {
	Owner    uint `json:"owner"`
	isDelete bool `json:"-"`
}

func (tg *TelegramSender) send(*ArticleView) error {
	return nil
}

func (tg *TelegramSender) IsDelete() bool { return tg.isDelete }

func (tg *TelegramSender) setDelete(flag bool) {
	tg.isDelete = flag
}

type RawMessage struct {
	From    string       `json:"from"`
	Content *ArticleView `json:"content"`
}
