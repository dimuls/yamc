package store

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("newTicker", func() {
	Specify("errors", func() {
		_, err := newTicker(time.Nanosecond, nil)
		Expect(err).To(HaveOccurred())
	})
	Specify("succeeds", func() {
		t, err := newTicker(time.Nanosecond, func() {})
		Expect(err).ToNot(HaveOccurred())
		Expect(t).ToNot(BeNil())
		Expect(t.period).To(Equal(time.Nanosecond))
		Expect(t.f).ToNot(BeNil())
		Expect(t.stopper).To(BeNil())
		Expect(t.isRunning()).To(BeFalse())
	})
})

var _ = Specify("Ticker", func() {
	By("creating ticker")
	ticks := 0
	t, err := newTicker(100*time.Millisecond, func() {
		By("tick occurred")
		ticks++
	})
	Expect(err).ToNot(HaveOccurred())
	Expect(t.isRunning()).To(BeFalse())

	By("starting ticker")
	Expect(t.start()).To(Succeed())
	Expect(t.isRunning()).To(BeTrue())

	time.Sleep(110 * time.Millisecond)

	By("after tick")
	Expect(ticks).To(Equal(1))

	By("trying to start one more time")
	Expect(t.start()).ToNot(Succeed())
	Expect(t.isRunning()).To(BeTrue())

	time.Sleep(110 * time.Millisecond)

	By("after tick")
	Expect(ticks).To(Equal(2))

	By("stopping timer")
	Expect(t.stop()).To(Succeed())
	Expect(t.isRunning()).To(BeFalse())

	time.Sleep(110 * time.Millisecond)

	By("after tick")
	Expect(ticks).To(Equal(2))

	By("stopping timer second time")
	Expect(t.stop()).ToNot(Succeed())
	Expect(t.isRunning()).To(BeFalse())
})
