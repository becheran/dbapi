package dbapi_test

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/becheran/dbapi"
	"github.com/stretchr/testify/assert"
)

func TestStationUnmarshal(t *testing.T) {
	input := `<stations><station p="1|2|3|4|5" name="Köln Airport-Businesspark" eva="8003370" ds100="KBP" db="true" creationts="21-12-22 11:46:29.268"/></stations>`
	var res dbapi.Stations
	err := xml.Unmarshal([]byte(input), &res)
	assert.Nil(t, err)
	assert.Len(t, res.Stations, 1)

	stat := res.Stations[0]
	assert.Equal(t, stat.Name, "Köln Airport-Businesspark")
	assert.Equal(t, stat.EvaNumber, "8003370")
	assert.Equal(t, stat.Ds100, "KBP")
	assert.True(t, *stat.DB)

	assert.Equal(t, stat.Platforms, "1|2|3|4|5")

	assert.Equal(t, stat.Creation.Day(), 22)
	assert.Equal(t, stat.Creation.Month(), time.Month(12))
	assert.Equal(t, stat.Creation.Year(), 2021)
	assert.Equal(t, stat.Creation.Hour(), 11)
	assert.Equal(t, stat.Creation.Minute(), 46)
	assert.Equal(t, stat.Creation.Second(), 29)
	assert.Equal(t, stat.Creation.Nanosecond(), 268000000)

	assert.True(t, stat.Updates.IsZero())
}

func TestTimetableUnmarshal(t *testing.T) {
	input := `<?xml version='1.0' encoding='UTF-8'?>
	<timetable station='Berlin Alexanderplatz (S)'>
	  <s id="-8834013167311099628-2112282111-14">
		<tl f="S" t="p" o="08" c="S" n="7144"/>
		<ar pt="2112282139" pp="4" l="7" ppth="Ahrensfelde (S)|Berlin Mehrower Allee|Berlin Raoul-Wallenberg-Str.|Berlin-Marzahn|Berlin Poelchaustr.|Berlin Springpfuhl|Berlin-Friedrichsfelde Ost|Berlin-Lichtenberg (S)|Berlin N&#246;ldnerplatz|Berlin Ostkreuz (S)|Berlin Warschauer Stra&#223;e|Berlin Ostbahnhof (S)|Berlin Jannowitzbr&#252;cke"/>
		<dp pt="2112282140" pp="4" l="7" ppth="Berlin Hackescher Markt|Berlin Friedrichstra&#223;e (S)|Berlin Hbf (S-Bahn)|Berlin Bellevue|Berlin-Tiergarten|Berlin Zoologischer Garten (S)|Berlin Savignyplatz|Berlin Charlottenburg (S)|Berlin Westkreuz|Berlin-Grunewald|Berlin-Nikolassee|Berlin Wannsee (S)|Potsdam Griebnitzsee (S)|Potsdam-Babelsberg|Potsdam Hbf (S)"/>
	  </s>
	  <s id="-426624563964702035-2112282049-15">
		<tl f="S" t="p" o="08" c="S" n="9132"/>
		<ar pt="2112282132" pp="4" l="9" ppth="Flughafen BER - Terminal 1-2 (S-Bahn)|Wa&#223;mannsdorf|Flughafen BER - Terminal 5 (Sch&#246;nefeld)|Berlin Gr&#252;nbergallee|Berlin-Altglienicke|Berlin-Adlershof|Berlin-Johannisthal|Berlin-Sch&#246;neweide (S)|Berlin Baumschulenweg|Berlin Pl&#228;nterwald|Berlin Treptower Park|Berlin Warschauer Stra&#223;e|Berlin Ostbahnhof (S)|Berlin Jannowitzbr&#252;cke"/>
		<dp pt="2112282133" pp="4" l="9" ppth="Berlin Hackescher Markt|Berlin Friedrichstra&#223;e (S)|Berlin Hbf (S-Bahn)|Berlin Bellevue|Berlin-Tiergarten|Berlin Zoologischer Garten (S)|Berlin Savignyplatz|Berlin Charlottenburg (S)|Berlin Westkreuz|Berlin Messe S&#252;d (Eichkamp)|Berlin Heerstra&#223;e|Berlin Olympiastadion|Berlin-Pichelsberg|Berlin-Stresow|Berlin-Spandau (S)"/>
	  </s>
	</timetable>`
	var res dbapi.Timetable
	err := xml.Unmarshal([]byte(input), &res)

	assert.Nil(t, err)
	assert.Equal(t, "Berlin Alexanderplatz (S)", res.Station)
	assert.Equal(t, 2021, res.Stops[0].Arrival.PlannedTime.Year())
}
