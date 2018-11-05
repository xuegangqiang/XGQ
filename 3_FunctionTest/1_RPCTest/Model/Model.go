package Model

import (
	"context"
	"fmt"
	"time"
)

type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

type ModelData struct {
}

func (md *ModelData) RpcxMethod(cxt context.Context, args *Args, reply *Reply) error {

	reply.C = args.A * args.B
	t := time.Now()

	str := fmt.Sprintf("%d年%d月%d日 %d:%d:%d.%d : %d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1e6, reply.C)
	fmt.Println(str)

	// fmt.Println("%d年%d月%d日 %d:%d:%d.%d : %d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1e6, reply.C)
	// log.Printf("%d年%i月%i日 %i:%i:%i.%i : %i", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1e6, reply.C)

	return nil
}
