// Copyright (C) 2015 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package util

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// SetDefaults sets struct values to the values specified in the `default`
// struct annotations
func SetDefaults(data interface{}) error {
	s := reflect.ValueOf(data).Elem()
	t := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		tag := t.Field(i).Tag

		v := tag.Get("default")
		if len(v) > 0 {
			switch f.Interface().(type) {
			case string:
				f.SetString(v)

			case int:
				i, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return err
				}
				f.SetInt(i)

			case float64:
				i, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return err
				}
				f.SetFloat(i)

			case bool:
				f.SetBool(v == "true")

			case []string:
				// We don't do anything with string slices here. Any default
				// we set will be appended to by the XML decoder, so we fill
				// those after decoding.

			default:
				panic(f.Type())
			}
		}
	}
	return nil
}

// FillNilSlices sets default value on slices that are still nil.
func FillNilSlices(data interface{}) error {
	s := reflect.ValueOf(data).Elem()
	t := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		tag := t.Field(i).Tag

		v := tag.Get("default")
		if len(v) > 0 {
			switch f.Interface().(type) {
			case []string:
				if f.IsNil() {
					// Treat the default as a comma separated slice
					vs := strings.Split(v, ",")
					for i := range vs {
						vs[i] = strings.TrimSpace(vs[i])
					}

					rv := reflect.MakeSlice(reflect.TypeOf([]string{}), len(vs), len(vs))
					for i, v := range vs {
						rv.Index(i).SetString(v)
					}
					f.Set(rv)
				}
			}
		}
	}
	return nil
}

// UniqueStrings returns unqiue strings from the source slice
func UniqueStrings(ss []string) []string {
	var m = make(map[string]bool, len(ss))
	for _, s := range ss {
		m[strings.Trim(s, " ")] = true
	}

	var us = make([]string, 0, len(m))
	for k := range m {
		us = append(us, k)
	}

	sort.Strings(us)

	return us
}