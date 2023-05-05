package main

import (
	"bufio"
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

//=============================== MAIN ===============================

// Main function
func main() {
	ListaDiscos := list.New()
	LlenarListaDisco(ListaDiscos)
	//LeerArchivo de Entrada
	var comando string = ""
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\n\n ")
	fmt.Println("# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #")
	fmt.Println("# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #")
	fmt.Println("# # #                                                                 # # #")
	fmt.Println("# # #                                                                 # # #")
	fmt.Println("# # #         S I S T E M A   D E   A R C H I V O S   E X T 2         # # #")
	fmt.Println("# # #                                                                 # # #")
	fmt.Println("# # #                                                                 # # #")
	fmt.Println("# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #")
	fmt.Println("# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #")
	fmt.Println("\n\n\n ")
	for comando != "exit" {
		fmt.Println("\n# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #")
		fmt.Print("Ingrese comando \n-> ")
		scanner.Scan()
		comando = scanner.Text()
		if comando != "" {
			LeerTexto(comando, ListaDiscos)
		}

	}
}
func LlenarListaDisco(ListaDiscos *list.List) {
	IdDisco := [26]string{"a", "b", "c", "d", "e", "f", "g", "h", "i",
		"j", "k", "l", "m", "n", "o", "p", "q",
		"r", "s", "t", "u", "v", "w", "x", "y", "z"}
	for i := 0; i < 26; i++ {
		disco := DISCO{}
		copy(disco.Estado[:], "0")
		copy(disco.Id[:], IdDisco[i])
		for j := 0; j < len(disco.Particiones); j++ {
			mount := MOUNT{}
			mount.Part_name = ""
			mount.Id = strconv.Itoa(j + 1)
			copy(mount.Estado[:], "0")
			disco.Particiones[j] = mount
		}
		ListaDiscos.PushBack(disco)
	}
}

func RecorrerListaDisco(id string, ListaDiscos *list.List) (string, string, string) {
	Id := strings.ReplaceAll(id, "81", "")
	IdDisco := Id[1:]
	pathDisco := ""
	Part_name := ""
	nombreDisco := ""
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disco DISCO
		disco = element.Value.(DISCO)
		if BytesToString(disco.Id) == IdDisco {
			for i := 0; i < len(disco.Particiones); i++ {
				var mountTemp = disco.Particiones[i]
				if mountTemp.Id == id {
					copy(mountTemp.EstadoMKS[:], "1")
					Part_name = mountTemp.Part_name
					pathDisco = disco.Path
					nombreDisco = disco.NombreDisco
					return pathDisco, Part_name, nombreDisco
				}
			}

		}
		element.Value = disco
	}
	return "", "", ""
}

//============================== FIN -> MAIN =============================

//=============================== STRUCTS ===============================

//Estructura para cada Comando y sus Propiedades
type Propiedad struct {
	Name string
	Val  string
}
type Comando struct {
	Name        string
	Propiedades []Propiedad
}

//Estructuras para el Disco y Particiones
type Particion struct {
	Part_status [1]byte
	Part_type   [1]byte
	Part_fit    [2]byte
	Part_start  int64
	Part_size   int64
	Part_name   [15]byte
}

//Struct para el MBR
type MBR struct {
	Mbr_tamano         int64
	Mbr_fecha_creacion [19]byte
	Mbr_dsk_signature  int64
	Dsk_fit            [2]byte
	Particiones        [4]Particion
}

//Struct para las particiones Logicas
type EBR struct {
	Part_status [1]byte
	Part_fit    [2]byte
	Part_start  int64
	Part_next   int64
	Part_size   int64
	Part_name   [15]byte
}

//EStruc de las particiones montadas
type MOUNT struct {
	Part_name string
	Id        string
	Estado    [1]byte
	EstadoMKS [1]byte
}

//Estruct Disco
type DISCO struct {
	NombreDisco string
	Path        string
	Id          [1]byte
	Estado      [1]byte
	Particiones [100]MOUNT
}

//57.51
type Integers struct {
	I1  uint16
	I2  int32
	I3  int64
	DOS byte
}

//Structuras Segunda Fase
//SuperBloque
type SB struct {
	Sb_nombre_hd                          [15]byte
	Sb_arbol_virtual_count                int64
	Sb_detalle_directorio_count           int64
	Sb_inodos_count                       int64
	Sb_bloques_count                      int64
	Sb_arbol_virtual_free                 int64
	Sb_detalle_directorio_free            int64
	Sb_inodos_free                        int64
	Sb_bloques_free                       int64
	Sb_date_creacion                      [19]byte
	Sb_date_ultimo_montaje                [19]byte
	Sb_montajes_count                     int64
	Sb_ap_bitmap_arbol_directorio         int64
	Sb_ap_arbol_directorio                int64
	Sb_ap_bitmap_detalle_directorio       int64
	Sb_ap_detalle_directorio              int64
	Sb_ap_bitmap_tabla_inodo              int64
	Sb_ap_tabla_inodo                     int64
	Sb_ap_bitmap_bloques                  int64
	Sb_ap_bloques                         int64
	Sb_ap_log                             int64
	Sb_size_struct_arbol_directorio       int64
	Sb_size_struct_Detalle_directorio     int64
	Sb_size_struct_inodo                  int64
	Sb_size_struct_bloque                 int64
	Sb_first_free_bit_arbol_directorio    int64
	Sb_first_free_bit_detalle_directoriio int64
	Sb_dirst_free_bit_tabla_inodo         int64
	Sb_first_free_bit_bloques             int64
	Sb_magic_num                          int64
	InicioCopiaSB                         int64
	ConteoAVD                             int64
	ConteoDD                              int64
	ConteoInodo                           int64
	ConteoBloque                          int64
}

//Arbol virtual de directorio
type AVD struct {
	Avd_fecha_creacion              [19]byte
	Avd_nomre_directotrio           [15]byte
	Avd_ap_array_subdirectoios      [6]int64
	Avd_ap_detalle_directorio       int64
	Avd_ap_arbol_virtual_directorio int64
	Avd_proper                      [10]byte
}

//Detalle de Directorio

type ArregloDD struct {
	Dd_file_nombre            [15]byte
	Dd_file_ap_inodo          int64
	Dd_file_date_creacion     [19]byte
	Dd_file_date_modificacion [19]byte
}
type DD struct {
	Dd_array_files           [5]ArregloDD
	Dd_ap_detalle_directorio int64
	Ocupado                  int8
}

//Cantidad de Inodos
type Inodo struct {
	I_count_inodo             int64
	I_size_archivo            int64
	I_count_bloques_asignados int64
	I_array_bloques           [4]int64
	I_ao_indirecto            int64
	I_id_proper               int64
}

//Bloque
type Bloque struct {
	//Db_data [25]byte
	Db_data [64]byte
}

//bitacora
type Bitacora struct {
	Log_tipo_operacion [19]byte
	Log_tipo           [1]byte
	Log_nombre         [35]byte
	Log_Contenido      [25]byte
	Log_fecha          [19]byte
	Size               int64
}

//
func BytesPart_name(data [15]byte) string {
	return string(data[:])
}

/*func ConvertData(data [25]byte) string {
	return string(data[:])
}*/
func ConvertData(data [64]byte) string {
	return string(data[:])
}

//=============================== FIN -> STRUCTS ===============================

//=============================== ANALIZADOR ===============================
var global string = ""
var globalIdLogin string = ""

//Funcion para leer y reconocer los comandos lleno la lista de comandos
func LeerTexto(dat string, ListaDiscos *list.List) {
	//Leendo la cadena de entrada
	ListaComandos := list.New()
	lineaComando := strings.Split(dat, "\n")
	var c Comando
	for i := 0; i < len(lineaComando); i++ {
		EsComentario := lineaComando[i][0:1]
		if EsComentario != "#" {
			comando := lineaComando[i]
			if strings.Contains(lineaComando[i], "\\*") {
				comando = strings.Replace(lineaComando[i], "\\*", " ", 1) + lineaComando[i+1]
				i = i + 1
			}
			propiedades := strings.Split(string(comando), " ")
			//Nombre Comando
			nombreComando := propiedades[0]
			//Struct para el Comando
			c = Comando{Name: strings.ToLower(nombreComando)}
			propiedadesTemp := make([]Propiedad, len(propiedades)-1)
			for i := 1; i < len(propiedades); i++ {
				if propiedades[i] == "" {
					continue
				} else if propiedades[i] == ">p" {
					propiedadesTemp[i-1] = Propiedad{Name: ">p",
						Val: ">p"}
				} else if propiedades[i] == ">r" {
					propiedadesTemp[i-1] = Propiedad{Name: ">r",
						Val: ">r"}
				} else {
					if strings.Contains(propiedades[i], "=") {
						valor_propiedad_Comando := strings.Split(propiedades[i], "=")
						propiedadesTemp[i-1] = Propiedad{Name: valor_propiedad_Comando[0],
							Val: valor_propiedad_Comando[1]}
					} else {
						propiedadesTemp[i-1] = Propiedad{Name: ">sigue",
							Val: propiedades[i]}
					}
				}
			}
			c.Propiedades = propiedadesTemp
			//Agregando el comando a la lista comandos
			ListaComandos.PushBack(c)
		}
	}
	RecorrerListaComando(ListaComandos, ListaDiscos)
}

//Funcion para recorrer la Lista de Comandos
func RecorrerListaComando(ListaComandos *list.List, ListaDiscos *list.List) {
	var ParamValidos bool = true
	var cont = 1
	for element := ListaComandos.Front(); element != nil; element = element.Next() {
		var comandoTemp Comando
		comandoTemp = element.Value.(Comando)
		//Lista de propiedades del Comando
		switch strings.ToLower(comandoTemp.Name) {
		case "mkdisk":
			ParamValidos = EjecutarComandoMKDISK(comandoTemp.Name, comandoTemp.Propiedades, cont)
			cont++
			if ParamValidos == false {
				fmt.Println("*Error: en parametros")
			}
		case "rmdisk":
			ParamValidos = EjecutarComandoRMDISK(comandoTemp.Name, comandoTemp.Propiedades)
			if ParamValidos == false {
				fmt.Println("*Error: en parametros")
			}
		case "fdisk":
			ParamValidos = EjecutarComandoFDISK(comandoTemp.Name, comandoTemp.Propiedades)
			if ParamValidos == false {
				fmt.Println("*Error: en parametros")
			}
		case "mount":
			if len(comandoTemp.Propiedades) != 0 {
				ParamValidos = EjecutarComandoMount(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
				if ParamValidos == false {
					fmt.Println("*Error: en parametros")
				}
			} else {
				EjecutarReporteMount(ListaDiscos)
			}

		case "exit":
			fmt.Println("Hasta la proxima :)")
		case "pause":
			fmt.Print("Presiona Enter para continuar...")
			exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
			os.Stdin.Read(make([]byte, 1))
			fmt.Println()
		case "cls":
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
		case "exec":
			ParamValidos = EjecutarComandoExec(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if ParamValidos == false {
				fmt.Println("*Error: en parametros")
			}
		case "mkdir":
			ParamValidos = EjecutarComandoMKDIR(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if ParamValidos == false {
				fmt.Println("*Error: en parametros")
			}
		case "mkfile":
			ParamValidos = EjecutarComandoMKFILE(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if ParamValidos == false {
				fmt.Println("*Error: en parametros")
			}
		case "mkfs":
			ParamValidos = EjecutarComandoMKFS(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if ParamValidos == false {
				fmt.Println("*Error: en parametros")
			}
		case "rep":
			ParamValidos = EjecutarComandoReporte(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if ParamValidos == false {
				fmt.Println("*Error: en parametros")
			}
		case "login":
			ParamValidos, global = EjecutarComandoLogin(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if ParamValidos == false {
				fmt.Println("*Error: en parametros")
			}
		case "logout":
			if global == "" {
				fmt.Println("No existe sesión iniciada")
			} else {
				global = ""
				globalIdLogin = ""
				fmt.Println("Sesion finalizada")
			}
		default:
			fmt.Println("*Error: Comando no reconocido")
		}
	}
}

//=============================== FIN -> ANALIZADOR ===============================

//=============================== MKDISK ===============================

func EjecutarComandoMKDISK(nombreComando string, propiedadesTemp []Propiedad, cont int) (ParamValidos bool) {
	dt := time.Now()
	mbr1 := MBR{}
	copy(mbr1.Mbr_fecha_creacion[:], dt.String())
	mbr1.Mbr_dsk_signature = int64(rand.Intn(100) + cont)
	fmt.Println("->Ejecutando MKDISK...")
	comandos := "dd if=/dev/zero "
	ParamValidos = true
	var propiedades [4]string
	if len(propiedadesTemp) >= 2 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			//Vector temporal de propiedades
			switch strings.ToLower(nombrePropiedad) {
			case ">size":
				propiedades[0] = propiedadTemp.Val
			case ">unit":
				propiedades[2] = strings.ToLower(propiedadTemp.Val)
			case ">fit":
				//propiedades[2] = strings.ToLower(propiedadTemp.Val)
			case ">path":
				propiedades[3] = propiedadTemp.Val
				arr_path := strings.Split(propiedades[3], "/")
				pathCompleta := ""
				for i := 0; i < len(arr_path)-1; i++ {
					pathCompleta += arr_path[i] + "/"
				}
				executeComand("mkdir " + pathCompleta)
				comandos += "of=" + propiedades[3]
			default:
				fmt.Println("Error al ejecutar el comando")
				ParamValidos = false
				return ParamValidos
			}
		}
		EsComilla := propiedades[3][0:1]
		if EsComilla == "\"" {
			propiedades[3] = propiedades[3][1 : len(propiedades[3])-1]
		}
		Part_size, _ := strconv.ParseInt(propiedades[0], 10, 64)
		if propiedades[2] == "k" {
			comandos += " bs=" + strconv.Itoa((int(Part_size))*1000) + " count=1"
			mbr1.Mbr_tamano = ((Part_size) - 1) * 1000
		} else {
			comandos += " bs=" + strconv.Itoa(int(Part_size)) + "MB" + " count=1"
			mbr1.Mbr_tamano = Part_size * 1000000
		}
		//Inicializando Particiones
		for i := 0; i < 4; i++ {
			copy(mbr1.Particiones[i].Part_status[:], "0")
			copy(mbr1.Particiones[i].Part_type[:], "")
			copy(mbr1.Particiones[i].Part_fit[:], "")
			mbr1.Particiones[i].Part_start = 0
			mbr1.Particiones[i].Part_size = 0
			copy(mbr1.Particiones[i].Part_name[:], "")
		}
		executeComand(comandos)
		//Escribir MBR
		f, err := os.OpenFile(propiedades[3], os.O_WRONLY, 0755)
		if err != nil {
			log.Fatalln(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatalln(err)
			}
		}()
		f.Seek(0, 0)
		err = binary.Write(f, binary.BigEndian, mbr1)
		if err != nil {
			log.Fatalln(err, propiedades[3])
		}
		fmt.Println("¡Disco creado correctamente!")
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}

func executeComand(comandos string) {
	args := strings.Split(comandos, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.CombinedOutput()
}
func BytesToString(data [1]byte) string {
	return string(data[:])
}
func CheckError(e error) {
	if e != nil {
		fmt.Println("Error - ----------")
		fmt.Println(e)
	}
}

//=============================== FIN -> MKDISK ===============================

//=============================== RMDISK ===============================
func EjecutarComandoRMDISK(nombreComando string, propiedadesTemp []Propiedad) (ParamValidos bool) {
	fmt.Println("->Ejecutando RMDISK...")
	ParamValidos = true
	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case ">path":
				f, err := os.OpenFile(propiedadTemp.Val, os.O_RDWR, 0755)
				if err != nil {
					fmt.Println("No existe el disco en la ruta " + propiedadTemp.Val)
					return false
				}
				scanner := bufio.NewScanner(os.Stdin)
				fmt.Print("¿Esta seguro de eliminar el disco? s/n >>")
				scanner.Scan()
				if scanner.Text() == "s" {
					executeComand("rm " + propiedadTemp.Val)
					fmt.Println("¡Disco eliminado correctamente!")
				} else {
					fmt.Println("Disco no eliminado")
				}
				defer f.Close()
			default:
				fmt.Println("Error al ejecutar el comando")
				ParamValidos = false
				return ParamValidos
			}
		}
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}

//=============================== FIN -> RMDISK ===============================

//=============================== FDISK ===============================
func EjecutarComandoFDISK(nombreComando string, propiedadesTemp []Propiedad) (ParamValidos bool) {
	fmt.Println("->Ejecutando FDISK...")
	ParamValidos = true
	mbr := MBR{}
	particion := Particion{}
	var startPart int64 = int64(unsafe.Sizeof(mbr))
	var propiedades [8]string
	if len(propiedadesTemp) >= 2 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case ">size":
				propiedades[0] = propiedadTemp.Val
			case ">fit":
				propiedades[1] = propiedadTemp.Val
			case ">unit":
				propiedades[2] = propiedadTemp.Val
			case ">path":
				propiedades[3] = propiedadTemp.Val
			case ">type":
				propiedades[4] = propiedadTemp.Val
			case ">delete":
				propiedades[5] = propiedadTemp.Val
			case ">name":
				propiedades[6] = propiedadTemp.Val
			case ">add":
				propiedades[7] = propiedadTemp.Val
			default:
				fmt.Println("Error al ejecutar el comando")
				ParamValidos = false
				return ParamValidos
			}
		}
		if propiedades[4] == "" {
			propiedades[4] = "p"
		}
		EsComilla := propiedades[3][0:1]
		if EsComilla == "\"" {
			propiedades[3] = propiedades[3][1 : len(propiedades[3])-1]
		}
		//Tamanio Particion
		var Part_sizeParticion int64 = 0
		if strings.ToLower(propiedades[2]) == "b" {
			TamanioParticion, _ := strconv.ParseInt(propiedades[0], 10, 64)
			Part_sizeParticion = TamanioParticion
		} else if strings.ToLower(propiedades[2]) == "k" {
			TamanioParticion, _ := strconv.ParseInt(propiedades[0], 10, 64)
			Part_sizeParticion = TamanioParticion * 1000
		} else if strings.ToLower(propiedades[2]) == "m" {
			TamanioParticion, _ := strconv.ParseInt(propiedades[0], 10, 64)
			Part_sizeParticion = TamanioParticion * 1000000
		} else {
			TamanioParticion, _ := strconv.ParseInt(propiedades[0], 10, 64)
			Part_sizeParticion = TamanioParticion * 1000
		}
		if propiedades[5] != "" {
			EliminarParticion(propiedades[3], propiedades[6], propiedades[5])
			return
		}

		//Obtener el MBR
		switch strings.ToLower(propiedades[4]) {
		case "p":
			var Particiones [4]Particion
			f, err := os.OpenFile(propiedades[3], os.O_RDWR, 0755)
			if err != nil {
				fmt.Println("No existe la ruta " + propiedades[3])
				return false
			}
			defer f.Close()
			f.Seek(0, 0)
			err = binary.Read(f, binary.BigEndian, &mbr)
			Particiones = mbr.Particiones
			if err != nil {
				fmt.Println("No existe el archivo en la ruta")
			}
			//El mbr ya se a leido, 2.Verificar si existe espacion disponible o que no lo rebase
			if HayEspacio(Part_sizeParticion, mbr.Mbr_tamano) {
				return false
			}

			//Verificar si ya hay 4 particiones creadas
			if BytesToString(Particiones[3].Part_status) == "1" {
				fmt.Println("*ERROR: Ya existen 4 particiones")
				return false
			}
			//Verificar si ya hay particiones
			if BytesToString(Particiones[0].Part_status) == "1" {
				for i := 0; i < 4; i++ {
					//Posicion en bytes del partstar de la n particion
					startPart += Particiones[i].Part_size
					if BytesToString(Particiones[i].Part_status) == "0" {
						break
					}
				}
			}
			if HayEspacio(startPart+Part_sizeParticion, mbr.Mbr_tamano) {
				return false
			}
			//dando valores a la particion
			copy(particion.Part_status[:], "1")
			copy(particion.Part_type[:], propiedades[4])
			copy(particion.Part_fit[:], propiedades[1])
			particion.Part_start = startPart
			particion.Part_size = Part_sizeParticion
			copy(particion.Part_name[:], propiedades[6])
			//Particion creada
			for i := 0; i < 4; i++ {
				if BytesToString(Particiones[i].Part_status) == "0" {
					Particiones[i] = particion
					break
				}
			}
			f.Seek(0, 0)
			mbr.Particiones = Particiones
			err = binary.Write(f, binary.BigEndian, mbr)
			ReadFile(propiedades[3])
			fmt.Println("¡Particion primaria creada correctamente!")
		case "l":
			if !HayExtendida(propiedades[3]) {
				fmt.Println("*ERROR: No existe una particion Extendida")
				return false
			}
			ebr := EBR{}
			copy(ebr.Part_status[:], "1")
			copy(ebr.Part_fit[:], propiedades[1])
			ebr.Part_start = startPart
			ebr.Part_next = 0
			ebr.Part_size = Part_sizeParticion
			copy(ebr.Part_name[:], propiedades[6])
			//Obteniendo el byte donde empezara la particion Logica
			InicioParticionLogica(propiedades[3], ebr)
			fmt.Println("¡Particion logica creada correctamente!")
		case "e":
			//Particiones Extendidas
			var Particiones [4]Particion
			f, err := os.OpenFile(propiedades[3], os.O_RDWR, 0755)
			if err != nil {
				fmt.Println("No existe la ruta " + propiedades[3])
				return false
			}
			defer f.Close()
			f.Seek(0, 0)
			err = binary.Read(f, binary.BigEndian, &mbr)
			Particiones = mbr.Particiones
			if err != nil {
				fmt.Println("*Error: No existe el archivo en la ruta")
			}
			//El mbr ya se a leido,2.Verificar si existe espacio disponible o que no lo rebase
			if HayEspacio(Part_sizeParticion, mbr.Mbr_tamano) {
				return false
			}

			//Verificar si ya hay 4 particiones creadas
			if BytesToString(Particiones[3].Part_status) == "1" {
				fmt.Println("*ERROR: Ya existen 4 particiones")
				return false
			}
			//Verificar si ya hay particiones
			if BytesToString(Particiones[0].Part_status) == "1" {
				for i := 0; i < 4; i++ {
					//Posicion en bytes del partstar de la n particion
					startPart += Particiones[i].Part_size
					if BytesToString(Particiones[i].Part_status) == "0" {
						break
					}
				}
			}
			if HayEspacio(startPart+Part_sizeParticion, mbr.Mbr_tamano) {
				return false
			}
			//dando valores a la particion
			copy(particion.Part_status[:], "1")
			copy(particion.Part_type[:], propiedades[4])
			copy(particion.Part_fit[:], propiedades[1])
			particion.Part_start = startPart
			particion.Part_size = Part_sizeParticion
			copy(particion.Part_name[:], propiedades[6])
			//Particion creada
			for i := 0; i < 4; i++ {
				if BytesToString(Particiones[i].Part_status) == "0" {
					Particiones[i] = particion
					break
				}
			}
			f.Seek(0, 0)
			mbr.Particiones = Particiones
			err = binary.Write(f, binary.BigEndian, mbr)
			ReadFile(propiedades[3])
			fmt.Println("¡Particion extendida creada correctamente!")
			ebr := EBR{}
			copy(ebr.Part_status[:], "1")
			copy(ebr.Part_fit[:], propiedades[1])
			ebr.Part_start = startPart
			ebr.Part_next = -1
			ebr.Part_size = Part_sizeParticion
			copy(ebr.Part_name[:], propiedades[6])
			f.Seek(ebr.Part_start, 0)
			err = binary.Write(f, binary.BigEndian, ebr)
		default:
			fmt.Println("Ocurrio un error")
		}
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}

func EscribirParticionLogica(path string, ebr EBR, inicioParticionLogica int64, inicioParticionExtendida int64) bool {
	ebr.Part_start = inicioParticionLogica
	ebr.Part_next = inicioParticionLogica + ebr.Part_size
	return true
}
func EliminarParticion(path string, name string, typeDelete string) bool {
	var name2 [15]byte
	Encontrada := false
	copy(name2[:], name)
	f, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta " + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	//Posiciono al inicio el Puntero
	f.Seek(0, 0)
	//Leo el mbr
	err = binary.Read(f, binary.BigEndian, &mbr)
	Particiones := mbr.Particiones
	for i := 0; i < 4; i++ {
		if strings.ToLower(BytesToString(Particiones[i].Part_type)) == "e" && BytesPart_name(Particiones[i].Part_name) == BytesPart_name(name2) {
			fmt.Println("Es una Extendida")
			Encontrada = true
		} else if strings.ToLower(BytesToString(Particiones[i].Part_type)) == "p" && BytesPart_name(Particiones[i].Part_name) == BytesPart_name(name2) {
			var partTemp = Particion{}
			copy(partTemp.Part_status[:], "0")
			copy(partTemp.Part_type[:], "")
			copy(partTemp.Part_fit[:], "")
			partTemp.Part_start = 0
			partTemp.Part_size = 0
			copy(partTemp.Part_name[:], "")
			Particiones[i] = partTemp
			mbr.Particiones = Particiones
			f.Seek(0, 0)
			err = binary.Write(f, binary.BigEndian, &mbr)
			fmt.Println("Particon Primaria Eliminada")
			ReadFile(path)
			Encontrada = true
		}
	}
	if Encontrada == false {
		for i := 0; i < 4; i++ {
			if strings.ToLower(BytesToString(Particiones[i].Part_type)) == "e" {
				var InicioExtendida int64 = Particiones[i].Part_start
				f.Seek(InicioExtendida, 0)
				ebrAnterior := EBR{}
				ebr := EBR{}
				ebrAnterior = ebr
				err = binary.Read(f, binary.BigEndian, &ebr)
				if ebr.Part_next == -1 {
					fmt.Println("No Hay particiones Logicas")
				} else {
					f.Seek(InicioExtendida, 0)
					err = binary.Read(f, binary.BigEndian, &ebr)
					for {
						if BytesPart_name(ebr.Part_name) == BytesPart_name(name2) {
							fmt.Println("Particion Logica Encontrada")
							if strings.ToLower(typeDelete) == "fast" {
								ebrAnterior.Part_next = ebr.Part_next
								f.Seek(ebrAnterior.Part_start, 0)
								err = binary.Write(f, binary.BigEndian, ebrAnterior)

							} else if strings.ToLower(typeDelete) == "full" {
								ebrAnterior.Part_next = ebr.Part_next
								f.Seek(ebrAnterior.Part_start, 0)
								err = binary.Write(f, binary.BigEndian, ebrAnterior)
							}
							Encontrada = true
						}
						if ebr.Part_next == -1 {
							break
						} else {
							f.Seek(ebr.Part_next, 0)
							ebrAnterior = ebr
							err = binary.Read(f, binary.BigEndian, &ebr)
						}
					}
				}
			}
		}
	}
	if Encontrada == false {
		fmt.Println("*Error: no se encontro la particion")
	}
	return false
}
func InicioParticionLogica(path string, ebr2 EBR) bool {
	f, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta " + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	Particiones := mbr.Particiones
	for i := 0; i < 4; i++ {
		if strings.ToLower(BytesToString(Particiones[i].Part_type)) == "e" {
			var InicioExtendida int64 = Particiones[i].Part_start
			f.Seek(InicioExtendida, 0)
			ebr := EBR{}
			err = binary.Read(f, binary.BigEndian, &ebr)
			if ebr.Part_next == -1 {
				ebr.Part_next = ebr.Part_start + int64(unsafe.Sizeof(ebr)) + ebr2.Part_size
				f.Seek(InicioExtendida, 0)
				err = binary.Write(f, binary.BigEndian, ebr)
				ebr2.Part_start = ebr.Part_next
				ebr2.Part_next = -1
				f.Seek(ebr2.Part_start, 0)
				err = binary.Write(f, binary.BigEndian, ebr2)

				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				return false
			} else {
				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				for {
					if ebr.Part_next == -1 {
						ebr.Part_next = ebr.Part_start + int64(unsafe.Sizeof(ebr)) + ebr2.Part_size
						f.Seek(ebr.Part_start, 0)
						err = binary.Write(f, binary.BigEndian, ebr)
						ebr2.Part_start = ebr.Part_next
						ebr2.Part_next = -1
						f.Seek(ebr2.Part_start, 0)
						err = binary.Write(f, binary.BigEndian, ebr2)
						break
					} else {
						f.Seek(ebr.Part_next, 0)
						err = binary.Read(f, binary.BigEndian, &ebr)
					}
				}
				return false
			}
		}
	}
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}
	return false
}
func HayExtendida(path string) bool {
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta " + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	Particiones := mbr.Particiones
	for i := 0; i < 4; i++ {
		if strings.ToLower(BytesToString(Particiones[i].Part_type)) == "e" {
			return true
		}
	}
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}
	return false
}

func ReadFile(path string) (funciona bool) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta " + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}
	return true
}

func HayEspacio(Part_sizeParticion int64, tamanioDisco int64) bool {
	if ((Part_sizeParticion) > tamanioDisco) || (Part_sizeParticion < 0) {
		fmt.Println("*ERROR: el tamanio de la particion es mayor al tamanio disponible del disco")
		return true
	}
	return false
}

//=============================== FIN -> FDISK ===============================

//=============================== MOUNT ===============================
func EjecutarComandoMount(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	fmt.Println("->Ejecutando MOUNT...")
	var propiedades [2]string
	var nombre [15]byte
	ParamValidos = true
	if len(propiedadesTemp) >= 2 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case ">name":
				propiedades[0] = propiedadTemp.Val
				copy(nombre[:], propiedades[0])
			case ">path":
				propiedades[1] = propiedadTemp.Val
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		//Empezar a montar las Particiones
		EjecutarComando(propiedades[1], nombre, ListaDiscos)
		EjecutarReporteMount(ListaDiscos)
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}
func EjecutarReporteMount(ListaDiscos *list.List) {
	fmt.Println("\n - - - - - - - - - - P A R T I C I O N E S  M O N T A D A S - - - - - - - - -\n")
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disco DISCO
		disco = element.Value.(DISCO)
		if disco.NombreDisco != "" {
			for i := 0; i < len(disco.Particiones); i++ {
				var mountTemp = disco.Particiones[i]
				if mountTemp.Part_name != "" {
					fmt.Println("ID:", mountTemp.Id, "  Disco:", disco.Path, "  Name:", mountTemp.Part_name)
				}
			}
		}
	}
}
func IdValido(id string, ListaDiscos *list.List) bool {
	esta := false
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disco DISCO
		disco = element.Value.(DISCO)
		if disco.NombreDisco != "" {
			for i := 0; i < len(disco.Particiones); i++ {
				var mountTemp = disco.Particiones[i]
				if mountTemp.Part_name != "" {
					fmt.Println("-z»", mountTemp.Id)
					fmt.Println("-x»", id)
					if mountTemp.Id == id {
						return true
					}
				}
			}
		}
	}
	return esta
}
func EjecutarComando(path string, Part_name [15]byte, ListaDiscos *list.List) bool {
	var encontrada = false
	lineaComando := strings.Split(path, "/")
	nombreDisco := lineaComando[len(lineaComando)-1]
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	Particiones := mbr.Particiones
	for i := 0; i < 4; i++ {
		if string(Particiones[i].Part_name[:]) == string(Part_name[:]) {
			encontrada = true
			if strings.ToLower(BytesToString(Particiones[i].Part_type)) == "e" {
				fmt.Println("*ERROR: no se puede montar una particion Extendida")
			} else {
				ParticionMontar(ListaDiscos, string(Part_name[:]), string(nombreDisco), path)
			}
		}
		if strings.ToLower(BytesToString(Particiones[i].Part_type)) == "e" {
			ebr := EBR{}
			f.Seek(Particiones[i].Part_start, 0)
			err = binary.Read(f, binary.BigEndian, &ebr)
			for {
				if ebr.Part_next == -1 {
					break
				} else {
					f.Seek(ebr.Part_next, 0)
					err = binary.Read(f, binary.BigEndian, &ebr)
				}
				var nombre string = string(ebr.Part_name[:])
				var nombre2 string = string(Part_name[:])
				if nombre == nombre2 {
					encontrada = true
					//Montar Particion
					ParticionMontar(ListaDiscos, string(Part_name[:]), string(nombreDisco), path)
				}
			}
		}
	}
	if encontrada == false {
		fmt.Println("*Error: no se encontro la particion")
	}
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}
	return true
}

func ParticionMontar(ListaDiscos *list.List, Part_name string, nombreDisco string, path string) {

	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disco DISCO
		disco = element.Value.(DISCO)
		if BytesToString(disco.Estado) == "0" && !ExisteDisco(ListaDiscos, nombreDisco) {
			disco.NombreDisco = nombreDisco
			disco.Path = path
			copy(disco.Estado[:], "1")
			//#id->vda1
			for i := 0; i < len(disco.Particiones); i++ {
				var mountTemp = disco.Particiones[i]
				if BytesToString(mountTemp.Estado) == "0" {
					mountTemp.Id = "81" + mountTemp.Id + BytesToString(disco.Id)
					mountTemp.Part_name = Part_name
					copy(mountTemp.Estado[:], "1")
					copy(mountTemp.EstadoMKS[:], "0")
					disco.Particiones[i] = mountTemp
					break
				} else if BytesToString(mountTemp.Estado) == "1" && mountTemp.Part_name == Part_name {
					fmt.Println("La Particion ya esta montada")
					break
				}
			}
			element.Value = disco
			break
		} else if BytesToString(disco.Estado) == "1" && ExisteDisco(ListaDiscos, nombreDisco) && nombreDisco == disco.NombreDisco {
			for i := 0; i < len(disco.Particiones); i++ {
				var mountTemp = disco.Particiones[i]
				if BytesToString(mountTemp.Estado) == "0" {
					mountTemp.Id = "81" + mountTemp.Id + BytesToString(disco.Id)
					mountTemp.Part_name = Part_name
					copy(mountTemp.Estado[:], "1")
					copy(mountTemp.EstadoMKS[:], "0")
					disco.Particiones[i] = mountTemp
					break
				} else if BytesToString(mountTemp.Estado) == "1" && mountTemp.Part_name == Part_name {
					fmt.Println("La Particion ya esta montada")
					break
				}
			}
			element.Value = disco
			break
		}
	}
}
func ExisteDisco(ListaDiscos *list.List, nombreDisco string) bool {
	Existe := false
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disco DISCO
		disco = element.Value.(DISCO)
		if disco.NombreDisco == nombreDisco {
			return true
		} else {
			Existe = false
		}
	}
	return Existe
}

//=============================== FIN -> MOUNT ===============================

//=============================== REPORTES ===============================

func EjecutarComandoReporte(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	fmt.Println("->Ejecutando REP...")
	ParamValidos = true
	var propiedades [4]string
	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case ">id":
				propiedades[0] = propiedadTemp.Val
			case ">path":
				propiedades[1] = propiedadTemp.Val
			case ">name":
				propiedades[2] = propiedadTemp.Val
			case ">ruta":
				propiedades[3] = propiedadTemp.Val
			case ">sigue":
				propiedades[1] += propiedadTemp.Val
			default:
				fmt.Println("Error al Ejecutar el Comando", nombrePropiedad)
			}
		}
		EsComilla := propiedades[1][0:1]
		if EsComilla == "\"" {
			if propiedades[3] != "" {
				propiedades[3] = propiedades[3][1 : len(propiedades[3])-1]
			}
			propiedades[1] = propiedades[1][1 : len(propiedades[1])-1]
		}
		carpetas_Graficar := strings.Split(propiedades[1], "/")
		var comando = ""
		for i := 1; i < len(carpetas_Graficar)-1; i++ {
			comando += carpetas_Graficar[i] + "/"
		}

		executeComand("mkdir " + comando[0:len(comando)-1])
		switch strings.ToLower(propiedades[2]) {
		case "disk":
			GraficarDisk(propiedades[0], ListaDiscos, propiedades[1])
		case "file":
			ReporteFile(propiedades[0], propiedades[1], propiedades[3], ListaDiscos)
		case "tree":
			GraficarTreeFull(propiedades[0], propiedades[1], propiedades[3], ListaDiscos)
		default:
			fmt.Println("*ERROR: name incorrecto para el reporte")
		}
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}

//Graficar Disco y calcular Porcentajes
func GraficarDisk(idParticion string, ListaDiscos *list.List, path string) bool {
	var Part_name [15]byte
	var buffer bytes.Buffer
	buffer.WriteString("digraph G{\ntbl [\nshape=box\nlabel=<\n<table border='0' cellborder='2' width='100' height=\"30\" color='orange'>\n<tr>")
	pathDisco, _, _ := RecorrerListaDisco(idParticion, ListaDiscos)
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta " + pathDisco)
		return false
	}
	defer f.Close()
	PorcentajeUtilizado := 0.0
	var EspacioUtilizado int64 = 0
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	TamanioDisco := mbr.Mbr_tamano
	Particiones := mbr.Particiones
	buffer.WriteString("<td height='30' width='75'> MBR </td>")
	for i := 0; i < 4; i++ {
		if convertName(Particiones[i].Part_name[:]) != convertName(Part_name[:]) && strings.ToLower(BytesToString(Particiones[i].Part_type)) == "p" {
			PorcentajeUtilizado = (float64(Particiones[i].Part_size) / float64(TamanioDisco)) * 100
			buffer.WriteString("<td height='30' width='75.0'>PRIMARIA <br/>" + convertName(Particiones[i].Part_name[:]) + " <br/> Ocupado: " + strconv.Itoa(int(PorcentajeUtilizado)) + "%</td>")
			EspacioUtilizado += Particiones[i].Part_size
		} else if convertName(Particiones[i].Part_status[:]) == "0" {
			buffer.WriteString("<td height='30' width='75.0'>Libre</td>")
		}
		if strings.ToLower(BytesToString(Particiones[i].Part_type)) == "e" {
			EspacioUtilizado += Particiones[i].Part_size
			PorcentajeUtilizado = (float64(Particiones[i].Part_size) / float64(TamanioDisco)) * 100
			buffer.WriteString("<td  height='30' width='15.0'>\n")
			buffer.WriteString("<table border='5'  height='30' WIDTH='15.0' cellborder='1'>\n")
			buffer.WriteString(" <tr>  <td height='60' colspan='100%'>EXTENDIDA <br/>" + convertName(Particiones[i].Part_name[:]) + " <br/> Ocupado: " + strconv.Itoa(int(PorcentajeUtilizado)) + "%</td>  </tr>")
			var InicioExtendida int64 = Particiones[i].Part_start
			f.Seek(InicioExtendida, 0)
			ebr := EBR{}
			err = binary.Read(f, binary.BigEndian, &ebr)
			if ebr.Part_next == -1 {
			} else {
				buffer.WriteString("\n<tr>")
				var EspacioUtilizado int64 = 0
				cont := 0
				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				for {
					if ebr.Part_next == -1 {
						break
					} else {
						f.Seek(ebr.Part_next, 0)
						err = binary.Read(f, binary.BigEndian, &ebr)
						EspacioUtilizado += ebr.Part_size
						PorcentajeUtilizado = (float64(ebr.Part_size) / float64(Particiones[i].Part_size)) * 100
						buffer.WriteString("<td height='30'>EBR</td><td height='30'> Logica:  " + convertName(ebr.Part_name[:]) + " " + strconv.Itoa(int(PorcentajeUtilizado)) + "%</td>")
						cont++
					}
				}
				if (Particiones[i].Part_size - EspacioUtilizado) > 0 {
					PorcentajeUtilizado = (float64(TamanioDisco-EspacioUtilizado) / float64(TamanioDisco)) * 100
					buffer.WriteString("<td height='30' width='100%'>Libre: " + strconv.Itoa(int(PorcentajeUtilizado)) + "%</td>")
				}
				buffer.WriteString("</tr>\n")
			}
			buffer.WriteString("</table>\n</td>")
		}
	}
	if (TamanioDisco - EspacioUtilizado) > 0 {
		PorcentajeUtilizado = (float64(TamanioDisco-EspacioUtilizado) / float64(TamanioDisco)) * 100
		buffer.WriteString("<td height='30' width='75.0'>Libre: " + strconv.Itoa(int(PorcentajeUtilizado)) + "%</td>")
	}
	buffer.WriteString("     </tr>\n</table>\n>];\n}")
	var datos string
	datos = string(buffer.String())
	CreateArchivo(path, datos)
	fmt.Println("¡Reporte Disk creado exitosamente!")
	return false
}

func GraficarTreeFull(idParticion string, pathCarpeta string, ruta string, ListaDiscos *list.List) bool {
	var buffer bytes.Buffer
	buffer.WriteString("digraph grafica{\nrankdir=TB;\nnode [shape = record, style=filled, fillcolor=orange];\n")
	sb := SB{}
	var dos [15]byte
	avd := AVD{}
	var strArray [100]string
	pathDisco, Part_name, _ := RecorrerListaDisco(idParticion, ListaDiscos)
	sb, _ = DevolverSuperBlque(pathDisco, Part_name)
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	/*
	   Graficar AVD's
	*/
	f.Seek(sb.Sb_ap_arbol_directorio, 0)
	for i := 0; i < int(sb.Sb_arbol_virtual_count); i++ {
		err = binary.Read(f, binary.BigEndian, &avd)
		if avd.Avd_nomre_directotrio == dos {
			break
		}
		for j := 0; j < 6; j++ {
			if avd.Avd_ap_array_subdirectoios[j] != -1 {
				buffer.WriteString("nodo" + strconv.Itoa(i) + ":f" + strconv.Itoa(j) + " -> nodo" + strconv.Itoa(int(avd.Avd_ap_array_subdirectoios[j])) + "\n")
			} else {
				break
			}
		}
		if avd.Avd_ap_arbol_virtual_directorio != -1 {
			buffer.WriteString("nodo" + strconv.Itoa(i) + ":f7" + " -> nodo" + strconv.Itoa(int(avd.Avd_ap_arbol_virtual_directorio)) + "\n")
		}
		if EstaLlenoDD(avd.Avd_ap_detalle_directorio, sb.Sb_ap_detalle_directorio, sb.Sb_detalle_directorio_count, pathDisco) {
			strArray[i] = convertName(avd.Avd_nomre_directotrio[:])
			buffer.WriteString("nodo" + strconv.Itoa(i) + ":f6 -> node" + strconv.Itoa(int(avd.Avd_ap_detalle_directorio)) + "\n")
		}
		buffer.WriteString("nodo" + strconv.Itoa(i) + "[ shape=record, label =\"" + "{" + convertName(avd.Avd_nomre_directotrio[:]) + "|{<f0> " + strconv.Itoa(int(avd.Avd_ap_array_subdirectoios[0])) + "|<f1>" + strconv.Itoa(int(avd.Avd_ap_array_subdirectoios[1])) + "|<f2> " + strconv.Itoa(int(avd.Avd_ap_array_subdirectoios[2])) + "|<f3> " + strconv.Itoa(int(avd.Avd_ap_array_subdirectoios[3])) + "|<f4> " + strconv.Itoa(int(avd.Avd_ap_array_subdirectoios[4])) + "|<f5>" + strconv.Itoa(int(avd.Avd_ap_array_subdirectoios[5])) + "|<f6>" + strconv.Itoa(int(avd.Avd_ap_detalle_directorio)) + "|<f7> " + strconv.Itoa(int(avd.Avd_ap_arbol_virtual_directorio)) + "}}\"];\n")
	}
	/*
	   Graficar DD's
	*/
	f.Seek(sb.Sb_ap_detalle_directorio, 0)
	dd := DD{}
	for i := 0; i < int(sb.Sb_detalle_directorio_count); i++ {
		err = binary.Read(f, binary.BigEndian, &dd)
		if dd.Ocupado == 0 {
			break
		}

		if EstaLlenoDD(int64(i), sb.Sb_ap_detalle_directorio, sb.Sb_detalle_directorio_count, pathDisco) {
			for j := 0; j < 5; j++ {
				if convertName(dd.Dd_array_files[j].Dd_file_nombre[:]) != convertName(dos[:]) {
					buffer.WriteString("node" + strconv.Itoa(i) + ":f" + strconv.Itoa(j+1) + "->  nodex" + strconv.Itoa(int(dd.Dd_array_files[j].Dd_file_ap_inodo)) + "\n")
				}
			}
			buffer.WriteString("node" + strconv.Itoa(i) + "[shape=record, label=\"" + "{ dd " + strArray[i] + "|")
			for j := 0; j < 5; j++ {
				if convertName(dd.Dd_array_files[j].Dd_file_nombre[:]) != convertName(dos[:]) {
					buffer.WriteString("{<f" + strconv.Itoa(j) + "> " + convertName(dd.Dd_array_files[j].Dd_file_nombre[:]) + "| <f" + strconv.Itoa(j+1) + "> " + strconv.Itoa(int(dd.Dd_array_files[j].Dd_file_ap_inodo)) + "} |")
				} else {
					buffer.WriteString("{-1 | } |")
				}

			}
			if dd.Dd_ap_detalle_directorio != -1 {
				buffer.WriteString("{" + strconv.Itoa(int(dd.Dd_ap_detalle_directorio)) + " | <f10>  }}\"];\n")
				buffer.WriteString("node" + strconv.Itoa(i) + ":f10 -> " + "node" + strconv.Itoa(int(dd.Dd_ap_detalle_directorio)))
			} else {
				buffer.WriteString("{*1 | <f10>  }}\"];\n")
			}
			buffer.WriteString("\n")
		}
	}
	/*
	   Graficar Inodo's
	   X para identificarlos
	*/
	f.Seek(sb.Sb_ap_tabla_inodo, 0)
	inodo := Inodo{}
	for i := 0; i < int(sb.Sb_inodos_count); i++ {
		err = binary.Read(f, binary.BigEndian, &inodo)
		if inodo.I_count_inodo == -1 {
			break
		}
		if inodo.I_ao_indirecto != -1 {
			buffer.WriteString("nodex" + strconv.Itoa(int(inodo.I_count_inodo)) + "[shape=record, label=\"{Inodo" + strconv.Itoa(int(inodo.I_count_inodo)) + "|{" + strconv.Itoa(int(inodo.I_array_bloques[0])) + "| <f0> }|{" + strconv.Itoa(int(inodo.I_array_bloques[1])) + "| <f1> }|{" + strconv.Itoa(int(inodo.I_array_bloques[2])) + " | <f2> }|{" + strconv.Itoa(int(inodo.I_array_bloques[3])) + "| <f3> }|{" + strconv.Itoa(int(inodo.I_ao_indirecto)) + " | <f4> }}\"];" + "\n")
			buffer.WriteString("nodex" + strconv.Itoa(int(inodo.I_count_inodo)) + " :f4 ->" + "nodex" + strconv.Itoa(int(inodo.I_ao_indirecto)) + "\n")
			for h := 0; h < 4; h++ {
				if inodo.I_array_bloques[h] == -1 {
					break
				} else {
					buffer.WriteString("nodex" + strconv.Itoa(int(inodo.I_count_inodo)) + " :f" + strconv.Itoa(h) + "-> data" + strconv.Itoa(int(inodo.I_array_bloques[h])) + "\n")
				}
			}
		} else {
			buffer.WriteString("nodex" + strconv.Itoa(int(inodo.I_count_inodo)) + "[shape=record, label=\"{Inodo" + strconv.Itoa(int(inodo.I_count_inodo)) + "|{" + strconv.Itoa(int(inodo.I_array_bloques[0])) + "| <f0> }|{" + strconv.Itoa(int(inodo.I_array_bloques[1])) + "| <f1> }|{" + strconv.Itoa(int(inodo.I_array_bloques[2])) + " | <f2> }|{" + strconv.Itoa(int(inodo.I_array_bloques[3])) + "| <f3> }|{*" + strconv.Itoa(int(inodo.I_ao_indirecto)) + " | <f4> }}\"];" + "\n")
			for h := 0; h < 4; h++ {
				if inodo.I_array_bloques[h] == -1 {
					break
				} else {
					buffer.WriteString("nodex" + strconv.Itoa(int(inodo.I_count_inodo)) + " :f" + strconv.Itoa(h) + "-> data" + strconv.Itoa(int(inodo.I_array_bloques[h])) + "\n")
				}
			}
		}
	}
	/*
	   Graficar Bloque's
	*/
	f.Seek(sb.Sb_ap_bloques, 0)
	data := Bloque{}
	for i := 0; i < int(sb.Sb_bloques_count); i++ {
		err = binary.Read(f, binary.BigEndian, &data)
		if data.Db_data[0] == 0 {
			break
		}
		buffer.WriteString("data" + strconv.Itoa(i) + "[shape=record, label=\"{File block | <f1> " + convertBloqueData(data.Db_data[:]) + "}}\"];\n")

	}
	buffer.WriteString("\n}")
	var datos string
	datos = string(buffer.String())
	CreateArchivo(pathCarpeta, datos)
	fmt.Println("¡Reporte Tree creado exitosamente!")
	return false
}

func ReporteFile(idParticion string, pathCarpeta string, ruta string, ListaDiscos *list.List) bool {
	var bloquesGraficar [100]int
	carpetas_Graficar := strings.Split(ruta, "/")
	var buffer bytes.Buffer
	var noDirectorio int64 = 0
	buffer.WriteString("digraph grafica{\nrankdir=TB;\nnode [shape = record, style=filled, fillcolor=orange];\n")
	sb := SB{}
	var dos [15]byte
	avd := AVD{}
	var strArray [100]string
	pathDisco, Part_name, _ := RecorrerListaDisco(idParticion, ListaDiscos)
	sb, _ = DevolverSuperBlque(pathDisco, Part_name)
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta " + pathDisco)
		return false
	}
	defer f.Close()
	f.Seek(sb.Sb_ap_arbol_directorio, 0)
	for i := 0; i < int(sb.Sb_arbol_virtual_count); i++ {
		err = binary.Read(f, binary.BigEndian, &avd)
		if avd.Avd_nomre_directotrio == dos {
			break
		}
		if convertName(avd.Avd_nomre_directotrio[:]) == carpetas_Graficar[len(carpetas_Graficar)-2] {

			noDirectorio = avd.Avd_ap_detalle_directorio
			if EstaLlenoDD(avd.Avd_ap_detalle_directorio, sb.Sb_ap_detalle_directorio, sb.Sb_detalle_directorio_count, pathDisco) {
				strArray[i] = convertName(avd.Avd_nomre_directotrio[:])
			}
		}
	}
	/*
	   Graficar DD's
	*/
	noInodoGraficar := 0
	f.Seek(sb.Sb_ap_detalle_directorio, 0)
	dd := DD{}
	for i := 0; i < int(sb.Sb_detalle_directorio_count); i++ {
		err = binary.Read(f, binary.BigEndian, &dd)
		if dd.Ocupado == 0 {
			break
		}
		if noDirectorio == int64(i) {
			if EstaLlenoDD(int64(i), sb.Sb_ap_detalle_directorio, sb.Sb_detalle_directorio_count, pathDisco) {
				for j := 0; j < 5; j++ {
					if convertName(dd.Dd_array_files[j].Dd_file_nombre[:]) == carpetas_Graficar[len(carpetas_Graficar)-1] {
						noInodoGraficar = int(dd.Dd_array_files[j].Dd_file_ap_inodo)
					}
				}

				if dd.Dd_ap_detalle_directorio != -1 {
					noDirectorio = dd.Dd_ap_detalle_directorio
				}
			}
		}
	}
	/*
	   Graficar Inodo's
	   X para identificarlos
	*/
	f.Seek(sb.Sb_ap_tabla_inodo, 0)
	inodo := Inodo{}
	cont1 := 0
	for i := 0; i < int(sb.Sb_inodos_count); i++ {
		err = binary.Read(f, binary.BigEndian, &inodo)
		if inodo.I_count_inodo == -1 {
			break
		}
		if noInodoGraficar == i {
			if inodo.I_ao_indirecto != -1 {
				noInodoGraficar = i + 1
				for h := 0; h < 4; h++ {
					if inodo.I_array_bloques[h] == -1 {
						break
					} else {
						bloquesGraficar[cont1] = int(inodo.I_array_bloques[h])
						cont1++
					}
				}
			} else {
				for h := 0; h < 4; h++ {
					if inodo.I_array_bloques[h] == -1 {
						break
					} else {
						bloquesGraficar[cont1] = int(inodo.I_array_bloques[h])
						cont1++
					}
				}
			}
		}
	}

	contenidoFile := ""
	carpetas := strings.Split(ruta, "/")
	nameFile := carpetas[len(carpetas)-1]
	cont1 = 0
	f.Seek(sb.Sb_ap_bloques, 0)
	data := Bloque{}
	for i := 0; i < int(sb.Sb_bloques_count); i++ {
		err = binary.Read(f, binary.BigEndian, &data)
		if data.Db_data[0] == 0 {
			break
		}
		if bloquesGraficar[cont1] == i {
			contenidoFile += convertBloqueData(data.Db_data[:])
			cont1++
		}
	}

	buffer.WriteString("File " + "[shape=record, label=\"{" + nameFile + "| <f1> " + contenidoFile + "}}\"];\n")

	//Crear Archivo
	buffer.WriteString("\n}")
	var datos string
	datos = string(buffer.String())
	CreateArchivo(pathCarpeta, datos)
	fmt.Println("¡Reporte File creado exitosamente!")
	return false
}

func EstaLlenoDD(posicion int64, inicioDD int64, cantidadDD int64, pathDisco string) bool {
	estaLleno := false
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	f.Seek(inicioDD, 0)
	dd := DD{}
	for i := 0; i < int(cantidadDD); i++ {
		err = binary.Read(f, binary.BigEndian, &dd)
		if dd.Ocupado == 0 {
			break
		}
		if i == int(posicion) {
			for j := 0; j < 5; j++ {
				if len(convertName(dd.Dd_array_files[j].Dd_file_nombre[:])) > 0 {
					estaLleno = true
					break
				} else {
					estaLleno = false
				}
			}
		}
	}
	return estaLleno
}

func CreateArchivo(path string, data string) {
	f, err := os.Create(path[0:len(path)-4] + ".dot")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(data)

	if err2 != nil {
		log.Fatal(err2)
	}

	extension := path[len(path)-3:]
	executeComand("dot -T" + extension + " " + path[0:len(path)-4] + ".dot" + " -o " + path)
}

func convertName(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

func convertBloqueData(c []byte) string {
	if c[0] == 32 {
		return " "
	}
	n := -1
	for i, b := range c {
		if b == 32 || b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

//=============================== FIN -> REPORTES ===============================

//=============================== EXEC ===============================

func EjecutarComandoExec(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	fmt.Println("->Ejecutando EXEC...")
	ParamValidos = true
	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case ">path":
				fmt.Println(propiedadTemp.Val)
				dat, err := ioutil.ReadFile(propiedadTemp.Val)
				if err != nil {
					fmt.Println("*Error")
				}
				fmt.Println(string(dat))
				fmt.Println("")
				LeerTexto(string(dat), ListaDiscos)
			default:
				fmt.Println("*Error al ejecutar el comando")
			}
		}
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}

//=============================== FIN -> EXEC ===============================

//=============================== MKFS ===============================
func EjecutarComandoMKFS(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	fmt.Println("->Ejecutando MKFS...")
	ParamValidos = true
	var propiedades [4]string
	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case ">id":
				propiedades[0] = propiedadTemp.Val
			case ">type":
				propiedades[1] = propiedadTemp.Val
			case ">add":
				propiedades[2] = propiedadTemp.Val
			case ">unit":
				propiedades[3] = propiedadTemp.Val
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		ExecuteMKFS(propiedades[0], ListaDiscos)
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}
func ExecuteMKFS(id string, ListaDiscos *list.List) bool {
	dt := time.Now()
	idValido := IdValido(id, ListaDiscos)
	if idValido == false {
		fmt.Println("El id no existe")
		return false
	}
	Id := strings.ReplaceAll(id, "81", "")
	IdDisco := Id[1:]
	pathDisco := ""
	Part_name := ""
	nombreDisco := ""
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disco DISCO
		disco = element.Value.(DISCO)
		if BytesToString(disco.Id) == IdDisco {
			for i := 0; i < len(disco.Particiones); i++ {
				var mountTemp = disco.Particiones[i]
				if mountTemp.Id == id {
					copy(mountTemp.EstadoMKS[:], "1")
					Part_name = mountTemp.Part_name
					pathDisco = disco.Path
					nombreDisco = disco.NombreDisco
					break
				}
			}

		}
		element.Value = disco
	}
	mbr, sizeParticion, InicioParticion := ReturnMBR(pathDisco, Part_name)
	superBloque := SB{}
	avd := AVD{}
	dd := DD{}
	inodo := Inodo{}
	bloque := Bloque{}
	bitacora := Bitacora{}
	noEstructuras := (sizeParticion - (2 * int64(unsafe.Sizeof(superBloque)))) /
		(27 + int64(unsafe.Sizeof(avd)) + int64(unsafe.Sizeof(dd)) + (5*int64(unsafe.Sizeof(inodo)) +
			(20 * int64(unsafe.Sizeof(bloque))) + int64(unsafe.Sizeof(bitacora))))

	//NO estructuras
	var cantidadAVD int64 = noEstructuras
	var cantidadDD int64 = noEstructuras
	var cantidadInodos int64 = noEstructuras * 5
	var cantidadBloques int64 = 4 * cantidadInodos
	var Bitacoras int64 = noEstructuras
	//Bitmaps
	var InicioBitmapAVD int64 = InicioParticion + int64(unsafe.Sizeof(superBloque))
	var InicioAVD int64 = InicioBitmapAVD + cantidadAVD
	var InicioBitmapDD int64 = InicioAVD + (int64(unsafe.Sizeof(avd)) * cantidadAVD)
	var InicioDD int64 = InicioBitmapDD + cantidadDD
	var InicioBitmapInodo int64 = InicioDD + (int64(unsafe.Sizeof(dd)) * cantidadDD)
	var InicioInodo int64 = InicioBitmapInodo + cantidadInodos
	var InicioBitmapBloque int64 = InicioInodo + (int64(unsafe.Sizeof(inodo)) * cantidadInodos)
	var InicioBLoque int64 = InicioBitmapBloque + cantidadBloques
	var InicioBitacora int64 = InicioBLoque + (int64(unsafe.Sizeof(bloque)) * cantidadBloques)
	var InicioCopiaSB int64 = InicioBitacora + (int64(unsafe.Sizeof(bitacora)) * Bitacoras)

	//Inicializando SuperBloque
	copy(superBloque.Sb_nombre_hd[:], nombreDisco)
	superBloque.Sb_arbol_virtual_count = cantidadAVD
	superBloque.Sb_detalle_directorio_count = cantidadDD
	superBloque.Sb_inodos_count = cantidadInodos
	superBloque.Sb_bloques_count = cantidadBloques
	//
	superBloque.Sb_arbol_virtual_free = cantidadAVD
	superBloque.Sb_detalle_directorio_free = cantidadDD
	superBloque.Sb_inodos_free = cantidadInodos
	superBloque.Sb_bloques_free = cantidadBloques
	copy(superBloque.Sb_date_creacion[:], dt.String())
	copy(superBloque.Sb_date_ultimo_montaje[:], dt.String())
	superBloque.Sb_montajes_count = 1
	//Bitmaps
	superBloque.Sb_ap_bitmap_arbol_directorio = InicioBitmapAVD
	superBloque.Sb_ap_arbol_directorio = InicioAVD
	superBloque.Sb_ap_bitmap_detalle_directorio = InicioBitmapDD
	superBloque.Sb_ap_detalle_directorio = InicioDD
	superBloque.Sb_ap_bitmap_tabla_inodo = InicioBitmapInodo
	superBloque.Sb_ap_tabla_inodo = InicioInodo
	superBloque.Sb_ap_bitmap_bloques = InicioBitmapBloque
	superBloque.Sb_ap_bloques = InicioBLoque
	superBloque.Sb_ap_log = InicioBitacora
	superBloque.Sb_size_struct_arbol_directorio = int64(unsafe.Sizeof(avd))
	superBloque.Sb_size_struct_Detalle_directorio = int64(unsafe.Sizeof(dd))
	superBloque.Sb_size_struct_inodo = int64(unsafe.Sizeof(inodo))
	superBloque.Sb_size_struct_bloque = int64(unsafe.Sizeof(bloque))
	superBloque.Sb_first_free_bit_arbol_directorio = InicioBitmapAVD
	superBloque.Sb_first_free_bit_detalle_directoriio = InicioBitmapDD
	superBloque.Sb_dirst_free_bit_tabla_inodo = InicioBitmapInodo
	superBloque.Sb_first_free_bit_bloques = InicioBitmapBloque
	superBloque.Sb_magic_num = 201700381
	superBloque.InicioCopiaSB = InicioCopiaSB
	superBloque.ConteoAVD = 0
	superBloque.ConteoDD = 0
	superBloque.ConteoInodo = 0
	superBloque.ConteoBloque = 0
	//Escribir en Particion
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	//Escribir Super Boot
	f.Seek(InicioParticion, 0)
	err = binary.Write(f, binary.BigEndian, &superBloque)
	//Escribir Bit Map Arbol Virtual de Directorio
	f.Seek(InicioBitmapAVD, 0)
	var otro int8 = 0
	var i int64 = 0
	for i = 0; i < cantidadAVD; i++ {
		err = binary.Write(f, binary.BigEndian, &otro)
	}
	//Escribir Arbol de Directorio
	f.Seek(InicioAVD, 0)
	i = 0
	for i = 0; i < cantidadAVD; i++ {
		err = binary.Write(f, binary.BigEndian, &avd)
	}
	//Escribir Bitmap Detalle Directorio
	f.Seek(InicioBitmapDD, 0)
	i = 0
	for i = 0; i < cantidadDD; i++ {
		err = binary.Write(f, binary.BigEndian, &otro)
	}
	//Escribir Detalle Directorio
	f.Seek(InicioDD, 0)
	i = 0
	dd.Dd_ap_detalle_directorio = -1
	for i = 0; i < cantidadDD; i++ {
		err = binary.Write(f, binary.BigEndian, &dd)
	}
	//Escribir Bitmap Tabla Inodo
	f.Seek(InicioBitmapInodo, 0)
	i = 0
	for i = 0; i < cantidadInodos; i++ {
		err = binary.Write(f, binary.BigEndian, &otro)
	}
	//Escribir Tabla Inodos
	f.Seek(InicioInodo, 0)
	i = 0
	inodo.I_count_inodo = -1
	for i = 0; i < cantidadInodos; i++ {
		err = binary.Write(f, binary.BigEndian, &inodo)
	}
	//Escribir Bitmap BLoque de datos
	f.Seek(InicioBitmapBloque, 0)
	i = 0
	for i = 0; i < cantidadBloques; i++ {
		err = binary.Write(f, binary.BigEndian, &otro)
	}
	//Escribir Bloque de datos
	f.Seek(InicioBLoque, 0)
	i = 0
	copy(bloque.Db_data[:], "")
	for i = 0; i < cantidadBloques; i++ {
		err = binary.Write(f, binary.BigEndian, &bloque)
	}
	//Escribir Bitacoras
	f.Seek(InicioBitacora, 0)
	i = 0
	bitacora.Size = -1
	for i = 0; i < Bitacoras; i++ {
		err = binary.Write(f, binary.BigEndian, &bitacora)
	}
	//Escribir Copia Super Boot
	f.Seek(InicioCopiaSB, 0)
	err = binary.Write(f, binary.BigEndian, &superBloque)

	//Crear Raiz  -----> /  y  archivo con usuarios
	CrearRaiz(pathDisco, InicioParticion)
	fmt.Println(mbr.Mbr_tamano)
	fmt.Println("¡Formateo realizado correctamente!")
	return false
}

func ReturnMBR(path string, Part_name string) (MBR, int64, int64) {
	mbr := MBR{}
	var Particiones [4]Particion
	var nombre2 [15]byte
	var size int64
	copy(nombre2[:], Part_name)
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return mbr, 0, 0
	}
	defer f.Close()

	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}
	Particiones = mbr.Particiones
	for i := 0; i < 4; i++ {
		if BytesPart_name(Particiones[i].Part_name) == BytesPart_name(nombre2) {
			size = Particiones[i].Part_size
			return mbr, size, Particiones[i].Part_start
		}
	}
	for i := 0; i < 4; i++ {
		if strings.ToLower(BytesToString(Particiones[i].Part_type)) == "e" {
			var InicioExtendida int64 = Particiones[i].Part_start
			f.Seek(InicioExtendida, 0)
			ebr := EBR{}
			err = binary.Read(f, binary.BigEndian, &ebr)
			if ebr.Part_next == -1 {
				fmt.Println("No Hay particiones Logicas")
			} else {
				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				for {
					if ebr.Part_next == -1 {
						break
					} else {
						f.Seek(ebr.Part_next, 0)
						err = binary.Read(f, binary.BigEndian, &ebr)
					}
					if BytesPart_name(ebr.Part_name) == BytesPart_name(nombre2) {
						fmt.Println("Logica Encontrada")
						return mbr, ebr.Part_size, ebr.Part_start
					}

				}
			}
		}
	}
	return mbr, 0, 0
}

func CrearRaiz(pathDisco string, InicioParticion int64) bool {
	dt := time.Now()
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	f.Seek(InicioParticion, 0)
	sb := SB{}
	err = binary.Read(f, binary.BigEndian, &sb)
	/*
		Escribir 1 en bitmap avd y escribir avd
	*/
	f.Seek(sb.Sb_ap_bitmap_arbol_directorio, 0)
	var otro int8 = 0
	otro = 1
	err = binary.Write(f, binary.BigEndian, &otro)
	bitLibre, _ := f.Seek(0, os.SEEK_CUR)
	sb.Sb_first_free_bit_arbol_directorio = bitLibre
	avd := AVD{}
	copy(avd.Avd_fecha_creacion[:], dt.String())
	copy(avd.Avd_nomre_directotrio[:], "/")
	for j := 0; j < 6; j++ {
		avd.Avd_ap_array_subdirectoios[j] = -1
	}
	avd.Avd_ap_detalle_directorio = 0
	avd.Avd_ap_arbol_virtual_directorio = -1
	copy(avd.Avd_proper[:], "root")
	f.Seek(sb.Sb_ap_arbol_directorio, 0)
	err = binary.Write(f, binary.BigEndian, &avd)

	sb.Sb_arbol_virtual_free = sb.Sb_arbol_virtual_free - 1
	/*
		Escribir 1 en bitmap detalleDirectorio y escribir detalleDirectorio
	*/
	f.Seek(sb.Sb_ap_bitmap_detalle_directorio, 0)
	otro = 1
	err = binary.Write(f, binary.BigEndian, &otro)
	otro = 0
	bitLibre, _ = f.Seek(0, os.SEEK_CUR)
	sb.Sb_first_free_bit_detalle_directoriio = bitLibre
	detalleDirectorio := DD{}
	arregloDD := ArregloDD{}
	copy(arregloDD.Dd_file_nombre[:], "users.txt")
	copy(arregloDD.Dd_file_date_creacion[:], dt.String())
	copy(arregloDD.Dd_file_date_modificacion[:], dt.String())
	arregloDD.Dd_file_ap_inodo = 0
	detalleDirectorio.Dd_array_files[0] = arregloDD
	detalleDirectorio.Ocupado = 1
	for j := 0; j < 5; j++ {
		if j == 0 {
			detalleDirectorio.Dd_array_files[j].Dd_file_ap_inodo = 0
		} else {
			detalleDirectorio.Dd_array_files[j].Dd_file_ap_inodo = -1
		}
	}
	detalleDirectorio.Dd_ap_detalle_directorio = -1
	f.Seek(sb.Sb_ap_detalle_directorio, 0)
	err = binary.Write(f, binary.BigEndian, &detalleDirectorio)

	sb.Sb_detalle_directorio_free = sb.Sb_detalle_directorio_free - 1
	/*
		Escribir 1 en bitmap tablaInodo y escribir Inodo
	*/
	//var cantidadBloque int64 = CantidadBloqueUsar("1,G,root\n1,U,root,root,201902308\n")
	var cantidadBloque int64 = CantidadBloqueUsar("1,G,root\n1,U,root,123\n")
	f.Seek(sb.Sb_ap_bitmap_tabla_inodo, 0)
	otro = 1
	err = binary.Write(f, binary.BigEndian, &otro)
	otro = 0
	bitLibre, _ = f.Seek(0, os.SEEK_CUR)
	sb.Sb_dirst_free_bit_tabla_inodo = bitLibre
	inodo := Inodo{}
	for j := 0; j < 4; j++ {
		inodo.I_array_bloques[j] = -1
	}
	inodo.I_count_inodo = 0
	inodo.I_size_archivo = 10
	inodo.I_count_bloques_asignados = cantidadBloque
	for h := 0; h < int(cantidadBloque); h++ {
		inodo.I_array_bloques[h] = int64(h)
	}
	inodo.I_ao_indirecto = -1
	inodo.I_id_proper = 201902308
	f.Seek(sb.Sb_ap_tabla_inodo, 0)
	err = binary.Write(f, binary.BigEndian, &inodo)
	sb.Sb_inodos_free = sb.Sb_inodos_free - 1
	/*
		Escribir 1 en bitmap bloqueDatos y escribir el bloque datos
	*/
	f.Seek(sb.Sb_ap_bitmap_bloques, 0)
	otro = 1
	for k := 0; k < int(cantidadBloque); k++ {
		err = binary.Write(f, binary.BigEndian, &otro)
	}
	otro = 0
	bitLibre, _ = f.Seek(0, os.SEEK_CUR)
	sb.Sb_first_free_bit_bloques = bitLibre
	f.Seek(sb.Sb_ap_bloques, 0)
	usesTxt := []byte("1,G,root\n1,U,root,123\n")
	for k := 0; k < int(cantidadBloque); k++ {
		if k == 0 {
			bloque := Bloque{}
			copy(bloque.Db_data[:], string([]byte(usesTxt)))
			err = binary.Write(f, binary.BigEndian, &bloque)
		} /*else {
			bloque := Bloque{}
			copy(bloque.Db_data[:], string([]byte(usesTxt[k*25:len(usesTxt)])))
			err = binary.Write(f, binary.BigEndian, &bloque)
		}*/
		sb.Sb_bloques_free = sb.Sb_bloques_free - 1
		sb.ConteoBloque = sb.ConteoBloque + int64(k)
	}
	/*
		Actualizar SB
	*/
	f.Seek(0, 0)
	f.Seek(InicioParticion, 0)
	err = binary.Write(f, binary.BigEndian, &sb)
	return false
}
func CantidadBloqueUsar(data string) int64 {
	var noBloque int64 = 0
	cont := 1
	var dataX []byte = []byte(data)
	for i := 0; i < len(dataX); i++ {
		if cont == 25 {
			noBloque = noBloque + 1
			cont = 0
		}
		cont++
	}
	if len(dataX)%25 != 0 {
		noBloque = noBloque + 1
	}
	return noBloque
}
func CantidadInodosUsar(data string) int64 {
	var noBloque int64 = 0
	cont := 0
	var dataX []byte = []byte(data)
	for i := 0; i < len(dataX); i++ {
		if cont == 25 {
			noBloque = noBloque + 1
			cont = 0
		}
		cont++
	}
	if len(dataX)%5 != 0 {
		noBloque = noBloque + 1
	}
	return noBloque
}
func DevolverSuperBlque(path string, Part_name string) (SB, int64) {
	mbr := MBR{}
	sb := SB{}
	var Particiones [4]Particion
	var nombre2 [15]byte
	copy(nombre2[:], Part_name)
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return sb, 0
	}
	defer f.Close()

	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}
	Particiones = mbr.Particiones
	for i := 0; i < 4; i++ {
		if BytesPart_name(Particiones[i].Part_name) == BytesPart_name(nombre2) {
			f.Seek(Particiones[i].Part_start, 0)
			err = binary.Read(f, binary.BigEndian, &sb)
			return sb, Particiones[i].Part_start
		}
	}
	for i := 0; i < 4; i++ {
		if strings.ToLower(BytesToString(Particiones[i].Part_type)) == "e" {
			var InicioExtendida int64 = Particiones[i].Part_start
			f.Seek(InicioExtendida, 0)
			ebr := EBR{}
			err = binary.Read(f, binary.BigEndian, &ebr)
			if ebr.Part_next == -1 {
				fmt.Println("No Hay particiones Logicas")
			} else {
				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				for {
					if ebr.Part_next == -1 {
						break
					} else {
						f.Seek(ebr.Part_next, 0)
						err = binary.Read(f, binary.BigEndian, &ebr)
					}
					if BytesPart_name(ebr.Part_name) == BytesPart_name(nombre2) {
						fmt.Println("Logica Encontrada")
						f.Seek(ebr.Part_start, 0)
						err = binary.Read(f, binary.BigEndian, &sb)
						return sb, ebr.Part_start
					}

				}
			}
		}
	}
	return sb, 0
}

//=============================== FIN -> MKFS ===============================

//=============================== LOGIN ===============================
func EjecutarComandoLogin(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (bool, string) {
	fmt.Println("->Ejecutando Login...")
	ParamValidos := true
	usuario := ""
	var propiedades [3]string
	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case ">user":
				propiedades[0] = propiedadTemp.Val
			case ">pwd":
				propiedades[1] = string(propiedadTemp.Val)
			case ">id":
				propiedades[2] = propiedadTemp.Val
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		ParamValidos, usuario = ExecuteLogin(propiedades[0], propiedades[1], propiedades[2], ListaDiscos)

		return ParamValidos, usuario
	} else {
		ParamValidos = false
		return ParamValidos, usuario
	}
}
func ExecuteLogin(usuario string, password string, id string, ListaDiscos *list.List) (bool, string) {
	idValido := IdValido(id, ListaDiscos)
	if idValido == false {
		fmt.Println("El id no existe, la particion no esta montada")
		return false, ""
	} else if global != "" {
		fmt.Println("Ya hay una sesion iniciada")
		return false, ""
	}
	pathDisco, Part_name, nombreDisco := RecorrerListaDisco(id, ListaDiscos)
	mbr, sizeParticion, InicioParticion := ReturnMBR(pathDisco, Part_name)
	superBloque := SB{}
	f, err := os.OpenFile(pathDisco, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta " + pathDisco)
		return false, ""
	}
	defer f.Close()
	f.Seek(InicioParticion, 0)
	err = binary.Read(f, binary.BigEndian, &superBloque)

	//Obtener avd raiz
	avd := AVD{}
	dd := DD{}
	inodo := Inodo{}
	bloque := Bloque{}
	f.Seek(superBloque.Sb_ap_arbol_directorio, 0)
	err = binary.Read(f, binary.BigEndian, &avd)
	apuntadorDD := avd.Avd_ap_detalle_directorio
	f.Seek(superBloque.Sb_ap_detalle_directorio, 0)
	for i := 0; i < int(superBloque.Sb_arbol_virtual_free); i++ {
		err = binary.Read(f, binary.BigEndian, &dd)
		if i == int(apuntadorDD) {
			break
		}
	}
	apuntadorInodo := dd.Dd_array_files[0].Dd_file_ap_inodo
	f.Seek(superBloque.Sb_ap_tabla_inodo, 0)
	for i := 0; i < int(superBloque.Sb_inodos_free); i++ {
		err = binary.Read(f, binary.BigEndian, &inodo)
		if i == int(apuntadorInodo) {
			break
		}
	}
	var userstxt string = ""
	//Leer Users.txt
	posicion := 0
	f.Seek(superBloque.Sb_ap_bloques, 0)
	for i := 0; i < int(superBloque.Sb_inodos_free); i++ {
		err = binary.Read(f, binary.BigEndian, &bloque)

		if int(inodo.I_array_bloques[posicion]) != -1 && int(inodo.I_array_bloques[posicion]) == i {
			userstxt += ConvertData(bloque.Db_data)
		} else if int(inodo.I_array_bloques[posicion]) == -1 {
			break
		} else {
			break
		}
		if posicion < 4 {
			posicion++
		} else if posicion == 4 {
			posicion = 0
		}
	}
	lineaUsuarioTxt := strings.Split(userstxt, "\n")
	for i := 0; i < len(lineaUsuarioTxt); i++ {
		fmt.Println("hola", len(lineaUsuarioTxt))
		if "root" == usuario && "123" == password {
			fmt.Println("¡Sesion iniciada con exito!")
			globalIdLogin = id
			//fmt.Println(globalIdLogin)
			return true, usuario
		} else {
			fmt.Printf("Error en usuario o password")
			fmt.Println(nombreDisco, mbr.Mbr_tamano, sizeParticion)
		}

	}
	return false, ""
}

//=============================== FIN -> LOGIN ===============================

//=============================== MKDIR ===============================
func EjecutarComandoMKDIR(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	fmt.Println("->Ejecutando MKDIR...")
	if globalIdLogin == "" {
		fmt.Println("Primero debe iniciar sesion")
		return false
	}
	ParamValidos = true
	var propiedades [2]string

	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			/*case "-id":
			propiedades[0] = propiedadTemp.Val*/
			case ">path":
				propiedades[0] = propiedadTemp.Val
			case ">r":
				propiedades[1] = propiedadTemp.Val
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}

		//ExecuteMKDIR(propiedades[0], propiedades[1], propiedades[2], ListaDiscos)
		ExecuteMKDIR(globalIdLogin, propiedades[0], propiedades[1], ListaDiscos)
		return ParamValidos

	} else {
		ParamValidos = false
		return ParamValidos
	}
}
func ExecuteMKDIR(id string, path string, p string, ListaDiscos *list.List) bool {

	//Si no existen las carpetas se crean

	//Escribir en bitacora
	dt := time.Now()
	sb := SB{}
	pathDisco, Part_name, _ := RecorrerListaDisco(id, ListaDiscos)
	sb, _ = DevolverSuperBlque(pathDisco, Part_name)
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)

	}
	defer f.Close()
	bitacora := Bitacora{}
	copy(bitacora.Log_tipo_operacion[:], "mkdir")
	copy(bitacora.Log_tipo[:], "0")
	copy(bitacora.Log_nombre[:], path)
	copy(bitacora.Log_Contenido[:], "")
	copy(bitacora.Log_fecha[:], dt.String())
	bitacora.Size = 1
	bitacoraTemp := Bitacora{}
	var bitBitacora int64 = 0
	f.Seek(sb.Sb_ap_log, 0)
	for i := 0; i < 3000; i++ {
		bitBitacora, _ = f.Seek(0, os.SEEK_CUR)
		err = binary.Read(f, binary.BigEndian, &bitacoraTemp)
		if bitacoraTemp.Size == -1 {
			f.Seek(bitBitacora, 0)
			err = binary.Write(f, binary.BigEndian, &bitacora)
			break
		}
	}
	/*
	   Ejecutando MKDIR
	*/

	RecorrePath(path, Part_name, pathDisco)
	fmt.Println("¡Directorio creado correctamente!")
	return true
}
func RecorrePath(path string, Part_name string, pathDisco string) {
	/*
		Quitar las comillas al path si tiene
	*/
	EsComilla := path[0:1]
	if EsComilla == "\"" {
		path = path[1 : len(path)-1]
	}
	//Ver si hay mas de una carpeta
	if strings.Contains(path, "/") {
		carpetas := strings.Split(path, "/")
		if len(carpetas) == 2 {
			if ExisteCarpeta(pathDisco, Part_name, carpetas[1]) == false {
				otroAvd, _ := ModificarCarpeta(pathDisco, Part_name, "/", "")
				if otroAvd == true {
					ModificarCarpeta(pathDisco, Part_name, "/", "/")
					CrearCarpeta(pathDisco, Part_name, carpetas[1])
				} else {
					if ExisteCarpeta(pathDisco, Part_name, carpetas[1]) == false {
						CrearCarpeta(pathDisco, Part_name, carpetas[1])
					}
				}
			}
		} else {
			for i := 1; i < len(carpetas); i++ {
				if ExisteCarpeta(pathDisco, Part_name, carpetas[i]) == false {
					if carpetas[i-1] == "" {
						carpetas[i-1] = "/"
					}
					otroAvd, _ := ModificarCarpeta(pathDisco, Part_name, carpetas[i-1], "")
					if otroAvd == true {
						ModificarCarpeta(pathDisco, Part_name, carpetas[i-1], carpetas[i-1])
						CrearCarpeta(pathDisco, Part_name, carpetas[i])
					} else {
						CrearCarpeta(pathDisco, Part_name, carpetas[i])
					}
				} else {
				}
			}
		}
	}
}
func ExisteCarpeta(pathDisco string, Part_name string, carpetaBuscar string) bool {
	sb := SB{}
	var nombre2 [15]byte
	copy(nombre2[:], carpetaBuscar)
	avd := AVD{}
	sb, _ = DevolverSuperBlque(pathDisco, Part_name)
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	f.Seek(sb.Sb_ap_arbol_directorio, 0)
	for i := 0; i < int(sb.Sb_arbol_virtual_count); i++ {
		err = binary.Read(f, binary.BigEndian, &avd)
		if BytesPart_name(avd.Avd_nomre_directotrio) == BytesPart_name(nombre2) {
			return true
		}
	}
	return false
}

/*
	Funcion para modifica Puntero de avd
*/
func ModificarCarpeta(pathDisco string, Part_name string, carpetaModificar string, nombreOpcional string) (bool, int64) {
	puntero_avd := true
	sb := SB{}
	avd := AVD{}
	var nombre2 [15]byte
	copy(nombre2[:], carpetaModificar)
	var bitLibre int64
	sb, _ = DevolverSuperBlque(pathDisco, Part_name)
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false, 0
	}
	defer f.Close()
	f.Seek(sb.Sb_ap_arbol_directorio, 0)
	bitLibre = sb.Sb_ap_arbol_directorio
	for i := 0; i < int(sb.Sb_arbol_virtual_count); i++ {
		err = binary.Read(f, binary.BigEndian, &avd)
		if BytesPart_name(avd.Avd_nomre_directotrio) == BytesPart_name(nombre2) {
			if avd.Avd_ap_arbol_virtual_directorio != -1 {
				bitLibre, _ = f.Seek(0, os.SEEK_CUR)
				continue
			}
			for i := 0; i < len(avd.Avd_ap_array_subdirectoios); i++ {
				if avd.Avd_ap_array_subdirectoios[i] == -1 {
					avd.Avd_ap_array_subdirectoios[i] = sb.ConteoAVD + 1
					puntero_avd = false
					break
				}
			}
			if puntero_avd != true {
				f.Seek(bitLibre, 0)
				err = binary.Write(f, binary.BigEndian, &avd)
				bitLibre = 0
				break
			} else {
				if estaLlenoAVD(pathDisco, Part_name, carpetaModificar) == false {
					avd.Avd_ap_arbol_virtual_directorio = sb.ConteoAVD + 1
					f.Seek(bitLibre, 0)
					err = binary.Write(f, binary.BigEndian, &avd)
					bitLibre = 0
					CrearCarpeta(pathDisco, Part_name, carpetaModificar)
					return true, avd.Avd_ap_arbol_virtual_directorio
				}
				break
			}
		}
		bitLibre, _ = f.Seek(0, os.SEEK_CUR)
	}
	return false, 0

}
func estaLlenoAVD(pathDisco string, Part_name string, carpeta string) bool {
	sb := SB{}
	avd := AVD{}
	estaLleno := true
	var nombre2 [15]byte
	copy(nombre2[:], carpeta)
	sb, _ = DevolverSuperBlque(pathDisco, Part_name)
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	f.Seek(sb.Sb_ap_arbol_directorio, 0)
	for i := 0; i < int(sb.Sb_arbol_virtual_count); i++ {
		err = binary.Read(f, binary.BigEndian, &avd)
		if BytesPart_name(avd.Avd_nomre_directotrio) == BytesPart_name(nombre2) {
			if avd.Avd_ap_array_subdirectoios[5] == -1 {
				estaLleno = true
			} else if avd.Avd_ap_array_subdirectoios[5] != -1 {
				estaLleno = false
			}
		}
	}
	return estaLleno
}
func CrearCarpeta(pathDisco string, Part_name string, carpetaHija string) bool {
	dt := time.Now()
	var nombre2 [15]byte
	copy(nombre2[:], "")
	sb := SB{}
	avd := AVD{}
	var InicioParticion int64
	sb, InicioParticion = DevolverSuperBlque(pathDisco, Part_name)
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	var bitLibre int64 = 0
	var bitLibreDD int64 = 0
	f.Seek(sb.Sb_ap_arbol_directorio, 0)
	for i := 0; i < int(sb.Sb_arbol_virtual_count); i++ {
		err = binary.Read(f, binary.BigEndian, &avd)
		if BytesPart_name(avd.Avd_nomre_directotrio) == BytesPart_name(nombre2) {
			avdTemp := AVD{}
			copy(avdTemp.Avd_fecha_creacion[:], dt.String())
			copy(avdTemp.Avd_nomre_directotrio[:], carpetaHija)
			for j := 0; j < 6; j++ {
				avdTemp.Avd_ap_array_subdirectoios[j] = -1
			}
			avdTemp.Avd_ap_detalle_directorio = sb.ConteoDD + 1
			avdTemp.Avd_ap_arbol_virtual_directorio = -1
			copy(avdTemp.Avd_proper[:], global)
			f.Seek(bitLibre, 0)
			/*
				Escribir AVD
			*/
			err = binary.Write(f, binary.BigEndian, &avdTemp)
			sb.Sb_arbol_virtual_free = sb.Sb_arbol_virtual_free - 1
			sb.ConteoAVD = sb.ConteoAVD + 1
			sb.ConteoDD = sb.ConteoDD + 1
			/*
				Marcar en bitmap
			*/
			f.Seek(sb.Sb_first_free_bit_arbol_directorio, 0)
			var otro int8 = 0
			otro = 1
			err = binary.Write(f, binary.BigEndian, &otro)
			bitLibre, _ := f.Seek(0, os.SEEK_CUR)
			sb.Sb_first_free_bit_arbol_directorio = bitLibre
			/*
				Escribir DD y marcar en bitmap
			*/
			f.Seek(sb.Sb_first_free_bit_detalle_directoriio, 0)
			otro = 1
			err = binary.Write(f, binary.BigEndian, &otro)
			otro = 0
			bitLibre, _ = f.Seek(0, os.SEEK_CUR)
			sb.Sb_first_free_bit_detalle_directoriio = bitLibre
			detalleDirectorio := DD{}
			f.Seek(sb.Sb_ap_detalle_directorio, 0)
			for i := 0; i < int(sb.Sb_detalle_directorio_count); i++ {
				err = binary.Read(f, binary.BigEndian, &detalleDirectorio)
				if detalleDirectorio.Ocupado == 0 {
					detalleDirectorioTemp := DD{}
					arregloDD := ArregloDD{}
					arregloDD.Dd_file_ap_inodo = -1
					for j := 0; j < 5; j++ {
						detalleDirectorioTemp.Dd_array_files[j] = arregloDD
					}
					detalleDirectorioTemp.Ocupado = 1
					detalleDirectorioTemp.Dd_ap_detalle_directorio = -1
					f.Seek(bitLibreDD, 0)
					err = binary.Write(f, binary.BigEndian, &detalleDirectorioTemp)
					sb.Sb_detalle_directorio_free = sb.Sb_detalle_directorio_free - 1
					bitLibreDD = 0
					break
				}
				bitLibreDD, _ = f.Seek(0, os.SEEK_CUR)
			}
			/*
				Actualizar SB
			*/
			f.Seek(InicioParticion, 0)
			err = binary.Write(f, binary.BigEndian, &sb)
			bitLibre = 0
			break
		}
		bitLibre, _ = f.Seek(0, os.SEEK_CUR)
	}

	return false
}

//=============================== FIN -> MKDIR ===============================

//=============================== MKFILE ===============================

func EjecutarComandoMKFILE(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	fmt.Println("->Ejecutando MKFILE...")
	if globalIdLogin == "" {
		fmt.Println("Primero debe iniciar sesion")
		return false
	}
	ParamValidos = true
	var propiedades [3]string
	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case ">path":
				propiedades[0] = propiedadTemp.Val
			case ">r":
				propiedades[1] = propiedadTemp.Val
			case ">size":
				propiedades[2] = propiedadTemp.Val
			case ">count":
				//propiedades[2] = propiedadTemp.Val
			default:
				fmt.Println(propiedadTemp)
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		size, _ := strconv.Atoi(propiedades[2])
		ExecuteMKFILE(globalIdLogin, propiedades[0], propiedades[1], size, ListaDiscos)
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}

func ExecuteMKFILE(idParticion string, pathArchivo string, _p string, size int, ListaDiscos *list.List) bool {

	contenido := ""
	contContenido := 0
	for i := 0; i < size; i++ {
		contenido = contenido + strconv.Itoa(contContenido)
		contContenido++
		if contContenido == 10 {
			contContenido = 0
		}
	}

	//Quitar las comillas al path si tiene
	EsComilla := pathArchivo[0:1]
	if EsComilla == "\"" {
		pathArchivo = pathArchivo[1 : len(pathArchivo)-1]
	}
	pathDisco, Part_name, _ := RecorrerListaDisco(idParticion, ListaDiscos)
	CrearArchivo(pathDisco, Part_name, pathArchivo, _p, size, contenido, -1)
	fmt.Println("¡Archivo creado correctamente!")
	return true

}
func CrearArchivo(pathDisco string, Part_name string, pathArchivo string, _p string, size int, contenido string, siguienteDD int) bool {
	/*
		Obtener el SB de la particion
	*/
	otroDD := true
	dt := time.Now()
	avd := AVD{}
	sb := SB{}
	encontrado := false
	dd := DD{}
	var InicioParticion int64 = 0
	var nombreArchivo = ""
	var carpetaPadre = ""
	if strings.Contains(pathArchivo, "/") {
		nuevaPath := "/"
		carpetas := strings.Split(pathArchivo, "/")
		nombreArchivo = carpetas[len(carpetas)-1]
		carpetaPadre = carpetas[len(carpetas)-2]
		for i := 1; i < len(carpetas)-1; i++ {
			nuevaPath += carpetas[i] + "/"
		}
		nuevaPath = nuevaPath[0 : len(nuevaPath)-1]
		//Se crean las carpetas si no estan creadas
		RecorrePath(nuevaPath, Part_name, pathDisco)
	}
	sb, InicioParticion = DevolverSuperBlque(pathDisco, Part_name)
	/*
		1.Buscar AVD
		2.Buscar DD
		3.Modificar DD
		4.Crear Inodo
		5.Crear BLoque
		6.Modificar AVD
	*/
	//Obtener AVD
	var nombre2 [15]byte
	var bitLibre int64 = 0
	copy(nombre2[:], carpetaPadre)
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	/*
		Escribit en bitacora
	*/
	if siguienteDD == -1 {
		f.Seek(sb.Sb_ap_log, 0)
		bitacora := Bitacora{}
		copy(bitacora.Log_tipo_operacion[:], "mkfile")
		copy(bitacora.Log_tipo[:], "1")
		copy(bitacora.Log_nombre[:], pathArchivo)
		copy(bitacora.Log_Contenido[:], contenido)
		copy(bitacora.Log_fecha[:], dt.String())
		bitacora.Size = int64(size)
		bitacoraTemp := Bitacora{}
		var bitBitacora int64 = 0
		for i := 0; i < 3000; i++ {
			bitBitacora, _ = f.Seek(0, os.SEEK_CUR)
			err = binary.Read(f, binary.BigEndian, &bitacoraTemp)
			if bitacoraTemp.Size == -1 {
				f.Seek(bitBitacora, 0)
				err = binary.Write(f, binary.BigEndian, &bitacora)
				break
			}
		}
	}
	//EScribir Arbol Directorio

	f.Seek(sb.Sb_ap_arbol_directorio, 0)
	for i := 0; i < int(sb.Sb_arbol_virtual_count); i++ {
		err = binary.Read(f, binary.BigEndian, &avd)
		if BytesPart_name(avd.Avd_nomre_directotrio) == BytesPart_name(nombre2) {
			//Leer DD
			f.Seek(sb.Sb_ap_detalle_directorio, 0)
			for i := 0; i < 20; i++ {
				err = binary.Read(f, binary.BigEndian, &dd)
				if dd.Dd_ap_detalle_directorio != -1 && dd.Dd_ap_detalle_directorio != 0 {
					siguienteDD = int(dd.Dd_ap_detalle_directorio)
					bitLibre, _ = f.Seek(0, os.SEEK_CUR)
					continue
				}
				if i == int(avd.Avd_ap_detalle_directorio) || i == siguienteDD {
					for i := 0; i < len(dd.Dd_array_files); i++ {
						if dd.Dd_array_files[i].Dd_file_ap_inodo == -1 {
							copy(dd.Dd_array_files[i].Dd_file_nombre[:], nombreArchivo)
							dd.Dd_array_files[i].Dd_file_ap_inodo = sb.ConteoInodo + 1
							copy(dd.Dd_array_files[i].Dd_file_date_creacion[:], dt.String())
							copy(dd.Dd_array_files[i].Dd_file_date_modificacion[:], dt.String())
							otroDD = false
							break
						}
					}
					if otroDD == false {
						f.Seek(bitLibre, 0)
						err = binary.Write(f, binary.BigEndian, &dd)
						bitLibre = 0
						encontrado = true
						EscribirInodo(pathDisco, sb, contenido, InicioParticion)
						break
					} else if otroDD == true {
						//Apuntador del dd anteriro al actual
						f.Seek(bitLibre, 0)
						dd.Dd_ap_detalle_directorio = sb.ConteoDD + 1
						err = binary.Write(f, binary.BigEndian, &dd)
						bitLibre = 0
						f.Seek(0, 0)
						//Crear otro Detalle de directorio
						nuevoDD := DD{}
						sb.ConteoDD = sb.ConteoDD + 1
						nuevoDD.Ocupado = 1
						//Marcar 1 bitmap DD
						sb.Sb_detalle_directorio_free = sb.Sb_detalle_directorio_free - 1
						f.Seek(sb.Sb_first_free_bit_detalle_directoriio, 0)
						var otro int8 = 1
						err = binary.Write(f, binary.BigEndian, &otro)
						otro = 0
						bitLibre, _ = f.Seek(0, os.SEEK_CUR)
						sb.Sb_first_free_bit_detalle_directoriio = bitLibre
						//Actualizar SB
						f.Seek(InicioParticion, 0)
						err = binary.Write(f, binary.BigEndian, &sb)
						EscribirDD(sb.Sb_ap_detalle_directorio, pathDisco, sb.Sb_detalle_directorio_count, nuevoDD)
						f.Seek(0, 0)
						CrearArchivo(pathDisco, Part_name, pathArchivo, _p, size, contenido, int(dd.Dd_ap_detalle_directorio))
						encontrado = true
						break
					}
				}
				bitLibre, _ = f.Seek(0, os.SEEK_CUR)
			}
		}
		if encontrado == true {
			break
		}
	}
	return false
}
func EscribirDD(InicioDD int64, pathDisco string, cantidadDD int64, ddNuevo DD) bool {
	dd := DD{}
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	f.Seek(InicioDD, 0)
	var bitLibre int64
	for i := 0; i < int(cantidadDD); i++ {
		err = binary.Read(f, binary.BigEndian, &dd)
		if dd.Ocupado == 0 {
			f.Seek(bitLibre, 0)
			for h := 0; h < 5; h++ {
				ddNuevo.Dd_array_files[h].Dd_file_ap_inodo = -1
			}
			ddNuevo.Dd_ap_detalle_directorio = -1
			err = binary.Write(f, binary.BigEndian, &ddNuevo)
			break
		}
		bitLibre, _ = f.Seek(0, os.SEEK_CUR)
	}
	return false
}
func EscribirInodo(pathDisco string, sb SB, contenido string, InicioParticion int64) bool {
	var otro int8 = 0
	var bitLibre int64 = 0
	var restoBloque int64 = 0
	contenido2 := ""
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	//Escribir bitmap de Inodo
	var cantidadBloque int64 = CantidadBloqueUsar(contenido)
	if cantidadBloque >= 5 {
		restoBloque = cantidadBloque - 4
		cantidadBloque = 4
		contenido2 = contenido[100:len(contenido)]
		contenido = contenido[0:100]
	}
	f.Seek(sb.Sb_dirst_free_bit_tabla_inodo, 0)
	otro = 1
	err = binary.Write(f, binary.BigEndian, &otro)
	otro = 0
	bitLibre, _ = f.Seek(0, os.SEEK_CUR)
	sb.Sb_dirst_free_bit_tabla_inodo = bitLibre
	inodo := Inodo{}
	for j := 0; j < 4; j++ {
		inodo.I_array_bloques[j] = -1
	}
	inodo.I_count_inodo = sb.ConteoInodo + 1
	inodo.I_size_archivo = 10
	inodo.I_count_bloques_asignados = cantidadBloque
	inodo.I_ao_indirecto = -1
	inodo.I_id_proper = 201902308

	inodo, sb.ConteoBloque, sb.Sb_first_free_bit_bloques, sb.Sb_bloques_free = EscribirBloque(sb, cantidadBloque, pathDisco, InicioParticion, inodo, contenido)
	f.Seek(sb.Sb_ap_tabla_inodo, 0)
	inodoTemp := Inodo{}
	sb.ConteoInodo = sb.ConteoInodo + 1
	sb.Sb_inodos_free = sb.Sb_inodos_free - 1
	if restoBloque != 0 {
		inodo.I_ao_indirecto = sb.ConteoInodo + 1
	}
	for i := 0; i < int(sb.Sb_inodos_count); i++ {
		err = binary.Read(f, binary.BigEndian, &inodoTemp)
		if inodoTemp.I_count_inodo == -1 {
			f.Seek(bitLibre, 0)
			err = binary.Write(f, binary.BigEndian, &inodo)
			break
		}
		bitLibre, _ = f.Seek(0, os.SEEK_CUR)
	}
	/*
		Actualizar SB
	*/
	f.Seek(InicioParticion, 0)
	err = binary.Write(f, binary.BigEndian, &sb)
	if restoBloque != 0 {
		EscribirInodo(pathDisco, sb, contenido2, InicioParticion)
	}
	return false
}

func EscribirBloque(sb SB, cantidadBloque int64, pathDisco string, InicioParticion int64, inodo Inodo, contenido string) (Inodo, int64, int64, int64) {
	var contenido2 [64]byte
	copy(contenido2[:], contenido)
	bloqueTemp := Bloque{}
	var bitLibre_BLoque int64 = 0
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return inodo, 0, 0, 0
	}
	defer f.Close()

	//EScribir en Bitmap en bloque de datos

	var otro int8 = 0
	f.Seek(sb.Sb_first_free_bit_bloques, 0)
	otro = 1
	for k := 0; k < int(cantidadBloque); k++ {
		err = binary.Write(f, binary.BigEndian, &otro)
	}
	otro = 0
	bitLibre, _ := f.Seek(0, os.SEEK_CUR)
	sb.Sb_first_free_bit_bloques = bitLibre

	//Escribir BLoques de datos

	f.Seek(sb.Sb_ap_bloques, 0)
	for i := 0; i < int(sb.Sb_bloques_count); i++ {
		err = binary.Read(f, binary.BigEndian, &bloqueTemp)
		if bloqueTemp.Db_data[0] == 0 {
			f.Seek(bitLibre_BLoque, 0)
			for h := 0; h < int(cantidadBloque); h++ {
				inodo.I_array_bloques[h] = sb.ConteoBloque + 1
				//EScribir BLoque
				if h == 0 {
					bloque := Bloque{}
					if len(contenido) >= 25 {
						copy(bloque.Db_data[:], string([]byte(contenido[0:25])))
						err = binary.Write(f, binary.BigEndian, &bloque)
					} else {
						bloque.Db_data = contenido2
						err = binary.Write(f, binary.BigEndian, &bloque)
					}
				} else {
					bloque := Bloque{}
					copy(bloque.Db_data[:], string([]byte(contenido[h*25:len(contenido)])))
					err = binary.Write(f, binary.BigEndian, &bloque)
				}
				sb.Sb_bloques_free = sb.Sb_bloques_free - 1
				sb.ConteoBloque = sb.ConteoBloque + 1
			}
			break
		}
		bitLibre_BLoque, _ = f.Seek(0, os.SEEK_CUR)
	}
	return inodo, sb.ConteoBloque, sb.Sb_first_free_bit_bloques, sb.Sb_bloques_free
}

//=============================== FIN -> MKFILE ===============================
