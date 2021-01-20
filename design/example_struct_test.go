package design

import (
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Struct mocking test", func() {
	defer GinkgoRecover()
	// mock 초기화.......
	// gomock을 사용하면 mock.go와 같이 자동으로 interface에 대한 mock struct와 method를 생성해줌
	var (
		mockCtrl *gomock.Controller
		m        *MockoneInterface //oneInterface의 mock interface
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		m = NewMockoneInterface(mockCtrl) // interface 별로 mock 생성자가 다름. oneInterface의 mock 생성자 function
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("[Two struct Test]", func() {
		It("should return test", func() {
			//mock struct의 resultOne()은 "test"를 return하도록 정의
			m.EXPECT().resultOne().Return("test").AnyTimes()

			// Newtwo를 통한 two struct 생성
			newTwo := NewTwo(m)
			printResult := newTwo.print()

			// result는 exec에서 받은 resultStdout를 그대로 return한 것!
			Expect(printResult).To(Equal("test"))
		})
	})

})
