package main

import (
	"errors"
	"fmt"
)
//faz tudo, armazenamento de todas as transações, tabela de bloqueio, grafo de espera e waitList dos itens
type TrManager struct {
	Transactions  []Transaction
	LockTable     LockTable
	WaitGraph     Graph
	WaitListItems WaitItemList
}

func (trm *TrManager) NewTransaction(trID, status int) Transaction {

	tr := Transaction{
		TrID:   trID,
		Status: 0,
	}

	trm.Transactions = append(trm.Transactions, tr)

	// sempre que criar uma transação colocar ela na lista de transações
	return tr
}

// ReadLock insere um bloqueio de leitura na tabela de bloqueios
func (trm *TrManager) ReadLock(TrID int, ItemID string, duration string) (int, error) {
	for _, transaction := range trm.Transactions {
		if transaction.TrID != TrID && transaction.Status == 0 {
			for _, lock := range trm.LockTable.Locks {
				if lock.OpType == "W" && lock.TrID != TrID && lock.ItemID == ItemID {
					trm.setWait(TrID, lock)
					return lock.TrID, errors.New("Bloqueio de escrita encontrado no mesmo objeto")
					// Se houver um bloqueio de escrita sobre o mesmo objeto, não liberamos
				}
			}

			fmt.Printf("      Transação %d - Obtém o bloqueio de Leitura sobre o Item %s\n", TrID, ItemID)

			// Criar novo objeto Lock
			lock, err := NewLock(ItemID, TrID, "O", duration, "R")
			if err != nil {
				return 0, err
			}

			// Adicionar lock à LockTable
			trm.LockTable.Locks = append(trm.LockTable.Locks, *lock)

			// Realizar outras operações conforme necessário
			// ... IF caso seja curto tiro o lock
			if duration == "C" {
				trm.ReleaseLock(TrID, ItemID)
			}
			return len(trm.LockTable.Locks), nil
		}
	}

	return 0, nil
}

func (trm *TrManager) WriteLock(TrID int, ItemID string, duration string) (int, error) {
	for _, transaction := range trm.Transactions {
		if transaction.TrID != TrID && transaction.Status == 0 {
			for _, lock := range trm.LockTable.Locks {
				if lock.TrID != TrID && lock.ItemID == ItemID {
					trm.setWait(TrID, lock)
					return lock.TrID, errors.New("Bloqueio encontrado no mesmo objeto.")
				}
			}

			fmt.Printf("      Transação %d - Obtém bloqueio de Escrita sobre o Item %s\n", TrID, ItemID)

			// Criar novo objeto Lock
			lock, err := NewLock(ItemID, TrID, "O", duration, "W")
			if err != nil {
				return 0, err
			}

			// Adicionar lock à LockTable
			trm.LockTable.Locks = append(trm.LockTable.Locks, *lock)

			// Realizar outras operações conforme necessário
			// ... IF caso seja curto tiro o lock
			if duration == "C" {
				trm.ReleaseLock(TrID, ItemID)
			}

			return len(trm.LockTable.Locks), nil
		}
	}

	return 0, nil
}

func (trm *TrManager) ReleaseLock(TrID int, ItemID string) {
	for i_lock, lock := range trm.LockTable.Locks {
		if ItemID != "" {
			if lock.ItemID == ItemID && lock.TrID == TrID {
				copy(trm.LockTable.Locks[i_lock:], trm.LockTable.Locks[i_lock+1:])     // Substitui elemento na posição i_lock pelo elemento subsequente
				trm.LockTable.Locks = trm.LockTable.Locks[:len(trm.LockTable.Locks)-1] // Reduz o tamanho da fatia em 1

				trm.escalonarWaitFor(lock.ItemID)

				fmt.Printf("      Transação %d - Libera bloqueio do tipo %s sobre o item %s\n", TrID, lock.OpType, ItemID)
				break
			}
		} else {
			if lock.TrID == TrID {
				if len(trm.LockTable.Locks) <= 1 {
					trm.LockTable.Locks = []Lock{}
				} else {
					copy(trm.LockTable.Locks[i_lock:], trm.LockTable.Locks[i_lock+1:])
					trm.LockTable.Locks = trm.LockTable.Locks[:len(trm.LockTable.Locks)-1]
				}
				fmt.Printf("      Transação %d - Libera bloqueio do tipo %s sobre o item %s\n", TrID, lock.OpType, ItemID)

			}
			trm.escalonarWaitFor(lock.ItemID)
		}
	}
}

func (trm *TrManager) Commit(TrID int) {

	for _, transaction := range trm.Transactions {
		if transaction.TrID == TrID {
			transaction.Status = 1
		}
	}
	trm.ReleaseLock(TrID, "")
}

func (trm *TrManager) setWait(TrWithLock int, requester Lock) Edge {

	e := Edge{-1, -1}

	if requester.TrID > TrWithLock {

		for _, transaction := range trm.Transactions {

			if transaction.TrID == requester.TrID {
				transaction.Status = 2
				// abortou a transação caso ela seja mais velha? ver esse depois
			}
		}

		fmt.Printf("      Transação %d - Abortada por estratégia Wait-Die (Transação %d possui bloqueio em item %s)\n", requester.TrID, TrWithLock, requester.ItemID)
		trm.ReleaseLock(requester.TrID, "")
		return e
	}

	for _, edge := range trm.WaitGraph.Edges {
		if edge.PriorityTransaction == requester.TrID && edge.WaitingTransaction == TrWithLock {
			// Existe uma relação de espera entre edge.PriorityTransaction -> WaitingTransaction
			// Só que o Waiting Transaction-> PriorityTransaction tá sendo solicitado
			// Duvida se relações mais "complexas" são encontradas aqui tipo
			// T1->T2->T3->T1, acho que não MAAAAAS não sei como resolver então deixa assim
			// DEADLOCK
			fmt.Printf("     Transação %d - Deadlock com Transação %d\n", edge.PriorityTransaction, edge.WaitingTransaction)
			return edge
		}
	}

	new_edge := Edge{TrWithLock, requester.TrID}

	trm.WaitGraph.Edges = append(trm.WaitGraph.Edges, new_edge)

	fmt.Printf("%d - Entra na Fila de Espera pela Liberação do Item %s pela Transação %d\n", requester.TrID, requester.ItemID, TrWithLock)

	for _, transaction := range trm.Transactions {
		if transaction.TrID == requester.TrID {
			transaction.Status = 3
		}
	}
	for _, wf_item := range trm.WaitListItems.list {
		if wf_item.ItemID == requester.ItemID {
			wf_item.WaitList = append(wf_item.WaitList, requester)
			return e
		}
	}

	var lt []Lock
	lt = append(lt, requester)

	wf_item := WaitForItem{
		ItemID:   requester.ItemID,
		WaitList: lt,
	}

	trm.WaitListItems.list = append(trm.WaitListItems.list, wf_item)

	return e

}

func (trm *TrManager) escalonarWaitFor(idItem string) {

	for _, waitForItem := range trm.WaitListItems.list {

		if waitForItem.ItemID == idItem {
			if len(waitForItem.WaitList) < 1 { // Se o item não tiver ninguém esperando a gente para aqui
				return
			}

			firstLock := waitForItem.WaitList[0]
			waitForItem.WaitList = waitForItem.WaitList[1:]

			for _, transaction := range trm.Transactions {
				if transaction.TrID == firstLock.TrID {
					transaction.Status = 0
				}
			}

			if firstLock.OpType == "W" {
				// fmt.Println(fmt.Sprintf("Transação %d - Solicita bloqueio de Escrita sobre o item %s", trID, idItem))TrID int, ItemID string, duration string
				result, erro := trm.WriteLock(firstLock.TrID, firstLock.ItemID, "L")
				if erro != nil {
					fmt.Println(erro)
				}

				if result != -1 {
					trm.setWait(result, firstLock)
				}

			} else {
				// fmt.Println(fmt.Sprintf("Transação %d - Solicita bloqueio de Escrita sobre o item %s", trID, idItem))
				result, erro := trm.WriteLock(firstLock.TrID, firstLock.ItemID, "L")

				if erro != nil {
					fmt.Println(erro)
				}

				if result != -1 {
					trm.setWait(result, firstLock)
				}
			}
		}
	}
}
