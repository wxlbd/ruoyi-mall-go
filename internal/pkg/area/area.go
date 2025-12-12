package area

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"strconv"
	"strings"
	"sync"

	"go.uber.org/zap"
)

//go:embed data/area.csv
var areaCSV []byte

//go:embed data/ip2region.xdb
var IP2RegionXDB []byte

// Area 地区节点
type Area struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Type     int     `json:"type"`
	Parent   *Area   `json:"-"`
	Children []*Area `json:"children,omitempty"`
}

// 常量定义
const (
	IDGlobal = 0 // 全球
	IDChina  = 1 // 中国
)

// AreaTypeEnum 区域类型枚举
const (
	AreaTypeCountry  = 1 // 国家
	AreaTypeProvince = 2 // 省份
	AreaTypeCity     = 3 // 城市
	AreaTypeDistrict = 4 // 区县
)

var (
	areas    map[int]*Area
	areaOnce sync.Once
	initErr  error
)

// Init 初始化地区数据（使用嵌入的 CSV 数据）
func Init(_ string) error {
	areaOnce.Do(func() {
		initErr = loadFromEmbeddedCSV()
	})
	return initErr
}

// loadFromEmbeddedCSV 从嵌入的 CSV 数据加载地区
func loadFromEmbeddedCSV() error {
	reader := csv.NewReader(bytes.NewReader(areaCSV))
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// 初始化 areas map
	areas = make(map[int]*Area)

	// 添加全球根节点
	areas[IDGlobal] = &Area{
		ID:       IDGlobal,
		Name:     "全球",
		Type:     0,
		Children: make([]*Area, 0),
	}

	// 跳过 header，解析所有行
	for i, row := range records {
		if i == 0 {
			continue // 跳过 header: id,name,type,parentId
		}
		if len(row) < 4 {
			continue
		}

		id, _ := strconv.Atoi(row[0])
		name := row[1]
		areaType, _ := strconv.Atoi(row[2])

		area := &Area{
			ID:       id,
			Name:     name,
			Type:     areaType,
			Children: make([]*Area, 0),
		}
		areas[id] = area
	}

	// 构建父子关系
	for i, row := range records {
		if i == 0 {
			continue
		}
		if len(row) < 4 {
			continue
		}

		id, _ := strconv.Atoi(row[0])
		parentID, _ := strconv.Atoi(row[3])

		area := areas[id]
		parent := areas[parentID]

		if area != nil && parent != nil {
			area.Parent = parent
			parent.Children = append(parent.Children, area)
		}
	}

	zap.L().Info("地区数据加载完成 (embedded)", zap.Int("count", len(areas)))
	return nil
}

// GetArea 获取指定编号的地区
func GetArea(id int) *Area {
	if areas == nil {
		return nil
	}
	return areas[id]
}

// GetAreaTree 获取中国地区树 (返回中国的子节点)
func GetAreaTree() []*Area {
	china := GetArea(IDChina)
	if china == nil {
		return nil
	}
	return china.Children
}

// Format 格式化地区名称
// 例如: id="静安区" 返回 "上海 上海市 静安区"
func Format(id int) string {
	return FormatWithSep(id, " ")
}

// FormatWithSep 使用指定分隔符格式化地区名称
func FormatWithSep(id int, sep string) string {
	area := GetArea(id)
	if area == nil {
		return ""
	}

	var parts []string
	for current := area; current != nil; current = current.Parent {
		// 跳过全球和中国
		if current.ID == IDGlobal || current.ID == IDChina {
			break
		}
		parts = append([]string{current.Name}, parts...)
	}

	return strings.Join(parts, sep)
}
