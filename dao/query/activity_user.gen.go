// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"api.backend.xjco2913/dao/model"
)

func newActivityUser(db *gorm.DB, opts ...gen.DOOption) activityUser {
	_activityUser := activityUser{}

	_activityUser.activityUserDo.UseDB(db, opts...)
	_activityUser.activityUserDo.UseModel(&model.ActivityUser{})

	tableName := _activityUser.activityUserDo.TableName()
	_activityUser.ALL = field.NewAsterisk(tableName)
	_activityUser.ActivityID = field.NewString(tableName, "activityId")
	_activityUser.UserID = field.NewString(tableName, "userId")
	_activityUser.CreatedAt = field.NewTime(tableName, "createdAt")
	_activityUser.UpdatedAt = field.NewTime(tableName, "updatedAt")
	_activityUser.FinalFee = field.NewInt32(tableName, "finalFee")

	_activityUser.fillFieldMap()

	return _activityUser
}

type activityUser struct {
	activityUserDo activityUserDo

	ALL        field.Asterisk
	ActivityID field.String
	UserID     field.String
	CreatedAt  field.Time
	UpdatedAt  field.Time
	FinalFee   field.Int32

	fieldMap map[string]field.Expr
}

func (a activityUser) Table(newTableName string) *activityUser {
	a.activityUserDo.UseTable(newTableName)
	return a.updateTableName(newTableName)
}

func (a activityUser) As(alias string) *activityUser {
	a.activityUserDo.DO = *(a.activityUserDo.As(alias).(*gen.DO))
	return a.updateTableName(alias)
}

func (a *activityUser) updateTableName(table string) *activityUser {
	a.ALL = field.NewAsterisk(table)
	a.ActivityID = field.NewString(table, "activityId")
	a.UserID = field.NewString(table, "userId")
	a.CreatedAt = field.NewTime(table, "createdAt")
	a.UpdatedAt = field.NewTime(table, "updatedAt")
	a.FinalFee = field.NewInt32(table, "finalFee")

	a.fillFieldMap()

	return a
}

func (a *activityUser) WithContext(ctx context.Context) *activityUserDo {
	return a.activityUserDo.WithContext(ctx)
}

func (a activityUser) TableName() string { return a.activityUserDo.TableName() }

func (a activityUser) Alias() string { return a.activityUserDo.Alias() }

func (a activityUser) Columns(cols ...field.Expr) gen.Columns {
	return a.activityUserDo.Columns(cols...)
}

func (a *activityUser) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (a *activityUser) fillFieldMap() {
	a.fieldMap = make(map[string]field.Expr, 5)
	a.fieldMap["activityId"] = a.ActivityID
	a.fieldMap["userId"] = a.UserID
	a.fieldMap["createdAt"] = a.CreatedAt
	a.fieldMap["updatedAt"] = a.UpdatedAt
	a.fieldMap["finalFee"] = a.FinalFee
}

func (a activityUser) clone(db *gorm.DB) activityUser {
	a.activityUserDo.ReplaceConnPool(db.Statement.ConnPool)
	return a
}

func (a activityUser) replaceDB(db *gorm.DB) activityUser {
	a.activityUserDo.ReplaceDB(db)
	return a
}

type activityUserDo struct{ gen.DO }

func (a activityUserDo) Debug() *activityUserDo {
	return a.withDO(a.DO.Debug())
}

func (a activityUserDo) WithContext(ctx context.Context) *activityUserDo {
	return a.withDO(a.DO.WithContext(ctx))
}

func (a activityUserDo) ReadDB() *activityUserDo {
	return a.Clauses(dbresolver.Read)
}

func (a activityUserDo) WriteDB() *activityUserDo {
	return a.Clauses(dbresolver.Write)
}

func (a activityUserDo) Session(config *gorm.Session) *activityUserDo {
	return a.withDO(a.DO.Session(config))
}

func (a activityUserDo) Clauses(conds ...clause.Expression) *activityUserDo {
	return a.withDO(a.DO.Clauses(conds...))
}

func (a activityUserDo) Returning(value interface{}, columns ...string) *activityUserDo {
	return a.withDO(a.DO.Returning(value, columns...))
}

func (a activityUserDo) Not(conds ...gen.Condition) *activityUserDo {
	return a.withDO(a.DO.Not(conds...))
}

func (a activityUserDo) Or(conds ...gen.Condition) *activityUserDo {
	return a.withDO(a.DO.Or(conds...))
}

func (a activityUserDo) Select(conds ...field.Expr) *activityUserDo {
	return a.withDO(a.DO.Select(conds...))
}

func (a activityUserDo) Where(conds ...gen.Condition) *activityUserDo {
	return a.withDO(a.DO.Where(conds...))
}

func (a activityUserDo) Order(conds ...field.Expr) *activityUserDo {
	return a.withDO(a.DO.Order(conds...))
}

func (a activityUserDo) Distinct(cols ...field.Expr) *activityUserDo {
	return a.withDO(a.DO.Distinct(cols...))
}

func (a activityUserDo) Omit(cols ...field.Expr) *activityUserDo {
	return a.withDO(a.DO.Omit(cols...))
}

func (a activityUserDo) Join(table schema.Tabler, on ...field.Expr) *activityUserDo {
	return a.withDO(a.DO.Join(table, on...))
}

func (a activityUserDo) LeftJoin(table schema.Tabler, on ...field.Expr) *activityUserDo {
	return a.withDO(a.DO.LeftJoin(table, on...))
}

func (a activityUserDo) RightJoin(table schema.Tabler, on ...field.Expr) *activityUserDo {
	return a.withDO(a.DO.RightJoin(table, on...))
}

func (a activityUserDo) Group(cols ...field.Expr) *activityUserDo {
	return a.withDO(a.DO.Group(cols...))
}

func (a activityUserDo) Having(conds ...gen.Condition) *activityUserDo {
	return a.withDO(a.DO.Having(conds...))
}

func (a activityUserDo) Limit(limit int) *activityUserDo {
	return a.withDO(a.DO.Limit(limit))
}

func (a activityUserDo) Offset(offset int) *activityUserDo {
	return a.withDO(a.DO.Offset(offset))
}

func (a activityUserDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *activityUserDo {
	return a.withDO(a.DO.Scopes(funcs...))
}

func (a activityUserDo) Unscoped() *activityUserDo {
	return a.withDO(a.DO.Unscoped())
}

func (a activityUserDo) Create(values ...*model.ActivityUser) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Create(values)
}

func (a activityUserDo) CreateInBatches(values []*model.ActivityUser, batchSize int) error {
	return a.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (a activityUserDo) Save(values ...*model.ActivityUser) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Save(values)
}

func (a activityUserDo) First() (*model.ActivityUser, error) {
	if result, err := a.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.ActivityUser), nil
	}
}

func (a activityUserDo) Take() (*model.ActivityUser, error) {
	if result, err := a.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.ActivityUser), nil
	}
}

func (a activityUserDo) Last() (*model.ActivityUser, error) {
	if result, err := a.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.ActivityUser), nil
	}
}

func (a activityUserDo) Find() ([]*model.ActivityUser, error) {
	result, err := a.DO.Find()
	return result.([]*model.ActivityUser), err
}

func (a activityUserDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.ActivityUser, err error) {
	buf := make([]*model.ActivityUser, 0, batchSize)
	err = a.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (a activityUserDo) FindInBatches(result *[]*model.ActivityUser, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return a.DO.FindInBatches(result, batchSize, fc)
}

func (a activityUserDo) Attrs(attrs ...field.AssignExpr) *activityUserDo {
	return a.withDO(a.DO.Attrs(attrs...))
}

func (a activityUserDo) Assign(attrs ...field.AssignExpr) *activityUserDo {
	return a.withDO(a.DO.Assign(attrs...))
}

func (a activityUserDo) Joins(fields ...field.RelationField) *activityUserDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Joins(_f))
	}
	return &a
}

func (a activityUserDo) Preload(fields ...field.RelationField) *activityUserDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Preload(_f))
	}
	return &a
}

func (a activityUserDo) FirstOrInit() (*model.ActivityUser, error) {
	if result, err := a.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.ActivityUser), nil
	}
}

func (a activityUserDo) FirstOrCreate() (*model.ActivityUser, error) {
	if result, err := a.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.ActivityUser), nil
	}
}

func (a activityUserDo) FindByPage(offset int, limit int) (result []*model.ActivityUser, count int64, err error) {
	result, err = a.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = a.Offset(-1).Limit(-1).Count()
	return
}

func (a activityUserDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = a.Count()
	if err != nil {
		return
	}

	err = a.Offset(offset).Limit(limit).Scan(result)
	return
}

func (a activityUserDo) Scan(result interface{}) (err error) {
	return a.DO.Scan(result)
}

func (a activityUserDo) Delete(models ...*model.ActivityUser) (result gen.ResultInfo, err error) {
	return a.DO.Delete(models)
}

func (a *activityUserDo) withDO(do gen.Dao) *activityUserDo {
	a.DO = *do.(*gen.DO)
	return a
}
