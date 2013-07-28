package csv_controller

import (
	"encoding/csv"
	//"io"
	"logger"
	"os"
)

func WriteCSV(path string, startline []string, data [][]string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f)
	w.Write(startline)

	for i := 0; i < len(data); i++ {
		w.Write(data[i])
	}

	w.Flush()
}

func ReadCSV(path string) (data [][]string) {
	file, err := os.Open(path)
	if err != nil {
		logger.Debugln("Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)

	data, err2 := reader.ReadAll()

	if err2 != nil {
		logger.Debugln("Error:", err2)
		return
	}
	return data
	//flag := 0
	//for {
	//	record, err := reader.Read()
	//	if  err == io.EOF {
	//		break
	//	} else if err != nil {
	//		logger.Debugln("Error:", err)
	//		return
	//	}
	//	data[flag] = record
	//	flag++
	//	//fmt.Println(record) // record has the type []string
	//}

}
