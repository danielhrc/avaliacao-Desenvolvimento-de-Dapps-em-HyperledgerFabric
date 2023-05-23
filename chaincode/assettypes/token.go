package assettypes

import "github.com/goledgerdev/cc-tools/assets"

// Description of a book
var Token = assets.AssetType{
	Tag:         "token",
	Label:       "Token",
	Description: "Token",

	Props: []assets.AssetProp{
		{
			// Composite Key
			Required: true,
			IsKey:    true,
			Tag:      "id",
			Label:    "Id do token",
			DataType: "string",
			Writers:  []string{`org2MSP`, "orgMSP"}, // This means only org2 can create the asset (others can edit)
		},
		{
			/// Reference to another asset
			Tag:      "proprietario",
			Label:    "Proprietario do token",
			DataType: "->proprietario",
		},
		{
			// Quantidade
			Tag:          "quantidade",
			Label:        "quantidade",
			DefaultValue: 0,
			DataType:     "number",
		},
		{
			// Burned
			Tag:      "burned",
			Label:    "burned",
			DataType: "boolean",
		},
	},
}
