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

		idProprietario, _ := req["proprietario"].(assets.Key)

		// Prepare couchdb query
		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType":   "token",
				"proprietario": idProprietario,
			},
		}

		var err error
		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching  proprietrio", 500)
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error marshaling response", 500)
		}

		return responseJSON, nil

	},
}
