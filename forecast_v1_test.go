package changeevidence

import (
	"errors"
	"math"
	"reflect"
	"strings"
	"testing"
)

func forecastIngestInputV1() ForecastInputV1 {
	zero, maximum := uint64(0), uint64(math.MaxUint64)
	return ForecastInputV1{
		AccountingPolicyDigest: DocumentDigest(strings.Repeat("a", 64)),
		InventoryDigest:        DocumentDigest(strings.Repeat("b", 64)),
		Revisions:              RevisionIdentity{Base: strings.Repeat("1", 40), Head: strings.Repeat("2", 40)},
		Metrics: []ForecastMetricV1{
			{Category: "z", Kind: ForecastNonCountableV1, Value: nil},
			{Category: "a", Kind: ForecastDeletionsV1, Value: &maximum},
			{Category: "z", Kind: ForecastAdditionsV1, Value: &zero},
		},
	}
}

func TestForecastIngestRetentionAndOwnership(t *testing.T) {
	input := forecastIngestInputV1()
	want := forecastIngestInputV1()
	forecast, err := IngestForecastV1(input)
	if err != nil {
		t.Fatal(err)
	}
	if forecast.AccountingPolicyDigest() != want.AccountingPolicyDigest || forecast.InventoryDigest() != want.InventoryDigest || forecast.Revisions() != want.Revisions {
		t.Fatal("declared binding changed")
	}
	if !reflect.DeepEqual(forecast.Metrics(), want.Metrics) {
		t.Fatal("values/order/unknown changed")
	}
	if string(ForecastAdditionsV1) != "additions" || string(ForecastDeletionsV1) != "deletions" || string(ForecastNonCountableV1) != "non_countable" {
		t.Fatal("axis names changed")
	}
	*input.Metrics[1].Value = 7
	*input.Metrics[2].Value = 8
	input.Metrics[0].Category = "changed"
	input.Metrics[1].Kind = ForecastAdditionsV1
	input.Revisions.Base = "changed"
	input.InventoryDigest = "changed"
	if !reflect.DeepEqual(forecast.Metrics(), want.Metrics) || forecast.Revisions() != want.Revisions || forecast.InventoryDigest() != want.InventoryDigest {
		t.Fatal("ingress aliases caller")
	}
	got := forecast.Metrics()
	*got[1].Value = 9
	*got[2].Value = 10
	got[0] = ForecastMetricV1{Category: "changed"}
	if !reflect.DeepEqual(forecast.Metrics(), want.Metrics) {
		t.Fatal("egress aliases retained input")
	}
}

func TestForecastIngestOmissionAndRevisionFormats(t *testing.T) {
	for _, size := range []int{40, 64} {
		for _, metrics := range [][]ForecastMetricV1{nil, {}} {
			input := forecastIngestInputV1()
			input.Revisions = RevisionIdentity{Base: strings.Repeat("a", size), Head: strings.Repeat("a", size)}
			input.Metrics = metrics
			forecast, err := IngestForecastV1(input)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(forecast.Metrics(), metrics) || forecast.Revisions() != input.Revisions {
				t.Fatal("omission/empty or valid identity changed")
			}
		}
	}
}

func TestForecastIngestInvalidContracts(t *testing.T) {
	cases := []struct {
		name, field, code string
		mutate            func(*ForecastInputV1)
	}{
		{"policy", "accounting_policy_sha256", "invalid_digest", func(i *ForecastInputV1) { i.AccountingPolicyDigest = "bad" }},
		{"inventory", "inventory_sha256", "invalid_digest", func(i *ForecastInputV1) { i.InventoryDigest = DocumentDigest(strings.Repeat("A", 64)) }},
		{"base", "revisions.base", "invalid_revision", func(i *ForecastInputV1) { i.Revisions.Base = "bad" }},
		{"head", "revisions.head", "invalid_revision", func(i *ForecastInputV1) { i.Revisions.Head = strings.Repeat("A", 40) }},
		{"mixed formats", "revisions.head", "invalid_revision", func(i *ForecastInputV1) { i.Revisions.Head = strings.Repeat("2", 64) }},
		{"empty category", "metrics.category", "invalid", func(i *ForecastInputV1) { i.Metrics[0].Category = "" }},
		{"control category", "metrics.category", "invalid", func(i *ForecastInputV1) { i.Metrics[0].Category = "private\nvalue" }},
		{"invalid UTF8", "metrics.category", "invalid", func(i *ForecastInputV1) { i.Metrics[0].Category = string([]byte{255}) }},
		{"unsupported axis", "metrics.kind", "invalid", func(i *ForecastInputV1) { i.Metrics[0].Kind = "actual" }},
		{"empty axis", "metrics.kind", "invalid", func(i *ForecastInputV1) { i.Metrics[0].Kind = "" }},
		{"duplicate", "metrics", "duplicate", func(i *ForecastInputV1) { i.Metrics = append(i.Metrics, i.Metrics[0]) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := forecastIngestInputV1()
			tc.mutate(&input)
			forecast, err := IngestForecastV1(input)
			var contract *ContractError
			if !errors.As(err, &contract) || contract.Field != tc.field || contract.Code != tc.code {
				t.Fatalf("expected %s/%s, got %v", tc.field, tc.code, err)
			}
			if !reflect.DeepEqual(forecast, ForecastV1{}) {
				t.Fatal("invalid input retained")
			}
		})
	}
}
