package storage

type DBType string

const (
	DBTypeMem DBType = "memory"
	DBTypeSQL DBType = "sql"
)

type Event struct {
	UUID  string
	Title string
	// TODO
}
