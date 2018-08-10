package geode

import (
	"github.com/paulsmith/gogeos/geos"
	"github.com/dhconnelly/rtreego"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkt"
	"encoding/json"
	"fmt"
)

type SpatialData struct {
	Id         interface{} `json:"id,omitempty"`
	Geom       *geos.Geometry
	Properties map[string]interface{} `json:"properties"`
}

// Bounds allows us to implement the Spatial interface - and store our Features in the rtree.
func (sd *SpatialData) Bounds() (*rtreego.Rect) {
	geomType, _ := sd.Geom.Type()
	if geomType == geos.POINT {
		coordinates, _ := sd.Geom.Coords()
		p := rtreego.Point{coordinates[0].X, coordinates[0].Y}
		rect, _ := rtreego.NewRect(p, []float64{0.000001, 0.000001})
		return rect
	}
	coords, err := sd.coords()
	if err != nil {

	}
	pmin, pmax := coords[0], coords[2]

	pt := rtreego.Point{pmin.X, pmin.Y}
	xlen := pmax.X - pmin.X
	ylen := pmax.Y - pmin.Y
	lens := []float64{xlen, ylen}
	rect, err := rtreego.NewRect(pt, lens)
	if err != nil {}

	return rect
}

func (sd *SpatialData) GetProperty(name string) (interface{}) {
	return sd.Properties[name]
}

// Intersects helper method on SpatialData type
func (sd *SpatialData) Intersects(other *geos.Geometry) (intersects bool, err error) {
	intersects, err = sd.Geom.Intersects(other)
	return intersects, err
}

// Intersection helper method on SpatialData type
func (sd *SpatialData) Intersection(other *geos.Geometry) (intersection *geos.Geometry, err error) {
	intersection, err = sd.Geom.Intersection(other)
	return intersection, err
}

// Helper function for calculating geometry bounds
func (sd *SpatialData) coords() ([]geos.Coord, error) {
	envelope, err := sd.Geom.Envelope()
	if err != nil {
		return nil, err
	}

	boundary, err := envelope.Boundary()
	if err != nil {
		return nil, err
	}

	coords, err := boundary.Coords()
	if err != nil {
		return nil, err
	}

	return coords, nil
}

// Allow us to create a new SpatialData object from JSON. We use the go-geom package for its helpful encoding/decoding tools - to
// convert our GeoJSON geometry into something that the geos package can use.
func (sd *SpatialData) UnmarshalJSON(data []byte) error {
	jf := &jsonFeature{}
	err := json.Unmarshal(data, jf)
	if err != nil {
		return err
	}
	if jf.Type != "Feature" {
		return fmt.Errorf("geojson: not a feature: type=%s", jf.Type)
	}
	decodedGeom, err := jf.Geometry.Decode()
	if err != nil {
		return err
	}
	wktGeom, _ := wkt.Marshal(decodedGeom)
	if err != nil {
		return err
	}

	geosGeom := geos.Must(geos.FromWKT(wktGeom))

	*sd = SpatialData{
		Id:         jf.Id,
		Properties: jf.Properties,
		Geom: geosGeom,
	}
	return nil
}

type jsonFeature struct {
	Id         interface{} `json:"id,omitempty"`
	Type       string      `json:"type"`
	Geometry   geojson.Geometry  `json:"geometry"`
	Properties map[string]interface{}  `json:"properties,omitempty"`
}


// Helper function to create a new SpatialData from the individual parts - WKT geometry, id, and properties.
func FromWKT(wkt string, id interface{}, properties map[string]interface{}) (*SpatialData) {
	geom := SpatialData{
		Geom:       geos.Must(geos.FromWKT(wkt)),
		Id:         id,
		Properties: properties,
	}
	return &geom
}

func mustString(str string, err error) (string) {
	if err != nil {
		panic(err)
	}
	return str
}