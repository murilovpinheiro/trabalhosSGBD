package main


type Transaction struct {
	TrID   int
	Status int // 0 ativa; 1 concluída; 2 abortada; 3 esperando.
}

//aresta do grafo de espera
type Edge struct {
	PriorityTransaction int
	WaitingTransaction  int
}

type Graph struct {
	Edges []Edge
}

type WaitItemList struct {
	list []WaitForItem
}

type WaitForItem struct {
	ItemID   string
	WaitList []Lock
}
