package dpl

import (
	"fmt"
	"net/url"
	"time"
)

func (p *PluginInstance) assembleContext(action Action) (PluginContext, error) {
	context := PluginContext{
		Tags:   p.createTagsLinker(),
		Action: p.ActionLambda(),
		Host:   fmt.Sprintf("DPL v0.1 (Host: %s)", p.Host.Identify()),
	}
	return context, nil
}

func (p *PluginInstance) ActionLambda() func(string) string {
	return func(a string) string {
		v := p.Actions[a]
		url, _ := p.Host.GetURLForAction(p, v)
		return url
	}
}

// Plugins have access to three top-level variables
// Tags - lists of messages that have been marked for the plugin
// Actions - list of actions that have been defined by the plugin, access to URLs and such
// Host - a string containing information about the host
type PluginContext struct {
	Tags   TagLinker
	Action func(string) string
	Host   string
}

// TagLinker manages message lists - this is a map indexed by the name of the tag
// The tag variables allow the plugin to join with other tags _or_ sort the tags
type TagLinker map[string]TagContext
type TagContext []*MessageContext

func (p *PluginInstance) createTagsLinker() TagLinker {
	linker := make(TagLinker)
	for _, v := range p.Plugin.Tags {
		messageList, err := p.Host.GetMessagesForTag(p, v)
		if err != nil {
			continue
		}

		mContext := make([]*MessageContext, 0)
		for _, v := range messageList {
			mContext = append(mContext, &MessageContext{v})
		}
		linker[v.Name] = mContext
	}
	return linker
}

type User struct {
	Name   string
	Avatar url.URL
}

type Message interface {
	Get(field string) ([]byte, error)
	Has(field string) bool
	Created() time.Time
	Sender() User
}

type MessageContext struct {
	inner Message
}

// The Get Function allows DPL templates to access different fields off of the current MessageContext
func (m *MessageContext) Get() func(string) string {
	return func(field string) string {
		val, err := m.inner.Get(field)
		if err != nil {
			return "Plugin Error (" + err.Error() + ") when accessing field " + field
		}
		return string(val)
	}
}

// The Date Function allows DPL templates to access the sent date of the current MessageContext
func (m *MessageContext) Created() func(string) string {
	return func(format string) string {
		val := m.inner.Created()
		return val.Format(format)
	}
}

// ActionLinker manages lists of actions - this is a map indexed by the name of the action
type ActionLinker map[string]ActionContext
type ActionContext struct {
	URL string
}

// Action Linkers _can_ be cached, actions only change on PluginUpdate
func (p *PluginInstance) createActionLinker() ActionLinker {
	actions := make(ActionLinker)
	for _, v := range p.Plugin.Action {
		url, _ := p.Host.GetURLForAction(p, v)
		actions[v.Name] = ActionContext{url}
	}
	return actions
}
