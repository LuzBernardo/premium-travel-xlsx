package main

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

var (
	data string = `Fairfield Inn & Suites Los Angeles LAX/El Segundo

FLEXIBLE RATE GUEST ROOM 1 KING.

Thu, May 9, 2024 – Mon, May 13, 2024
2 travelers | 1 room | 4 nights

Room description
1 King bed
Mini fridge
255 square feet/23 square meters
Wireless internet
Complimentary coffee/tea maker
Maximum occupancy of 2 guests

Rate details
A00Rega: Flexible Rate Guest Room 1 King

Cancellation policy
Refundable until May 7, 2024
* 178.08 usd cancel fee person room cancellation permitted up to 1 days before arrival.

Nightly rate
Night 1 (2024-05-09)     $159.00
Night 2 (2024-05-10)     $174.00
Night 3 (2024-05-11)     $164.00
Night 4 (2024-05-12)     $159.00

Price details
1 Room x 4 Nights     $656.00
Taxes and Fees        $80.00
Total USD             $736.00
Including all known taxes and fees`

	regexPrices          = regexp.MustCompile(`(\d+) Nights.*\$(\d+\.?\d*)`)
	regexAfterCifraValue = regexp.MustCompile(`\$(\d+\.?\d*)`)
)

func Create() {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Data")
	if err != nil {
		fmt.Printf("Failed to create sheet: %s\n", err)
		return
	}

	headerStyle := xlsx.NewStyle()
	headerStyle.Font = *xlsx.NewFont(12, "Arial")
	headerStyle.Font.Bold = true
	headerStyle.Fill = *xlsx.NewFill("solid", "D9D9D9", "FF000000")
	headerStyle.Alignment.Horizontal = "left"
	headerStyle.ApplyFont = true
	headerStyle.ApplyFill = true
	headerStyle.ApplyBorder = true
	headerStyle.ApplyAlignment = true

	normalStyle := xlsx.NewStyle()
	normalStyle.Font = *xlsx.NewFont(12, "Arial")
	normalStyle.Border = *xlsx.NewBorder("thin", "thin", "thin", "thin")
	normalStyle.Alignment.Horizontal = "center"
	normalStyle.ApplyFont = true
	normalStyle.ApplyBorder = true
	normalStyle.ApplyAlignment = true

	totalStyle := xlsx.NewStyle()
	totalStyle.Font = *xlsx.NewFont(12, "Arial")
	totalStyle.Fill = *xlsx.NewFill("solid", "FFC000", "FF000000")
	totalStyle.Border = *xlsx.NewBorder("thin", "thin", "thin", "thin")
	totalStyle.Alignment.Horizontal = "left"
	totalStyle.Font.Bold = true
	totalStyle.ApplyFont = true
	totalStyle.ApplyFill = true
	totalStyle.ApplyBorder = true
	totalStyle.ApplyAlignment = true

	dataInsert := make([][]string, 0)

	scanner := bufio.NewScanner(strings.NewReader(data))

	in := make([]string, 0)
	for scanner.Scan() {
		text := scanner.Text()

		if text == "" {
			dataInsert = append(dataInsert, in)
			in = make([]string, 0)
			continue
		}

		in = append(in, text)
	}

	dataInsert = append(dataInsert, in)

	mxLines := 0

	for _, column := range dataInsert {
		if len(column) > mxLines {
			mxLines = len(column)
		}
	}

	for i := 1; i <= mxLines; i++ {
		sheet.AddRow()
	}

	xInit := 2
	yInit := 6
	add := 1

	avgRate, err := avgRate(dataInsert[len(dataInsert)-1][1])
	if err != nil {
		fmt.Printf("Failed to create get avgRate: %s\n", err)
		return
	}

	priceDetails := getSection(dataInsert, "Price details")
	if priceDetails == nil {
		fmt.Printf("[priceDetails == nil] There is no price details!!!")
		return
	}

	rateDetails := getSection(dataInsert, "Rate details")
	if priceDetails == nil {
		fmt.Printf("[priceDetails == nil] There is no rate details!!!")
		return
	}

	cancellationPolicy := getSection(dataInsert, "Cancellation policy")
	if priceDetails == nil {
		fmt.Printf("[priceDetails == nil] There is no cancellation policy!!!")
		return
	}

	totalRateRoom := getRate(priceDetails[1])
	taxesAndFees := getRate(priceDetails[2])
	totalUSD := getRate(priceDetails[3])

	SetValueX(sheet, yInit, xInit, "HOTEL", dataInsert[0][0])
	SetValueX(sheet, yInit+add, xInit, "RATE", dataInsert[1][0])
	add++
	SetValueX(sheet, yInit+add, xInit, "CANCELLATION POLICY")
	add++
	SetValueX(sheet, yInit+add, xInit, "ROOM CATEGORY")
	add++
	SetValueX(sheet, yInit+add, xInit, "AVG. NIGHTLY RATE", avgRate)
	add++
	SetValueX(sheet, yInit+add, xInit, "TOTAL ROOM RATE", totalRateRoom)
	add++
	SetValueX(sheet, yInit+add, xInit, "TAXES AND FEES", taxesAndFees)
	add++
	SetValueX(sheet, yInit+add, xInit, "TOTAL USD", totalUSD)
	add++
	SetValueX(sheet, yInit+add, xInit, "Including all known taxes and fees")

	SetValueY(sheet, yInit+add+3, xInit, cancellationPolicy[0:]...)
	SetValueY(sheet, yInit+add+3, xInit+1, rateDetails[0:]...)

	for i := yInit; i < yInit+add-1; i++ {
		style := *headerStyle
		style.Border = *xlsx.NewBorder("medium", "thin", "thin", "thin")
		if i == yInit {
			style.Border = *xlsx.NewBorder("medium", "thin", "medium", "thin")
		}

		cell := sheet.Cell(i, xInit)
		cell.SetStyle(&style)

		styleB := *normalStyle
		styleB.Border = *xlsx.NewBorder("thin", "medium", "thin", "thin")
		if i == yInit {
			styleB.Border = *xlsx.NewBorder("thin", "medium", "medium", "thin")
		}

		cellB := sheet.Cell(i, xInit+1)
		cellB.SetStyle(&styleB)
	}

	for j := xInit; j < xInit+2; j++ {
		style := *totalStyle
		style.Border = *xlsx.NewBorder("medium", "thin", "thin", "medium")
		if j == xInit+1 {
			style.Border = *xlsx.NewBorder("thin", "medium", "thin", "medium")
			style.Alignment.Horizontal = "center"
		}

		cell := sheet.Cell(yInit+add-1, j)
		cell.SetStyle(&style)
	}

	if err = file.Save("structured_output.xlsx"); err != nil {
		fmt.Printf("Failed to save file: %s\n", err)
		return
	}

	fmt.Println("Excel file created successfully with structured data")
}

func SetValueX(sheet *xlsx.Sheet, yInit int, xInit int, values ...string) {
	cellA := sheet.Cell(yInit, xInit)
	cellA.SetValue(values[0])

	for i := 1; i < len(values); i++ {
		cellA1 := sheet.Cell(yInit, xInit+i)
		cellA1.SetValue(values[i])
	}
}

func SetValueY(sheet *xlsx.Sheet, yInit int, xInit int, values ...string) {
	cellA := sheet.Cell(yInit, xInit)
	cellA.SetValue(values[0])

	for i := 1; i < len(values); i++ {
		cellA1 := sheet.Cell(yInit+i, xInit)
		cellA1.SetValue(values[i])
	}
}

func avgRate(s string) (string, error) {
	// Regular expression to find numbers and dollar amounts
	matches := regexPrices.FindStringSubmatch(s)

	if len(matches) < 3 {
		fmt.Println("Failed to extract data")
		return "", errors.New("[avgRate] len(matches) < 3 is not true")
	}

	// Extracting number of nights
	nights, err := strconv.Atoi(matches[1])
	if err != nil {
		fmt.Printf("Error converting nights: %s\n", err)
		return "", err
	}

	// Extracting total amount
	amountStr := strings.Replace(matches[2], ",", "", -1)
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Printf("Error converting amount: %s\n", err)
		return "", err
	}

	// Calculating price per night
	pricePerNight := amount / float64(nights)

	return fmt.Sprintf("$ %.2f", pricePerNight), nil
}

func getRate(s string) string {
	matches := regexAfterCifraValue.FindStringSubmatch(s)

	return fmt.Sprintf("$ %s", matches[1])
}

func getSection(data [][]string, s string) []string {
	regex := regexp.MustCompile(fmt.Sprintf("\\b%s\\b", s))

	for _, d := range data {
		if regex.MatchString(d[0]) {
			return d
		}
	}

	return nil
}
