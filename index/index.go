package index

import (
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/graph/g"
	"github.com/open-falcon/graph/rrdtool"
	"log"
)

// 初始化索引功能模块
func Start() {
	InitCache()
	go StartIndexUpdateIncrTask()
	log.Println("index:Start, ok")
}

// index收到一条新上报的监控数据,尝试用于增量更新索引
func ReceiveItem(item *model.GraphItem, md5 string) {
	if item == nil {
		return
	}

	uuid := item.UUID()

	// 已上报过的数据
	if indexedItemCache.ContainsKey(md5) {
		old := indexedItemCache.Get(md5).(*IndexCacheItem)
		if uuid == old.UUID { // dsType+step没有发生变化,只更新缓存
			old.Item = item
		} else { // dsType+step变化了,当成一个新的增量来处理(甚至,不用rrd文件来过滤)
			//indexedItemCache.Remove(md5)
			unIndexedItemCache.Put(md5, NewIndexCacheItem(uuid, item))
		}
		return
	}

	// 是否有rrdtool文件存在,如果有 认为已建立索引
	// 针对 索引缓存重建场景 做的优化, 结合索引全量更新 来保证一致性
	rrdFileName := rrdtool.RrdFileName(g.Config().RRD.Storage, item, md5)
	if rrdtool.IsRrdFileExist(rrdFileName) {
		indexedItemCache.Put(md5, NewIndexCacheItem(uuid, item))
		return
	}

	// 缓存未命中, 放入增量更新队列
	unIndexedItemCache.Put(md5, NewIndexCacheItem(uuid, item))
}