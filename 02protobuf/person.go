package main

import (
	"fmt"
	"log"

	// 这里是创建一个目录，把生成个person.pb.go放里面
	person "./protobuf"
	proto "github.com/golang/protobuf/proto"
)

func main() {

	// 这里不加上Name: Age: 就要报错
	p := &person.Person{Name: "zhangsan", Age: 99}

	data, err := proto.Marshal(p)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data, string(data))

	newp := new(person.Person)

	err = proto.Unmarshal(data, newp)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(newp.GetName(), newp.GetAge())

}
