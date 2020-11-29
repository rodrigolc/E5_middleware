package MyRPC

//Requestor
type Requestor struct {
	CRH ClientRequestHandler
}

type Call struct {
	Method string
	Args   []interface{}
}

type Invocation struct {
	AOR  AbsoluteObjectReference
	Call Call
}

func (r *Requestor) Init() {

}

func (r *Requestor) Request(invocation Invocation) ([]interface{}, error) {
	var err error
	m := Marshaller{}
	r.CRH, err = r.CRH.SetUp(invocation.AOR.Address)
	if err != nil {
		return nil, err
	}
	data, err := m.Marshal(invocation)
	if err != nil {
		return nil, err
	}
	response, err := r.CRH.SendReceive(data)
	if err != nil {
		return nil, err
	}
	var obj interface{}
	err = m.Unmarshal(response, &obj)
	return []interface{}{obj}, err
}
