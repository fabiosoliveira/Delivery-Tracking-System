package domain

type StatusDelivery uint8

const (
	StatusPendente StatusDelivery = iota
	StatusEmAndamento
	StatusEntregue
)

var statusDeliverys = [3]string{StatusPendente: "pendente", StatusEmAndamento: "em andamento", StatusEntregue: "entregue"}

func (u StatusDelivery) String() *string {
	return &statusDeliverys[u]
}
