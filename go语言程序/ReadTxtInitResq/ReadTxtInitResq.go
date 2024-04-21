package ReadTxtInitResq

import (
	"2024_4_15/Dingyi"
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// 项目变更后只需要修改结构体的内容以及结构体值的设置这两个地方

// 结构体的构造函数
func NewTextInfo() *Dingyi.TextInfo {
	return &Dingyi.TextInfo{}
}

// 将文本信息结构体录入到map中去
func AddToAllInfoOfTextMap(AllInfoOfTextMap *map[string]*Dingyi.TextInfo, A *Dingyi.TextInfo, Target string) {
	tmp := GetFieldValue(A, Target)
	use := tmp.(string)
	(*AllInfoOfTextMap)[use] = A
}

func ReadTxtInitResq(dir string, WhoisMapKey string) (*map[string]*Dingyi.TextInfo, []*Dingyi.TextInfo, error) {
	file, err := os.OpenFile(dir, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println("文件打开失败:", err)
		return nil, nil, errors.New("文件打开失败")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("文件关闭错误")
			return
		}
	}(file)
	var allInfoOfTextMap = make(map[string]*Dingyi.TextInfo)
	var data []*Dingyi.TextInfo // 存储读取的每行数据的 TextInfo 实例
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "----")
		// 创建一个 TextInfo 实例，并一次性设置多个字段值
		info_struct := NewTextInfo()
		info_struct.SetValues(parts)
		// 将实例添加到 data 切片中
		data = append(data, info_struct)
		AddToAllInfoOfTextMap(&allInfoOfTextMap, info_struct, WhoisMapKey)
	}
	return &allInfoOfTextMap, data, nil
	// 处理数据...
}

// GetFieldValue 函数用于获取结构体中指定字段的值
func GetFieldValue(info *Dingyi.TextInfo, fieldName string) interface{} {
	// 使用反射获取结构体类型信息
	t := reflect.TypeOf(*info)
	// 使用反射获取结构体字段值
	v := reflect.ValueOf(*info)

	// 遍历结构体字段，查找指定字段名
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == fieldName {
			return v.Field(i).Interface()
		}
	}
	// 如果没有找到指定字段名，返回 nil
	return nil
}
