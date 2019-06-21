package post

// Post ...
type Post struct {
	Username      string
	ID            string
	Title         string
	Created       int64
	Updated       int64
	Publish       int
	UserPic       string
	Content       string
	LikesCount    int64
	CommentsCount int64
}
