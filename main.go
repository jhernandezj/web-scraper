package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/xuri/excelize/v2"
)

func main() {
	// Crear archivo Excel
	excelFile := excelize.NewFile()
	sheetName := "Awardees"
	excelFile.NewSheet(sheetName)
	index, err := excelFile.GetSheetIndex(sheetName)
	if err != nil {
		log.Fatalf("Error al obtener el índice de la hoja: %v\n", err)
	}
	excelFile.SetActiveSheet(index)

	// Escribir encabezados
	headers := []string{"#", "First Name", "Surname", "Gender", "University", "Programme", "Origin"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		excelFile.SetCellValue(sheetName, cell, header)
		// Ajustar el ancho de la columna
		excelFile.SetColWidth(sheetName, string(rune('A'+i)), string(rune('A'+i)), 30)
	}

	// Crear contexto de Chrome con timeout más largo
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // Para ver el navegador mientras trabaja
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-background-networking", false),
		chromedp.Flag("enable-logging", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// Establecer timeout más largo (30 minutos)
	ctx, cancel = context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	// URL base
	baseURL := "https://goal.edu.gy/goal-awardees/"
	row := 2 // Comenzamos desde la fila 2 después de los encabezados

	// Navegar a la página inicial
	if err := chromedp.Run(ctx,
		chromedp.Navigate(baseURL),
		chromedp.Sleep(5*time.Second),
		chromedp.WaitVisible("table.ninja_footable", chromedp.ByQuery),
	); err != nil {
		log.Fatalf("Error al navegar a la página inicial: %v", err)
	}

	// Procesar cada página
	for pageNum := 1; pageNum <= 156; pageNum++ {
		fmt.Printf("Procesando página %d\n", pageNum)

		// Extraer datos de la página actual
		var rows [][]string
		err := chromedp.Run(ctx,
			chromedp.Sleep(2*time.Second), // Esperar a que se carguen los datos
			chromedp.Evaluate(`
				Array.from(document.querySelectorAll("tr[class^='ninja_table_row_']")).map(row => {
					return Array.from(row.querySelectorAll("td")).map(td => td.textContent.trim());
				})
			`, &rows),
		)
		if err != nil {
			log.Printf("Error al extraer datos de la página %d: %v\n", pageNum, err)
			continue
		}

		// Escribir datos en Excel
		for _, fields := range rows {
			if len(fields) == 7 { // Asegurarse de que tengamos todas las columnas
				// Verificar que no sea una fila vacía
				hasData := false
				for _, field := range fields {
					if field != "" {
						hasData = true
						break
					}
				}

				if hasData {
					for i, field := range fields {
						cell := fmt.Sprintf("%s%d", string(rune('A'+i)), row)
						excelFile.SetCellValue(sheetName, cell, field)
					}
					fmt.Printf("Procesando registro #%s: %s %s\n", fields[0], fields[1], fields[2])
					row++
				}
			}
		}

		// Guardar periódicamente
		if pageNum%10 == 0 {
			tempFileName := fmt.Sprintf("goal_awardees_temp_%d.xlsx", pageNum)
			if err := excelFile.SaveAs(tempFileName); err != nil {
				log.Printf("Error al guardar archivo temporal en página %d: %v\n", pageNum, err)
			} else {
				fmt.Printf("Guardado temporal en página %d: %s\n", pageNum, tempFileName)
			}
		}

		// Si no es la última página, hacer clic en el siguiente botón
		if pageNum < 156 {
			// Hacer clic en el siguiente número de página
			nextPage := pageNum + 1
			err = chromedp.Run(ctx,
				chromedp.Sleep(1*time.Second),
				chromedp.Evaluate(fmt.Sprintf(`
					(() => {
						const pageLink = document.querySelector('li.footable-page[data-page="%d"] a.footable-page-link');
						if (pageLink) {
							pageLink.click();
							return true;
						}
						return false;
					})()
				`, nextPage), nil),
				chromedp.Sleep(2*time.Second), // Esperar a que se actualice la tabla
			)
			if err != nil {
				log.Printf("Error al navegar a la siguiente página: %v\n", err)
				continue
			}
		}
	}

	// Guardar Excel final
	excelFileName := "goal_awardees.xlsx"
	if err := excelFile.SaveAs(excelFileName); err != nil {
		log.Fatalf("Error al guardar Excel: %v\n", err)
	}

	fmt.Printf("Archivo Excel generado exitosamente: %s\n", excelFileName)
}
