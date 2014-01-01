package dpl

import (
	"airdispat.ch/common"
	"html/template"
	"net/url"
)

type MessageType int

// Message Types
const MessagesPublic MessageType = 0
const MessagesPrivate MessageType = 1

// Special Addresses
var AllAddresses *common.ADAddress = nil

// Notifications are specially named actions
type Notification struct {
	Text  template.HTML
	Link  url.URL
	Image url.URL
}

type Host interface {
	GetMessagesForTag(plugin *PluginInstance, tag Tag) ([]Message, error)
	GetMessagesForTagWithPredicate(plugin *PluginInstance, tag Tag, p Predicate) ([]Message, error)
	GetMessagesForTagWithUser(plugin *PluginInstance, tag Tag, user *common.ADAddress) ([]Message, error)
	GetURLForAction(plugin *PluginInstance, action Action) (string, error)
	RunNotification(plugin *PluginInstance, n *Notification)
	Identify() string
}

type Predicate string
