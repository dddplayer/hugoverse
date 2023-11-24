package entity

// TextTemplate is the representation of a parsed template. The *parse.Tree
// field is exported only for use by html/template and should be treated
// as unexported by all other clients.
type TextTemplate struct {
	name string
}
