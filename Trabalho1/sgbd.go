package main

import (
	"bytes"
	"errors"
	"fmt"
)

type SGBD struct {
	PrimeiraPagina *Pagina
	UltimaPagina   *Pagina
	NumPaginas     int
}

//construtor do SGBD
func NewSGBD() *SGBD {
	sgbd := &SGBD{}
	pagina := NewPagina(0) // Pagina 0
	sgbd.PrimeiraPagina = pagina
	sgbd.UltimaPagina = pagina
	sgbd.NumPaginas = 1
	return sgbd
}

func (s *SGBD) updatePageId(pagina *Pagina) {
	// Inicializa o contador de page_id's
	if pagina == nil {
		return
	}
	pageID := pagina.page_id
	pageID--
	// Percorre todas as páginas a partir da página atual
	for p := pagina; p != nil; p = p.next {
		// Atualiza o page_id da página
		p.page_id = pageID
		p.AtualizarPageID(pageID)

		// Incrementa o contador de page_id's
		pageID++
	}
}


// Insert de um documento(data) no banco de dados
func (s *SGBD) Insert(data []byte) bool {
	doc, err := NewDocumento(data)
  // Utiliza o construtor do documento, caso dê erro ele mostra o erro e retorna falso
	if err != nil {
		fmt.Println(err)
		return false
	}
  // Caso o documento tenha sido criado então ele vai iterar nas páginas até encontrar um lugar livre
	paginaAtual := s.PrimeiraPagina
	for {
		if paginaAtual.arm_livre >= doc.did.tam {
			paginaAtual.AddDocumento(doc)
			return true
		}
		if paginaAtual.next == nil {
			break
		}
		paginaAtual = paginaAtual.next
	}

  // Como não encontrou um lugar livre, tenta criar uma nova página para inserir o documento nela
	if s.NumPaginas < 20 {
		pagina := NewPagina(s.NumPaginas)
		pagina.AddDocumento(doc)
		s.UltimaPagina.next = pagina
		s.UltimaPagina = pagina
		s.NumPaginas++
		return true
	}
  // Caso tenha 20 página não podemos fazer isso, para por aqui
  fmt.Println("O SGBD está cheio, não foi possível fazer a inserção")
	return false
}

func (s *SGBD) Scan() []*Documento {
  // O Scan não tem muito segredo, basicamente itera e retorna todos os documentos em um array
	var docs []*Documento
	paginaAtual := s.PrimeiraPagina
	for {
		for _, doc := range paginaAtual.Docs {
			docs = append(docs, doc)
		}
		if paginaAtual.next == nil {
			break
		}
		paginaAtual = paginaAtual.next
	}
	return docs
}

func (s *SGBD) Seek(content []byte) (DID, error) {
  
  // Seek itera sobe as páginas comparando o content com o dado de um documento
  // Caso seja igual, ele captura o DID do documento e retorna o DID
	paginaAtual := s.PrimeiraPagina
	for {
		for _, doc := range paginaAtual.Docs {
			if bytes.Equal(doc.dados, content) {
				return doc.did, nil
			}
		}
		if paginaAtual.next == nil {
			break
		}
		paginaAtual = paginaAtual.next
	}
	return DID{}, errors.New("Documento não encontrado")
}

func (s *SGBD) Delete(content []byte) error {
	paginaAtual := s.PrimeiraPagina
	var paginaAnterior *Pagina
  // Delete é mais complicado
	for {
    // Ele itera sobre todo o SGBD e itera sobre os docs de cada página
		for i, doc := range paginaAtual.Docs {
      // Vai comparando dado por dado, caso seja igual ele tira o documento da página atual e faz todas as modificações para manter integridade
      // das informações do banco
			if bytes.Equal(doc.dados, content) {
				paginaAtual.Docs = append(paginaAtual.Docs[:i], paginaAtual.Docs[i+1:]...)
				paginaAtual.num_docs--
        
        // Caso a página fique vazia, reorganizamos os ponteiros para o SGBD continuar consistente
				if paginaAtual.num_docs == 0 {
					if paginaAnterior == nil {
						s.PrimeiraPagina = paginaAtual.next
					} else {
						paginaAnterior.next = paginaAtual.next
					}
					if paginaAtual == s.UltimaPagina {
						s.UltimaPagina = paginaAnterior
					}
					s.NumPaginas--
					s.updatePageId(paginaAtual.next)
				} else {
          // Caso a página continue existindo, nós atualizamos tudo para manter a integridade do Banco
					paginaAtual.arm_livre = paginaAtual.arm_livre + len(doc.dados)
					paginaAtual.updateSeq()
				}
				return nil
			}
		}
		if paginaAtual.next == nil {
			break
		}
		paginaAnterior = paginaAtual
		paginaAtual = paginaAtual.next
	}
	return errors.New("Documento não encontrado")
}