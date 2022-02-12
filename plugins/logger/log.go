package logger

import (
	"fmt"
	"time"
)

func WriteText(level int, text string) {
	fmt.Printf("[%s] (%s) * %s\n", time.Now(), GetLevelStr(level), text)
}

func WriteError(level int, text string, err error) {
	fmt.Printf("[%s] (%s) * %s:%v\n", time.Now(), GetLevelStr(level), text, err)
}
