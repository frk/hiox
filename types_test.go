package hxio

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/frk/compare"
	"github.com/frk/route"
)

func TestReaders(t *testing.T) {
	Bp := func(v bool) *bool { return &v }
	Ip := func(v int) *int { return &v }
	I8p := func(v int8) *int8 { return &v }
	I16p := func(v int16) *int16 { return &v }
	I32p := func(v int32) *int32 { return &v }
	I64p := func(v int64) *int64 { return &v }
	Up := func(v uint) *uint { return &v }
	U8p := func(v uint8) *uint8 { return &v }
	U16p := func(v uint16) *uint16 { return &v }
	U32p := func(v uint32) *uint32 { return &v }
	U64p := func(v uint64) *uint64 { return &v }
	F32p := func(v float32) *float32 { return &v }
	F64p := func(v float64) *float64 { return &v }
	Sp := func(v string) *string { return &v }

	t.Run("PathReaders", func(t *testing.T) {
		tests := []struct {
			path   route.Params
			reader PathReader
			result PathReader
			err    error
		}{{
			path:   route.NewParams("foo", "true"),
			reader: Bool{"foo": new(bool), "bar": new(bool)},
			result: Bool{"foo": Bp(true), "bar": Bp(false)},
		}, {
			path:   route.NewParams("foo", "2147483647", "bar", "-2147483648"),
			reader: Int{"foo": new(int), "bar": new(int)},
			result: Int{"foo": Ip(2147483647), "bar": Ip(-2147483648)},
		}, {
			path:   route.NewParams("foo", "127", "bar", "-128"),
			reader: Int8{"foo": new(int8), "bar": new(int8)},
			result: Int8{"foo": I8p(127), "bar": I8p(-128)},
		}, {
			path:   route.NewParams("foo", "32767", "bar", "-32768"),
			reader: Int16{"foo": new(int16), "bar": new(int16)},
			result: Int16{"foo": I16p(32767), "bar": I16p(-32768)},
		}, {
			path:   route.NewParams("foo", "2147483647", "bar", "-2147483648"),
			reader: Int32{"foo": new(int32), "bar": new(int32)},
			result: Int32{"foo": I32p(2147483647), "bar": I32p(-2147483648)},
		}, {
			path:   route.NewParams("foo", "9223372036854775807", "bar", "-9223372036854775808"),
			reader: Int64{"foo": new(int64), "bar": new(int64)},
			result: Int64{"foo": I64p(9223372036854775807), "bar": I64p(-9223372036854775808)},
		}, {
			path:   route.NewParams("foo", "4294967295", "bar", "-2147483648"),
			reader: Uint{"foo": new(uint), "bar": new(uint)},
			result: Uint{"foo": Up(4294967295), "bar": Up(0)},
		}, {
			path:   route.NewParams("foo", "255", "bar", "-128"),
			reader: Uint8{"foo": new(uint8), "bar": new(uint8)},
			result: Uint8{"foo": U8p(255), "bar": U8p(0)},
		}, {
			path:   route.NewParams("foo", "65535", "bar", "-32768"),
			reader: Uint16{"foo": new(uint16), "bar": new(uint16)},
			result: Uint16{"foo": U16p(65535), "bar": U16p(0)},
		}, {
			path:   route.NewParams("foo", "4294967295", "bar", "-2147483648"),
			reader: Uint32{"foo": new(uint32), "bar": new(uint32)},
			result: Uint32{"foo": U32p(4294967295), "bar": U32p(0)},
		}, {
			path:   route.NewParams("foo", "18446744073709551615", "bar", "-9223372036854775808"),
			reader: Uint64{"foo": new(uint64), "bar": new(uint64)},
			result: Uint64{"foo": U64p(18446744073709551615), "bar": U64p(0)},
		}, {
			path:   route.NewParams("foo", "21474.83647", "bar", "-0.23092222"),
			reader: Float32{"foo": new(float32), "bar": new(float32)},
			result: Float32{"foo": F32p(21474.83647), "bar": F32p(-0.23092222)},
		}, {
			path:   route.NewParams("foo", "21474.83647", "bar", "-0.23092222"),
			reader: Float64{"foo": new(float64), "bar": new(float64)},
			result: Float64{"foo": F64p(21474.83647), "bar": F64p(-0.23092222)},
		}, {
			path:   route.NewParams("foo", "hello world", "bar", "-0.23092222"),
			reader: String{"foo": new(string), "bar": new(string)},
			result: String{"foo": Sp("hello world"), "bar": Sp("-0.23092222")},
		}}

		for _, tt := range tests {
			err := tt.reader.ReadPath(tt.path)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(tt.reader, tt.result); e != nil {
				t.Error(e)
			}
		}
	})

	t.Run("QueryReaders", func(t *testing.T) {
		tests := []struct {
			query  url.Values
			reader QueryReader
			result QueryReader
			err    error
		}{{
			query:  url.Values{"foo": {"true"}},
			reader: Bool{"foo": new(bool), "bar": new(bool)},
			result: Bool{"foo": Bp(true), "bar": Bp(false)},
		}, {
			query:  url.Values{"foo": {"2147483647"}, "bar": {"-2147483648"}},
			reader: Int{"foo": new(int), "bar": new(int)},
			result: Int{"foo": Ip(2147483647), "bar": Ip(-2147483648)},
		}, {
			query:  url.Values{"foo": {"127"}, "bar": {"-128"}},
			reader: Int8{"foo": new(int8), "bar": new(int8)},
			result: Int8{"foo": I8p(127), "bar": I8p(-128)},
		}, {
			query:  url.Values{"foo": {"32767"}, "bar": {"-32768"}},
			reader: Int16{"foo": new(int16), "bar": new(int16)},
			result: Int16{"foo": I16p(32767), "bar": I16p(-32768)},
		}, {
			query:  url.Values{"foo": {"2147483647"}, "bar": {"-2147483648"}},
			reader: Int32{"foo": new(int32), "bar": new(int32)},
			result: Int32{"foo": I32p(2147483647), "bar": I32p(-2147483648)},
		}, {
			query:  url.Values{"foo": {"9223372036854775807"}, "bar": {"-9223372036854775808"}},
			reader: Int64{"foo": new(int64), "bar": new(int64)},
			result: Int64{"foo": I64p(9223372036854775807), "bar": I64p(-9223372036854775808)},
		}, {
			query:  url.Values{"foo": {"4294967295"}, "bar": {"-2147483648"}},
			reader: Uint{"foo": new(uint), "bar": new(uint)},
			result: Uint{"foo": Up(4294967295), "bar": Up(0)},
		}, {
			query:  url.Values{"foo": {"255"}, "bar": {"-128"}},
			reader: Uint8{"foo": new(uint8), "bar": new(uint8)},
			result: Uint8{"foo": U8p(255), "bar": U8p(0)},
		}, {
			query:  url.Values{"foo": {"65535"}, "bar": {"-32768"}},
			reader: Uint16{"foo": new(uint16), "bar": new(uint16)},
			result: Uint16{"foo": U16p(65535), "bar": U16p(0)},
		}, {
			query:  url.Values{"foo": {"4294967295"}, "bar": {"-2147483648"}},
			reader: Uint32{"foo": new(uint32), "bar": new(uint32)},
			result: Uint32{"foo": U32p(4294967295), "bar": U32p(0)},
		}, {
			query:  url.Values{"foo": {"18446744073709551615"}, "bar": {"-9223372036854775808"}},
			reader: Uint64{"foo": new(uint64), "bar": new(uint64)},
			result: Uint64{"foo": U64p(18446744073709551615), "bar": U64p(0)},
		}, {
			query:  url.Values{"foo": {"21474.83647"}, "bar": {"-0.23092222"}},
			reader: Float32{"foo": new(float32), "bar": new(float32)},
			result: Float32{"foo": F32p(21474.83647), "bar": F32p(-0.23092222)},
		}, {
			query:  url.Values{"foo": {"21474.83647"}, "bar": {"-0.23092222"}},
			reader: Float64{"foo": new(float64), "bar": new(float64)},
			result: Float64{"foo": F64p(21474.83647), "bar": F64p(-0.23092222)},
		}, {
			query:  url.Values{"foo": {"hello world"}, "bar": {"-0.23092222"}},
			reader: String{"foo": new(string), "bar": new(string)},
			result: String{"foo": Sp("hello world"), "bar": Sp("-0.23092222")},
		}}

		for _, tt := range tests {
			err := tt.reader.ReadQuery(tt.query)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(tt.reader, tt.result); e != nil {
				t.Error(e)
			}
		}
	})

	t.Run("HeaderReaders", func(t *testing.T) {
		tests := []struct {
			header http.Header
			reader HeaderReader
			result HeaderReader
			err    error
		}{{
			header: http.Header{"Foo": {"true"}},
			reader: Bool{"foo": new(bool), "bar": new(bool)},
			result: Bool{"foo": Bp(true), "bar": Bp(false)},
		}, {
			header: http.Header{"Foo": {"2147483647"}, "Bar": {"-2147483648"}},
			reader: Int{"foo": new(int), "bar": new(int)},
			result: Int{"foo": Ip(2147483647), "bar": Ip(-2147483648)},
		}, {
			header: http.Header{"Foo": {"127"}, "Bar": {"-128"}},
			reader: Int8{"foo": new(int8), "bar": new(int8)},
			result: Int8{"foo": I8p(127), "bar": I8p(-128)},
		}, {
			header: http.Header{"Foo": {"32767"}, "Bar": {"-32768"}},
			reader: Int16{"foo": new(int16), "bar": new(int16)},
			result: Int16{"foo": I16p(32767), "bar": I16p(-32768)},
		}, {
			header: http.Header{"Foo": {"2147483647"}, "Bar": {"-2147483648"}},
			reader: Int32{"foo": new(int32), "bar": new(int32)},
			result: Int32{"foo": I32p(2147483647), "bar": I32p(-2147483648)},
		}, {
			header: http.Header{"Foo": {"9223372036854775807"}, "Bar": {"-9223372036854775808"}},
			reader: Int64{"foo": new(int64), "bar": new(int64)},
			result: Int64{"foo": I64p(9223372036854775807), "bar": I64p(-9223372036854775808)},
		}, {
			header: http.Header{"Foo": {"4294967295"}, "Bar": {"-2147483648"}},
			reader: Uint{"foo": new(uint), "bar": new(uint)},
			result: Uint{"foo": Up(4294967295), "bar": Up(0)},
		}, {
			header: http.Header{"Foo": {"255"}, "Bar": {"-128"}},
			reader: Uint8{"foo": new(uint8), "bar": new(uint8)},
			result: Uint8{"foo": U8p(255), "bar": U8p(0)},
		}, {
			header: http.Header{"Foo": {"65535"}, "Bar": {"-32768"}},
			reader: Uint16{"foo": new(uint16), "bar": new(uint16)},
			result: Uint16{"foo": U16p(65535), "bar": U16p(0)},
		}, {
			header: http.Header{"Foo": {"4294967295"}, "Bar": {"-2147483648"}},
			reader: Uint32{"foo": new(uint32), "bar": new(uint32)},
			result: Uint32{"foo": U32p(4294967295), "bar": U32p(0)},
		}, {
			header: http.Header{"Foo": {"18446744073709551615"}, "Bar": {"-9223372036854775808"}},
			reader: Uint64{"foo": new(uint64), "bar": new(uint64)},
			result: Uint64{"foo": U64p(18446744073709551615), "bar": U64p(0)},
		}, {
			header: http.Header{"Foo": {"21474.83647"}, "Bar": {"-0.23092222"}},
			reader: Float32{"foo": new(float32), "bar": new(float32)},
			result: Float32{"foo": F32p(21474.83647), "bar": F32p(-0.23092222)},
		}, {
			header: http.Header{"Foo": {"21474.83647"}, "Bar": {"-0.23092222"}},
			reader: Float64{"foo": new(float64), "bar": new(float64)},
			result: Float64{"foo": F64p(21474.83647), "bar": F64p(-0.23092222)},
		}, {
			header: http.Header{"Foo": {"hello world"}, "Bar": {"-0.23092222"}},
			reader: String{"foo": new(string), "bar": new(string)},
			result: String{"foo": Sp("hello world"), "bar": Sp("-0.23092222")},
		}}

		for _, tt := range tests {
			err := tt.reader.ReadHeader(tt.header)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(tt.reader, tt.result); e != nil {
				t.Error(e)
			}
		}
	})
}
