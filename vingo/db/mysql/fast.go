package mysql

// Get 通过主键id获取记录
func Get[T any, I string | int | uint](id I) (data T) {
	NotExistsErr(&data, "id=?", id)
	return
}

// GetByColumn 通过条件获取记录
func GetByColumn[T any](condition ...any) (data T) {
	NotExistsErr(&data, condition...)
	return
}
