package err

import "fmt"

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
