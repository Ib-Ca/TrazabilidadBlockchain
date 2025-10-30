package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Tolva struct {
	DocType     string `json:"docType"`
	ID          string `json:"id"`
	Fecha       string `json:"fecha"`
	NOrden      string `json:"nOrden"`
	NChapa      string `json:"nChapa"`
	Chofer      string `json:"chofer"`
	Origen      string `json:"origen"`
	Variedad    string `json:"variedad"`
	HoraInicio  string `json:"horaInicio"`
	HoraSalida  string `json:"horaSalida"`
	Observacion string `json:"observacion"`
}

func (s *SmartContract) RegistrarTolva(ctx contractapi.TransactionContextInterface, tolvaJSON string) error {
	var t Tolva
	if err := json.Unmarshal([]byte(tolvaJSON), &t); err != nil {
		return fmt.Errorf("JSON inválido: %v", err)
	}
	if t.ID == "" {
		return fmt.Errorf("el campo 'id' es obligatorio")
	}
	exists, err := s.TolvaExiste(ctx, t.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("la tolva %s ya existe", t.ID)
	}
	if t.DocType == "" {
		t.DocType = "tolva"
	}

	data, _ := json.Marshal(t)
	if err := ctx.GetStub().PutState("TOLVA_"+t.ID, data); err != nil {
		return fmt.Errorf("error al guardar tolva: %v", err)
	}
	return ctx.GetStub().SetEvent("RegistrarTolva", data)
}

func (s *SmartContract) EditarTolva(ctx contractapi.TransactionContextInterface, tolvaJSON string) error {
	var nueva Tolva
	if err := json.Unmarshal([]byte(tolvaJSON), &nueva); err != nil {
		return fmt.Errorf("JSON inválido: %v", err)
	}
	if nueva.ID == "" {
		return fmt.Errorf("el campo 'id' es obligatorio")
	}
	exists, err := s.TolvaExiste(ctx, nueva.ID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("no se puede editar, la tolva %s no existe", nueva.ID)
	}
	if nueva.DocType == "" {
		nueva.DocType = "tolva"
	}

	data, _ := json.Marshal(nueva)
	if err := ctx.GetStub().PutState("TOLVA_"+nueva.ID, data); err != nil {
		return err
	}
	return ctx.GetStub().SetEvent("EditarTolva", data)
}

func (s *SmartContract) EliminarTolva(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.TolvaExiste(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("no se puede eliminar, la tolva %s no existe", id)
	}
	if err := ctx.GetStub().DelState("TOLVA_" + id); err != nil {
		return fmt.Errorf("error al eliminar tolva %s: %v", id, err)
	}
	return ctx.GetStub().SetEvent("EliminarTolva", []byte(id))
}

func (s *SmartContract) ConsultarTolva(ctx contractapi.TransactionContextInterface, id string) (*Tolva, error) {
	data, err := ctx.GetStub().GetState("TOLVA_" + id)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("tolva %s no encontrada", id)
	}
	var t Tolva
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *SmartContract) ListarTolvas(ctx contractapi.TransactionContextInterface) ([]Tolva, error) {
	q := `{"selector":{"docType":"tolva"}}`
	return s.queryTolvas(ctx, q)
}

func (s *SmartContract) BuscarTolvas(ctx contractapi.TransactionContextInterface, filtrosJSON string) ([]Tolva, error) {
	var filtros map[string]interface{}
	if err := json.Unmarshal([]byte(filtrosJSON), &filtros); err != nil {
		return nil, fmt.Errorf("filtros inválidos: %v", err)
	}

	selector := map[string]interface{}{
		"docType": "tolva",
	}
	for k, v := range filtros {
		selector[k] = v
	}

	query := map[string]interface{}{
		"selector": selector,
	}
	queryBytes, _ := json.Marshal(query)

	return s.queryTolvas(ctx, string(queryBytes))
}

func (s *SmartContract) TolvaExiste(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState("TOLVA_" + id)
	if err != nil {
		return false, err
	}
	return data != nil, nil
}

func (s *SmartContract) queryTolvas(ctx contractapi.TransactionContextInterface, query string) ([]Tolva, error) {
	it, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer it.Close()

	var out []Tolva
	for it.HasNext() {
		kv, err := it.Next()
		if err != nil {
			return nil, err
		}
		var t Tolva
		if err := json.Unmarshal(kv.Value, &t); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}
