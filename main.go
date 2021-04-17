package main

import (
	tracker "commit-tracker/api/github"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/week/{week}", WeeklyCommits).Methods("GET")
	http.ListenAndServe(":8080", r)
}

// WeeklyCommits takes for which week number to scan commits
func WeeklyCommits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	whichWeek := vars["week"]
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
	var commits []tracker.CommitTracker

	for _, user := range users {
		repos := tracker.GetAllReposOfUser(user)
		//repos := []tracker.Repo{tracker.Repo{Name: "NetworkScanner"}}

		// for each repo created by the user check for any new commits
		//for _, repo := range repos {
		//if repo.Name == "NetworkScanner" {
		allCommits, commitTracker := tracker.GetAllCommitsForRepo(user.UserName, repos, whichWeek)

		//fmt.Printf("All Commits for this Week %s \n", string(allCommits))
		allUsers = append(allUsers, string(allCommits))
		commits = append(commits, commitTracker)
	}
	fmt.Println("***************ALL USERS***********************")
	/*for _, user := range allUsers {
		fmt.Printf("All Commits for this Week %s \n", user)
	}*/
	fmt.Println("************************DONE*******************")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allUsers)
	template, _ := template.ParseFiles("tracker.html")
	template.Execute(os.Stdout, commits)
}
