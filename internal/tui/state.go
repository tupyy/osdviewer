package tui

type state int

const (
	LoadingState state = iota
	ReadyState
	ErrorState

	// state for table view
	TableState
	ClusterState
)

type ViewState struct {
	State state
	Err   error
}
