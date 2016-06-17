package gpxgo

import (
	"github.com/bmizerany/assert"
	"log"
	"math"
	"os"
	"testing"
)

var g *Gpx

func init() {
	log.Println("gpx test init")
}

func TestParseWithPath(t *testing.T) {
	var err error
	path := "testdata/St_Louis_Zoo_sample.gpx"
	g, err = ParseWithPath(path)
	assert.NotEqual(t, nil, g)
	assert.Equal(t, nil, err)

	assert.Equal(t, "St Louis Zoo sample", g.Metadata.Name)
	assert.Equal(t, "2008-02-26T19:49:13", g.Metadata.Time)

	assert.Equal(t, 38.63473, g.Waypoints[0].Lat)
	assert.Equal(t, -90.29408, g.Waypoints[0].Lon)

	assert.Equal(t, 10, len(g.Waypoints))
}

func TestParseWithContent(t *testing.T) {
	buf := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="no" ?>
        <gpx xmlns="http://www.topografix.com/GPX/1/1" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3" xmlns:gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" creator="Oregon 400t" version="1.1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd">
          <metadata>
            <link href="http://www.garmin.com">
              <text>Garmin International</text>
            </link>
            <time>2009-10-17T22:58:43Z</time>
          </metadata>
          <trk>
            <name>Example GPX Document</name>
            <trkseg>
              <trkpt lat="47.644548" lon="-122.326897">
                <ele>4.46</ele>
                <time>2009-10-17T18:37:26Z</time>
              </trkpt>
              <trkpt lat="47.644548" lon="-122.326897">
                <ele>4.94</ele>
                <time>2009-10-17T18:37:31Z</time>
              </trkpt>
              <trkpt lat="47.644548" lon="-122.326897">
                <ele>6.87</ele>
                <time>2009-10-17T18:37:34Z</time>
              </trkpt>
            </trkseg>
          </trk>
        </gpx>`)
	gpx, err := ParseWithContent(buf)
	assert.NotEqual(t, nil, gpx)
	assert.Equal(t, nil, err)

	assert.Equal(t, "http://www.garmin.com", gpx.Metadata.Link[0].Href)
	assert.Equal(t, "2009-10-17T22:58:43Z", gpx.Metadata.Time)

	assert.Equal(t, "Example GPX Document", gpx.Tracks[0].Name)
	assert.Equal(t, 3, len(gpx.Tracks[0].Segments[0].Waypoints))
}

func TestParseWithReader(t *testing.T) {
	path := "testdata/St_Louis_Zoo_sample.gpx"
	file, err := os.Open(path)
	assert.Equal(t, nil, err)
	defer file.Close()

	gpx, err := ParseWithReader(file)
	assert.NotEqual(t, nil, gpx)
	assert.Equal(t, nil, err)
}

func TestCloneGpx(t *testing.T) {
	newgpx := g.Clone()
	assert.Equal(t, g.Metadata.Time, newgpx.Metadata.Time)
	assert.Equal(t, newgpx.ToXML(), g.ToXML())
}

func TestNewXml(t *testing.T) {
	gpx := NewGpx()
	gpxTrack := Trk{}

	gpxSegment := Trkseg{}
	gpxSegment.Waypoints = append(gpxSegment.Waypoints, Wpt{Lat: 32.1234, Lon: 121.1233, Ele: 1233})
	gpxSegment.Waypoints = append(gpxSegment.Waypoints, Wpt{Lat: 32.1235, Lon: 121.1234, Ele: 1234})
	gpxSegment.Waypoints = append(gpxSegment.Waypoints, Wpt{Lat: 32.1236, Lon: 121.1235, Ele: 1235})

	gpxTrack.Segments = append(gpxTrack.Segments, gpxSegment)
	gpx.Tracks = append(gpx.Tracks, gpxTrack)

	gpx.Waypoints = append(gpx.Waypoints, Wpt{Lat: 1.1111, Lon: 9.9999, Ele: 1111})
	gpx.Waypoints = append(gpx.Waypoints, Wpt{Lat: 2.2222, Lon: 8.8888, Ele: 2222})
	gpx.Waypoints = append(gpx.Waypoints, Wpt{Lat: 3.3333, Lon: 7.7777, Ele: 3333})
	actualXML := string(toXML(gpx))
	expectedXML := `<gpx xmlns="http://www.topografix.com/GPX/1/1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd" version="1.1" creator="https://github.com/pikeszfish/gpxgo">
  <wpt lat="1.1111" lon="9.9999">
    <ele>1111</ele>
  </wpt>
  <wpt lat="2.2222" lon="8.8888">
    <ele>2222</ele>
  </wpt>
  <wpt lat="3.3333" lon="7.7777">
    <ele>3333</ele>
  </wpt>
  <trk>
    <trkseg>
      <trkpt lat="32.1234" lon="121.1233">
        <ele>1233</ele>
      </trkpt>
      <trkpt lat="32.1235" lon="121.1234">
        <ele>1234</ele>
      </trkpt>
      <trkpt lat="32.1236" lon="121.1235">
        <ele>1235</ele>
      </trkpt>
    </trkseg>
  </trk>
</gpx>`
	assert.Equal(t, expectedXML, actualXML)
}

func TestRemoveElevation(t *testing.T) {
	g.RemoveElevation()
	for _, ts := range g.Tracks {
		for _, seg := range ts.Segments {
			for _, wp := range seg.Waypoints {
				assert.Equal(t, 0.0, wp.Ele)
			}
		}
	}
}

func TestBounds(t *testing.T) {
	sourceBounds := g.Bounds()
	expectBounds := &Bounds{
		MaxLat: 18.233625,
		MinLat: 18.231815,
		MaxLon: 109.521163,
		MinLon: 109.520261,
	}
	assert.Equal(t, sourceBounds, expectBounds)
}

func TestMoveWpt(t *testing.T) {
	wp1 := &Wpt{Lat: 32.11111, Lon: 121.22222}
	wp2 := &Wpt{Lat: 33.33333, Lon: 123.33333}

	assert.Equal(t, math.Abs(wp1.Lat-wp2.Lat) < 0.0000001, false)
	assert.Equal(t, math.Abs(wp1.Lon-wp2.Lon) < 0.0000001, false)

	ld := wp1.DistanceAngle(wp2)
	wp1.Move(ld)

	assert.Equal(t, math.Abs(wp1.Lat-wp2.Lat) < 0.0000001, true)
	assert.Equal(t, math.Abs(wp1.Lon-wp2.Lon) < 0.0000001, true)
}
