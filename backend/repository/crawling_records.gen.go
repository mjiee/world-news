// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/mjiee/world-news/backend/repository/model"
)

func newCrawlingRecord(db *gorm.DB, opts ...gen.DOOption) crawlingRecord {
	_crawlingRecord := crawlingRecord{}

	_crawlingRecord.crawlingRecordDo.UseDB(db, opts...)
	_crawlingRecord.crawlingRecordDo.UseModel(&model.CrawlingRecord{})

	tableName := _crawlingRecord.crawlingRecordDo.TableName()
	_crawlingRecord.ALL = field.NewAsterisk(tableName)
	_crawlingRecord.ID = field.NewUint(tableName, "id")
	_crawlingRecord.RecordType = field.NewString(tableName, "record_type")
	_crawlingRecord.Date = field.NewTime(tableName, "date")
	_crawlingRecord.Quantity = field.NewInt64(tableName, "quantity")
	_crawlingRecord.Status = field.NewString(tableName, "status")
	_crawlingRecord.Config = field.NewString(tableName, "config")
	_crawlingRecord.CreatedAt = field.NewTime(tableName, "created_at")
	_crawlingRecord.UpdatedAt = field.NewTime(tableName, "updated_at")

	_crawlingRecord.fillFieldMap()

	return _crawlingRecord
}

type crawlingRecord struct {
	crawlingRecordDo crawlingRecordDo

	ALL        field.Asterisk
	ID         field.Uint
	RecordType field.String
	Date       field.Time
	Quantity   field.Int64
	Status     field.String
	Config     field.String
	CreatedAt  field.Time
	UpdatedAt  field.Time

	fieldMap map[string]field.Expr
}

func (c crawlingRecord) Table(newTableName string) *crawlingRecord {
	c.crawlingRecordDo.UseTable(newTableName)
	return c.updateTableName(newTableName)
}

func (c crawlingRecord) As(alias string) *crawlingRecord {
	c.crawlingRecordDo.DO = *(c.crawlingRecordDo.As(alias).(*gen.DO))
	return c.updateTableName(alias)
}

func (c *crawlingRecord) updateTableName(table string) *crawlingRecord {
	c.ALL = field.NewAsterisk(table)
	c.ID = field.NewUint(table, "id")
	c.RecordType = field.NewString(table, "record_type")
	c.Date = field.NewTime(table, "date")
	c.Quantity = field.NewInt64(table, "quantity")
	c.Status = field.NewString(table, "status")
	c.Config = field.NewString(table, "config")
	c.CreatedAt = field.NewTime(table, "created_at")
	c.UpdatedAt = field.NewTime(table, "updated_at")

	c.fillFieldMap()

	return c
}

func (c *crawlingRecord) WithContext(ctx context.Context) *crawlingRecordDo {
	return c.crawlingRecordDo.WithContext(ctx)
}

func (c crawlingRecord) TableName() string { return c.crawlingRecordDo.TableName() }

func (c crawlingRecord) Alias() string { return c.crawlingRecordDo.Alias() }

func (c crawlingRecord) Columns(cols ...field.Expr) gen.Columns {
	return c.crawlingRecordDo.Columns(cols...)
}

func (c *crawlingRecord) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := c.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (c *crawlingRecord) fillFieldMap() {
	c.fieldMap = make(map[string]field.Expr, 8)
	c.fieldMap["id"] = c.ID
	c.fieldMap["record_type"] = c.RecordType
	c.fieldMap["date"] = c.Date
	c.fieldMap["quantity"] = c.Quantity
	c.fieldMap["status"] = c.Status
	c.fieldMap["config"] = c.Config
	c.fieldMap["created_at"] = c.CreatedAt
	c.fieldMap["updated_at"] = c.UpdatedAt
}

func (c crawlingRecord) clone(db *gorm.DB) crawlingRecord {
	c.crawlingRecordDo.ReplaceConnPool(db.Statement.ConnPool)
	return c
}

func (c crawlingRecord) replaceDB(db *gorm.DB) crawlingRecord {
	c.crawlingRecordDo.ReplaceDB(db)
	return c
}

type crawlingRecordDo struct{ gen.DO }

func (c crawlingRecordDo) Debug() *crawlingRecordDo {
	return c.withDO(c.DO.Debug())
}

func (c crawlingRecordDo) WithContext(ctx context.Context) *crawlingRecordDo {
	return c.withDO(c.DO.WithContext(ctx))
}

func (c crawlingRecordDo) ReadDB() *crawlingRecordDo {
	return c.Clauses(dbresolver.Read)
}

func (c crawlingRecordDo) WriteDB() *crawlingRecordDo {
	return c.Clauses(dbresolver.Write)
}

func (c crawlingRecordDo) Session(config *gorm.Session) *crawlingRecordDo {
	return c.withDO(c.DO.Session(config))
}

func (c crawlingRecordDo) Clauses(conds ...clause.Expression) *crawlingRecordDo {
	return c.withDO(c.DO.Clauses(conds...))
}

func (c crawlingRecordDo) Returning(value interface{}, columns ...string) *crawlingRecordDo {
	return c.withDO(c.DO.Returning(value, columns...))
}

func (c crawlingRecordDo) Not(conds ...gen.Condition) *crawlingRecordDo {
	return c.withDO(c.DO.Not(conds...))
}

func (c crawlingRecordDo) Or(conds ...gen.Condition) *crawlingRecordDo {
	return c.withDO(c.DO.Or(conds...))
}

func (c crawlingRecordDo) Select(conds ...field.Expr) *crawlingRecordDo {
	return c.withDO(c.DO.Select(conds...))
}

func (c crawlingRecordDo) Where(conds ...gen.Condition) *crawlingRecordDo {
	return c.withDO(c.DO.Where(conds...))
}

func (c crawlingRecordDo) Order(conds ...field.Expr) *crawlingRecordDo {
	return c.withDO(c.DO.Order(conds...))
}

func (c crawlingRecordDo) Distinct(cols ...field.Expr) *crawlingRecordDo {
	return c.withDO(c.DO.Distinct(cols...))
}

func (c crawlingRecordDo) Omit(cols ...field.Expr) *crawlingRecordDo {
	return c.withDO(c.DO.Omit(cols...))
}

func (c crawlingRecordDo) Join(table schema.Tabler, on ...field.Expr) *crawlingRecordDo {
	return c.withDO(c.DO.Join(table, on...))
}

func (c crawlingRecordDo) LeftJoin(table schema.Tabler, on ...field.Expr) *crawlingRecordDo {
	return c.withDO(c.DO.LeftJoin(table, on...))
}

func (c crawlingRecordDo) RightJoin(table schema.Tabler, on ...field.Expr) *crawlingRecordDo {
	return c.withDO(c.DO.RightJoin(table, on...))
}

func (c crawlingRecordDo) Group(cols ...field.Expr) *crawlingRecordDo {
	return c.withDO(c.DO.Group(cols...))
}

func (c crawlingRecordDo) Having(conds ...gen.Condition) *crawlingRecordDo {
	return c.withDO(c.DO.Having(conds...))
}

func (c crawlingRecordDo) Limit(limit int) *crawlingRecordDo {
	return c.withDO(c.DO.Limit(limit))
}

func (c crawlingRecordDo) Offset(offset int) *crawlingRecordDo {
	return c.withDO(c.DO.Offset(offset))
}

func (c crawlingRecordDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *crawlingRecordDo {
	return c.withDO(c.DO.Scopes(funcs...))
}

func (c crawlingRecordDo) Unscoped() *crawlingRecordDo {
	return c.withDO(c.DO.Unscoped())
}

func (c crawlingRecordDo) Create(values ...*model.CrawlingRecord) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Create(values)
}

func (c crawlingRecordDo) CreateInBatches(values []*model.CrawlingRecord, batchSize int) error {
	return c.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (c crawlingRecordDo) Save(values ...*model.CrawlingRecord) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Save(values)
}

func (c crawlingRecordDo) First() (*model.CrawlingRecord, error) {
	if result, err := c.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.CrawlingRecord), nil
	}
}

func (c crawlingRecordDo) Take() (*model.CrawlingRecord, error) {
	if result, err := c.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.CrawlingRecord), nil
	}
}

func (c crawlingRecordDo) Last() (*model.CrawlingRecord, error) {
	if result, err := c.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.CrawlingRecord), nil
	}
}

func (c crawlingRecordDo) Find() ([]*model.CrawlingRecord, error) {
	result, err := c.DO.Find()
	return result.([]*model.CrawlingRecord), err
}

func (c crawlingRecordDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.CrawlingRecord, err error) {
	buf := make([]*model.CrawlingRecord, 0, batchSize)
	err = c.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (c crawlingRecordDo) FindInBatches(result *[]*model.CrawlingRecord, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return c.DO.FindInBatches(result, batchSize, fc)
}

func (c crawlingRecordDo) Attrs(attrs ...field.AssignExpr) *crawlingRecordDo {
	return c.withDO(c.DO.Attrs(attrs...))
}

func (c crawlingRecordDo) Assign(attrs ...field.AssignExpr) *crawlingRecordDo {
	return c.withDO(c.DO.Assign(attrs...))
}

func (c crawlingRecordDo) Joins(fields ...field.RelationField) *crawlingRecordDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Joins(_f))
	}
	return &c
}

func (c crawlingRecordDo) Preload(fields ...field.RelationField) *crawlingRecordDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Preload(_f))
	}
	return &c
}

func (c crawlingRecordDo) FirstOrInit() (*model.CrawlingRecord, error) {
	if result, err := c.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.CrawlingRecord), nil
	}
}

func (c crawlingRecordDo) FirstOrCreate() (*model.CrawlingRecord, error) {
	if result, err := c.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.CrawlingRecord), nil
	}
}

func (c crawlingRecordDo) FindByPage(offset int, limit int) (result []*model.CrawlingRecord, count int64, err error) {
	result, err = c.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = c.Offset(-1).Limit(-1).Count()
	return
}

func (c crawlingRecordDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = c.Count()
	if err != nil {
		return
	}

	err = c.Offset(offset).Limit(limit).Scan(result)
	return
}

func (c crawlingRecordDo) Scan(result interface{}) (err error) {
	return c.DO.Scan(result)
}

func (c crawlingRecordDo) Delete(models ...*model.CrawlingRecord) (result gen.ResultInfo, err error) {
	return c.DO.Delete(models)
}

func (c *crawlingRecordDo) withDO(do gen.Dao) *crawlingRecordDo {
	c.DO = *do.(*gen.DO)
	return c
}
