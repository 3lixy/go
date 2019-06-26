// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// 把json的反序列化"encoding/json"库变更为"github.com/json-iterator/go"，性能上升100%
// 贡献人：郑鹏林

package render

import (
	"bytes"
	"net/http"

	"azoya/nova/json"
)

type JSON struct {
	Data interface{}
}

type IndentedJSON struct {
	Data interface{}
}

type SecureJSON struct {
	Prefix string
	Data   interface{}
}

type SecureJSONPrefix string

var jsonContentType = []string{"application/json; charset=utf-8"}

func (r JSON) Render(w http.ResponseWriter) (err error) {
	if err = WriteJSON(w, r.Data); err != nil {
		panic(err)
	}
	return
}

func (r JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func WriteJSON(w http.ResponseWriter, obj interface{}) error {
	writeContentType(w, jsonContentType)
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	w.Write(jsonBytes)
	return nil
}

func (r IndentedJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	jsonBytes, err := json.MarshalIndent(r.Data, "", "    ")
	if err != nil {
		return err
	}
	w.Write(jsonBytes)
	return nil
}

func (r IndentedJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func (r SecureJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	jsonBytes, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}
	// if the jsonBytes is array values
	if bytes.HasPrefix(jsonBytes, []byte("[")) && bytes.HasSuffix(jsonBytes, []byte("]")) {
		w.Write([]byte(r.Prefix))
	}
	w.Write(jsonBytes)
	return nil
}

func (r SecureJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}
