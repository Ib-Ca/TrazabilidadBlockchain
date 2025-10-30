package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// === Modelo ===
type Tolva struct {
	ID          string `json:"id"`     // identificador único
	Fecha       string `json:"fecha"`  // yyyy-mm-dd
	NOrden      string `json:"nOrden"` // número de orden
	NChapa      string `json:"nChapa"` // patente del camión
	Chofer      string `json:"chofer"`
	Origen      string `json:"origen"`
	Variedad    string `json:"variedad"`
	HoraInicio  string `json:"horaInicio"` // HH:MM
	HoraSalida  string `json:"horaSalida"` // HH:MM
	Observacion string `json:"observacion"`
}

// === Contrato ===
type SmartContract struct {
	contractapi.Contract
}

// RegistrarTolva: alta de una tolva
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

	b, _ := json.Marshal(t)
	if err := ctx.GetStub().PutState("TOLVA_"+t.ID, b); err != nil {
		return fmt.Errorf("error al guardar tolva: %v", err)
	}
	return ctx.GetStub().SetEvent("RegistrarTolva", b)
}

// ConsultarTolva: obtiene por ID
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

// ListarTolvas: devuelve todas (búsqueda por prefijo)
func (s *SmartContract) ListarTolvas(ctx contractapi.TransactionContextInterface) ([]Tolva, error) {
	it, err := ctx.GetStub().GetStateByRange("TOLVA_", "TOLVA_~") // ~ > zzz...
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

// helper
func (s *SmartContract) TolvaExiste(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState("TOLVA_" + id)
	if err != nil {
		return false, err
	}
	return data != nil, nil
}

func main() {
	cc, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error create chaincode: %s\n", err.Error())
		return
	}
	if err := cc.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s\n", err.Error())
	}
}
