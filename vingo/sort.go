package vingo

import (
	"reflect"
	"sort"
	"time"
)

const (
	Ascending  = "ASC"
	Descending = "DESC"
)

type ByField struct {
	data       interface{}
	sortOrders [][]string
}

func NewByField(data interface{}, sortOrders [][]string) *ByField {
	return &ByField{
		data:       data,
		sortOrders: sortOrders,
	}
}

func (bf *ByField) Sort() {
	val := reflect.ValueOf(bf.data)
	if val.Kind() != reflect.Slice {
		panic("data must be a slice")
	}

	sort.SliceStable(bf.data, func(i, j int) bool {
		v1 := val.Index(i)
		v2 := val.Index(j)

		for _, item := range bf.sortOrders {
			field := item[0]
			sortOrder := item[1]
			f1 := v1.FieldByName(field)
			f2 := v2.FieldByName(field)

			if !f1.IsValid() || !f2.IsValid() {
				panic("Invalid field name: " + field)
			}

			switch f1.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				v1, v2 := f1.Int(), f2.Int()
				if v1 != v2 {
					if sortOrder == Ascending {
						return v1 < v2
					} else {
						return v1 > v2
					}
				}
			case reflect.Float32, reflect.Float64:
				v1, v2 := f1.Float(), f2.Float()
				if v1 != v2 {
					if sortOrder == Ascending {
						return v1 < v2
					} else {
						return v1 > v2
					}
				}
			case reflect.String:
				v1, v2 := f1.String(), f2.String()
				if v1 != v2 {
					if sortOrder == Ascending {
						return v1 < v2
					} else {
						return v1 > v2
					}
				}
			case reflect.Struct:
				if f1.Type() == reflect.TypeOf(time.Time{}) {
					t1, t2 := f1.Interface().(time.Time), f2.Interface().(time.Time)
					if !t1.Equal(t2) {
						if sortOrder == Ascending {
							return t1.Before(t2)
						} else {
							return t1.After(t2)
						}
					}
				}
			case reflect.Ptr:
				if f1.Type() == reflect.TypeOf(&LocalTime{}) {

					if f1.IsNil() && f1.IsNil() {
						return false // 如果都为nil，则相等
					} else if f1.IsNil() {
						return true // ti为nil，tj不为nil，tj较大
					} else if f2.IsNil() {
						return false // tj为nil，ti不为nil，ti较大
					}

					var t1 = f1.Elem().Interface().(LocalTime).Time()
					var t2 = f2.Elem().Interface().(LocalTime).Time()
					if !t1.Equal(t2) {
						if sortOrder == Ascending {
							return t1.Before(t2)
						} else {
							return t1.After(t2)
						}
					}
				}
			}
		}
		return false
	})
}
