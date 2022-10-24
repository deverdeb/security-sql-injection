package tasks

type Priority string

const (
    Urgent Priority = "Urgent"
    Normal Priority = "Normal"
    Basse  Priority = "Basse"
)

type Status string

const (
    EnAttente  Status = "EnAttente"
    AFaire     Status = "AFaire"
    EnCours    Status = "EnCours"
    Terminee   Status = "Terminee"
    Abandonnee Status = "Abandonnee"
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
