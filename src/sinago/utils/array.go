package utils

import "sort"

func ArrayDiff(array1 []string, array2 []string) (diff []string) {
	if len(array1) < len(array2) {
		array1, array2 = array2, array1
	}
	for _, value := range array1 {
		pos := ArrayFind(value, array2)
		if pos == -1 {
			diff = append(diff, value)
		}
	}
	return diff
}

func ArrayPop(array []string) (elem string, arrayPoped []string) {
	if len(array) <= 0 {
		return "", arrayPoped
	}
	if len(array) == 1 {
		return array[0], arrayPoped
	}
	elem = array[len(array) - 1]
	arrayPoped = array[:len(array) -2]
	return elem, arrayPoped
}

func ArrayFind(target string, slice []string) (pos int) {
	for index, value := range slice {
		if target == value {
			return index
		}
	}
	return -1
}


func RemoveDupAndEmpty(a []string) (ret []string){
	sort.Strings(a)
	if len(a) == 1 {
		return a
	}
	a_len := len(a)
	for i:=0; i < a_len; i++{
		if (i > 0 && a[i-1] == a[i]) || len(a[i])==0{
			continue
		}
		arrValue := a[i]
		ret = append(ret, arrValue)
	}
	return
}


