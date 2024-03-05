package decoder

import (
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"strings"
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

type TransportInfo struct {
	Chars   string `json:"chars"`
	Numbers string `json:"numbers"`
	Region  string `json:"region"`
}

type CameraInfo struct {
	ID string `json:"id"`
}

type ViolationInfo struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type SkillInfo struct {
	Value int `json:"value"`
}

type DatetimeInfo struct {
	Year      int    `json:"year"`
	Month     int    `json:"month"`
	Day       int    `json:"day"`
	Hour      int    `json:"hour"`
	Minute    int    `json:"minute"`
	Seconds   int    `json:"seconds"`
	UtcOffset string `json:"utc_offset"`
}

type CaseDataType2 struct {
	Transport TransportInfo `json:"transport"`
	Camera    CameraInfo    `json:"camera"`
	Violation ViolationInfo `json:"violation"`
	Skill     SkillInfo     `json:"skill"`
	Datetime  DatetimeInfo  `json:"datetime"`
}

type CaseDataType3 struct {
	Transport string        `json:"transport"`
	Camera    CameraInfo    `json:"camera"`
	Violation ViolationInfo `json:"violation"`
	Skill     int           `json:"skill"`
	Datetime  string        `json:"datetime"`
}

func (c CaseDataType1) CameraDataToCaseBase() (models.CaseBase, error) {
	var caseData models.CaseBase

	parsedTime, err := time.Parse(time.RFC3339, c.Datetime)
	if err != nil {
		return models.CaseBase{}, fmt.Errorf("Недопустимый тип времени")
	}
	parsedTime = parsedTime.UTC()

	charsRunes := []rune(c.TransportChars)
	caseData.CameraID = c.CameraID
	caseData.Transport = string(charsRunes[:1]) + c.TransportNumbers + string(charsRunes[1:]) + c.TransportRegion
	caseData.ViolationID = c.ViolationID
	caseData.ViolationValue = c.ViolationValue
	caseData.Level = c.SkillValue
	caseData.Datetime = parsedTime

	return caseData, nil
}

func (c CaseDataType2) CameraDataToCaseBase() (models.CaseBase, error) {
	var caseData models.CaseBase

	c.Datetime.UtcOffset = strings.Replace(c.Datetime.UtcOffset, "+", "", 1)

	datetime := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d+0%s:00",
		c.Datetime.Year,
		c.Datetime.Month,
		c.Datetime.Day,
		c.Datetime.Hour,
		c.Datetime.Minute,
		c.Datetime.Seconds,
		c.Datetime.UtcOffset,
	)

	parsedTime, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		return models.CaseBase{}, fmt.Errorf("Недопустимый тип времени: %v", err)
	}
	parsedTime = parsedTime.UTC()

	charsRunes := []rune(c.Transport.Chars)
	caseData.CameraID = c.Camera.ID
	caseData.Transport = string(charsRunes[:1]) + c.Transport.Numbers + string(charsRunes[1:]) + c.Transport.Region
	caseData.ViolationID = c.Violation.ID
	caseData.ViolationValue = c.Violation.Value
	caseData.Level = c.Skill.Value
	caseData.Datetime = parsedTime

	return caseData, nil
}

func (c CaseDataType3) CameraDataToCaseBase() (models.CaseBase, error) {
	var caseData models.CaseBase

	parsedTime, err := time.Parse(time.RFC3339, c.Datetime)
	if err != nil {
		return models.CaseBase{}, fmt.Errorf("Недопустимый тип времени")
	}
	parsedTime = parsedTime.UTC()

	caseData.CameraID = c.Camera.ID
	caseData.Transport = c.Transport
	caseData.ViolationID = c.Violation.ID
	caseData.ViolationValue = c.Violation.Value
	caseData.Level = c.Skill
	caseData.Datetime = parsedTime

	return caseData, nil
}
