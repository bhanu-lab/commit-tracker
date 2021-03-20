package main

import (
	tracker "commit-tracker/api/github"
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	// Read All User Names for whom github commits has to be tracked
	// to-do
	whichWeek := os.Args[1]
	var users []tracker.User

	file, err := os.Open("users.csv")
	if err != nil {
		fmt.Println("error occured while reading csv file")
		panic(err)
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()
	var userNames []string
	for _, userName := range records {
		userNames = append(userNames, userName[0])
		users = append(users, tracker.User{UserName: userName[0]})
	}
	var allUsers []string
	// Get All  Repos for which user is the owner
	//user := tracker.User{UserName: userNames}
	//user := tracker.User{UserName: "bhanu-lab"}
	//user := tracker.User{UserName: "RajeshReddyG"}
	//user := tracker.User{UserName: "Gundupalli"}

	for _, user := range users {
		repos := tracker.GetAllReposOfUser(user)
		//repos := []tracker.Repo{tracker.Repo{Name: "NetworkScanner"}}

		// for each repo created by the user check for any new commits
		//for _, repo := range repos {
		//if repo.Name == "NetworkScanner" {
		allCommits := tracker.GetAllCommitsForRepo(user.UserName, repos, whichWeek)

		//fmt.Printf("All Commits for this Week %s \n", string(allCommits))
		allUsers = append(allUsers, string(allCommits))
	}
	//}
	//}
	fmt.Println("***************ALL USERS***********************")
	for _, user := range allUsers {
		fmt.Printf("All Commits for this Week %s \n", user)
	}
	fmt.Println("************************DONE*******************")

}
