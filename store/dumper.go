package store

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"time"
)

// Dumper is a items dumper
type Dumper interface {
	dump(items) error
	load() (items, error)
}

// path is file path string
type path string

// FileDumper is file dumper
type FileDumper path

// dump dumps items to file
func (fd FileDumper) dump(items items) error {
	f, err := os.Create(string(fd))
	if err != nil {
		return ErrFailOpenDumpFile.detailed(err.Error())
	}
	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(items); err != nil {
		return ErrFailToDumpItems.detailed(err.Error())
	}
	if err := f.Close(); err != nil {
		return ErrFailToCloseDumpFile.detailed(err.Error())
	}
	return nil
}

// load loads items from file. If file not exists it returns without error
func (fd FileDumper) load() (items, error) {
	items := items{}
	if _, err := os.Stat(string(fd)); os.IsNotExist(err) {
		return items, nil
	}
	f, err := os.Open(string(fd))
	if err != nil {
		return items, ErrFailOpenDumpFile.detailed(err.Error())
	}
	decoder := gob.NewDecoder(f)
	if err := decoder.Decode(&items); err != nil {
		return items, ErrFailToDecodeDumpFile.detailed(err.Error())
	}
	if err := f.Close(); err != nil {
		return items, ErrFailToCloseDumpFile.detailed(err.Error())
	}
	return items, nil
}

type gobItem struct {
	Expiry int64
	Value  interface{}
}

// MarshalBinary implements gob marshaling for keyItem
func (ki keyItem) MarshalBinary() ([]byte, error) {
	var out bytes.Buffer
	enc := gob.NewEncoder(&out)
	err := enc.Encode(gobItem{ki.expiry.UnixNano(), ki.value})
	if err != nil {
		return nil, fmt.Errorf(`fail to encode key item "%v": %s`, ki, err.Error())
	}
	return out.Bytes(), nil
}

// UnmarshalBinary implements gob unmarshaling for keyItem
func (ki *keyItem) UnmarshalBinary(data []byte) error {
	in := bytes.NewReader(data)
	dec := gob.NewDecoder(in)
	var i gobItem
	if err := dec.Decode(&i); err != nil {
		return fmt.Errorf("fail to decode key item")
	}
	v, ok := i.Value.(string)
	if !ok {
		return errors.New("fail to cast value to string")
	}
	ki.value = v
	ki.expiry = time.Unix(0, i.Expiry)
	return nil
}

// MarshalBinary implements gob marshaling for listItem
func (li listItem) MarshalBinary() ([]byte, error) {
	var out bytes.Buffer
	enc := gob.NewEncoder(&out)
	err := enc.Encode(gobItem{li.expiry.UnixNano(), li.list})
	if err != nil {
		return nil, fmt.Errorf(`fail to encode list item "%v": %s`, li, err.Error())
	}
	return out.Bytes(), nil
}

// UnmarshalBinary implements gob unmarshaling for listItem
func (li *listItem) UnmarshalBinary(data []byte) error {
	in := bytes.NewReader(data)
	dec := gob.NewDecoder(in)
	var i gobItem
	if err := dec.Decode(&i); err != nil {
		return fmt.Errorf("fail to decode key item")
	}
	l, ok := i.Value.([]string)
	if !ok {
		return errors.New("fail to cast value to []string")
	}
	li.list = l
	li.expiry = time.Unix(0, i.Expiry)
	return nil
}

// MarshalBinary implements gob marshaling for dictItem
func (di dictItem) MarshalBinary() ([]byte, error) {
	var out bytes.Buffer
	enc := gob.NewEncoder(&out)
	err := enc.Encode(gobItem{di.expiry.UnixNano(), di.dict})
	if err != nil {
		return nil, fmt.Errorf(`fail to encode dict item "%v": %s`, di, err.Error())
	}
	return out.Bytes(), nil
}

// UnmarshalBinary implements gob unmarshaling for dictItem
func (di *dictItem) UnmarshalBinary(data []byte) error {
	in := bytes.NewReader(data)
	dec := gob.NewDecoder(in)
	var i gobItem
	if err := dec.Decode(&i); err != nil {
		return fmt.Errorf("fail to decode key item")
	}
	d, ok := i.Value.(map[string]string)
	if !ok {
		return errors.New("fail to cast value to map[string]string")
	}
	di.dict = d
	di.expiry = time.Unix(0, i.Expiry)
	return nil
}

// init registers gob structures
func init() {
	gob.Register(map[string]string{})
	gob.Register(keyItem{})
	gob.Register(listItem{})
	gob.Register(dictItem{})
}
