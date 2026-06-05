package impl

import (
	"agenthub/backend/internal/dao"
	"agenthub/backend/internal/model"
	"agenthub/backend/internal/service"

	"github.com/google/uuid"
)

type ContactGroupService struct {
	dao dao.ContactGroupDao
}

func NewContactGroupService(contactGroupDao dao.ContactGroupDao) *ContactGroupService {
	return &ContactGroupService{dao: contactGroupDao}
}

func (svc *ContactGroupService) ListGroups() (*service.ListGroupsResponse, error) {
	groups, err := svc.dao.ListGroups()
	if err != nil {
		return nil, err
	}

	groupedSet := make(map[string]bool)
	result := make([]service.GroupWithItems, 0, len(groups))
	for _, group := range groups {
		items, err := svc.dao.ListItemsByGroupID(group.GroupID)
		if err != nil {
			return nil, err
		}

		groupItems := make([]service.GroupItem, 0, len(items))
		for _, item := range items {
			groupItems = append(groupItems, service.GroupItem{
				TaskID:    item.TaskID,
				SortOrder: item.SortOrder,
			})
			groupedSet[item.TaskID] = true
		}

		result = append(result, service.GroupWithItems{
			GroupID:   group.GroupID,
			Name:      group.Name,
			SortOrder: group.SortOrder,
			Items:     groupItems,
		})
	}

	taskIDs, err := svc.dao.ListActiveTaskIDs()
	if err != nil {
		return nil, err
	}
	ungrouped := make([]string, 0)
	for _, taskID := range taskIDs {
		if !groupedSet[taskID] {
			ungrouped = append(ungrouped, taskID)
		}
	}

	return &service.ListGroupsResponse{
		Groups:           result,
		UngroupedTaskIDs: ungrouped,
	}, nil
}

func (svc *ContactGroupService) CreateGroup(name string) (*model.ContactGroup, error) {
	return svc.dao.CreateGroup(model.ContactGroup{
		GroupID: uuid.New().String()[:8],
		Name:    name,
	})
}

func (svc *ContactGroupService) UpdateGroup(groupID, name string) error {
	updated, err := svc.dao.UpdateGroupName(groupID, name)
	if err != nil {
		return err
	}
	if !updated {
		return service.ErrNotFound("group not found")
	}
	return nil
}

func (svc *ContactGroupService) DeleteGroup(groupID string) error {
	deleted, err := svc.dao.DeleteGroupWithItems(groupID)
	if err != nil {
		return err
	}
	if !deleted {
		return service.ErrNotFound("group not found")
	}
	return nil
}

func (svc *ContactGroupService) AddItem(groupID, taskID string) (*model.ContactGroupItem, error) {
	return svc.dao.CreateItem(model.ContactGroupItem{
		GroupID: groupID,
		TaskID:  taskID,
	})
}

func (svc *ContactGroupService) RemoveItem(groupID, taskID string) error {
	deleted, err := svc.dao.DeleteItem(groupID, taskID)
	if err != nil {
		return err
	}
	if !deleted {
		return service.ErrNotFound("item not found")
	}
	return nil
}
