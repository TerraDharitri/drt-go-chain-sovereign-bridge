package testscommon

// HasherMock -
type HasherMock struct {
	ComputeCalled   func(s string) []byte
	EmptyHashCalled func() []byte
	SizeCalled      func() int
}

// Compute -
func (mock *HasherMock) Compute(s string) []byte {
	if mock.ComputeCalled != nil {
		return mock.ComputeCalled(s)
	}

	return nil
}

// Size -
func (mock *HasherMock) Size() int {
	if mock.SizeCalled != nil {
		return mock.SizeCalled()
	}

	return 0
}

// IsInterfaceNil -
func (mock *HasherMock) IsInterfaceNil() bool {
	return mock == nil
}
