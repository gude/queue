package ddutil

import (
	"testing"
)

/*
 *  基本的出入队列测试
 */
func Test_Enqueue(t *testing.T) {
	testqueue = New(5)
	testqueue.Enqueue("key1", "value1")
	testqueue.Enqueue("key2", "value2")
	testqueue.Enqueue("key3", "value3")
	testqueue.Enqueue("key4", "value4")
	testqueue.Enqueue("key5", "value5")
	testqueue.Enqueue("key6", "value6")

	k, v, ok := testqueue.Dequeue()
	if !ok {
		t.Error("dequeue failed")
		return
	}
	if k.(string) != "key1" || v.(string) != "value1" {
		t.Error("dequeue expect key1 and value1, fact %s, %s", k.(string), v.(string))
	}
	testqueue.Dequeue()
	testqueue.Dequeue()
	testqueue.Dequeue()

	k, v, ok = testqueue.Dequeue()
	if !ok {
		t.Error("dequeue failed")
		return
	}
	if k.(string) != "key5" || v.(string) != "value5" {
		t.Error("dequeue expect key1 and value1, fact %s, %s", k.(string), v.(string))
	}
}

/*
   出队先进先出
   相同key，后入数据覆盖先入key，同时位置到队尾
*/
func Test_EnqueueOrder(t *testing.T) {
	testqueue = New(5)
	testqueue.Enqueue("key1", "value1")
	if testqueue.Len() != 1 {
		t.Error("queue len expected 1, fact %d", testqueue.Len())
	}
	testqueue.Enqueue("key2", "value2")
	testqueue.Enqueue("key3", "value3")
	testqueue.Enqueue("key4", "value4")
	testqueue.Enqueue("key5", "value5")
	testqueue.Enqueue("key6", "value6")
	if testqueue.Len() != 5 {
		t.Error("queue len expected 5, fact %d", testqueue.Len())
	}
	//出队   先进先出
	k, v, ok := testqueue.Dequeue()
	if !ok {
		t.Error("dequeue failed")
	}
	if k.(string) != "key1" || v.(string) != "value1" {
		t.Error("dequeue expect key1 and value1, fact %s, %s", k.(string), v.(string))
	}

	//相同key, 后入数据会覆盖先入数据，但位置会变
	testqueue.Enqueue("key2", "testvalue2")
	testqueue.Dequeue()
	testqueue.Dequeue()
	testqueue.Dequeue()
	k, v, ok = testqueue.Dequeue()
	if !ok {
		t.Error("dequeue failed")
	}
	if k.(string) != "key2" || v.(string) != "testvalue2" {
		t.Error("dequeue expect key1 and value1, fact %s, %s", k.(string), v.(string))
	}
}

/*
   基本的pushfront，数据放入队头
*/
func Test_PushFront(t *testing.T) {
	testqueue = New(5)
	testqueue.Enqueue("key1", "value1")
	if testqueue.Len() != 1 {
		t.Error("queue len expected 1, fact %d", testqueue.Len())
	}
	testqueue.Enqueue("key2", "value2")
	testqueue.Enqueue("key3", "value3")
	testqueue.Enqueue("key4", "value4")
	testqueue.Enqueue("key5", "value5")
	testqueue.Enqueue("key6", "value6")
	if testqueue.Len() != 5 {
		t.Error("queue len expected 5, fact %d", testqueue.Len())
	}
	//出队   先进先出
	k, v, ok := testqueue.Dequeue()
	if !ok {
		t.Error("dequeue failed")
	}
	if k.(string) != "key1" || v.(string) != "value1" {
		t.Error("dequeue expect key1 and value1, fact %s, %s", k.(string), v.(string))
	}

	//PushFront数据会在队头
	testqueue.PushFront(k, v)
	k, v, ok = testqueue.Dequeue()
	if !ok {
		t.Error("dequeue failed")
	}
	if k.(string) != "key1" || v.(string) != "value1" {
		t.Error("dequeue expect key1 and value1, fact %s, %s", k.(string), v.(string))
	}

}

/*
	pushfront已有key到队列中，忽略pushfront数据
*/
func Test_PushFrontExistsKeys(t *testing.T) {
	testqueue = New(5)
	testqueue.Enqueue("key1", "value1")
	if testqueue.Len() != 1 {
		t.Error("queue len expected 1, fact %d", testqueue.Len())
	}
	testqueue.Enqueue("key2", "value2")
	testqueue.Enqueue("key3", "value3")
	testqueue.Enqueue("key4", "value4")
	testqueue.Enqueue("key5", "value5")
	testqueue.Enqueue("key6", "value6")
	if testqueue.Len() != 5 {
		t.Error("queue len expected 5, fact %d", testqueue.Len())
	}
	//出队
	k, v, _ := testqueue.Dequeue()
	testqueue.Enqueue(k, "testvalue1")
	//pushfront时，如果队列中存在相同key, 则pushfront数据被丢弃
	testqueue.PushFront(k, v)
	testqueue.Dequeue()
	testqueue.Dequeue()
	testqueue.Dequeue()
	testqueue.Dequeue()

	nk, nv, _ := testqueue.Dequeue()
	if nk != k || nv != "testvalue1" {
		t.Error("dequeue expect key %s, value %s, fact key %s, value %s", k, "testvalue1", nk, nv)
	}
}
