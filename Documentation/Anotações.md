# Arquivo de Configuração
 # ConnectionMongoDB
    Array de string de conexão para o MongoDB
        
 # FileLog
    Array de string para os arquivos de log.
   * LogBlockchain.txt para salvar o ultimo bloco salvo
   * LogIndiceEndereco.txt" para salvar o ultimo indice do Endereco salvo
   * LogEndereco.txt para salvar o último endereço salvo
   * LogEnderecosSemDados para salvar todos os enderecos que a api retorna sem dados
   * LogIndiceMultiEndereco.txt para salvar o ultimo indice do MultiEndereco salvo
   * LogMultiEnderecosSemDados.txt para salvar todos os multienderecos que a api retorna sem dados

 # DataBase
    Array de databases criados no mongoDB

 # Collection
    Array de todas as collections que tem no bando de dados
   * Adresses terá as transações de cada endereco.
   * MapeandoEnderecos usada para salvar os enderecos que estão mapeados
   * blockchain usada para salvar os blocos da blockchain
   * teste,testeMultiAdress usadas para realizar testes
      
 # UrlAPI
    Array de todas as uri de api 
 # O restante dos dados são os end points da api que está sendo usada nesse projeto
    

 # Controllers
  * AuxAddrControllers.go contêm implementações de funções auxiliares para os arquivos MultiEnderecoController.go e EnderecoController.go