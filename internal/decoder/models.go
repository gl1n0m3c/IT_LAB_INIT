package decoder

import (
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"time"
)

type CameraModel interface {
	CameraDataToCaseBase() (models.CaseBase, error)
}

type CaseDataType1 struct {
	TransportChars   string `json:"transport_chars"`
	TransportNumbers string `json:"transport_numbers"`
	TransportRegion  string `json:"transport_region"`
	CameraID         string `json:"camera_id"`
	ViolationID      string `json:"violation_id"`
	ViolationValue   string `json:"violation_value"`
	SkillValue       int    `json:"skill_value"`
	Datetime         string `json:"datetime"`
}

type CaseDataType2 struct {
	Transport struct {
		Chars   string `json:"chars"`
		Numbers string `json:"numbers"`
		Region  string `json:"region"`
	} `json:"transport"`
	Camera struct {
		ID string `json:"id"`
	} `json:"camera"`
	Violation struct {
		ID    string `json:"id"`
		Value string `json:"value"`
	} `json:"violation"`
	Skill struct {
		Value int `json:"value"`
	} `json:"skill"`
	Datetime struct {
		Year      int    `json:"year"`
		Month     int    `json:"month"`
		Day       int    `json:"day"`
		Hour      int    `json:"hour"`
		Minute    int    `json:"minute"`
		Seconds   int    `json:"seconds"`
		UtcOffset string `json:"utc_offset"`
	} `json:"datetime"`
}

type CaseDataType3 struct {
	Transport string `json:"transport"`
	Camera    struct {
		ID string `json:"id"`
	} `json:"camera"`
	Violation struct {
		ID    string `json:"id"`
		Value string `json:"value"`
	} `json:"violation"`
	Skill    int    `json:"skill"`
	Datetime string `json:"datetime"`
}

func (c CaseDataType1) CameraDataToCaseBase() (models.CaseBase, error) {
	var caseData models.CaseBase

	parsedTime, err := time.Parse(time.RFC3339, c.Datetime)
	if err != nil {
		return models.CaseBase{}, fmt.Errorf("Недопустимый тип времени")
	}

	caseData.CameraID = c.CameraID
	caseData.Transport = c.TransportChars + c.TransportNumbers + c.TransportRegion
	caseData.ViolationID = c.ViolationID
	caseData.ViolationValue = c.ViolationValue
	caseData.Level = c.SkillValue
	caseData.Datetime = parsedTime

	return caseData, nil
}

func (c CaseDataType2) CameraDataToCaseBase() (models.CaseBase, error) {
	var caseData models.CaseBase

	return caseData, nil
}

func (c CaseDataType3) CameraDataToCaseBase() (models.CaseBase, error) {
	var caseData models.CaseBase

	return caseData, nil
}
