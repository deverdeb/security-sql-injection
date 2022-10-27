package tasks

type Priority string

const (
	Highest Priority = "Highest"
	High    Priority = "High"
	Medium  Priority = "Medium"
	Low     Priority = "Low"
	Lowest  Priority = "Lowest"
)

type Status string

const (
	Draft      Status = "Draft"
	Open       Status = "Open"
	InProgress Status = "InProgress"
	Done       Status = "Done"
	Closed     Status = "Closed"
	Abandoned  Status = "Abandoned"
)

type Task struct {
	Id          int64
	UserId      int64
	Name        string
	Description string
	Priority    Priority
	Status      Status
	Archived    bool
}
