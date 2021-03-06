package Arrays

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func RemoveDuplicados(lista []string) []string {
	var temp []string

	for _, x := range lista {
		if !Contains(temp, x) {
			temp = append(temp, x)
		}
	}

	return temp
}

//Remove Duplicados e string vazia
func RemoveDuplicadosStringVazia(lista []string) ([]string, int) {
	var temp []string

	for _, x := range lista {
		if !Contains(temp, x) && len(x) > 0 {
			temp = append(temp, x)
		}
	}

	return temp, len(temp)
}

func UnionArray(input []string, out []string) []string {
	return append(input, out...)
}
