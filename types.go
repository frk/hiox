package hxio

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/frk/route"
)

// Bool is a map that can be used to read the values, using the map's keys,
// from an incoming request's header, path or query parameters and indirectly
// set them to the bools pointed to by the map's values.
type Bool map[string]*bool

// Bool implements the PathReader interface.
func (rr Bool) ReadPath(params route.Params) error {
	for k, v := range rr {
		*v = params.GetBool(k)
	}
	return nil
}

// Bool implements the QueryReader interface.
func (rr Bool) ReadQuery(query url.Values) error {
	for k, v := range rr {
		*v, _ = strconv.ParseBool(query.Get(k))
	}
	return nil
}

// Bool implements the HeaderReader interface.
func (rr Bool) ReadHeader(header http.Header) error {
	for k, v := range rr {
		*v, _ = strconv.ParseBool(header.Get(k))
	}
	return nil
}

// Int is a map that can be used to read the values, using the map's keys,
// from an incoming request's header, path or query parameters and indirectly
// set them to the ints pointed to by the map's values.
type Int map[string]*int

// Int implements the PathReader interface.
func (rr Int) ReadPath(params route.Params) error {
	for k, v := range rr {
		*v = params.GetInt(k)
	}
	return nil
}

// Int implements the QueryReader interface.
func (rr Int) ReadQuery(query url.Values) error {
	for k, v := range rr {
		*v, _ = strconv.Atoi(query.Get(k))
	}
	return nil
}

// Int implements the HeaderReader interface.
func (rr Int) ReadHeader(header http.Header) error {
	for k, v := range rr {
		*v, _ = strconv.Atoi(header.Get(k))
	}
	return nil
}

// Int8 is a map that can be used to read the values, using the map's keys,
// from an incoming request's header, path or query parameters and indirectly
// set them to the int8s pointed to by the map's values.
type Int8 map[string]*int8

// Int8 implements the PathReader interface.
func (rr Int8) ReadPath(params route.Params) error {
	for k, v := range rr {
		i := params.GetInt(k)
		*v = int8(i)
	}
	return nil
}

// Int8 implements the QueryReader interface.
func (rr Int8) ReadQuery(query url.Values) error {
	for k, v := range rr {
		i64, _ := strconv.ParseInt(query.Get(k), 10, 8)
		*v = int8(i64)
	}
	return nil
}

// Int8 implements the HeaderReader interface.
func (rr Int8) ReadHeader(header http.Header) error {
	for k, v := range rr {
		i64, _ := strconv.ParseInt(header.Get(k), 10, 8)
		*v = int8(i64)
	}
	return nil
}

// Int16 is a map that can be used to read the values, using the map's keys,
// from an incoming request's header, path or query parameters and indirectly
// set them to the int16s pointed to by the map's values.
type Int16 map[string]*int16

// Int16 implements the PathReader interface.
func (rr Int16) ReadPath(params route.Params) error {
	for k, v := range rr {
		i := params.GetInt(k)
		*v = int16(i)
	}
	return nil
}

// Int16 implements the QueryReader interface.
func (rr Int16) ReadQuery(query url.Values) error {
	for k, v := range rr {
		i64, _ := strconv.ParseInt(query.Get(k), 10, 16)
		*v = int16(i64)
	}
	return nil
}

// Int16 implements the HeaderReader interface.
func (rr Int16) ReadHeader(header http.Header) error {
	for k, v := range rr {
		i64, _ := strconv.ParseInt(header.Get(k), 10, 16)
		*v = int16(i64)
	}
	return nil
}

// Int32 is a map that can be used to read the values, using the map's keys,
// from an incoming request's header, path or query parameters and indirectly
// set them to the int32s pointed to by the map's values.
type Int32 map[string]*int32

// Int32 implements the PathReader interface.
func (rr Int32) ReadPath(params route.Params) error {
	for k, v := range rr {
		i := params.GetInt(k)
		*v = int32(i)
	}
	return nil
}

// Int32 implements the QueryReader interface.
func (rr Int32) ReadQuery(query url.Values) error {
	for k, v := range rr {
		i64, _ := strconv.ParseInt(query.Get(k), 10, 32)
		*v = int32(i64)
	}
	return nil
}

// Int32 implements the HeaderReader interface.
func (rr Int32) ReadHeader(header http.Header) error {
	for k, v := range rr {
		i64, _ := strconv.ParseInt(header.Get(k), 10, 32)
		*v = int32(i64)
	}
	return nil
}

// Int64 is a map that can be used to read the values, using the map's keys,
// from an incoming request's header, path or query parameters and indirectly
// set them to the int64s pointed to by the map's values.
type Int64 map[string]*int64

// Int64 implements the PathReader interface.
func (rr Int64) ReadPath(params route.Params) error {
	for k, v := range rr {
		*v = params.GetInt64(k)
	}
	return nil
}

// Int64 implements the QueryReader interface.
func (rr Int64) ReadQuery(query url.Values) error {
	for k, v := range rr {
		*v, _ = strconv.ParseInt(query.Get(k), 10, 64)
	}
	return nil
}

// Int64 implements the HeaderReader interface.
func (rr Int64) ReadHeader(header http.Header) error {
	for k, v := range rr {
		*v, _ = strconv.ParseInt(header.Get(k), 10, 64)
	}
	return nil
}

// Uint is a map that can be used to read the values, using the map's keys,
// from an incoming request's header, path or query parameters and indirectly
// set them to the uints pointed to by the map's values.
type Uint map[string]*uint

// Uint implements the PathReader interface.
func (rr Uint) ReadPath(params route.Params) error {
	for k, v := range rr {
		*v = params.GetUint(k)
	}
	return nil
}

// Uint implements the QueryReader interface.
func (rr Uint) ReadQuery(query url.Values) error {
	for k, v := range rr {
		u64, _ := strconv.ParseUint(query.Get(k), 10, 64)
		*v = uint(u64)
	}
	return nil
}

// Uint implements the HeaderReader interface.
func (rr Uint) ReadHeader(header http.Header) error {
	for k, v := range rr {
		u64, _ := strconv.ParseUint(header.Get(k), 10, 64)
		*v = uint(u64)
	}
	return nil
}

// Uint8 is a map that can be used to read the values, using the map's keys,
// from an incoming request's header, path or query parameters and indirectly
// set them to the uint8s pointed to by the map's values.
type Uint8 map[string]*uint8

// Uint8 implements the PathReader interface.
func (rr Uint8) ReadPath(params route.Params) error {
	for k, v := range rr {
		u := params.GetUint(k)
		*v = uint8(u)
	}
	return nil
}

// Uint8 implements the QueryReader interface.
func (rr Uint8) ReadQuery(query url.Values) error {
	for k, v := range rr {
		u64, _ := strconv.ParseUint(query.Get(k), 10, 8)
		*v = uint8(u64)
	}
	return nil
}

// Uint8 implements the HeaderReader interface.
func (rr Uint8) ReadHeader(header http.Header) error {
	for k, v := range rr {
		u64, _ := strconv.ParseUint(header.Get(k), 10, 8)
		*v = uint8(u64)
	}
	return nil
}

// Uint16 is a map that can be used to read the values, using the map's keys,
// from an incoming request's header, path or query parameters and indirectly
// set them to the uint16s pointed to by the map's values.
type Uint16 map[string]*uint16

// Uint16 implements the PathReader interface.
func (rr Uint16) ReadPath(params route.Params) error {
	for k, v := range rr {
		u := params.GetUint(k)
		*v = uint16(u)
	}
	return nil
}

// Uint16 implements the QueryReader interface.
func (rr Uint16) ReadQuery(query url.Values) error {
	for k, v := range rr {
		u64, _ := strconv.ParseUint(query.Get(k), 10, 16)
		*v = uint16(u64)
	}
	return nil
}

// Uint16 implements the HeaderReader interface.
func (rr Uint16) ReadHeader(header http.Header) error {
	for k, v := range rr {
		u64, _ := strconv.ParseUint(header.Get(k), 10, 16)
		*v = uint16(u64)
	}
	return nil
}

// Uint32 is a map that can be used to read the values, using the map's keys,
// from an incoming request's header, path or query parameters and indirectly
// set them to the uint32s pointed to by the map's values.
type Uint32 map[string]*uint32

// Uint32 implements the PathReader interface.
func (rr Uint32) ReadPath(params route.Params) error {
	for k, v := range rr {
		u := params.GetUint(k)
		*v = uint32(u)
	}
	return nil
}

// Uint32 implements the QueryReader interface.
func (rr Uint32) ReadQuery(query url.Values) error {
	for k, v := range rr {
		u64, _ := strconv.ParseUint(query.Get(k), 10, 32)
		*v = uint32(u64)
	}
	return nil
}

// Uint32 implements the HeaderReader interface.
func (rr Uint32) ReadHeader(header http.Header) error {
	for k, v := range rr {
		u64, _ := strconv.ParseUint(header.Get(k), 10, 32)
		*v = uint32(u64)
	}
	return nil
}

// Uint64 is a map that can be used to read the values, using the map's keys,
// from an incoming request's header, path or query parameters and indirectly
// set them to the uint64s pointed to by the map's values.
type Uint64 map[string]*uint64

// Uint64 implements the PathReader interface.
func (rr Uint64) ReadPath(params route.Params) error {
	for k, v := range rr {
		*v = params.GetUint64(k)
	}
	return nil
}

// Uint64 implements the QueryReader interface.
func (rr Uint64) ReadQuery(query url.Values) error {
	for k, v := range rr {
		*v, _ = strconv.ParseUint(query.Get(k), 10, 64)
	}
	return nil
}

// Uint64 implements the HeaderReader interface.
func (rr Uint64) ReadHeader(header http.Header) error {
	for k, v := range rr {
		*v, _ = strconv.ParseUint(header.Get(k), 10, 64)
	}
	return nil
}

// Float32 is a map that can be used to read the values, using the map's keys,
// from an incoming request's header, path or query parameters and indirectly
// set them to the float32s pointed to by the map's values.
type Float32 map[string]*float32

// Float32 implements the PathReader interface.
func (rr Float32) ReadPath(params route.Params) error {
	for k, v := range rr {
		*v = float32(params.GetFloat(k))
	}
	return nil
}

// Float32 implements the QueryReader interface.
func (rr Float32) ReadQuery(query url.Values) error {
	for k, v := range rr {
		f64, _ := strconv.ParseFloat(query.Get(k), 64)
		*v = float32(f64)
	}
	return nil
}

// Float32 implements the HeaderReader interface.
func (rr Float32) ReadHeader(header http.Header) error {
	for k, v := range rr {
		f64, _ := strconv.ParseFloat(header.Get(k), 64)
		*v = float32(f64)
	}
	return nil
}

// Float64 is a map that can be used to read the values, using the map's keys,
// from an incoming request's header, path or query parameters and indirectly
// set them to the float64s pointed to by the map's values.
type Float64 map[string]*float64

// Float64 implements the PathReader interface.
func (rr Float64) ReadPath(params route.Params) error {
	for k, v := range rr {
		*v = params.GetFloat(k)
	}
	return nil
}

// Float64 implements the QueryReader interface.
func (rr Float64) ReadQuery(query url.Values) error {
	for k, v := range rr {
		*v, _ = strconv.ParseFloat(query.Get(k), 64)
	}
	return nil
}

// Float64 implements the HeaderReader interface.
func (rr Float64) ReadHeader(header http.Header) error {
	for k, v := range rr {
		*v, _ = strconv.ParseFloat(header.Get(k), 64)
	}
	return nil
}

// String is a map that can be used to read the values, using the map's keys,
// from an incoming request's header, path or query parameters and indirectly
// set them to the strings pointed to by the map's values.
type String map[string]*string

// String implements the PathReader interface.
func (rr String) ReadPath(params route.Params) error {
	for k, v := range rr {
		*v = params.GetString(k)
	}
	return nil
}

// String implements the QueryReader interface.
func (rr String) ReadQuery(query url.Values) error {
	for k, v := range rr {
		*v = query.Get(k)
	}
	return nil
}

// String implements the HeaderReader interface.
func (rr String) ReadHeader(header http.Header) error {
	for k, v := range rr {
		*v = header.Get(k)
	}
	return nil
}
