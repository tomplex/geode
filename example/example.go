package main

import (
	"github.com/tomplex/geode"
	"github.com/tomplex/wktfile"
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"
)


func main() {
	townsFilePath, _ := filepath.Abs("example/VT_Hartford_Hartland.geojson")
	fmt.Println(townsFilePath)
	vtTownsFile, _ := os.Open(townsFilePath)
	buildingsBytes, _ := ioutil.ReadAll(vtTownsFile)

	fmt.Println("Creating towns dataset...")
	vtTowns, err := geode.NewDatasetFromGeoJSON(buildingsBytes)
	fmt.Println("Done")
	if err != nil {
		fmt.Errorf("Error creating towns dataset")
		panic(err)
	}

	markersFilePath, _ := filepath.Abs("example/VT_Roadside_Historic_Markers.wkt")
	markersWktFile, err := wktfile.Read(markersFilePath)
	if err != nil {
		fmt.Errorf("Error loading WKT file")
		panic(err)
	}
	markersDataset, err := geode.NewDatasetFromWKTFile(markersWktFile, "id", "wkt")

	fmt.Println("Built markers dataset")
	if err != nil {
		panic(err)
	}

	for _, town := range vtTowns.Features {
		fmt.Println("Processing", town.GetProperty("TOWNNAME"))
		matchingMarkers := markersDataset.SearchIntersect(town)
		fmt.Println("Found", len(matchingMarkers), "markers")
		for _, marker := range matchingMarkers {
			fmt.Println(marker.GetProperty("name"))
		}
	}
}