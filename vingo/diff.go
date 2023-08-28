package vingo

import (
	"fmt"
	"reflect"
)

type DiffBox struct {
	Old    any
	New    any
	Result *map[string]DiffItem
}

type DiffItem struct {
	Column   string
	OldValue any
	NewValue any
	Message  string
}

func (s *DiffItem) SetMessage() {
	s.Message = fmt.Sprintf("将%v的值[%v]变更为[%v]；", s.Column, s.OldValue, s.NewValue)
}

// 比较
func (s *DiffBox) Compare() {
	result := map[string]DiffItem{}
	oldVal := reflect.ValueOf(s.Old)
	newVal := reflect.ValueOf(s.New)
	if oldVal.Kind() != reflect.Struct || newVal.Kind() != reflect.Struct {
		return
	}
	oldType := oldVal.Type()
	newType := newVal.Type()
	if oldType != newType {
		return
	}
	for i := 0; i < oldVal.NumField(); i++ {
		oldField := oldVal.Field(i)
		newField := newVal.Field(i)

		if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
			name := oldType.Field(i).Name
			if IsInSlice(name, []string{"CreatedAt", "UpdatedAt", "DeletedAt"}) {
				continue
			}
			diffItem := DiffItem{
				Column:   name,
				OldValue: oldField.Interface(),
				NewValue: newField.Interface(),
			}
			diffItem.SetMessage()
			result[diffItem.Column] = diffItem
		}
	}
	s.Result = &result
}

// 判断指定字段是否被修改，被修改返回true
func (s *DiffBox) IsChange(column string) (DiffItem, bool) {
	if s.Result == nil {
		s.Compare()
	}
	item, ok := (*s.Result)[column]
	return item, ok
}
