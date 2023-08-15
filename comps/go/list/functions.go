package list

// ArrayStringToArrayInterface convert a array of string into array of interface
func ArrayStringToArrayInterface(stringArray []string) []interface{} {
	var interfaceArray []interface{}

	for _, value := range stringArray {
		interfaceArray = append(interfaceArray, value)
	}
	return interfaceArray
}

// ArrayPointStringToArrayString convert a pointer string's array to normal string's array
func ArrayStringPointToArrayString(stringArray []*string) []string {
	var result []string

	for _, value := range stringArray {
		result = append(result, *value)
	}

	return result
}

// ContainsString check if a string contains in array of string
func ContainsString(array []string, str string) bool {
	for _, v := range array {
		if v == str {
			return true
		}
	}
	return false
}

// RemoveIndex remove um item no array/slice baseado no index
func RemoveIndex(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}

// RemoveIndexFromInterface remove um item no array/slice baseado no index cujo parametro é uma interface
func RemoveIndexFromInterface(s []interface{}, index int) []interface{} {
	return append(s[:index], s[index+1:]...)
}

/*
ArrayStringReverse Retorna um array de strings com os mesmos dados do array passado como parametro mas com a ordem
invertida
*/
func ArrayStringReverse(s []string) (ret []string) {
	for i := len(s); i > 0; i-- {
		ret = append(ret, s[i-1])
	}
	return
}

// RemoveStringItem remove um item de um array de string
func RemoveStringItem(items []string, item string) (newItems []string) {
	for _, i := range items {
		if i != item {
			newItems = append(newItems, i)
		}
	}
	return newItems
}

// Difference retorna a diferenca entre duas listas de string
func Difference(slice1 []string, slice2 []string) []string {
	var diffStr []string
	m := map[string]int{}

	for _, s1Val := range slice1 {
		m[s1Val] = 1
	}
	for _, s2Val := range slice2 {
		m[s2Val] = m[s2Val] + 1
	}

	for mKey, mVal := range m {
		if mVal == 1 {
			diffStr = append(diffStr, mKey)
		}
	}

	return diffStr
}

// RemoveDuplicate remove items duplicado de um slice
func RemoveDuplicate[T string | int | float64](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
