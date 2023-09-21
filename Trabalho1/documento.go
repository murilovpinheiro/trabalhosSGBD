package main

import (
	"errors"
)

//Documentos e DID e construtor do Documento

type DID struct {
	page_id int
	seq     int
	tam     int
}

type Documento struct {
	did   DID    // -> Meio que os metadados do documento
	dados []byte // -> Dados: "Sequência de caracteres"
}

func NewDocumento(dados []byte) (*Documento, error) {
	if len(dados) < 1 || len(dados) > 5 {
		return nil, errors.New("Os dados do documento devem ter entre 1 e 5 caracteres")
	}
	// Cria um novo objeto Documento com o campo dados inicializado com um slice vazio de bytes do mesmo tamanho da sequência de caracteres.
	documento := &Documento{
		did: DID{
			page_id: -1,
			seq:     -1,
			tam:     len(dados),
		},
		dados: make([]byte, len(dados)),
	}

	// Copia os dados fornecidos para o campo dados do objeto Documento.
	copy(documento.dados, dados)

	// Retorna um ponteiro para o novo objeto Documento e nenhum erro.
	return documento, nil
}
