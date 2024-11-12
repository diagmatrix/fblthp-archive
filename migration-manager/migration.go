package mm

type Migration struct {
	ID       int    // Migration version id
	Name     string // Migration description
	Executed bool   // Migration status
}
