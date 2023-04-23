package analizador

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// Struct Mbr
type Mbr struct {
	Mbr_tamano         int32
	Mbr_fecha_creacion [20]byte
	Mbr_dsk_signature  int32
	Dsk_fit            [2]byte
	Mbr_partition_1    Partition
	Mbr_partition_2    Partition
	Mbr_partition_3    Partition
	Mbr_partition_4    Partition
}

// Struct Partition
type Partition struct {
	Part_status [1]byte
	Part_type   [1]byte
	Part_fit    [2]byte
	Part_start  int32
	Part_size   int32
	Part_name   [20]byte
}

//Struct EBR
type Ebr struct {
	Part_status [1]byte
	Part_fit    [2]byte
	Part_start  int32
	Part_size   int32
	Part_next   int32
	Part_name   [20]byte
}

// Separa comandos por línea
func Command(command string) string {
	command = strings.Replace(command, "\n", "", 1)
	comando := strings.Split(command, " ")
	paramm := comando[1:]
	comando[0] = strings.ToUpper(comando[0])
	//Llama metodos por comando
	switch comando[0] {
	case "MKDISK":
		mkdisk(paramm)
	case "RMDISK":
		rmdisk(paramm)
	case "FDISK":
		fdisk(paramm)
	default:
		fmt.Println("Comando no reconocido")
	}
	return comando[0]
}

// Lee línea completa desde consolas
func GetCommand() {
	entrada := bufio.NewReader(os.Stdin) //Obtiene todos los items después de 0
	lectura, _ := entrada.ReadString('\n')
	salida := Command(lectura)
	fmt.Println(salida)
}

// Comando MKDISK
func mkdisk(paramm []string) {
	// variables para parametros
	var (
		unit, path, fit, pivP string
		size                  int
		piv                   []string
		alerta                bool
	)
	// Obtenemos los valores de cada parametro
	for i := 0; i < len(paramm); i++ {
		piv = strings.Split(paramm[i], "=")
		pivP = strings.TrimPrefix(piv[0], ">")
		pivP = strings.ToUpper(pivP)

		if pivP == "SIZE" {
			size, _ = strconv.Atoi(piv[1])
		} else if pivP == "FIT" {
			fit = strings.ToUpper(piv[1])
		} else if pivP == "UNIT" {
			unit = strings.ToUpper(piv[1])
		} else if pivP == "PATH" {
			path = piv[1]
		} else {
			fmt.Println("Parametro incorrecto " + pivP)
			alerta = true
		}
	}
	// Validando parametros obligatorios
	// Validando size mayor a 0
	if size < 0 {
		fmt.Println("Parametro SIZE debe ser mayor a cero")
		alerta = true
	}
	//Validando que el path no venga vacío
	if path == "" {
		fmt.Println("Parametro PATH es requerido")
		alerta = true
	}
	// Unidades
	if unit == "K" {
		size = size * 1024
	} else if unit == "M" || unit == "" {
		size = size * (1024 * 1024)
	} else {
		fmt.Println("Unidades incorrecta")
		alerta = true
	}
	// Ajuste
	if fit == "" {
		fit = "FF"
	}
	if fit != "FF" && fit != "BF" && fit != "WF" {
		fmt.Println("Ajuste ingresado incorrectamente")
		alerta = true
	}

	// Creando disco
	if !alerta {
		dataMBR := Mbr{}
		//Agregando valores a MBR
		dataMBR.Mbr_tamano = int32(size)
		copy(dataMBR.Mbr_fecha_creacion[:], hora())
		dataMBR.Mbr_dsk_signature = random()
		copy(dataMBR.Dsk_fit[:], fit)
		// Asignando Particiones vacías para control de FDISK
		stt := "-"
		copy(dataMBR.Mbr_partition_1.Part_status[:], stt)
		copy(dataMBR.Mbr_partition_2.Part_status[:], stt)
		copy(dataMBR.Mbr_partition_3.Part_status[:], stt)
		copy(dataMBR.Mbr_partition_4.Part_status[:], stt)

		//Validando/Creando carpetas
		findingRuta(path)

		file, err := os.Create(path)
		if err != nil { // Valida que el disco esté vacío
			log.Fatal(err)
		}
		defer file.Close()
		// Llenar el disco con ceros
		var temp int8 = 0
		temporal := &temp
		var binario bytes.Buffer
		binary.Write(&binario, binary.BigEndian, temporal)
		for i := 0; i < size; i++ {
			writeB(file, binario.Bytes())
		}
		// Agregar Mbr a disco
		file.Seek(0, 0) // Posiciona al inicio del archivo
		var bufferMbr bytes.Buffer
		binary.Write(&bufferMbr, binary.BigEndian, &dataMBR)
		writeB(file, bufferMbr.Bytes())
	} else {
		fmt.Println("Disco no creado")
	}
}

// Obtiene hora del sistema
func hora() []byte {
	t := time.Now()
	fmt.Println(t.Format("2006-01-02 15:04:05"))
	b := []byte(t.String())
	return b
}

// Id de disco
func random() int32 {
	random := rand.Int()
	b := int32(random)
	fmt.Println("ID", random)
	return b
}

// Escribe bytes en el disco
func writeB(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}

// Comando RMDISK
func rmdisk(paramm []string) {
	// Variables de parametros
	var (
		path string
	)
	// Obtenemos el path
	piv := strings.Split(paramm[0], "=")
	path = piv[1]
	fmt.Println("Eliminando archivo...")
	err := os.Remove(path)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Disco eliminado con exito")
	}
}

// Comando FDISK
func fdisk(paramm []string) {
	// variables para parametros
	var (
		unit, path, fit, tipo, nombre, pivP string
		size                                int
		piv                                 []string
		alerta                              bool
	)
	// Obtenemos los valores de cada parametro
	for i := 0; i < len(paramm); i++ {
		piv = strings.Split(paramm[i], "=")
		pivP = strings.TrimPrefix(piv[0], ">")
		pivP = strings.ToUpper(pivP)
		if pivP == "SIZE" {
			size, _ = strconv.Atoi(piv[1])
		} else if pivP == "FIT" {
			fit = strings.ToUpper(piv[1])
		} else if pivP == "UNIT" {
			unit = strings.ToUpper(piv[1])
		} else if pivP == "PATH" {
			path = piv[1]
		} else if pivP == "TYPE" {
			tipo = strings.ToUpper(piv[1])
		} else if pivP == "NAME" {
			nombre = piv[1]
		} else {
			fmt.Println("Parametro incorrecto" + pivP)
			alerta = true
		}
	}
	// Validamos comandos obligatorios
	// Validando size mayor a 0
	if size < 0 {
		fmt.Println("Parametro SIZE debe ser mayor a cero")
		alerta = true
	}
	// Validando Tamaño partition
	if unit == "B" {
	}
	if unit == "K" || unit == "" {
		size = size * 1024
	} else if unit == "M" {
		size = size * (1024 * 1024)
	} else {
		fmt.Println("Unidades incorrectas")
		alerta = true
	}
	// Tipo de partición
	if tipo == "" {
		tipo = "P"
	}
	if tipo != "E" && tipo != "P" && tipo != "L" {
		fmt.Println("Tipo de partición incorrecto")
		alerta = true
	}
	//Validando que no exista extendida
	if tipo == "E" {
		if findingExt(path) == true {
			alerta = true
			fmt.Println("Ya existe una partición extendida")
		}
	}
	///Si viene type L, validamos que sí exista partición extendida
	if tipo == "L" {
		if findingExt(path) != true {
			alerta = true
			fmt.Println("No existe partición extendida para escribir particiones lógicas...")
		}
	}

	// Ajuste
	if fit == "" {
		fit = "WF"
	}
	if fit != "FF" && fit != "BF" && fit != "WF" {
		fmt.Println("Ajuste ingresado incorrectamente")
		alerta = true
	}

	//Validando que el disco a modificar exista
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		alerta = true
		fmt.Println("Disco no existente")
	}
	defer file.Close()

	//Verificando que las particiones aún no están llenas
	if partitionFull(path) == true {
		alerta = true
		fmt.Println("Particiones llenas")
	}

	//Verificando que el nombre no se repita
	if nombre == "" {
		alerta = true
		fmt.Println("Error, nombre requerido")
	}
	if namePartition(path, nombre) == true {
		alerta = true
		fmt.Println("El nombre: " + nombre + " ya ha sido ingresado a las particiones.")
	}

	//Creando partición
	if alerta != true {
		particion := Partition{}
		status := "0"
		//Agregando valores al struct Partition
		copy(particion.Part_status[:], status)
		copy(particion.Part_type[:], tipo)
		copy(particion.Part_fit[:], fit)
		copy(particion.Part_name[:], nombre)
		particion.Part_size = int32(size)

		//Enviando la información hacia el tipo de partición
		if tipo == "P" {
			makePrinary(path, particion)
		} else if tipo == "E" {
			makeExtended(path, particion)
		} else if tipo == "L" {

		}

	} else {
		fmt.Println("Error: la partición no ha sido creada")
	}
}

// Lee bytes del disco
func readB(file *os.File, number int) []byte {
	bytes := make([]byte, number) //arreglo de bytes
	_, err := file.Read(bytes)    // bytes leidos
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

//Función para validar si existe partición extendida
func findingExt(path string) bool {
	dataMBR := Mbr{}
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var mSize int = int(unsafe.Sizeof(dataMBR))
	data := readB(file, mSize)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &dataMBR) //Error
	if err != nil {
		log.Fatal("Falló lectura binaria", err)
	}
	if dataMBR.Mbr_partition_1.Part_type[0] != 'E' && dataMBR.Mbr_partition_2.Part_type[0] != 'E' && dataMBR.Mbr_partition_3.Part_type[0] != 'E' && dataMBR.Mbr_partition_4.Part_type[0] != 'E' {
		return false
	}
	return true
}

//Función para validar si las particiones están llenas
func partitionFull(path string) bool {
	dataMBR := Mbr{}
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var mSize int = int(unsafe.Sizeof(dataMBR))
	data := readB(file, mSize)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &dataMBR) //Error
	if err != nil {
		log.Fatal("Falló lectura binaria", err)
	}
	if dataMBR.Mbr_partition_1.Part_status[0] == '-' || dataMBR.Mbr_partition_2.Part_status[0] == '-' || dataMBR.Mbr_partition_3.Part_status[0] == '-' || dataMBR.Mbr_partition_4.Part_status[0] == '-' {
		return false
	}
	return true
}

//Función para validar que el nombre de la partición no se repita
func namePartition(path string, nombre string) bool {
	dataMBR := Mbr{}
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var mSize int = int(unsafe.Sizeof(dataMBR))
	data := readB(file, mSize)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &dataMBR) //Error
	if err != nil {
		log.Fatal("Falló lectura binaria", err)
	}
	if string(dataMBR.Mbr_partition_1.Part_name[:len(nombre)]) != nombre && string(dataMBR.Mbr_partition_2.Part_name[:len(nombre)]) != nombre && string(dataMBR.Mbr_partition_3.Part_name[:len(nombre)]) != nombre && string(dataMBR.Mbr_partition_4.Part_name[:len(nombre)]) != nombre {
		return false
	}
	return true
}

//Función para crear particiones primarias
func makePrinary(path string, particion Partition) {
	dataMBR := Mbr{}
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var mSize int = int(unsafe.Sizeof(dataMBR))
	data := readB(file, mSize)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &dataMBR) //Error
	if err != nil {
		log.Fatal("Falló lectura binaria", err)
	}
	numPart := 0
	if dataMBR.Mbr_partition_1.Part_status[0] == '-' {
		numPart = 1
	} else if dataMBR.Mbr_partition_2.Part_status[0] == '-' {
		numPart = 2
	} else if dataMBR.Mbr_partition_3.Part_status[0] == '-' {
		numPart = 3
	} else if dataMBR.Mbr_partition_4.Part_status[0] == '-' {
		numPart = 4
	}
	fmt.Println(numPart)
	//Primera partición en el disco
	sizeMBR := int32(mSize + 1)
	if numPart == 1 {
		particion.Part_start = sizeMBR
		dataMBR.Mbr_partition_1 = particion
	} else if numPart == 2 {
		sizeMBR += dataMBR.Mbr_partition_1.Part_size
		particion.Part_start = sizeMBR
		dataMBR.Mbr_partition_2 = particion
	} else if numPart == 3 {
		sizeMBR += dataMBR.Mbr_partition_1.Part_size + dataMBR.Mbr_partition_2.Part_size
		particion.Part_start = sizeMBR
		dataMBR.Mbr_partition_3 = particion
	} else if numPart == 4 {
		sizeMBR += dataMBR.Mbr_partition_1.Part_size + dataMBR.Mbr_partition_2.Part_size + dataMBR.Mbr_partition_3.Part_size
		particion.Part_start = sizeMBR
		dataMBR.Mbr_partition_4 = particion
	}
	//Actualizando MBR
	file.Seek(0, 0) // Posiciona al inicio del archivo
	var bufferMbr bytes.Buffer
	binary.Write(&bufferMbr, binary.BigEndian, &dataMBR)
	writeB(file, bufferMbr.Bytes()) //Error acá
	fmt.Println("MBR actualizado")
}

//Función para crear particiones primarias
func makeExtended(path string, particion Partition) {
	dataMBR := Mbr{}
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var mSize int = int(unsafe.Sizeof(dataMBR))
	data := readB(file, mSize)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &dataMBR) //Error
	if err != nil {
		log.Fatal("Falló lectura binaria", err)
	}
	numPart := 0
	if dataMBR.Mbr_partition_1.Part_status[0] == '-' {
		numPart = 1
	} else if dataMBR.Mbr_partition_2.Part_status[0] == '-' {
		numPart = 2
	} else if dataMBR.Mbr_partition_3.Part_status[0] == '-' {
		numPart = 3
	} else if dataMBR.Mbr_partition_4.Part_status[0] == '-' {
		numPart = 4
	}
	fmt.Println(numPart)
	//Primera partición en el disco
	sizeMBR := int32(mSize + 1)
	if numPart == 1 {
		particion.Part_start = sizeMBR
		dataMBR.Mbr_partition_1 = particion
	} else if numPart == 2 {
		sizeMBR += dataMBR.Mbr_partition_1.Part_size
		particion.Part_start = sizeMBR
		dataMBR.Mbr_partition_2 = particion
	} else if numPart == 3 {
		sizeMBR += dataMBR.Mbr_partition_1.Part_size + dataMBR.Mbr_partition_2.Part_size
		particion.Part_start = sizeMBR
		dataMBR.Mbr_partition_3 = particion
	} else if numPart == 4 {
		sizeMBR += dataMBR.Mbr_partition_1.Part_size + dataMBR.Mbr_partition_2.Part_size + dataMBR.Mbr_partition_3.Part_size
		particion.Part_start = sizeMBR
		dataMBR.Mbr_partition_4 = particion
	}
	//Actualizando MBR
	file.Seek(0, 0) // Posiciona al inicio del archivo
	var bufferMbr bytes.Buffer
	binary.Write(&bufferMbr, binary.BigEndian, &dataMBR)
	writeB(file, bufferMbr.Bytes())
	fmt.Println("MBR actualizado")

	//Creando EBR inicial
	dataEBR := Ebr{}
	statusEbr := "-"
	copy(dataEBR.Part_status[:], statusEbr)
	dataEBR.Part_fit = particion.Part_fit
	dataEBR.Part_start = particion.Part_start
	dataEBR.Part_size = particion.Part_size
	dataEBR.Part_next = -1
	dataEBR.Part_name = particion.Part_name

	//Escribiendo EBR en el disco
	sizeMBR++
	file.Seek(int64(sizeMBR), os.SEEK_SET) // Posiciona al inicio del archivo
	var bufferEbr bytes.Buffer
	binary.Write(&bufferEbr, binary.BigEndian, &dataEBR)
	writeB(file, bufferEbr.Bytes())
	fmt.Println("EBR actualizado")
}

//Función para buscar path
func findingRuta(ruta string) {
	err := crearRuta(ruta)
	if err != nil {
		fmt.Println(err)

	}
	//fmt.Println("La ruta ha sido creada correctamente")
}

//Función para crear path
func crearRuta(ruta string) error {
	_, err := os.Stat(ruta)
	if os.IsNotExist(err) {
		// La ruta no existe, se debe crear
		err = os.MkdirAll(filepath.Dir(ruta), 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		// Ocurrió un error al verificar la existencia de la ruta
		return err
	}
	return nil
}
