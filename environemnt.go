package main

type Environment struct {
	name        string
	shortName   string
	description string
}

// implement list.Item interface from bubbletea
func (e Environment) FilterValue() string {
	return e.shortName
}

func (e Environment) Name() string {
	return e.name
}

func (e Environment) ShortName() string {
	return e.shortName
}

func (e Environment) Description() string {
	return e.description
}
