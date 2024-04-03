package rules

import (
	"testing"
)

func Test_ValidSha256(t *testing.T) {
	type test struct {
		name       string
		identifier string
		isValid    bool
	}
	tests := []test{
		{
			name:       "4cd1fce53a8b3e67e174859e6672ca29bc1e16585859c53a116e7f53d04350b7",
			identifier: "4cd1fce53a8b3e67e174859e6672ca29bc1e16585859c53a116e7f53d04350b7",
			isValid:    true,
		},
		{
			name:       "1507564a650077bdc8e155b2a4ba8bd85a55dc347a34ddf4e78a836c17d81bb1",
			identifier: "1507564a650077bdc8e155b2a4ba8bd85a55dc347a34ddf4e78a836c17d81bb1",
			isValid:    true,
		},
		{
			name:       "1507564a650077bdc8e155b2a4ba8bd85a55dc347a34ddf4e78a836c17d81bb",
			identifier: "1507564a650077bdc8e155b2a4ba8bd85a55dc347a34ddf4e78a836c17d81bb",
			isValid:    false,
		},
		{
			name:       "588d84953ae992c5de61d3774ce86e710ed42d29",
			identifier: "588d84953ae992c5de61d3774ce86e710ed42d29",
			isValid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidSha256(tt.identifier)
			if got != tt.isValid {
				t.Errorf("ValidSha256() got = %v, want %v", got, tt.isValid)
				return
			}
		})
	}
}

func Test_ValidTeamID(t *testing.T) {
	type test struct {
		name       string
		identifier string
		isValid    bool
	}
	tests := []test{
		{
			name:       "EQHXZ8M8AV",
			identifier: "EQHXZ8M8AV",
			isValid:    true,
		},
		{
			name:       "APPLE",
			identifier: "APPLE",
			isValid:    true,
		},
		{
			name:       "EQHXZ8M8AVAAAAA",
			identifier: "EQHXZ8M8AVAAAAA",
			isValid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidTeamID(tt.identifier)
			if got != tt.isValid {
				t.Errorf("ValidTeamID() got = %v, want %v", got, tt.isValid)
				return
			}
		})
	}
}

func Test_ValidSigningID(t *testing.T) {
	type test struct {
		name       string
		identifier string
		isValid    bool
	}
	tests := []test{
		{
			name:       "EQHXZ8M8AV:com.google.Chrome",
			identifier: "EQHXZ8M8AV:com.google.Chrome",
			isValid:    true,
		},
		{
			name:       "EQHXZ8M8AVAAAAA:com.google.Chrome",
			identifier: "EQHXZ8M8AVAAAAA:com.google.Chrome",
			isValid:    false,
		},
		{
			name:       "com.google.Chrome",
			identifier: "com.google.Chrome",
			isValid:    false,
		},
		{
			name:       "platform:com.apple.curl",
			identifier: "platform:com.apple.curl",
			isValid:    true,
		},
		{
			name:       ":com.apple.curl",
			identifier: ":com.apple.curl",
			isValid:    false,
		},
		{
			name:       "APPLE:com.apple.curl",
			identifier: "APPLE:com.apple.curl",
			isValid:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidSigningID(tt.identifier)
			if got != tt.isValid {
				t.Errorf("ValidSigningID() got = %v, want %v", got, tt.isValid)
				return
			}
		})
	}
}
