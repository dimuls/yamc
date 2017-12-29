package store

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("baseItem", func() {
	Specify("keyValue error", func() {
		_, err := baseItem{}.keyValue()
		Expect(err).To(MatchError(ErrNotKeyItem))
	})
	Specify("listValue error", func() {
		_, err := baseItem{}.listValue(0)
		Expect(err).To(MatchError(ErrNotListItem))
	})
	Specify("dictValue error", func() {
		_, err := baseItem{}.dictValue("")
		Expect(err).To(MatchError(ErrNotDictItem))
	})
	Describe("expired", func() {
		var t time.Time
		BeforeEach(func() {
			t = time.Now()
		})
		Specify("when expiry equals now", func() {
			Expect(baseItem{expiry: t}.expired(t)).To(BeTrue())
		})
		Specify("when expiry less than now", func() {
			t := time.Now()
			Expect(baseItem{expiry: t}.expired(t.Add(time.Nanosecond))).To(BeTrue())
		})
		Specify("when expiry greater than now", func() {
			t := time.Now()
			Expect(baseItem{expiry: t}.expired(t.Add(-time.Nanosecond))).To(BeFalse())
		})
	})
})

var _ = Describe("newKeyItem", func() {
	Specify("empty value", func() {
		ki := newKeyItem("", time.Time{})
		Expect(ki.baseItem.expiry).To(Equal(time.Time{}))
		Expect(ki.value).To(Equal(""))
	})
	Specify("not empty value", func() {
		ki := newKeyItem("a", time.Time{})
		Expect(ki.baseItem.expiry).To(Equal(time.Time{}))
		Expect(ki.value).To(Equal("a"))
	})
})

var _ = Describe("keyItem", func() {
	Specify("keyValue returns value", func() {
		Expect(keyItem{value: "test-value"}.keyValue()).To(Equal("test-value"))
	})
	Specify("listValue error", func() {
		_, err := keyItem{}.listValue(0)
		Expect(err).To(MatchError(ErrNotListItem))
	})
	Specify("dictValue error", func() {
		_, err := keyItem{}.dictValue("")
		Expect(err).To(MatchError(ErrNotDictItem))
	})
	Describe("expired", func() {
		var t time.Time
		BeforeEach(func() {
			t = time.Now()
		})
		Specify("when expiry equals now", func() {
			Expect(newKeyItem("", t).expired(t)).To(BeTrue())
		})
		Specify("when expiry less than now", func() {
			t := time.Now()
			Expect(newKeyItem("", t).expired(t.Add(time.Nanosecond))).To(BeTrue())
		})
		Specify("when expiry greater than now", func() {
			t := time.Now()
			Expect(newKeyItem("", t).expired(t.Add(-time.Nanosecond))).To(BeFalse())
		})
	})
})

var _ = Describe("newListItem", func() {
	Specify("nil list", func() {
		li := newListItem(nil, time.Time{})
		Expect(li.baseItem.expiry).To(Equal(time.Time{}))
		Expect(li.list).To(BeNil())
	})
	Specify("empty list", func() {
		li := newListItem([]string{}, time.Time{})
		Expect(li.baseItem.expiry).To(Equal(time.Time{}))
		Expect(li.list).To(Equal([]string{}))
	})
	Specify("not empty list", func() {
		li := newListItem([]string{"a", "b"}, time.Time{})
		Expect(li.baseItem.expiry).To(Equal(time.Time{}))
		Expect(li.list).To(Equal([]string{"a", "b"}))
	})
})

var _ = Describe("listItem", func() {
	Specify("listItem returns value", func() {
		li := newListItem([]string{"a", "b"}, time.Time{})
		Expect(li.listValue(0)).To(Equal("a"))
		Expect(li.listValue(1)).To(Equal("b"))
	})
	Describe("listValue error", func() {
		Specify("when nil list", func() {
			li := newListItem(nil, time.Time{})
			_, err := li.listValue(2)
			Expect(err).To(MatchError(ErrListIndexNotExists))
		})
		Specify("when empty list", func() {
			li := newListItem([]string{}, time.Time{})
			_, err := li.listValue(2)
			Expect(err).To(MatchError(ErrListIndexNotExists))
		})
		Specify("when too big index", func() {
			li := newListItem([]string{"a", "b"}, time.Time{})
			_, err := li.listValue(2)
			Expect(err).To(MatchError(ErrListIndexNotExists))
		})
		Specify("when negative index", func() {
			li := newListItem([]string{"a", "b"}, time.Time{})
			_, err := li.listValue(-1)
			Expect(err).To(MatchError(ErrInvalidListIndex))
		})
	})
	Specify("keyValue error", func() {
		_, err := listItem{}.keyValue()
		Expect(err).To(MatchError(ErrNotKeyItem))
	})
	Specify("dictValue error", func() {
		_, err := listItem{}.dictValue("")
		Expect(err).To(MatchError(ErrNotDictItem))
	})
	Describe("expired", func() {
		var t time.Time
		BeforeEach(func() {
			t = time.Now()
		})
		Specify("when expiry equals now", func() {
			Expect(newListItem(nil, t).expired(t)).To(BeTrue())
		})
		Specify("when expiry less than now", func() {
			t := time.Now()
			Expect(newListItem(nil, t).expired(t.Add(time.Nanosecond))).To(BeTrue())
		})
		Specify("when expiry greater than now", func() {
			t := time.Now()
			Expect(newListItem(nil, t).expired(t.Add(-time.Nanosecond))).To(BeFalse())
		})
	})
})

var _ = Describe("newDictItem", func() {
	Specify("nil dict", func() {
		li := newDictItem(nil, time.Time{})
		Expect(li.baseItem.expiry).To(Equal(time.Time{}))
		Expect(li.dict).To(BeNil())
	})
	Specify("empty dict", func() {
		li := newDictItem(map[string]string{}, time.Time{})
		Expect(li.baseItem.expiry).To(Equal(time.Time{}))
		Expect(li.dict).To(Equal(map[string]string{}))
	})
	Specify("not empty dict", func() {
		li := newDictItem(map[string]string{"a": "b"}, time.Time{})
		Expect(li.baseItem.expiry).To(Equal(time.Time{}))
		Expect(li.dict).To(Equal(map[string]string{"a": "b"}))
	})
})

var _ = Describe("dictItem", func() {
	Specify("dictItem returns value", func() {
		li := newDictItem(map[string]string{"a": "b"}, time.Time{})
		Expect(li.dictValue("a")).To(Equal("b"))
	})
	Describe("dictValue error", func() {
		Specify("when nil dict", func() {
			li := newDictItem(nil, time.Time{})
			_, err := li.dictValue("a")
			Expect(err).To(MatchError(ErrDictKeyNotExists))
		})
		Specify("when empty dict", func() {
			li := newDictItem(map[string]string{}, time.Time{})
			_, err := li.dictValue("a")
			Expect(err).To(MatchError(ErrDictKeyNotExists))
		})
		Specify("when dict key not exists", func() {
			li := newDictItem(map[string]string{"a": "b"}, time.Time{})
			_, err := li.dictValue("c")
			Expect(err).To(MatchError(ErrDictKeyNotExists))
		})
	})
	Specify("keyValue error", func() {
		_, err := listItem{}.keyValue()
		Expect(err).To(MatchError(ErrNotKeyItem))
	})
	Specify("listValue error", func() {
		_, err := keyItem{}.listValue(0)
		Expect(err).To(MatchError(ErrNotListItem))
	})
	Describe("expired", func() {
		var t time.Time
		BeforeEach(func() {
			t = time.Now()
		})
		Specify("when expiry equals now", func() {
			Expect(newDictItem(nil, t).expired(t)).To(BeTrue())
		})
		Specify("when expiry less than now", func() {
			t := time.Now()
			Expect(newDictItem(nil, t).expired(t.Add(time.Nanosecond))).To(BeTrue())
		})
		Specify("when expiry greater than now", func() {
			t := time.Now()
			Expect(newDictItem(nil, t).expired(t.Add(-time.Nanosecond))).To(BeFalse())
		})
	})
})
