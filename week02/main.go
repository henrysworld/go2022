package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func main() {
	if _, err := os.Open("non-existing"); err != nil {
		var pathError *fs.PathError

		err2 := fmt.Errorf("error2: [%w]", err)
		err3 := fmt.Errorf("error3: [%w]", err2)
		e := errors.As(err2, &pathError)
		e1 := errors.Is(err3, pathError)
		fmt.Println(e)
		fmt.Println(e1)
		if errors.As(err, &pathError) {
			fmt.Println("Failed at path:", pathError.Path)
		} else {
			fmt.Println(err)
		}
	}

}
