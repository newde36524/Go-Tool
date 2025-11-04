package sqlx

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/lann/builder"
	builder2 "github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

func isNoColumns(b squirrel.SelectBuilder) bool {
	v, ok := builder.Get(b, "Columns")
	data, ok2 := v.([]squirrel.Sqlizer)
	return !ok || ok2 && len(data) == 0
}

type ExplainItem struct {
	Id           int64          `db:"id"`
	SelectType   string         `db:"select_type"`
	Table        string         `db:"table"`
	Partitions   sql.NullString `db:"partitions"`
	Type         string         `db:"type"`
	PossibleKeys sql.NullString `db:"possible_keys"`
	Key          sql.NullString `db:"key"`
	KeyLen       sql.NullInt64  `db:"key_len"`
	Ref          sql.NullString `db:"ref"`
	Rows         int            `db:"rows"`
	Filtered     float64        `db:"filtered"`
	Extra        sql.NullString `db:"extra"`
}

func (e ExplainItem) Check() error {
	if !e.Key.Valid || e.Type == "ALL" {
		return fmt.Errorf("[%s] type:%s Extra:%s", e.Table, e.Type, e.Extra.String)
	}
	return nil
}

func Explain(ctx context.Context, session sqlx.Session, b squirrel.SelectBuilder) error {
	query, values, err := b.ToSql()
	if err != nil {
		return err
	}
	var explains []ExplainItem
	sql := fmt.Sprintf("explain %s", query)
	err = session.QueryRowsCtx(ctx, &explains, sql, values...)
	if err != nil {
		return err
	}
	var result error
	for _, v := range explains {
		if err := v.Check(); err != nil {
			result = errors.Join(result, err)
		}
	}
	if result != nil {
		sqlInfo := fmt.Errorf("%s %v", sql, values)
		result = errors.Join(sqlInfo, result)
	}
	return result
}

// Get 查询并返回一条数据(不支持map类型)
func Get[T any](ctx context.Context, session sqlx.Session, builder squirrel.SelectBuilder) (*T, error) {
	var resp T
	if isNoColumns(builder) {
		feilds := rawFieldNames(resp)
		builder = builder.Columns(strings.Join(feilds, ","))
	}
	query, values, err := builder.Limit(1).ToSql()
	if err != nil {
		return &resp, err
	}
	err = session.QueryRowCtx(ctx, &resp, query, values...)
	return &resp, err
}

// List 查询并返回所有数据(不支持map类型)
func List[T any](ctx context.Context, session sqlx.Session, builder squirrel.SelectBuilder) ([]T, error) {
	if isNoColumns(builder) {
		feilds := rawFieldNames(CreateInstance[T]())
		builder = builder.Columns(strings.Join(feilds, ","))
	}
	query, values, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	var resp []T
	err = session.QueryRowsCtx(ctx, &resp, query, values...)
	return resp, err
}

// Page 分页查询并返回分页数据
func Page[T any](ctx context.Context, session sqlx.Session, builder squirrel.SelectBuilder, pageNum, pageSize int) ([]T, error) {
	if isNoColumns(builder) {
		feilds := rawFieldNames(CreateInstance[T]())
		builder = builder.Columns(strings.Join(feilds, ","))
	}
	offset := (pageNum - 1) * pageSize
	if pageNum > 0 && pageSize > 0 {
		builder = builder.Offset(uint64(offset)).Limit(uint64(pageSize))
	}
	query, values, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	var resp []T
	err = session.QueryRowsCtx(ctx, &resp, query, values...)
	return resp, err
}

// Count 查询数量
func Count(ctx context.Context, session sqlx.Session, b squirrel.SelectBuilder) (int64, error) {
	columns, ok := builder.Get(b, "Columns")
	var b2 squirrel.SelectBuilder
	if ok {
		b2 = builder.Set(b, "Columns", nil).(squirrel.SelectBuilder)
		defer func() {
			builder.Set(b, "Columns", columns)
		}()
	}
	b2 = b2.Columns("COUNT(0)")

	query, values, err := b2.ToSql()
	if err != nil {
		return 0, err
	}
	var total int64
	err = session.QueryRowCtx(ctx, &total, query, values...)
	if err != nil {
		return 0, err
	}
	return total, err
}

// Exist 判断是否存在
func Exist(ctx context.Context, session sqlx.Session, b squirrel.SelectBuilder) (bool, error) {
	num, err := Count(ctx, session, b.Limit(1))
	if err != nil {
		return false, err
	}
	return num > 0, nil
}

// PageCount 分页查询并返回列表和总数
func PageCount[T any](ctx context.Context, session sqlx.Session, b squirrel.SelectBuilder, pageNum, pageSize int) (list []T, count int64, err error) {
	count, err = Count(ctx, session, b)
	if err != nil {
		return
	}
	list, err = Page[T](ctx, session, b, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return
}

// Add 创建一条记录
func Add[T any](ctx context.Context, session sqlx.Session, tableName string, data T) (sql.Result, error) {
	fieldNames := rawFieldNames(data)
	values := strings.Split(strings.Repeat("?", len(fieldNames)), "")
	query := fmt.Sprintf("insert into %s (%s) values (%s)", tableName, strings.Join(fieldNames, ","), strings.Join(values, ","))
	mp := toMap(data)
	args := make([]any, 0)
	for _, fieldName := range fieldNames {
		if v, ok := mp[fieldName].(time.Time); ok {
			zero := time.Time{}
			if v.Equal(zero) {
				mp[fieldName] = time.Unix(0, 0).Local()
				if fieldName == "`create_time`" || fieldName == "`update_time`" {
					mp[fieldName] = time.Now().Local()
				}
			}
		}
		args = append(args, mp[fieldName])
	}
	sqlResult, err := session.ExecCtx(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return sqlResult, nil
}

// GetOrSet 查询或插入数据 @hadExist 表示已存在老数据
func GetOrSet[T any](ctx context.Context, session sqlx.Session, builder squirrel.SelectBuilder, data T) (t *T, hadExist bool, err error) {
	tableName := getTableName(builder)
	t, err = Get[T](ctx, session, builder)
	hadExist = err == nil
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			_, err = Add(ctx, session, tableName, data)
			if err != nil {
				return t, false, err
			}
			t, err = Get[T](ctx, session, builder)
			return
		}
		return t, false, err
	}
	return
}

// ExistOrSet 判断是否存在或插入数据
func ExistOrSet[T any](ctx context.Context, session sqlx.Session, builder squirrel.SelectBuilder, data T) error {
	exist, err := Exist(ctx, session, builder)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	tableName := getTableName(builder)
	_, err = Add(ctx, session, tableName, data)
	return err
}

// ForceIndex 强制使用索引
func ForceIndex(builder squirrel.SelectBuilder, indexName string) squirrel.SelectBuilder {
	tableName := getTableName(builder)
	return builder.From(fmt.Sprintf("%s FORCE INDEX(%s)", tableName, indexName))
}

// AddWithWhere 根据条件插入或更新单表数据
func AddWithWhere[T any](ctx context.Context, session sqlx.Session, builder squirrel.SelectBuilder, data T) error {
	tableName := getTableName(builder)
	v, err := Get[T](ctx, session, builder)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			_, err = Add(ctx, session, tableName, data)
			if err != nil {
				return err
			}
			_, err = Get[T](ctx, session, builder)
			return err
		}
		return err
	}
	fieldNames := rawFieldNames(data)
	fieldNames = stringx.Remove(fieldNames, "`id`", "`create_time`", "`update_time`")
	mp1 := toMap(v)
	mp2 := toMap(data)
	for _, fieldName := range fieldNames {
		mp1[fieldName] = mp2[fieldName]
	}
	if _, ok := mp1["`update_time`"]; ok {
		mp1["`update_time`"] = time.Now().Local()
	}
	args := make([]any, 0)
	for _, fieldName := range fieldNames {
		if v, ok := mp1[fieldName].(time.Time); ok {
			zero := time.Time{}
			if v.Equal(zero) {
				mp1[fieldName] = time.Unix(0, 0).Local()
			}
		}
		args = append(args, mp1[fieldName])
	}
	args = append(args, mp1["`id`"])
	query := fmt.Sprintf("update %s set %s where `id` = ?", tableName, strings.Join(stringx.Remove(fieldNames, "`id`"), "=?,")+"=?")
	_, err = session.ExecCtx(ctx, query, args...)
	return err
}

// UpdateSet 执行更新语句
func UpdateSet(ctx context.Context, session sqlx.Session, builder squirrel.UpdateBuilder) (sql.Result, error) {
	query, values, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	return session.ExecCtx(ctx, query, values...)
}

// Update 更新一条记录(id字段必须)
func Update[T any](ctx context.Context, session sqlx.Session, tableName string, data T) (sql.Result, error) {
	fieldNames := rawFieldNames(data)
	fieldNames = stringx.Remove(fieldNames, "`id`", "`create_time`", "`update_time`")
	fieldNamesWithPlaceHolder := strings.Join(fieldNames, "=?,") + "=?"
	mp := toMap(data)
	query := fmt.Sprintf("update %s set %s where `id` = ?", tableName, fieldNamesWithPlaceHolder)
	args := make([]any, 0)
	for _, fieldName := range fieldNames {
		if v, ok := mp[fieldName].(time.Time); ok {
			zero := time.Time{}
			if v.Equal(zero) {
				mp[fieldName] = time.Unix(0, 0).Local()
			}
		}
		args = append(args, mp[fieldName])
	}
	args = append(args, mp["`id`"])
	return session.ExecCtx(ctx, query, args...)
}

// Upsert 插入或更新单表数据(如果存在字段设置了唯一索引, 则只会修改那一条数据)
func Upsert[T any](ctx context.Context, session sqlx.Session, tableName string, data T) error {
	fieldNames := rawFieldNames(data)
	fieldNames = stringx.Remove(fieldNames, "`id`", "`create_time`", "`update_time`")
	values := strings.Split(strings.Repeat("?", len(fieldNames)), "")
	mp := toMap(data)
	args := make([]any, 0)
	for _, fieldName := range fieldNames {
		if v, ok := mp[fieldName].(time.Time); ok {
			zero := time.Time{}
			if v.Equal(zero) {
				mp[fieldName] = time.Unix(0, 0).Local()
			}
		}
		args = append(args, mp[fieldName])
	}
	args = append(args, args...)
	query := fmt.Sprintf("insert into %s (%s) values (%s) ON DUPLICATE KEY UPDATE %s;", tableName, strings.Join(fieldNames, ","), strings.Join(values, ","), strings.Join(stringx.Remove(fieldNames, "`id`"), "=?,")+"=?")
	_, err := session.ExecCtx(ctx, query, args...)
	return err
}

// OptimisticLock 乐观锁
func OptimisticLock[T any](ctx context.Context, session sqlx.Session, b squirrel.SelectBuilder, fn func(ctx context.Context, session sqlx.Session, in *T) error) error {
	data, err := Get[T](ctx, session, b)
	if err != nil {
		return err
	}
	mp := toMap(data)
	if err := fn(ctx, session, data); err != nil {
		return err
	}
	id, ok := mp["`id`"]
	if !ok {
		return sqlx.ErrNotSettable
	}
	v, ok := mp["`version`"]
	if !ok {
		return sqlx.ErrNotSettable
	}
	version := v.(int64)
	tableName := getTableName(b)
	builder := squirrel.Update(tableName).
		Where(`del_state = 0`).
		Where(`id = ? and version = ?`, id, version).
		Set(`version`, version+1)
	result, err := UpdateSet(ctx, session, builder)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sqlx.ErrNotSettable
	}
	return err
}

// Delete 删除指定查询条件的记录
func Delete(ctx context.Context, session sqlx.Session, b squirrel.SelectBuilder) error {
	tableName := getTableName(b)
	whereParts, ok := builder.Get(b, "WhereParts")
	if !ok {
		return sqlx.ErrNotSettable
	}
	sqlizers := whereParts.([]squirrel.Sqlizer)
	wheres := make([]string, 0)
	values := make([]any, 0)
	for _, sqlizer := range sqlizers {
		where, value, err := sqlizer.ToSql()
		if err != nil {
			return err
		}
		wheres = append(wheres, where)
		values = append(values, value...)
	}
	delStr := fmt.Sprintf("delete from %s where %s", tableName, strings.Join(wheres, " and "))
	if limit, ok := builder.Get(b, "Limit"); ok {
		delStr = fmt.Sprintf("%s limit %s", delStr, limit)
	}
	_, err := session.ExecCtx(ctx, delStr, values...)
	return err
}

// toMap 将给定对象的字段名和字段值映射到map中
func toMap[T any](in T) map[string]any {
	getDbTag := func(fi reflect.StructField) string {
		dbTag := "db"
		tagv := fi.Tag.Get(dbTag)
		switch tagv {
		case "-":
			return ""
		case "":
			return fmt.Sprintf("`%s`", fi.Name)
		default:
			// get tag name with the tag option, e.g.:
			// `db:"id"`
			// `db:"id,type=char,length=16"`
			// `db:",type=char,length=16"`
			// `db:"-,type=char,length=16"`
			if strings.Contains(tagv, ",") {
				tagv = strings.TrimSpace(strings.Split(tagv, ",")[0])
			}
			if tagv == "-" {
				return ""
			}
			if len(tagv) == 0 {
				tagv = fi.Name
			}
			return fmt.Sprintf("`%s`", tagv)
		}
	}
	var typ reflect.Type
	var val = reflect.ValueOf(in)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = val.Type()
	} else {
		val = reflect.ValueOf(in)
		typ = val.Type()
	}
	switch val.Kind() {
	case reflect.Struct:
		// 创建一个map来存储字段名和字段值
		m := make(map[string]any)
		// 遍历结构体的所有字段
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			// 使用reflect.ValueOf获取字段的值
			fieldVal := val.Field(i)
			// 将字段名和值添加到map中
			m[getDbTag(field)] = fieldVal.Interface()
		}
		return m
	case reflect.Map:
		// 如果是map类型，则创建一个新的map，key为"key"，value为原map的值
		newMap := make(map[string]any)
		for _, k := range val.MapKeys() {
			str := fmt.Sprintf("%v", k.Interface())
			str = strings.Trim(str, "`")
			newMap[fmt.Sprintf("`%v`", str)] = val.MapIndex(k).Interface()
		}
		return newMap
	default:
		panic(fmt.Sprintf("ToMap expects a struct or map, err kind:%s", typ.Kind().String()))
	}
}

// getTableName 获取表名
func getTableName(b squirrel.SelectBuilder) string {
	from, ok := builder.Get(b, "From")
	if !ok {
		return ""
	}
	sql, ok := from.(squirrel.Sqlizer)
	if !ok {
		return ""
	}
	query, _, err := sql.ToSql()
	if err != nil {
		return ""
	}
	return query
}

// rawFieldNames 获取类型的字段名称列表
//
//nolint:errchkjson
func rawFieldNames[T any](in T) []string {
	var typ reflect.Type
	var val = reflect.ValueOf(in)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = val.Type()
	} else {
		val = reflect.ValueOf(in)
		typ = val.Type()
	}
	switch val.Kind() {
	case reflect.Struct:
		fieldNames := builder2.RawFieldNames(in)
		return fieldNames
	case reflect.Map:
		// 准备一个切片来存储所有的键
		keys := make([]string, 0, val.Len())
		// 遍历map中的每个键值对
		for _, k := range val.MapKeys() {
			// 将键添加到切片中
			str := fmt.Sprintf("%v", k.Interface())
			str = strings.Trim(str, "`")
			keys = append(keys, fmt.Sprintf("`%v`", str))
		}
		return keys
	default:
		bs, _ := json.Marshal(in)
		panic(fmt.Sprintf("rawFieldNames expects a struct or map, err kind:%s, value:%s", typ.Kind().String(), string(bs)))
	}
}

func CreateInstance[T any](size ...int) (t T) {
	s := 0
	if len(size) > 0 {
		s = size[0]
	}
	switch reflect.TypeFor[T]().Kind() {
	case reflect.Array:
		return reflect.New(reflect.TypeFor[T]()).Elem().Interface().(T)
	case reflect.Map:
		return reflect.MakeMap(reflect.TypeFor[T]()).Interface().(T)
	case reflect.Pointer:
		return reflect.New(reflect.TypeFor[T]().Elem()).Interface().(T)
	case reflect.Slice:
		return reflect.MakeSlice(reflect.TypeFor[T](), s, s).Interface().(T)
	case reflect.Chan:
		return reflect.MakeChan(reflect.TypeFor[T](), s).Interface().(T)
	}
	return
}
