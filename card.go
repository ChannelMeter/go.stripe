package stripe

// Credit Card Types accepted by the Stripe API.
const (
	AmericanExpress = "American Express"
	DinersClub      = "Diners Club"
	Discover        = "Discover"
	JCB             = "JCB"
	MasterCard      = "MasterCard"
	Visa            = "Visa"
	UnknownCard     = "Unknown"
)

type Card struct {
    Id                string `json:"id"`
    Name              string `json:"name"`                // Cardholder name
    Type              string `json:"type"`                // Card brand. Can be Visa, American Express, MasterCard, Discover, JCB, Diners Club, or Unknown
    ExpMonth          int    `json:"exp_month"`
    ExpYear           int    `json:"exp_year"`
    Last4             int    `json:"last4"`
    Fingerprint       string `json:"fingerprint"`         // Uniquely identifies this particular card number. You can use this attribute to check whether two customers who've signed up with you are using the same card number
    Country           string `json:"country"`             // Two-letter ISO code representing the country of the card (as accurately as we can determine it). You could use this attribute to get a sense of the international breakdown of cards you've collected.
    Address1          string `json:"address_line1"`
	Address2          string `json:"address_line2"`
    AddressCountry    string `json:"address_country"`     // Billing address country, if provided when creating card
	AddressState      string `json:"address_state"`
    AddressZip        string `json:"address_zip"`
    AddressLine1Check string `json:"address_line1_check"` // If address_line1 was provided, results of the check: pass, fail, or unchecked
	AddressZipCheck   string `json:"address_zip_check"`   // If address_zip was provided, results of the check: pass, fail, or unchecked 
    CVCCheck          string `json:"cvc_check"`           // If a CVC was provided, results of the check: pass, fail, or unchecked
}

// TODO handle A common source of error is an invalid or expired card, or a valid card with insufficient available balance.
func (self *Card) IsExpired() bool {
	return false
}


// LuhnValid uses the Luhn Algorithm (also known as the Mod 10 algorithm) to
// verify a credit cards checksum, which helps flag accidental data entry
// errors.
//
// see http://en.wikipedia.org/wiki/Luhn_algorithm
func LuhnValid(card string) (bool, error) {

	var sum = 0
	var digits = strings.Split(card, "")

	// iterate through the digits in reverse order
	for i, even :=len(digits)-1, false; i>=0; i, even = i-1, !even {

		// convert the digit to an integer
		digit, err := strconv.Atoi(digits[i])
		if err != nil {
			return false, err
		}

		// we multiply every other digit by 2, adding the product to the sum.
		// note: if the product is double digits (i.e. 14) we add the two digits
		//       to the sum (14 -> 1+4 = 5). A simple shortcut is to subtract 9
		//       from a double digit product (14 -> 14 - 9 = 5).
		switch {
		case  even && digit > 4 : sum += (digit * 2) - 9
		case  even : sum += digit * 2
		case !even : sum += digit
		}
	}

	// if the sum is divisible by 10, it passes the check
	return sum % 10 == 0, nil
}

// CardType is a simple algorithm to determine the Card Type (ie Visa, Discover)
// based on the Credit Card Number. If the Number is not recognized, a value
// of "Unknown" will be returned.
func CardType(card string) string {
	
	switch card[0:1] {
	case "4" : return Visa
	case "2", "1" :
		switch card[0:4] {
		case "2131", "1800" : return JCB
		}
	case "6" : 
		switch card[0:4] {
		case "6011" : return Discover
		}
	case "5" :
		switch card[0:2] {
		case "51", "52", "53", "54", "55" : return MasterCard
		}
	case "3" :
		switch card[0:2] {
		case "34", "37" : return AmericanExpress
		case "36" : return DinersClub
		case "30" :
			switch card[0:3] {
			case "300", "301", "302", "303", "304", "305" : return DinersClub
			}
		default : return JCB
		}
	}

	return UnknownCard
}