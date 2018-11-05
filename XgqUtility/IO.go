package XgqUtility

import "os"

type Utility struct {
}

func (u *Utility) FileExist(path string) bool {
	stat, err := os.Stat(path)
	if err == nil && !stat.IsDir() {
		return true
	}

	return false
}

func (u *Utility) DirectoryExist(path string) bool {
	stat, err := os.Stat(path)
	if err == nil && stat.IsDir() {
		return true
	}

	return false
}
