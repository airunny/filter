package assignment

// set value to data
type Setter interface {
	AssignmentSet(key string, value interface{})
}

// merge value to data
type Merger interface {
	AssignmentMerge(key string, value interface{})
}

// delete key from data
type Deleter interface {
	AssignmentDelete(key string, value interface{})
}

// increase value to data key failed
type Increaser interface {
	AssignmentIncrease(key string, value float64)
}
