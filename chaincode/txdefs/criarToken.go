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
var CreateNewToken = tx.Transaction{
	Tag:         "criarToken",
	Label:       "CriarToken",
	Description: "Cria um novo token",
	Method:      "POST",

	Args: []tx.Argument{
		{

			Required:    true,
			Tag:         "id",
			Label:       "Id do token",
			Description: "id do Token",
			DataType:    "string",
		},
		{
			/// Reference to another asset
			Tag:      "proprietario",
			Label:    "Proprietario do token",
			DataType: "->proprietario",
		},
		{
			// Quantidade
			Tag:      "quantidade",
			Label:    "quantidade",
			DataType: "number",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		id, _ := req["id"].(string)
		proprietario, _ := req["proprietario"]
		quantidade, _ := req["quantidade"].(json.Number)

		tokenMap := make(map[string]interface{})
		tokenMap["@assetType"] = "Token"
		tokenMap["id"] = id
		tokenMap["proprietario"] = proprietario
		tokenMap["quantidade"] = quantidade

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
