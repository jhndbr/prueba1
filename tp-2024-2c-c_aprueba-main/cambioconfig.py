import json

# Rutas de los archivos de configuración
kernel_path = "/home/utnso/TPSO2C/tp-2024-2c-c_aprueba/kernel/config.json"
memoria_path = "/home/utnso/TPSO2C/tp-2024-2c-c_aprueba/memoria/config.json"
filesystem_path = "/home/utnso/TPSO2C/tp-2024-2c-c_aprueba/filesystem/config.json"

# Configuraciones predefinidas
configurations = {
    "plani_proc": {
        "Kernel": {"sheduler_algorithm": "FIFO", "quantum": 1750},
        "Memoria": {
            "memory_size": 1024,
            "response_delay": 1000,
            "scheme": "FIJAS",
            "search_algorithm": "FIRST",
            "partitions": [32, 32, 32, 32, 32, 32, 32, 32],
        },
        "FileSystem": {
            "block_size": 16,
            "block_count": 1024,
            "block_access_delay": 2500,
        },
    },
    "recurso_mutex": {
        "Kernel": {"sheduler_algorithm": "CMN", "quantum": 3750},
        "Memoria": {
            "memory_size": 1024,
            "response_delay": 1000,
            "scheme": "FIJAS",
            "search_algorithm": "FIRST",
            "partitions": [32, 32, 32, 32, 32, 32, 32, 32],
        },
        "FileSystem": {
            "block_size": 16,
            "block_count": 1024,
            "block_access_delay": 2500,
        },
    },
    "mem_fija_base": {
        "Kernel": {"sheduler_algorithm": "CMN", "quantum": 3750},
        "Memoria": {
            "memory_size": 256,
            "response_delay": 1000,
            "scheme": "FIJAS",
            "search_algorithm": "FIRST",
            "partitions": [32, 16, 64, 128, 16],
        },
        "FileSystem": {
            "block_size": 16,
            "block_count": 1024,
            "block_access_delay": 2500,
        },
    },
    "mem_dinamica_base": {
        "Kernel": {"sheduler_algorithm": "CMN", "quantum": 500},
        "Memoria": {
            "memory_size": 1024,
            "response_delay": 200,
            "scheme": "DINAMICAS",
            "search_algorithm": "BEST",
        },
        "FileSystem": {
            "block_size": 32,
            "block_count": 4096,
            "block_access_delay": 2500,
        },
    },
    "prueba_fs": {
        "Kernel": {"sheduler_algorithm": "CMN", "quantum": 500},
        "Memoria": {
            "memory_size": 2048,
            "response_delay": 200,
            "scheme": "DINAMICAS",
            "search_algorithm": "BEST",
        },
        "FileSystem": {
            "block_size": 32,
            "block_count": 200,
            "block_access_delay": 500,
        },
    },
    "the_emptiness_machine": {
        "Kernel": {"sheduler_algorithm": "CMN", "quantum": 250},
        "Memoria": {
            "memory_size": 8192,
            "response_delay": 100,
            "scheme": "DINAMICAS",
            "search_algorithm": "BEST",
        },
        "FileSystem": {
            "block_size": 64,
            "block_count": 1024,
            "block_access_delay": 250,
        },
    },
}

# Función para modificar valores específicos sin eliminar el resto
def modify_config(file_path, changes):
    try:
        # Leer el archivo actual
        with open(file_path, "r") as file:
            config = json.load(file)
        # Actualizar solo las claves especificadas
        for key, value in changes.items():
            config[key] = value
        # Guardar los cambios
        with open(file_path, "w") as file:
            json.dump(config, file, indent=4)
        print(f"Archivo {file_path} actualizado correctamente.")
    except Exception as e:
        print(f"Error al modificar {file_path}: {e}")

# Función para aplicar una configuración predefinida
def apply_config(config_name):
    if config_name not in configurations:
        print(f"La configuración '{config_name}' no existe.")
        return
    config = configurations[config_name]
    try:
        # Modificar cada archivo con los cambios correspondientes
        modify_config(kernel_path, config["Kernel"])
        modify_config(memoria_path, config["Memoria"])
        modify_config(filesystem_path, config["FileSystem"])
        print(f"Configuración '{config_name}' aplicada correctamente.")
    except Exception as e:
        print(f"Error al aplicar la configuración '{config_name}': {e}")

# Menú interactivo
def main():
    print("Selecciona una configuración para aplicar:")
    print("1. plani_proc")
    print("2. recurso_mutex")
    print("3. mem_fija_base")
    print("4. mem_dinamica_base")
    print("5. prueba_fs")
    print("6. the_emptiness_machine")
    try:
        choice = int(input("Ingresa el número de la configuración: ")) - 1
        config_name = list(configurations.keys())[choice]
        apply_config(config_name)
    except (ValueError, IndexError):
        print("Opción inválida. Por favor, inténtalo de nuevo.")

if __name__ == "__main__":
    main()
