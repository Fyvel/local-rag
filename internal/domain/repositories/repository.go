package repositories

type File struct {
	Path    string
	Name    string
	Type    string
	Content string
}

type Repository struct {
	URL   string
	Name  string
	Files []File
	SHA   string
}

func (r *Repository) FileCount() int {
	return len(r.Files)
}
