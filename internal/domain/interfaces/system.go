package interfaces

type System interface {
	MakeDir(folder string) string
	DownloadFile(url, filepath string, number int) error
	SaveURL(url, filepath string) error
}
