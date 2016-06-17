package gpxgo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"math"
	"os"
)

type Waypoints []Wpt

type Person struct {
	XMLName xml.Name `xml:"author"`
	Name    string   `xml:"name,omitempty"`
	Email   *Email   `xml:"email,omitempty"`
	Link    *Link    `xml:"link,omitempty"`
}

type Email struct {
	XMLName xml.Name `xml:"email"`
	Id      string   `xml:"id,attr"`
	Domain  string   `xml:"domain,attr"`
}

type Link struct {
	XMLName xml.Name `xml:"link"`
	Href    string   `xml:"href,attr"`
	Text    string   `xml:"text,omitempty"`
	Type    string   `xml:"type,omitempty"`
}

type Extensions struct {
	XMLName xml.Name `xml:"extensions"`
	Info    string   `xml:",innerxml"`
}

type Copyright struct {
	XMLName xml.Name `xml:"copyright,omitempty"`
	Author  string   `xml:"author,attr"`
	Year    string   `xml:"year,omitempty"`
	License string   `xml:"license,omitempty"`
}

type Bounds struct {
	XMLName xml.Name `xml:"bounds"`
	MinLat  float64  `xml:"minlat,omitempty"`
	MinLon  float64  `xml:"minlon,omitempty"`
	MaxLat  float64  `xml:"maxlat,omitempty"`
	MaxLon  float64  `xml:"maxlon,omitempty"`
}

type Metadata struct {
	XMLName    xml.Name    `xml:"metadata"`
	Name       string      `xml:"name,omitempty"`
	Desc       string      `xml:"desc,omitempty"`
	Author     *Person     `xml:"author,omitempty"`
	Copyright  *Copyright  `xml:"copyright,omitempty"`
	Link       []Link      `xml:"link,omitempty"`
	Time       string      `xml:"time,omitempty"`
	Keywords   string      `xml:"keywords,omitempty"`
	Bounds     *Bounds     `xml:"bounds"`
	Extensions *Extensions `xml:"extensions,omitempty"`
}

type Wpt struct {
	// XMLName xml.Name `xml:"wpt"`
	// attr
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	//
	Ele           float64     `xml:"ele,omitempty"`
	Time          string      `xml:"time,omitempty"`
	Magvar        string      `xml:"magvar,omitempty"`
	Geoidheight   string      `xml:"geoidheight,omitempty"`
	Name          string      `xml:"name,omitempty"`
	Cmt           string      `xml:"cmt,omitempty"`
	Desc          string      `xml:"desc,omitempty"`
	Src           string      `xml:"src,omitempty"`
	Link          []Link      `xml:"link,omitempty"`
	Sym           string      `xml:"sym,omitempty"`
	Type          string      `xml:"type,omitempty"`
	Fix           string      `xml:"fix,omitempty"`
	Sat           uint        `xml:"sat,omitempty"`
	Hdop          float64     `xml:"hdop,omitempty"`
	Vdop          float64     `xml:"vdop,omitempty"`
	Pdop          float64     `xml:"pdop,omitempty"`
	Ageofdgpsdata float64     `xml:"ageofdgpsdata,omitempty"`
	Dgpsid        int         `xml:"dgpsid,omitempty"`
	Extensions    *Extensions `xml:"extensions,omitempty"`
}

type Rte struct {
	XMLName    xml.Name  `xml:"rte"`
	Name       string    `xml:"name,omitempty"`
	Cmt        string    `xml:"cmt,omitempty"`
	Desc       string    `xml:"desc,omitempty"`
	Src        string    `xml:"src,omitempty"`
	Link       []Link    `xml:"link,omitempty"`
	Number     uint      `xml:"number,omitempty"`
	Type       string    `xml:"type,omitempty"`
	Extensions string    `xml:"extensions,omitempty"`
	Waypoints  Waypoints `xml:"rtept,omitempty"`
}

type Trkseg struct {
	XMLName    xml.Name    `xml:"trkseg"`
	Waypoints  Waypoints   `xml:"trkpt"`
	Extensions *Extensions `xml:"extensions,omitempty"`
}

type Trk struct {
	XMLName    xml.Name `xml:"trk"`
	Name       string   `xml:"name,omitempty"`
	Cmt        string   `xml:"cmt,omitempty"`
	Desc       string   `xml:"desc,omitempty"`
	Src        string   `xml:"src,omitempty"`
	Link       []Link   `xml:"link,omitempty"`
	Number     uint     `xml:"number,omitempty"`
	Type       string   `xml:"type,omitempty"`
	Extensions string   `xml:"extensions,omitempty"`
	Segments   []Trkseg `xml:"trkseg,omitempty"`
}

type Gpx struct {
	XMLName      xml.Name  `xml:"gpx"`
	XMLNs        string    `xml:"xmlns,attr"`
	XMLNsXsi     string    `xml:"xmlns:xsi,attr"`
	XMLSchemaLoc string    `xml:"xsi:schemaLocation,attr"`
	Version      string    `xml:"version,attr"`
	Creator      string    `xml:"creator,attr"`
	Metadata     *Metadata `xml:"metadata,omitempty"`
	Waypoints    Waypoints `xml:"wpt,omitempty"`
	Routes       []Rte     `xml:"rte,omitempty"`
	Tracks       []Trk     `xml:"trk,omitempty"`
	Extensions   string    `xml:"extensions,omitempty"`
}

type TimeBounds struct {
	StartTime float64
	EndTime   float64
}

/*==========================================================*/
// Static
func ParseWithContent(content []byte) (*Gpx, error) {
	gpx := NewGpx()
	err := xml.Unmarshal(content, gpx)
	if err != nil {
		return nil, err
	}
	return gpx, nil
}

func ParseWithReader(o io.Reader) (*Gpx, error) {
	// data, err := ioutil.ReadAll(o)
	// if err != nil {
	// 	return nil, err
	// }

	// return ParseWithContent(data)
	gpx := NewGpx()
	d := xml.NewDecoder(o)
	d.CharsetReader = charset.NewReaderLabel
	err := d.Decode(gpx)
	if err != nil {
		return nil, err
	}
	return gpx, nil
}

func ParseWithPath(path string) (*Gpx, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ParseWithReader(file)
}

func NewGpx() *Gpx {
	return &Gpx{
		XMLNs:        "http://www.topografix.com/GPX/1/1",
		XMLNsXsi:     "http://www.w3.org/2001/XMLSchema-instance",
		XMLSchemaLoc: "http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd",
		Version:      "1.1",
		Creator:      "https://github.com/pikeszfish/gpxgo",
	}
}

/*==========================================================*/
// Bounds
func minBounds() *Bounds {
	return &Bounds{
		MaxLat: -math.MaxFloat64,
		MinLat: math.MaxFloat64,
		MaxLon: -math.MaxFloat64,
		MinLon: math.MaxFloat64,
	}
}

func (b *Bounds) merge(b2 *Bounds) {
	b.MaxLat = math.Max(b.MaxLat, b2.MaxLat)
	b.MinLat = math.Min(b.MinLat, b2.MinLat)
	b.MaxLon = math.Max(b.MaxLon, b2.MaxLon)
	b.MinLon = math.Min(b.MinLon, b2.MinLon)
}

func (b Bounds) String() string {
	return fmt.Sprintf("Min: %+v, %+v Max: %+v, %+v",
		b.MinLat, b.MinLon, b.MaxLat, b.MaxLon)
}

/*==========================================================*/
// Gpx
func toXML(n interface{}) []byte {
	content, _ := xml.MarshalIndent(n, "", "  ")
	return content
}

func (g *Gpx) ToXML() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(xml.Header)
	buffer.Write(toXML(g))
	return buffer.Bytes()
}

func (g *Gpx) Clone() *Gpx {
	newgpx := new(Gpx)
	newgpx.XMLNs = g.XMLNs
	newgpx.XMLNsXsi = g.XMLNsXsi
	newgpx.XMLSchemaLoc = g.XMLSchemaLoc
	newgpx.Version = g.Version
	newgpx.Creator = g.Creator

	if g.Metadata != nil {
		newgpx.Metadata = &Metadata{
			Name: g.Metadata.Name,
			Desc: g.Metadata.Desc,
			// Author:     g.Metadata.Author,
			// Copyright:  g.Metadata.Copyright,
			Link:     make([]Link, len(g.Metadata.Link)),
			Time:     g.Metadata.Time,
			Keywords: g.Metadata.Keywords,
			// Bounds:     g.Metadata.Bounds,
			Extensions: g.Metadata.Extensions,
		}
		copy(newgpx.Metadata.Link, g.Metadata.Link)
		if g.Metadata.Author != nil {
			newgpx.Metadata.Author = &Person{
				Name: g.Metadata.Author.Name,
			}
			if g.Metadata.Author.Email != nil {
				newgpx.Metadata.Author.Email = &Email{
					Id:     g.Metadata.Author.Email.Id,
					Domain: g.Metadata.Author.Email.Domain,
				}
			}
			if g.Metadata.Author.Link != nil {
				newgpx.Metadata.Author.Link = &Link{
					Href: g.Metadata.Author.Link.Href,
					Text: g.Metadata.Author.Link.Text,
					Type: g.Metadata.Author.Link.Type,
				}
			}
		}
		if g.Metadata.Copyright != nil {
			newgpx.Metadata.Copyright = &Copyright{
				Author:  g.Metadata.Copyright.Author,
				Year:    g.Metadata.Copyright.Year,
				License: g.Metadata.Copyright.License,
			}
		}
		if g.Metadata.Bounds != nil {
			newgpx.Metadata.Bounds = &Bounds{
				MaxLat: g.Metadata.Bounds.MaxLat,
				MinLat: g.Metadata.Bounds.MinLat,
				MaxLon: g.Metadata.Bounds.MaxLon,
				MinLon: g.Metadata.Bounds.MinLon,
			}
		}
	}
	newgpx.Waypoints = make([]Wpt, len(g.Waypoints))
	newgpx.Routes = make([]Rte, len(g.Routes))
	newgpx.Tracks = make([]Trk, len(g.Tracks))
	copy(newgpx.Waypoints, g.Waypoints)
	copy(newgpx.Routes, g.Routes)
	copy(newgpx.Tracks, g.Tracks)

	return newgpx
}

func (g *Gpx) Bounds() *Bounds {
	b := minBounds()
	for _, trk := range g.Tracks {
		b.merge(trk.Bounds())
	}
	return b
}

func (g *Gpx) Length2D() float64 {
	var length2d float64
	for _, trk := range g.Tracks {
		length2d += trk.Length2D()
	}
	return length2d
}

func (g *Gpx) Length3D() float64 {
	var length3d float64
	for _, trk := range g.Tracks {
		length3d += trk.Length3D()
	}
	return length3d
}

func (g *Gpx) UphillDownhill() (float64, float64) {
	var (
		uphill   float64
		downhill float64
	)
	for _, trk := range g.Tracks {
		u, d := trk.UphillDownhill()
		uphill += u
		downhill += d
	}
	return uphill, downhill
}

func (g *Gpx) ElevationExtremes() (float64, float64) {
	// elevations := []int{}
	var (
		elevations []float64
		min        float64
		max        float64
	)
	for _, trk := range g.Tracks {
		min, max = trk.ElevationExtremes()
		elevations = append(elevations, min)
		elevations = append(elevations, max)
	}
	min = elevations[0]
	max = elevations[0]
	for _, elevation := range elevations {
		if elevation < min {
			min = elevation
		}
		if elevation > max {
			max = elevation
		}
	}
	return min, max
}

func (g *Gpx) RemoveTime() {
	for _, trk := range g.Tracks {
		trk.RemoveTime()
	}
}

func (g *Gpx) RemoveElevation() {
	for _, trk := range g.Tracks {
		trk.RemoveElevation()
	}
}

/*==========================================================*/
// Routes
func (r *Rte) Length2D() float64 {
	var length2d float64
	for i := 1; i < len(r.Waypoints); i++ {
		length2d += r.Waypoints[i].Length2D(&r.Waypoints[i-1])
	}
	return length2d
}

func (r *Rte) Length3D() float64 {
	var length3d float64
	for i := 1; i < len(r.Waypoints); i++ {
		length3d += r.Waypoints[i].Length3D(&r.Waypoints[i-1])
	}
	return length3d
}

func (r *Rte) RemoveTime() {
	for _, waypoint := range r.Waypoints {
		waypoint.RemoveTime()
	}
}

func (r *Rte) RemoveElevation() {
	for _, waypoint := range r.Waypoints {
		waypoint.RemoveElevation()
	}
}

/*==========================================================*/
// Tracks
func (t *Trk) Bounds() *Bounds {
	b := minBounds()
	for _, seg := range t.Segments {
		b.merge(seg.Bounds())
	}
	return b
}

func (t *Trk) Length2D() float64 {
	var length2d float64
	for _, seg := range t.Segments {
		length2d += seg.Length2D()
	}
	return length2d
}

func (t *Trk) Length3D() float64 {
	var length3d float64
	for _, seg := range t.Segments {
		length3d += seg.Length3D()
	}
	return length3d
}

func (t *Trk) RemoveTime() {
	for _, seg := range t.Segments {
		seg.RemoveTime()
	}
}

func (t *Trk) UphillDownhill() (float64, float64) {
	var (
		uphill   float64
		downhill float64
	)
	for _, seg := range t.Segments {
		u, d := seg.UphillDownhill()
		uphill += u
		downhill += d
	}
	return uphill, downhill
}

func (t *Trk) ElevationExtremes() (float64, float64) {
	var (
		elevations []float64
		min        float64
		max        float64
	)
	for _, seg := range t.Segments {
		min, max = seg.ElevationExtremes()
		elevations = append(elevations, min)
		elevations = append(elevations, max)
	}
	min = elevations[0]
	max = elevations[0]
	for _, elevation := range elevations {
		if elevation < min {
			min = elevation
		}
		if elevation > max {
			max = elevation
		}
	}
	return min, max
}

func (t *Trk) RemoveElevation() {
	for _, seg := range t.Segments {
		seg.RemoveElevation()
	}
}

func (t Trk) String() string {
	return fmt.Sprintf("Name: %s Segment Count: %d", t.Name, len(t.Segments))
}

/*==========================================================*/
// Trkseg
func (ts *Trkseg) Bounds() *Bounds {
	b := minBounds()
	b2 := minBounds()
	for _, wp := range ts.Waypoints {
		b2.MaxLat = wp.Lat
		b2.MaxLon = wp.Lon
		b2.MinLat = wp.Lat
		b2.MinLon = wp.Lon
		b.merge(b2)
	}
	return b
}

func (ts *Trkseg) Length2D() float64 {
	var length2d float64
	for i := 1; i < len(ts.Waypoints); i++ {
		length2d += ts.Waypoints[i].Length2D(&ts.Waypoints[i-1])
	}
	return length2d
}

func (ts *Trkseg) Length3D() float64 {
	var length3d float64
	for i := 1; i < len(ts.Waypoints); i++ {
		length3d += ts.Waypoints[i].Length3D(&ts.Waypoints[i-1])
	}
	return length3d
}

func (ts Trkseg) String() string {
	return fmt.Sprintf("Waypoints Count: %d", len(ts.Waypoints))
}

func (ts *Trkseg) UphillDownhill() (float64, float64) {
	return 0.0, 0.0
	// var (
	// 	uphill             float64
	// 	downhill           float64
	// 	smoothedElevations []float64
	// )
	// // for _, wp := range ts.Waypoints {

	// // }
	// for i := 0; i < len(ts.Waypoints); i++ {
	// 	if i > 0 && i < len(ts.Waypoints)-1 {
	// 		previousWp := ts.Waypoints[i-1]
	// 		nextWp := ts.Waypoints[i+1]
	// 		if previousWp.Ele && ts.Waypoints[i].Ele && nextWp.Ele {
	// 			smoothedElevations = append(smoothedElevations, previousWp*0.3 + )
	// 		}
	// 	}

	// }
}

func (ts *Trkseg) ElevationExtremes() (float64, float64) {
	var (
		min float64
		max float64
	)
	for _, wp := range ts.Waypoints {
		if wp.Ele < min {
			min = wp.Ele
		}
		if wp.Ele > max {
			max = wp.Ele
		}
	}
	return min, max
}

func (ts *Trkseg) RemoveTime() {
	for _, wp := range ts.Waypoints {
		wp.RemoveTime()
	}
}

func (ts *Trkseg) RemoveElevation() {
	for _, wp := range ts.Waypoints {
		wp.RemoveElevation()
	}
}

/*==========================================================*/
// Waypoints / []Wpt
// func (w Waypoints) Bounds() *Bounds {
// 	b := minBounds()
// 	for _, wp := range w {
// 		b.merge(&Bounds{
// 			MaxLat: wp.Lat,
// 			MinLat: wp.Lat,
// 			MaxLon: wp.Lon,
// 			MinLon: wp.Lon,
// 		})
// 	}
// 	return b
// }

// func (w Waypoints) Length2D() float64 {
// 	var length2d float64
// 	for i := 1; i < len(w); i++ {
// 		length2d += w[i].Length2D(&w[i-1])
// 	}
// 	return length2d
// }

// func (w Waypoints) Length3D() float64 {
// 	var length3d float64
// 	for i := 1; i < len(w); i++ {
// 		length3d += w[i].Length3D(&w[i-1])
// 	}
// 	return length3d
// }

// func (w Waypoints) RemoveTime() {
// 	for _, wp := range w {
// 		wp.RemoveTime()
// 	}
// }

// func (w Waypoints) RemoveElevation() {
// 	for _, wp := range w {
// 		wp.RemoveElevation()
// 	}
// }

/*==========================================================*/
// Wpt
func (wp Wpt) String() string {
	return fmt.Sprintf("Wpt Lat: %f Lon: %f Name: %s \n", wp.Lat, wp.Lon, wp.Name)
}

func (wp *Wpt) Length2D(wp2 *Wpt) float64 {
	return Distance(wp.Lat, wp.Lon, 0.0, wp2.Lat, wp2.Lon, 0.0, false, false)
}

func (wp *Wpt) Length3D(wp2 *Wpt) float64 {
	return Distance(wp.Lat, wp.Lon, wp.Ele, wp2.Lat, wp2.Lon, wp2.Ele, true, false)
}

func (wp *Wpt) RemoveTime() {
	wp.Time = ""
}

func (wp *Wpt) RemoveElevation() {
	wp.Ele = 0.0
}

func (wp *Wpt) DeepCopy() *Wpt {
	newWpt := &Wpt{
		Lat:           wp.Lat,
		Lon:           wp.Lon,
		Ele:           wp.Ele,
		Time:          wp.Time,
		Magvar:        wp.Magvar,
		Geoidheight:   wp.Geoidheight,
		Name:          wp.Name,
		Cmt:           wp.Cmt,
		Desc:          wp.Desc,
		Src:           wp.Src,
		Link:          make([]Link, len(wp.Link)),
		Sym:           wp.Sym,
		Type:          wp.Type,
		Fix:           wp.Fix,
		Sat:           wp.Sat,
		Hdop:          wp.Hdop,
		Vdop:          wp.Vdop,
		Pdop:          wp.Pdop,
		Ageofdgpsdata: wp.Ageofdgpsdata,
		Dgpsid:        wp.Dgpsid,
		Extensions:    wp.Extensions,
	}
	copy(newWpt.Link, wp.Link)
	return newWpt
}

func (wp *Wpt) Move(ld *LocationDelta) {
	ld.Move(wp)
}

func (wp *Wpt) DistanceAngle(wp2 *Wpt) *LocationDelta {
	return &LocationDelta{
		Angle:    Bearing(wp.Lat, wp.Lon, wp2.Lat, wp2.Lon),
		Distance: Distance(wp.Lat, wp.Lon, 0, wp2.Lat, wp2.Lon, 0, false, false),
	}
}
