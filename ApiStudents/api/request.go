package api //pacote para validação dos dados.

import (
	"fmt"
)

//struct studentRequest -> terá todos os dados que esta no schemas.go, exceção do gorm.model
type StudentRequest struct {

	Name string `json:"Name"` //temos que dizer como é nome do campo vai querer atrelar no json que recebe no post
	CPF int `json:"CPF"`
	Email string `json:"Email"`
	Age int `json:"Age"`
	Active *bool `json:"Active"` //usa o bool como ponteiro para forçar o bool
}	
//cria uma função para informar o erro caso o campo não seja preenchido conforme adquado, para ser chamado na proxima função.
func errParamRequired(param, typ string) error { //param e typ tipo string -> retorna error
	return fmt.Errorf("param '%s' of type '%s' is required", param, typ) //caso o usuario não preencha o campo, msg será exibida
}
//função para validação para os campos de students que irão chegar no Request 
func (s *StudentRequest) Validate() error { //atrela a struct studentRequest para poder usar os dados.
	//verifica se o campo foi passado ou não (esta em branco ou não).
	if s.Name == "" { //se o campo Name estiver vazio
		return errParamRequired("Name", "String") //exibe a msg na função errParamRequired, solicitando o preenchimento corretamente
	}

	if s.CPF == 0 { //se o campo cpf estiver vazio
		return errParamRequired("CPF", "int") //exibe a msg na função errParamRequired, solicitando o preenchimento corretamente
	}

	if s.Email == "" { //se o campo email estiver vazio
		return errParamRequired("Email", "String") //exibe a msg na função errParamRequired, solicitando o preenchimento corretamente
	}

	if s.Age == 0 { //se o campo age estiver vazio
		return errParamRequired("Age", "int") //exibe a msg na função errParamRequired, solicitando o preenchimento corretamente
	}

	if s.Active == nil { //se o campo email estiver vazio
		return errParamRequired("Active", "bool") //exibe a msg na função errParamRequired, solicitando o preenchimento corretamente
	}
	return nil

}