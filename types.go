package hxio

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/frk/route"
)

type Bool map[string]*bool

func (rr Bool) ReadRoute(params route.Params) error {
	for k, v := range rr {
		*v = params.GetBool(k)
	}
	return nil
}

func (rr Bool) ReadQuery(query url.Values) error {
	for k, v := range rr {
		*v, _ = strconv.ParseBool(query.Get(k))
	}
	return nil
}

func (rr Bool) ReadHeader(header http.Header) error {
	for k, v := range rr {
		*v, _ = strconv.ParseBool(header.Get(k))
	}
	return nil
}

type Int map[string]*int

func (rr Int) ReadRoute(params route.Params) error {
	for k, v := range rr {
		*v = params.GetInt(k)
	}
	return nil
}

func (rr Int) ReadQuery(query url.Values) error {
	for k, v := range rr {
		*v, _ = strconv.Atoi(query.Get(k))
	}
	return nil
}

func (rr Int) ReadHeader(header http.Header) error {
	for k, v := range rr {
		*v, _ = strconv.Atoi(header.Get(k))
	}
	return nil
}

type Int8 map[string]*int8

func (rr Int8) ReadRoute(params route.Params) error {
	for k, v := range rr {
		i := params.GetInt(k)
		*v = int8(i)
	}
	return nil
}

func (rr Int8) ReadQuery(query url.Values) error {
	for k, v := range rr {
		i64, _ := strconv.ParseInt(query.Get(k), 10, 8)
		*v = int8(i64)
	}
	return nil
}

func (rr Int8) ReadHeader(header http.Header) error {
	for k, v := range rr {
		i64, _ := strconv.ParseInt(header.Get(k), 10, 8)
		*v = int8(i64)
	}
	return nil
}

type Int16 map[string]*int16

func (rr Int16) ReadRoute(params route.Params) error {
	for k, v := range rr {
		i := params.GetInt(k)
		*v = int16(i)
	}
	return nil
}

func (rr Int16) ReadQuery(query url.Values) error {
	for k, v := range rr {
		i64, _ := strconv.ParseInt(query.Get(k), 10, 16)
		*v = int16(i64)
	}
	return nil
}

func (rr Int16) ReadHeader(header http.Header) error {
	for k, v := range rr {
		i64, _ := strconv.ParseInt(header.Get(k), 10, 16)
		*v = int16(i64)
	}
	return nil
}

type Int32 map[string]*int32

func (rr Int32) ReadRoute(params route.Params) error {
	for k, v := range rr {
		i := params.GetInt(k)
		*v = int32(i)
	}
	return nil
}

func (rr Int32) ReadQuery(query url.Values) error {
	for k, v := range rr {
		i64, _ := strconv.ParseInt(query.Get(k), 10, 32)
		*v = int32(i64)
	}
	return nil
}

func (rr Int32) ReadHeader(header http.Header) error {
	for k, v := range rr {
		i64, _ := strconv.ParseInt(header.Get(k), 10, 32)
		*v = int32(i64)
	}
	return nil
}

type Int64 map[string]*int64

func (rr Int64) ReadRoute(params route.Params) error {
	for k, v := range rr {
		*v = params.GetInt64(k)
	}
	return nil
}

func (rr Int64) ReadQuery(query url.Values) error {
	for k, v := range rr {
		*v, _ = strconv.ParseInt(query.Get(k), 10, 64)
	}
	return nil
}

func (rr Int64) ReadHeader(header http.Header) error {
	for k, v := range rr {
		*v, _ = strconv.ParseInt(header.Get(k), 10, 64)
	}
	return nil
}

type Uint map[string]*uint

func (rr Uint) ReadRoute(params route.Params) error {
	for k, v := range rr {
		*v = params.GetUint(k)
	}
	return nil
}

func (rr Uint) ReadQuery(query url.Values) error {
	for k, v := range rr {
		u64, _ := strconv.ParseUint(query.Get(k), 10, 64)
		*v = uint(u64)
	}
	return nil
}

func (rr Uint) ReadHeader(header http.Header) error {
	for k, v := range rr {
		u64, _ := strconv.ParseUint(header.Get(k), 10, 64)
		*v = uint(u64)
	}
	return nil
}

type Uint8 map[string]*uint8

func (rr Uint8) ReadRoute(params route.Params) error {
	for k, v := range rr {
		u := params.GetUint(k)
		*v = uint8(u)
	}
	return nil
}

func (rr Uint8) ReadQuery(query url.Values) error {
	for k, v := range rr {
		u64, _ := strconv.ParseUint(query.Get(k), 10, 8)
		*v = uint8(u64)
	}
	return nil
}

func (rr Uint8) ReadHeader(header http.Header) error {
	for k, v := range rr {
		u64, _ := strconv.ParseUint(header.Get(k), 10, 8)
		*v = uint8(u64)
	}
	return nil
}

type Uint16 map[string]*uint16

func (rr Uint16) ReadRoute(params route.Params) error {
	for k, v := range rr {
		u := params.GetUint(k)
		*v = uint16(u)
	}
	return nil
}

func (rr Uint16) ReadQuery(query url.Values) error {
	for k, v := range rr {
		u64, _ := strconv.ParseUint(query.Get(k), 10, 16)
		*v = uint16(u64)
	}
	return nil
}

func (rr Uint16) ReadHeader(header http.Header) error {
	for k, v := range rr {
		u64, _ := strconv.ParseUint(header.Get(k), 10, 16)
		*v = uint16(u64)
	}
	return nil
}

type Uint32 map[string]*uint32

func (rr Uint32) ReadRoute(params route.Params) error {
	for k, v := range rr {
		u := params.GetUint(k)
		*v = uint32(u)
	}
	return nil
}

func (rr Uint32) ReadQuery(query url.Values) error {
	for k, v := range rr {
		u64, _ := strconv.ParseUint(query.Get(k), 10, 32)
		*v = uint32(u64)
	}
	return nil
}

func (rr Uint32) ReadHeader(header http.Header) error {
	for k, v := range rr {
		u64, _ := strconv.ParseUint(header.Get(k), 10, 32)
		*v = uint32(u64)
	}
	return nil
}

type Uint64 map[string]*uint64

func (rr Uint64) ReadRoute(params route.Params) error {
	for k, v := range rr {
		*v = params.GetUint64(k)
	}
	return nil
}

func (rr Uint64) ReadQuery(query url.Values) error {
	for k, v := range rr {
		*v, _ = strconv.ParseUint(query.Get(k), 10, 64)
	}
	return nil
}

func (rr Uint64) ReadHeader(header http.Header) error {
	for k, v := range rr {
		*v, _ = strconv.ParseUint(header.Get(k), 10, 64)
	}
	return nil
}

type Float32 map[string]*float32

func (rr Float32) ReadRoute(params route.Params) error {
	for k, v := range rr {
		*v = float32(params.GetFloat(k))
	}
	return nil
}

func (rr Float32) ReadQuery(query url.Values) error {
	for k, v := range rr {
		f64, _ := strconv.ParseFloat(query.Get(k), 64)
		*v = float32(f64)
	}
	return nil
}

func (rr Float32) ReadHeader(header http.Header) error {
	for k, v := range rr {
		f64, _ := strconv.ParseFloat(header.Get(k), 64)
		*v = float32(f64)
	}
	return nil
}

type Float64 map[string]*float64

func (rr Float64) ReadRoute(params route.Params) error {
	for k, v := range rr {
		*v = params.GetFloat(k)
	}
	return nil
}

func (rr Float64) ReadQuery(query url.Values) error {
	for k, v := range rr {
		*v, _ = strconv.ParseFloat(query.Get(k), 64)
	}
	return nil
}

func (rr Float64) ReadHeader(header http.Header) error {
	for k, v := range rr {
		*v, _ = strconv.ParseFloat(header.Get(k), 64)
	}
	return nil
}

type String map[string]*string

func (rr String) ReadRoute(params route.Params) error {
	for k, v := range rr {
		*v = params.GetString(k)
	}
	return nil
}

func (rr String) ReadQuery(query url.Values) error {
	for k, v := range rr {
		*v = query.Get(k)
	}
	return nil
}

func (rr String) ReadHeader(header http.Header) error {
	for k, v := range rr {
		*v = header.Get(k)
	}
	return nil
}
