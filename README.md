## Raft

This project is based on MIT 6.5840 Distributed Systems.

###  Lab 2: Key/Value Server

Key Requirements and Corresponding Solutions:
1) Handle dropped messages.
   - Client side: Add ClientID, RequestID in the requests, so the server can track if the request is duplicate or not.
    - Server side: Use `lastClientOp` to record the last request operation of each client. If the request have the same ID in the lastClientOp, then we can filter it.(This lab assumes that a client will make only one call into a Clerk at a time.)
2) Reduce memory usage at the server.
   - When the client use `GET` operation, we can delete the record in `lastClientOp`. Another method may be also useful if a client receive the reply of Put and Append, then the client can send another request to ask for removal the record in `lastClientOp`. 
