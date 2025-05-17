package data

import (
	"go-stock/backend/db"
	"gorm.io/gorm"
)

// @Author spark
// @Date 2025/4/3 11:18
// @Desc
// -----------------------------------------------------------------------------------
type Group struct {
	gorm.Model
	Name string `json:"name" gorm:"index"`
	Sort int    `json:"sort"`
}

func (Group) TableName() string {
	return "stock_groups"
}

type GroupStock struct {
	gorm.Model
	StockCode string `json:"stockCode" gorm:"index"`
	GroupId   int    `json:"groupId" gorm:"index"`
	GroupInfo Group  `json:"groupInfo" gorm:"foreignKey:GroupId;references:ID"`
}

func (GroupStock) TableName() string {
	return "group_stock_info"
}

type StockGroupApi struct {
	dao *gorm.DB
}

func NewStockGroupApi(dao *gorm.DB) *StockGroupApi {
	return &StockGroupApi{dao: db.Dao}
}

func (receiver StockGroupApi) AddGroup(group Group) bool {
	err := receiver.dao.Where("name = ?", group.Name).FirstOrCreate(&group).Updates(&Group{
		Name: group.Name,
		Sort: group.Sort,
	}).Error
	return err == nil
}
func (receiver StockGroupApi) GetGroupList() []Group {
	var groups []Group
	receiver.dao.Find(&groups)
	return groups
}
func (receiver StockGroupApi) GetGroupStockByGroupId(groupId int) []GroupStock {
	var stockGroup []GroupStock
	receiver.dao.Preload("GroupInfo").Where("group_id = ?", groupId).Find(&stockGroup)
	return stockGroup
}

func (receiver StockGroupApi) AddStockGroup(groupId int, stockCode string) bool {
	err := receiver.dao.Where("group_id = ? and stock_code = ?", groupId, stockCode).FirstOrCreate(&GroupStock{
		GroupId:   groupId,
		StockCode: stockCode,
	}).Updates(&GroupStock{
		GroupId:   groupId,
		StockCode: stockCode,
	}).Error
	return err == nil
}

func (receiver StockGroupApi) RemoveStockGroup(code string, name string, id int) bool {
	err := receiver.dao.Where("group_id = ? and stock_code = ?", id, code).Delete(&GroupStock{}).Error
	return err == nil
}

func (receiver StockGroupApi) RemoveGroup(id int) bool {
	err := receiver.dao.Where("id = ?", id).Delete(&Group{}).Error
	err = receiver.dao.Where("group_id = ?", id).Delete(&GroupStock{}).Error
	return err == nil

}
