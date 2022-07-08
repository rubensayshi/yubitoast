package notifier

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed imgs/gnupg.png
var icon []byte

type Notifier interface {
	Run()
}

type PopupEvent struct {
	Toggle bool
}

func createTmpIcon() string {
	// create tmpfile holding the icon
	iconFile, err := os.CreateTemp(os.TempDir(), "yubitoast-icon-*.png")
	if err != nil {
		fmt.Printf("Failed to create icon in tmpdir: %v \n", err)
	}
	iconPath := iconFile.Name()

	_, err = iconFile.Write(icon)
	if err != nil {
		fmt.Printf("Failed to create icon in tmpdir: %v \n", err)
	}

	return iconPath
}
