package dpl

import (
	"encoding/xml"
	"errors"
	"io"
)

type Tag struct {
	Name      string  `xml:"name,attr"`
	Type      string  `xml:"type,attr"`
	OnReceive string  `xml:"onreceive,attr"`
	Fields    []Field `xml:"field"`
}

type Field struct {
	Name           string `xml:",innerxml"`
	Optional       bool   `xml:"optional,attr"`
	Representation string `xml:"repr,attr"` // Allows the plugin to access the value in other ways
}

type Action struct {
	HTML    string `xml:",innerxml"`
	Name    string `xml:"name,attr"`
	Default bool   `xml:"default,attr"`
}

func ParseDPLStream(r io.Reader) (*Plugin, error) {
	coder := xml.NewDecoder(r)

	var d Plugin
	err := coder.Decode(&d)
	if err != nil {
		return nil, err
	}

	return &d, verifyPlugin(&d)
}

func ParseDPL(r []byte) (*Plugin, error) {
	var d Plugin
	err := xml.Unmarshal(r, &d)
	if err != nil {
		return nil, err
	}

	return &d, verifyPlugin(&d)
}

func verifyPlugin(p *Plugin) error {
	if p.Path == "" {
		return errors.New("Plugin lacks a path.")
	}
	if p.Name == "" {
		p.Name = p.Path
	}

	if p.Tag == nil {
		return errors.New("Plugin lacks tags.")
	}
	for _, v := range p.Tag {
		if p.Tags == nil {
			p.Tags = make(map[string]Tag)
		}
		p.Tags[v.Name] = v
	}

	if p.Action == nil {
		return errors.New("Plugin lacks actions.")
	}
	for _, v := range p.Action {
		if p.Actions == nil {
			p.Actions = make(map[string]Action)
		}
		if v.Default {
			p.DefaultAction = v.Name
		}
		p.Actions[v.Name] = v
	}

	return nil
}
