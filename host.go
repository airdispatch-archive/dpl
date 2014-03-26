package dpl

import (
	"airdispat.ch/identity"
	"html/template"
	"net/url"
)

type MessageType int

// Message Types
const MessagesPublic MessageType = 0
const MessagesPrivate MessageType = 1

// Special Addresses
var AllAddresses *identity.Address = nil

// Notifications are specially named actions
type Notification struct {
	Text  template.HTML
	Link  url.URL
	Image url.URL
}

// Host Interface describes what methods a Plugin Needs to Communicate with the Mailserver
type Host interface {
	// Data Fetching Methods
	GetMessages(plugin *PluginInstance, tag Tag, predicate *Predicate, limit int) ([]Message, error)     // Get a list of Messages that are Tagged with a specified tag, and follow a predicate
	GetURLForAction(plugin *PluginInstance, action Action, message Message, user User) (*url.URL, error) // Gets the URL for an Action
	SendURL(plugin *PluginInstance) *url.URL

	// Action Methods
	RunNotification(plugin *PluginInstance, n *Notification)

	// Meta Methods
	Identify() string
}

type Predicate struct {
	Sort   Sorter
	Search Searcher
}

type Sorter struct {
	FieldName string
	Ascending bool
}

type Searcher struct {
	Value string
	Data  interface{}
}
