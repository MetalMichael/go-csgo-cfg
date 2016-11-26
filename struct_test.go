// Copyright 2014 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package csgo_cfg

import (
	"bytes"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type testStruct struct {
	Name         string `csgo:"NAME"`
	Age          int
	Male         bool
	Money        float64
	Born         time.Time
	Time         time.Duration `csgo:"Duration"`
	Unused       int `csgo:"-"`
	Unsigned     uint
	Omitted      bool `csgo:"omitthis,omitempty"`
}

const _CONF_DATA_STRUCT = `
NAME Unknwon
Age 21
Male true
Money 1.25
Born 1993-10-07T20:17:05Z
Duration 2h45m
Unsigned 3
omitthis true
Cities "HangZhou|Boston"
Visits "1993-10-07T20:17:05Z, 1993-10-07T20:17:05Z"
Years "1993,1994"
Numbers "10010,10086"
Ages "18,19"
Populations "12345678,98765432"
Coordinates "192.168,10.11"
Note "Hello world"

GPA 2.8

Here "there"
When "then"
`

type unsupport struct {
	Byte byte
}

type unsupport2 struct {
	Others struct {
		Cities byte
	}
}

type unsupport3 struct {
	Cities byte
}

type unsupport4 struct {
	*unsupport3 `csgo:"Others"`
}

type defaultValue struct {
	Name   string
	Age    int
	Male   bool
	Money  float64
	Born   time.Time
	Cities []string
}

type fooBar struct {
	Here, When string
}

const _INVALID_DATA_CONF_STRUCT = `
Name 
Age age
Male 123
Money money
Born nil
Cities 
`

func Test_Struct(t *testing.T) {
	Convey("Map to struct", t, func() {
		Convey("Map file to struct", func() {
			ts := new(testStruct)
			So(MapTo(ts, []byte(_CONF_DATA_STRUCT)), ShouldBeNil)

			So(ts.Name, ShouldEqual, "Unknwon")
			So(ts.Age, ShouldEqual, 21)
			So(ts.Male, ShouldBeTrue)
			So(ts.Money, ShouldEqual, 1.25)
			So(ts.Unsigned, ShouldEqual, 3)

			t, err := time.Parse(time.RFC3339, "1993-10-07T20:17:05Z")
			So(err, ShouldBeNil)
			So(ts.Born.String(), ShouldEqual, t.String())

			dur, err := time.ParseDuration("2h45m")
			So(err, ShouldBeNil)
			So(ts.Time.Seconds(), ShouldEqual, dur.Seconds())
		})

		Convey("Map section to struct", func() {
			foobar := new(fooBar)
			f, err := Load([]byte(_CONF_DATA_STRUCT))
			So(err, ShouldBeNil)

			So(f.Section("").MapTo(foobar), ShouldBeNil)
			So(foobar.Here, ShouldEqual, "there")
			So(foobar.When, ShouldEqual, "then")
		})

		Convey("Map to non-pointer struct", func() {
			cfg, err := Load([]byte(_CONF_DATA_STRUCT))
			So(err, ShouldBeNil)
			So(cfg, ShouldNotBeNil)

			So(cfg.MapTo(testStruct{}), ShouldNotBeNil)
		})

		Convey("Map to unsupported type", func() {
			cfg, err := Load([]byte(_CONF_DATA_STRUCT))
			So(err, ShouldBeNil)
			So(cfg, ShouldNotBeNil)

			cfg.NameMapper = func(raw string) string {
				if raw == "Byte" {
					return "NAME"
				}
				return raw
			}
			So(cfg.MapTo(&unsupport{}), ShouldNotBeNil)
			//So(cfg.MapTo(&unsupport2{}), ShouldNotBeNil)
			//So(cfg.MapTo(&unsupport4{}), ShouldNotBeNil)
		})

		Convey("Map to omitempty field", func() {
			ts := new(testStruct)
			So(MapTo(ts, []byte(_CONF_DATA_STRUCT)), ShouldBeNil)

			So(ts.Omitted, ShouldEqual, true)
		})

		Convey("Map from invalid data source", func() {
			So(MapTo(&testStruct{}, "hi"), ShouldNotBeNil)
		})

		Convey("Map to wrong types and gain default values", func() {
			cfg, err := Load([]byte(_INVALID_DATA_CONF_STRUCT))
			So(err, ShouldBeNil)

			t, err := time.Parse(time.RFC3339, "1993-10-07T20:17:05Z")
			So(err, ShouldBeNil)
			dv := &defaultValue{"Joe", 10, true, 1.25, t, []string{"HangZhou", "Boston"}}
			So(cfg.MapTo(dv), ShouldBeNil)
			So(dv.Name, ShouldEqual, "Joe")
			So(dv.Age, ShouldEqual, 10)
			So(dv.Male, ShouldBeTrue)
			So(dv.Money, ShouldEqual, 1.25)
			So(dv.Born.String(), ShouldEqual, t.String())
			So(strings.Join(dv.Cities, ","), ShouldEqual, "HangZhou,Boston")
		})
	})

	Convey("Reflect from struct", t, func() {
		type Author struct {
			Name      string `csgo:"NAME"`
			Male      bool
			Age       int
			Height    uint
			GPA       float64
			Date      time.Time
			NeverMind string `csgo:"-"`
		}

		t, err := time.Parse(time.RFC3339, "1993-10-07T20:17:05Z")
		So(err, ShouldBeNil)
		a := &Author{"Unknwon", true, 21, 100, 2.8, t, ""}
		cfg := Empty()
		So(ReflectFrom(cfg, a), ShouldBeNil)

		var buf bytes.Buffer
		_, err = cfg.WriteTo(&buf)
		So(err, ShouldBeNil)

		So(buf.String(), ShouldEqual,
            "NAME   Unknwon" + LineBreak +
            "Male   true" + LineBreak +
            "Age    21" + LineBreak +
            "Height 100" + LineBreak +
            "GPA    2.8" + LineBreak +
            "Date   1993-10-07T20:17:05Z" + LineBreak)


		Convey("Reflect from non-point struct", func() {
			So(ReflectFrom(cfg, Author{}), ShouldNotBeNil)
		})

		Convey("Reflect from struct with omitempty", func() {
			cfg := Empty()
			type SpecialStruct struct {
				FirstName  string    `csgo:"first_name"`
				LastName   string    `csgo:"last_name"`
				JustOmitMe string    `csgo:"omitempty"`
				LastLogin  time.Time `csgo:"last_login,omitempty"`
				LastLogin2 time.Time `csgo:",omitempty"`
				NotEmpty   int       `csgo:"omitempty"`
			}

			So(ReflectFrom(cfg, &SpecialStruct{FirstName: "John", LastName: "Doe", NotEmpty: 9}), ShouldBeNil)

			var buf bytes.Buffer
			_, err = cfg.WriteTo(&buf)
			So(buf.String(), ShouldEqual, 
            "first_name John" + LineBreak +
            "last_name  Doe" + LineBreak +
            "omitempty  9" + LineBreak)
		})
	})
}

type testMapper struct {
	PackageName string
}

func Test_NameGetter(t *testing.T) {
	Convey("Test name mappers", t, func() {
		So(MapToWithMapper(&testMapper{}, TitleUnderscore, []byte("package_name cfg")), ShouldBeNil)

		cfg, err := Load([]byte("PACKAGE_NAME cfg"))
		So(err, ShouldBeNil)
		So(cfg, ShouldNotBeNil)

		cfg.NameMapper = AllCapsUnderscore
		tg := new(testMapper)
		So(cfg.MapTo(tg), ShouldBeNil)
		So(tg.PackageName, ShouldEqual, "cfg")
	})
}
