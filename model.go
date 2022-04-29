package dbapi

import (
	"encoding/xml"
	"time"
)

// The time, in ten digit 'YYMMddHHmm' format.
// e.g. '1404011437' for 14:37 on April the 1st of 2014.
type Time struct {
	time.Time
}

func (t *Time) UnmarshalXMLAttr(attr xml.Attr) error {
	// Format: 21-12-22 10:10:36.633
	res, err := time.Parse("0601021504", attr.Value)
	if err != nil {
		return err
	}
	*t = Time{res}
	return nil
}

// Format: 21-12-22 10:10:36.633
type DateTime struct {
	time.Time
}

func (dt *DateTime) UnmarshalXMLAttr(attr xml.Attr) error {
	res, err := time.Parse("06-01-02 15:04:05.000", attr.Value)
	if err != nil {
		return err
	}
	*dt = DateTime{res}
	return nil
}

type Stations struct {
	XMLName  xml.Name  `xml:"stations"`
	Stations []Station `xml:"station"`
}

type Station struct {
	XMLName xml.Name `xml:"station"`
	// List of platforms. A sequence of platforms
	// separated by the pipe symbols (“|”).
	// Optional.
	Platforms string `xml:"p,attr"`
	// List of meta stations. A sequence of station
	// names separated by the pipe symbols (“|”).
	// Optional.
	Meta string `xml:"meta,attr"`
	// Station name.
	Name string `xml:"name,attr"`
	// EVA station number.
	EvaNumber string `xml:"eva,attr"`
	// DS100 station code.
	Ds100 string `xml:"ds100,attr"`
	// Flag for Stations of Deutsche Bahn. Optional
	DB *bool `xml:"db,attr"`
	// Creation Time, Optional
	Creation DateTime `xml:"creationts,attr"`
	// Creation Time, Optional
	Updates DateTime `xml:"updatets,attr"`
}

type Timetable struct {
	XMLName   xml.Name        `xml:"timetable"`
	Station   string          `xml:"station,attr"`
	EvaNumber int             `xml:"eva,attr"`
	Stops     []TimetableStop `xml:"s"`
	Messages  []Message       `xml:"m"`
}

type TimetableStop struct {
	XMLName                xml.Name                 `xml:"s"`
	ID                     string                   `xml:"id,attr"`
	EvaNumber              int                      `xml:"eva,attr"`
	TripLabel              TripLabel                `xml:"tl"`
	TripReference          TripReference            `xml:"ref"`
	Arrival                Event                    `xml:"ar"`
	Departure              Event                    `xml:"dp"`
	Message                Message                  `xml:"m"`
	HistoricDelay          []HistoricDelay          `xml:"hd"`
	HistoricPlatformChange []HistoricPlatformChange `xml:"hpc"`
	Connection             []Connection             `xml:"conn"`
	ReferenceTripRelation  []ReferenceTripRelation  `xml:"rtr"`
}

type Message struct {
	XMLName            xml.Name             `xml:"m"`
	ID                 string               `xml:"id,attr"`
	Type               MessageType          `xml:"t,attr"`
	From               Time                 `xml:"from,attr"`
	To                 Time                 `xml:"to,attr"`
	Code               int                  `xml:"c,attr"`
	InternalText       string               `xml:"int,attr"`
	ExternalText       string               `xml:"ext,attr"`
	Category           string               `xml:"cat,attr"`
	ExternalCategory   string               `xml:"ec,attr"`
	Timestamp          Time                 `xml:"ts,attr"`
	Priority           Priority             `xml:"pr,attr"`
	Owner              string               `xml:"o,attr"`
	ExternalLink       string               `xml:"elnk,attr"`
	Deleted            int                  `xml:"del,attr"`
	DistributorMessage []DistributorMessage `xml:"dm"`
	TripLabel          []TripLabel          `xml:"tl"`
}

type TripLabel struct {
	XMLName     xml.Name `xml:"tl"`
	FilterFlags string   `xml:"f,attr"`
	TripType    TripType `xml:"t,attr"`
	// A unique short-form and only intended to map a trip to specific evu.
	Owner             string `xml:"o,attr"`
	TripOrTrainNumber string `xml:"n,attr"`
	Category          string `xml:"c,attr"`
}

type TripReference struct {
	XMLName       xml.Name    `xml:"ref"`
	TripLabel     TripLabel   `xml:"tl"`
	TripReference []TripLabel `xml:"tr"`
}

type Event struct {
	PlannedPath     string      `xml:"ppth,attr"`
	ChangedPath     string      `xml:"cpth,attr"`
	PlannedPlatform string      `xml:"pp,attr"`
	ChangedPlatform string      `xml:"cp,attr"`
	PlannedTime     Time        `xml:"pt,attr"`
	ChangedTime     Time        `xml:"ct,attr"`
	PlanedStatus    EventStatus `xml:"ps,attr"`
	/*
			Changed status. The status of this event, a one-character indicator that is one of:

		'a' = this event was added
		'c' = this event was cancelled
		'p' = this event was planned (also used when the cancellation of an event has been revoked)
			  The status applies to the event, not to the trip as a whole. Insertion or removal of a single stop
			  will usually affect two events at once: one arrival and one departure event. Note that these two events do
			  not have to belong to the same stop. For example, removing the last stop of a trip will result in arrival
			  cancellation for the last stop and of departure cancellation for the stop before the last.
			  So asymmetric cancellations of just arrival or departure for a stop can occur.
	*/
	ChangedStatus EventStatus `xml:"cs,attr"`
	// 1 if the event should not be shown on WBT because travellers are not supposed to enter or exit the train at this stop.
	Hidden           int  `xml:"hi,attr"`
	CancellationTime Time `xml:"clt,attr"`
	// Wing. A sequence of trip id separated by the pipe symbols ('|'). E.g. '-906407760000782942-1403311431'.
	Wings string `xml:"wings,attr"`
	// Transition. Trip id of the next or previous train of a shared train. At the start stop this references the previous trip, at the last stop it references the next trip. E.g. '2016448009055686515-1403311438-1'
	Transitions            string    `xml:"tra,attr"`
	PlannedDistantEndpoint string    `xml:"pde,attr"`
	ChangedDistantEndpoint string    `xml:"cde,attr"`
	DistantChanged         int       `xml:"dc,attr"`
	Line                   string    `xml:"l,attr"`
	Messages               []Message `xml:"m"`
}

type HistoricChange struct {
	Timestamp Time `xml:"ts"`
}

type HistoricDelay struct {
	HistoricChange
	ArrivalEvent   Time        `xml:"ar,attr"`
	DepartureEvent Time        `xml:"dp,attr"`
	Source         DelaySource `xml:"src,attr"`
	Description    string      `xml:"cod,attr"`
}

type HistoricPlatformChange struct {
	HistoricChange
	ArrivalPlatform   string `xml:"ar,attr"`
	DeparturePlatform string `xml:"dp,attr"`
	// Detailed cause of track change
	Description string `xml:"cot,attr"`
}

type Connection struct {
	ID               string           `xml:"id,attr"`
	Timestamp        Time             `xml:"ts,attr"`
	EvaNumber        int              `xml:"eva,attr"`
	ConnectionStatus ConnectionStatus `xml:"cs,attr"`
	StopOfMissedTrip TimetableStop    `xml:"ref"`
	TametableStop    TimetableStop    `xml:"s"`
}

type ReferenceTripRelation struct {
	ReferenceTrip               ReferenceTrip               `xml:"rt"`
	ReferenceTripRelationToStop ReferenceTripRelationToStop `xml:"rts"`
}

type MessageType struct{}

type Priority struct{}

type DistributorMessage struct {
	Type         DistributorType `xml:"t,attr"`
	Name         string          `xml:"n,attr"`
	InternalText string          `xml:"int,attr"`
	Timestamp    Time            `xml:"ts,attr"`
}

// Trip type for example "p"
type TripType string

type EventStatus struct{}

type DelaySource struct{}

type ConnectionStatus struct{}

type ReferenceTrip struct {
	ID                                       string                 `xml:"id,attr"`
	Cancelation                              bool                   `xml:"c,attr"`
	ReferenceTripLabel                       ReferenceTripLabel     `xml:"rtl"`
	ReferencTripStopLabelStartDepartureEvent ReferenceTripStopLabel `xml:"sd"`
	ReferencTripStopLabelEndArrivalEvent     ReferenceTripStopLabel `xml:"ea"`
}

type ReferenceTripRelationToStop struct{}

type DistributorType struct{}

type ReferenceTripLabel struct {
	TripTrainNumber string `xml:"n,attr"`
	Category        string `xml:"c,attr"`
}

type ReferenceTripStopLabel struct {
	Index       int    `xml:"i,attr"`
	PlannedTime Time   `xml:"pt,attr"`
	EvaNumber   int    `xml:"eva,attr"`
	Name        string `xml:"n,attr"`
}
