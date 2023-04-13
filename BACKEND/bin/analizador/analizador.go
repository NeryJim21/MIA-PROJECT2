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
		path, confirma string
	)
	// Obtenemos el path
	piv := strings.Split(paramm[0], "=")
	path = piv[1]
	fmt.Println("Eliminar archivo...")
	fmt.Println("S/N ")
	fmt.Print("-> ")
	fmt.Scanln(&confirma)
	if confirma == "S" || confirma == "s" {
		err := os.Remove(path)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Disco eliminado con exito")
		}
	} else {
		fmt.Println("Disco no eliminado")
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
	datos := Mbr{}
	particion := Partition{}
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
		if piv[0] == ">TYPE" {
			tipo = piv[1]
		}
		if piv[0] == ">NAME" {
			nombre = piv[1]
		}
	}
	// Validamos comandos obligatorios
	if size > 0 && path != "" && nombre != "" {
		alerta = false
		// Unidades
		if unit == "B" {
			pivS := strconv.Itoa(size)
			copy(particion.Part_size[:], pivS)
		} else if unit == "K" || unit == "" {
			size = size * 1024
			pivS := strconv.Itoa(size)
			copy(particion.Part_size[:], pivS)
		} else if unit == "M" {
			size = size * (1024 * 1024)
			pivS := strconv.Itoa(size)
			copy(particion.Part_size[:], pivS)
		} else {
			fmt.Println("Valor de unidades incorrecta")
			alerta = true
		}
		// Tipo de partición
		if tipo == "" {
			tipo = "P"
			copy(particion.Part_type[:], tipo)
		} else if tipo == "E" || tipo == "P" || tipo == "L" {
			copy(particion.Part_type[:], tipo)
		} else {
			fmt.Println("Tipo de partición incorrecto")
			alerta = true
		}
		// Ajuste
		if fit == "" {
			fit = "WF"
			copy(particion.Part_fit[:], fit)
		} else if fit == "FF" || fit == "BF" || fit == "WF" {
			copy(particion.Part_fit[:], fit)
		} else {
			fmt.Println("Ajuste ingresado incorrectamente")
			alerta = true
		}
	} else {
		fmt.Println("No cumple con los parámetros obligatorios")
	}
	//Terminando de llenar struct de partición
	copy(particion.Part_name[:], nombre)
	stt := "1"
	copy(particion.Part_status[:], stt)
	copy(particion.Part_start[:], nombre) // Modificar por valor real
	// Buscando particiones disponibles
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var mSize int = int(unsafe.Sizeof(datos))
	data := readB(file, mSize)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &datos) //Error
	if err != nil {
		log.Fatal("Falló lectura binaria", err)
	}
	// Escogiendo partición libre
	pivSt := 0
	partSt := string(datos.Mbr_partition_1.Part_status[:])
	if partSt == "0" {
		pivSt = 1
	} else {
		partSt = string(datos.Mbr_partition_2.Part_status[:])
		if partSt == "0" {
			pivSt = 2
		} else {
			partSt = string(datos.Mbr_partition_3.Part_status[:])
			if partSt == "0" {
				pivSt = 3
			} else {
				partSt = string(datos.Mbr_partition_4.Part_status[:])
				if partSt == "0" {
					pivSt = 4
				} else {
					alerta = true
					fmt.Println("Particiones llenas")
				}
			}
		}
	}
	//Validando nombre repetido
	if string(datos.Mbr_partition_1.Part_name[:]) == nombre {
		alerta = true
	}
	if string(datos.Mbr_partition_2.Part_name[:]) == nombre {
		alerta = true
	}
	if string(datos.Mbr_partition_3.Part_name[:]) == nombre {
		alerta = true
	}
	if string(datos.Mbr_partition_4.Part_name[:]) == nombre {
		alerta = true
	}
	//Validar particiones extendidas
	if tipo == "E" {
		if string(datos.Mbr_partition_1.Part_type[:]) == "E" {
			alerta = true
		}
		if string(datos.Mbr_partition_2.Part_type[:]) == "E" {
			alerta = true
		}
		if string(datos.Mbr_partition_3.Part_type[:]) == "E" {
			alerta = true
		}
		if string(datos.Mbr_partition_4.Part_type[:]) == "E" {
			alerta = true
		}
	}
	//Validar particiones logicas
	if tipo == "L" {
		pivL1 := string(datos.Mbr_partition_1.Part_type[:])
		pivL2 := string(datos.Mbr_partition_2.Part_type[:])
		pivL3 := string(datos.Mbr_partition_3.Part_type[:])
		pivL4 := string(datos.Mbr_partition_4.Part_type[:])
		if pivL1 == "E" || pivL2 == "E" || pivL3 == "E" || pivL4 == "E" {
			alerta = false
		} else {
			alerta = true
			fmt.Println("No existe partición extendida para crear particiones lógicas")
		}
	}
	//Creando partición
	if !alerta {
		switch pivSt {
		case 1:
			fmt.Println("Espacio en la partición 1")
		case 2:
			fmt.Println("Espacio en la partición 2")
		case 3:
			fmt.Println("Espacio en la partición 3")
		case 4:
			fmt.Println("Espacio en la partición 4")
		default:
			fmt.Println("Partición no disponible")
		}
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