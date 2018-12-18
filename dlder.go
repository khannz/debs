package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type MainMap struct {
	sync.Mutex
	M map[string]uint
}

func main() {
	var i uint
	pack := "docker.io" //From console key
	packagesMap := &MainMap{M: make(map[string]uint)}

	Put(pack, i, packagesMap) // put "main" package
	i++
	//
	RecurseDependens(packagesMap, i) // Get dependens for all packages (exclude some os pack)
	//
	folderName := "packages_for_" + pack //make new folder end move to it
	os.Mkdir(folderName, 0700)
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	packagesFullPath := exPath + "/" + folderName
	err := os.Chdir(packagesFullPath)
	//
	if err != nil {
		log.Fatal("Cant enter directory: ", packagesFullPath)
	}
	///
	for keyPack := range packagesMap.M { // download all packeges
		andso, _ := exec.Command("apt", "download", keyPack).Output()
		fmt.Println(string(andso))
	}
}

func RecurseDependens(ma *MainMap, i uint) { //pepreopredeleniya peremennih
	for key := range ma.M {
		depends, err := ListPackageDepends(key)
		if err != nil {
			log.Fatal("Cant check dependens: ", err)
		}
		if depends != nil {
			for _, dep := range depends {
				ma.Lock()
				keyExist := IsKeyExist(dep, ma)
				if !keyExist {
					Put(dep, i, ma)
					i++
				}
				ma.Unlock()
			}
		}
		// fmt.Println("for pack: ", key, " have depends: ", depends)
	}
	fmt.Println("Itog:", ma, ". Total packeges: ", i)
}

func Put(key string, val uint, ma *MainMap) {
	if !(strings.Contains(key, "lib")) && !(strings.Contains(key, "iptables")) && !(strings.Contains(key, "passwd")) && !(strings.Contains(key, "adduser")) { //avoid standart libs
		ma.M[key] = val
	}
}

func IsKeyExist(key string, ma *MainMap) bool { //Check key exist
	_, ok := ma.M[key]
	if !ok {
		return false
	}
	return true
}

func ListPackageDepends(pak string) ([]string, error) { //Get list of depends
	outAptCache, err := exec.Command("apt-cache", "depends", pak).Output()
	if err != nil {
		return nil, err
	}
	stdoutresult := string(outAptCache)
	arstr := strings.Split(stdoutresult, "\n")

	var dependsArrStr []string
	for _, da := range arstr {
		if !(strings.Contains(da, "<") || strings.Contains(da, ">")) {
			if strings.Contains(da, "PreDepends") {
				da = string(da[13:])
				dependsArrStr = append(dependsArrStr, da)
			} else if strings.Contains(da, "Depends") {
				da = string(da[11:])
				dependsArrStr = append(dependsArrStr, da)
			}

		}
	}
	return dependsArrStr, nil
}
