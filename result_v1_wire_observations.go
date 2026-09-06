package changeevidence

import "encoding/json"

type resultWireAvailableLineThreshold struct {
	Kind      ResultObservationKind `json:"kind"`
	Category  string                `json:"category"`
	Limit     uint64                `json:"limit"`
	Available bool                  `json:"available"`
	Actual    uint64                `json:"actual"`
	Exceeded  bool                  `json:"exceeded"`
}

type resultWireUnavailableLineThreshold struct {
	Kind      ResultObservationKind `json:"kind"`
	Category  string                `json:"category"`
	Limit     uint64                `json:"limit"`
	Available bool                  `json:"available"`
	Exceeded  bool                  `json:"exceeded"`
}

type resultWireAvailableRatio struct {
	Kind        ResultObservationKind `json:"kind"`
	Category    string                `json:"category"`
	Reference   string                `json:"reference"`
	Limit       string                `json:"limit"`
	Available   bool                  `json:"available"`
	Numerator   uint64                `json:"numerator"`
	Denominator uint64                `json:"denominator"`
	Exceeded    bool                  `json:"exceeded"`
}

type resultWireUnavailableRatio struct {
	Kind      ResultObservationKind `json:"kind"`
	Category  string                `json:"category"`
	Reference string                `json:"reference"`
	Limit     string                `json:"limit"`
	Available bool                  `json:"available"`
	Exceeded  bool                  `json:"exceeded"`
}

func resultWireObservation(observation ResultObservationV1) (json.RawMessage, error) {
	switch observation.Kind {
	case ResultLineThresholdObservationV1:
		if observation.Available {
			return json.Marshal(resultWireAvailableLineThreshold{
				observation.Kind, observation.Category, observation.LineThresholdLimit, true, observation.Actual, observation.Exceeded,
			})
		}
		return json.Marshal(resultWireUnavailableLineThreshold{
			observation.Kind, observation.Category, observation.LineThresholdLimit, false, false,
		})
	case ResultRatioObservationV1:
		if observation.Available {
			return json.Marshal(resultWireAvailableRatio{
				observation.Kind, observation.Category, observation.Reference, observation.RatioLimit, true, observation.Numerator, observation.Denominator, observation.Exceeded,
			})
		}
		return json.Marshal(resultWireUnavailableRatio{
			observation.Kind, observation.Category, observation.Reference, observation.RatioLimit, false, false,
		})
	default:
		return nil, &ContractError{"observations.kind", "invalid"}
	}
}
