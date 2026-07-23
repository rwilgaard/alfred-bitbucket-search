package util

import (
	"fmt"
	"os"

	aw "github.com/deanishe/awgo"
)

func GetIcon(name string) *aw.Icon {
	iconPath := fmt.Sprintf("icons/%s.png", name)
	if _, err := os.Stat(iconPath); err == nil {
		return &aw.Icon{Value: iconPath}
	}
	return aw.IconWorkflow
}
