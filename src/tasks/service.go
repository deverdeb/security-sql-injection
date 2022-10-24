package tasks

import (
    "fmt"
    "web-appli/src/users"
)

var Service = TasksService{
    dao: &Dao,
}

type TasksService struct {
    dao *TasksDao
}

func (service TasksService) Save(task *Task) (*Task, error) {
    if task == nil {
        return nil, fmt.Errorf("we cannot save 'nil' task")
    } else if task.Id <= 0 {
        // Création
        return service.dao.Create(task)
    } else {
        // Mise à jour
        return service.dao.Update(task)
    }
}

func (service TasksService) Delete(task *Task) error {
    if task == nil {
        return fmt.Errorf("we cannot save 'nil' task")
    } else {
        // Suppression
        return service.dao.Delete(task)
    }
}

func (service TasksService) FindAll() ([]*Task, error) {
    return service.dao.FindAll()
}

func (service TasksService) FindByUser(user *users.User) ([]*Task, error) {
    return service.dao.FindByUserId(user.Id)
}

// FindByIdAndUser appel une méthode vulnérable - injection SQL sur le champ d'identifiant.
// Pour sécuriser un minimum, nous devrions passer l'identifiant sous forme de "int64".
func (service TasksService) FindByIdAndUser(id string, user *users.User) (*Task, error) {
    return service.dao.FindByIdAndUserId(id, user.Id)
}

func (service TasksService) SearchByText(user *users.User, search string) ([]*Task, error) {
    return service.dao.SearchByText(user.Id, search)
}
