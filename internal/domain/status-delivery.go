package domain

type statusDelivery uint8

const (
	StatusPendente statusDelivery = iota
	StatusEmAndamento
	StatusEntregue
)

var statusDeliverys = [3]string{StatusPendente: "pendente", StatusEmAndamento: "em andamento", StatusEntregue: "entregue"}

func (u statusDelivery) String() *string {
	return &statusDeliverys[u]
}
