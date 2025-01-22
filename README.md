# Web Scraper para Goal Awardees

Este es un web scraper desarrollado en Go que extrae información de la tabla de beneficiarios de GOAL (https://goal.edu.gy/goal-awardees/) y la guarda en un archivo Excel.

## Características

- Utiliza chromedp para interactuar con la página web
- Maneja paginación dinámica de la tabla
- Guarda los datos en formato Excel (.xlsx)
- Realiza guardados temporales cada 10 páginas
- Muestra el progreso en tiempo real

## Requisitos

- Go 1.21 o superior
- Chrome/Chromium instalado

## Instalación

1. Clona el repositorio:
```bash
git clone https://github.com/jhernandezj/web-scraper.git
cd web-scraper
```

2. Instala las dependencias:
```bash
go mod tidy
```

## Uso

Para ejecutar el scraper:

```bash
go run main.go
```

El programa generará un archivo `goal_awardees.xlsx` con los datos extraídos.

## Estructura del Proyecto

- `main.go`: Código principal del scraper
- `go.mod`: Definición del módulo y dependencias
- `go.sum`: Checksums de las dependencias

## Notas

- El scraper está configurado para mostrar el navegador mientras trabaja (modo no headless)
- Se recomienda tener una conexión estable a internet durante la ejecución
- Los archivos temporales se guardan con el formato `goal_awardees_temp_[número].xlsx`
