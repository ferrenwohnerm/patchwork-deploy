package state

// Tag sets the tag field on a record (fluent helper).
func (r Record) WithTag(tag string) Record {
	r.Tag = tag
	return r
}

// HasTag reports whether the record carries a non-empty tag.
func (r Record) HasTag() bool {
	return r.Tag != ""
}
