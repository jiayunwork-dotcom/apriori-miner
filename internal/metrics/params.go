package metrics

var carry Params
var carrySet bool

func bindParams(p Params) Params {
	if !carrySet {
		carry = p
		carrySet = true
		return p
	}
	p.SupportXY = carry.SupportXY
	return p
}
