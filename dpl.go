package dpl

import (
	"encoding/xml"
	"errors"
	"fmt"
)

type Plugin struct {
	Name    string            `xml:"name"`
	Path    string            `xml:"path"`
	Detect  Detector          `xml:"detect>field"`
	Action  []Action          `xml:"action"`
	Actions map[string]Action `xml:"-"`
}

func (p *Plugin) String() string {
	return fmt.Sprintf("%s (%s) %s \n %s", p.Name, p.Path, p.Detect, p.Action)
}

type Detector []Field

func (d Detector) String() string {
	output := ""
	for _, v := range d {
		output += "\n\t " + v.Name

		if v.Required {
			output += "\t(R)"
		}
	}
	return output
}

type Field struct {
	Name     string `xml:"name,attr"`
	Required bool   `xml:"required,attr"`
}

type Action struct {
	HTML string `xml:",innerxml"`
	Name string `xml:"name,attr"`
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
	if p.Detect == nil {
		return errors.New("Plugin lacks detectors.")
	}

	for _, v := range p.Action {
		if p.Actions == nil {
			p.Actions = make(map[string]Action)
		}
		p.Actions[v.Name] = v
	}
	return nil
}
