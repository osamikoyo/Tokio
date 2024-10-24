package visual

import (
	"fmt"
)

func NewMessage(from string, user string, message string){
    colorReset := "\033[0m"
    colorCyan := "\033[36m"

    fmt.Println(colorReset + colorCyan + from + colorReset + " - " + message)
}