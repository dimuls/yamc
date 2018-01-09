package store

import (
	"reflect"
	"runtime"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("NewStore", func() {
	Specify("invalid params error", func() {
		_, err := NewStore(Params{
			CleaningPeriod: 100*time.Millisecond - time.Nanosecond,
			DumpingPeriod:  60 * time.Second,
		}, testClock{}, &testDumper{})
		Expect(err).To(MatchError(ErrInvalidParams.detailed("too small cleaning period, must be >= 100ms")))
		_, err = NewStore(Params{
			CleaningPeriod: 100 * time.Millisecond,
			DumpingPeriod:  60*time.Second - time.Nanosecond,
		}, testClock{}, &testDumper{})
		Expect(err).To(MatchError(ErrInvalidParams.detailed("too small dumping period, must be >= 60s")))
	})
	Specify("succeeds", func() {
		p := Params{
			CleaningPeriod: 100 * time.Millisecond,
			DumpingPeriod:  60 * time.Second,
		}
		c := testClock{}
		d := &testDumper{}
		si, err := NewStore(p, c, d)
		Expect(err).ToNot(HaveOccurred())
		d.expectLoad()
		s := si.(*store)
		Expect(s).ToNot(BeNil())
		Expect(s.params).To(Equal(p))
		Expect(s.clock).To(Equal(c))
		Expect(s.items).ToNot(BeNil())
		Expect(s.items).To(BeEmpty())
		Expect(s.cleaning).To(BeNil())
	})
})

var _ = Describe("store", func() {
	var (
		c testClock
		d *testDumper
		s *store
	)
	BeforeEach(func() {
		c = testClock(time.Now())
		d = &testDumper{}
		si, err := NewStore(Params{
			CleaningPeriod: 100 * time.Millisecond,
			DumpingPeriod:  60 * time.Second,
		}, c, d)
		Expect(err).ToNot(HaveOccurred())
		d.expectLoad()
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
			Expect(s.items["a"]).To(Equal(newKeyItem("aa", c.now().Add(time.Nanosecond))))
		})
		Specify("creating new empty key", func() {
			s.Set("", "aa", time.Nanosecond)
			Expect(s.items).To(HaveKey(""))
			Expect(s.items[""]).To(Equal(newKeyItem("aa", c.now().Add(time.Nanosecond))))
		})
		Specify("creating new key with empty value", func() {
			s.Set("a", "", time.Nanosecond)
			Expect(s.items).To(HaveKey("a"))
			Expect(s.items["a"]).To(Equal(newKeyItem("", c.now().Add(time.Nanosecond))))
		})
		Specify("rewriting same type", func() {
			s.items["a"] = newKeyItem("aa", c.now())
			s.Set("a", "bb", time.Nanosecond)
			Expect(s.items["a"]).To(Equal(newKeyItem("bb", c.now().Add(time.Nanosecond))))
		})
		Specify("rewriting other type", func() {
			s.items["a"] = newListItem(nil, c.now())
			s.Set("a", "bb", time.Nanosecond)
			Expect(s.items["a"]).To(Equal(newKeyItem("bb", c.now().Add(time.Nanosecond))))
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
			Expect(s.items["a"]).To(Equal(newListItem([]string{"a"}, c.now().Add(time.Nanosecond))))
		})
		Specify("creating new list with empty key", func() {
			s.ListSet("", []string{"a"}, time.Nanosecond)
			Expect(s.items).To(HaveKey(""))
			Expect(s.items[""]).To(Equal(newListItem([]string{"a"}, c.now().Add(time.Nanosecond))))
		})
		Specify("creating new nil list", func() {
			s.ListSet("a", nil, time.Nanosecond)
			Expect(s.items).To(HaveKey("a"))
			Expect(s.items["a"]).To(Equal(newListItem(nil, c.now().Add(time.Nanosecond))))
		})
		Specify("creating new empty list", func() {
			s.ListSet("a", []string{}, time.Nanosecond)
			Expect(s.items).To(HaveKey("a"))
			Expect(s.items["a"]).To(Equal(newListItem([]string{}, c.now().Add(time.Nanosecond))))
		})
		Specify("rewriting same type", func() {
			s.items["a"] = newListItem(nil, c.now())
			s.ListSet("a", []string{"a"}, time.Nanosecond)
			Expect(s.items["a"]).To(Equal(newListItem([]string{"a"}, c.now().Add(time.Nanosecond))))
		})
		Specify("rewriting other type", func() {
			s.items["a"] = newKeyItem("aa", c.now())
			s.ListSet("a", []string{"a"}, time.Nanosecond)
			Expect(s.items["a"]).To(Equal(newListItem([]string{"a"}, c.now().Add(time.Nanosecond))))
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
			Expect(s.items["a"]).To(Equal(newDictItem(map[string]string{"b": "aa"}, c.now().Add(time.Nanosecond))))
		})
		Specify("creating new dict with empty key", func() {
			s.DictSet("", map[string]string{"b": "aa"}, time.Nanosecond)
			Expect(s.items).To(HaveKey(""))
			Expect(s.items[""]).To(Equal(newDictItem(map[string]string{"b": "aa"}, c.now().Add(time.Nanosecond))))
		})
		Specify("creating new nil dict", func() {
			s.DictSet("a", nil, time.Nanosecond)
			Expect(s.items).To(HaveKey("a"))
			Expect(s.items["a"]).To(Equal(newDictItem(nil, c.now().Add(time.Nanosecond))))
		})
		Specify("creating new empty dict", func() {
			s.DictSet("a", map[string]string{}, time.Nanosecond)
			Expect(s.items).To(HaveKey("a"))
			Expect(s.items["a"]).To(Equal(newDictItem(map[string]string{}, c.now().Add(time.Nanosecond))))
		})
		Specify("rewriting same type", func() {
			s.items["a"] = newDictItem(map[string]string{"b": "aa"}, c.now())
			s.DictSet("a", map[string]string{"cc": "dd"}, time.Nanosecond)
			Expect(s.items["a"]).To(Equal(newDictItem(map[string]string{"cc": "dd"}, c.now().Add(time.Nanosecond))))
		})
		Specify("rewriting other type", func() {
			s.items["a"] = newKeyItem("aa", c.now())
			s.DictSet("a", map[string]string{"cc": "dd"}, time.Nanosecond)
			Expect(s.items["a"]).To(Equal(newDictItem(map[string]string{"cc": "dd"}, c.now().Add(time.Nanosecond))))
		})
	})
	Describe("Remove", func() {
		Specify("key not exists error", func() {
			Expect(s.Remove("a")).To(MatchError(ErrKeyNotExists))
		})
		Specify("expired key", func() {
			s.items["a"] = baseItem{expiry: c.now().Add(-time.Nanosecond)}
			Expect(s.Remove("a")).To(Succeed())
			Expect(s.items).ToNot(HaveKey("a"))
		})
		Specify("not expired key", func() {
			s.items["a"] = baseItem{expiry: c.now().Add(time.Nanosecond)}
			Expect(s.Remove("a")).To(Succeed())
			Expect(s.items).ToNot(HaveKey("a"))
		})
	})
	Describe("Keys", func() {
		Specify("when empty store", func() {
			Expect(s.Keys()).To(BeEmpty())
		})
		Specify("when all keys expired", func() {
			s.items["a"] = baseItem{expiry: c.now().Add(-time.Nanosecond)}
			s.items["b"] = baseItem{expiry: c.now()}
			Expect(s.Keys()).To(BeEmpty())
		})
		Specify("when not expired keys exists", func() {
			s.items["a"] = baseItem{expiry: c.now().Add(-time.Nanosecond)}
			s.items["b"] = baseItem{expiry: c.now()}
			s.items["c"] = baseItem{expiry: c.now().Add(time.Nanosecond)}
			s.items["d"] = baseItem{expiry: c.now().Add(2 * time.Nanosecond)}
			Expect(s.Keys()).To(ConsistOf("c", "d"))
		})
	})
	Specify("StartCleaning and StopCleaning", func() {
		defer s.StopCleaning()
		s.items["a"] = baseItem{expiry: c.now().Add(-time.Nanosecond)}
		s.items["b"] = baseItem{expiry: c.now()}
		s.items["c"] = baseItem{expiry: c.now().Add(time.Nanosecond)}

		By("not stated before")
		Expect(s.cleaning).To(BeNil())

		By("starting first time")
		Expect(s.StartCleaning()).To(Succeed())
		Expect(s.cleaning).ToNot(BeNil())
		Expect(s.cleaning.isRunning()).To(BeTrue())

		By("trying to start second time")
		Expect(s.StartCleaning()).To(MatchError(ErrFailToStartCleaning.detailed("already started")))

		By("waiting cleaning tick")
		time.Sleep(110 * time.Millisecond)

		By("checking expired items removed")
		Expect(s.items).To(HaveLen(1))
		Expect(s.items).To(HaveKey("c"))

		By("stopping cleaning first time")
		Expect(s.StopCleaning()).To(Succeed())
		Expect(s.cleaning).ToNot(BeNil())
		Expect(s.cleaning.isRunning()).To(BeFalse())

		By("checking expired not removed after tick time")
		s.items["a"] = baseItem{expiry: c.now().Add(-time.Nanosecond)}
		time.Sleep(110 * time.Millisecond)
		Expect(s.items).To(HaveLen(2))
		Expect(s.items).To(HaveKey("a"))
		Expect(s.items).To(HaveKey("c"))

		By("stopping cleaning second time")
		Expect(s.StopCleaning()).To(MatchError(ErrFailToStopCleaning.detailed("already stopped")))
	})
	Specify("StartDumping and StopDumping", func() {
		defer s.StopDumping()
		s.params.DumpingPeriod = 100 * time.Millisecond

		s.Set("k", "v", time.Second)
		s.ListSet("lk", []string{"a", "b"}, 2*time.Second)
		s.DictSet("dk", map[string]string{"dk": "dv"}, 3*time.Second)

		By("not stated before")
		Expect(s.dumping).To(BeNil())
		d.expectNoCalls()

		By("starting first time")
		Expect(s.StartDumping()).To(Succeed())
		Expect(s.dumping).ToNot(BeNil())
		Expect(s.dumping.isRunning()).To(BeTrue())
		d.expectNoCalls()

		By("trying to start second time")
		Expect(s.StartDumping()).To(MatchError(ErrFailToStartDumping.detailed("already started")))

		By("waiting dumping tick")
		time.Sleep(110 * time.Millisecond)

		By("checking not expired items are dumped")
		d.expectDump(s.items)
		d.expectNoCalls()

		By("stopping dumping first time")
		Expect(s.StopDumping()).To(Succeed())
		Expect(s.dumping).ToNot(BeNil())
		Expect(s.dumping.isRunning()).To(BeFalse())

		By("checking that after tick time no dump occured")
		time.Sleep(110 * time.Millisecond)
		d.expectNoCalls()

		By("stopping dumping second time")
		Expect(s.StopDumping()).To(MatchError(ErrFailToStopDumping.detailed("already stopped")))
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
	Specify("dump", func() {
		s.Set("k", "v", time.Second)
		s.ListSet("lk", []string{"a", "b"}, 2*time.Second)
		s.DictSet("dk", map[string]string{"dk": "dv"}, 3*time.Second)
		s.dump()
		d.expectDump(s.items)
		d.expectNoCalls()
	})
	Specify("expiry", func() {
		Expect(s.expiry(100 * time.Nanosecond)).To(Equal(c.now().Add(100 * time.Nanosecond)))
		Expect(s.expiry(0)).To(Equal(c.now()))
		Expect(s.expiry(-100 * time.Nanosecond)).To(Equal(c.now().Add(-100 * time.Nanosecond)))
	})
})

type testClock time.Time

func (tc testClock) now() time.Time {
	return time.Time(tc)
}

type testDumper struct {
	calls []call
	items items
	error error
}

func (td *testDumper) newCall(f interface{}, args ...interface{}) {
	td.calls = append(td.calls, call{
		method: funcToName(f),
		args:   args,
	})
}

func (td *testDumper) popCall() call {
	if len(td.calls) == 0 {
		panic("no calls are left")
	}
	var c call
	c, td.calls = td.calls[0], td.calls[1:]
	return c
}

func (td *testDumper) expectNoCalls() {
	ExpectWithOffset(1, td.calls).To(BeEmpty())
}

func (td *testDumper) dump(items items) error {
	td.newCall(td.dump, items)
	return td.error
}

func (td *testDumper) expectDump(items items) {
	ExpectWithOffset(1, td.calls).ToNot(BeEmpty())
	ExpectWithOffset(1, td.popCall()).To(beCall(td.dump, items))
}

func (td *testDumper) load() (items, error) {
	td.newCall(td.load)
	if td.items == nil {
		td.items = items{}
	}
	return td.items, td.error
}

func (td *testDumper) expectLoad() {
	ExpectWithOffset(1, td.calls).ToNot(BeEmpty())
	ExpectWithOffset(1, td.popCall()).To(beCall(td.load))
}

type call struct {
	method string
	args   []interface{}
}

func beCall(method interface{}, args ...interface{}) types.GomegaMatcher {
	methodMatcher := WithTransform(func(c interface{}) string {
		return c.(call).method
	}, BeIdenticalTo(funcToName(method)))
	argsMatcher := BeEmpty()
	if len(args) > 0 {
		argsMatcher = Equal(args)
	}
	return And(methodMatcher, WithTransform(func(c interface{}) []interface{} {
		return c.(call).args
	}, argsMatcher))
}

func funcToName(f interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	name = strings.TrimPrefix(name, "github.com/someanon/yamc/")
	name = strings.TrimSuffix(name, "-fm")
	return name
}
