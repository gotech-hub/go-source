package logger

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"net"
	"time"
)

type Context struct {
	l Logger
}

func (c Context) Logger() Logger {
	return Logger{logger: c.l.logger}
}

func (c Context) Fields(fields interface{}) Context {
	return Context{Logger{c.l.logger.With().Fields(fields).Logger()}}
}

func (c Context) Array(key string, arr zerolog.LogArrayMarshaler) Context {
	return Context{Logger{c.l.logger.With().Array(key, arr).Logger()}}
}

func (c Context) Object(key string, obj zerolog.LogObjectMarshaler) Context {
	return Context{Logger{c.l.logger.With().Object(key, obj).Logger()}}
}

func (c Context) EmbedObject(obj zerolog.LogObjectMarshaler) Context {
	return Context{Logger{c.l.logger.With().EmbedObject(obj).Logger()}}
}

func (c Context) Str(key, val string) Context {
	return Context{Logger{c.l.logger.With().Str(key, val).Logger()}}
}

func (c Context) Strs(key string, vals []string) Context {
	return Context{Logger{c.l.logger.With().Strs(key, vals).Logger()}}
}

func (c Context) Stringer(key string, val fmt.Stringer) Context {
	return Context{Logger{c.l.logger.With().Stringer(key, val).Logger()}}
}

func (c Context) Bytes(key string, val []byte) Context {
	return Context{Logger{c.l.logger.With().Bytes(key, val).Logger()}}
}

func (c Context) Hex(key string, val []byte) Context {
	return Context{Logger{c.l.logger.With().Hex(key, val).Logger()}}
}

func (c Context) RawJSON(key string, b []byte) Context {
	return Context{Logger{c.l.logger.With().RawJSON(key, b).Logger()}}
}

func (c Context) AnErr(key string, err error) Context {
	return Context{Logger{c.l.logger.With().AnErr(key, err).Logger()}}
}

func (c Context) Errs(key string, errs []error) Context {
	return Context{Logger{c.l.logger.With().Errs(key, errs).Logger()}}
}

func (c Context) Err(err error) Context {
	return Context{Logger{c.l.logger.With().Err(err).Logger()}}
}

func (c Context) Ctx(ctx context.Context) Context {
	return Context{Logger{c.l.logger.With().Ctx(ctx).Logger()}}
}

func (c Context) Bool(key string, b bool) Context {
	return Context{Logger{c.l.logger.With().Bool(key, b).Logger()}}
}

func (c Context) Bools(key string, b []bool) Context {
	return Context{Logger{c.l.logger.With().Bools(key, b).Logger()}}
}

func (c Context) Int(key string, i int) Context {
	return Context{Logger{c.l.logger.With().Int(key, i).Logger()}}
}

func (c Context) Ints(key string, i []int) Context {
	return Context{Logger{c.l.logger.With().Ints(key, i).Logger()}}
}

func (c Context) Int8(key string, i int8) Context {
	return Context{Logger{c.l.logger.With().Int8(key, i).Logger()}}
}

func (c Context) Ints8(key string, i []int8) Context {
	return Context{Logger{c.l.logger.With().Ints8(key, i).Logger()}}
}

func (c Context) Int16(key string, i int16) Context {
	return Context{Logger{c.l.logger.With().Int16(key, i).Logger()}}
}

func (c Context) Ints16(key string, i []int16) Context {
	return Context{Logger{c.l.logger.With().Ints16(key, i).Logger()}}
}

func (c Context) Int32(key string, i int32) Context {
	return Context{Logger{c.l.logger.With().Int32(key, i).Logger()}}
}

func (c Context) Ints32(key string, i []int32) Context {
	return Context{Logger{c.l.logger.With().Ints32(key, i).Logger()}}
}

func (c Context) Int64(key string, i int64) Context {
	return Context{Logger{c.l.logger.With().Int64(key, i).Logger()}}
}

func (c Context) Ints64(key string, i []int64) Context {
	return Context{Logger{c.l.logger.With().Ints64(key, i).Logger()}}
}

func (c Context) Uint(key string, i uint) Context {
	return Context{Logger{c.l.logger.With().Uint(key, i).Logger()}}
}

func (c Context) Uints(key string, i []uint) Context {
	return Context{Logger{c.l.logger.With().Uints(key, i).Logger()}}
}

func (c Context) Uint8(key string, i uint8) Context {
	return Context{Logger{c.l.logger.With().Uint8(key, i).Logger()}}
}

func (c Context) Uints8(key string, i []uint8) Context {
	return Context{Logger{c.l.logger.With().Uints8(key, i).Logger()}}
}

func (c Context) Uint16(key string, i uint16) Context {
	return Context{Logger{c.l.logger.With().Uint16(key, i).Logger()}}
}

func (c Context) Uints16(key string, i []uint16) Context {
	return Context{Logger{c.l.logger.With().Uints16(key, i).Logger()}}
}

func (c Context) Uint32(key string, i uint32) Context {
	return Context{Logger{c.l.logger.With().Uint32(key, i).Logger()}}
}

func (c Context) Uints32(key string, i []uint32) Context {
	return Context{Logger{c.l.logger.With().Uints32(key, i).Logger()}}
}

func (c Context) Uint64(key string, i uint64) Context {
	return Context{Logger{c.l.logger.With().Uint64(key, i).Logger()}}
}

func (c Context) Uints64(key string, i []uint64) Context {
	return Context{Logger{c.l.logger.With().Uints64(key, i).Logger()}}
}

func (c Context) Float32(key string, f float32) Context {
	return Context{Logger{c.l.logger.With().Float32(key, f).Logger()}}
}

func (c Context) Floats32(key string, f []float32) Context {
	return Context{Logger{c.l.logger.With().Floats32(key, f).Logger()}}
}

func (c Context) Float64(key string, f float64) Context {
	return Context{Logger{c.l.logger.With().Float64(key, f).Logger()}}
}

func (c Context) Floats64(key string, f []float64) Context {
	return Context{Logger{c.l.logger.With().Floats64(key, f).Logger()}}
}

func (c Context) Timestamp() Context {
	return Context{Logger{c.l.logger.With().Timestamp().Logger()}}
}

func (c Context) Time(key string, t time.Time) Context {
	return Context{Logger{c.l.logger.With().Time(key, t).Logger()}}
}

func (c Context) Times(key string, t []time.Time) Context {
	return Context{Logger{c.l.logger.With().Times(key, t).Logger()}}
}

func (c Context) Dur(key string, d time.Duration) Context {
	return Context{Logger{c.l.logger.With().Dur(key, d).Logger()}}
}

func (c Context) Durs(key string, d []time.Duration) Context {
	return Context{Logger{c.l.logger.With().Durs(key, d).Logger()}}
}

func (c Context) Interface(key string, i interface{}) Context {
	return Context{Logger{c.l.logger.With().Interface(key, i).Logger()}}
}

func (c Context) Type(key string, val interface{}) Context {
	return Context{Logger{c.l.logger.With().Type(key, val).Logger()}}
}

func (c Context) Any(key string, i interface{}) Context {
	return Context{Logger{c.l.logger.With().Any(key, i).Logger()}}
}

func (c Context) Caller() Context {
	return Context{Logger{c.l.logger.With().Caller().Logger()}}
}

func (c Context) CallerWithSkipFrameCount(skipFrameCount int) Context {
	return Context{Logger{c.l.logger.With().CallerWithSkipFrameCount(skipFrameCount).Logger()}}
}

func (c Context) Stack() Context {
	return Context{Logger{c.l.logger.With().Stack().Logger()}}
}

func (c Context) IPAddr(key string, ip net.IP) Context {
	return Context{Logger{c.l.logger.With().IPAddr(key, ip).Logger()}}
}

func (c Context) IPPrefix(key string, pfx net.IPNet) Context {
	return Context{Logger{c.l.logger.With().IPPrefix(key, pfx).Logger()}}
}

func (c Context) MACAddr(key string, ha net.HardwareAddr) Context {
	return Context{Logger{c.l.logger.With().MACAddr(key, ha).Logger()}}
}
