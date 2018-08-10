package geode

import (
	"github.com/dhconnelly/rtreego"
	"encoding/json"
	"github.com/tomplex/wktfile"
	"fmt"
)

type SpatialDataset struct {
	Features 	[]*SpatialData `json:"features"`
	index 		*rtreego.Rtree
}

func (sds *SpatialDataset) SearchIntersect(feat *SpatialData) []*SpatialData {
	results := sds.index.SearchIntersect(feat.Bounds())
	features := make([]*SpatialData, len(results))
	for i, v := range results {
		features[i] = v.(*SpatialData)
	}
	return features
}

func (sds *SpatialDataset) Size() (int) {
	return len(sds.Features)
}

func NewDataset(features []*SpatialData) (*SpatialDataset, error) {
	tree := rtreego.NewTree(3, 25, 50)
	collection := &SpatialDataset{
		Features: features,
		index:    tree,
	}
	collection.insertFeaturesIntoTree()

	return collection, nil
}
func (sds *SpatialDataset) insertFeaturesIntoTree() {
	for i := range sds.Features {
		sds.index.Insert(sds.Features[i])
	}
}

func NewDatasetFromGeoJSON(data []byte) (*SpatialDataset, error) {
	dataset := SpatialDataset{}
	err := json.Unmarshal(data, &dataset)
	if err != nil {
		return nil, err
	}

	tree := rtreego.NewTree(3, 25, 50)
	dataset.index = tree

	dataset.insertFeaturesIntoTree()

	return &dataset, nil
}

func NewDatasetFromWKTFile(wktFile *wktfile.WKTFile, idColumn, geomColumn string, columns ...string) (*SpatialDataset, error) {
	if len(wktFile.Header) == 0 && len(columns) == 0 {
		return nil, fmt.Errorf("Cannot create a SpatialDataset from a WKTFile with no header and no columns specified")
	}
	if len(columns) == 0 {
		columns = wktFile.Header
	}

	geomColumnIndex := -1
	idColumnIndex := -1

	for i, obj := range columns {
		if obj == geomColumn {
			geomColumnIndex = i
		}
		if obj == idColumn {
			idColumnIndex = i
		}
	}

	if geomColumnIndex == -1 {
		return nil, fmt.Errorf("Specified geom column not found in header or columns")
	}

	if idColumnIndex == -1 {
		return nil, fmt.Errorf("Specified id column not found in header or columns")
	}

	features := make([]*SpatialData, len(wktFile.Rows))
	for i, row := range wktFile.Rows {
		id := row[idColumnIndex]
		wktGeometry := row[geomColumnIndex]
		properties := createMapFromRow(row, columns, idColumnIndex, geomColumnIndex)
		features[i] = FromWKT(wktGeometry, id, properties)
	}

	collection, err := NewDataset(features)

	if err != nil {
		return nil, err
	}
	return collection, nil
}

func createMapFromRow(row, columns []string, idIndex, geomIndex int) (map[string]interface{}) {
	mapRow := make(map[string]interface{})
	for i, column := range columns {
		if i == idIndex || i == geomIndex {
			continue
		}
		mapRow[column] = row[i]
	}
	return mapRow
}
