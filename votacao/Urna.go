/*
 * Filename: /mnt/c/Aulas/PCD/E5_Middleware/votacao/Urna.go
 * Path: /mnt/c/Aulas/PCD/E5_Middleware/votacao
 * Created Date: Friday, November 20th 2020, 2:53:43 pm
 * Author: Rodrigo Lopes (rlc2)
 *
 * Copyright (c) 2020 Your Company
 */

package votacao

type Eleicao struct {
	ID    int
	Cargo string
}

type Candidato struct {
	ID        int
	Nome      string
	PartidoId int
}

type Partido struct {
	ID    int
	Nome  string
	Sigla string
}

type UrnaInterface interface { //interface para a Urna, software de votacao, será implementado como
	RegisterEleicao(newEleicao Eleicao) Eleicao         // registra e retorna objeto com info de cadastro da eleicao. deve ter a mesma ID enviada no objeto, ou aleatória se vazio. Falha se já existir
	RegisterPartido(newPartido Partido) Partido         // registra e retorna objeto com info de cadastro do partido. deve ter a mesma ID enviada no objeto, ou aleatória se vazio. Falha se já existir
	RegisterCandidato(newCandidato Candidato) Candidato // registra e retorna objeto com info de cadastro do candidato. deve ter a mesma ID enviada no objeto, ou aleatória se vazio. Falha se já existir
	RegisterVoto(eleicaoID, candidatoID int)            // registra voto na eleicao para o candidato.
	RegisterVotoPartido(eleicaoID, partidoID int)       // registra voto na eleicao para o partido.

	GetCandidato(id int) Candidato                  //retorna um candidato específico
	GetCandidatos(ids []int) []Candidato            //retorna candidatos específicos
	GetAllCandidatos() []Candidato                  //retorna todos os candidatos
	GetCandidatosPartido(partidoID int) []Candidato //retorna candidatos do partido

	RemoveCandidato(id int) Candidato //deve retornar o candidato, ou nil se nao existe. não deve dar erro.
	RemovePartido(id int) Partido     //idem.
	RemoveEleicao(id int) Eleicao     //idem.

	UpdateCandidato(newCandidato Candidato) int //modifica e retorna o ID. deve ser o mesmo ID, e deve ja estar registrado
	UpdatePartido(newPartido Partido) int       //idem, partido
	UpdateEleicao(newEleicao Eleicao) int       //idem, eleicao

	GetVotos(idCandidato int) int
	GetVotosPartido(idPartido int) int
}
