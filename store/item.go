package store

import "time"

// item is store item
type item interface {
	expired(now time.Time) bool
	keyValue() (string, error)
	listValue(i int) (string, error)
	dictValue(k string) (string, error)
}

// baseItem is general store item
type baseItem struct {
	expiry time.Time
}

// expired determines if item is expired
func (bi baseItem) expired(now time.Time) bool {
	return bi.expiry.Before(now) || bi.expiry.Equal(now)
}

// keyValue is default returns keyItem value or error if item is not keyItem
func (_ baseItem) keyValue() (string, error) {
	return "", ErrNotKeyItem
}

// listValue returns listItem value or error if item is not listItem
func (_ baseItem) listValue(_ int) (string, error) {
	return "", ErrNotListItem
}

// dictValue returns dict value  or error if item is not dictItem
func (_ baseItem) dictValue(_ string) (string, error) {
	return "", ErrNotDictItem
}

// keyItem is a simple string scalar item
type keyItem struct {
	baseItem
	value string
}

// newKeyItem is a keyItem constructor
func newKeyItem(value string, expiry time.Time) keyItem {
	return keyItem{
		baseItem: baseItem{
			expiry: expiry,
		},
		value: value,
	}
}

// keyValue is default returns keyItem value
func (ki keyItem) keyValue() (string, error) {
	return ki.value, nil
}

// listItem is a strings list item
type listItem struct {
	baseItem
	list []string
}

// newListItem is a listItem constructor
func newListItem(list []string, expiry time.Time) listItem {
	return listItem{
		baseItem: baseItem{
			expiry: expiry,
		},
		list: list,
	}
}

// listValue returns listItem value by index i
func (li listItem) listValue(i int) (string, error) {
	if i < 0 {
		return "", ErrInvalidListIndex
	}
	if len(li.list)-1 < i {
		return "", ErrListIndexNotExists
	}
	return li.list[i], nil
}

// dictItem is a strings to strings map item
type dictItem struct {
	baseItem
	dict map[string]string
}

// newDictItem is a dictItem constructor
func newDictItem(dict map[string]string, expiry time.Time) dictItem {
	return dictItem{
		baseItem: baseItem{
			expiry: expiry,
		},
		dict: dict,
	}
}

// dictValue returns dictItem value by key k
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
