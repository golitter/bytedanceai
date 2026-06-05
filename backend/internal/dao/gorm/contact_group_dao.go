package gormdao

import (
	"errors"

	"agenthub/backend/internal/model"
	"agenthub/backend/pkg/db"

	"gorm.io/gorm"
)

type ContactGroupDao struct{}

func NewContactGroupDao() *ContactGroupDao {
	return &ContactGroupDao{}
}

func (dao *ContactGroupDao) ListGroups() ([]model.ContactGroup, error) {
	var groups []model.ContactGroup
	if err := db.GetDB().Order("sort_order ASC, created_at ASC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (dao *ContactGroupDao) ListItemsByGroupID(groupID string) ([]model.ContactGroupItem, error) {
	var items []model.ContactGroupItem
	if err := db.GetDB().Where("group_id = ?", groupID).Order("sort_order ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (dao *ContactGroupDao) ListActiveTaskIDs() ([]string, error) {
	var taskIDs []string
	if err := db.GetDB().Model(&model.Task{}).Where("status = ?", "active").Pluck("task_id", &taskIDs).Error; err != nil {
		return nil, err
	}
	return taskIDs, nil
}

func (dao *ContactGroupDao) CreateGroup(group model.ContactGroup) (*model.ContactGroup, error) {
	if err := db.GetDB().Create(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (dao *ContactGroupDao) UpdateGroupName(groupID, name string) (bool, error) {
	result := db.GetDB().Model(&model.ContactGroup{}).Where("group_id = ?", groupID).Update("name", name)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (dao *ContactGroupDao) DeleteGroupWithItems(groupID string) (bool, error) {
	found := true
	err := db.GetDB().Transaction(func(tx *gorm.DB) error {
		tx.Where("group_id = ?", groupID).Delete(&model.ContactGroupItem{})
		result := tx.Where("group_id = ?", groupID).Delete(&model.ContactGroup{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			found = false
			return nil
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return found, nil
}

func (dao *ContactGroupDao) CreateItem(item model.ContactGroupItem) (*model.ContactGroupItem, error) {
	if err := db.GetDB().Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (dao *ContactGroupDao) DeleteItem(groupID, taskID string) (bool, error) {
	result := db.GetDB().Where("group_id = ? AND task_id = ?", groupID, taskID).Delete(&model.ContactGroupItem{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

var _ = errors.Is
