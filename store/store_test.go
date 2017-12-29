package store

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type testClock time.Time

func (tc testClock) Now() time.Time {
	return time.Time(tc)
}

var _ = Describe("NewStore", func() {
	Specify("invalid params error", func() {
		_, err := NewStore(Params{CleaningPeriod: 100*time.Millisecond - time.Nanosecond}, testClock{})
		Expect(err).To(MatchError(ErrInvalidParams.detailed("too small cleaning period")))
	})
	Specify("nil clock error", func() {
		_, err := NewStore(Params{CleaningPeriod: 100 * time.Millisecond}, nil)
		Expect(err).To(MatchError(ErrNilClock))
	})
	Specify("succeeds", func() {
		p := Params{CleaningPeriod: 100 * time.Millisecond}
		c := testClock{}
		si, err := NewStore(p, c)
		s := si.(*store)
		Expect(err).ToNot(HaveOccurred())
		Expect(s).ToNot(BeNil())
		Expect(s.params).To(Equal(p))
		Expect(s.clock).To(Equal(c))
		Expect(s.items).ToNot(BeNil())
		Expect(s.items).To(BeEmpty())
		Expect(s.cleaner).To(BeNil())
	})
})

var _ = Describe("store", func() {
	var (
		c Clock
		s *store
	)
	BeforeEach(func() {
		c = testClock(time.Now())
		si, err := NewStore(Params{CleaningPeriod: 100 * time.Millisecond}, c)
		Expect(err).ToNot(HaveOccurred())
		s = si.(*store)
	})
	Describe("Get", func() {
		Specify("expired item error", func() {
			s.items["a"] = newKeyItem("a", time.Now())
			_, err := s.Get("a")
			Expect(err).To(MatchError(ErrKeyNotExists))
		})
		Specify("not existed item error", func() {
			_, err := s.Get("a")
			Expect(err).To(MatchError(ErrKeyNotExists))
		})
		Specify("not key item error", func() {
			s.items["a"] = newListItem(nil, time.Now().Add(time.Nanosecond))
			_, err := s.Get("a")
			Expect(err).To(MatchError(ErrNotKeyItem))
		})
		Specify("succeeds", func() {
			s.items["a"] = newKeyItem("a", time.Now().Add(time.Nanosecond))
			v, err := s.Get("a")
			Expect(err).ToNot(HaveOccurred())
			Expect(v).To(Equal("a"))
		})
	})
	Describe("Set", func() {
		Specify("creating new key", func() {
			s.Set("a", "aa", time.Nanosecond)
			Expect(s.items).To(HaveKey("a"))
			Expect(s.items["a"]).To(Equal(newKeyItem("aa", c.Now().Add(time.Nanosecond))))
		})
		Specify("creating new empty key", func() {
			s.Set("", "aa", time.Nanosecond)
			Expect(s.items).To(HaveKey(""))
			Expect(s.items[""]).To(Equal(newKeyItem("aa", c.Now().Add(time.Nanosecond))))
		})
		Specify("creating new key with empty value", func() {
			s.Set("a", "", time.Nanosecond)
			Expect(s.items).To(HaveKey("a"))
			Expect(s.items["a"]).To(Equal(newKeyItem("", c.Now().Add(time.Nanosecond))))
		})
		Specify("rewriting same type", func() {
			s.items["a"] = newKeyItem("aa", c.Now())
			s.Set("a", "bb", time.Nanosecond)
			Expect(s.items["a"]).To(Equal(newKeyItem("bb", c.Now().Add(time.Nanosecond))))
		})
		Specify("rewriting other type", func() {
			s.items["a"] = newListItem(nil, c.Now())
			s.Set("a", "bb", time.Nanosecond)
			Expect(s.items["a"]).To(Equal(newKeyItem("bb", c.Now().Add(time.Nanosecond))))
		})
	})
	Describe("ListGet", func() {
		Specify("expired item error", func() {
			s.items["a"] = newListItem([]string{"a"}, time.Now())
			_, err := s.ListGet("a", 0)
			Expect(err).To(MatchError(ErrKeyNotExists))
		})
		Specify("not existed item error", func() {
			_, err := s.ListGet("a", 0)
			Expect(err).To(MatchError(ErrKeyNotExists))
		})
		Specify("not list item error", func() {
			s.items["a"] = newKeyItem("a", time.Now().Add(time.Nanosecond))
			_, err := s.ListGet("a", 0)
			Expect(err).To(MatchError(ErrNotListItem))
		})
		Specify("index not exists error", func() {
			s.items["a"] = newListItem([]string{"a"}, time.Now().Add(time.Nanosecond))
			_, err := s.ListGet("a", 1)
			Expect(err).To(MatchError(ErrListIndexNotExists))
		})
		Specify("succeeds", func() {
			s.items["a"] = newListItem([]string{"a"}, time.Now().Add(time.Nanosecond))
			v, err := s.ListGet("a", 0)
			Expect(err).ToNot(HaveOccurred())
			Expect(v).To(Equal("a"))
		})
	})
	Describe("ListSet", func() {
		Specify("creating new list", func() {
			s.ListSet("a", []string{"a"}, time.Nanosecond)
			Expect(s.items).To(HaveKey("a"))
			Expect(s.items["a"]).To(Equal(newListItem([]string{"a"}, c.Now().Add(time.Nanosecond))))
		})
		Specify("creating new list with empty key", func() {
			s.ListSet("", []string{"a"}, time.Nanosecond)
			Expect(s.items).To(HaveKey(""))
			Expect(s.items[""]).To(Equal(newListItem([]string{"a"}, c.Now().Add(time.Nanosecond))))
		})
		Specify("creating new nil list", func() {
			s.ListSet("a", nil, time.Nanosecond)
			Expect(s.items).To(HaveKey("a"))
			Expect(s.items["a"]).To(Equal(newListItem(nil, c.Now().Add(time.Nanosecond))))
		})
		Specify("creating new empty list", func() {
			s.ListSet("a", []string{}, time.Nanosecond)
			Expect(s.items).To(HaveKey("a"))
			Expect(s.items["a"]).To(Equal(newListItem([]string{}, c.Now().Add(time.Nanosecond))))
		})
		Specify("rewriting same type", func() {
			s.items["a"] = newListItem(nil, c.Now())
			s.ListSet("a", []string{"a"}, time.Nanosecond)
			Expect(s.items["a"]).To(Equal(newListItem([]string{"a"}, c.Now().Add(time.Nanosecond))))
		})
		Specify("rewriting other type", func() {
			s.items["a"] = newKeyItem("aa", c.Now())
			s.ListSet("a", []string{"a"}, time.Nanosecond)
			Expect(s.items["a"]).To(Equal(newListItem([]string{"a"}, c.Now().Add(time.Nanosecond))))
		})
	})
	Describe("DictGet", func() {
		Specify("expired item error", func() {
			s.items["a"] = newDictItem(map[string]string{"b": "aa"}, time.Now())
			_, err := s.DictGet("a", "b")
			Expect(err).To(MatchError(ErrKeyNotExists))
		})
		Specify("not existed item error", func() {
			_, err := s.DictGet("a", "b")
			Expect(err).To(MatchError(ErrKeyNotExists))
		})
		Specify("not dict item error", func() {
			s.items["a"] = newKeyItem("a", time.Now().Add(time.Nanosecond))
			_, err := s.DictGet("a", "b")
			Expect(err).To(MatchError(ErrNotDictItem))
		})
		Specify("dict key not exists error", func() {
			s.items["a"] = newDictItem(map[string]string{"b": "aa"}, time.Now().Add(time.Nanosecond))
			_, err := s.DictGet("a", "c")
			Expect(err).To(MatchError(ErrDictKeyNotExists))
		})
		Specify("succeeds", func() {
			s.items["a"] = newDictItem(map[string]string{"b": "aa"}, time.Now().Add(time.Nanosecond))
			v, err := s.DictGet("a", "b")
			Expect(err).ToNot(HaveOccurred())
			Expect(v).To(Equal("aa"))
		})
	})
	Describe("DictSet", func() {
		Specify("creating new dict", func() {
			s.DictSet("a", map[string]string{"b": "aa"}, time.Nanosecond)
			Expect(s.items).To(HaveKey("a"))
			Expect(s.items["a"]).To(Equal(newDictItem(map[string]string{"b": "aa"}, c.Now().Add(time.Nanosecond))))
		})
		Specify("creating new dict with empty key", func() {
			s.DictSet("", map[string]string{"b": "aa"}, time.Nanosecond)
			Expect(s.items).To(HaveKey(""))
			Expect(s.items[""]).To(Equal(newDictItem(map[string]string{"b": "aa"}, c.Now().Add(time.Nanosecond))))
		})
		Specify("creating new nil dict", func() {
			s.DictSet("a", nil, time.Nanosecond)
			Expect(s.items).To(HaveKey("a"))
			Expect(s.items["a"]).To(Equal(newDictItem(nil, c.Now().Add(time.Nanosecond))))
		})
		Specify("creating new empty dict", func() {
			s.DictSet("a", map[string]string{}, time.Nanosecond)
			Expect(s.items).To(HaveKey("a"))
			Expect(s.items["a"]).To(Equal(newDictItem(map[string]string{}, c.Now().Add(time.Nanosecond))))
		})
		Specify("rewriting same type", func() {
			s.items["a"] = newDictItem(map[string]string{"b": "aa"}, c.Now())
			s.DictSet("a", map[string]string{"cc": "dd"}, time.Nanosecond)
			Expect(s.items["a"]).To(Equal(newDictItem(map[string]string{"cc": "dd"}, c.Now().Add(time.Nanosecond))))
		})
		Specify("rewriting other type", func() {
			s.items["a"] = newKeyItem("aa", c.Now())
			s.DictSet("a", map[string]string{"cc": "dd"}, time.Nanosecond)
			Expect(s.items["a"]).To(Equal(newDictItem(map[string]string{"cc": "dd"}, c.Now().Add(time.Nanosecond))))
		})
	})
	Describe("Remove", func() {
		Specify("key not exists error", func() {
			Expect(s.Remove("a")).To(MatchError(ErrKeyNotExists))
		})
		Specify("expired key", func() {
			s.items["a"] = baseItem{expiry: c.Now().Add(-time.Nanosecond)}
			Expect(s.Remove("a")).To(Succeed())
			Expect(s.items).ToNot(HaveKey("a"))
		})
		Specify("not expired key", func() {
			s.items["a"] = baseItem{expiry: c.Now().Add(time.Nanosecond)}
			Expect(s.Remove("a")).To(Succeed())
			Expect(s.items).ToNot(HaveKey("a"))
		})
	})
	Describe("Keys", func() {
		Specify("when empty store", func() {
			Expect(s.Keys()).To(BeEmpty())
		})
		Specify("when all keys expired", func() {
			s.items["a"] = baseItem{expiry: c.Now().Add(-time.Nanosecond)}
			s.items["b"] = baseItem{expiry: c.Now()}
			Expect(s.Keys()).To(BeEmpty())
		})
		Specify("when not expired keys exists", func() {
			s.items["a"] = baseItem{expiry: c.Now().Add(-time.Nanosecond)}
			s.items["b"] = baseItem{expiry: c.Now()}
			s.items["c"] = baseItem{expiry: c.Now().Add(time.Nanosecond)}
			s.items["d"] = baseItem{expiry: c.Now().Add(2 * time.Nanosecond)}
			Expect(s.Keys()).To(ConsistOf("c", "d"))
		})
	})
	Specify("StartCleaner and StopCleaner", func() {
		s.items["a"] = baseItem{expiry: c.Now().Add(-time.Nanosecond)}
		s.items["b"] = baseItem{expiry: c.Now()}
		s.items["c"] = baseItem{expiry: c.Now().Add(time.Nanosecond)}

		By("not stated before")
		Expect(s.cleaner).To(BeNil())

		By("starting first time")
		Expect(s.StartCleaner()).To(Succeed())
		Expect(s.cleaner).ToNot(BeNil())
		Expect(s.cleaner.running()).To(BeTrue())

		By("trying to start second time")
		Expect(s.StartCleaner()).To(MatchError(ErrFailedToStartCleaner.detailed("already started")))

		By("waiting cleaning tick")
		time.Sleep(110 * time.Millisecond)

		By("checking expired items removed")
		Expect(s.items).To(HaveLen(1))
		Expect(s.items).To(HaveKey("c"))

		By("stopping cleaner first time")
		Expect(s.StopCleaner()).To(Succeed())
		Expect(s.cleaner).ToNot(BeNil())
		Expect(s.cleaner.running()).To(BeFalse())

		By("checking expired not removed after tick time")
		s.items["a"] = baseItem{expiry: c.Now().Add(-time.Nanosecond)}
		time.Sleep(110 * time.Millisecond)
		Expect(s.items).To(HaveLen(2))
		Expect(s.items).To(HaveKey("a"))
		Expect(s.items).To(HaveKey("c"))

		By("stopping cleaner second time")
		Expect(s.StopCleaner()).To(MatchError(ErrFailedToStopCleaner.detailed("already stopped")))
	})
	Describe("get", func() {
		Specify("expired item error", func() {
			s.items["a"] = baseItem{expiry: time.Now()}
			_, err := s.get("a")
			Expect(err).To(MatchError(ErrKeyNotExists))
		})
		Specify("not existed item error", func() {
			_, err := s.get("a")
			Expect(err).To(MatchError(ErrKeyNotExists))
		})
		Specify("succeeds", func() {
			s.items["a"] = baseItem{expiry: time.Now().Add(time.Nanosecond)}
			i, err := s.get("a")
			Expect(err).ToNot(HaveOccurred())
			Expect(i).To(Equal(s.items["a"]))
		})
	})
	Specify("clean", func() {
		s.items["a"] = baseItem{expiry: time.Now()}
		s.items["b"] = baseItem{expiry: time.Now().Add(-time.Nanosecond)}
		s.items["c"] = baseItem{expiry: time.Now().Add(time.Nanosecond)}
		s.clean()
		Expect(s.items).To(HaveLen(1))
		Expect(s.items).To(HaveKey("c"))
	})
	Specify("expiry", func() {
		Expect(s.expiry(100 * time.Nanosecond)).To(Equal(c.Now().Add(100 * time.Nanosecond)))
		Expect(s.expiry(0)).To(Equal(c.Now()))
		Expect(s.expiry(-100 * time.Nanosecond)).To(Equal(c.Now().Add(-100 * time.Nanosecond)))
	})
})
