package post

// Post ...
type Post struct {
	ID            string
	Cover         string
	Title         string
	Author        string
	Username      string
	Created       int64
	Updated       int64
	Publish       int
	UserPic       string
	Content       string
	ReadingTime   int
	LikesCount    int64
	CommentsCount int64
}
