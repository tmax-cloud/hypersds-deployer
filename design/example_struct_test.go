package design

import (
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Area mocking test", func() {
	defer GinkgoRecover()
	// mock 초기화
	// gomock을 사용하면 mock.go와 같이 자동으로 interface에 대한 mock struct와 method를 생성해줌
	var (
		mockCtrl *gomock.Controller
		m        *MockShape // Shape의 mock struct
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		m = NewMockShape(mockCtrl) // interface 별로 mock 생성자가 다름. Shape의 mock 생성자 function
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("[Area Test]", func() {
		It("should return true", func() {
			//mock struct의 Area() 함수를 25.0을 return하도록 정의
			m.EXPECT().Area().Return(25.).AnyTimes()

			// SquareManager가 자신 square의 넓이를 36.0으로 갖도록 정의
			sm := SquareManager{6.}

			// 결과는 true여야 함
			cmpResult := sm.IsWiderThan(m)
			Expect(cmpResult).To(Equal(true))
		})
	})

})
