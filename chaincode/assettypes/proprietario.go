package assettypes

import (
	"fmt"

	"github.com/goledgerdev/cc-tools/assets"
)

var Proprietario = assets.AssetType{
	Tag:         "proprietario",
	Label:       "Propietario",
	Description: "Proprietario",

	Props: []assets.AssetProp{
		{
			// Primary key
			Required: true,
			IsKey:    true,
			Tag:      "id",
			Label:    "id do proprietario",
			DataType: "string",                      // Datatypes are identified at datatypes folder
			Writers:  []string{`org1MSP`, "orgMSP"}, // This means only org1 can create the asset (others can edit)
		},
		{
			// Mandatory property
			Required: true,
			Tag:      "nome",
			Label:    "Nome do proprietario",
			DataType: "string",
			// Validate funcion
			Validate: func(name interface{}) error {
				nameStr := name.(string)
				if nameStr == "" {
					return fmt.Errorf("name must be non-empty")
				}
				return nil
			},
		},
	},
}
