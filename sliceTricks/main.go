package main

import (
	"fmt"
	"sort"
)

// #1 【...】 在slice中的使用

func addUsersV1(users []string) {
	for _, user := range users {
		fmt.Println(user)
	}
}

func addUsersV2(users ...string) {
	for _, user := range users {
		fmt.Println(user)
	}
}

var usersT = []string{}

func addUserV1(user string) {
	usersT = append(usersT, user)
}

func addUserV2(user ...string) {
	usersT = append(usersT, user...)
}

// #2  remove from slice by index

func removeFromSliceV1(slice []int, index int) []int {
	slice[index] = slice[len(slice)-1] // 用最后一位元素覆盖要删除的元素
	return slice[:len(slice)-1]        //返回除了最后一位以外的slice
}

func removeFromSliceV2(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...) // idx前后两个slice粘贴
}

type MySlice []int

func (s MySlice) removeFromSliceV3(index int) []int {
	return append(s[:index], s[index+1:]...)
}

// #3 sorting  ,  Len() Swap(i, j int) Less(i, j int) 是实现sort interface的三个函数，需要预先指定
type Numbers []int
type byDec struct {
	Numbers
}

func (n byDec) Len() int           { return len(n.Numbers) }
func (n byDec) Swap(i, j int)      { n.Numbers[i], n.Numbers[j] = n.Numbers[j], n.Numbers[i] }
func (n byDec) Less(i, j int) bool { return n.Numbers[i] > n.Numbers[j] }

type byInc struct {
	Numbers
}

func (n byInc) Len() int           { return len(n.Numbers) }
func (n byInc) Swap(i, j int)      { n.Numbers[i], n.Numbers[j] = n.Numbers[j], n.Numbers[i] }
func (n byInc) Less(i, j int) bool { return n.Numbers[i] < n.Numbers[j] }

func main() {
	addUsersV1([]string{"Alice", "Bob", "Cha"})
	addUsersV2("Alice", "Bob", "Cha")
	addUserV1("alice")
	addUserV1("Bob")
	addUserV1("Cha")
	addUserV2("Alice", "Bob", "Cha")

	numbers := []int{1, 2, 3, 4, 5}
	fmt.Println(removeFromSliceV1(numbers, 2)) //[1 2 5 4]  乱序但是快，不用实际copy元素
	fmt.Println(removeFromSliceV2(numbers, 2)) //[1 2 4 5]  有序但是慢，发生copy
	numbers3 := MySlice{1, 2, 3, 4, 5}
	numbers3 = numbers3.removeFromSliceV3(2) //[1 2 4 5] 更清晰
	fmt.Println(numbers3)

	numberssort := Numbers{3, 4, 7, 15, 8, 3, 7}
	fmt.Println(numberssort)
	sort.Sort(byInc{numberssort})
	fmt.Println(numberssort)
	sort.Sort(byDec{numberssort})
	fmt.Println(numberssort)
}
