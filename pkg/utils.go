package getcode

import "fmt"

func ErrorCheck(err error, message string) {

	if err != nil {
		fmt.Println(message)
		panic(err)
	}
}

func isMapOfListKey(key string, mapOfList map[string][]string) bool {

	_, ok := mapOfList[key]
	if ok {
		return true
	} else {
		return false
	}
}

func isValueinList(list []string, target string) bool {

	for _, item := range list {
		if target == item {
			return true
		}
	}
	return false
}
