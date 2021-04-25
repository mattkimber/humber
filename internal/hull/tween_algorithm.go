package hull

import (
	"encoding/json"
	"fmt"
)

type TweenAlgorithm int

const (
	TweenAlgorithmLinear = iota
	TweenAlgorithmSquareRoot
	TweenAlgorithmReverseSquareRoot
	TweenAlgorithmSquare
)

func (ta TweenAlgorithm) String() string {
	switch ta {
	case TweenAlgorithmLinear:
		return "linear"
	case TweenAlgorithmSquareRoot:
		return "square_root"
	case TweenAlgorithmSquare:
		return "square"
	case TweenAlgorithmReverseSquareRoot:
		return "reverse_square_root"
	default:
		return "unknown"
	}
}

func GetTweenAlgorithmFromName(input string) TweenAlgorithm {
	switch input {
	case "linear":
		return TweenAlgorithmLinear
	case "square_root":
		return TweenAlgorithmSquareRoot
	case "reverse_square_root":
		return TweenAlgorithmReverseSquareRoot
	case "square":
		return TweenAlgorithmSquare
	default:
		return TweenAlgorithmLinear
	}
}

func (ta TweenAlgorithm) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", ta.String())), nil
}

func (ta *TweenAlgorithm) UnmarshalJSON(b []byte) error {
	var algName string
	if err := json.Unmarshal(b, &algName); err != nil {
		return err
	}

	*ta = GetTweenAlgorithmFromName(algName)
	return nil
}