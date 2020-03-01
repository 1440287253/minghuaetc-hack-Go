package main

import (
	UserUtils "Minghuaetc/api"
	"fmt"
)

func main() {
	var answer string
	fmt.Printf("本工具仅供交流与学习使用 禁止商业和违法用途 使用此工具造成的一切后果由本人承担 是否同意(y同意/n不同意)? ")
	fmt.Scanf("%s\n", &answer)
	switch answer {
	case "y":
		UserUtils.GetSchool()
		UserUtils.Login()
	case "n":
		fmt.Errorf("bye bye")
	default:
		fmt.Errorf("bye bye")
	}
}
