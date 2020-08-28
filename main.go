package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
	"os/exec"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/jedib0t/go-pretty/table"
)

func main() {
	cls()
	courses := getCoursesJson("schedule.json")
	s := gocron.NewScheduler()

	for _, course := range courses.Courses {
		time := hhToHH(course.Time[0])
		for i := 0; i < len(course.Days); i++ {
			switch strings.ToUpper(string(course.Days[i])) {
				case "M":
					s.Every(1).Monday().At(time).Do(course.classStarting, courses)
				case "T":
					s.Every(1).Tuesday().At(time).Do(course.classStarting, courses)
				case "W":
					s.Every(1).Wednesday().At(time).Do(course.classStarting, courses)
				case "R":
					s.Every(1).Thursday().At(time).Do(course.classStarting, courses)
				case "F":
					s.Every(1).Friday().At(time).Do(course.classStarting, courses)
			}
		}
	}
	printFullSchedule(courses)
	printDaySchedule(courses, time.Now().Weekday().String())
	<- s.Start()

}

type Courses struct {
    Courses []Course `json:"courses"`
}

type Course struct {
    Name string `json:"name"`
    Days string `json:"days"`
    Time []string `json:"time"`
    AttendCode bool `json:"attendCode"`
    AutoZoom bool `json:"autoZoom"`
	Zoom string `json:"zoom"`
    Links []Link `json:"links"`
}

type Link struct {
    Label string `json:"label"`
    Url  string `json:"url"`
}

func (course Course) classStarting(courses Courses) {
	if course.AutoZoom {
		openZoom(course)
	} else {
		askToOpenZoom(course)
	}

	if course.AttendCode {
		attend(course)
	}

	cls()
	printFullSchedule(courses)
	printDaySchedule(courses, time.Now().Weekday().String())
	if len(course.Links) > 0 {
		printLinks(course)
	}
}

func askToOpenZoom(course Course) {
	// Check if the user want to join the zoom call.
	fmt.Printf("Your %s course is starting! Want to open zoom link? (y/n): ",
	course.Name)
	var yn string
	fmt.Scanf("%s", &yn)
	if strings.ToLower(yn) == "y" || strings.ToLower(yn) == "yes" {
		// Open Zoom!
		openZoom(course)
	} else {
		// Display course name and Zoom URL in case they need it later.
		fmt.Printf("You declined to open zoom for %s: %s\n\n",
			course.Name, course.Zoom)
	}
}

func openZoom(course Course) {
	// Opens Zoom call URL for the given course.
	openZoom := exec.Command("open", course.Zoom)
	err := openZoom.Run()
	checkErr(err, fmt.Sprintf("Error while opening zoom link for %s:\n\t%s\n",
		course.Name, course.Zoom))
}

func attend(course Course) {
	// Asks for the attendance code and opens URL to check into the course.
	fmt.Print("\nWhat's the attendance code for today? (enter to skip): \n--> ")
	var code string
	fmt.Scanf("%s", &code)
	if strings.ToLower(code) != "n" && code != "" {
		attendClass := exec.Command("open", "https://www.makeschool.com/attend/" + code)
		err := attendClass.Run()
		checkErr(err, fmt.Sprintf("Error while opening: 'https://www.makeschool.com/attend/%s/'\n",
			code))
	}
}

func printLinks(course Course) {
	// Prints course.Links (if any)
	fmt.Printf("Links for %s:\n", course.Name)
	for _, link := range course.Links {
		fmt.Printf("\t%s: %s\n", link.Label, link.Url)
	}
	fmt.Println()
}

func prettyDays(days string) []string {
	// Tool to format days.
	// Ex. days="mwF" returns ["Mon", "Wed", "Friday"]
	dayKey := map[string] string {
	    "M": "Monday", "m": "Mon",
	    "T": "Tuesday", "t": "Tues",
	    "W": "Wednesday", "w": "Wed",
	    "R": "Thursday", "r": "Thur",
	    "F": "Friday", "f": "Fri",
	}
	var prettyDays []string
	for i := 0; i < len(days); i++ {
		prettyDays = append(prettyDays, dayKey[string(days[i])])
	}
	return prettyDays
}

func printFullSchedule(courses Courses) {
	// Prints full schedule table for every course in courses
	t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
    t.AppendHeader(table.Row{"Course", "Day(s)", "Time"})
	for _, course := range courses.Courses {
	    t.AppendRow([]interface{}{
			course.Name,
			strings.Join(prettyDays(strings.ToLower(course.Days)), ", "),
			strings.Join(course.Time, " - "),
		})
	}
	t.SetStyle(table.StyleLight)
    t.Render()
}

func printDaySchedule(courses Courses, day string) {
	// Prints schedule table for each course in courses occuring on day
	dayKey := map[string] string {
	    "Monday": "m",
	    "Tuesday": "t",
	    "Wednesday": "w",
	    "Thursday": "r",
	    "Friday": "f",
	}
	t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
	t.SetTitle(fmt.Sprintf("%s Schedule", day))
	for _, course := range courses.Courses {
		if strings.Contains(strings.ToLower(course.Days), dayKey[day]) {
			t.AppendRow([]interface{}{
				course.Name,
				strings.Join(course.Time, " - "),
			})
		}
	}
	t.SetStyle(table.StyleLight)
    t.Render()
}

func getCoursesJson(jsonFilename string) Courses {
	// Opens JSON file at jsonFilename
	scheduleJson, err := os.Open(jsonFilename)
    checkErr(err, fmt.Sprintf("Error while opening file: %s", jsonFilename))
	// Reads JSON data into bytes
	scheduleByte, err := ioutil.ReadAll(scheduleJson)
	checkErr(err, fmt.Sprintf("Error while reading file: %s", jsonFilename))
	// Outputs a Courses object
	// Courses is a stuct that holds an array of Course objects
	var courses Courses
	err = json.Unmarshal(scheduleByte, &courses)
	checkErr(err, fmt.Sprintf("Error while parsing JSON file: %s", jsonFilename))
	return courses
}

func hhToHH(time string) string {
	// Converts/validates the time string returns a time string in 24hr format.
	i := strings.Index(time, ":")
	if i == 1 {
		// 4:20PM --> 04:20PM
		time = "0" + time
	}
	if strings.ToUpper(time[len(time)-2:]) == "AM" {
		// handles AM
		if time[:2] == "12" {
			// if 12:??AM --> 00:??
			return "00" + time[2:len(time)-2]
		}
		return time[:len(time)-2]
	}
	if strings.ToUpper(time[len(time)-2:]) == "PM" {
		// Handles PM by adding 12 hrs and removing "PM"
		hr, err := strconv.Atoi(time[:2])
		checkErr(err, fmt.Sprintf("Error while reading time: %s", time))
		return strconv.Itoa(hr+12) + time[2:len(time)-2]
	}
	// time="4:20PM" returns
	return time
}

func cls() {
	// Clears the terminal screen (This only works on unix systems!)
	cls := exec.Command("clear")
	cls.Stdout = os.Stdout
	err := cls.Run()
	checkErr(err, "Error while running cmd: 'clear'\n")
}

func checkErr(err error, msg string) {
	// Standard check error func
	if err != nil {
		fmt.Println(msg)
	    panic(err)
	}
}
