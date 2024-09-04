package engine

import (
	"context"

	"github.com/sirupsen/logrus"
)

// PlanDelete takes a Package and creates a Plan to delete it from Kubernetes
// To perform the delete simply apply the plan :)
func (pkg *Package) PlanDelete(ctx context.Context) (*Plan, error) {

	plan := &Plan{
		Valid:   false,
		Applied: true, // Starts off as applied, if anything needs updating it will be set to false
	}

	// For every Podinate resource, get the objects and create a plan for each
	for _, resource := range pkg.Resources {
		var objectChanges *[]ObjectChange
		objects, err := resource.GetObjects(ctx)
		if err != nil {
			return nil, err
		}

		for _, object := range objects {
			objectChange, err := GetDeleteChangeForObject(ctx, object)
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"object": object,
				"change": objectChange,
				"error":  err,
			}).Debug("Planned delete for object")
			if err != nil {
				return nil, err
			}

			if objectChanges == nil {
				objectChanges = new([]ObjectChange)
			}

			*objectChanges = append(*objectChanges, *objectChange)

		}

		changeType := ChangeTypeNoop
		for _, objectChange := range *objectChanges {
			if objectChange.ChangeType != ChangeTypeNoop {
				changeType = objectChange.ChangeType
			}
		}

		change := Change{
			ResourceType: resource.GetType(),
			ResourceID:   resource.GetName(),
			ChangeType:   changeType,
			Changes:      objectChanges,
		}
		plan.Changes = append(plan.Changes, change)
	}

	// If no errors so far, this plan must be valid
	plan.Valid = true

	// If any changes are not noops, the plan is not applied
	for _, change := range plan.Changes {
		if change.ChangeType != ChangeTypeNoop {
			plan.Applied = false
			break
		}
	}

	return plan, nil

}
