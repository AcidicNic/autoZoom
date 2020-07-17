# MS Utility

[![Go Report Card](https://goreportcard.com/badge/github.com/acidicnic/msutility)](https://goreportcard.com/report/github.com/acidicnic/msutility)

Easily join your zoom classes and see your schedule.


## Getting Started
```bash
git clone https://github.com/AcidicNic/msutility.git
cd msutility
```

- Setup your schedule.json file before proceeding! (see [schedule.json setup](#schedule-setup) below)

```bash
go run msutility.go
```


## Schedule Setup

schedule.json is where all of your course data will be pulled from. You must set this up before using msUtility!

**_If you have any issues with your json file use this free online [JSON validator](https://jsonlint.com/)!_**


#### Basic Structure
```json
{
    "courses": [
        {
            "name": "Course Name",
            "days": "mw",
            "time": ["9:00AM", "11:30AM"],
            "attendCode": true,
            "zoom": "https://URL-TO-ZOOM/",
            "links": [
                {
                    "label": "some class link",
                    "url": "https://your-class-related-URL/"
                },
                {
                    "label": "another class link",
                    "url": "https://your-class-related-URL/"
                }
            ]
        },
        {
            "...": "..."
        }
    ]
}
```
- **"courses"**: *(list)* A list of JSON  objects containing the following:
    - **"name"**: *(string)* Whatever you'd like this course to be called.
    - **"days"**: *(string)* Any combination of "MTWRF", case doesn't matter.
        - M = Monday
        - T = Tuesday
        - W = Wednesday
        - R = Thursday
        - F = Friday
    - **"time"**: *(list of 2 strings)* The first is the start time of the course, the second is the end time.
    - **"attendCode"**: *(bool)* Does this course use attendance codes? (true/false)
    - **"zoom"**: *(string)* URL to the Zoom room for the course.
    - **"links"**: *(list)* [Optional] A list of links you want displayed when the course starts.
        - **"label"**: *(string)* Whatever you'd like this link to be called.
        - **"url"**: *(string)* URL for link.
