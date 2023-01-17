package main

func remove(slice []string, need any) []string {
	for index, value := range slice {
		if value == need {
			return append(slice[:index], slice[index+1:]...)
		}
	}

	return nil
}

func find(array []string, need any) bool {
	for _, value := range array {
		if value == need {
			return true
		}
	}

	return false
}
