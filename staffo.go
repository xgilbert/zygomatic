package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// TimeIn : returns the current time in a given location such as Europe/London
func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

// buildReqURL : helper function to build request url from configuration (conf struct)
func buildReqURL(s string) string {
	return fmt.Sprintf("%s%s/", conf.Hostname, conf.Subdomain) + s
}

// get : wrapper for http requests
func get(s string) []byte {
	reqURL := buildReqURL(s)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.SetBasicAuth(conf.Username, conf.Password)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

// getTimeZone : get current time in a given time zone
func getTimeZone() {
	body := get("account.json")

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	// tz : string, e.g. "Europe/London"
	tz := fmt.Sprintf("%v", result["time_zone"])

	t, err := TimeIn(time.Now(), tz)
	if err == nil {
		fmt.Println(t.Location(), t.Format("15:04"))
	} else {
		fmt.Println(tz, "<time unknown>")
	}
}

// getLocationByID : return a location (top tree) for a given locationID
func getLocationByID(locationID int) {
	body := get("locations/" + strconv.Itoa(locationID) + ".json")

	log.Println(string(body))
}

// getDepartments : return all departments (for all locations)
func getDepartments() {
	body := get("departments.json")

	log.Println(string(body))
}

// getDepartmentByID : return department details for a given department ID
func getDepartmentByID(departmentID int) {
	body := get("departments/" + strconv.Itoa(departmentID) + ".json")
	// fmt.Println(buildReqURL("departments/" + strconv.Itoa(departmentID) + ".json"))

	log.Println(string(body))
}

// @FIXME
// func getDepartmentByIDByLocationID(locationID int, departmentID int) {
// 	body := get("locations/" + strconv.Itoa(locationID) + "/departments/" + strconv.Itoa(departmentID) + ".json")
//
// 	log.Println(string(body))
// }

// getSchedulesByLocationID : return all schedules for a given location ID for the whole time range
func getSchedulesByLocationID(locationID int) {
	body := get("locations/" + strconv.Itoa(locationID) + "/" + "schedules.json")

	log.Println(string(body))
}

// Get schedules for a given location (by ID) and given time range
// The body is an array of json objects; the important fields are:
// "id":662108 -schedule ID
// "location_id":37759 - location id
// "bop":"2020-07-27T00:00:00+02:00" - beginning of rota
// "eop":"2020-08-02T23:59:59+02:00" - end of rota
func getSchedulesIDs(locationID int, startDate string, endDate string) []string {
	// format string to YYYY-MM-DDTHH:MM:SS+01:00
	startDateTimeTz := fmt.Sprintf("%sT00:00:00+01:00", startDate)
	endDateTimeTz := fmt.Sprintf("%sT23:59:59+01:00", endDate)

	reqURL := buildReqURL("locations/" + strconv.Itoa(locationID) + "/" + "schedules.json")

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.SetBasicAuth(conf.Username, conf.Password)

	q := req.URL.Query()
	q.Add("from", startDateTimeTz)
	q.Add("until", endDateTimeTz)
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var result []map[string]interface{}
	json.Unmarshal(body, &result)

	var schedulesIDs []string

	for _, res := range result {
		schedulesIDs = append(schedulesIDs, fmt.Sprintf("%v", res["id"]))
	}

	return schedulesIDs
}

// TODO: add pagination
// getShiftsByScheduleID : return all shifts for a given schedule ID
func getShiftsByScheduleID(scheduleID int) {
	body := get("schedules/" + strconv.Itoa(scheduleID) + "/" + "shifts.json")

	fmt.Println(string(body))
}

// TODO: make an array version for deparments
// getShiftsByDepartmentID :
func getShiftsByDepartmentID(departmentID int) {

	reqURL := buildReqURL("shifts.json")

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.SetBasicAuth(conf.Username, conf.Password)

	q := req.URL.Query()
	q.Add("department_ids[]", strconv.Itoa(departmentID))
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

	// var result []map[string]interface{}
	// json.Unmarshal(body, &result)

	// var schedulesIDs []string

	// for _, res := range result {
	// 	schedulesIDs = append(schedulesIDs, fmt.Sprintf("%v", res["id"]))
	// }

	// return schedulesIDs
}
