package MyRPC

//Client Proxy
type EchoerProxy struct {
	aor       AbsoluteObjectReference
	requestor Requestor
}

func (e *EchoerProxy) Echo(line string) string {

	call := Call{"Echo", []interface{}{line}}
	newInv := Invocation{e.aor, call}
	newLine, err := e.requestor.Request(newInv)
	if err != nil {
		panic(err)
	}
	return newLine[0].(string)
}

func (e *EchoerProxy) ReverseEcho(line string) string {

	call := Call{"ReverseEcho", []interface{}{line}}
	newInv := Invocation{e.aor, call}
	newLine, err := e.requestor.Request(newInv)
	if err != nil {
		panic(err)
	}
	return newLine[0].(string)
}
