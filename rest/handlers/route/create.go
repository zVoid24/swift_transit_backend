package route

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Property struct {
	Name string `json:"Name"`
}

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"` // Handle both types with an interface{}
}

type Feature struct {
	Type       string   `json:"type"`
	Properties Property `json:"properties"`
	Geom       Geometry `json:"geometry"`
}

type Geojson struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

// Handle incoming request and unmarshal JSON
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"Message": "Failed to read body",
		})
		return
	}

	// Print the raw body for debugging
	fmt.Println("=== Received Raw JSON ===")
	fmt.Println(string(body)) // Print the raw JSON text

	// Unmarshal the raw JSON body into a Geojson object
	var jsonData Geojson
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		// If unmarshaling fails, print the error to help identify the problem
		fmt.Printf("Error unmarshaling JSON: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"Message": "Failed to parse JSON",
		})
		return
	}

	// Debugging: Print the entire parsed JSON structure
	//fmt.Printf("Parsed JSON Data: %+v\n", jsonData.Features[1].Geom.Coordinates)
	if coords, ok := jsonData.Features[1].Geom.Coordinates.([]interface{}); ok {
		// Accessing latitude, longitude, and altitude
		if len(coords) == 3 {
			longitude := coords[0].(float64)
			latitude := coords[1].(float64)
			altitude := coords[2].(float64)

			fmt.Printf("Longitude: %f, Latitude: %f, Altitude: %f\n", longitude, latitude, altitude)
		}
	}

	// Loop through features and handle coordinates based on geometry type
	for _, feature := range jsonData.Features {
		fmt.Println("Feature:", feature.Properties.Name)
		if feature.Geom.Type == "Point" {
			// If it's a Point, coordinates should be a single 3-item array
			if coords, ok := feature.Geom.Coordinates.([]interface{}); ok {
				fmt.Println("Point Coordinates:", coords)
			} else {
				fmt.Println("Error: Point coordinates not in expected format")
			}
		} else if feature.Geom.Type == "LineString" {
			// If it's a LineString, coordinates should be an array of arrays of 3 items
			if coords, ok := feature.Geom.Coordinates.([][]interface{}); ok {
				for _, coord := range coords {
					fmt.Println("LineString Coordinate:", coord)
				}
			} else {
				fmt.Println("Error: LineString coordinates not in expected format")
			}
		}
	}
	h.utilHandler.SendData(w, jsonData, http.StatusOK)
}
