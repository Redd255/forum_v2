package server

type comment struct {
	Id       int
	Username string
	Content  string
	Creation string
	Like     int
	Dislike  int
	PostId int
}

type post struct {
	Id       int
	Username string
	Content  string
	Creation string
	Topic    string
	Like     int
	Dislike  int
	Commentcount int
}

type data struct {
	Username string
	Posts    []post
	Logout   bool
}

type Data struct {
	Username string
	PostId   int
	Comments []comment
	Post     post
	Logout   bool
	Error    bool
}
type Rreaction struct {
	PostId   string `json:"postId"`
	Reaction string `json:"reaction"`
}

type Creation struct {
	CommentId string `json:"commentId"`
	Rreaction string `json:"reaction"`
}
