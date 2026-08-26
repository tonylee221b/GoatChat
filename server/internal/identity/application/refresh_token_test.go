package application

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateRefreshToken(t *testing.T) {
	for _, tt := range generateRefreshTokenTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			token, err := generateRefreshToken()
			if err != nil {
				require.ErrorContains(t, err, tt.wantErrMsg)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantPlainLength, len(token.Plain))
			require.Equal(t, tt.wantDecodedLength, len(token.Hash.Bytes()))
		})
	}
}

type generateRefreshTokenTestCase struct {
	name              string
	wantPlainLength   int
	wantDecodedLength int
	wantErrMsg        string
}

func generateRefreshTokenTestCases() []generateRefreshTokenTestCase {
	return []generateRefreshTokenTestCase{
		{
			name:              "success",
			wantPlainLength:   43,
			wantDecodedLength: 32,
			wantErrMsg:        "",
		},
	}
}
