package ui

// Category represents an activity category
type Category struct {
	Name        string
	Description string
	Activities  []ActivityEntry
}

// ActivityEntry represents a single activity in the UI
type ActivityEntry struct {
	ID          string
	Key         string
	Description string
}
