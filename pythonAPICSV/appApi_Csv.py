import pandas as pd
import requests

# Identificação das APIs
apiURLProdutosEstatisticas = "https://olinda.bcb.gov.br/olinda/servico/PTAX/versao/v1/odata/CotacaoDolarPeriodo(dataInicial=@dataInicial,dataFinalCotacao=@dataFinalCotacao)?@dataInicial='01-01-2023'&@dataFinalCotacao='12-31-2023'&$top=100&$format=json&$select=cotacaoCompra,cotacaoVenda,dataHoraCotacao"

# Identificação dos arquivos
arquivoCSV = "01_CotacaoDolar2023.csv"

## Extract
response = requests.get(apiURLProdutosEstatisticas)

# Verifica se a request na API deu certo
if response.status_code == 200:
    data = response.json()

else:
    print(f"Falha em buscar a API. Código do status: {response.status_code}")
    data = []
    exit()

## Transform
# Criando um dataframe com os dados da API
df = pd.DataFrame(data["value"])

## Load
df.to_csv(arquivoCSV, index=False)

print(f"Dados salvos com sucesso no arquivo {arquivoCSV}")
