package MyRPC

//Requestor
type Requestor struct {
	crh ClientRequestHandler
}

type Call struct {
	Method string
	Args   []interface{}
}

type Invocation struct {
	aor  AbsoluteObjectReference
	call Call
}

func (r *Requestor) Request(invocation Invocation) ([]interface{}, error) {
	var err error
	m := Marshaller{}
	r.crh, err = r.crh.SetUp(invocation.aor.address)
	if err != nil {
		return nil, err
	}
	data, err := m.Marshal(invocation)
	if err != nil {
		return nil, err
	}
	response, err := r.crh.SendReceive(data)
	if err != nil {
		return nil, err
	}
	var obj interface{}
	err = m.Unmarshal(response, &obj)
	return []interface{}{obj}, err
}
