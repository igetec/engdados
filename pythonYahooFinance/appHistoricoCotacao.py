# Mais informações sobre a yfinance:
#  https://algotrading101.com/learn/yfinance-guide/

import yfinance as yf
import pandas as pd

arquicoCSV = "01_CotacaoAcoes.csv"

# Define o símbolo da ação
acao_simbolo = 'PETR4.SA'

# Extract
acao_dados = yf.Ticker(acao_simbolo)

# Transform
acao_historico = acao_dados.history(period='1mo')  # Período de 1 mês 
df_acao_historico = pd.DataFrame(acao_historico)

# Load
print("*** Dados históricos de 1 mês *** \n")
print(df_acao_historico[["Open", "High", "Close"]])
df_acao_historico[["Open", "High", "Close"]].to_csv(arquicoCSV, index=False)

