package analizador

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// Struct Mbr
type Mbr struct {
	Mbr_tamano         [10]byte
	Mbr_fecha_creacion [20]byte
	Mbr_dsk_signature  [10]byte
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
	Part_start  [10]byte
	Part_size   [10]byte
	Part_name   [20]byte
}

// Separa comandos por línea
func Command(command string) string {
	command = strings.Replace(command, "\n", "", 1)
	comando := strings.Split(command, " ")
	paramm := comando[1:]
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
		unit, path, fit, pivS string
		size                  int
		piv                   []string
		alerta                bool
	)
	datos := Mbr{}
	// Obtenemos los valores de cada parametro
	for i := 0; i < len(paramm); i++ {
		piv = strings.Split(paramm[i], "=")
		if piv[0] == ">SIZE" {
			size, _ = strconv.Atoi(piv[1])
		}
		if piv[0] == ">FIT" {
			fit = piv[1]
		}
		if piv[0] == ">UNIT" {
			unit = piv[1]
		}
		if piv[0] == ">PATH" {
			path = piv[1]
		}
	}
	// Validando parametros obligatorios
	if path != "" && size > 0 {
		alerta = false
		// Unidades
		if unit == "K" {
			size = size * 1024
			pivS = strconv.Itoa(size)
			copy(datos.Mbr_tamano[:], pivS)
		} else if unit == "M" || unit == "" {
			size = size * (1024 * 1024)
			pivS = strconv.Itoa(size)
			copy(datos.Mbr_tamano[:], pivS)
		} else {
			fmt.Println("Valor de unidades incorrecta")
			alerta = true
		}
		// Ajuste
		if fit == "" {
			fit = "FF"
			copy(datos.Dsk_fit[:], fit)
		} else if fit == "FF" || fit == "BF" || fit == "WF" {
			copy(datos.Dsk_fit[:], fit)
		} else {
			fmt.Println("Ajuste ingresado incorrectamente")
			alerta = true
		}
	} else {
		fmt.Println("No cumple con los parámetros mínimos")
		alerta = true
	}
	// Terminando de llenar el Mbr
	copy(datos.Mbr_fecha_creacion[:], hora())
	copy(datos.Mbr_dsk_signature[:], random())
	// Asignando Particiones vacías para control de FDISK
	stt := "0"
	copy(datos.Mbr_partition_1.Part_status[:], stt)
	copy(datos.Mbr_partition_2.Part_status[:], stt)
	copy(datos.Mbr_partition_3.Part_status[:], stt)
	copy(datos.Mbr_partition_4.Part_status[:], stt)
	// Creando disco
	if !alerta {
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
		binary.Write(&bufferMbr, binary.BigEndian, &datos)
		writeB(file, bufferMbr.Bytes())
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
func random() []byte {
	random := rand.Int()
	b := []byte(strconv.Itoa(random))
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
		unit, path, fit, tipo, nombre string
		size                          int
		piv                           []string
		alerta                        bool
	)
	// Obtenemos los valores de cada parametro
	for i := 0; i < len(paramm); i++ {
		piv = strings.Split(paramm[i], "=")
		piv[0] = strings.ToUpper(piv[0])
		if piv[0] == ">SIZE" {
			size, _ = strconv.Atoi(piv[1])
		} else if piv[0] == ">FIT" {
			fit = piv[1]
		} else if piv[0] == ">UNIT" {
			unit = strings.ToUpper(piv[1])
		} else if piv[0] == ">PATH" {
			path = piv[1]
		} else if piv[0] == ">TYPE" {
			tipo = strings.ToUpper(piv[1])
		} else if piv[0] == ">NAME" {
			nombre = piv[1]
		} else {
			fmt.Println("Parametro incorrecto" + piv[0])
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
		pivS := strconv.Itoa(size)
		copy(particion.Part_size[:], pivS)

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
	if dataMBR.Mbr_partition_1.Part_status[0] != '-' || dataMBR.Mbr_partition_2.Part_status[0] != '-' || dataMBR.Mbr_partition_3.Part_status[0] != '-' || dataMBR.Mbr_partition_4.Part_status[0] != '-' {
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
	if string(dataMBR.Mbr_partition_1.Part_name[:]) != nombre && string(dataMBR.Mbr_partition_2.Part_name[:]) != nombre && string(dataMBR.Mbr_partition_3.Part_name[:]) != nombre && string(dataMBR.Mbr_partition_4.Part_name[:]) != nombre {
		return false
	}
	return true
}

//Función para crear particiones primarias
func makePrinary(path string, particion Partition) {

}

//Función para crear particiones primarias
func makeExtended(path string, particion Partition) {

}
