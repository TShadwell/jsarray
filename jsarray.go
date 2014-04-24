//Package jsarray performs serialisation and deserialisation like the
//the Javascript array based format I've seen Google use before.
//
//For an example of the format, check the example for Marshal.
//
//It is important to note that in versions of Javascript where
//the global Array constructor is not frozen, this can be used to break
//SOP, since it can be evaluated as Javascript and the Array constructor
//can be hooked, making this insecure for some older browsers. You can
//fix this by prepending `];` or some other invalid Javascript so that
//it cannot be evaluated as javascript without removing that part.
//
//I say this, but I haven't yet seen a JSON marshaler that notes this point
//when producing arrays, so you can just ignore that if you want. Everybody else does.
package jsarray

import (
	"bytes"
	"encoding"
	"encoding/json"
	"fmt"
	. "reflect"
	"strings"
)

var rm = TypeOf(json.RawMessage(""))
var empty = []byte("")

//Function Marshal returns a []byte representing the given value in
//array encoding.
//The json tags 'string' and '-' are honoured.
//
//The format encodes structs into Javascript arrays in field order.
//Zero values are 'undefined', which is empty in Javascript
//array notation.
func Marshal(i interface{}) ([]byte, error) {
	return reduce(i)
}

func reduce(i interface{}) (j json.RawMessage, err error) {
	v := ValueOf(i)
	t := v.Type()

	if um, ok := i.(json.Marshaler); ok {
		j, err = um.MarshalJSON()
		return
	}

	if tm, ok := i.(encoding.TextMarshaler); ok {
		var t []byte
		t, err = tm.MarshalText()
		if err != nil {
			return nil, err
		}
		j, err = reduce(string(t))
		return
	}

	if DeepEqual(i, Zero(t).Interface()) {
		return []byte(""), nil
	}

	switch t.Kind() {
	case Slice, Array:
		if v.Len() == 0 {
			j = []byte("")
			return
		}
		j, err = reduceIndexable(v.Len, v.Index, nil)
	case Map:
		// map[T]T -> map[T]RawMessage
		mp := MakeMap(MapOf(t.Key(), rm))
		for _, v := range v.MapKeys() {
			var mg json.RawMessage
			mg, err = reduce(v.MapIndex(v).Interface())
			if err != nil {
				return
			}
			mp.SetMapIndex(v, ValueOf(mg))
		}
		j, err = reduce(mp.Interface())
	case Struct:
		j, err = reduceIndexable(
			v.NumField,
			v.Field,
			func(n int) string {
				return t.Field(n).Tag.Get("json")
			},
		)
	default:
		j, err = json.Marshal(i)
	}

	return
}

func reduceIndexable(Len func() int, Index func(n int) Value, tag func(n int) string) (j json.RawMessage, err error) {
	ln := Len()

	ar := make([]json.RawMessage, 0, ln)

	for i := 0; i < ln; i++ {
		var opts tagOptions
		if tag != nil {
			opts = parseTag(tag(i))
		}

		if opts.ignore == true {
			continue
		}

		var bt []byte
		bt, err = reduce(Index(i).Interface())
		if err != nil {
			return
		}
		if opts.quote {
			bt, err = json.Marshal(string(bt))
			if err != nil {
				return
			}
		}

		ar = append(ar, json.RawMessage(bt))
	}
	j = renderArray(tag == nil, ar)
	return

}

func renderArray(noprune bool, rm []json.RawMessage) json.RawMessage {
	return []byte("[" + dumpMessage(noprune, rm...) + "]")
}

func dumpMessage(noprune bool, m ...json.RawMessage) string {
	//the last value that is non-zero.
	var last int
	if noprune {
		last = len(m) - 1
	}
	sa := make([]string, len(m))
	for i := len(sa) - 1; i >= 0; i-- {
		s := string(m[i])
		switch {
		case last == 0:
			if s == "" {
				break
			} else {
				last = i
			}
			fallthrough
		default:
			sa[i] = s
		}
	}
	return strings.Join(sa[:last+1], ",")
}

var bytesComma = []byte(",")

//TODO
func increase(data []byte, vl interface{}) (err error) {
	v := ValueOf(vl)
	t := v.Type()

	if !v.CanSet() {
		err = fmt.Errorf(
			"Cannot set value of %s, is not addressable.",
			t,
		)
		return
	}

	v = Indirect(v)
	t = v.Type()

	if bytes.Equal(data, empty) {
		//leave as zero
		return
	}

	switch tp := t.Kind(); tp {
	case Slice, Array:
		//read the items of the array into v.
		//et := t.Elem()

		sl := bytes.Split(bytes.Trim(data, "[]"), bytesComma)
		ln := v.Len()
		if tp == Array && len(sl) < ln {
			err = fmt.Errorf(
				"Array %s too short, need %v elements, have %v.",
				t,
				len(sl),
				ln,
			)
			return
		}
	case Map:
	default:
	}
	return

}
