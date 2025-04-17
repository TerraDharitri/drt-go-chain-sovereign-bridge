package txSender

// WalletConfig holds wallet config
type WalletConfig struct {
	Path     string
	Password string
}

// TxSenderConfig holds tx sender config
type TxSenderConfig struct {
	HeaderVerifierSCAddress string
	DcdtSafeSCAddress       string
	Proxy                   string
	IntervalToSend          int
	Hasher                  string
}
