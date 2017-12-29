package store

import "time"

type item interface {
	expired(now time.Time) bool
	keyValue() (string, error)
	listValue(i int) (string, error)
	dictValue(k string) (string, error)
}

type baseItem struct {
	expiry time.Time
}

func (bi baseItem) expired(now time.Time) bool {
	return bi.expiry.Before(now) || bi.expiry.Equal(now)
}

func (_ baseItem) keyValue() (string, error) {
	return "", ErrNotKeyItem
}

func (_ baseItem) listValue(_ int) (string, error) {
	return "", ErrNotListItem
}

func (_ baseItem) dictValue(_ string) (string, error) {
	return "", ErrNotDictItem
}

type keyItem struct {
	baseItem
	value string
}

func newKeyItem(value string, expiry time.Time) keyItem {
	return keyItem{
		baseItem: baseItem{
			expiry: expiry,
		},
		value: value,
	}
}

func (ki keyItem) keyValue() (string, error) {
	return ki.value, nil
}

type listItem struct {
	baseItem
	list []string
}

func newListItem(list []string, expiry time.Time) listItem {
	return listItem{
		baseItem: baseItem{
			expiry: expiry,
		},
		list: list,
	}
}

func (li listItem) listValue(i int) (string, error) {
	if i < 0 {
		return "", ErrInvalidListIndex
	}
	if len(li.list)-1 < i {
		return "", ErrListIndexNotExists
	}
	return li.list[i], nil
}

type dictItem struct {
	baseItem
	dict map[string]string
}

func newDictItem(dict map[string]string, expiry time.Time) dictItem {
	return dictItem{
		baseItem: baseItem{
			expiry: expiry,
		},
		dict: dict,
	}
}

func (di dictItem) dictValue(k string) (string, error) {
	if di.dict == nil {
		return "", ErrDictKeyNotExists
	}
	v, exists := di.dict[k]
	if !exists {
		return "", ErrDictKeyNotExists
	}
	return v, nil
}
