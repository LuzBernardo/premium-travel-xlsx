package main

import (
	fl "text-formatter/src"
)

const data string = `Fairfield Inn & Suites Los Angeles LAX/El Segundo

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

func main() {
	fl.Create(data)
}
