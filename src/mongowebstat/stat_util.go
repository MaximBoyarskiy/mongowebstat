package mongowebstat

import (
	"./jsonq"
	"bytes"
	"encoding/json"
	"reflect"
)

var (
	data map[string]interface{}
	dec  *json.Decoder
	jq   *jsonq.JsonQuery
)

func byteArrayToMongoStatRaw(b []byte) (*MongoStatRaw, error) {
	data = map[string]interface{}{}
	dec = json.NewDecoder(bytes.NewReader(b))
	err := dec.Decode(&data)
	if err != nil {
		return &MongoStatRaw{}, err
	}

	return mapToMongoStatRaw(data)
}

func mapToMongoStatRaw(data map[string]interface{}) (*MongoStatRaw, error) {
	jq = jsonq.NewQuery(data)
	r := new(MongoStatRaw)
	rv := reflect.ValueOf(r).Elem()
	for _, field := range MongoStatRawFields {
		if len(field.jsonqPath) < 2 {
			continue
		}
		jsonqPath := field.jsonqPath
		switch field.kind {
		case reflect.Int32:
			if i, err := jq.Int(jsonqPath...); err == nil {
				i32 := (int32)(i)
				rv.FieldByName(field.name).Set(reflect.ValueOf(&i32))
			} else {
				debug.Print("jq.Int32(jsonqPath...) err : ", err.Error())
			}
		case reflect.Int64:
			if i, err := jq.Int64(jsonqPath...); err == nil {
				rv.FieldByName(field.name).Set(reflect.ValueOf(&i))
			} else {
				debug.Print("jq.Int64(jsonqPath...) err : ", err.Error())
			}
		case reflect.Float64:
			if i, err := jq.Float(jsonqPath...); err == nil {
				rv.FieldByName(field.name).Set(reflect.ValueOf(&i))
			} else {
				debug.Print("jq.Float(jsonqPath...) err : ", err.Error())
			}
		case reflect.String:
			if s, err := jq.String(jsonqPath...); err == nil {
				rv.FieldByName(field.name).Set(reflect.ValueOf(&s))
			} else {
				debug.Print("jq.String(jsonqPath...) err: ", err.Error())
			}
		case reflect.Map:
			if m, err := jq.Object(jsonqPath...); err == nil {
				rv.FieldByName(field.name).Set(reflect.MakeMap(reflect.TypeOf(m)))
				for k, v := range m {
					rv.FieldByName(field.name).SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
				}
			} else {
				debug.Print("jq.Object(jsonqPath...) err: ", err.Error())
			}
		case reflect.Bool:
			if b, err := jq.Bool(jsonqPath...); err == nil {
				rv.FieldByName(field.name).Set(reflect.ValueOf(&b))
			} else {
				debug.Print("jq.Bool(jsonqPath...) err: ", err.Error())
			}
		}
	}
	return r, nil
}
