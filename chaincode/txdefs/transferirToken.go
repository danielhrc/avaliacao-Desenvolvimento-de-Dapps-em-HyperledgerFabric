package txdefs

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

// Create a new Library on channel
// POST Method
var TransferirToken = tx.Transaction{
	Tag:         "transferirToken",
	Label:       "TransferirToken",
	Description: "Transfere token",
	Method:      "PUT",
	Callers:     []string{"$org3MSP", "$orgMSP"}, // Only org3 can call this transaction

	Args: []tx.Argument{

		{
			/// token origem
			Tag:      "token",
			Label:    "Token origem",
			DataType: "->token",
		},
		{
			/// Proprietario
			Tag:      "proprietario",
			Label:    "Proprietario do token",
			DataType: "->proprietario",
		},
		{
			// Quantidade
			Tag:      "quantidade",
			Label:    "quantidade transferida",
			DataType: "number",
		},
			// Id novo token
		{
			Tag:         "id",
			Label:       "ID Novo Token",
			Description: "ID Novo Token",
			DataType:    "string",
			Required:    true,
		},
			// Id token origem 
		{
			Tag:         "novoId",
			Label:       "Novo ID Token Origem",
			Description: "Novo ID Token Origem",
			DataType:    "string",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		tokenKey, ok := req["token"].(assets.Key)

		if !ok {
			return nil, errors.WrapError(nil, "Parameter  must be an asset")
		}

		proprietario, ok := req["proprietario"].(assets.Key)

		if !ok {
			return nil, errors.WrapError(nil, "Parameter  must be an asset")
		}

		quantidade, ok := req["quantidade"].(float64)

		if !ok {
			return nil, errors.WrapError(nil, "Parameter  must be an asset")
		}

		id, ok := req["id"].(string)

		if !ok {
			return nil, errors.WrapError(nil, "Parameter  must be an asset")
		}

		// Returns token from channel
		tokenAsset, err := tokenKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to get asset from the ledger")
		}
		
		novoId, _ := req["novoId"].(string)
		
		tokenMap := (map[string]interface{})(*tokenAsset)
		
		if tokenMap["burned"].(bool) {
			return nil, errors.WrapError(err, "Already burned.")
		}
		
		// Update burned
		tokenMap["burned"] = true
		
		tokenMap, err = tokenAsset.Update(stub, tokenMap)
		if err != nil {
			return nil, errors.WrapError(err, "Falha ao atualizar ativo 'token'.")
		}

		novaQuantidade := tokenMap["quantidade"].(float64) - quantidade

		if novaQuantidade < 0 {
			return nil, errors.WrapError(err, "Saldo de token insuficiente.")
		}

		novoTokenOrigemMap := make(map[string]interface{})
		novoTokenOrigemMap["@assetType"] = "token"
		novoTokenOrigemMap["id"] = novoId
		novoTokenOrigemMap["proprietario"] = tokenMap["proprietario"]
		novoTokenOrigemMap["quantidade"] = novaQuantidade

		novoTokenOrigemAsset, err := assets.NewAsset(novoTokenOrigemMap)
		if err != nil {
			return nil, errors.WrapError(err, "Falha ao criar ativo 'novo token de origem'.")
		}

		_, err = novoTokenOrigemAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Erro ao salvar ativo 'novo token de origem' na blochchain.")
		}

		proprietarioKey, ok := req["destino"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parametro 'destino' deve ser um ativo.")
		}

		proprietarioAsset, err := proprietarioKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Falha ao obter ativo 'destino'.")
		}
		proprietarioMap := (map[string]interface{})(*proprietarioAsset)

		updatedProprietarioKey := make(map[string]interface{})
		updatedProprietarioKey["@assetType"] = "proprietario"
		updatedProprietarioKey["@key"] = proprietarioMap["@key"]

		novoTokenMap := make(map[string]interface{})
		novoTokenMap["@assetType"] = "token"
		novoTokenMap["id"] = id
		novoTokenMap["proprietario"] = updatedProprietarioKey
		novoTokenMap["quantidade"] = quantidade
		novoTokenMap["burned"] = false

		novoTokenAsset, err := assets.NewAsset(novoTokenMap)
		
		if err != nil {
			return nil, errors.WrapError(err, "Falha ao criar ativo 'token de destino'.")
		}
		// Save the new Token on channel
		_, err = TokenAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Error saving asset on blockchain")
		}

		// Marshal asset back to JSON format
		tokenJSON, nerr := json.Marshal(TokenAsset)
		if nerr != nil {
			return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		}

		return tokenJSON, nil
	},
}
