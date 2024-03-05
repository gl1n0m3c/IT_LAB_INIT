package tests

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/decoder"
	leb1282 "github.com/jcalabro/leb128"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func EncodeString(buf *bytes.Buffer, key, value string) {
	keyBytes := []byte(key)
	valueBytes := []byte(value)

	// Длина ключа (2 байта)
	binary.Write(buf, binary.BigEndian, uint16(len(keyBytes)))
	// Длина значения (2 байта)
	binary.Write(buf, binary.BigEndian, uint16(len(valueBytes)))
	// Тип значения (1 байт, 0x00 для строк)
	buf.WriteByte(0x00)
	// Ключ (UTF-8 строка в байтах)
	buf.Write(keyBytes)
	// Значение (UTF-8 строка в байтах)
	buf.Write(valueBytes)
}

func EncodeInt(buf *bytes.Buffer, key string, value int) {
	keyBytes := []byte(key)
	valueBytes := leb1282.EncodeS64(int64(value))

	// Длина ключа (2 байта)
	binary.Write(buf, binary.BigEndian, uint16(len(keyBytes)))
	// Длина значения (2 байта)
	binary.Write(buf, binary.BigEndian, uint16(len(valueBytes)))
	// Тип значения (1 байт, 0x01 для целых чисел)
	buf.WriteByte(0x01)
	// Ключ (UTF-8 строка в байтах)
	buf.Write(keyBytes)
	// Значение (UTF-8 строка в байтах)
	buf.Write(valueBytes)
}

func TestCameraType1(t *testing.T) {
	var buf, ans bytes.Buffer

	EncodeString(&buf, "transport_chars", caseDataType1.TransportChars)
	EncodeString(&buf, "transport_numbers", caseDataType1.TransportNumbers)
	EncodeString(&buf, "transport_region", caseDataType1.TransportRegion)
	EncodeString(&buf, "camera_id", caseDataType1.CameraID)
	EncodeString(&buf, "violation_id", caseDataType1.ViolationID)
	EncodeString(&buf, "violation_value", caseDataType1.ViolationValue)
	EncodeInt(&buf, "skill_value", caseDataType1.SkillValue)
	EncodeString(&buf, "datetime", caseDataType1.Datetime)

	ans.Write(buf.Bytes())

	res, err := decoder.Decoder(&ans)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v\n", res)
	}

	decodedCamera, err := decoder.MapToStruct(res)
	if err != nil {
		fmt.Println(err)
	}

	resCase, err := decodedCamera.CameraDataToCaseBase()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, caseDataType1, decodedCamera)
	assert.Equal(t, resultCase1, resCase)
}

func TestCameraType2(t *testing.T) {
	var dataBytes []byte
	for i := 0; i+8 <= len(byteString2); i += 8 {
		byteString := byteString2[i : i+8]
		byteValue, _ := strconv.ParseUint(byteString, 2, 8)
		dataBytes = append(dataBytes, byte(byteValue))
	}

	result, err := decoder.Decoder(bytes.NewBuffer(dataBytes[2:]))
	if err != nil {
		t.Fatal(err)
	}

	cameraModel, err := decoder.MapToStruct(result)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(cameraModel)

	resCase, err := cameraModel.CameraDataToCaseBase()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, caseDataType2, cameraModel)
	assert.Equal(t, resultCase2, resCase)
}

func TestCameraType3(t *testing.T) {
	var dataBytes []byte
	for i := 0; i+8 <= len(byteString3); i += 8 {
		byteString := byteString3[i : i+8]
		byteValue, _ := strconv.ParseUint(byteString, 2, 8)
		dataBytes = append(dataBytes, byte(byteValue))
	}

	result, err := decoder.Decoder(bytes.NewBuffer(dataBytes[2:]))
	if err != nil {
		t.Fatal(err)
	}

	cameraModel, err := decoder.MapToStruct(result)
	if err != nil {
		t.Fatal(err)
	}

	resCase, err := cameraModel.CameraDataToCaseBase()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, caseDataType3, cameraModel)
	assert.Equal(t, resultCase3, resCase)
}
