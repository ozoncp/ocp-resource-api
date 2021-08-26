package saver_test

import (
	"context"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/ozoncp/ocp-resource-api/internal/mocks"
	"github.com/ozoncp/ocp-resource-api/internal/models"
	s "github.com/ozoncp/ocp-resource-api/internal/saver"
	"time"
)

var _ = Describe("Saver", func() {
	var (
		ctrl        *gomock.Controller
		mockFlusher *mocks.MockFlusher
		saver       s.Saver
		args        []models.Resource
		ctx         context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("default", func() {
		BeforeEach(func() {
			mockFlusher = mocks.NewMockFlusher(ctrl)
			args = []models.Resource{
				models.NewResource(1, 1, 1, 1),
				models.NewResource(2, 2, 2, 2),
				models.NewResource(3, 3, 3, 3),
				models.NewResource(4, 4, 4, 4),
			}
			saver = s.NewSaver(5, mockFlusher, 1*time.Second)
		})

		It("Flush several times", func() {
			saver.Init()
			mockFlusher.EXPECT().Flush(ctx, gomock.Eq(args[0:2]), nil).Times(1)
			mockFlusher.EXPECT().Flush(ctx, gomock.Eq(args[2:]), nil).Times(1)
			mockFlusher.EXPECT().Flush(ctx, gomock.Eq([]models.Resource{}), nil)
			for _, res := range args[0:2] {
				saver.Save(res)
			}
			time.Sleep(1 * time.Second)
			for _, res := range args[2:] {
				saver.Save(res)
			}
			time.Sleep(1 * time.Second)
			saver.Close()
		})

		It("Flush once", func() {
			saver.Init()
			mockFlusher.EXPECT().Flush(ctx, gomock.Eq(args), nil).Times(1)
			mockFlusher.EXPECT().Flush(ctx, gomock.Eq([]models.Resource{}), nil)
			for _, res := range args {
				saver.Save(res)
			}
			time.Sleep(1 * time.Second)
			saver.Close()
		})

		It("Flush failed", func() {
			saver.Init()
			mockFlusher.EXPECT().Flush(ctx, gomock.Eq(args), nil).Return(args).Times(2)
			for _, res := range args {
				saver.Save(res)
			}
			time.Sleep(1*time.Second + 500*time.Millisecond)
			saver.Close()
		})

		It("Close twice", func() {
			mockFlusher.EXPECT().Flush(ctx, gomock.Any(), nil).Times(1).Return([]models.Resource{})
			defer GinkgoRecover()
			saver.Init()
			saver.Close()
			gomega.Expect(saver.Close).Should(gomega.PanicWith("saver must not be closed"))
		})

		It("Save after close", func() {
			mockFlusher.EXPECT().Flush(ctx, gomock.Any(), nil).Times(1).Return([]models.Resource{})
			defer GinkgoRecover()
			saver.Init()
			saver.Close()
			gomega.Expect(func() { saver.Save(args[0]) }).Should(gomega.PanicWith("saver must not be closed"))
		})

		It("Save before init", func() {
			defer GinkgoRecover()
			gomega.Expect(func() { saver.Save(args[0]) }).Should(gomega.PanicWith("saver must be initiated"))
		})

		It("Close before close", func() {
			defer GinkgoRecover()
			gomega.Expect(saver.Close).Should(gomega.PanicWith("saver must be initiated"))
		})
	})
})
