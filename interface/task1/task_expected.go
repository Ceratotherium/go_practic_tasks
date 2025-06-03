package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Message map[string]string

func (m *Message) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber() // Включаем режим json.Number

	var rawStruct map[string]interface{}
	if err := dec.Decode(&rawStruct); err != nil {
		return err
	}

	for k, v := range rawStruct {
		switch v := v.(type) {
		case json.Number:
			// Пытаемся преобразовать в int64 если возможно
			if i, err := v.Int64(); err == nil {
				(*m)[k] = strconv.FormatInt(i, 10)
			} else { // Если не int, пробуем как float
				originStr := v.String()
				if f, err := v.Float64(); err == nil && strings.ContainsRune(originStr, '.') { // Проверка на наличие точки нужна, чтобы отловить очень большие числа
					(*m)[k] = strconv.FormatFloat(f, 'f', -1, 64)
				} else {
					// Если ничего не получилось, сохраняем как строку
					(*m)[k] = v.String()
				}
			}
		case string:
			(*m)[k] = v
		case bool:
			(*m)[k] = strconv.FormatBool(v)
		case nil:
			(*m)[k] = ""
		case map[string]interface{}, []interface{}:
			str, err := json.Marshal(v)
			if err != nil {
				return err
			}
			(*m)[k] = string(str)
		default:
			return fmt.Errorf("unsupported type %T", v)
		}
	}

	return nil
}
