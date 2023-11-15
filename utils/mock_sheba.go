package utils

var bankCode = map[string]string{
	"1000": "ملت",
	"1011": "ملی",
	"1022": "Eghtesad Novin",
	"1033": "Tejart",
	"1044": "Pasargad",
}

func GetIdentityAndBank(sheba string) (identity, bankName string) {
	if len(sheba) != 24 {
		return "", ""
	}
	bankName = bankCode[sheba[:4]]
	identity = sheba[4:14]
	return identity, bankName
}
