package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Secado struct {
	DocType     string  `json:"docType"`
	ID          string  `json:"id"`
	Fecha       string  `json:"fecha"`
	Hora        string  `json:"hora"`
	NroSecada   string  `json:"nrosecada"`
	VolumenKGR  float64 `json:"volumenkgr"`
	TempAire    float64 `json:"tempAire"`
	TempGrano   float64 `json:"tempGrano"`
	HumGrano    float64 `json:"humgrano"`
	Var         float64 `json:"var"`
	Destino     string  `json:"destino"`
	Observacion string  `json:"observacion"`
}

func (s *SmartContract) RegistrarSecado(ctx contractapi.TransactionContextInterface, secadoJSON string) error {
	var sc Secado
	if err := json.Unmarshal([]byte(secadoJSON), &sc); err != nil {
		return fmt.Errorf("JSON inválido: %v", err)
	}
	if sc.ID == "" {
		return fmt.Errorf("el campo 'id' es obligatorio")
	}
	if sc.DocType == "" {
		sc.DocType = "secado"
	}

	exists, err := s.SecadoExiste(ctx, sc.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("el secado %s ya existe", sc.ID)
	}

	data, _ := json.Marshal(sc)
	if err := ctx.GetStub().PutState("SECADO_"+sc.ID, data); err != nil {
		return fmt.Errorf("error al guardar secado: %v", err)
	}
	return ctx.GetStub().SetEvent("RegistrarSecado", data)
}

func (s *SmartContract) EditarSecado(ctx contractapi.TransactionContextInterface, secadoJSON string) error {
	var sc Secado
	if err := json.Unmarshal([]byte(secadoJSON), &sc); err != nil {
		return fmt.Errorf("JSON inválido: %v", err)
	}
	if sc.ID == "" {
		return fmt.Errorf("el campo 'id' es obligatorio")
	}
	exists, err := s.SecadoExiste(ctx, sc.ID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("no se puede editar, el secado %s no existe", sc.ID)
	}
	if sc.DocType == "" {
		sc.DocType = "secado"
	}

	data, _ := json.Marshal(sc)
	if err := ctx.GetStub().PutState("SECADO_"+sc.ID, data); err != nil {
		return err
	}
	return ctx.GetStub().SetEvent("EditarSecado", data)
}

func (s *SmartContract) EliminarSecado(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.SecadoExiste(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("no se puede eliminar, el secado %s no existe", id)
	}
	if err := ctx.GetStub().DelState("SECADO_" + id); err != nil {
		return fmt.Errorf("error al eliminar secado %s: %v", id, err)
	}
	return ctx.GetStub().SetEvent("EliminarSecado", []byte(id))
}

func (s *SmartContract) ConsultarSecado(ctx contractapi.TransactionContextInterface, id string) (*Secado, error) {
	b, err := ctx.GetStub().GetState("SECADO_" + id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, fmt.Errorf("secado %s no encontrado", id)
	}
	var sc Secado
	if err := json.Unmarshal(b, &sc); err != nil {
		return nil, err
	}
	return &sc, nil
}

//filtro couchdb

func (s *SmartContract) ListarSecados(ctx contractapi.TransactionContextInterface) ([]Secado, error) {
	q := `{"selector":{"docType":"secado"}}`
	return s.querySecados(ctx, q)
}

func (s *SmartContract) BuscarSecados(ctx contractapi.TransactionContextInterface, filtrosJSON string) ([]Secado, error) {
	var filtros map[string]interface{}
	if err := json.Unmarshal([]byte(filtrosJSON), &filtros); err != nil {
		return nil, fmt.Errorf("filtros inválidos: %v", err)
	}
	selector := map[string]interface{}{"docType": "secado"}
	for k, v := range filtros {
		selector[k] = v
	}
	query := map[string]interface{}{"selector": selector}
	qb, _ := json.Marshal(query)
	return s.querySecados(ctx, string(qb))
}

func (s *SmartContract) SecadoExiste(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	b, err := ctx.GetStub().GetState("SECADO_" + id)
	if err != nil {
		return false, err
	}
	return b != nil, nil
}

func (s *SmartContract) querySecados(ctx contractapi.TransactionContextInterface, q string) ([]Secado, error) {
	it, err := ctx.GetStub().GetQueryResult(q)
	if err != nil {
		return nil, err
	}
	defer it.Close()

	var out []Secado
	for it.HasNext() {
		kv, err := it.Next()
		if err != nil {
			return nil, err
		}
		var sc Secado
		if err := json.Unmarshal(kv.Value, &sc); err != nil {
			return nil, err
		}
		out = append(out, sc)
	}
	return out, nil
}
