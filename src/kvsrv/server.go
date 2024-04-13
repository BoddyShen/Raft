package kvsrv

import (
	"log"
	"sync"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}


type KVServer struct {
	mu sync.Mutex
	db map[string]string
}


func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	value, ok := kv.db[args.Key]
	if !ok {
		reply.Value = ""
	} else {
		reply.Value = value
	}
}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.db[args.Key] = args.Value
	reply.Value = args.Value
}

func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	reply.Value = kv.db[args.Key]
	kv.db[args.Key] += args.Value
}

func StartKVServer() *KVServer {
	kv := new(KVServer)
	kv.db = make(map[string]string)
	return kv
}
