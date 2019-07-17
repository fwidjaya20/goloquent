package goloquent

import "fmt"

// Reference is a struct that is used to store information about schema relations
type Reference struct {
	key          string
	relatedTable string
	relatedKey   string
	onUpdate     ReferenceAction
	onDelete     ReferenceAction
}

// Reference is a setter for related key
func (r *Reference) Reference(key string) *Reference {
	r.relatedKey = key

	return r
}

// On is a setter for related table
func (r *Reference) On(table string) *Reference {
	r.relatedTable = table

	return r
}

// OnUpdate is a setter for on update action
func (r *Reference) OnUpdate(action ReferenceAction) *Reference {
	r.onUpdate = action

	return r
}

// OnDelete is a setter for on delete action
func (r *Reference) OnDelete(action ReferenceAction) *Reference {
	r.onDelete = action

	return r
}

// Verbose is a function for print schema reference detail
func (r *Reference) Verbose() {
	fmt.Printf("    Key           : %v\n", r.key)
	fmt.Printf("    Related Table : %v\n", r.relatedTable)
	fmt.Printf("    Related Key   : %v\n", r.relatedKey)
	fmt.Printf("    OnUpdate      : %v\n", r.onUpdate)
	fmt.Printf("    OnDelete      : %v\n", r.onDelete)
}
