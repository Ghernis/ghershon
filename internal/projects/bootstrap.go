package bootstrap
//package main

import (
	"ghershon/internal/projects/templates"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
	"log"
	"strings"
)

func write_file(content string, filepath string){
	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
        log.Fatalf("failed to write file: %v", err)
    }

    log.Printf("✅ %v created\n",filepath)
}
func write_template(content []byte, filepath string){
	if err := os.WriteFile(filepath, content, 0644); err != nil {
        log.Fatalf("failed to write file: %v", err)
    }

    log.Printf("✅ %v created\n",filepath)
}

func Python_boot(basepath string,name string) {
	pyMainTemplate, err := templates.FS.ReadFile("main.py.tmpl")
	if err != nil{
		fmt.Errorf("Loading template: %w",err)
	}
	adoPipeTemplate, err := templates.FS.ReadFile("pipelines.yml.tmpl")
	if err != nil{
		fmt.Errorf("Loading template: %w",err)
	}
	content_md :=`# Titulo App
`
	// TODO tener el anio en config object, tambien si es de company o privado
    content_utils := "from .utils import util"
    content_internal := "from .internal import modelo"
    content_gitignore := `.env
venv/**
__pycache__/**
utils/__pycache__/**
internal/__pycache__/**
`
	home, err := os.UserHomeDir()
	if err != nil{
		log.Fatal("could not resolve home directory: %v",err)
	}
	//target := filepath.Join(home,"kode","proyects","test",name)
	rutas:=strings.Split(basepath,"/")
	target := home
	for _,v := range rutas{
		target = filepath.Join(target,v)
	}
	target = filepath.Join(target,name)

    // 1. Create folders
    dirs := []string{
        target,
        filepath.Join(target, "utils"),
        filepath.Join(target, "pipelines"),
        filepath.Join(target, "internal"),
    }
    for _, d := range dirs {
        if err := os.MkdirAll(d, 0755); err != nil {
            fmt.Fprintf(os.Stderr, "❌ cannot create %s: %v\n", d, err)
            os.Exit(1)
        }
    }
    fmt.Println(" Project structure created.")
	// 1.5 Create files
	write_template(pyMainTemplate,filepath.Join(target,"main.py"))
	write_template(adoPipeTemplate,filepath.Join(target,"pipelines","run.yml"))
	write_file(content_utils,filepath.Join(target,"utils","__init__.py"))
	write_file("",filepath.Join(target,"utils","utils.py"))
	write_file("",filepath.Join(target,"internal","internal.py"))
	write_file(content_internal,filepath.Join(target,"internal","__init__.py"))
	//write_file(content_pipeline,filepath.Join(target,"pipelines","run.yml"))
	write_file(content_gitignore,filepath.Join(target,".gitignore"))
	write_file("",filepath.Join(target,".env"))
	write_file("",filepath.Join(target,"requirements.txt"))
	write_file(content_md,filepath.Join(target,"readme.md"))
    fmt.Println(" Individual files created.")

    // 2. Create a Python venv
    pyCmd := "python"
    if runtime.GOOS == "windows" {
        pyCmd = "python" // typically installed as python.exe
    }
    venvPath := filepath.Join(target, "venv")
    cmd := exec.Command(pyCmd, "-m", "venv", venvPath)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    fmt.Printf(" Initializing venv with %q...\n", pyCmd)
    if err := cmd.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "❌ venv setup failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("✅ Python venv ready at", venvPath)
}
