package tests

//import (
//	"bytes"
//	"encoding/binary"
//	"fmt"
//	"github.com/gl1n0m3c/IT_LAB_INIT/internal/decoder"
//	leb1282 "github.com/jcalabro/leb128"
//	"github.com/stretchr/testify/assert"
//	"strings"
//	"testing"
//	"time"
//)
//
//func EncodeString(buf *bytes.Buffer, key, value string) {
//	keyBytes := []byte(key)
//	valueBytes := []byte(value)
//
//	// Длина ключа (2 байта)
//	binary.Write(buf, binary.BigEndian, uint16(len(keyBytes)))
//	// Длина значения (2 байта)
//	binary.Write(buf, binary.BigEndian, uint16(len(valueBytes)))
//	// Тип значения (1 байт, 0x00 для строк)
//	buf.WriteByte(0x00)
//	// Ключ (UTF-8 строка в байтах)
//	buf.Write(keyBytes)
//	// Значение (UTF-8 строка в байтах)
//	buf.Write(valueBytes)
//}
//
//func EncodeInt(buf *bytes.Buffer, key string, value int) {
//	keyBytes := []byte(key)
//	valueBytes := leb1282.EncodeS64(int64(value))
//
//	// Длина ключа (2 байта)
//	binary.Write(buf, binary.BigEndian, uint16(len(keyBytes)))
//	// Длина значения (2 байта)
//	binary.Write(buf, binary.BigEndian, uint16(len(valueBytes)))
//	// Тип значения (1 байт, 0x01 для целых чисел)
//	buf.WriteByte(0x01)
//	// Ключ (UTF-8 строка в байтах)
//	buf.Write(keyBytes)
//	// Значение (UTF-8 строка в байтах)
//	buf.Write(valueBytes)
//}
//
//func EncodeTime(buf *bytes.Buffer, key string, value time.Time) {
//	// Форматируем время в строку согласно указанному формату
//	formattedTime := value.Format(time.RFC3339)
//	formattedTime = strings.Replace(formattedTime, "Z", "+", 1)
//	fmt.Printf("%#v\n", formattedTime)
//
//	// Преобразуем ключ в байты
//	keyBytes := []byte(key)
//	// Преобразуем отформатированное время в байты
//	valueBytes := []byte(formattedTime)
//
//	// Длина ключа (2 байта)
//	binary.Write(buf, binary.BigEndian, uint16(len(keyBytes)))
//	// Длина значения (2 байта)
//	binary.Write(buf, binary.BigEndian, uint16(len(valueBytes)))
//	// Тип значения (1 байт, 0x00 для строк)
//	buf.WriteByte(0x00)
//	// Ключ
//	buf.Write(keyBytes)
//	// Значение
//	buf.Write(valueBytes)
//}
//
//func TestCameraType1(t *testing.T) {
//	var buf, ans bytes.Buffer
//
//	now := time.Now()
//	newNow := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.Local)
//	camera := decoder.CaseDataType1{
//		TransportChars:   "ABC",
//		TransportNumbers: "1234",
//		TransportRegion:  "77",
//		CameraID:         "CAM123",
//		ViolationID:      "VIO456",
//		ViolationValue:   "Speeding",
//		SkillValue:       5,
//		Datetime:         newNow,
//	}
//
//	fmt.Printf("%#v\n", camera.Datetime)
//
//	EncodeString(&buf, "transport_chars", camera.TransportChars)
//	EncodeString(&buf, "transport_numbers", camera.TransportNumbers)
//	EncodeString(&buf, "transport_region", camera.TransportRegion)
//	EncodeString(&buf, "camera_id", camera.CameraID)
//	EncodeString(&buf, "violation_id", camera.ViolationID)
//	EncodeString(&buf, "violation_value", camera.ViolationValue)
//	EncodeInt(&buf, "skill_value", camera.SkillValue)
//	EncodeTime(&buf, "datetime", camera.Datetime)
//
//	binary.Write(&ans, binary.BigEndian, uint16(len(buf.Bytes())+2))
//	ans.Write(buf.Bytes())
//
//	fmt.Println(len(ans.Bytes()), ans.Bytes())
//
//	res, err := decoder.Decoder(&ans)
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Printf("%#v\n", res)
//	}
//
//	decodedCamera, err := MapToCameraType1(res)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	assert.Equal(t, camera, decodedCamera)
//	fmt.Println(camera.Datetime)
//	fmt.Println(decodedCamera.Datetime)
//
//}
