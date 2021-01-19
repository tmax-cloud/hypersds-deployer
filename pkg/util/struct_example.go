package util

import "fmt"

// interface를 인자로 받는 function들에 변수를 넘겨줄 때 왜 &를 써서 pointer로 넘겨주었는지는 아래 url 참고
// https://stackoverflow.com/questions/40823315/x-does-not-implement-y-method-has-a-pointer-receiver

/*
특정 struct의 내부에서 다른 struct의 method를 호출 경우 unit test시 mocking을 위한 방법은 두가지가 존재함.
//two: 다른 struct의 method 호출하는 struct, one: 호출되는 struct

1. two의 method가 인자로 one을 받음, 이 떄 one의 method에 대한 inetrface가 구현되어 있어서, 인자로 interface를 받고 inteface의 method 호출

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

func TestStruct() {
	s := one{name: "ttt"}
	result := two{name: "xxx"}
	result.print(&s)
}

2. two가 one을 struct 내부에 가지는데 interface로 가짐. 그리고 two의 method 내에서 inetrface의 method 호출
당연히 one의 method에 대한 interface 구현되어 있다는 가정임
*/
// 2번째 방법 예시
// one struct method의 Interface인 oneInterface는 interface.go 참고
type one struct {
	name string
}
type two struct {
	a oneInterface
}

func (o *one) resultOne() string {
	return o.name + "out"
}

// two 생성시 인자로 one의 interface를 받고, 이를 가지게 하기 위해 Newtwo란 생성자 function 정의
// 실제 코드 상에서는 one struct가 인자로 넘겨지고,
// unit test시에는 Newtwo에 넘겨주는 변수는 one의 interface 정의를 따라 생성한 mock struct를 넘겨서 mocking함
func Newtwo(a oneInterface) two {
	return two{a: a}
}

func (t *two) print() string {
	fmt.Printf("%s\n", t.a.resultOne())
	return t.a.resultOne()
}

func TestStruct() {
	s := one{name: "ttt"}
	result := Newtwo(&s)
	_ = result.print()
}
