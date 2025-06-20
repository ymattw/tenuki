package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func DP(format string, args ...any) {
	f, err := os.OpenFile("/tmp/gote.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("failed to open /tmp/gote.log")
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf(time.Now().Format(time.RFC3339)+" "+format, args...) + "\n")
}

func formatObject(obj any) string {
	var out bytes.Buffer
	data, _ := json.Marshal(obj)
	if json.Indent(&out, []byte(data), "", "  ") != nil {
		return ""
	}
	return out.String()
}
