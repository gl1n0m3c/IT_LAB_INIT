package cameras

import "github.com/gl1n0m3c/IT_LAB_INIT/internal/models"

var testcaseCameraCreate = []models.CameraBase{
	{
		Type:        "camerus1",
		Coordinates: [2]float64{0, 0},
		Description: "test description1",
	},
	{
		Type:        "camerus2",
		Coordinates: [2]float64{0.1315465165, 0.49841653256},
		Description: "test description2",
	},
	{
		Type:        "camerus3",
		Coordinates: [2]float64{0.25656633, 0},
		Description: "test description3",
	},
}
