# Sistemas Gerenciadores de Bancos de Dados - Trabalho 2

## Especificação do Trabalho
O objetivo do trabalho consiste em desenvolver um gerenciador de transações de banco de dados com o suporte para o controle de concorrência baseado na técnica de bloqueio em duas fases.

O programa ter a duas classes principais: TrManager e LockTable. A classe TrManager gerencia o estado de cada uma das transações concorrentes. Essa classe deve ter dois atributos. O atributo TrId, do tipo inteiro, será o identificador da transação e representará um “timestamp”. Assim, a transação de TrId=2 é mais nova que a transação de TrId=1. O atributo Status irá armazenar o estado da transação: ativa, concluída, abortada, esperando. A classe LockTable gerencia a tabela de bloqueios. Essa classe deve ter cinco atributos. O atributo IdItem representa o identificador do item bloqueado. O atributo TrId representa o identificador da transação que bloqueou o item de identificador IdItem. O atributo Escopo representa o escopo do bloqueio: objeto ou predicado. O atributo Duração representa a duração do bloqueio: curta ou longa. O atributo Tipo representa o tipo do bloqueio: leitura ou escrita. A classe LockTable deve ter pelo menos as seguintes funções:
- RL(Tr, D) : Insere um bloqueio de leitura na LockTable sobre o item D para a 
transação Tr, se for possível.
- WL(Tr, D) : Insere um bloqueio de escrita na LockTable sobre o item D para a
transação Tr, se for possível.
- UL(Tr, D) : Apaga o bloqueio da transação Tr sobre o item D na LockTable.

Implemente um Grafo de Espera (Wait For) para identificar deadlocks.

Implemente uma estrutura de dados chamada WaitItem para manter, para cada item de dado bloqueado “i”, uma lista FIFO com os identificadores de transações que estão esperando pelo item “i”.

Implemente uma classe chamada Scheduler para realizar o controle de concorrência. Essa classe deve implementar a técnica de bloqueio em duas fases. Deve ainda realizar o controle de deadlocks baseado na estratégia de Wait Die. 

A interface de entrada do programa principal dever a permitir ao usuário: 
1. Definir o nível de isolamento a ser utilizado (read uncommitted, read committed, repeatable read ou serializable)
2. Entrar com escalonamento do tipo BT(1)r1(x)BT(2)w2(x)r2(y)r1(y)C(1)r2(z)C(2) onde:
    - BT(X): inicia a transação X
    - r1(x): transação 1 deseja ler o item x, portanto solicita bloqueio de leitura sobre o item x
    - w1(x): transação 1 deseja escrever o item x, portanto solicita bloqueio de escrita sobre o item x
    - C(X): Validação da transação X, quando todos os seus bloqueios (de longa duração) devem ser liberados
3. Seguir o escalonamento recebido, exibindo as operações realizadas, bem como o estado das estruturas de dados utilizadas, à cada operação.
