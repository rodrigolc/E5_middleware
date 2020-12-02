package MyRPC

import (
	"errors"
	"fmt"
)

type AbsoluteObjectReference struct {
	Address string
	ID      string
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
	//println("lookup! register! string:", serviceName, " Address: ", reference.Address, "-", reference.ID)
	aor, achou := lookup.Services[serviceName]

	if achou {
		//println("lookup! register! string:", serviceName, " found!: ", aor.Address, "-", aor.ID)
		return aor, errors.New("Já tem servico com o nome")
	}
	//println("lookup! register! string:", serviceName, " not found, go on.")
	lookup.Services[serviceName] = reference
	//println("lookup! register! string:", serviceName, " registered: ", reference.Address, "-", reference.ID)
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
func (lookup *LookUp) CreateReference(address string, id string) AbsoluteObjectReference {
	return AbsoluteObjectReference{address, id}
}

func (lookup *LookUp) Init(address string) (LookUpInvoker, error) {

	if lookup == nil {
		lookup = new(LookUp)
		lookup.Services = nil
		lookup.SRH = ServerRequestHandlerTCP{}
	}
	if lookup.SRH == nil {
		lookup.SRH = ServerRequestHandlerTCP{}
	}
	if lookup.Services == nil {
		lookup.Services = make(map[string]AbsoluteObjectReference)
	}
	lookup.Register("Lookup", lookup.CreateReference(address, "1"))
	var inv Invoker = LookUpInvoker{lookup}
	var err error
	lookup.SRH, err = lookup.SRH.SetUp(&inv, address)
	if err != nil {
		panic(err)
	}
	lookup.SRH.Listen()
	return inv.(LookUpInvoker), err
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
		args := op.Call.Args[1].(map[string]interface{})
		tAOR := AbsoluteObjectReference{args["Address"].(string), args["ID"].(string)}
		aor, err := (*inv.Lookup).Register(op.Call.Args[0].(string), tAOR)
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
			return nil, err
		}
		return m.Marshal(aor)
	case "List":
		aors, err := (*inv.Lookup).List()
		if err != nil {
			return nil, err
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
	aor := lookup.CreateReference(address, "1")
	*lookup = LookUpProxy{aor, Requestor{}} //ID fixo do lookup
	lookup.Requestor.CRH = ClientRequestHandlerTCP{}
	return lookup
}

func (lookup *LookUpProxy) Register(serviceName string, reference AbsoluteObjectReference) (AbsoluteObjectReference, error) {
	//println("lookup proxy! register?")
	call := Call{"Register", []interface{}{serviceName, reference}}
	newInv := Invocation{lookup.AOR, call}
	//println("lookup proxy! register? request?")
	ret, err := lookup.Requestor.Request(newInv)

	//println("lookup proxy! register? request!", ret, err)
	if err != nil {
		//println("lookup proxy! register? request! error!!!", ret, err)
		fmt.Println(err)
	}
	var aor map[string]interface{} = ret[0].(map[string]interface{})
	return AbsoluteObjectReference{aor["Address"].(string), aor["ID"].(string)}, nil
}

func (lookup *LookUpProxy) Remove(serviceName string) error {

	call := Call{"Remove", []interface{}{serviceName}}
	newInv := Invocation{lookup.AOR, call}
	_, err := lookup.Requestor.Request(newInv)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (lookup *LookUpProxy) LookUp(serviceName string) (AbsoluteObjectReference, error) {

	call := Call{"LookUp", []interface{}{serviceName}}
	newInv := Invocation{lookup.AOR, call}
	ret, err := lookup.Requestor.Request(newInv)
	if err != nil {
		fmt.Println(err)
	}
	var aor map[string]interface{} = ret[0].(map[string]interface{})
	return AbsoluteObjectReference{aor["Address"].(string), aor["ID"].(string)}, nil
}

func (lookup *LookUpProxy) List() ([]string, error) {
	call := Call{"List", []interface{}{}}
	newInv := Invocation{lookup.AOR, call}
	newLine, err := lookup.Requestor.Request(newInv)
	if err != nil {
		fmt.Println(err)
	}
	ret := make([]string, len(newLine))
	for i, s := range newLine {
		ret[i] = s.(string)
	}
	return ret, nil
}

//retorna localmente.
func (lookup *LookUpProxy) CreateReference(address string, id string) AbsoluteObjectReference {
	return AbsoluteObjectReference{address, id}
}

func (lookup *LookUpProxy) Init(address string) (LookUpInvoker, error) {
	return LookUpInvoker{}, errors.New("init() não pode ser chamado pelo cliente")
}

func (lookup *LookUpProxy) Stop() error {
	return errors.New("Stop() não pode ser chamado pelo cliente")
}
