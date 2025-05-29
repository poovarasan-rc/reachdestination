package destination

import (
	"appsec/db"
	"appsec/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ResponseStruct struct {
	Resp []model.Destination `json:"data"`
	model.Debug
}

func GetAccDestination(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	log.Println("GetAccDestination(+)")

	var lRespRec ResponseStruct

	lRespRec.Sts = "S"
	lRespRec.Msg = ""

	var lErr error

	// Validate the request
	lErr = validateReq(r)
	if lErr != nil {
		lRespRec.Sts = "E"
		lRespRec.Msg = "DGAD01 : " + lErr.Error()
		log.Println(lRespRec)
		goto Marshal
	}
	// Retrieve the destinations of the account
	lRespRec.Resp, lErr = fetchDestination(r.Header.Get("ACC_ID"))
	if lErr != nil {
		lRespRec.Sts = "E"
		lRespRec.Msg = "DGAD02 : " + lErr.Error()
		log.Println(lRespRec)
		goto Marshal
	}

Marshal:
	data, err := json.Marshal(lRespRec)
	if err != nil {
		fmt.Fprintf(w, "Error taking data"+err.Error())
	} else {
		fmt.Fprintf(w, string(data))
	}
	log.Println("GetAccDestination(-)")

}

// Validate the request
func validateReq(r *http.Request) error {

	if r.Header.Get("ACC_ID") == "" {
		return fmt.Errorf("Missing ACC_ID in header or empty")
	}

	if r.Method != "GET" {
		return fmt.Errorf("GET method only allowed")
	}

	return nil
}

// Retrieve the destinations of the account
func fetchDestination(pReq string) ([]model.Destination, error) {
	var lArr []model.Destination

	query := `
        SELECT d.Desid,d.url,d.method
		from account_holders ah,destination d 
		where ah.id = d.accid
		and ah.accid = ?;
    `
	rows, err := db.GDBCon.Query(query, pReq)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var lRec model.Destination

		err := rows.Scan(&lRec.DesID, &lRec.URL, &lRec.HTTPMethod)
		if err != nil {
			return nil, err
		}
		lArr = append(lArr, lRec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(lArr) == 0 {
		return nil, fmt.Errorf("Destinations are not available for the account - " + pReq)
	}

	return lArr, nil
}
