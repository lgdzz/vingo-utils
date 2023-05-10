package vingo

import (
	"fmt"
	"reflect"
	"strconv"
)

// 通用切片去重
// s := []string{"a", "b", "c", "a", "b", "a"}
// SliceUnique(&s)
// print：["a", "b", "c"]
//
// i := []int{1, 1, 2, 3, 3, 1, 3}
// SliceUnique(&i)
// print：[1, 2, 3]
func SliceUnique(slice interface{}) {
	uniqueMap := make(map[interface{}]struct{})
	valueOfSlice := reflect.ValueOf(slice).Elem()

	for i := 0; i < valueOfSlice.Len(); i++ {
		uniqueMap[valueOfSlice.Index(i).Interface()] = struct{}{}
	}

	valueOfSlice.Set(reflect.MakeSlice(valueOfSlice.Type(), 0, 0))

	for k := range uniqueMap {
		valueOfSlice.Set(reflect.Append(valueOfSlice, reflect.ValueOf(k)))
	}
}

// 将[]string数据去重返回
func SliceUniqueString(slice []string) []string {
	uniqueMap := make(map[string]interface{})
	for _, v := range slice {
		uniqueMap[v] = nil
	}
	var uniqueSlice []string
	for k := range uniqueMap {
		uniqueSlice = append(uniqueSlice, k)
	}
	return uniqueSlice
}

// 将[]uint数据去重返回
func SliceUniqueUint(slice []uint) []uint {
	uniqueMap := make(map[uint]interface{})
	for _, v := range slice {
		uniqueMap[v] = nil
	}
	var uniqueSlice []uint
	for k := range uniqueMap {
		uniqueSlice = append(uniqueSlice, k)
	}
	return uniqueSlice
}

// 将[]int数据去重返回
func SliceUniqueInt(slice []int) []int {
	uniqueMap := make(map[int]interface{})
	for _, v := range slice {
		uniqueMap[v] = nil
	}
	var uniqueSlice []int
	for k := range uniqueMap {
		uniqueSlice = append(uniqueSlice, k)
	}
	return uniqueSlice
}

// []string转[]int
func SliceStringToInt(s []string) []int {
	slice := make([]int, 0)
	for _, v := range s {
		num, _ := strconv.Atoi(v)
		slice = append(slice, num)
	}
	return slice
}

// []string转[]uint
func SliceStringToUint(s []string) []uint {
	slice := make([]uint, 0)
	for _, v := range s {
		num, _ := strconv.Atoi(v)
		slice = append(slice, uint(num))
	}
	return slice
}

// []string转[]float64
func SliceStringToFloat64(s []string) []float64 {
	slice := make([]float64, 0)
	for _, v := range s {
		num, _ := strconv.Atoi(v)
		slice = append(slice, float64(num))
	}
	return slice
}

// []int转[]string
func SliceIntToString(s []int) []string {
	slice := make([]string, 0)
	for _, v := range s {
		slice = append(slice, strconv.Itoa(v))
	}
	return slice
}

// []uint转[]string
func SliceUintToString(s []uint) []string {
	slice := make([]string, 0)
	for _, v := range s {
		slice = append(slice, strconv.Itoa(int(v)))
	}
	return slice
}

// []float64转[]string
func SliceFloat64ToString(s []float64) []string {
	slice := make([]string, 0)
	for _, v := range s {
		slice = append(slice, fmt.Sprintf("%f", v))
	}
	return slice
}

// 判断一个节点是否在切片中，与IsInSlice函数不同，该函数支持更多场景，而IsInSlice只适合切片类型
// 判断字符串是否在字符串切片中
// 判断数字是否在整型切片中
// 判断字符串是否在字符串字典中
// 判断结构体是否在结构体切片中
func IsInSliceAny(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

// 判断一个节点是否在切片中
func IsInSlice(item interface{}, items interface{}) bool {
	s := reflect.ValueOf(items)
	if s.Kind() != reflect.Slice {
		panic("not a slice")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func GetSliceElement(slice interface{}, index int) interface{} {
	value := reflect.ValueOf(slice)
	if value.Kind() != reflect.Slice {
		panic("not a slice")
	}
	if index >= value.Len() {
		panic(fmt.Sprintf("Index out of range: %d", index))
	}
	element := value.Index(index)
	if !element.IsValid() {
		panic(fmt.Sprintf("Element does not exist: %d", index))
	}
	return element.Interface()
}

// 从切片中删除元素
func SliceDelItem(item interface{}, items interface{}) {
	value := reflect.ValueOf(items)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Slice {
		panic("not a slice pointer")
	}

	sliceValue := value.Elem()
	for i := 0; i < sliceValue.Len(); i++ {
		if reflect.DeepEqual(sliceValue.Index(i).Interface(), item) {
			// 将要删除的元素移到最后一个元素位置
			lastIndex := sliceValue.Len() - 1
			lastElement := sliceValue.Index(lastIndex)
			sliceValue.Index(i).Set(lastElement)

			// 切片长度减一
			newSliceValue := sliceValue.Slice(0, lastIndex)
			sliceValue.Set(newSliceValue)

			return
		}
	}
}

// 在切片中搜索元素，返回索引，-1未找到
func IndexOf(item interface{}, items interface{}) int {
	value := reflect.ValueOf(items)
	if value.Kind() != reflect.Slice {
		panic("not a slice")
	}

	for i := 0; i < value.Len(); i++ {
		if reflect.DeepEqual(value.Index(i).Interface(), item) {
			return i
		}
	}
	return -1
}
