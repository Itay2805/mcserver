package minecraft

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"reflect"
	"strings"
	"sync"
)

const (
	nbtTagEnd = 0
	nbtTagByte = 1
	nbtTagShort = 2
	nbtTagInt = 3
	nbtTagLong = 4
	nbtTagFloat = 5
	nbtTagDouble = 6
	nbtTagByteArray = 7
	nbtTagString = 8
	nbtTagList = 9
	nbtTagCompound = 10
	nbtTagIntArray = 11
	nbtTagLongArray = 12
)

// stream writer (easier to use in some cases)

type NbtWriter struct {
	w 					*Writer
	hierarchy			[]uint8
	listSizeStack		[]int
	listSizeOffsetStack	[]int
}

func (writer *NbtWriter) emitString(val string) {
	writer.w.WriteShort(int16(len(val)))
	writer.w.WriteBytes([]byte(val))
}

func (writer *NbtWriter) pushListElement() {
	if len(writer.hierarchy) != 0 {
		if writer.hierarchy[len(writer.hierarchy) - 1] == nbtTagList {
			writer.listSizeStack[len(writer.listSizeStack) - 1]++
		}
	}
}

func (writer *NbtWriter) emitTagHeader(t uint8, name string) {
	if len(writer.hierarchy) == 0 || writer.hierarchy[len(writer.hierarchy) - 1] != nbtTagEnd {
		writer.w.WriteByte(t)
		writer.emitString(name)
	}
	writer.pushListElement()
}

func (writer *NbtWriter) PushByte(val int8, name string) {
	writer.emitTagHeader(nbtTagByte, name)
	writer.w.WriteByte(byte(val))
}

func (writer *NbtWriter) PushBool(val bool, name string) {
	writer.emitTagHeader(nbtTagByte, name)
	if val {
		writer.w.WriteByte(1)
	} else {
		writer.w.WriteByte(0)
	}
}


func (writer *NbtWriter) PushShort(val int16, name string) {
	writer.emitTagHeader(nbtTagByte, name)
	writer.w.WriteShort(val)
}

func (writer *NbtWriter) PushInt(val int32, name string) {
	writer.emitTagHeader(nbtTagInt, name)
	writer.w.WriteInt(val)
}

func (writer *NbtWriter) PushLong(val int64, name string) {
	writer.emitTagHeader(nbtTagLong, name)
	writer.w.WriteLong(val)
}

func (writer *NbtWriter) PushFloat(val float32, name string) {
	writer.emitTagHeader(nbtTagFloat, name)
	writer.w.WriteFloat(val)
}

func (writer *NbtWriter) PushDouble(val float64, name string) {
	writer.emitTagHeader(nbtTagDouble, name)
	writer.w.WriteDouble(val)
}

func (writer *NbtWriter) PushByteArray(data []byte, name string) {
	writer.emitTagHeader(nbtTagByteArray, name)
	writer.w.WriteInt(int32(len(data)))
	writer.w.WriteBytes(data)
}

func (writer *NbtWriter) PushIntArray(data []int32, name string) {
	writer.emitTagHeader(nbtTagIntArray, name)
	writer.w.WriteInt(int32(len(data)))
	for _, val := range data {
		writer.w.WriteInt(val)
	}
}

func (writer *NbtWriter) PushLongArray(data []int64, name string) {
	writer.emitTagHeader(nbtTagLongArray, name)
	writer.w.WriteInt(int32(len(data)))
	for _, val := range data {
		writer.w.WriteLong(val)
	}
}

func (writer *NbtWriter) PushString(val string, name string) {
	writer.emitTagHeader(nbtTagString, name)
	writer.emitString(val)
}

func (writer *NbtWriter) StartCompound(name string) {
	writer.emitTagHeader(nbtTagCompound, name)
	writer.hierarchy = append(writer.hierarchy, nbtTagCompound)
}

func (writer *NbtWriter) EndCompound() {
	writer.w.WriteByte(nbtTagEnd)
	writer.hierarchy = writer.hierarchy[:len(writer.hierarchy) - 1]
}

// TODO: List

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// nbt marshal and unmarshal
// take from https://github.com/Tnze/go-mc/tree/master/nbt
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func NbtMarshal(w io.Writer, v interface{}) error {
	return NbtNewEncoder(w).Encode(v)
}

func NbtMarshalCompound(w io.Writer, v interface{}, rootTagName string) error {
	enc := NbtNewEncoder(w)
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		log.Panicln("Must be struct!")
	}
	return enc.marshal(val, nbtTagCompound, rootTagName)
}

type NbtEncoder struct {
	w io.Writer
}

func NbtNewEncoder(w io.Writer) *NbtEncoder {
	return &NbtEncoder{w: w}
}

func (e *NbtEncoder) Encode(v interface{}) error {
	val := reflect.ValueOf(v)
	return e.marshal(val, getTagType(val.Type()), "")
}

func (e *NbtEncoder) marshal(val reflect.Value, tagType byte, tagName string) error {
	if err := e.writeHeader(val, tagType, tagName); err != nil {
		return err
	}
	return e.writeValue(val, tagType)
}

func (e *NbtEncoder) writeHeader(val reflect.Value, tagType byte, tagName string) (err error) {
	if tagType == nbtTagList {
		eleType := getTagType(val.Type().Elem())
		err = e.writeListHeader(eleType, tagName, val.Len())
	} else {
		err = e.writeTag(tagType, tagName)
	}
	return err
}

func (e *NbtEncoder) writeValue(val reflect.Value, tagType byte) error {
	switch tagType {
	default:
		return errors.New("unsupported type " + val.Type().Kind().String())
	case nbtTagByte:
		_, err := e.w.Write([]byte{byte(val.Uint())})
		return err
	case nbtTagShort:
		return e.writeInt16(int16(val.Int()))
	case nbtTagInt:
		return e.writeInt32(int32(val.Int()))
	case nbtTagFloat:
		return e.writeInt32(int32(math.Float32bits(float32(val.Float()))))
	case nbtTagLong:
		return e.writeInt64(val.Int())
	case nbtTagDouble:
		return e.writeInt64(int64(math.Float64bits(val.Float())))
	case nbtTagByteArray, nbtTagIntArray, nbtTagLongArray:
		n := val.Len()
		if err := e.writeInt32(int32(n)); err != nil {
			return err
		}

		if tagType == nbtTagByteArray {
			_, err := e.w.Write(val.Bytes())
			return err
		} else {
			for i := 0; i < n; i++ {
				v := val.Index(i).Int()

				var err error
				if tagType == nbtTagIntArray {
					err = e.writeInt32(int32(v))
				} else if tagType == nbtTagLongArray {
					err = e.writeInt64(v)
				}
				if err != nil {
					return err
				}
			}
		}

	case nbtTagList:
		for i := 0; i < val.Len(); i++ {
			arrVal := val.Index(i)
			err := e.writeValue(arrVal, getTagType(arrVal.Type()))
			if err != nil {
				return err
			}
		}

	case nbtTagString:
		if err := e.writeInt16(int16(val.Len())); err != nil {
			return err
		}
		_, err := e.w.Write([]byte(val.String()))
		return err

	case nbtTagCompound:
		if val.Kind() == reflect.Map {
			for _, key := range val.MapKeys() {
				value := val.MapIndex(key)
				err := e.marshal(value, getTagType(value.Type()), key.Interface().(string))
				if err != nil {
					return err
				}
			}
		} else {
			if val.Kind() == reflect.Interface {
				val = reflect.ValueOf(val.Interface())
			}

			n := val.NumField()
			for i := 0; i < n; i++ {
				f := val.Type().Field(i)
				tag := f.Tag.Get("nbt")
				if (f.PkgPath != "" && !f.Anonymous) || tag == "-" {
					continue // Private field
				}

				tagProps := parseTag(f, tag)
				err := e.marshal(val.Field(i), tagProps.Type, tagProps.Name)
				if err != nil {
					return err
				}
			}
		}

		_, err := e.w.Write([]byte{nbtTagEnd})
		return err
	}
	return nil
}

func getTagType(vk reflect.Type) byte {
	switch vk.Kind() {
	case reflect.Uint8:
		return nbtTagByte
	case reflect.Int16, reflect.Uint16:
		return nbtTagShort
	case reflect.Int32, reflect.Uint32:
		return nbtTagInt
	case reflect.Float32:
		return nbtTagFloat
	case reflect.Int64, reflect.Uint64:
		return nbtTagLong
	case reflect.Float64:
		return nbtTagDouble
	case reflect.String:
		return nbtTagString
	case reflect.Struct, reflect.Interface:
		return nbtTagCompound
	case reflect.Array, reflect.Slice:
		switch vk.Elem().Kind() {
		case reflect.Uint8: // Special types for these values
			return nbtTagByteArray
		case reflect.Int32:
			return nbtTagIntArray
		case reflect.Int64:
			return nbtTagLongArray
		default:
			return nbtTagList
		}
	case reflect.Map:
		return nbtTagCompound
	default:
		log.Panicln("Invalid type", vk)
	}
	return 0
}

type tagProps struct {
	Name string
	Type byte
}

func isArrayTag(ty byte) bool {
	return ty == nbtTagByteArray || ty == nbtTagIntArray || ty == nbtTagLongArray
}

func parseTag(f reflect.StructField, tagName string) tagProps {
	result := tagProps{}
	result.Name = tagName
	if result.Name == "" {
		result.Name = f.Name
	}

	nbtType := f.Tag.Get("nbt_type")
	result.Type = getTagType(f.Type)
	if strings.Contains(nbtType, "list") {
		if isArrayTag(result.Type) {
			result.Type = nbtTagList // for expanding the array to a standard list
		} else {
			panic("list is only supported for array types (byte, int, long)")
		}
	}

	return result
}

func (e *NbtEncoder) writeTag(tagType byte, tagName string) error {
	if _, err := e.w.Write([]byte{tagType}); err != nil {
		return err
	}
	bName := []byte(tagName)
	if err := e.writeInt16(int16(len(bName))); err != nil {
		return err
	}
	_, err := e.w.Write(bName)
	return err
}

func (e *NbtEncoder) writeListHeader(elementType byte, tagName string, n int) (err error) {
	if err = e.writeTag(nbtTagList, tagName); err != nil {
		return
	}
	if _, err = e.w.Write([]byte{elementType}); err != nil {
		return
	}
	// Write length of strings
	if err = e.writeInt32(int32(n)); err != nil {
		return
	}
	return nil
}

func (e *NbtEncoder) writeInt16(n int16) error {
	_, err := e.w.Write([]byte{byte(n >> 8), byte(n)})
	return err
}

func (e *NbtEncoder) writeInt32(n int32) error {
	_, err := e.w.Write([]byte{byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n)})
	return err
}

func (e *NbtEncoder) writeInt64(n int64) error {
	_, err := e.w.Write([]byte{
		byte(n >> 56), byte(n >> 48), byte(n >> 40), byte(n >> 32),
		byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n)})
	return err
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type NbtDecoderReader = interface {
	io.ByteScanner
	io.Reader
}
type NbtDecoder struct {
	r NbtDecoderReader
}

func NewNbtDecoder(r io.Reader) *NbtDecoder {
	d := new(NbtDecoder)
	if br, ok := r.(NbtDecoderReader); ok {
		d.r = br
	} else {
		d.r = bufio.NewReader(r)
	}
	return d
}

func NbtUnmarshal(data []byte, v interface{}) error {
	return NewNbtDecoder(bytes.NewReader(data)).Decode(v)
}

func (d *NbtDecoder) Decode(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return errors.New("nbt: non-pointer passed to NbtUnmarshal")
	}

	//start read NBT
	tagType, tagName, err := d.readTag()
	if err != nil {
		return fmt.Errorf("nbt: %w", err)
	}

	if c := d.checkCompressed(tagType); c != "" {
		return fmt.Errorf("nbt: unknown Tag, maybe need %s", c)
	}

	err = d.unmarshal(val.Elem(), tagType, tagName)
	if err != nil {
		return fmt.Errorf("nbt: fail to decode tag %q: %w", tagName, err)
	}
	return nil
}

// check the first byte and return if it use compress
func (d *NbtDecoder) checkCompressed(head byte) (compress string) {
	if head == 0x1f { //gzip
		compress = "gzip"
	} else if head == 0x78 { //zlib
		compress = "zlib"
	}
	return
}

// ErrEND error will be returned when reading a NBT with only Tag_End
var ErrEND = errors.New("NBT with only Tag_End")

type typeInfo struct {
	tagName     string
	nameToIndex map[string]int
}

var tInfoMap sync.Map

func getTypeInfo(typ reflect.Type) *typeInfo {
	if ti, ok := tInfoMap.Load(typ); ok {
		return ti.(*typeInfo)
	}

	tInfo := new(typeInfo)
	tInfo.nameToIndex = make(map[string]int)
	if typ.Kind() == reflect.Struct {
		n := typ.NumField()
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			tag := f.Tag.Get("nbt")
			if (f.PkgPath != "" && !f.Anonymous) || tag == "-" {
				continue // Private field
			}

			tInfo.nameToIndex[tag] = i
			if _, ok := tInfo.nameToIndex[f.Name]; !ok {
				tInfo.nameToIndex[f.Name] = i
			}
		}
	}

	ti, _ := tInfoMap.LoadOrStore(typ, tInfo)
	return ti.(*typeInfo)
}

func (t *typeInfo) findIndexByName(name string) int {
	i, ok := t.nameToIndex[name]
	if !ok {
		return -1
	}
	return i
}

func (d *NbtDecoder) unmarshal(val reflect.Value, tagType byte, tagName string) error {
	switch tagType {
	default:
		return fmt.Errorf("unknown Tag 0x%02x", tagType)

	case nbtTagEnd:
		return ErrEND

	case nbtTagByte:
		value, err := d.r.ReadByte()
		if err != nil {
			return err
		}
		switch vk := val.Kind(); vk {
		default:
			return errors.New("cannot parse TagByte as " + vk.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val.SetInt(int64(value))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val.SetUint(uint64(value))
		case reflect.Interface:
			val.Set(reflect.ValueOf(value))
		}

	case nbtTagShort:
		value, err := d.readInt16()
		if err != nil {
			return err
		}
		switch vk := val.Kind(); vk {
		default:
			return errors.New("cannot parse TagShort as " + vk.String())
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			val.SetInt(int64(value))
		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val.SetUint(uint64(value))
		case reflect.Interface:
			val.Set(reflect.ValueOf(value))
		}

	case nbtTagInt:
		value, err := d.readInt32()
		if err != nil {
			return err
		}
		switch vk := val.Kind(); vk {
		default:
			return errors.New("cannot parse TagInt as " + vk.String())
		case reflect.Int, reflect.Int32, reflect.Int64:
			val.SetInt(int64(value))
		case reflect.Uint, reflect.Uint32, reflect.Uint64:
			val.SetUint(uint64(value))
		case reflect.Interface:
			val.Set(reflect.ValueOf(value))
		}

	case nbtTagFloat:
		vInt, err := d.readInt32()
		if err != nil {
			return err
		}
		value := math.Float32frombits(uint32(vInt))
		switch vk := val.Kind(); vk {
		default:
			return errors.New("cannot parse TagFloat as " + vk.String())
		case reflect.Float32:
			val.Set(reflect.ValueOf(value))
		case reflect.Float64:
			val.Set(reflect.ValueOf(float64(value)))
		case reflect.Interface:
			val.Set(reflect.ValueOf(value))
		}

	case nbtTagLong:
		value, err := d.readInt64()
		if err != nil {
			return err
		}
		switch vk := val.Kind(); vk {
		default:
			return errors.New("cannot parse TagLong as " + vk.String())
		case reflect.Int, reflect.Int64:
			val.SetInt(int64(value))
		case reflect.Uint, reflect.Uint64:
			val.SetUint(uint64(value))
		case reflect.Interface:
			val.Set(reflect.ValueOf(value))
		}

	case nbtTagDouble:
		vInt, err := d.readInt64()
		if err != nil {
			return err
		}
		value := math.Float64frombits(uint64(vInt))

		switch vk := val.Kind(); vk {
		default:
			return errors.New("cannot parse TagDouble as " + vk.String())
		case reflect.Float64:
			val.Set(reflect.ValueOf(value))
		case reflect.Interface:
			val.Set(reflect.ValueOf(value))
		}

	case nbtTagString:
		s, err := d.readString()
		if err != nil {
			return err
		}
		switch vk := val.Kind(); vk {
		default:
			return errors.New("cannot parse TagString as " + vk.String())
		case reflect.String:
			val.SetString(s)
		case reflect.Interface:
			val.Set(reflect.ValueOf(s))
		}

	case nbtTagByteArray:
		aryLen, err := d.readInt32()
		if err != nil {
			return err
		}
		if aryLen < 0 {
			return errors.New("byte array len less than 0")
		}
		ba := make([]byte, aryLen)
		if _, err = io.ReadFull(d.r, ba); err != nil {
			return err
		}

		switch vt := val.Type(); {
		default:
			return errors.New("cannot parse TagByteArray to " + vt.String() + ", use []byte in this instance")
		case vt == reflect.TypeOf(ba):
			val.SetBytes(ba)
		case vt.Kind() == reflect.Interface:
			val.Set(reflect.ValueOf(ba))
		}

	case nbtTagIntArray:
		aryLen, err := d.readInt32()
		if err != nil {
			return err
		}
		vt := val.Type() //receiver must be []int or []int32
		if vt.Kind() == reflect.Interface {
			vt = reflect.TypeOf([]int32{}) // pass
		} else if vt.Kind() != reflect.Slice {
			return errors.New("cannot parse TagIntArray to " + vt.String() + ", it must be a slice")
		} else if tk := val.Type().Elem().Kind(); tk != reflect.Int && tk != reflect.Int32 {
			return errors.New("cannot parse TagIntArray to " + vt.String())
		}

		buf := reflect.MakeSlice(vt, int(aryLen), int(aryLen))
		for i := 0; i < int(aryLen); i++ {
			value, err := d.readInt32()
			if err != nil {
				return err
			}
			buf.Index(i).SetInt(int64(value))
		}
		val.Set(buf)

	case nbtTagLongArray:
		aryLen, err := d.readInt32()
		if err != nil {
			return err
		}
		vt := val.Type() //receiver must be []int or []int64
		if vt.Kind() == reflect.Interface {
			vt = reflect.TypeOf([]int64{}) // pass
		} else if vt.Kind() != reflect.Slice {
			return errors.New("cannot parse TagLongArray to " + vt.String() + ", it must be a slice")
		} else if val.Type().Elem().Kind() != reflect.Int64 {
			return errors.New("cannot parse TagLongArray to " + vt.String())
		}

		buf := reflect.MakeSlice(vt, int(aryLen), int(aryLen))
		for i := 0; i < int(aryLen); i++ {
			value, err := d.readInt64()
			if err != nil {
				return err
			}
			buf.Index(i).SetInt(value)
		}
		val.Set(buf)

	case nbtTagList:
		listType, err := d.r.ReadByte()
		if err != nil {
			return err
		}
		listLen, err := d.readInt32()
		if err != nil {
			return err
		}
		if listLen < 0 {
			return errors.New("list length less than 0")
		}

		// If we need parse TAG_List into slice, make a new with right length.
		// Otherwise if we need parse into array, we check if len(array) are enough.
		var buf reflect.Value
		vk := val.Kind()
		switch vk {
		default:
			return errors.New("cannot parse TagList as " + vk.String())
		case reflect.Interface:
			buf = reflect.ValueOf(make([]interface{}, listLen))
		case reflect.Slice:
			buf = reflect.MakeSlice(val.Type(), int(listLen), int(listLen))
		case reflect.Array:
			if vl := val.Len(); vl < int(listLen) {
				return fmt.Errorf(
					"TagList %s has len %d, but array %v only has len %d",
					tagName, listLen, val.Type(), vl)
			}
			buf = val
		}
		for i := 0; i < int(listLen); i++ {
			if err := d.unmarshal(buf.Index(i), listType, ""); err != nil {
				return err
			}
		}

		if vk != reflect.Array {
			val.Set(buf)
		}

	case nbtTagCompound:
		switch vk := val.Kind(); vk {
		default:
			return errors.New("cannot parse TagCompound as " + vk.String())
		case reflect.Struct:
			tinfo := getTypeInfo(val.Type())
			for {
				tt, tn, err := d.readTag()
				if err != nil {
					return err
				}
				if tt == nbtTagEnd {
					break
				}
				field := tinfo.findIndexByName(tn)
				if field != -1 {
					err = d.unmarshal(val.Field(field), tt, tn)
					if err != nil {
						return fmt.Errorf("fail to decode tag %q: %w", tn, err)
					}
				} else {
					if err := d.rawRead(tt); err != nil {
						return err
					}
				}
			}
		case reflect.Map:
			if val.Type().Key().Kind() != reflect.String {
				return errors.New("cannot parse TagCompound as " + val.Type().String())
			}
			if val.IsNil() {
				val.Set(reflect.MakeMap(val.Type()))
			}
			for {
				tt, tn, err := d.readTag()
				if err != nil {
					return err
				}
				if tt == nbtTagEnd {
					break
				}
				v := reflect.New(val.Type().Elem())
				if err = d.unmarshal(v.Elem(), tt, tn); err != nil {
					return fmt.Errorf("fail to decode tag %q: %w", tn, err)
				}
				val.SetMapIndex(reflect.ValueOf(tn), v.Elem())
			}
		case reflect.Interface:
			buf := make(map[string]interface{})
			for {
				tt, tn, err := d.readTag()
				if err != nil {
					return err
				}
				if tt == nbtTagEnd {
					break
				}
				var value interface{}
				if err = d.unmarshal(reflect.ValueOf(&value).Elem(), tt, tn); err != nil {
					return fmt.Errorf("fail to decode tag %q: %w", tn, err)
				}
				buf[tn] = value
			}
			val.Set(reflect.ValueOf(buf))
		}
	}

	return nil
}

func (d *NbtDecoder) rawRead(tagType byte) error {
	var buf [8]byte
	switch tagType {
	default:
		return fmt.Errorf("unknown to read 0x%02x", tagType)
	case nbtTagByte:
		_, err := d.r.ReadByte()
		return err
	case nbtTagString:
		_, err := d.readString()
		return err
	case nbtTagShort:
		_, err := io.ReadFull(d.r, buf[:2])
		return err
	case nbtTagInt, nbtTagFloat:
		_, err := io.ReadFull(d.r, buf[:4])
		return err
	case nbtTagLong, nbtTagDouble:
		_, err := io.ReadFull(d.r, buf[:8])
		return err
	case nbtTagByteArray:
		aryLen, err := d.readInt32()
		if err != nil {
			return err
		}

		if _, err = io.CopyN(ioutil.Discard, d.r, int64(aryLen)); err != nil {
			return err
		}
	case nbtTagIntArray:
		aryLen, err := d.readInt32()
		if err != nil {
			return err
		}
		for i := 0; i < int(aryLen); i++ {
			if _, err := d.readInt32(); err != nil {
				return err
			}
		}

	case nbtTagLongArray:
		aryLen, err := d.readInt32()
		if err != nil {
			return err
		}
		for i := 0; i < int(aryLen); i++ {
			if _, err := d.readInt64(); err != nil {
				return err
			}
		}

	case nbtTagList:
		listType, err := d.r.ReadByte()
		if err != nil {
			return err
		}
		listLen, err := d.readInt32()
		if err != nil {
			return err
		}
		for i := 0; i < int(listLen); i++ {
			if err := d.rawRead(listType); err != nil {
				return err
			}
		}
	case nbtTagCompound:
		for {
			tt, _, err := d.readTag()
			if err != nil {
				return err
			}
			if tt == nbtTagEnd {
				break
			}
			err = d.rawRead(tt)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *NbtDecoder) readTag() (tagType byte, tagName string, err error) {
	tagType, err = d.r.ReadByte()
	if err != nil {
		return
	}

	if tagType != nbtTagEnd { //Read Tag
		tagName, err = d.readString()
	}
	return
}

func (d *NbtDecoder) readInt16() (int16, error) {
	var data [2]byte
	_, err := io.ReadFull(d.r, data[:])
	return int16(data[0])<<8 | int16(data[1]), err
}

func (d *NbtDecoder) readInt32() (int32, error) {
	var data [4]byte
	_, err := io.ReadFull(d.r, data[:])
	return int32(data[0])<<24 | int32(data[1])<<16 |
		int32(data[2])<<8 | int32(data[3]), err
}

func (d *NbtDecoder) readInt64() (int64, error) {
	var data [8]byte
	_, err := io.ReadFull(d.r, data[:])
	return int64(data[0])<<56 | int64(data[1])<<48 |
		int64(data[2])<<40 | int64(data[3])<<32 |
		int64(data[4])<<24 | int64(data[5])<<16 |
		int64(data[6])<<8 | int64(data[7]), err
}

func (d *NbtDecoder) readString() (string, error) {
	length, err := d.readInt16()
	if err != nil {
		return "", err
	} else if length < 0 {
		return "", errors.New("string length less than 0")
	}

	var str string
	if length > 0 {
		buf := make([]byte, length)
		_, err = io.ReadFull(d.r, buf)
		str = string(buf)
	}
	return str, err
}
