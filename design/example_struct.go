package design

import (
	"fmt"
	"math"
)

// 특정 struct의 내부에서 다른 struct의 method를 호출하는 경우의 예시

// Square: Circle의 넓이와 자신의 Sqaure의 넓이를 비교하고자 하는 struct
// Circle: Shape를 구현한 struct

// Square의 method(IsWiderThan)가 인자로 Circle을 받아야 하는 상황이라면, Circle의 interface Shape를 인자로 받고 inteface의 method 호출
// 그래야 unit test 작성 시 mocking 가능

type Circle struct {
	radius float64
}

func (this *Circle) Area() float64 {
	return math.Pi * this.radius * this.radius
}

type SquareManager struct {
	side float64
}

func (this SquareManager) IsWiderThan(s Shape) bool {
	return this.side*this.side > s.Area()
}

// 왜 &를 써서 pointer로 넘겨주었는지는 아래 url 참고
// https://stackoverflow.com/questions/40823315/x-does-not-implement-y-method-has-a-pointer-receiver
func CompareSquareAndCircle() {
	c := Circle{3.}
	sm := SquareManager{7.}

	if sm.IsWiderThan(&c) {
		fmt.Println("square is wider than circle")
	} else {
		fmt.Println("circle is wider than square")
	}
}
