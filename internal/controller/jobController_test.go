package controller

import (
	"strings"
	"testing"

	"github.com/BUIHUUHAIDANG/freelance_market/internal/model"
)

func TestValidateAddJobRequest(t *testing.T) {
	valid := model.AddJobRequest{
		HireManagerID: 1,
		DurationID:    1,
		ComplexityID:  1,
		MainSkillID:   1,
		PaymentTypeID: 1,
		Title:         "Frontend developer",
	}
	if err := validateAddJobRequest(valid); err != nil {
		t.Fatalf("valid job request rejected: %v", err)
	}

	tests := []struct {
		name    string
		input   model.AddJobRequest
		wantErr string
	}{
		{
			name: "invalid reference ID",
			input: func() model.AddJobRequest {
				input := valid
				input.HireManagerID = 0
				return input
			}(),
			wantErr: "reference IDs",
		},
		{
			name: "blank title",
			input: func() model.AddJobRequest {
				input := valid
				input.Title = "  "
				return input
			}(),
			wantErr: "title",
		},
		{
			name: "negative payment",
			input: func() model.AddJobRequest {
				amount := -1.0
				input := valid
				input.PaymentAmount = &amount
				return input
			}(),
			wantErr: "payment_amount",
		},
		{
			name: "invalid required skill ID",
			input: func() model.AddJobRequest {
				input := valid
				input.RequiredSkillIDs = []int64{0}
				return input
			}(),
			wantErr: "required_skill_ids",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateAddJobRequest(test.input)
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("validateAddJobRequest() error = %v, want text %q", err, test.wantErr)
			}
		})
	}
}
