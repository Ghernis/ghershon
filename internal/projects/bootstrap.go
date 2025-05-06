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

func sanitize_name(name string) string{
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.Replace(name,"-","_",-1)
	name = strings.Replace(name," ","",-1)
	name = strings.Replace(name,",","_",-1)
	return name

}
func Django_boot(basepath string,name string){
	// 0. Home
	home, err := os.UserHomeDir()
	if err != nil{
		log.Fatal("could not resolve home directory: %v",err)
	}
	rutas:=strings.Split(basepath,"/")
	target := home
	for _,v := range rutas{
		target = filepath.Join(target,v)
	}
	target = filepath.Join(target,name)
    // 1. Create folders
    dirs := []string{
        target,
        filepath.Join(target, "app"),
        filepath.Join(target, "app","openshift"),
        filepath.Join(target, "app","ci"),
        filepath.Join(target, "cd"),
        filepath.Join(target, "cd","templates"),
    }
    for _, d := range dirs {
        if err := os.MkdirAll(d, 0755); err != nil {
            fmt.Fprintf(os.Stderr, "❌ cannot create %s: %v\n", d, err)
            os.Exit(1)
        }
    }
    fmt.Println(" Project structure created.")

    // 2. Create a Python venv
    pyCmd := "python"
	venvActivatePath:= "bin"
    if runtime.GOOS == "windows" {
        pyCmd = "python" // typically installed as python.exe
		venvActivatePath = "Scritps"
    }
	djangoPath := filepath.Join(target,"app")
    venvPath := filepath.Join(djangoPath, "venv")
    cmd := exec.Command(pyCmd, "-m", "venv", venvPath)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    fmt.Printf(" Initializing venv with %q...\n", pyCmd)
    if err := cmd.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "❌ venv setup failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("✅ Python venv ready at", venvPath)
	// 3. Install python libs
	pythonVenv := filepath.Join(venvPath,venvActivatePath,"python")
	pipVenv := filepath.Join(venvPath,venvActivatePath,"pip")
	fmt.Println(pythonVenv)

	cmd2 := exec.Command(pipVenv,"install","django","gunicorn")
    cmd2.Stdout = os.Stdout
    cmd2.Stderr = os.Stderr

    fmt.Printf(" Installing libraries venv with %q...\n", pyCmd)
    if err := cmd2.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "❌ libs setup failed: %v\n", err)
        os.Exit(1)
    }
	cmd3 := exec.Command(pipVenv,"list")
    cmd3.Stdout = os.Stdout
    cmd3.Stderr = os.Stderr

    fmt.Printf(" Listing libraries in venv %q...\n", pyCmd)
    if err := cmd3.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "❌ libs listing failed: %v\n", err)
        os.Exit(1)
    }

	//django-admin startproject <service_name> .

	djangoVenv := filepath.Join(venvPath,venvActivatePath,"django-admin")
	cmd4 := exec.Command(djangoVenv,"startproject", name, ".")
	cmd4.Dir = djangoPath
    cmd4.Stdout = os.Stdout
    cmd4.Stderr = os.Stderr

    fmt.Printf(" Starting django project %q...\n", pyCmd)
    if err := cmd4.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "❌ libs setup failed: %v\n", err)
        os.Exit(1)
    }

	// CI Pipeline
	
	ciPipe, err := templates.FS.ReadFile("ado_oc_build.yml.tmpl")
	if err != nil{
		fmt.Errorf("Loading template: %w",err)
	}
	write_template(ciPipe,filepath.Join(djangoPath,"ci","build.yml"))
	// gitignore

	gitignore, err := templates.FS.ReadFile("gitignore.tmpl")
	if err != nil{
		fmt.Errorf("Loading template: %w",err)
	}
	write_template(gitignore,filepath.Join(djangoPath,".gitignore"))

	// buildConfig template

	buildConfig, err := templates.ParseTemplate("django_buildconfig.yml.tmpl", map[string]string{
		"ProjectName": name,
	})
	if err != nil{
		fmt.Errorf("Template generation failed: %v",err)
	}
	write_file(buildConfig,filepath.Join(djangoPath,"openshift","buildconfig.yml"))

	// helm chart
	cdPath := filepath.Join(target,"cd")
	helmChart, err := templates.ParseTemplate("chart.yml.tmpl", map[string]string{
		"ServiceName": name,
	})
	if err != nil{
		fmt.Errorf("Chart generation failed: %v",err)
	}
	write_file(helmChart,filepath.Join(cdPath,"Chart.yml"))

	valuesChart, err := templates.ParseTemplate("chart_values.yml.tmpl", map[string]string{
		"ServiceName": name,
	})
	if err != nil{
		fmt.Errorf("Chart values generation failed: %v",err)
	}
	write_file(valuesChart,filepath.Join(cdPath,"values.yml"))

	deploymentChart, err := templates.ParseTemplate("chart_deployment.yml.tmpl", map[string]string{
		"ServiceName": name,
	})
	if err != nil{
		fmt.Errorf("Chart deployment generation failed: %v",err)
	}
	write_file(deploymentChart,filepath.Join(cdPath,"templates","deployment.yml"))

	serviceChart, err := templates.ParseTemplate("chart_service.yml.tmpl", map[string]string{
		"ServiceName": name,
	})
	if err != nil{
		fmt.Errorf("Chart service generation failed: %v",err)
	}
	write_file(serviceChart,filepath.Join(cdPath,"templates","service.yml"))

	routeChart, err := templates.ParseTemplate("chart_route.yml.tmpl", map[string]string{
		"ServiceName": name,
	})
	if err != nil{
		fmt.Errorf("Chart route generation failed: %v",err)
	}
	write_file(routeChart,filepath.Join(cdPath,"templates","route.yml"))
	

	// Git Init

	cmd5 := exec.Command("git","init")
	cmd5.Dir = djangoPath
    cmd5.Stdout = os.Stdout
    cmd5.Stderr = os.Stderr
    fmt.Printf(" Initializing git with %q...\n", "git")
    if err := cmd5.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "❌ git init failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("✅ Git ready at", target)
	

}

func Python_boot(basepath string,name string) {
	//pyMainTemplate, err := templates.FS.ReadFile("main.py.tmpl")
	//if err != nil{
	//	fmt.Errorf("Loading template: %w",err)
	//}
	pyMainTemplate, err := templates.ParseTemplate("main.py.tmpl", map[string]string{
		"Author": "\"Hernan Gomez\"",
		"Year": "2025",
		"Mail": "\"hernan.gomez@set.ypf.com\"",
	})
	if err != nil{
		fmt.Errorf("Template generation failed: %v",err)
	}
	readmeTemplate, err := templates.FS.ReadFile("readme.md.tmpl")
	if err != nil{
		fmt.Errorf("Loading template: %w",err)
	}
	adoPipeTemplate, err := templates.FS.ReadFile("pipelines.yml.tmpl")
	if err != nil{
		fmt.Errorf("Loading template: %w",err)
	}
	gitIgnoreTemplate, err := templates.FS.ReadFile("gitignore.tmpl")
	if err != nil{
		fmt.Errorf("Loading template: %w",err)
	}
	// TODO tener el anio en config object, tambien si es de company o privado
    content_utils := "from .utils import util"
    content_internal := "from .internal import modelo"
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

	write_file(pyMainTemplate,filepath.Join(target,"main.py"))
	write_template(adoPipeTemplate,filepath.Join(target,"pipelines","run.yml"))
	write_template(gitIgnoreTemplate,filepath.Join(target,".gitignore"))
	write_template(readmeTemplate,filepath.Join(target,"readme.md"))
	write_file(content_utils,filepath.Join(target,"utils","__init__.py"))
	write_file("",filepath.Join(target,"utils","utils.py"))
	write_file("",filepath.Join(target,"internal","internal.py"))
	write_file(content_internal,filepath.Join(target,"internal","__init__.py"))
	write_file("",filepath.Join(target,".env"))
	write_file("",filepath.Join(target,"requirements.txt"))

    fmt.Println(" Individual files created.")

    // 2. Create a Python venv
    pyCmd := "python"
	venvActivatePath:= "bin"
    if runtime.GOOS == "windows" {
        pyCmd = "python" // typically installed as python.exe
		venvActivatePath = "Scritps"
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

	pythonPath := filepath.Join(target, "venv", venvActivatePath, "python")
	pipPath := filepath.Join(target, "venv", venvActivatePath, "pip")
	fmt.Println(pythonPath)
	cmd3 := exec.Command(pipPath,"list")
	//cmd3.Dir=target
    cmd3.Stdout = os.Stdout
    cmd3.Stderr = os.Stderr
    if err := cmd3.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "❌ fail to enter venv: %v\n", err)
        os.Exit(1)
    }

	cmd2 := exec.Command("git","init")
	cmd2.Dir = target
    cmd2.Stdout = os.Stdout
    cmd2.Stderr = os.Stderr
    fmt.Printf(" Initializing git with %q...\n", "git")
    if err := cmd2.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "❌ git init failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("✅ Git ready at", target)
}
