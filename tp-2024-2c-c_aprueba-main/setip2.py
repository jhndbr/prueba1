import json

# Rutas de los archivos de configuración
kernel_path = "/home/utnso/TPSO2C/tp-2024-2c-c_aprueba/kernel/config.json"
memoria_path = "/home/utnso/TPSO2C/tp-2024-2c-c_aprueba/memoria/config.json"
filesystem_path = "/home/utnso/TPSO2C/tp-2024-2c-c_aprueba/filesystem/config.json"
cpu_path = "/home/utnso/TPSO2C/tp-2024-2c-c_aprueba/cpu/config.json"

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

# Función para solicitar al usuario los datos de configuración
def request_config():
    print("Ingrese las configuraciones necesarias:")
    return {
        "ipkernel": input("IP del kernel (ej. 127.0.0.1): "),
        "portkernel": int(input("Puerto del kernel (ej. 3061): ")),
        "ipmemoria": input("IP de la memoria (ej. 127.0.0.1): "),
        "portmemoria": int(input("Puerto de la memoria (ej. 3062): ")),
        "ipfilesystem": input("IP del filesystem (ej. 127.0.0.1): "),
        "portfilesystem": int(input("Puerto del filesystem (ej. 3063): ")),
        "ipcpu": input("IP de la CPU (ej. 127.0.0.1): "),
        "portcpu": int(input("Puerto de la CPU (ej. 3064): "))
    }

# Función principal
def main():
    # Solicitar datos al usuario
    config = request_config()

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
