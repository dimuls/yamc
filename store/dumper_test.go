package store

import (
	"io/ioutil"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("FileDumper", func() {
	var (
		file *os.File
		d    FileDumper
	)
	BeforeEach(func() {
		var err error
		file, err = ioutil.TempFile("", "file-dump-")
		if err != nil {
			Fail("temp file creating error: " + err.Error())
		}
		d = FileDumper(file.Name())
	})
	Describe("dump and load", func() {
		Specify("success empty items if file not exists", func() {
			d = FileDumper("!--- NOT EXISTS ---!")
			items, err := d.load()
			Expect(items).To(BeEmpty())
			Expect(err).ToNot(HaveOccurred())
		})
		Specify("fail to decode invalid data error", func() {
			_, err := file.Write([]byte("!--- INVALID DATA ---!"))
			Expect(err).ToNot(HaveOccurred())
			Expect(file.Sync()).To(Succeed())
			_, err = d.load()
			Expect(err).To(MatchError(ErrFailToDecodeDumpFile.detailed("unexpected EOF")))
		})
		Specify("success load", func() {
			exp := items{
				"key":  newKeyItem("value", time.Now().Add(10*time.Second)),
				"lkey": newListItem([]string{"l1", "l2", "l3"}, time.Now()),
				"dkey": newDictItem(map[string]string{"dk1": "dv1", "dk2": "dv2"}, time.Now().Add(-10*time.Second)),
			}
			err := d.dump(exp)
			Expect(err).ToNot(HaveOccurred())
			got, err := d.load()
			Expect(err).ToNot(HaveOccurred())
			Expect(got).To(HaveLen(3))
			Expect(got).To(HaveKey("key"))
			Expect(got).To(HaveKey("lkey"))
			Expect(got).To(HaveKey("dkey"))
			Expect(got["key"]).To(beKeyItem(exp["key"].(keyItem)))
			Expect(got["lkey"]).To(beListItem(exp["lkey"].(listItem)))
			Expect(got["dkey"]).To(beDictItem(exp["dkey"].(dictItem)))
		})
	})
})

func matchExpiry(e time.Time) types.GomegaMatcher {
	return WithTransform(func(i interface{}) bool {
		switch ii := i.(type) {
		case keyItem:
			return ii.expiry.Equal(e)
		case listItem:
			return ii.expiry.Equal(e)
		case dictItem:
			return ii.expiry.Equal(e)
		default:
			Fail("not item", 2)
		}
		return false
	}, BeTrue())
}

func beKeyItem(ki keyItem) types.GomegaMatcher {
	matchValue := WithTransform(func(i interface{}) string {
		ii, ok := i.(keyItem)
		if !ok {
			Fail("not a key item")
		}
		return ii.value
	}, Equal(ki.value))
	return And(matchExpiry(ki.expiry), matchValue)
}

func beListItem(li listItem) types.GomegaMatcher {
	matchValue := WithTransform(func(i interface{}) []string {
		ii, ok := i.(listItem)
		if !ok {
			Fail("not a list item")
		}
		return ii.list
	}, Equal(li.list))
	return And(matchExpiry(li.expiry), matchValue)
}

func beDictItem(di dictItem) types.GomegaMatcher {
	matchValue := WithTransform(func(i interface{}) map[string]string {
		ii, ok := i.(dictItem)
		if !ok {
			Fail("not a dict item")
		}
		return ii.dict
	}, Equal(di.dict))
	return And(matchExpiry(di.expiry), matchValue)
}
