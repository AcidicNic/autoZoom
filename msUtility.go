package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/jedib0t/go-pretty/table"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
	"os/exec"
	"time"
)

func main() {
	courses := getCoursesJson("schedule.json")
	s := gocron.NewScheduler()

	for _, course := range courses.Courses {
		time := hhToHH(course.Time[0])
		for i := 0; i < len(course.Days); i++ {
			switch strings.ToUpper(string(course.Days[i])) {
				case "M":
					s.Every(1).Monday().At(time).Do(course.classStarting)
				case "T":
					s.Every(1).Tuesday().At(time).Do(course.classStarting)
				case "W":
					s.Every(1).Wednesday().At(time).Do(course.classStarting)
				case "R":
					s.Every(1).Thursday().At(time).Do(course.classStarting)
				case "F":
					s.Every(1).Friday().At(time).Do(course.classStarting)
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
	Zoom string `json:"zoom"`
    Links []Link `json:"links"`
}

type Link struct {
    Label string `json:"label"`
    Url  string `json:"url"`
}

func (course Course) classStarting() {
	openZoom(course)
	if course.AttendCode {
		attend(course)
	}
	printLinks(course)
}

func openZoom(course Course) {
	openZoom := exec.Command("open", course.Zoom)

	fmt.Printf("Your %s course is starting! Want to open zoom link? (y/n): ",
	course.Name)
    var yn string
    fmt.Scanf("%s", &yn)
	if strings.ToLower(yn) == "y" || strings.ToLower(yn) == "yes" {
		err := openZoom.Run()
		checkErr(err, fmt.Sprintf("Error while opening zoom link for %s:\n\t%s\n",
			course.Name, course.Zoom))
	} else {
		fmt.Printf("You declined to open zoom for %s: %s\n\n",
			course.Name, course.Zoom)
	}
}

func attend(course Course) {
	fmt.Print("\nWhat's the attendance code for today? (enter to skip): \n--> ")
	var code string
	fmt.Scanf("%s", &code)
	if strings.ToLower(code) != "n" && code != "" {
		attendClass := exec.Command("open", "https://make.sc/attend/" + code)
		err := attendClass.Run()
		checkErr(err, fmt.Sprintf("Error while opening: 'https://make.sc/attend/%s/'\n",
			code))
	}
}

func printLinks(course Course) {
	fmt.Printf("Links for %s:\n", course.Name)
	for _, link := range course.Links {
		fmt.Printf("\t%s: %s\n", link.Label, link.Url)
	}
}

func prettyDays(days string) []string {
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
	scheduleJson, err := os.Open(jsonFilename)
    checkErr(err, fmt.Sprintf("Error while opening file: %s", jsonFilename))

	scheduleByte, err := ioutil.ReadAll(scheduleJson)
	checkErr(err, fmt.Sprintf("Error while reading file: %s", jsonFilename))

	var courses Courses
	err = json.Unmarshal(scheduleByte, &courses)
	checkErr(err, fmt.Sprintf("Error while parsing JSON file: %s", jsonFilename))
	return courses
}

func hhToHH(time string) string {
	// 4:20PM --> 16:20
	i := strings.Index(time, ":")
	if i == 1 {
		time = "0" + time
	}
	if strings.ToUpper(time[len(time)-2:]) == "AM" {
		if time[:2] == "12" {
			return "00" + time[2:len(time)-2]
		}
		return time[:len(time)-2]
	}
	if strings.ToUpper(time[len(time)-2:]) == "PM" {
		hr, err := strconv.Atoi(time[:2])
		checkErr(err, fmt.Sprintf("Error while reading time: %s", time))
		return strconv.Itoa(hr+12) + time[2:len(time)-2]
	}
	return time
}

func checkErr(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
	    panic(err)
	}
}
