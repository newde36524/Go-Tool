package arraytool

//RevertArray 反转数组
func RevertArray(arr []interface{}) []interface{} {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

//CopySlice 数组拷贝
func CopySlice(bs []interface{}) []interface{} {
	result := append(bs, struct{}{})
	return result[:len(result)-1]
}
