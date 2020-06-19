package idworker

import (
	"errors"
	"fmt"
	"genids/config"
	"genids/utils/logs"
	"github.com/getsentry/sentry-go"
	"sync"
	"time"
)

const (
	twepoch      = int64(1525705533) // 默认起始的时间戳 1449473700000 。计算时，减去这个值
	NodeIdBits   = uint(2)           //节点 所占位置00 01 10 11，对应0， 1， 2， 3
	sequenceBits = uint(11)          //自增ID 所占用位置

	/*
	 * 1 符号位  |  38 时间戳                                   | 2 节点   | 11 （秒内）自增ID
	 * 0        |  000000 00000000 00000000 00000000 00000000 |   00   | 000000 00000
	 * 按照此方法，每秒每个节点可以产生的ID数量为: 2048个，即2^11个ID，三个节点每秒可以产生 2048 * 3 = 6144个ID
	 * 此方法可以确保生成的ID位数为12位的数字，若想在每秒生成更多的ID，则产生的ID的位数会大于12位
	 */
	maxNodeId          = -1 ^ (-1 << NodeIdBits) //节点 ID 最大范围
	nodeIdShift        = sequenceBits            //左移次数
	timestampLeftShift = sequenceBits + NodeIdBits
	sequenceMask       = -1 ^ (-1 << sequenceBits)
	maxNextIdsNum      = 100 //单次获取ID的最大数量
)

type IdWorker struct {
	sequence      int64 //序号
	lastTimestamp int64 //最后时间戳
	nodeId        int64 //节点ID
	twepoch       int64
	mutex         sync.Mutex
}

var idw *IdWorker = nil

// 预生产ID的channel队列，最多预生产10W个
var preGenChn = make(chan int64, 500000)

func GetIdWokrer() *IdWorker {
	if idw == nil && config.Config.BaseConf.NodeID > -1 {
		tidw, err := newIdWorker(config.Config.BaseConf.NodeID)
		if err != nil {
			logs.LogSystem.Errorf("GetIdWorker: %s\n", err.Error())
		} else {
			idw = tidw
			go preGen()
		}
	}
	return idw
}

// NewIdWorker new a idworker id generator object.
func newIdWorker(NodeId int64) (*IdWorker, error) {
	idWorker := &IdWorker{}
	if NodeId > maxNodeId || NodeId < 0 {
		_ = fmt.Sprintf("NodeId Id can't be greater than %d or less than 0", maxNodeId)
		return nil, errors.New(fmt.Sprintf("NodeId Id: %d error", NodeId))
	}

	idWorker.nodeId = NodeId
	idWorker.lastTimestamp = -1
	idWorker.sequence = 0
	idWorker.twepoch = twepoch
	idWorker.mutex = sync.Mutex{}
	_ = fmt.Sprintf("worker starting. timestamp left shift %d, worker id bits %d, sequence bits %d, workerid %d", timestampLeftShift, NodeIdBits, sequenceBits, NodeId)
	return idWorker, nil
}

func preGen() {
	defer func() {
		if err := recover(); err != nil {
			sentry.CurrentHub().Recover(err)
			sentry.Flush(time.Second * 3)
			logs.LogSystem.Error("idworker:pregen:", err)
		}
	}()
	for {
		if idw != nil {
			if config.Config.BaseConf.PreGen {
				id, err := idw.genNextId()
				if err != nil {
					logs.LogSystem.Error("idworker:pregen:genNextId:", err)
				} else {
					preGenChn <- id
				}
			} else {
				time.Sleep(time.Second * 30)
			}
		}
	}
}

// timeGen generate a unix millisecond.
func timeGen() int64 {
	return time.Now().UnixNano() / int64(time.Second)
}

// tilNextMillis spin wait till next millisecond.
func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeGen()
	for timestamp <= lastTimestamp {
		timestamp = timeGen()
	}
	return timestamp
}

func (id *IdWorker) GetNextId() (int64, error) {
	select {
	case id := <-preGenChn:
		return id, nil
	default:
		break
	}
	return id.genNextId()
}

// genNextId get a idworker id.
func (id *IdWorker) genNextId() (int64, error) {
	id.mutex.Lock()
	defer id.mutex.Unlock()
	return id.nextid()
}

// genNextIds get idworker ids.
func (id *IdWorker) genNextIds(num int) ([]int64, error) {
	if num > maxNextIdsNum || num < 0 {
		_ = fmt.Sprintf("genNextIds num can't be greater than %d or less than 0", maxNextIdsNum)
		return nil, errors.New(fmt.Sprintf("genNextIds num: %d error", num))
	}
	ids := make([]int64, num)
	id.mutex.Lock()
	defer id.mutex.Unlock()
	for i := 0; i < num; i++ {
		ids[i], _ = id.nextid()
	}
	return ids, nil
}

func (id *IdWorker) nextid() (int64, error) {
	timestamp := timeGen()
	if timestamp < id.lastTimestamp {
		//    fmt.Sprintf("clock is moving backwards.  Rejecting requests until %d.", id.lastTimestamp)
		return 0, errors.New(fmt.Sprintf("Clock moved backwards.  Refusing to generate id for %d milliseconds", id.lastTimestamp-timestamp))
	}
	if id.lastTimestamp == timestamp {
		id.sequence = (id.sequence + 1) & sequenceMask
		if id.sequence == 0 {
			timestamp = tilNextMillis(id.lastTimestamp)
		}
	} else {
		id.sequence = 0
	}
	id.lastTimestamp = timestamp
	return ((timestamp - id.twepoch) << timestampLeftShift) | (id.nodeId << nodeIdShift) | id.sequence, nil
}
