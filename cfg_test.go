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
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Version(t *testing.T) {
	Convey("Get version", t, func() {
		So(Version(), ShouldEqual, _VERSION)
	})
}

const _CONF_DATA = `
cash_team_elimination_hostage_map_ct 2000
cash_team_elimination_hostage_map_t 1000
cash_team_win_by_time_running_out_bomb 3250
cash_team_win_by_time_running_out_hostage 3250
ff_damage_reduction_grenade 0.85            // How much to reduce damage done to teammates by a thrown grenade.  Range is from 0 - 1 (with 1 being damage equal to what is done to an enemy)
ff_damage_reduction_bullets 0.33            // How much to reduce damage done to teammates when shot.  Range is from 0 - 1 (with 1 being damage equal to what is done to an enemy)
ff_damage_reduction_other 0.4               // How much to reduce damage done to teammates by things other than bullets and grenades.  Range is from 0 - 1 (with 1 being damage equal to what is done to an enemy)
ff_damage_reduction_grenade_self 1          // How much to damage a player does to himself with his own grenade.  Range is from 0 - 1 (with 1 being damage equal to what is done to an enemy)
mp_afterroundmoney 0                // amount of money awared to every player after each round
mp_autokick 0                   // Kick idle/team-killing players
mp_autoteambalance 0
`

func Test_Load(t *testing.T) {
	Convey("Load from data sources", t, func() {

		Convey("Load with empty data", func() {
			So(Empty(), ShouldNotBeNil)
		})

		Convey("Load with multiple data sources", func() {
			cfg, err := Load([]byte(_CONF_DATA), "testdata/conf.cfg")
			So(err, ShouldBeNil)
			So(cfg, ShouldNotBeNil)

			f, err := Load([]byte(_CONF_DATA), "testdata/404.cfg")
			So(err, ShouldNotBeNil)
			So(f, ShouldBeNil)
		})
	})

	Convey("Bad load process", t, func() {

		Convey("Load from invalid data sources", func() {
			_, err := Load(_CONF_DATA)
			So(err, ShouldNotBeNil)

			f, err := Load("testdata/404.cfg")
			So(err, ShouldNotBeNil)
			So(f, ShouldBeNil)

			_, err = Load(1)
			So(err, ShouldNotBeNil)

			_, err = Load([]byte(""), 1)
			So(err, ShouldNotBeNil)
		})

		Convey("Load with bad section name", func() {
			_, err := Load([]byte("[]"))
			So(err, ShouldNotBeNil)

			_, err = Load([]byte("["))
			So(err, ShouldNotBeNil)
		})

		Convey("Load with bad keys", func() {
			_, err := Load([]byte(`"""name`))
			So(err, ShouldNotBeNil)

			_, err = Load([]byte(`"""name"""`))
			So(err, ShouldNotBeNil)

			_, err = Load([]byte(`""=1`))
			So(err, ShouldNotBeNil)

			_, err = Load([]byte(`=`))
			So(err, ShouldNotBeNil)

			_, err = Load([]byte(`name`))
			So(err, ShouldNotBeNil)
		})

		Convey("Load with bad values", func() {
			_, err := Load([]byte(`name="""Unknwon`))
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Get section and key insensitively", t, func() {
		cfg, err := InsensitiveLoad([]byte(_CONF_DATA), "testdata/conf.cfg")
		So(err, ShouldBeNil)
		So(cfg, ShouldNotBeNil)
	})

	Convey("Load with ignoring continuation lines", t, func() {
		cfg, err := LoadSources(LoadOptions{IgnoreContinuation: true}, []byte(`key1 a\b\ 
key2 c\d\`))
		So(err, ShouldBeNil)
		So(cfg, ShouldNotBeNil)

		So(cfg.Section("").Key("key1").String(), ShouldEqual, `a\b\`)
		So(cfg.Section("").Key("key2").String(), ShouldEqual, `c\d\`)
	})
}

func Test_LooseLoad(t *testing.T) {
	Convey("Loose load from data sources", t, func() {
		Convey("Loose load mixed with nonexistent file", func() {
			cfg, err := LooseLoad("testdata/404.cfg")
			So(err, ShouldBeNil)
			So(cfg, ShouldNotBeNil)
			var fake struct {
				Name string `csgo:"name"`
			}
			So(cfg.MapTo(&fake), ShouldBeNil)

			cfg, err = LooseLoad([]byte("name Unknwon"), "testdata/404.cfg")
			So(err, ShouldBeNil)
			So(cfg.Section("").Key("name").String(), ShouldEqual, "Unknwon")
			So(cfg.MapTo(&fake), ShouldBeNil)
			So(fake.Name, ShouldEqual, "Unknwon")
		})
	})

}

func Test_File_Append(t *testing.T) {
	Convey("Append data sources", t, func() {
		cfg, err := Load([]byte(""))
		So(err, ShouldBeNil)
		So(cfg, ShouldNotBeNil)

		So(cfg.Append([]byte(""), []byte("")), ShouldBeNil)

		Convey("Append bad data sources", func() {
			So(cfg.Append(1), ShouldNotBeNil)
			So(cfg.Append([]byte(""), 1), ShouldNotBeNil)
		})
	})
}

func Test_File_WriteTo(t *testing.T) {
	Convey("Write to somewhere", t, func() {
		var buf bytes.Buffer
		cfg := Empty()
		cfg.WriteTo(&buf)
	})
}

func Test_File_SaveTo_WriteTo(t *testing.T) {
	Convey("Save file", t, func() {
		cfg, err := Load([]byte(_CONF_DATA), "testdata/conf.cfg")
		So(err, ShouldBeNil)
		So(cfg, ShouldNotBeNil)

		cfg.Section("").Key("NAME").Comment = "Package name"
		cfg.Section("author").Comment = `Information about package author
# Bio can be written in multiple lines.`
		cfg.Section("advanced").Key("val w/ pound").SetValue("my#password")
		cfg.Section("advanced").Key("longest key has a colon : yes/no").SetValue("yes")
		So(cfg.SaveTo("testdata/conf_out.cfg"), ShouldBeNil)

		cfg.Section("author").Key("NAME").Comment = "This is author name"
	})
}

// Helpers for slice tests.
func float64sEqual(values []float64, expected ...float64) {
	So(values, ShouldHaveLength, len(expected))
	for i, v := range expected {
		So(values[i], ShouldEqual, v)
	}
}

func intsEqual(values []int, expected ...int) {
	So(values, ShouldHaveLength, len(expected))
	for i, v := range expected {
		So(values[i], ShouldEqual, v)
	}
}

func int64sEqual(values []int64, expected ...int64) {
	So(values, ShouldHaveLength, len(expected))
	for i, v := range expected {
		So(values[i], ShouldEqual, v)
	}
}

func uintsEqual(values []uint, expected ...uint) {
	So(values, ShouldHaveLength, len(expected))
	for i, v := range expected {
		So(values[i], ShouldEqual, v)
	}
}

func uint64sEqual(values []uint64, expected ...uint64) {
	So(values, ShouldHaveLength, len(expected))
	for i, v := range expected {
		So(values[i], ShouldEqual, v)
	}
}

func timesEqual(values []time.Time, expected ...time.Time) {
	So(values, ShouldHaveLength, len(expected))
	for i, v := range expected {
		So(values[i].String(), ShouldEqual, v.String())
	}
}
