package acm

type AcmFieldValueListener struct {
	*BaseAcmListener
	fields []string
	values []string
}

func NewAcmFieldListener() *AcmFieldValueListener {
	return &AcmFieldValueListener{}
}

func (l *AcmFieldValueListener) GetFields() []string {
	return l.fields
}

func (l *AcmFieldValueListener) GetValues() []string {
	return l.values
}

func (l *AcmFieldValueListener) EnterFieldName(ctx *FieldNameContext) {
	str := ctx.GetText()
	for _, field := range l.fields {
		if field == str {
			return
		}
	}
	l.fields = append(l.fields, str)
}

func (l *AcmFieldValueListener) EnterValue(ctx *ValueContext) {
	str := ctx.GetText()
	for _, value := range l.values {
		if value == str {
			return
		}
	}
	l.values = append(l.values, str)
}
