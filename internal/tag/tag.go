package tag


type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
}

func NewTag(name string) *Tag {
	return &Tag{
		Name:      name,
	}
}
