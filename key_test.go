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
	"fmt"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Key(t *testing.T) {
	Convey("Test getting and setting values", t, func() {
		cfg, err := Load([]byte(_CONF_DATA), "testdata/conf.cfg")
		So(err, ShouldBeNil)
		So(cfg, ShouldNotBeNil)

		Convey("Get values in default section", func() {
			sec := cfg.Section(DEFAULT_SECTION)
			So(sec, ShouldNotBeNil)
			So(sec.Key("ammo_grenade_limit_default").Value(), ShouldEqual, "1")
			So(sec.Key("ammo_grenade_limit_default").String(), ShouldEqual, "1")
			So(sec.Key("ammo_grenade_limit_default").Validate(func(in string) string {
				return in
			}), ShouldEqual, "1")
			So(sec.Key("ammo_grenade_limit_default").Comment, ShouldEqual, "// test comment")
			So(sec.Key("cash_player_bomb_planted").String(), ShouldEqual, "300")
            So(sec.Key("cash_player_bomb_planted").Value(), ShouldEqual, "300")
		})

		Convey("Get auto-increment key names", func() {
			keys := cfg.Section("features").Keys()
			for i, k := range keys {
				So(k.Name(), ShouldEqual, fmt.Sprintf("#%d", i+1))
			}
		})

		Convey("Get parent-keys that are available to the child section", func() {
			parentKeys := cfg.Section("package.sub").ParentKeys()
			for _, k := range parentKeys {
				So(k.Name(), ShouldEqual, "CLONE_URL")
			}
		})

		Convey("Get overwrite value", func() {
			So(cfg.Section("").Key("ammo_grenade_limit_total").String(), ShouldEqual, "5")
		})

		Convey("Get sections", func() {
			sections := cfg.Sections()
			So(sections[0].Name(), ShouldEqual, DEFAULT_SECTION)
            So(len(sections), ShouldEqual, 1)
		})

		Convey("Get key hash", func() {
			cfg.Section("").KeysHash()
		})

		Convey("Set key value", func() {
			k := cfg.Section("author").Key("NAME")
			k.SetValue("无闻")
			So(k.String(), ShouldEqual, "无闻")
		})

		Convey("Delete a key", func() {
			cfg.Section(DEFAULT_SECTION).DeleteKey("UNUSED_KEY")
			_, err := cfg.Section(DEFAULT_SECTION).GetKey("UNUSED_KEY")
			So(err, ShouldNotBeNil)
		})

		Convey("Has Key (backwards compatible)", func() {
			sec := cfg.Section(DEFAULT_SECTION)
			haskey1 := sec.Haskey("ammo_grenade_limit_default")
			haskey2 := sec.Haskey("bot_quota")
			haskey3 := sec.Haskey("CLONE_URL_NO")
			So(haskey1, ShouldBeTrue)
			So(haskey2, ShouldBeTrue)
			So(haskey3, ShouldBeFalse)
		})

		Convey("Has Key", func() {
			sec := cfg.Section(DEFAULT_SECTION)
			haskey1 := sec.HasKey("ammo_grenade_limit_default")
			haskey2 := sec.HasKey("bot_quota")
			haskey3 := sec.HasKey("CLONE_URL_NO")
			So(haskey1, ShouldBeTrue)
			So(haskey2, ShouldBeTrue)
			So(haskey3, ShouldBeFalse)
		})

		Convey("Has Value", func() {
			sec := cfg.Section(DEFAULT_SECTION)
			hasvalue1 := sec.HasValue("2")
			hasvalue2 := sec.HasValue("doc")
			So(hasvalue1, ShouldBeTrue)
			So(hasvalue2, ShouldBeFalse)
		})
	})

	Convey("Test getting and setting bad values", t, func() {
		cfg, err := Load([]byte(_CONF_DATA), "testdata/conf.cfg")
		So(err, ShouldBeNil)
		So(cfg, ShouldNotBeNil)

		Convey("Create new key with empty name", func() {
			k, err := cfg.Section("").NewKey("", "")
			So(err, ShouldNotBeNil)
			So(k, ShouldBeNil)
		})

		Convey("Create new section with empty name", func() {
			s, err := cfg.NewSection("")
			So(err, ShouldNotBeNil)
			So(s, ShouldBeNil)
		})

		Convey("Create new sections with empty name", func() {
			So(cfg.NewSections(""), ShouldNotBeNil)
		})

		Convey("Get section that not exists", func() {
			s, err := cfg.GetSection("404")
			So(err, ShouldNotBeNil)
			So(s, ShouldBeNil)

			s = cfg.Section("404")
			So(s, ShouldNotBeNil)
		})
	})

	Convey("Test key hash clone", t, func() {
		cfg, err := Load([]byte(strings.Replace("network tcp,addr 127.0.0.1:6379,db 4,pool_size 100,idle_timeout 180", ",", "\n", -1)))
		So(err, ShouldBeNil)
		for _, v := range cfg.Section("").KeysHash() {
			So(len(v), ShouldBeGreaterThan, 0)
		}
	})

	Convey("Key has empty value", t, func() {
		_conf := `key1 
key2  // comment`
		cfg, err := Load([]byte(_conf))
		So(err, ShouldBeNil)
		So(cfg.Section("").Key("key1").Value(), ShouldBeEmpty)
	})
}

func newTestFile(block bool) *File {
	c, _ := Load([]byte(_CONF_DATA))
	c.BlockMode = block
	return c
}

func Benchmark_Key_Value(b *testing.B) {
	c := newTestFile(true)
	for i := 0; i < b.N; i++ {
		c.Section("").Key("NAME").Value()
	}
}

func Benchmark_Key_Value_NonBlock(b *testing.B) {
	c := newTestFile(false)
	for i := 0; i < b.N; i++ {
		c.Section("").Key("NAME").Value()
	}
}

func Benchmark_Key_Value_ViaSection(b *testing.B) {
	c := newTestFile(true)
	sec := c.Section("")
	for i := 0; i < b.N; i++ {
		sec.Key("NAME").Value()
	}
}

func Benchmark_Key_Value_ViaSection_NonBlock(b *testing.B) {
	c := newTestFile(false)
	sec := c.Section("")
	for i := 0; i < b.N; i++ {
		sec.Key("NAME").Value()
	}
}

func Benchmark_Key_Value_Direct(b *testing.B) {
	c := newTestFile(true)
	key := c.Section("").Key("NAME")
	for i := 0; i < b.N; i++ {
		key.Value()
	}
}

func Benchmark_Key_Value_Direct_NonBlock(b *testing.B) {
	c := newTestFile(false)
	key := c.Section("").Key("NAME")
	for i := 0; i < b.N; i++ {
		key.Value()
	}
}

func Benchmark_Key_String(b *testing.B) {
	c := newTestFile(true)
	for i := 0; i < b.N; i++ {
		_ = c.Section("").Key("NAME").String()
	}
}

func Benchmark_Key_String_NonBlock(b *testing.B) {
	c := newTestFile(false)
	for i := 0; i < b.N; i++ {
		_ = c.Section("").Key("NAME").String()
	}
}

func Benchmark_Key_String_ViaSection(b *testing.B) {
	c := newTestFile(true)
	sec := c.Section("")
	for i := 0; i < b.N; i++ {
		_ = sec.Key("NAME").String()
	}
}

func Benchmark_Key_String_ViaSection_NonBlock(b *testing.B) {
	c := newTestFile(false)
	sec := c.Section("")
	for i := 0; i < b.N; i++ {
		_ = sec.Key("NAME").String()
	}
}

func Benchmark_Key_SetValue(b *testing.B) {
	c := newTestFile(true)
	for i := 0; i < b.N; i++ {
		c.Section("").Key("NAME").SetValue("10")
	}
}

func Benchmark_Key_SetValue_VisSection(b *testing.B) {
	c := newTestFile(true)
	sec := c.Section("")
	for i := 0; i < b.N; i++ {
		sec.Key("NAME").SetValue("10")
	}
}
