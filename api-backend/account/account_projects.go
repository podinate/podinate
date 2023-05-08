package account

import (
	"github.com/johncave/podinate/api-backend/config"
	"github.com/johncave/podinate/api-backend/project"
)

// GetProjects returns the projects of the account
func (a *Account) GetProjects(page int32, limit int32) ([]project.Project, error) {
	rows, err := config.DB.Query("SELECT uuid, id, name FROM project WHERE account_uuid = $1 OFFSET $2 LIMIT $3", a.Uuid, page, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Read all the projects for the account
	projects := make([]project.Project, 0)
	for rows.Next() {
		var project project.Project
		err = rows.Scan(&project.Uuid, &project.ID, &project.Name)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil

}
