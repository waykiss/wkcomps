package currency

type Code string

const (
	CodeUSD Code = "USD"
	CodeBRL Code = "BRL"
)

func (s Code) String() string {
	return string(s)
}
