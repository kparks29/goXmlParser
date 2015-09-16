package main

import (
	"fmt"
	"github.com/kparks29/homdna-xml-parser/appraisal"
	"log"
	"os"
)

func main() {

	file, err := appraisal.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	appraisalResponse, err := appraisal.UpdateHomdnaModel(file)
	if err != nil {
		log.Fatalln(err)
	}

	// result := appraisal.ParseXml(file)
	// homdnaVersion := &appraisal.HomdnaResponse{}
	// if appraisalConfig, err := appraisal.LoadAppraisalConfig("./appraisal.conf"); err != nil {
	// 	log.Fatalln("Error loading HOMDNA appraisal configuration file.", err)
	// }
	// if err = appraisalConfig.Validate(); err != nil {
	// 	log.Fatalln(err)
	// }

	// serviceApiKey := &appraisalConfig.ServiceApiKey

	// address := appraisal.CreateAddressModel(result)
	// appliances := appraisal.CreateApplianceModels(result.Property.Structure.KitchenAppliances)
	// lotFeatures := appraisal.CreateLotFeatureModels(result.Property.Structure.LotFeatures)
	// structureFeatures := appraisal.CreateStructureFeatureModels(result.Property.Structure.StructureFeatures)
	// lotSize := appraisal.GetSize(result.Property.LotInfo.LotSize)
	// lot := appraisal.CreateLot(&lotSize, lotFeatures)
	// structures := appraisal.CreateStructures(result.Property.Structure, structureFeatures, appliances)
	// homdna := appraisal.CreateHomdnaModel(address, lot, structures)

	//post homdna
	// homdnaVersion, err = appraisal.PostNewHomdna(homdna, serviceApiKey)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// md5Hash := appraisal.GetFileMd5(&result.Report.Document.File)
	// fmt.Printf("\n%x\n\n", md5Hash)

	// res := appraisal.PostFile(homdnaVersion, md5Hash, &result.Report.Document, serviceApiKey)
}
