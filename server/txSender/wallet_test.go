package txSender

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadWallet(t *testing.T) {
	type testScenario struct {
		cfg             WalletConfig
		expectedError   error
		expectedAddress string
	}

	scenarios := []testScenario{
		{
			cfg: WalletConfig{
				Path:     "testData/alice.pem",
				Password: "",
			},
			expectedError:   nil,
			expectedAddress: "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf",
		},
		{
			cfg: WalletConfig{
				Path:     "testData/bob.json",
				Password: "password",
			},
			expectedError:   nil,
			expectedAddress: "drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c",
		},
		{
			cfg: WalletConfig{
				Path:     "testData/alice.ledger",
				Password: "",
			},
			expectedError:   errInvalidWalletType,
			expectedAddress: "",
		},
	}

	for _, scenario := range scenarios {
		log.Info("executing test scenario", "wallet", scenario.cfg.Path)

		wallet, err := LoadWallet(scenario.cfg)
		if scenario.expectedError == nil {
			require.Nil(t, err)
			require.Equal(t, scenario.expectedAddress, wallet.GetBech32())
		} else {
			require.ErrorIs(t, err, scenario.expectedError)
			require.Nil(t, wallet)
		}
	}
}
