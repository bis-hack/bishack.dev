package post

// Post ...
type Post struct {
	ID            string
	Cover         string
	Title         string
	Username      string
	Created       int64
	Updated       int64
	Publish       int
	UserPic       string
	Content       string
	LikesCount    int64
	CommentsCount int64
}
