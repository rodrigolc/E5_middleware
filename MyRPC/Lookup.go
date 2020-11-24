package MyRPC

import (
	"errors"
)

type AbsoluteObjectReference struct {
	address string
	id      uint
}

//interface do LookUp
type LookUper interface {
	Register(serviceName string, reference AbsoluteObjectReference) (AbsoluteObjectReference, error)
	Remove(serviceName string) error
	LookUp(serviceName string) (AbsoluteObjectReference, error)

	CreateReference(address string, id uint) AbsoluteObjectReference
	Init(address string) (LookUpInvoker, error)
	Stop() error
}

type LookUp struct {
	services map[string]AbsoluteObjectReference
	srh      ServerRequestHandler
}

func (lookup *LookUp) Register(serviceName string, reference AbsoluteObjectReference) (AbsoluteObjectReference, error) {
	aor, achou := lookup.services[serviceName]
	if achou {
		return aor, errors.New("Já tem servico com o nome")
	}
	lookup.services[serviceName] = reference
	return reference, nil
}
func (lookup *LookUp) Remove(serviceName string) error {
	_, achou := lookup.services[serviceName]
	if !achou {
		return errors.New("Não tem servico com o nome")
	}
	delete(lookup.services, serviceName)
	return nil
}

func (lookup *LookUp) LookUp(serviceName string) (AbsoluteObjectReference, error) {
	aor, achou := lookup.services[serviceName]
	if !achou {
		return aor, errors.New("Não tem servico com o nome")
	}
	return aor, nil
}

func (lookup *LookUp) CreateReference(address string, id uint) AbsoluteObjectReference {
	return AbsoluteObjectReference{address, id}
}

func (lookup *LookUp) Init(address string) (LookUpInvoker, error) {
	if lookup == nil {
		lookup = new(LookUp)
		lookup.services = nil
		lookup.srh = ServerRequestHandlerTCP{}
	}
	if lookup.services == nil {
		lookup.services = make(map[string]AbsoluteObjectReference)
	}
	lookup.Register("Lookup", lookup.CreateReference(address, 1))
	var inv Invoker = LookUpInvoker{lookup}
	lookup.srh.SetUp(&inv, address)
	println("UHUL")
	lookup.srh.Listen()
	println("UHUL")
	return inv.(LookUpInvoker), nil
}

func (lookup *LookUp) Stop() error {
	lookup.srh.Close()
	return nil
}
func (l *LookUp) Close() error {
	l.srh.Close()
	l.srh.TearDown()
	return nil
}

type LookUpInvoker struct {
	lookup *LookUp
}

func (inv LookUpInvoker) Invoke(message []byte) ([]byte, error) {
	m := Marshaller{}
	op := Invocation{}
	err := m.Unmarshal(message, &op)
	if err != nil {
		return nil, err
	}

	switch op.call.method {
	case "Register":
		aor, err := (*inv.lookup).Register(op.call.args[0].(string), op.aor)
		if err != nil {
			return nil, err
		}
		return m.Marshal(aor)
	case "Remove":
		err := (*inv.lookup).Remove(op.call.args[0].(string))
		if err != nil {
			return nil, err
		}
		return nil, nil //parece errado, mas é isso mesmo
	case "LookUp":
		aor, err := (*inv.lookup).LookUp(op.call.args[0].(string))
		if err != nil {
			return nil, errors.New("Servico nao encontrado.")
		}
		return m.Marshal(aor)
	case "CreateReference":
		return nil, errors.New("CreateReference() deve ser implementada no cliente.")
	case "Init":
		return nil, errors.New("Init() não deve ser chamada via invoker.")
	case "Stop":
		return nil, errors.New("Stop() não deve ser chamada via invoker.")
	default:
		return nil, errors.New("Operação não reconhecida.")
	}
}

type LookUpProxy struct {
	aor       AbsoluteObjectReference
	requestor Requestor
}

func (e *LookUpProxy) Register(serviceName string, reference AbsoluteObjectReference) (AbsoluteObjectReference, error) {

	call := Call{"Register", []interface{}{serviceName, reference}}
	newInv := Invocation{e.aor, call}
	ret, err := e.requestor.Request(newInv)
	if err != nil {
		panic(err)
	}
	return ret[0].(AbsoluteObjectReference), nil
}

func (e *LookUpProxy) Remove(serviceName string) error {

	call := Call{"Remove", []interface{}{serviceName}}
	newInv := Invocation{e.aor, call}
	_, err := e.requestor.Request(newInv)
	if err != nil {
		panic(err)
	}
	return nil
}

func (e *LookUpProxy) LookUp(serviceName string) (AbsoluteObjectReference, error) {

	call := Call{"Register", []interface{}{serviceName}}
	newInv := Invocation{e.aor, call}
	newLine, err := e.requestor.Request(newInv)
	if err != nil {
		panic(err)
	}
	return newLine[0].(AbsoluteObjectReference), nil
}

//retorna localmente.
func (e *LookUpProxy) CreateReference(address string, id uint) AbsoluteObjectReference {
	return AbsoluteObjectReference{address, id}
}

func (e *LookUpProxy) Init(address string) (LookUpInvoker, error) {
	panic(errors.New("init() não pode ser chamado pelo cliente"))
}

func (e *LookUpProxy) Stop() error {
	panic(errors.New("Stop() não pode ser chamado pelo cliente"))
}
