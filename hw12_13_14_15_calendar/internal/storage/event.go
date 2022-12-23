package storage

type DBType string

const (
	DBTypeMem DBType = "memory"
	DBTypeSQL DBType = "sql"
)

type Event struct {
	ID    string
	Title string
	// TODO
}
