package check

// ICheck 当json反序化到dto时，将调用ICheck接口
type ICheck interface {
	// Check 自定义检查
	Check()
}
