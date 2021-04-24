package main

import (
	tracker "commit-tracker/api/github"
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/week/{week}", WeeklyCommits).Methods("GET")
	http.ListenAndServe(":8080", r)
}

// WeeklyCommits takes for which week number to scan commits
func WeeklyCommits(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
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

	for i, user := range users {
		wg.Add(1)
		go CommitTrackerWorker(user, whichWeek, i, allUsers, &commits, &wg)
	}
	wg.Wait()
	SortBasedOnTotalCommits(commits)
	template, _ := template.ParseFiles("html/tracker.html")
	template.Execute(w, commits)
}

// CommitTrackerWorker worker function to spawn go routines
func CommitTrackerWorker(user tracker.User, whichWeek string, i int, allUsers []string, commits *[]tracker.CommitTracker, wg *sync.WaitGroup) {
	defer wg.Done()
	repos := tracker.GetAllReposOfUser(user)

	// for each repo created by the user check for any new commits
	allCommits, commitTracker := tracker.GetAllCommitsForRepo(user.UserName, repos, whichWeek)
	commitTracker.Sno = i
	allUsers = append(allUsers, string(allCommits))
	*commits = append(*commits, commitTracker)
}

// CreateFile creates file and returns file pointer and error if any
func CreateFile(fileName string) (*os.File, error) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Println("create file: ", err)
		return nil, err
	}
	return f, nil
}

// SortBasedOnTotalCommits sorts based on total number of commits
func SortBasedOnTotalCommits(commits []tracker.CommitTracker) {
	sort.SliceStable(commits, func(i, j int) bool {
		return commits[i].TotalCommits > commits[j].TotalCommits
	})
}
