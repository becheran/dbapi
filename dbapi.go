// Deutsche Bahn API
package dbapi

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	timetableAPI = "https://api.deutschebahn.com/timetables/v1/"
)

type ErrRequestFailed struct {
	StatusCode   int
	ErrorMessage string
}

func (m *ErrRequestFailed) Error() string {
	return fmt.Sprintf("request failed (%d). %s", m.StatusCode, m.ErrorMessage)
}

type API struct {
	Bearer string
}

// StationInfo allows access to information about a station.
func (api *API) StationInfo(pattern string) ([]Station, error) {
	res, err := api.get("/station/" + pattern)
	if err != nil {
		return nil, err
	}

	var stations Stations
	err = xml.Unmarshal(res, &stations)
	return stations.Stations, err
}

// Plan returns planned data for the specified station within an hourly time slice
func (api *API) Plan(EvaNumber int, dateHour time.Time) (timetable Timetable, err error) {
	res, err := api.get(fmt.Sprintf("/plan/%d/%s/%d", EvaNumber, dateHour.Format("060102"), dateHour.Hour()))
	if err != nil {
		return
	}
	err = xml.Unmarshal(res, &timetable)
	return
}

func (api *API) get(path string) (resp []byte, err error) {
	req, _ := http.NewRequest(http.MethodGet, timetableAPI+path, nil)
	req.Header.Add("Accept", "application/xml")
	req.Header.Add("Authorization", "Bearer "+api.Bearer)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	resp, err = io.ReadAll(res.Body)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		return nil, &ErrRequestFailed{StatusCode: res.StatusCode, ErrorMessage: string(resp)}
	}

	return resp, nil
}
