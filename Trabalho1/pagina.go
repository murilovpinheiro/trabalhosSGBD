package main

//Pagina e construtor para pagina

type Pagina struct {
	page_id   int
	arm_livre int
	num_docs  int
	Docs      []*Documento
	next      *Pagina // ponteiro para próxima página da lista encadeada
}

func NewPagina(page_id int) *Pagina {
	return &Pagina{
		page_id:   page_id,
		arm_livre: 5,
		num_docs:  0,
		Docs:      []*Documento{},
		next:      nil,
	}
}
// Add um documento doc numa página
func (p *Pagina) AddDocumento(doc *Documento) bool {
	// checar se documento não um ponteiro vazio
  if doc == nil{
    return false
  }
	p.num_docs = p.num_docs + 1

	doc.did.seq = p.num_docs - 1
	doc.did.page_id = p.page_id

	p.Docs = append(p.Docs, doc)
	p.arm_livre = p.arm_livre - doc.did.tam
	return true
}

// Basicamente reconta a Seq e atualiza para manter a integridade dos valores de Seq da página
func(pagina *Pagina) updateSeq() {
	i := 0
	for _, doc := range pagina.Docs {
		doc.did.seq = i
		i++
	}
}

// Mesmo que acima, porém com o pageID
func (p *Pagina) AtualizarPageID(pageID int) {
	for _, doc := range p.Docs {
		doc.did.page_id = pageID
	}
	if p.next != nil {
		p.next.AtualizarPageID(pageID + 1)
	}
}