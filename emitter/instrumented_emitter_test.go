package emitter_test

import (
	"github.com/cloudfoundry-incubator/dropsonde/emitter"
	"github.com/cloudfoundry-incubator/dropsonde/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"code.google.com/p/gogoprotobuf/proto"
)

var _ = Describe("InstrumentedUdpEmitter", func() {
		Describe("Emit()", func() {
				var instrumentedEmitter *emitter.InstrumentedEmitter
				var testEvent *events.DropsondeStatus
				var fakeEmitter *emitter.FakeEmitter

				BeforeEach(func() {
					testEvent = &events.DropsondeStatus{SentCount: proto.Uint64(1), ErrorCount: proto.Uint64(0)}
					fakeEmitter = emitter.NewFake()
					instrumentedEmitter, _ = emitter.NewInstrumentedEmitter(fakeEmitter)
				})
				It("calls the concrete emitter", func() {
						Expect(fakeEmitter.Messages).To(HaveLen(0))

						err := instrumentedEmitter.Emit(testEvent)
						Expect(err).ToNot(HaveOccurred())

						Expect(fakeEmitter.Messages).To(HaveLen(1))
					})
				It("increments the ReceivedMetricsCounter", func() {
						Expect(instrumentedEmitter.ReceivedMetricsCounter).To(BeNumerically("==", 0))

						err := instrumentedEmitter.Emit(testEvent)
						Expect(err).ToNot(HaveOccurred())

						Expect(instrumentedEmitter.ReceivedMetricsCounter).To(BeNumerically("==", 1))
					})
				Context("when the concrete Emitter returns no error on Emit()", func() {
						It("increments the SentMetricsCounter", func() {
								Expect(instrumentedEmitter.SentMetricsCounter).To(BeNumerically("==", 0))

								err := instrumentedEmitter.Emit(testEvent)
								Expect(err).ToNot(HaveOccurred())

								Expect(instrumentedEmitter.SentMetricsCounter).To(BeNumerically("==", 1))
							})
					})
				Context("when the concrete Emitter returns an error on Emit()", func() {
						BeforeEach(func(){
							fakeEmitter.ReturnError = true
						})
						It("increments the ErrorCounter", func() {
								Expect(instrumentedEmitter.ErrorCounter).To(BeNumerically("==", 0))
								Expect(instrumentedEmitter.ReceivedMetricsCounter).To(BeNumerically("==", 0))
								Expect(instrumentedEmitter.SentMetricsCounter).To(BeNumerically("==", 0))

								err := instrumentedEmitter.Emit(testEvent)
								Expect(err).To(HaveOccurred())

								Expect(instrumentedEmitter.ErrorCounter).To(BeNumerically("==", 1))
								Expect(instrumentedEmitter.ReceivedMetricsCounter).To(BeNumerically("==", 1))
								Expect(instrumentedEmitter.SentMetricsCounter).To(BeNumerically("==", 0))
							})
					})
			})


	})