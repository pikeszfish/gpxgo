package gpxgo

import (
	"math"
)

const (
	// One degree in meters:
	ONE_DEGREE   = 1000. * 10000.8 / 90.
	EARTH_RADIUS = 6371 * 1000
)

type LocationDelta struct {
	Distance float64
	Angle    float64
}

type Location struct {
	Latitude  float64
	Longitude float64
	Elevation float64
}

/*==========================================================*/
// utils

func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180.
}

func Degrees(radians float64) float64 {
	return radians * 180. / math.Pi
}

func Bearing(lat1, lon1, lat2, lon2 float64) float64 {
	lat1r := Radians(lat1)
	lat2r := Radians(lat2)
	dlon := Radians(lon2 - lon1)
	y := math.Sin(dlon) * math.Cos(lat2r)
	x := math.Cos(lat1r)*math.Sin(lat2r) - math.Sin(lat1r)*math.Cos(lat2r)*math.Cos(dlon)
	return Degrees(math.Atan2(y, x))
}

func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	d_lat := Radians(lat1 - lat2)
	d_lon := Radians(lon1 - lon2)
	lat1 = Radians(lat1)
	lat2 = Radians(lat2)

	a := math.Sin(d_lat/2)*math.Sin(d_lat/2) + math.Sin(d_lon/2)*math.Sin(d_lon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := EARTH_RADIUS * c

	return d
}

func Distance(lat1, lon1, ele1, lat2, lon2, ele2 float64, threeD, haversine bool) float64 {
	if haversine || (math.Abs(lat1-lat2) > .2 || math.Abs(lon1-lon2) > .2) {
		return HaversineDistance(lat1, lon1, lat2, lon2)
	}

	coef := math.Cos(Radians(lat1))
	x := lat1 - lat2
	y := (lon1 - lon2) * coef

	distance2d := math.Sqrt(x*x+y*y) * ONE_DEGREE

	if !threeD || ele1 == ele2 {
		return distance2d
	}

	return math.Sqrt(math.Pow(distance2d, 2) + math.Pow((ele1-ele2), 2))
}

func ElevationAngle(l1, l2 *Location, radians bool) float64 {
	if l1.Elevation-l2.Elevation < 0.00001 {
		return 0.0
	}

	b := l2.Elevation - l1.Elevation
	a := l2.Distance2d(l1)

	if a < 0.00001 {
		return 0.0
	}

	angle := math.Atan(b / a)

	if radians {
		return angle
	}

	return Degrees(angle)
}

/*==========================================================*/
// LocationDelta

// http://www.movable-type.co.uk/scripts/latlong.html
func (ld *LocationDelta) Move(wp *Wpt) {
	l := ld.Distance / 6371.0 / 1000.0
	p_lat := Radians(wp.Lat)
	p_lon := Radians(wp.Lon)
	bearing := Radians(ld.Angle)

	lat := math.Asin(math.Sin(p_lat)*math.Cos(l) + math.Cos(p_lat)*math.Sin(l)*math.Cos(bearing))
	lon := p_lon + math.Atan2(math.Sin(bearing)*math.Sin(l)*math.Cos(p_lat), math.Cos(l)-math.Sin(p_lat)*math.Sin(lat))

	wp.Lat = Degrees(lat)
	wp.Lon = Degrees(lon)
}

/*==========================================================*/
// Location

func (l *Location) Distance2d(l2 *Location) float64 {
	return Distance(l.Latitude, l.Longitude, 0.0, l2.Latitude, l2.Longitude, 0.0, false, false)
}

func (l *Location) Distance3d(l2 *Location) float64 {
	return Distance(l.Latitude, l.Longitude, l.Elevation, l2.Latitude, l2.Longitude, l2.Elevation, true, false)
}
