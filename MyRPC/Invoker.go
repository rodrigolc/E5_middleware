package MyRPC

//Invoker ->
type Invoker interface {
	Invoke(message []byte) ([]byte, error)
}

//Absolute Object Reference (AOR)
//Client Proxy
//Requestor
//Marshaller
//Invoker
//Lookup
