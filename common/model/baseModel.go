package model

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseModel struct {
	BaseTableName string          `gorm:"-"`
	DB            *gorm.DB        `gorm:"-"`
	Ctx           context.Context `gorm:"-"`
}

func (m *BaseModel) TableName() string {
	return m.BaseTableName
}

// 批量插入
func (m *BaseModel) BatchCreate(ctx context.Context, db *gorm.DB, rows []map[string]interface{}) (int64, error) {
	res := db.Table(m.TableName()).WithContext(ctx).Create(rows)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

// 更新冻结申请订单
func (l *BaseModel) UpdateData(ctx context.Context, db *gorm.DB, queryCond map[string]interface{}, updateData map[string]interface{}) (affectRowNum int64, err error) {
	res := db.Table(l.TableName()).WithContext(ctx).Where(queryCond).Updates(updateData)

	if res.Error != nil {
		return affectRowNum, res.Error
	}

	return res.RowsAffected, nil
}

// 批量更新数据
func (m *BaseModel) BatchUpdateData(ctx context.Context, db *gorm.DB, pk string, updateFields []string, updateData interface{}) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: pk}},
		DoUpdates: clause.AssignmentColumns(updateFields),
	}).Table(m.TableName()).WithContext(ctx).CreateInBatches(updateData, 100).Error
}

// 查询数据列表
func (m *BaseModel) QueryDataList(ctx context.Context, db *gorm.DB, queryFields map[string]interface{}, simpleSql string, orderBy string) (list []map[string]interface{}, err error) {
	query := db.Table(m.TableName()).WithContext(ctx).Where(queryFields)
	if len(orderBy) > 0 {
		query = query.Order(orderBy)
	}
	if len(simpleSql) > 0 {
		query = query.Where(simpleSql)
	}
	err = query.Find(&list).Error

	if err != nil {
		return nil, err
	}

	return list, nil
}

// 查询单行数据
func (m *BaseModel) FindData(ctx context.Context, db *gorm.DB, queryFields map[string]interface{}, simpleSql string, orderBy string, selectField string) (result map[string]interface{}, err error) {
	query := db.Table(m.TableName()).WithContext(ctx)
	if len(queryFields) > 0 {
		query = query.Where(queryFields)
	}
	if len(orderBy) > 0 {
		query = query.Order(orderBy)
	}
	if len(simpleSql) > 0 {
		query = query.Where(simpleSql)
	}
	if len(selectField) > 0 {
		query = query.Select(selectField)
	}

	result = map[string]interface{}{}
	err = query.Take(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

// 数量统计
func (m *BaseModel) Count(ctx context.Context, db *gorm.DB, queryFields map[string]interface{}, simpleSql string, groupBySql string) (num int64, err error) {
	query := db.Table(m.TableName()).WithContext(ctx)
	if len(queryFields) > 0 {
		query = query.Where(queryFields)
	}
	if len(simpleSql) > 0 {
		query = query.Where(simpleSql)
	}
	if len(groupBySql) > 0 {
		query = query.Group(groupBySql)
	}

	err = query.Count(&num).Error
	if err != nil {
		return 0, err
	}
	return
}
