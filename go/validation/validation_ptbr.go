package validation

var messagesPtBR = map[string]string{
	isByteLength:  "%s (%s) precisa ter entre %d caracteres e %d, ele possui %d",
	inRangeFloat:  "%s precisa ser entre entre %f e %f, valor atual %f",
	inRangeInt:    "%s precisa ser entre entre %d e %d, valor atual %d",
	equalsFloat:   "%s precisa ser igual a %f, valor atual é %f",
	isRequired:    "%s é obrigatório",
	isObjectId:    "%s nao é um ObjectID válido",
	isEmailValid:  "Campo %s \"%s\" nao é um email válido",
	isGreaterThan: "Campo %s precisa ser maior que \"%s\"",
	isOnlyNumber:  "Campo %s só pode conter números ",
	isNotValid:    "Campo %s nao é válido com o valor \"%s\" ",
	isNotUrlValid: "Campo %s nao é uma URL válida com o valor \"%s\" ",
}
