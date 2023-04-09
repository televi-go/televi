package pages

import "testing"

type SomeBigPage struct {
	AState  State[int]
	BState  State[int]
	CState  State[string]
	DState  State[string]
	EState  State[int]
	FState  State[int]
	GState  State[string]
	HState  State[string]
	IState  State[float64]
	JState  State[float64]
	KState  State[bool]
	LState  State[bool]
	MState  State[float64]
	NState  State[float64]
	OState  State[string]
	PState  State[string]
	QState  State[int]
	RState  State[int]
	SState  State[float64]
	TState  State[float64]
	UState  State[bool]
	VState  State[bool]
	WState  State[float64]
	XState  State[float64]
	YState  State[string]
	ZState  State[string]
	AAState State[int]
	BBState State[int]
	CCState State[float64]
	DDState State[float64]
	EEState State[bool]
	FFState State[bool]
	GGState State[float64]
	HHState State[float64]
	IIState State[string]
	JJState State[string]
	KKState State[int]
	LLState State[int]
	MMState State[float64]
	NNState State[float64]
	OOState State[bool]
	PPState State[bool]
	QQState State[float64]
	RRState State[float64]
	SSState State[string]
	TTState State[string]
	UUState State[int]
	VVState State[int]
	WWState State[float64]
	XXState State[float64]
	YYState State[bool]
	ZZState State[bool]
}

func (s SomeBigPage) View(ctx PageBuildContext) {
	//TODO implement me
	panic("implement me")
}

type SomeUsualPage struct {
	State      State[string]
	OtherState State[string]
	Property   string
}

func (s SomeUsualPage) View(ctx PageBuildContext) {
	//TODO implement me
	panic("implement me")
}

func BenchmarkMountStates(b *testing.B) {
	var p Scene = SomeBigPage{}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		MountStates(&p, nil)
	}
}

func BenchmarkMountStates2(b *testing.B) {
	var p Scene = SomeUsualPage{}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		MountStates(&p, nil)
	}
}

func BenchmarkMountStatesWithInit(b *testing.B) {
	var p Scene = SomeUsualPage{
		State:      StateOf(""),
		OtherState: StateOf("123"),
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		MountStates(&p, nil)
	}
}
