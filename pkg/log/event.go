package logger

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"go-source/pkg/utils"
	"net"
	"time"
)

type Event struct {
	event *zerolog.Event
}

func (e *Event) Enabled() bool {
	return e.event.Enabled()
}

func (e *Event) Discard() *Event {
	e.event.Discard()
	return e
}

func (e *Event) Msg(msg string) {
	e.event.Msg(msg)
}

func (e *Event) Send() {
	e.event.Send()
}

func (e *Event) Msgf(format string, v ...interface{}) {
	e.event.Msgf(format, v...)
}

func (e *Event) MsgFunc(createMsg func() string) {
	e.event.MsgFunc(createMsg)
}

func (e *Event) Fields(fields interface{}) *Event {
	e.event.Fields(fields)
	return e
}

func (e *Event) Dict(key string, dict *Event) *Event {
	e.event.Dict(key, dict.event)
	return e
}

func (e *Event) Array(key string, arr zerolog.LogArrayMarshaler) *Event {
	e.event.Array(key, arr)
	return e
}

func (e *Event) Object(key string, obj zerolog.LogObjectMarshaler) *Event {
	e.event.Object(key, obj)
	return e
}

func (e *Event) Func(f func(e *Event)) *Event {
	e.event.Func(func(event *zerolog.Event) {
		f(&Event{event})
	})
	return e
}

func (e *Event) EmbedObject(obj zerolog.LogObjectMarshaler) *Event {
	e.event.EmbedObject(obj)
	return e
}

func (e *Event) Str(key, val string) *Event {
	e.event.Str(key, val)
	return e
}

func (e *Event) Strs(key string, vals []string) *Event {
	e.event.Strs(key, vals)
	return e
}

func (e *Event) Stringer(key string, val fmt.Stringer) *Event {
	e.event.Stringer(key, val)
	return e
}

func (e *Event) Stringers(key string, vals []fmt.Stringer) *Event {
	e.event.Stringers(key, vals)
	return e
}

func (e *Event) Bytes(key string, val []byte) *Event {
	e.event.Bytes(key, val)
	return e
}

func (e *Event) Hex(key string, val []byte) *Event {
	e.event.Hex(key, val)
	return e
}

func (e *Event) RawJSON(key string, b []byte) *Event {
	e.event.RawJSON(key, b)
	return e
}

func (e *Event) RawCBOR(key string, b []byte) *Event {
	e.event.RawCBOR(key, b)
	return e
}

func (e *Event) AnErr(key string, err error) *Event {
	e.event.AnErr(key, err)
	return e
}

func (e *Event) Errs(key string, errs []error) *Event {
	e.event.Errs(key, errs)
	return e
}

func (e *Event) Err(err error) *Event {
	e.event.Err(err)
	return e
}

func (e *Event) Stack() *Event {
	e.event.Stack()
	return e
}

func (e *Event) Ctx(ctx context.Context) *Event {
	e.event.Ctx(ctx)
	return e
}

func (e *Event) GetCtx() context.Context {
	return e.event.GetCtx()
}

func (e *Event) Bool(key string, b bool) *Event {
	e.event.Bool(key, b)
	return e
}

func (e *Event) Bools(key string, b []bool) *Event {
	e.event.Bools(key, b)
	return e
}

func (e *Event) Int(key string, i int) *Event {
	e.event.Int(key, i)
	return e
}

func (e *Event) Ints(key string, i []int) *Event {
	e.event.Ints(key, i)
	return e
}

func (e *Event) Int8(key string, i int8) *Event {
	e.event.Int8(key, i)
	return e
}

func (e *Event) Ints8(key string, i []int8) *Event {
	e.event.Ints8(key, i)
	return e
}

func (e *Event) Int16(key string, i int16) *Event {
	e.event.Int16(key, i)
	return e
}

func (e *Event) Ints16(key string, i []int16) *Event {
	e.event.Ints16(key, i)
	return e
}

func (e *Event) Int32(key string, i int32) *Event {
	e.event.Int32(key, i)
	return e
}

func (e *Event) Ints32(key string, i []int32) *Event {
	e.event.Ints32(key, i)
	return e
}

func (e *Event) Int64(key string, i int64) *Event {
	e.event.Int64(key, i)
	return e
}

func (e *Event) Ints64(key string, i []int64) *Event {
	e.event.Ints64(key, i)
	return e
}

func (e *Event) Uint(key string, i uint) *Event {
	e.event.Uint(key, i)
	return e
}

func (e *Event) Uints(key string, i []uint) *Event {
	e.event.Uints(key, i)
	return e
}

func (e *Event) Uint8(key string, i uint8) *Event {
	e.event.Uint8(key, i)
	return e
}

func (e *Event) Uints8(key string, i []uint8) *Event {
	e.event.Uints8(key, i)
	return e
}

func (e *Event) Uint16(key string, i uint16) *Event {
	e.event.Uint16(key, i)
	return e
}

func (e *Event) Uints16(key string, i []uint16) *Event {
	e.event.Uints16(key, i)
	return e
}

func (e *Event) Uint32(key string, i uint32) *Event {
	e.event.Uint32(key, i)
	return e
}

func (e *Event) Uints32(key string, i []uint32) *Event {
	e.event.Uints32(key, i)
	return e
}

func (e *Event) Uint64(key string, i uint64) *Event {
	e.event.Uint64(key, i)
	return e
}

func (e *Event) Uints64(key string, i []uint64) *Event {
	e.event.Uints64(key, i)
	return e
}

func (e *Event) Float32(key string, f float32) *Event {
	e.event.Float32(key, f)
	return e
}

func (e *Event) Floats32(key string, f []float32) *Event {
	e.event.Floats32(key, f)
	return e
}

func (e *Event) Float64(key string, f float64) *Event {
	e.event.Float64(key, f)
	return e
}

func (e *Event) Floats64(key string, f []float64) *Event {
	e.event.Floats64(key, f)
	return e
}

func (e *Event) Timestamp() *Event {
	e.event.Timestamp()
	return e
}

func (e *Event) Time(key string, t time.Time) *Event {
	e.event.Time(key, t)
	return e
}

func (e *Event) Times(key string, t []time.Time) *Event {
	e.event.Times(key, t)
	return e
}

func (e *Event) Dur(key string, d time.Duration) *Event {
	e.event.Dur(key, d)
	return e
}

func (e *Event) Durs(key string, d []time.Duration) *Event {
	e.event.Durs(key, d)
	return e
}

func (e *Event) TimeDiff(key string, t time.Time, start time.Time) *Event {
	e.event.TimeDiff(key, t, start)
	return e
}

func (e *Event) Any(key string, i interface{}) *Event {
	e.event.Any(key, i)
	return e
}

func (e *Event) Interface(key string, i interface{}) *Event {
	e.event.Interface(key, i)
	return e
}

func (e *Event) Type(key string, val interface{}) *Event {
	e.event.Type(key, val)
	return e
}

func (e *Event) CallerSkipFrame(skip int) *Event {
	e.event.CallerSkipFrame(skip)
	return e
}

func (e *Event) Caller(skip ...int) *Event {
	e.event.Caller(skip...)
	return e
}

func (e *Event) IPAddr(key string, ip net.IP) *Event {
	e.event.IPAddr(key, ip)
	return e
}

func (e *Event) IPPrefix(key string, pfx net.IPNet) *Event {
	e.event.IPPrefix(key, pfx)
	return e
}

func (e *Event) MACAddr(key string, ha net.HardwareAddr) *Event {
	e.event.MACAddr(key, ha)
	return e
}

func (e *Event) StrEncrypt(key, val string) *Event {
	if encr, err := Encrypt(val); err == nil {
		e.event.Str(key, encr)
	} else {
		e.event.Str(key, val)
	}
	return e
}

func (e *Event) StructEncrypt(key string, val interface{}) *Event {
	if encr, err := EncryptInterface(val); err == nil {
		e.event.Interface(key, encr)
	} else {
		e.event.Interface(key, val)
	}
	return e
}

func (e *Event) StructSliceEncrypt(key string, val interface{}) *Event {
	if encr, err := EncryptInterface(val); err == nil {
		e.event.Interface(key, encr)
	} else {
		e.event.Interface(key, val)
	}
	return e
}

func (e *Event) InterfaceEncrypt(key string, val interface{}) *Event {
	if str, err := utils.AnyToString(val); err == nil {
		if encr, err := Encrypt(str); err == nil {
			e.event.Str(key, encr)
		}
	} else {
		e.event.Interface(key, val)
	}
	return e
}
