package dpl

import (
	"html/template"
	"net/url"
	"time"
)

type PluginContext struct {
	Host    string
	Version int
	User    *User
	Message *MessageContext
}

func (p *PluginInstance) createFuncMap() template.FuncMap {
	return template.FuncMap{
		"action":  p.ActionLambda(),
		"actionc": p.ActionLambdaContext(),
		"tag":     p.TagLambda(),
	}
}

// ActionLinker manages lists of actions - this is a map indexed by the name of the action
type ActionLinker func(string) *url.URL
type ActionLinkerContext func(string, interface{}) *url.URL

func (p *PluginInstance) ActionLambdaContext() ActionLinkerContext {
	return func(a string, b interface{}) *url.URL {
		v, ok := p.Actions[a]
		if !ok {
			return nil
		}

		var url *url.URL

		switch t := b.(type) {
		case *MessageContext:
			url, _ = p.Host.GetURLForAction(p, v, t.inner, nil)
		case User:
			url, _ = p.Host.GetURLForAction(p, v, nil, t)
		default:
			url, _ = p.Host.GetURLForAction(p, v, nil, nil)
		}

		return url
	}
}

func (p *PluginInstance) ActionLambda() ActionLinker {
	return func(a string) *url.URL {
		return p.ActionLambdaContext()(a, nil)
	}
}

// TagLinker manages message lists - this is a map indexed by the name of the tag
// The tag variables allow the plugin to join with other tags _or_ sort the tags
type TagLinker func(string, *Predicate, int) []*MessageContext

func (p *PluginInstance) TagLambda() TagLinker {
	return func(tagName string, search *Predicate, limit int) []*MessageContext {
		tag, ok := p.Tags[tagName]
		if !ok {
			panic("Couldn't find tag named " + tagName)
		}

		messages, err := p.Host.GetMessages(p, tag, search, limit)
		if err != nil {
			// Not quite sure what to do with this error
			// Maybe we shouldn't pass an error and just Send an emtpy list instead
			panic(err)
		}

		var contexts []*MessageContext
		for _, v := range messages {
			contexts = append(contexts, &MessageContext{
				v,
				v.Sender(),
			})
		}

		return contexts
	}
}

type User interface {
	// Full Name of the Usre
	Name() string

	// Address Information about the User
	DisplayAddress() string
	Address() string

	// Profile Image Location
	Avatar() *url.URL

	// Profile Location
	Profile() *url.URL
}

type Message interface {
	Get(field string) ([]byte, error)
	Has(field string) bool
	Created() time.Time
	Sender() User
}

type MessageContext struct {
	inner  Message
	Sender User
}

// Action on MessageContext retrieves the URL for an Action that has a Message Attached
func (m *MessageContext) Action(f string) url.URL {
	return url.URL{}
}

// The Get Function allows DPL templates to access different fields off of the current MessageContext
func (m *MessageContext) Get(field string) string {
	val, err := m.inner.Get(field)
	if err != nil {
		return "Plugin Error (" + err.Error() + ") when accessing field " + field
	}
	return string(val)
}

// The Date Function allows DPL templates to access the sent date of the current MessageContext
func (m *MessageContext) Created(format string) string {
	val := m.inner.Created()
	return val.Format(format)
}
