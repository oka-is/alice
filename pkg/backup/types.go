package backup

const (
	_                            = byte(0)
	MarkerWhoAmI                 = byte(1)
	MarkerListWorkspacesResponse = byte(2)
	MarkerListCardsResponse      = byte(3)
)

//go:generate mockgen -destination ../backup_mock/flush_mock.go -source types.go -package backup_mock -mock_names IFlush=MockFlush
type IFlush interface {
	Flush()
}
