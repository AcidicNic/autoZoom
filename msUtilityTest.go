package main

import "testing"

func TesthhToHH(t *testing.T) {
    timeMap := map[string] string {
        "1:00AM": "01:00",
        "1:00PM": "13:00",
        "12:00PM": "12:00",
        "12:00AM": "00:00",
        "4:20PM": "16:20",
        "4:20AM": "04:20",
        "8:34": "08:34",
        "23:45": "23:45",
	}
    for k, v := range timeMap {
        time24hr := hhToHH("k")
        if time24hr != v {
           t.Errorf("Time was incorrect, got: %d, want: %d.", time24hr, v)
        }
    }

}
