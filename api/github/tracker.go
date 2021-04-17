package tracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	//ReposAPI is the API URL for getting all repos owned by user
	ReposAPI = "https://api.github.com/users/{{.UserName}}/repos"

	// CommitsAPI is the API for fetching commits made by user for a particular repository
	CommitsAPI = "https://api.github.com/repos/:user/{repositoryName}/commits?author=:user"

	// AuthKey for github APi TODO: to be delted
	AuthKey = "token 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"
)

// GetAllReposOfUser gets all repositories for the user
func GetAllReposOfUser(user User) []Repo {
	var repos []Repo

	filledAPI := strings.Replace(ReposAPI, "{{.UserName}}", user.UserName, -1)
	fmt.Println(filledAPI)
	body, err := GetReponse(filledAPI)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(body, &repos)
	fmt.Printf("unmarshalled json repo %#v \n", repos)

	return repos

}

// GetAllCommitsForRepo get all commits for the repository name and owner
func GetAllCommitsForRepo(userName string, repos []Repo, whichWeek string) ([]byte, CommitTracker) {
	var commits []CommitInfo
	var weekCommits []CommitInfo
	var commitTracker CommitTracker
	var repoNames []string
	var currentWeek int

	for _, repo := range repos {

		tempStr := strings.ReplaceAll(CommitsAPI, ":user", userName)
		filledAPI := strings.ReplaceAll(tempStr, "{repositoryName}", repo.Name)
		fmt.Printf("\n API call for getting commits %s \n", filledAPI)

		body, err := GetReponse(filledAPI)
		if err == nil {
			//s, _ := strconv.Unquote(string(body))
			err = json.Unmarshal(body, &commits)
			if err != nil {
				fmt.Printf("error occured while unmarshalling %#v \n", err)
				//panic(err)
			}
			fmt.Printf("unmarshalled commit info %#v \n\n\n", commits)
			for _, commit := range commits {

				year, weekNum, date := GetCommitYearAndWeek(commit.Commit.Author.Date)
				commit.Commit.CommitDate = date
				var currentYear int
				currentYear, currentWeek = time.Now().ISOWeek()

				if strings.EqualFold(whichWeek, "current") {

				} else if strings.EqualFold(whichWeek, "previous") {
					currentWeek = currentWeek - 1
				}
				fmt.Printf("commit is matching for this week currentYear: [%d], currentWeek: [%d], Year: [%d], Week: [%d] \n\n", currentYear, currentWeek, year, weekNum)
				if currentYear == year && (currentWeek) == weekNum {

					commit.RepoName = repo.Name
					weekCommits = append(weekCommits, commit)
					repoNames = append(repoNames, repo.Name)
				}
			}
			fmt.Println("***************WEEk COMMITS START******************")
			fmt.Printf("REPO NAME IS %s \n", repo.Name)
			fmt.Printf("WeekCommits from json %#v \n\n\n", weekCommits)
			fmt.Println("***************WEEk COMMITS END******************")

		} else {
			fmt.Println(err)
		}
	}
	commitTracker.TotalCommits = len(weekCommits)
	commitTracker.UserName = userName
	commitTracker.WeekNum = currentWeek

	for _, weekCommit := range weekCommits {
		commitTracker.Email = weekCommit.Commit.Author.Email
		commitTracker.ProfilePic = weekCommit.Author.ProfilePic
		commitDetail := CommitDetail{}
		//commitDetail.Link = "https://github.com/" + commitTracker.UserName + "/" + weekCommit.RepoName + "/commit/" + weekCommit.Sha
		commitDetail.Link = weekCommit.Link
		commitDetail.CommitDate = weekCommit.Commit.CommitDate
		commitDetail.Sha = weekCommit.Sha
		commitDetail.RepoName = weekCommit.RepoName
		commitDetail.Message = weekCommit.Commit.Message
		commitTracker.CommitDetails = append(commitTracker.CommitDetails, commitDetail)
	}

	jsonData, err := json.MarshalIndent(commitTracker, "", "\t")
	if err != nil {
		fmt.Println("failed while converting data to json")
		panic(err)
	}
	return jsonData, commitTracker
}

// GetReponse returns response for a request
func GetReponse(filledAPI string) ([]byte, error) {
	//resp, err := http.Get(filledAPI)
	client := &http.Client{}
	req, err := http.NewRequest("GET", filledAPI, nil)
	req.Header.Set("Authorization", AuthKey)
	resp, err := client.Do(req)

	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Response received is %#v \n", string(body))
	return body, nil
}

// GetCommitYearAndWeek ... Convert RFC3339 date format to date and then to year and week number
func GetCommitYearAndWeek(dateString string) (int, int, time.Time) {
	layout := time.RFC3339
	commitDate, err := time.Parse(layout, dateString)

	if err != nil {
		fmt.Println("error while parsing date")
		panic(err)
	}
	y, w := commitDate.ISOWeek()

	return y, w, commitDate
}
