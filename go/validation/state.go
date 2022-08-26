package validation

type States struct {
	Name string
	Code string
}

var statesBr = map[string]States{
	"RO": {"Rondônia", "RO"},
	"AC": {"Acre", "AC"},
	"AM": {"Amazonas", "AM"},
	"RR": {"Roraima", "PR"},
	"PA": {"Pará", "PA"},
	"AP": {"Amapá", "AP"},
	"TO": {"Tocantins", "TO"},
	"MA": {"Maranhão", "MA"},
	"PI": {"Piauí", "PI"},
	"CE": {"Ceará", "CE"},
	"RN": {"Rio Grande do Norte", "RN"},
	"PB": {"Paraíba", "PB"},
	"PE": {"Pernambuco", "PE"},
	"AL": {"Alagoas", "AL"},
	"SE": {"Sergipe", "SE"},
	"BA": {"Bahia", "BA"},
	"MG": {"Minas Gerais", "MG"},
	"ES": {"Espírito Santo", "ES"},
	"RJ": {"Rio de Janeiro", "RJ"},
	"SP": {"São Paulo", "SP"},
	"PR": {"Paraná", "PR"},
	"SC": {"Santa Catarina", "SC"},
	"RS": {"Rio Grande do Sul", "RS"},
	"MS": {"Mato Grosso do Sul", "MS"},
	"MT": {"Mato Grosso", "MT"},
	"GO": {"Goiás", "GO"},
	"DF": {"Distrito Federal", "DF"},
}

//IsStateBR check if the state is valid for Brazil
func IsStateBR(state string) bool {
	_, ok := statesBr[state]
	return ok
}
