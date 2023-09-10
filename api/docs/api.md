# Veneza

## Header Parser

```
    curl -X 'POST' \
    'http://localhost:8080/parser/item' \
    -H 'Content-Type: application/json' \
    -d '{
        "text": "210 - Emitir Pedido de Compra - Com saldo de pedido                                                                                      Data : 03/05/2023                                                  Página 1 de 4
		*** PEDIDO EMITIDO NA UNID. MASTER ***                      Número do Pedido:                   48740        N.VerbaBonific 0            Filial: 1
																COMPRA MERC.                                                       Nro da Negociação: 0         Tipo de carga: Paletizado
		
		Empresa:           FAKE                                                                        CNPJ/CPF:           00000000000000
		Endereço:          FAKE                      915   PARQUE DEZ                                                  Inscrição Estadual: 042133815
		Cidade:            MANAUS                          UF: AM     CEP: FAKE                                                 Telefone:           0000000000
		
		Fornecedor: FAKE S/A                                                  Codigo:  114                      CNPJ/CPF:00000000000000      Inscrição Estadual: 000000000
		Endereço: AV FAKE                                                 FAKE FAKE FAKE                        Telefone: 0000000000
																																							Informações Bancárias:
		Cidade:     FAKE                                                     UF: AM    CEP: 00000000                          Fax:
		Representante: FAKE                                                                                                    Telefone: 00 000000000       Banco :          Agência:     C.Corrente:
		Comprador:     FAKE DA SILVA                                                                                 Fax:
		Faturamento Embarque       Entrega                                          Condições de Pagamento : 30                                        Dias a partir da data de FATURAMENTO
		03/05/2023    03/05/2023   17/05/2023                                       Frete : CIF Transportadora :"
    }
```

## Item Parser

```
    curl -X 'POST' \
    'http://localhost:8080/parser/item' \
    -H 'accept: */*' \
    -H 'Content-Type: application/json' \
    -d '{
        "text": "25982   DEL
         - VALLE FRESH LARANJA + UVA 1,5L     7894900558852 PT 2X1,5L
                                                                PT         10,00        0,00    6,150000   0,00   6,150000   0,00   0,00   0,00   0,00   6,150000   0,00   0,00   0,00    0,00   0,00   0,00   0,00
30374   ENERGETICO
         -         MONSTER DRAGON PESSEGO 473ML
                                           1220000250413 PT 6X473ML
                                                             PT             5,00        0,00   30,829998   0,00 30,829998    0,00   0,00   0,00   0,00 30,829998    0,00   0,00   0,00    0,00   0,00   0,00   0,00
"
    }'
```
