package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Str_EquipeNasESF struct {
	EquipeNasESFSK      int64   `gorm:"column:equipenasesfsk;primaryKey"`
	CoMunicipio         string  `gorm:"column:co_municipio;type:bigint" csv:"CO_MUNICIPIO"`
	CoArea              string  `gorm:"column:co_area" csv:"CO_AREA"`
	SeqEquipe           string  `gorm:"column:seq_equipe;type:bigint" csv:"SEQ_EQUIPE"`
	CoMunicipioEsf      string  `gorm:"column:co_municipio_esf;type:bigint" csv:"CO_MUNICIPIO_ESF"`
	CoUnidade           string  `gorm:"column:co_unidade;type:bigint" csv:"CO_UNIDADE"`
	NoFantasiaEsf       string  `gorm:"column:no_fantasia_esf" csv:"NO_FANTASIA_ESF"`
	CoSegmentoEsf       string  `gorm:"column:co_segmento_esf" csv:"CO_SEGMENTO_ESF"`
	DsSegmentoEsf       string  `gorm:"column:ds_segmento_esf" csv:"DS_SEGMENTO_ESF"`
	DsAreaEsf           string  `gorm:"column:ds_area_esf" csv:"DS_AREA_ESF"`
	DtAtualizacao       string  `gorm:"column:dt_atualizacao;type:date" csv:"DT_ATUALIZACAO"`
	DtAtualizacaoOrigem *string `gorm:"column:dt_atualizacao_origem;type:date" csv:"DT_ATUALIZACAO_ORIGEM"`
}

func (Str_EquipeNasESF) TableName() string {
	return "bi_equipenasesf"
}

func parseDataFormat(datePar string) string {
	formatedDate, err := time.Parse("02/01/2006", datePar)

	if err != nil {
		log.Fatalln("Data com formato inconsistente")
	}

	return formatedDate.Format("2006-01-02")
}

func parseNullableString(value string) *string {
	if value == "" {
		return nil
	}
	value = parseDataFormat(value)
	return &value
}

func ConnectDB() *gorm.DB {
	DBHost := "172.17.0.1"
	DBPort := 5432
	DBUser := "postgres"
	DBPassword := "postdba"
	DBName := "datasus"

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBPassword, DBName)

	// Initializando GORM com PostgreSQL driver
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Falha ao conectar com o banco de dados " + DBName)
	}

	// Se a tabela já existir, eu apago antes.
	if db.Migrator().HasTable("bi_equipenasesf") {
		//@ renomei a bi_ftleads com dados válidos para depois apagá-la caso o rename da stage dê certo
		if err := db.Migrator().DropTable("bi_equipenasesf"); err != nil {
			log.Println("Houve erro em APAGAR a tabela de dados")
			os.Exit(1)
		}
	}

	// Auto-migrate para criar a tabela de Stage
	if err := db.AutoMigrate(&Str_EquipeNasESF{}); err != nil {
		log.Fatalf("Falha ao criar a tabela de dados: %v", err)
	}
	return db
}

func Extract() [][]string {
	file, err := os.Open("rlEquipeNasfEsf202411.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	readerUTF := transform.NewReader(file, charmap.Windows1252.NewDecoder())
	reader := csv.NewReader(readerUTF)
	reader.Comma = ';'
	rows, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	return rows
}

func Load(dbPar *gorm.DB, dfPar [][]string) {

	for i, row := range dfPar {
		if i == 0 {
			continue // Ignora o header
		}

		equipeNasESF := Str_EquipeNasESF{
			CoMunicipio:         row[0],
			CoArea:              row[1],
			SeqEquipe:           row[2],
			CoMunicipioEsf:      row[3],
			CoUnidade:           row[7],
			NoFantasiaEsf:       row[10],
			CoSegmentoEsf:       row[11],
			DsSegmentoEsf:       row[12],
			DsAreaEsf:           row[13],
			DtAtualizacao:       parseDataFormat(row[14]),
			DtAtualizacaoOrigem: parseNullableString(row[16]),
		}

		dbPar.Create(&equipeNasESF)
	}

}

func ETL() {
	fmt.Println("Conectando ao BD")
	db := ConnectDB()

	fmt.Println("Realizando a Extração")
	df := Extract()

	fmt.Println("Executando a carga de dados")
	Load(db, df)
}

func main() {

	ETL()
}
