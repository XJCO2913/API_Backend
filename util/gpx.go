package util

import (
	"fmt"
	"regexp"

	"github.com/tkrajina/gpxgo/gpx"
)

// Return a GPX Parser with gpx data bytes
func GPXParser(gpxDataBytes []byte) (*gpx.GPX, error) {
	gpxHandler, err := gpx.ParseBytes(gpxDataBytes)
	if err != nil {
		return nil, err
	}

	return gpxHandler, nil
}

// Convert gpx data to [[lon, lat]...] array
func GPXToLonLat(gpxDataBytes []byte) ([]string, error) {
	gpxHandler, err := gpx.ParseBytes(gpxDataBytes)
	if err != nil {
		return nil, err
	}

	res := []string{}
	for _, track := range gpxHandler.Tracks {
		for _, segment := range track.Segments {
			for _, point := range segment.Points {
				res = append(res, fmt.Sprintf("%v %v", point.Longitude, point.Latitude))
			}
		}
	}

	return res, nil
}

// Convert LINESTRING(x x, y y, z z,...) to x x, y y, z z,...
func GPXRoute(linestring string) (string, error) {
	re := regexp.MustCompile(`\((.*?)\)`)

	matches := re.FindStringSubmatch(linestring)
	if len(matches) > 1 {
		// match
		return matches[1], nil
	}

	// not match
	return "", fmt.Errorf("invalid linestring format")
}
