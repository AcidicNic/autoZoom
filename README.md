# Auto Zoom

[![Go Report Card](https://goreportcard.com/badge/github.com/acidicnic/autoZoom)](https://goreportcard.com/report/github.com/acidicnic/autoZoom)

Easily join your zoom classes and see your schedule.

[A video demo!](https://www.youtube.com/watch?v=J8eeF-AAvSw)


## Getting Started
```bash
git clone https://github.com/AcidicNic/autoZoom.git
cd autoZoom
```

- Setup your .schedule.json file before proceeding! (see [.schedule.json setup](#schedule-setup) below)

```bash
go run main.go
```


## Schedule Setup

.schedule.json is where all of your course data will be pulled from. You must set this up before using autoZoom!

It should be located in your home directory!

**_If you have any issues with your json file use this free online [JSON validator](https://jsonlint.com/)!_**


#### Example .schedule.json
```json
{
    "courses": [
        {
            "name": "Course Name",
            "days": "mw",
            "time": ["9:00AM", "11:30AM"],
            "attendCode": true,
            "autoZoom": true,
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

### How To Create Your Schedule File

- **"courses"**: *(list)* A list of JSON  objects containing the following:
    - **"name"**: *(string)* Whatever you'd like this course to be called.
    - **"days"**: *(string)* Any combination of "MTWRF", this is not case sensitive.
        - M = Monday
        - T = Tuesday
        - W = Wednesday
        - R = Thursday
        - F = Friday
    - **"time"**: *(list of 2 strings)* The first is the start time of the course, the second is the end time. This can be in 24hr format (13:30) or 12hr (1:30PM)
    - **"attendCode"**: *(bool)* Does this course use attendance codes? This is only for Make School students! (true/false)
    - **"autoZoom"**: *(bool)* Do you want zoom links to automatically open? (true/false)
    - **"zoom"**: *(string)* URL to the Zoom room for the course.
    - **"links"**: *(list)* [Optional] A list of links you want displayed when the course starts.
        - **"label"**: *(string)* Whatever you'd like this link to be called.
        - **"url"**: *(string)* URL for link.
