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
var TransferToken = tx.Transaction{
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
		{
			// Composite Key
			Tag:      "id",
			Label:    "Id do token",
			DataType: "string",
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

		quantidade, ok := req["quantidade"].(json.Number)

		if !ok {
			return nil, errors.WrapError(nil, "Parameter  must be an asset")
		}

		id, ok := req["id"].(json.Number)

		if !ok {
			return nil, errors.WrapError(nil, "Parameter  must be an asset")
		}

		// Returns token from channel
		tokenAsset, err := tokenKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to get asset from the ledger")
		}

		tokenMap := (map[string]interface{})(*tokenAsset)
		tokenMap["@assetType"] = "Token"
		tokenMap["token"] = tokenKey
		tokenMap["id"] = id
		tokenMap["proprietario"] = proprietario
		tokenMap["quantidade"] = quantidade
		tokenMap["burned"] = false

		updateTokenKey := make(map[string]interface{})
		updateTokenKey["@assetType"] = "token"
		updateTokenKey["@key"] = tokenMap["@key"]

		// Returns old token from channel
		oldTokenAsset, err := tokenKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to get asset from the ledger")
		}

		oldTokenMap := (map[string]interface{})(*oldTokenAsset)
		oldTokenMap["@assetType"] = "token"
		oldTokenMap["@key"] = tokenKey

		// Update burned
		oldTokenMap["burned"] = true

		oldTokenMap, err = tokenAsset.Update(stub, oldTokenMap)
		if err != nil {
			return nil, errors.WrapError(err, "failed to update asset")
		}

		TokenAsset, err := assets.NewAsset(tokenMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create a new asset")
		}

		if tokenMap["quantidade"] == 0 {
			return nil, errors.WrapError(err, "Quantidade não pode ser 0")
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
