package design

import "fmt"

// 특정 struct의 내부에서 다른 struct의 method를 호출하는 경우의 예시.
// two: 다른 struct의 method를 호출하는 struct
// one: two에게 호출되는 struct

// two의 method (print)가 인자로 one을 받아야 하는 상황이라면, one의 method에 대한 inetrface oneInterface를 인자로 받고 inteface의 method 호출
// 그래야 unit test 작성 시 mocking 가능

type oneInterface interface {
	resultOne() string
}

type one struct {
	name string
}
func (o *one) resultOne() string {
	return o.name + "out"
}

type two struct {
	a oneInterface
}
func (t *two) print(a oneInterface) {
	fmt.Printf("%s\n", a.resultOne())
}
func NewTwo(a oneInterface) two {
    return two{a: a}
}

// 왜 &를 써서 pointer로 넘겨주었는지는 아래 url 참고
// https://stackoverflow.com/questions/40823315/x-does-not-implement-y-method-has-a-pointer-receiver
func DoSomethingByTwo() {
	s := one{name: "ttt"}
	twoInstance := two{name: "xxx"}

	twoInstance.print(&s)
}
