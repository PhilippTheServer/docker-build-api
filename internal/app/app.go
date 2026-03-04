package app // Declares this file belongs to package "app" (an importable library)

// type introduces a new named type in Go.
// defines a struct named Config to hold configuration values.
type Config struct {
	Addr string // Addr is a field of type string, e.g. ":8080" or "localhost:8080"
}
