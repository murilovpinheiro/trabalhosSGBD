package main

import (
	"fmt"
	"strconv"
	"strings"
)

func GetIsolationLevel(level int) (string, string) {
	switch level {
	case 0:
		return "C", "C"
	case 1:
		return "C", "L"
	case 2:
		return "L", "L"
	case 3:
		return "L", "L"
	default:
		return "", ""
	}
}
func main() {
	I_level := 0
	schedule := ""

	fmt.Println("Sistema de Gerenciamento de Transações")
	fmt.Println("Qual o Nível de Isolamento Desejado:\n0:Read Uncommitted\n1:Read Committed\n2:Repeatable Read\n3:Serializable")
	fmt.Scanln(&I_level)
	durationRead, durationWrite := GetIsolationLevel(I_level)

	if durationRead == "" && durationWrite == "" {
		fmt.Println("Nível de isolamento inadequado")
		return
	}

	fmt.Println("Durações das Transações do Nível de Isolamento Escolhido: Leitura:", durationRead, "Escrita:", durationWrite)

	fmt.Println("Digite o Schedule de Transações que deseja rodar:")
	fmt.Scanln(&schedule)

	schedule = strings.ToUpper(schedule)
	trManager := TrManager{
		Transactions: []Transaction{},
		LockTable: LockTable{
			Locks:          []Lock{},
			IsolationLevel: 0,
		},
		WaitGraph:     Graph{Edges: []Edge{}},
		WaitListItems: WaitItemList{list: []WaitForItem{}},
	}

	operations := strings.Split(schedule, ")")
	operations = operations[:(len(operations) - 1)]

	for _, operation := range operations {

		if string(operation[0]) == "B" {
			trID, _ := strconv.Atoi(string(operation[len(operation)-1]))

			fmt.Println("Criando a Transação Número", trID)
			trManager.NewTransaction(trID, 0)

			fmt.Println()

		} else if string(operation[0]) == "W" {
			trID, _ := strconv.Atoi(string(operation[1]))
			idItem := string(operation[len(operation)-1])

			for _, transaction := range trManager.Transactions {

				if transaction.TrID == trID && transaction.Status != 2 {

					fmt.Printf("Transação %d - Solicita bloqueio de Escrita sobre o item %s\n", trID, idItem)
					lockID, err := trManager.WriteLock(trID, idItem, durationWrite)

					if err != nil {
						fmt.Println("Erro ao obter o bloqueio de escrita:", err)
					} else {
						fmt.Printf("Obteve o bloqueio de escrita com ID: %d\n", lockID)
					}

					fmt.Println()
				}
			}

		} else if string(operation[0]) == "R" {
			trID, _ := strconv.Atoi(string(operation[1]))
			idItem := string(operation[len(operation)-1])

			for _, transaction := range trManager.Transactions {

				if transaction.TrID == trID && transaction.Status != 2 {

					fmt.Printf("Transação %d - Solicita bloqueio de Leitura sobre o item %s\n", trID, idItem)
					lockID, err := trManager.ReadLock(trID, idItem, durationRead)

					if err != nil {
						fmt.Println("Erro ao obter o bloqueio de leitura:", err)
					} else {
						fmt.Printf("Obteve o bloqueio de leitura com ID: %d\n", lockID)
					}

					fmt.Println()
				}
			}

		} else if string(operation[0]) == "C" {
			trID, _ := strconv.Atoi(string(operation[len(operation)-1]))

			for _, transaction := range trManager.Transactions {

				if transaction.TrID == trID && transaction.Status != 2 {

					fmt.Printf("Transação %d - Solicita Commit.\n", trID)

					trManager.Commit(trID)

					fmt.Println("Commit efetuado.")

					fmt.Println()
				}
			}

		}
	}
	fmt.Println(trManager)
}
