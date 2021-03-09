package main

//https://stackoverflow.com/questions/61011873/cant-add-new-cobra-cli-command-when-the-file-is-inside-a-folder
import (
	"golang-base/cmd"
	_ "golang-base/cmd/aliyun"
)

func main() {
	cmd.Execute()
}
