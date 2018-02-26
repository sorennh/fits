package main

import (
	"bytes"
	"github.com/GeoNet/kit/weft"
	"net/http"
)

func method(r *http.Request, h http.Header, b *bytes.Buffer) error {
	err := weft.CheckQuery(r, []string{"GET"}, []string{}, []string{"typeID"})
	if err != nil {
		return err
	}

	h.Set("Content-Type", "application/json;version=1")

	typeID := r.URL.Query().Get("typeID")

	if typeID != "" {
		err = validType(typeID)
		if err != nil {
			return err
		}
	}

	var d string

	switch typeID {
	case "":
		err = db.QueryRow(
			`select row_to_json(fc) from (select array_to_json(array_agg(m)) as method  
		             from (select methodid as "methodID", method.name, method.description, method.reference 
		             from 
		             fits.type join fits.type_method using (typepk) 
			join fits.method using (methodpk)) as m) as fc`).Scan(&d)
	default:
		err = db.QueryRow(
			`select row_to_json(fc) from (select array_to_json(array_agg(m)) as method  
		             from (select methodid as "methodID", method.name, method.description, method.reference 
		             from 
		             fits.type join fits.type_method using (typepk) 
			join fits.method using (methodpk) 
			where type.typeID = $1) as m) as fc`, typeID).Scan(&d)
	}
	if err != nil {
		return err

	}

	b.WriteString(d)

	return nil
}
