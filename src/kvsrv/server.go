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

type OperationResult struct {
	RequestID int64
	Value string
}

type KVServer struct {
	mu sync.Mutex
	db map[string]string // key-value database
	lastClientOp map[int64]*OperationResult // last operation result for each client
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

	// Release useless lastClientOp memory
	delete(kv.lastClientOp, args.ClientID)
}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	// Handle duplicate request
	lastOp, ok := kv.lastClientOp[args.ClientID]
	if ok && lastOp.RequestID == args.RequestID {
		reply.Value = lastOp.Value
        return
	}
	
	kv.db[args.Key] = args.Value
	reply.Value = ""
	kv.lastClientOp[args.ClientID] = &OperationResult{RequestID: args.RequestID, Value: reply.Value}
}

func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	// Handle duplicate request
	lastOp, ok := kv.lastClientOp[args.ClientID]
	if ok && lastOp.RequestID == args.RequestID {
		reply.Value = lastOp.Value
        return
	}

	reply.Value = kv.db[args.Key]
	kv.db[args.Key] += args.Value
	kv.lastClientOp[args.ClientID] = &OperationResult{RequestID: args.RequestID, Value: reply.Value}
}

func StartKVServer() *KVServer {
	kv := new(KVServer)
	kv.db = make(map[string]string)
	kv.lastClientOp = make(map[int64]*OperationResult)
	return kv
}
