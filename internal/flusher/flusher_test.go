package flusher_test

import (
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	f "github.com/ozoncp/ocp-resource-api/internal/flusher"
	"github.com/ozoncp/ocp-resource-api/internal/mocks"
	"github.com/ozoncp/ocp-resource-api/internal/models"
	"golang.org/x/net/context"
)

var _ = Describe("Flusher", func() {
	var (
		ctrl     *gomock.Controller
		mockRepo *mocks.MockRepo
		flusher  f.Flusher
		args     []models.Resource
		ctx      context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("chunk of two elements", func() {
		BeforeEach(func() {
			mockRepo = mocks.NewMockRepo(ctrl)
			args = []models.Resource{
				models.NewResource(1, 1, 1, 1),
				models.NewResource(2, 2, 2, 2),
				models.NewResource(3, 3, 3, 3),
				models.NewResource(4, 4, 4, 4),
			}
			flusher = f.NewFlusher(2, mockRepo)
		})

		It("Valid two chunks", func() {
			mockRepo.EXPECT().AddEntities(ctx, gomock.Eq(args[0:2])).Times(1)
			mockRepo.EXPECT().AddEntities(ctx, gomock.Eq(args[2:])).Times(1)
			actual := flusher.Flush(ctx, args)
			gomega.Expect(actual).Should(gomega.BeEmpty())
		})

		Describe("Failures", func() {
			It("Fail on first chunk", func() {
				ErrNetwork := errors.New("network error")
				mockRepo.EXPECT().AddEntities(ctx, gomock.Eq(args[0:2])).Return(ErrNetwork)
				actual := flusher.Flush(ctx, args)
				gomega.Expect(actual).Should(gomega.BeEquivalentTo(args))
			})

			It("Fail on second chunk", func() {
				ErrNetwork := errors.New("network error")
				mockRepo.EXPECT().AddEntities(ctx, gomock.Eq(args[0:2])).Times(1)
				mockRepo.EXPECT().AddEntities(ctx, gomock.Eq(args[2:])).Return(ErrNetwork)
				actual := flusher.Flush(ctx, args)
				gomega.Expect(actual).Should(gomega.BeEquivalentTo(args[2:]))
			})
		})
	})

	Context("chunk of 5 elements", func() {
		BeforeEach(func() {
			mockRepo = mocks.NewMockRepo(ctrl)
			args = []models.Resource{
				models.NewResource(1, 1, 1, 1),
				models.NewResource(2, 2, 2, 2),
				models.NewResource(3, 3, 3, 3),
				models.NewResource(4, 4, 4, 4),
			}
			flusher = f.NewFlusher(5, mockRepo)
		})

		It("Valid single chunks", func() {
			mockRepo.EXPECT().AddEntities(ctx, gomock.Eq(args)).Times(1)
			actual := flusher.Flush(ctx, args)
			gomega.Expect(actual).Should(gomega.BeEmpty())
		})

		Describe("Failures", func() {
			It("Fail on first chunk", func() {
				ErrNetwork := errors.New("network error")
				mockRepo.EXPECT().AddEntities(ctx, gomock.Eq(args)).Return(ErrNetwork)
				actual := flusher.Flush(ctx, args)
				gomega.Expect(actual).Should(gomega.BeEquivalentTo(args))
			})
		})
	})
})
