package decoder

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	leb1282 "github.com/jcalabro/leb128"
)

func Decoder(buf *bytes.Buffer) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	reader := bytes.NewReader(buf.Bytes())

	for reader.Len() > 0 {
		var keyLen, valueLen uint16
		var valueType uint8

		// Чтение длины ключа
		if err := binary.Read(reader, binary.BigEndian, &keyLen); err != nil {
			return nil, err
		}

		// Чтение длины значения
		if err := binary.Read(reader, binary.BigEndian, &valueLen); err != nil {
			return nil, err
		}

		// Чтение типа значения
		if err := binary.Read(reader, binary.BigEndian, &valueType); err != nil {
			return nil, err
		}

		// Чтение ключа
		key := make([]byte, keyLen)
		if _, err := reader.Read(key); err != nil {
			return nil, err
		}

		// Чтение значения в зависимости от его типа
		value := make([]byte, valueLen)
		if _, err := reader.Read(value); err != nil {
			return nil, err
		}

		switch valueType {
		case 0x00: // Строка
			result[string(key)] = string(value)
		case 0x01: // Целое число
			intValueReader := bytes.NewReader(value)
			intValue, err := leb1282.DecodeS64(intValueReader)
			if err != nil {
				return nil, fmt.Errorf("Не удалось декодировать LEB128 значение для ключа %s: %v", string(key), err)
			}
			result[string(key)] = intValue
		case 0x02:
			nestedBuf := bytes.NewBuffer(value)
			nestedResult, err := Decoder(nestedBuf)
			if err != nil {
				return nil, err
			}
			result[string(key)] = nestedResult
		default:
			return nil, fmt.Errorf("неизвестный тип значения: %v", valueType)
		}
	}

	return result, nil
}

func MapToStruct(data map[string]interface{}) (CameraModel, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return CaseDataType1{}, err
	}

	// Попытка десериализации в каждую структуру
	var ct1 CaseDataType1
	if err := json.Unmarshal(jsonData, &ct1); err == nil && ct1.CameraID != "" {
		return ct1, nil
	}

	var ct2 CaseDataType2
	if err := json.Unmarshal(jsonData, &ct2); err == nil && ct2.Camera.ID != "" {
		return ct2, nil
	}

	var ct3 CaseDataType3
	if err := json.Unmarshal(jsonData, &ct3); err == nil && ct3.Camera.ID != "" {
		return ct3, nil
	}

	return CaseDataType1{}, fmt.Errorf("Не удалось определить тип структуры")
}
