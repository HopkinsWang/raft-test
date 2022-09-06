package main

import (
	pb "RingAllReduce_29server/ondisk/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strings"
	"time"
)


var (
	hostaddrlist = []string{
		"192.168.0.232:61001",
		"192.168.0.232:61002",
		"192.168.0.232:61003",}
)

//func ParseHttp(clientUrl string) {
//
//	u, _ := url.Parse(clientUrl)                //将string解析成*URL格式
//	fmt.Printf("u =%v\n", u)                    // go?a=123&b=456
//	fmt.Printf("u.RawQuery = %v\n", u.RawQuery) //编码后的查询字符串，没有'?' a=123&b=456
//	values, _ := url.ParseQuery(u.RawQuery)     //返回Values类型的字典
//	fmt.Println(values)                         // map[a:[123] b:[456]]
//	fmt.Println(values["a"])                    //[123]
//	fmt.Println(values.Get("b"))                //456
//}

func main() {
	//ParseHttp("http://localhost:8080/go?a=123&b=456")
	main_test()

}

func main_test(){
	Ondisktest()
}



func Ondisktest() {
	//var clients map[string]pb.RaftdServiceClient
	clients := make(map[string]pb.RaftdServiceClient,3)
	for _,addr:=range hostaddrlist{
		clients[addr]= startTestclient(addr)
	}
	for i :=1;i<4;i++{
		key:= fmt.Sprintf("%v",i)
		val:= strings.Repeat(key,4)
		log.Printf("key = %v, val = %v",key,val)
		_,err:=clients[hostaddrlist[i-1]].Put(context.Background(), &pb.PutRequest{Key:key,Val: val})
		if err!=nil{
			log.Printf("Put failed, Err: %v",err)
		}
	}

	time.Sleep(3*time.Second)
	for i :=3;i<0;i--{
		res, err:=clients[hostaddrlist[i]].Get(context.Background(), &pb.GetRequest{Key:fmt.Sprintf("%v",i)})
		if err!=nil{
			log.Printf("clients[hostaddrlist[%v]].Get faile, %v",i,err)
		}
		log.Printf("Host: %v  Get! key=%v , val=%v ",hostaddrlist[i],res.Key,res.Val)
	}
	log.Printf("Test Finish!")

}

func startTestclient(addr string) pb.RaftdServiceClient{
	rcv_size := 512 * 1024 * 1024
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(rcv_size), grpc.MaxCallSendMsgSize(rcv_size)))
	if err != nil {
		log.Printf(" Net.connect err: %v", err)
	}
	return pb.NewRaftdServiceClient(conn)
}







//开始分层分组的ringAllResuce
