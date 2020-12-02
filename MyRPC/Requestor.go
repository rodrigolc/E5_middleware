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
	//println("Requestor! request! invocation?")
	var err error
	m := Marshaller{}
	//println("Requestor! request! setup crh?")
	r.CRH, err = r.CRH.SetUp(invocation.AOR.Address)
	//println("Requestor! request! setup crh!")

	if err != nil {
		return nil, err
	}
	data, err := m.Marshal(invocation)
	if err != nil {
		return nil, err
	}
	//println("Requestor! request! sendreceive?")
	response, err := r.CRH.SendReceive(data)
	//println("Requestor! request! sendreceive!")
	if err != nil {
		//println("Requestor! request! sendreceive! errorrr!")
		return nil, err
	}
	//println("Requestor! request! response!", response, string(response), err)
	var obj interface{}
	err = m.Unmarshal(response, &obj)
	//println("Requestor! request! unmarshall!", obj, err)

	return []interface{}{obj}, err
}
