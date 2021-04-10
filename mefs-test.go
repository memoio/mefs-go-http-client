//该文件用于测试mefs各种功能
package shell

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	TESTBUCKET     = "b0"
	RandomDataSize = 1024 * 1024 * 100
)

//测试时空值的计算，流程：没过一段时间，user发送一段数据，keeper进行一次计算，将实际值与理论值对比
//目前使用metb进行测试，传入进行测试的user和keeper端口
func ReslultSumaryTest(userPort, keeperPort string) {
	theory := int64(0)                //时空值的理论值
	testLastSize := int64(0)          //上一次put文件的大小
	testLastTime := time.Now().Unix() //上一次测试的时间
	userSh := NewShell("localhost:" + userPort)
	keeperSh := NewShell("localhost:" + keeperPort)

	userSh.CreateBucket(TESTBUCKET) //先建bucket

	for i := 1; ; i++ { //每过30分钟进行一次测试
		fmt.Println("======================")
		testThisTime := time.Now().Unix() //本次测试时间
		fmt.Println("本次测试时间：", time.Unix(testThisTime, 0).In(time.Local))
		testThisSize := rand.Int63n(RandomDataSize)
		fmt.Println("传入文件大小:", testThisSize)
		userSh.putRandomObject(testThisSize)
		theory += (testThisTime - testLastTime) * testLastSize
		actual := keeperSh.ResultSummary() //时空支付的实际值
		fmt.Println("实际值：", actual)
		fmt.Println("理论值", theory/3)
		testLastTime = testThisTime
		testLastSize += testThisSize
		fmt.Println("======================\n")
		time.Sleep(30 * time.Minute)
	}

}

//获取节点信息的操作
func (s *Shell) TestLocalinfo() {
	var sl StringList
	rb := s.Request("test/localinfo")
	rb.Exec(context.Background(), &sl)
	fmt.Println(sl.String())
}

//keeper计算时空值命令，用于测试，返回计算好的时空值
func (s *Shell) ResultSummary() int {
	var il IntList
	rb := s.Request("test/resultsummary")
	rb.Exec(context.Background(), &il)
	if len(il.ChildLists) < 1 {
		fmt.Println("计算失败")
		return 0
	}
	return il.ChildLists[0]
}

//put一个指定大小的测试文件
func (s *Shell) putRandomObject(size int64) {
	data := make([]byte, size)
	fillRandom(data)
	buf := bytes.NewBuffer(data)
	objectName := "test_" + strconv.Itoa(int(size))
	_, err := s.PutObject(buf, objectName, TESTBUCKET)
	if err != nil {
		fmt.Println("PutObject err!", err)
	}
}
func fillRandom(p []byte) {
	for i := 0; i < len(p); i += 7 {
		val := rand.Int63()
		for j := 0; i+j < len(p) && j < 7; j++ {
			p[i+j] = byte(val)
			val >>= 8
		}
	}
}
