package openrtb

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Validation errors
var (
	ErrInvalidBidNoID    = errors.New("openrtb: bid is missing ID")
	ErrInvalidBidNoImpID = errors.New("openrtb: bid is missing impression ID")
)

type MultiString string

func (s *MultiString) UnmarshalJSON(data []byte) error {
	var value interface{}

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch value.(type) {
	case string:
		*s = MultiString(value.(string))
	case float64:
		*s = MultiString(strconv.FormatFloat(value.(float64), 'f', -1, 64))
	default:
		return errors.New("unknown type: " + reflect.TypeOf(value).String())
	}

	return nil
}

func (s *MultiString) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", s)), nil
}

func (s MultiString) String() string {
	return string(s)
}

// ID, ImpID and Price are required; all other optional.
// If the bidder wins the impression, the exchange calls notice URL (nurl)
// a) to inform the bidder of the win;
// b) to convey certain information using substitution macros.
// Adomain can be used to check advertiser block list compliance.
// Cid can be used to block ads that were previously identified as inappropriate.
// Substitution macros may allow a bidder to use a static notice URL for all of its bids.
type Bid struct {
	ID             string      `json:"id"`
	ImpID          string      `json:"impid"`                    // Required string ID of the impression object to which this bid applies.
	Price          float64     `json:"price"`                    // Bid price in CPM. Suggests using integer math for accounting to avoid rounding errors.
	AdID           string      `json:"adid,omitempty"`           // References the ad to be served if the bid wins.
	NURL           string      `json:"nurl,omitempty"`           // Win notice URL.
	AdMarkup       string      `json:"adm,omitempty"`            // Actual ad markup. XHTML if a response to a banner object, or VAST XML if a response to a video object.
	AdvDomain      []string    `json:"adomain,omitempty"`        // Advertiser’s primary or top-level domain for advertiser checking; or multiple if imp rotating.
	Bundle         string      `json:"bundle,omitempty"`         // A platform-specific application identifier intended to be unique to the app and independent of the exchange.
	IURL           string      `json:"iurl,omitempty"`           // Sample image URL.
	CampaignID     MultiString `json:"cid,omitempty"`            // Campaign ID that appears with the Ad markup.
	CreativeID     string      `json:"crid,omitempty"`           // Creative ID for reporting content issues or defects. This could also be used as a reference to a creative ID that is posted with an exchange.
	Cat            []string    `json:"cat,omitempty"`            // IAB content categories of the creative. Refer to List 5.1
	Attr           []int       `json:"attr,omitempty"`           // Array of creative attributes.
	API            int         `json:"api,omitempty"`            // API required by the markup if applicable
	Protocol       int         `json:"protocol,omitempty"`       // Video response protocol of the markup if applicable
	QAGMediaRating int         `json:"qagmediarating,omitempty"` // Creative media rating per IQG guidelines.
	DealID         string      `json:"dealid,omitempty"`         // DealID extension of private marketplace deals
	H              int         `json:"h,omitempty"`              // Height of the ad in pixels.
	W              int         `json:"w,omitempty"`              // Width of the ad in pixels.
	Exp            int         `json:"exp,omitempty"`            // Advisory as to the number of seconds the bidder is willing to wait between the auction and the actual impression.
	Ext            Extension   `json:"ext,omitempty"`
}

// Validate required attributes
func (bid *Bid) Validate() error {
	if bid.ID == "" {
		return ErrInvalidBidNoID
	} else if bid.ImpID == "" {
		return ErrInvalidBidNoImpID
	}

	return nil
}
