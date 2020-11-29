package MyRPC

import (
	"errors"
)

type AbsoluteObjectReference struct {
	Address string
	ID      uint
}

//interface do LookUp
type LookUper interface {
	Register(serviceName string, reference AbsoluteObjectReference) (AbsoluteObjectReference, error)
	Remove(serviceName string) error
	LookUp(serviceName string) (AbsoluteObjectReference, error)
	List() ([]string, error)
	CreateReference(address string, id uint) AbsoluteObjectReference
	Init(address string) (LookUpInvoker, error)
	Stop() error
}

type LookUp struct {
	Services map[string]AbsoluteObjectReference
	SRH      ServerRequestHandler
}

func (lookup *LookUp) Register(serviceName string, reference AbsoluteObjectReference) (AbsoluteObjectReference, error) {
	aor, achou := lookup.Services[serviceName]
	if achou {
		return aor, errors.New("Já tem servico com o nome")
	}
	lookup.Services[serviceName] = reference
	return reference, nil
}
func (lookup *LookUp) Remove(serviceName string) error {
	_, achou := lookup.Services[serviceName]
	if !achou {
		return errors.New("Não tem servico com o nome")
	}
	delete(lookup.Services, serviceName)
	return nil
}

func (lookup *LookUp) LookUp(serviceName string) (AbsoluteObjectReference, error) {
	aor, achou := lookup.Services[serviceName]
	if !achou {
		return aor, errors.New("Não tem servico com o nome")
	}
	return aor, nil
}

func (lookup *LookUp) List() ([]string, error) {
	keys := make([]string, len(lookup.Services))

	i := 0
	for k := range lookup.Services {
		keys[i] = k
		i++
	}
	if i == 0 {
		return nil, errors.New("Não há serviço para listar")
	}
	return keys, nil
}
func (lookup *LookUp) CreateReference(address string, id uint) AbsoluteObjectReference {
	return AbsoluteObjectReference{address, id}
}

func (lookup *LookUp) Init(address string) (LookUpInvoker, error) {
	if lookup == nil {
		lookup = new(LookUp)
		lookup.Services = nil
		lookup.SRH = ServerRequestHandlerTCP{}
	}
	if lookup.Services == nil {
		lookup.Services = make(map[string]AbsoluteObjectReference)
	}
	lookup.Register("Lookup", lookup.CreateReference(address, 1))
	var inv Invoker = LookUpInvoker{lookup}
	lookup.SRH.SetUp(&inv, address)
	println("UHUL")
	lookup.SRH.Listen()
	println("UHUL")
	return inv.(LookUpInvoker), nil
}

func (lookup *LookUp) Stop() error {
	lookup.SRH.Close()
	return nil
}
func (lookup *LookUp) Close() error {
	lookup.SRH.Close()
	lookup.SRH.TearDown()
	return nil
}

type LookUpInvoker struct {
	Lookup *LookUp
}

func (inv LookUpInvoker) Invoke(message []byte) ([]byte, error) {
	m := Marshaller{}
	op := Invocation{}
	err := m.Unmarshal(message, &op)
	if err != nil {
		return nil, err
	}

	switch op.Call.Method {
	case "Register":
		aor, err := (*inv.Lookup).Register(op.Call.Args[0].(string), op.AOR)
		if err != nil {
			return nil, err
		}
		return m.Marshal(aor)
	case "Remove":
		err := (*inv.Lookup).Remove(op.Call.Args[0].(string))
		if err != nil {
			return nil, err
		}
		return m.Marshal(nil) //parece errado, mas é isso mesmo
	case "LookUp":
		aor, err := (*inv.Lookup).LookUp(op.Call.Args[0].(string))
		if err != nil {
			return nil, errors.New("servico nao encontrado")
		}
		return m.Marshal(aor)
	case "List":
		aors, err := (*inv.Lookup).List()
		if err != nil {
			return nil, errors.New("nenhum servico encontrado")
		}
		return m.Marshal(aors)
	case "CreateReference":
		return nil, errors.New("CreateReference() deve ser implementada no client")
	case "Init":
		return nil, errors.New("Init() não deve ser chamada via invoker")
	case "Stop":
		return nil, errors.New("Stop() não deve ser chamada via invoker")
	default:
		return nil, errors.New("Operação não reconhecida")
	}
}

type LookUpProxy struct {
	AOR       AbsoluteObjectReference
	Requestor Requestor
}

func (lookup *LookUpProxy) New(address string) *LookUpProxy {
	aor := lookup.CreateReference(address, 1)
	*lookup = LookUpProxy{aor, Requestor{}} //ID fixo do lookup
	lookup.Requestor.CRH = ClientRequestHandlerTCP{}
	return lookup
}

func (lookup *LookUpProxy) Register(serviceName string, reference AbsoluteObjectReference) (AbsoluteObjectReference, error) {

	call := Call{"Register", []interface{}{serviceName, reference}}
	newInv := Invocation{lookup.AOR, call}
	ret, err := lookup.Requestor.Request(newInv)
	if err != nil {
		panic(err)
	}
	return ret[0].(AbsoluteObjectReference), nil
}

func (lookup *LookUpProxy) Remove(serviceName string) error {

	call := Call{"Remove", []interface{}{serviceName}}
	newInv := Invocation{lookup.AOR, call}
	_, err := lookup.Requestor.Request(newInv)
	if err != nil {
		panic(err)
	}
	return nil
}

func (lookup *LookUpProxy) LookUp(serviceName string) (AbsoluteObjectReference, error) {

	call := Call{"LookUp", []interface{}{serviceName}}
	newInv := Invocation{lookup.AOR, call}
	newLine, err := lookup.Requestor.Request(newInv)
	if err != nil {
		panic(err)
	}
	return newLine[0].(AbsoluteObjectReference), nil
}

func (lookup *LookUpProxy) List() ([]string, error) {
	call := Call{"List", []interface{}{}}
	newInv := Invocation{lookup.AOR, call}
	newLine, err := lookup.Requestor.Request(newInv)
	if err != nil {
		panic(err)
	}
	ret := make([]string, len(newLine))
	for i, s := range newLine {
		ret[i] = s.(string)
	}
	return ret, nil
}

//retorna localmente.
func (lookup *LookUpProxy) CreateReference(address string, id uint) AbsoluteObjectReference {
	return AbsoluteObjectReference{address, id}
}

func (lookup *LookUpProxy) Init(address string) (LookUpInvoker, error) {
	panic(errors.New("init() não pode ser chamado pelo cliente"))
}

func (lookup *LookUpProxy) Stop() error {
	panic(errors.New("Stop() não pode ser chamado pelo cliente"))
}
