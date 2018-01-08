package model

import "fmt"

type DomainRecord struct {
	ID       string `json:"id,omitempty",yaml:"id,omitempty"`
	Type     string `json:"type,omitempty",yaml:"type,omitempty"`
	Name     string `json:"name,omitempty",yaml:"name,omitempty"`
	Data     string `json:"data,omitempty",yaml:"data,omitempty"`
	Priority int    `json:"priority,omitempty",yaml:"priority,omitempty"`
	Port     int    `json:"port,omitempty",yaml:"port,omitempty"`
	TTL      int    `json:"ttl,omitempty",yaml:"ttl,omitempty"`
	Weight   int    `json:"weight,omitempty",yaml:"weight,omitempty"`
	Flags    int    `json:"flags,omitempty",yaml:"flags,omitempty"`
	Tag      string `json:"tag,omitempty",yaml:"tag,omitempty"`
}

func (d *DomainRecord) isTypeNameEqual(to DomainRecord) bool {
	return d.Type == to.Type && d.Name == to.Name
}

func (d *DomainRecord) isEqual(to DomainRecord) bool {
	fmt.Println(d.Name, to.Name, "...", d.Priority == to.Priority, d.Port == to.Port, d.TTL, "==", to.TTL, d.Weight == to.Weight, d.Flags == to.Flags, d.Tag == to.Tag)
	return d.Type == to.Type && d.Name == to.Name && d.Data == to.Data && d.Priority == to.Priority && d.Port == to.Port && d.TTL == to.TTL && d.Weight == to.Weight && d.Flags == to.Flags && d.Tag == to.Tag
}