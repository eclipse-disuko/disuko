// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package reviewremarks

import (
	"slices"
	"testing"
)

func TestCarryOverStateFrom_AddsSystemEvent(t *testing.T) {
	previous := &Remark{
		SBOMId: "sbom1", Status: Open,
		Events: []Event{{Type: CarriedOverEvent}, {Type: ChangedLevelEvent}},
	}
	current := &Remark{SBOMId: "sbom2"}
	current.CarryOverStateFrom(previous)
	if current.Status != Open || len(current.Events) != 3 {
		t.Fatalf("expected status and inherited events plus carry-over, got %+v", current)
	}
	if current.Events[2].Type != CarriedOverEvent || current.Events[2].Author != "SYSTEM" || current.Events[2].Key == "" {
		t.Errorf("expected a new SYSTEM CARRIED_OVER event with an entity key, got %+v", current.Events[2])
	}
	if len(previous.Events) != 2 {
		t.Fatal("original event history must remain unchanged")
	}
}

func TestCarryOverStateFrom_DoesNotRecordSameSBOMAsCarryOver(t *testing.T) {
	previous := &Remark{SBOMId: "sbom2", Status: Open, Events: []Event{{Type: CommentEvent}}}
	current := &Remark{SBOMId: "sbom2"}
	current.CarryOverStateFrom(previous)
	if len(current.Events) != 1 || current.Events[0].Type != CommentEvent {
		t.Fatalf("same-SBOM state reuse should not create a carry-over event: %+v", current.Events)
	}
}

func TestMergeOpenRemarksWithoutEvents_AllowsOnlyCarryOverEvents(t *testing.T) {
	items := []*Remark{
		{Status: Open, Components: []ComponentMeta{{ComponentName: "A"}}, Events: []Event{{Type: CarriedOverEvent}}},
		{Status: Open, Components: []ComponentMeta{{ComponentName: "B"}}, Events: []Event{{Type: CarriedOverEvent}, {Type: CarriedOverEvent}}},
		{Status: Open, Components: []ComponentMeta{{ComponentName: "C"}}},
		{Status: Open, Components: []ComponentMeta{{ComponentName: "D"}}, Events: []Event{{Type: ChangedLevelEvent}}},
		{Status: InProgress, Components: []ComponentMeta{{ComponentName: "E"}}, Events: []Event{{Type: CarriedOverEvent}}},
	}
	merged := mergeOpenRemarksWithoutEvents(items)
	if len(merged) != 3 {
		t.Fatalf("want merged ABC and separate D, E; got %d remarks", len(merged))
	}
	if len(merged[0].Components) != 3 || len(merged[0].Events) != 3 ||
		!slices.ContainsFunc(merged[0].Components, func(c ComponentMeta) bool { return c.ComponentName == "C" }) {
		t.Errorf("want ABC with three preserved CARRIED_OVER events, got %+v", merged[0])
	}
	if merged[1] != items[3] || merged[2] != items[4] {
		t.Error("business event history or status must prevent merge")
	}
}

func TestSplitAndCarryOverFrom_MergePreservesCarryOverEvents(t *testing.T) {
	component := func(name string) ComponentMeta { return ComponentMeta{ComponentName: name, ComponentVersion: "1"} }
	current := &Remark{
		SBOMId: "sbom2", TemplateKey: "template", Status: Open,
		Components: []ComponentMeta{component("A"), component("B"), component("C")},
	}
	previous := []*Remark{
		{SBOMId: "sbom1", TemplateKey: "template", Status: Open, Components: []ComponentMeta{component("A")}, Events: []Event{{Type: CarriedOverEvent}}},
		{SBOMId: "sbom1", TemplateKey: "template", Status: InProgress, Components: []ComponentMeta{component("B")}},
	}
	result := current.SplitAndCarryOverFrom(previous)
	if len(result) != 2 {
		t.Fatalf("want merged {A,C} and separate {B}, got %d", len(result))
	}
	if result[0].Status != Open || len(result[0].Components) != 2 || len(result[0].Events) != 2 {
		t.Errorf("want {A,C} with historical and new carry-over event, got %+v", result[0])
	}
	if result[1].Status != InProgress || len(result[1].Events) != 1 || result[1].Events[0].Type != CarriedOverEvent {
		t.Errorf("want {B} as IN_PROGRESS with carry-over event, got %+v", result[1])
	}
}
