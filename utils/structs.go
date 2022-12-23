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
	Id          int    `json:"id"`
	ProjectName string `json:"project_name"`
	Folder      string `json:"folder"`
	Files       []File `json:"files"`
}

type FileStatus struct {
	ProjectName string
	Modified    int
	New         int
	Excluded    int
}
