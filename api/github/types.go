package tracker

import "time"

// User for storing all username related details
type User struct {
	UserName string
}

// Repo for storing repo details
type Repo struct {
	Name          string `json:"name"`
	Fork          bool   `json:"fork"`
	DefaultBranch string `json:"default_branch"`
}

// Struct for storing html page details
type Page struct {
	Title string
	Body  []byte
}

// CommitInfo each commit info
type CommitInfo struct {
	Sha      string `json:"sha"`
	Commit   Commit `json:"commit"`
	RepoName string `json:"-"`
	Link     string `json:"html_url"`
	Author   Author `json:"author"`
}

// Author specific info of author
type Author struct {
	ProfilePic string `json:"avatar_url"`
}

// Commit each CommitInfo details
type Commit struct {
	Message    string    `json:"message"`
	Author     Committer `json:"author"`
	CommitDate time.Time `jsone:"-"`
}

// Committer each commit author detailsdcx
type Committer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

//CommitTracker used to store commit related info
type CommitTracker struct {
	Email         string         `json:"email"`
	TotalCommits  int            `json:"total_commits"`
	CommitDetails []CommitDetail `json:"commit_links"`
	UserName      string         `json:"user_name"`
	WeekNum       int            `json:"week_num"`
	ProfilePic    string         `json:"profilepic"`
	Sno           int            `json:"s_no"`
}

//CommitDetail each commit detail
type CommitDetail struct {
	Link       string    `json:"commit_url"`
	CommitDate time.Time `json:"commit_date"`
	Sha        string    `json:"sha"`
	RepoName   string    `json:"repo_name"`
	Message    string    `json:"commit_message"`
}
