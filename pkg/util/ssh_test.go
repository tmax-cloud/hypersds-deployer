package util

import (
	"bytes"

    "hypersds-provisioner/pkg/common/wrapper"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os/user"
)

var _ = Describe("SSH Command Test", func() {
	defer GinkgoRecover()
	// mock 초기화.......
	// gomock을 사용하면 mock.go와 같이 자동으로 interface에 대한 mock struct와 method를 생성해줌
	var (
		mockCtrl *gomock.Controller
		m        *wrapper.MockExecInterface //ExecInterface의 Mock interface, interface 별로 다 생성됨
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		m = wrapper.NewMockExecInterface(mockCtrl) //ExecInterface의 Mock 생성자 function
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("[RunSSHCmd Test]", func() {
		It("should return same username", func() {
			currentUser, err := user.Current()
			testCommand := []string{"printf", "$USER"}
			// mock struct는 정의된 function 호출시 들어올 인자값과 return 값을 정의할 수 있음.
			//  이는 정해진 값일 수도 있고, 아래와 같이 함수를 정의해서 계산값을 return받거나 인자로 들어온 변수를 수정할 수도 있음
			// commandExecute function은 interface 정의에 따라 5개의 변수를 인자로 받음, 여기서는 어떤 값이 오든 상관 없게 하려고, gomock.Any()를 넣어줌
			// return 시 인자로 들어온 변수를 수정하기 위해 DoAndReturn을 써서 function을 정의함. 정의된 function도 마찬가지로 5개의 변수를 인자로 받고,
			// function 내에서 사용할 변수만 type 제대로 명시하고,나머지는 사용하지 않으므로 interface{}로 명시하였음
			// 자세한 gomock 사용법은 인터넷 참고!
			m.EXPECT().CommandExecute(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
				func(resultStdout, resultStderr *bytes.Buffer, ctx, name interface{}, arg ...string) error {
					// resultStdout에 값 push
					resultStdout.WriteString(arg[len(arg)-1])
					return nil
				}).AnyTimes()

			// 정해진 값 return하고 싶으면 이와 같이 간단하게 return 값 정의 가능
			//m.EXPECT().commandExecute(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			result, err := RunSSHCmd(m, currentUser.Username, "localhost", testCommand...)
			Expect(err).NotTo(HaveOccurred())
			// result는 exec에서 받은 resultStdout를 그대로 return한 것!
			Expect(result.String()).To(Equal(testCommand[len(testCommand)-1]))
		})
	})

})
