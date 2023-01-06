package storage

type DBType string

const (
	DBTypeMem DBType = "memory"
	DBTypeSQL DBType = "sql"
)

type Event struct {
	UUID  string `json:"uuid"`
	Title string `json:"title"`
	User  string `json:"user"`
	Date  string `json:"date"`
	Begin string `json:"start"`
	End   string `json:"end"`
}
