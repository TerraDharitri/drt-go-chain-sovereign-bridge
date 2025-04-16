package txSender

import "errors"

var errInvalidWalletType = errors.New("invalid/unknown wallet type")

var errNilWallet = errors.New("nil wallet provided")

var errNilProxy = errors.New("nil proxy provided")

var errNilTxInteractor = errors.New("nil tx interactor provided")

var errNilDataFormatter = errors.New("nil data formatter provided")

var errNilNonceHandler = errors.New("nil nonce handler provided")

var errNoHeaderVerifierSCAddress = errors.New("no header verifier sc address provided")

var errNoDcdtSafeSCAddress = errors.New("no dcdt safe sc address provided")
