package txdefs

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

// Return the number of tokens from a proprietario
// GET method
var GetNumberOfTokensFromProprietario = tx.Transaction{
	Tag:         "getNumberOftokensFromProprietario",
	Label:       "Get Number Oftokens from proprietarios",
	Description: "Return the number of tokens from proprietarios",
	Method:      "GET",

	Args: []tx.Argument{
		{
			Tag:         "Proprietario",
			Label:       "Proprietario",
			Description: "Proprietario",
			DataType:    "->proprietario",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		idProprietario, ok := req["proprietario"].(assets.Key)

		if !ok {
			return nil, errors.WrapError(nil, "Parametro proprietario deve ser um ativo.")
		}

		proprietarioAsset, errKey := idProprietario.Get(stub)
	
		
		if errKey != nil {
			return nil, errors.WrapError(errKey, "Falha ao obter proprietario.")
		}
		proprietarioMap := (map[string]interface{})(*proprietarioAsset)

		searchProp := make(map[string]interface{})
		searchProp["@assetType"] = "proprietario"
		searchProp["@key"] = proprietarioMap["@key"]

		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType":   "token",
				"proprietario": searchProp,
			},
		}

		var err error
		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "Erro searching proprietário.", 500)
		}

		tokens := response.Result

		var quantidade float64 = 0

		for i := 0; i < len(tokens); i++ {
			if !tokens[i]["burned"].(bool) {
				quantidade = quantidade + tokens[i]["quantidade"].(float64)
			}
		}

		balance := make(map[string]interface{})
		
		balance["quantidade"] = quantidade
		
		responseJSON, err := json.Marshal(balance)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error marshaling response", 500)
		}

		return responseJSON, nil

	},
}
