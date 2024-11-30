import json

# Rutas de los archivos de configuración
kernel_path = "/home/utnso/TPSO2C/tp-2024-2c-c_aprueba/kernel/config.json"
memoria_path = "/home/utnso/TPSO2C/tp-2024-2c-c_aprueba/memoria/config.json"
filesystem_path = "/home/utnso/TPSO2C/tp-2024-2c-c_aprueba/filesystem/config.json"
cpu_path = "/home/utnso/TPSO2C/tp-2024-2c-c_aprueba/cpu/config.json"

# IPs y puertos hardcodeados
config = {
    "ipkernel": "127.0.0.1",
    "portkernel": 3061,
    "ipmemoria": "127.0.0.1",
    "portmemoria": 3062,
    "ipfilesystem": "127.0.0.1",
    "portfilesystem": 3063,
    "ipcpu": "127.0.0.1",
    "portcpu": 3064
}

# Función para modificar un archivo de configuración
def modify_config(file_path, updates):
    try:
        # Leer archivo existente
        with open(file_path, "r") as file:
            data = json.load(file)
        
        # Actualizar las claves indicadas
        for key, value in updates.items():
            if key in data:
                data[key] = value
        
        # Guardar los cambios
        with open(file_path, "w") as file:
            json.dump(data, file, indent=4)
        
        print(f"Archivo {file_path} actualizado correctamente.")
    except Exception as e:
        print(f"Error al modificar {file_path}: {e}")

# Función principal
def main():
    # Actualizar configuraciones específicas
    modify_config(memoria_path, {
        "port": config["portmemoria"],
        "ip_kernel": config["ipkernel"],
        "port_kernel": config["portkernel"],
        "ip_cpu": config["ipcpu"],
        "port_cpu": config["portcpu"],
        "ip_filesystem": config["ipfilesystem"],
        "port_filesystem": config["portfilesystem"]
    })
    
    modify_config(kernel_path, {
        "port": config["portkernel"],
        "ip_memory": config["ipmemoria"],
        "port_memory": config["portmemoria"],
        "ip_cpu": config["ipcpu"],
        "port_cpu": config["portcpu"]
    })
    
    modify_config(cpu_path, {
        "port": config["portcpu"],
        "ip_memory": config["ipmemoria"],
        "port_memory": config["portmemoria"],
        "ip_kernel": config["ipkernel"],
        "port_kernel": config["portkernel"]
    })
    
    modify_config(filesystem_path, {
        "port": config["portfilesystem"],
        "ip_memory": config["ipmemoria"],
        "port_memory": config["portmemoria"]
    })

if __name__ == "__main__":
    main()
