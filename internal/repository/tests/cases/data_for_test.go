package cases

import (
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"time"
)

var testCases = []models.CaseBase{
	{
		CameraID:       "c70a6f9c-cb9b-4cd9-9425-77e8c8a5b072",
		Transport:      "A123BC97",
		ViolationID:    "e2f7e5b2-ea16-4c48-8a8a-5e2d1a2b8410",
		ViolationValue: "Speeding",
		Level:          1,
		Datetime:       time.Now().UTC(),
		PhotoUrl:       "http://example.com/photo1.jpg",
	},
	{
		CameraID:       "d95a3f0c-cb9b-4cd9-9425-77e8c8a5b072",
		Transport:      "B234CD98",
		ViolationID:    "f3a4b5c6-d7e8-9fa0-b1c2-d3e4f5a6b7c8",
		ViolationValue: "Parking",
		Level:          1,
		Datetime:       time.Now().UTC(),
		PhotoUrl:       "http://example.com/photo2.jpg",
	},
	{
		CameraID:       "a85a4f9c-cb9b-4cd9-9425-77e8c8a5b072",
		Transport:      "C345DE99",
		ViolationID:    "a1b2c3d4-e5f6-a7b8-c9d0-e1f2a3b4c5d6",
		ViolationValue: "Red Light",
		Level:          1,
		Datetime:       time.Now().UTC(),
		PhotoUrl:       "http://example.com/photo3.jpg",
	},
	{
		CameraID:       "b60a6f9c-cb9b-4cd9-9425-77e8c8a5b072",
		Transport:      "D456EF00",
		ViolationID:    "b1c2d3e4-f5a6-b7c8-d9e0-a1f2b3c4d5e6",
		ViolationValue: "Stop Sign",
		Level:          1,
		Datetime:       time.Now().UTC(),
		PhotoUrl:       "http://example.com/photo4.jpg",
	},
	{
		CameraID:       "e70a6f9c-cb9b-4cd9-9425-77e8c8a5b072",
		Transport:      "E567FG01",
		ViolationID:    "c1d2e3f4-a5b6-c7d8-e9f0-a1b2c3d4e5f6",
		ViolationValue: "Crosswalk",
		Level:          2,
		Datetime:       time.Now().UTC(),
		PhotoUrl:       "http://example.com/photo5.jpg",
	},
}

var testCasesRated = []models.RatedBase{
	{
		RatedCreate: models.RatedCreate{
			CaseID: 0,
			Choice: true,
		},
		SpecialistID: 3,
		Date:         time.Now().UTC(),
		Status:       "Correct",
	},
	{
		RatedCreate: models.RatedCreate{
			CaseID: 0,
			Choice: false,
		},
		SpecialistID: 4,
		Date:         time.Now().UTC(),
		Status:       "Incorrect",
	},
	{
		RatedCreate: models.RatedCreate{
			CaseID: 0,
			Choice: true,
		},
		SpecialistID: 5,
		Date:         time.Now().UTC(),
		Status:       "Correct",
	},
	{
		RatedCreate: models.RatedCreate{
			CaseID: 0,
			Choice: true,
		},
		SpecialistID: 1,
		Date:         time.Now().UTC(),
		Status:       "Unknown",
	},
	{
		RatedCreate: models.RatedCreate{
			CaseID: 0,
			Choice: false,
		},
		SpecialistID: 2,
		Date:         time.Now().UTC(),
		Status:       "Unknown",
	},
}
