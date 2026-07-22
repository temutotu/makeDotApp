package input

type SelectField struct {
	ID            string
	Name          string
	Label         string
	Options       []SelectOption
	SelectedValue string
}

type SelectOption struct {
	Value string
	Label string
}
