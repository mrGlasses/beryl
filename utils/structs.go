package utils

type File struct {
	FilePath         string `json:"filePath"`
	LastModification string `json:"lastModification"`
	Modified         bool   `json:"modified"`
	// Exists           bool   `json:"Exists"`
	NewFile  bool `json:"newFile"`
	Excluded bool `json:"excluded"`
}

type Project struct {
	Id               int    `json:"id"`
	ProjectName      string `json:"projectName"`
	Folder           string `json:"folder"`
	LastVerification string `json:"lastVerification"`
	Files            []File `json:"files"`
}

type FileStatus struct {
	ProjectName string
	Modified    int
	New         int
	Excluded    int
}

type ExternalVariables struct {
	Old []byte
	New []byte
}

type ConnString struct {
	DbsName  string
	User     string
	Password string
	Server   string
	Database string
	Port     string
}
