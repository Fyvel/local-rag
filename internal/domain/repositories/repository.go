package repositories

type File struct {
	Path    string
	Name    string
	Type    string
	Content string
}

func NewFile(path, name, fileType, content string) (*File, error) {
	return &File{
		Path:    path,
		Name:    name,
		Type:    fileType,
		Content: content,
	}, nil
}

type Repository struct {
	URL   string
	Name  string
	Files []File
	SHA   string
}

func NewRepository(url, name, sha string, files []File) (*Repository, error) {
	return &Repository{
		URL:   url,
		Name:  name,
		SHA:   sha,
		Files: files,
	}, nil
}

func (r *Repository) FileCount() int {
	return len(r.Files)
}

func (r *Repository) HasFiles() bool {
	return len(r.Files) > 0
}
