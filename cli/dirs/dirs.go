package dirs

import "github.com/arelate/vangogh_local_data"

var (
	tempDir = ""
)

func GetStateDir() string {
	return vangogh_local_data.Pwd()
}

func SetStateDir(d string) {
	vangogh_local_data.ChRoot(d)
}

func GetTempDir() string {
	return tempDir
}

func SetTempDir(d string) {
	tempDir = d
}
