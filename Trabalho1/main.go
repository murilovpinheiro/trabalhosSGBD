package main

import (
	"bufio"
	"fmt"
	"os"
)

// Main que executa seleciona as funções e mostra tudo direitinho

func main() {
	sgbd := NewSGBD()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("SGBD em execução...")
	for {
		fmt.Println("Escolha uma opção:")
		fmt.Println("1 - Inserir")
		fmt.Println("2 - Deletar")
		fmt.Println("3 - Procurar")
		fmt.Println("4 - FULL SCAN")
		fmt.Println("5 - Sair")

		choice, _ := reader.ReadString('\n')
		choice = choice[:len(choice)-1]

		switch choice {
		case "1":
			fmt.Println("Digite o conteúdo do documento a ser inserido:")
			content, _ := reader.ReadString('\n')
			content = content[:len(content)-1]
			if sgbd.Insert([]byte(content)) {
				fmt.Println("Documento inserido com sucesso")
			} else {
				fmt.Println("\nNão foi possível inserir o documento")
			}
		case "2":
			fmt.Println("Digite o conteúdo do documento a ser deletado:")
			content, _ := reader.ReadString('\n')
			content = content[:len(content)-1]
			if len(content) > 5 {
				fmt.Println("Conteúdo do documento deve ter no máximo 5 caracteres")
				continue
			}
			if err := sgbd.Delete([]byte(content)); err == nil {
				fmt.Println("Documento deletado com sucesso")
			} else {
				fmt.Println("Não foi possível deletar o documento:", err)
			}
		case "3":
			fmt.Println("Digite o conteúdo do documento a ser procurado:")
			content, _ := reader.ReadString('\n')
			content = content[:len(content)-1]
			if len(content) > 5 {
				fmt.Println("Conteúdo do documento deve ter no máximo 5 caracteres")
				continue
			}
			if did, err := sgbd.Seek([]byte(content)); err == nil {
				fmt.Println("Documento encontrado:", did)
			} else {
				fmt.Println("Documento não encontrado")
			}
		case "4":
			docs := sgbd.Scan()
			fmt.Println("\nDocumentos armazenados no SGBD:")
			for _, doc := range docs {
				fmt.Printf("DID: %d - Caracteres Armazenados: '%s'\n", doc.did, doc.dados)
			}
		case "5":
			fmt.Println("Saindo...")
			return
		default:
			fmt.Println("Opção inválida")
		}
	}
}
