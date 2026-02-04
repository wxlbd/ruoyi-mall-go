package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// BitBool 映射到数据库 BIT(1) 类型的布尔值
type BitBool bool

// Scan 实现 Scanner 接口
func (b *BitBool) Scan(value interface{}) error {
	if value == nil {
		*b = false
		return nil
	}

	switch v := value.(type) {
	case []uint8:
		if len(v) > 0 {
			*b = BitBool(v[0] == 1)
		} else {
			*b = false
		}
	case int64:
		*b = BitBool(v == 1)
	case bool:
		*b = BitBool(v)
	default:
		return errors.New("incompatible type for BitBool")
	}
	return nil
}

// Value 实现 driver.Valuer 接口
func (b BitBool) Value() (driver.Value, error) {
	if b {
		return int64(1), nil
	}
	return int64(0), nil
}

// QueryClauses 实现 GORM 软删除查询子句
func (BitBool) QueryClauses(f *schema.Field) []clause.Interface {
	// 仅为标有 softDelete 标签的字段返回 QueryClause
	if f != nil && f.TagSettings != nil {
		if _, hasSoftDelete := f.TagSettings["SOFTDELETE"]; hasSoftDelete {
			return []clause.Interface{BitBoolQueryClause{Field: f}}
		}
	}
	return nil
}

// DeleteClauses 实现 GORM 软删除删除子句
func (BitBool) DeleteClauses(f *schema.Field) []clause.Interface {
	// 仅为标有 softDelete 标签的字段返回 DeleteClause
	if f != nil && f.TagSettings != nil {
		if _, hasSoftDelete := f.TagSettings["SOFTDELETE"]; hasSoftDelete {
			return []clause.Interface{BitBoolDeleteClause{Field: f}}
		}
	}
	return nil
}

// UpdateClauses 实现 GORM 软删除更新子句
func (BitBool) UpdateClauses(f *schema.Field) []clause.Interface {
	// 仅为标有 softDelete 标签的字段返回 UpdateClause
	if f != nil && f.TagSettings != nil {
		if _, hasSoftDelete := f.TagSettings["SOFTDELETE"]; hasSoftDelete {
			return []clause.Interface{BitBoolUpdateClause{Field: f}}
		}
	}
	return nil
}

// BitBoolQueryClause 查询子句 - 自动过滤已删除记录
type BitBoolQueryClause struct {
	Field *schema.Field
}

func (sd BitBoolQueryClause) Name() string {
	return ""
}

func (sd BitBoolQueryClause) Build(clause.Builder) {
}

func (sd BitBoolQueryClause) MergeClause(*clause.Clause) {
}

func (sd BitBoolQueryClause) ModifyStatement(stmt *gorm.Statement) {
	if _, ok := stmt.Clauses["soft_delete_enabled"]; !ok && !stmt.Statement.Unscoped {
		// 检查字段是否启用了软删除 (通过 softDelete tag)
		if sd.Field != nil && sd.Field.TagSettings != nil {
			if _, hasSoftDelete := sd.Field.TagSettings["SOFTDELETE"]; hasSoftDelete {
				// 添加 WHERE deleted = 0 条件
				stmt.AddClause(clause.Where{Exprs: []clause.Expression{
					clause.Eq{
						Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName},
						Value:  false, // BitBool false = 0
					},
				}})
				stmt.Clauses["soft_delete_enabled"] = clause.Clause{}
			}
		}
	}
}

// BitBoolDeleteClause 删除子句 - 软删除时设置 deleted = 1
type BitBoolDeleteClause struct {
	Field *schema.Field
}

func (sd BitBoolDeleteClause) Name() string {
	return ""
}

func (sd BitBoolDeleteClause) Build(clause.Builder) {
}

func (sd BitBoolDeleteClause) MergeClause(*clause.Clause) {
}

func (sd BitBoolDeleteClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() == 0 && !stmt.Statement.Unscoped {
		// 设置 deleted = 1
		set := clause.Set{{
			Column: clause.Column{Name: sd.Field.DBName},
			Value:  true, // BitBool true = 1
		}}
		stmt.SetColumn(sd.Field.DBName, true, true)
		stmt.AddClause(set)

		// 添加 WHERE 条件（基于主键）
		if stmt.Schema != nil {
			_, queryValues := schema.GetIdentityFieldValuesMap(stmt.Context, stmt.ReflectValue, stmt.Schema.PrimaryFields)
			column, values := schema.ToQueryValues(stmt.Table, stmt.Schema.PrimaryFieldDBNames, queryValues)

			if len(values) > 0 {
				stmt.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
			}
		}

		// 添加已删除过滤
		BitBoolQueryClause(sd).ModifyStatement(stmt)

		stmt.AddClauseIfNotExists(clause.Update{})
		stmt.Build(stmt.DB.Callback().Update().Clauses...)
	}
}

// BitBoolUpdateClause 更新子句 - 更新时自动过滤已删除记录
type BitBoolUpdateClause struct {
	Field *schema.Field
}

func (sd BitBoolUpdateClause) Name() string {
	return ""
}

func (sd BitBoolUpdateClause) Build(clause.Builder) {
}

func (sd BitBoolUpdateClause) MergeClause(*clause.Clause) {
}

func (sd BitBoolUpdateClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() == 0 && !stmt.Statement.Unscoped {
		BitBoolQueryClause(sd).ModifyStatement(stmt)
	}
}

func NewBitBool(b bool) BitBool {
	return BitBool(b)
}

// ListFromCSV 处理逗号分隔的列表（兼容 MyBatis TypeHandler）
// 支持 "1,2,3" / "a,b,c" 格式和 JSON "[1,2,3]" / "[\"a\",\"b\"]" 格式
// T 可以是数值类型 (int, int64, float64 等) 或 string
type ListFromCSV[T int | int64 | int32 | uint | uint64 | uint32 | float64 | float32 | string] []T

func (l *ListFromCSV[T]) Scan(value any) error {
	if value == nil {
		*l = nil
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return errors.New("incompatible type for ListFromCSV")
	}

	if len(data) == 0 {
		*l = nil
		return nil
	}

	str := strings.TrimSpace(string(data))
	if str == "" {
		*l = nil
		return nil
	}

	// 优先尝试 JSON 格式解析
	if strings.HasPrefix(str, "[") {
		var result []T
		if err := json.Unmarshal(data, &result); err == nil {
			*l = result
			return nil
		}
	}

	// 按逗号分隔解析
	parts := strings.Split(str, ",")
	result := make([]T, 0, len(parts))

	// 使用反射判断 T 的具体类型
	var zero T
	elemType := reflect.TypeOf(zero)

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		var val T
		switch elemType.Kind() {
		case reflect.String:
			// 字符串类型直接使用
			val = any(p).(T)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// 整数类型解析
			i, err := strconv.ParseInt(p, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse '%s' as integer: %w", p, err)
			}
			rv := reflect.New(elemType).Elem()
			rv.SetInt(i)
			val = rv.Interface().(T)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			// 无符号整数类型解析
			i, err := strconv.ParseUint(p, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse '%s' as unsigned integer: %w", p, err)
			}
			rv := reflect.New(elemType).Elem()
			rv.SetUint(i)
			val = rv.Interface().(T)
		case reflect.Float32, reflect.Float64:
			// 浮点数类型解析
			f, err := strconv.ParseFloat(p, 64)
			if err != nil {
				return fmt.Errorf("failed to parse '%s' as float: %w", p, err)
			}
			rv := reflect.New(elemType).Elem()
			rv.SetFloat(f)
			val = rv.Interface().(T)
		default:
			return fmt.Errorf("unsupported type: %v", elemType.Kind())
		}
		result = append(result, val)
	}
	*l = result
	return nil
}

func (l ListFromCSV[T]) Value() (driver.Value, error) {
	if len(l) == 0 {
		return "", nil
	}
	parts := make([]string, len(l))
	for i, v := range l {
		parts[i] = fmt.Sprintf("%v", v)
	}
	return strings.Join(parts, ","), nil
}

func (l ListFromCSV[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal([]T(l))
}

func (l *ListFromCSV[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*l = nil
		return nil
	}
	// 尝试解析为数组
	var arr []T
	if err := json.Unmarshal(data, &arr); err == nil {
		*l = arr
		return nil
	}
	// 尝试解析为单个值
	var single T
	if err := json.Unmarshal(data, &single); err == nil {
		*l = []T{single}
		return nil
	}
	// 尝试解析为字符串 (CSV 格式)
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		return l.Scan(str)
	}
	return fmt.Errorf("cannot unmarshal %s into ListFromCSV", string(data))
}

// ParseListFromCSV 从 "1,2,3" 或 "[1,2,3]" 格式的字符串解析为 ListFromCSV[T]
func ParseListFromCSV[T int | int64 | int32 | uint | uint64 | uint32 | float64 | float32 | string](s string) (ListFromCSV[T], error) {
	var result ListFromCSV[T]
	if err := result.Scan(s); err != nil {
		return nil, err
	}
	return result, nil
}

// IntListFromCSV 向后兼容的别名 (元素类型为 int)
type IntListFromCSV = ListFromCSV[int]

// Int64ListFromCSV 处理逗号分隔的 int64 列表 (对齐 Java: List<Long>)
type Int64ListFromCSV = ListFromCSV[int64]

// StringListFromCSV 处理逗号分隔的字符串列表
type StringListFromCSV = ListFromCSV[string]

// JsonDateTime 用于 JSON 序列化时输出毫秒时间戳（对齐 Java 的 TimestampLocalDateTimeSerializer）
type JsonDateTime time.Time

const jsonDateTimeLayout = "2006-01-02 15:04:05"

func (t JsonDateTime) MarshalJSON() ([]byte, error) {
	st := time.Time(t)
	if st.IsZero() {
		return []byte("null"), nil
	}
	// 对齐 Java: 输出毫秒时间戳
	return []byte(fmt.Sprintf("%d", st.UnixMilli())), nil
}

func (t *JsonDateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	s := string(data)
	// 尝试解析毫秒时间戳
	if ms, err := strconv.ParseInt(s, 10, 64); err == nil {
		*t = JsonDateTime(time.UnixMilli(ms))
		return nil
	}
	// 兼容旧的字符串格式
	now, err := time.ParseInLocation("\""+jsonDateTimeLayout+"\"", s, time.Local)
	*t = JsonDateTime(now)
	return err
}

func (t JsonDateTime) Value() (driver.Value, error) {
	st := time.Time(t)
	if st.IsZero() {
		return nil, nil
	}
	return st, nil
}

func (t *JsonDateTime) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	vt, ok := v.(time.Time)
	if !ok {
		return errors.New("invalid type for JsonDateTime")
	}
	*t = JsonDateTime(vt)
	return nil
}

func (t JsonDateTime) String() string {
	return time.Time(t).Format(jsonDateTimeLayout)
}

func ToJsonDateTime(t time.Time) JsonDateTime {
	return JsonDateTime(t)
}

func ToJsonDateTimePtr(t *time.Time) *JsonDateTime {
	if t == nil {
		return nil
	}
	jt := JsonDateTime(*t)
	return &jt
}

// TimeOfDay 处理数据库 TIME 类型（仅存储 HH:MM:SS）
// 存储为 "15:04:05" 格式的字符串
type TimeOfDay string

const timeOfDayLayout = "15:04:05"

// Scan 实现 Scanner 接口
func (t *TimeOfDay) Scan(value interface{}) error {
	if value == nil {
		*t = ""
		return nil
	}

	switch v := value.(type) {
	case []uint8:
		*t = TimeOfDay(string(v))
	case string:
		*t = TimeOfDay(v)
	case time.Time:
		*t = TimeOfDay(v.Format(timeOfDayLayout))
	default:
		return fmt.Errorf("incompatible type for TimeOfDay: %T", v)
	}
	return nil
}

// Value 实现 driver.Valuer 接口
func (t TimeOfDay) Value() (driver.Value, error) {
	if t == "" {
		return nil, nil
	}
	return string(t), nil
}

// MarshalJSON 实现 JSON 序列化
func (t TimeOfDay) MarshalJSON() ([]byte, error) {
	if t == "" {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", t)), nil
}

// UnmarshalJSON 实现 JSON 反序列化
func (t *TimeOfDay) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || len(data) == 0 {
		*t = ""
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*t = TimeOfDay(s)
	return nil
}

// String 返回字符串表示
func (t TimeOfDay) String() string {
	return string(t)
}
